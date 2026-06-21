package admin

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
)

// Repo menangani operasi storage untuk admin.
type Repo struct {
	store          *astra.AdminStore
	supabaseClient *supabase.Client
}

func NewRepo(store *astra.AdminStore, supabaseClient *supabase.Client) *Repo {
	return &Repo{
		store:          store,
		supabaseClient: supabaseClient,
	}
}

func (r *Repo) UseSupabase() bool {
	return r.supabaseClient != nil && r.supabaseClient.DB != nil
}

func (r *Repo) FindByEmail(ctx context.Context, email string) (*authmodel.Admin, error) {
	if r.UseSupabase() {
		var a authmodel.Admin
		var createdAt time.Time
		query := `
			SELECT id, nama, email, password, is_walas, kelas, role, created_at
			FROM akun_admin
			WHERE email = $1
		`
		var isWalas bool
		err := r.supabaseClient.DB.QueryRowContext(ctx, query, email).Scan(&a.ID, &a.Nama, &a.Email, &a.Password, &isWalas, &a.Kelas, &a.Role, &createdAt)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		a.IsWalas = isWalas
		a.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
		return &a, nil
	}
	return r.store.FindByEmail(ctx, email)
}

func (r *Repo) FindByID(ctx context.Context, id string) (*authmodel.Admin, error) {
	if r.UseSupabase() {
		var a authmodel.Admin
		var createdAt time.Time
		query := `
			SELECT id, nama, email, password, is_walas, kelas, role, created_at
			FROM akun_admin
			WHERE id = $1
		`
		var isWalas bool
		err := r.supabaseClient.DB.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Nama, &a.Email, &a.Password, &isWalas, &a.Kelas, &a.Role, &createdAt)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		a.IsWalas = isWalas
		a.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
		return &a, nil
	}
	return nil, fmt.Errorf("FindByID not supported on AstraDB")
}

func (r *Repo) Upsert(ctx context.Context, data *authmodel.Admin) error {
	if r.UseSupabase() {
		query := `
			INSERT INTO akun_admin (id, nama, email, password, is_walas, kelas, role)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (email) DO UPDATE 
			SET nama = EXCLUDED.nama,
			    password = EXCLUDED.password,
			    is_walas = EXCLUDED.is_walas,
			    kelas = EXCLUDED.kelas,
			    role = EXCLUDED.role
		`
		_, err := r.supabaseClient.DB.ExecContext(ctx, query, data.ID, data.Nama, data.Email, data.Password, data.IsWalas, data.Kelas, data.Role)
		return err
	}
	return r.store.Upsert(ctx, data)
}

func (r *Repo) UpdateFields(ctx context.Context, id string, isWalas bool, kelas string, role string) error {
	if r.UseSupabase() {
		tx, err := r.supabaseClient.DB.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		// Update
		_, err = tx.ExecContext(ctx, `
			UPDATE akun_admin
			SET role = $1, is_walas = $2, kelas = $3
			WHERE id = $4
		`, role, isWalas, kelas, id)
		if err != nil {
			return err
		}

		// Assignment Class
		_, err = tx.ExecContext(ctx, `
			UPDATE kelas
			SET walas_id = NULL
			WHERE walas_id = $1
		`, id)
		if err != nil {
			return err
		}

		// Kelas Sync
		if isWalas && kelas != "" {
			classes := strings.Split(kelas, ",")
			for _, cls := range classes {
				clsName := strings.TrimSpace(cls)
				if clsName == "" {
					continue
				}
				_, err = tx.ExecContext(ctx, `
					UPDATE kelas
					SET walas_id = $1
					WHERE nama = $2
				`, id, clsName)
				if err != nil {
					return err
				}
			}
		}

		// Walas Sync
		_, err = tx.ExecContext(ctx, `
			UPDATE akun_admin a
			SET is_walas = (a.role = 'walas' OR a.role = 'guru_wali' OR EXISTS (SELECT 1 FROM kelas k WHERE k.walas_id = a.id)),
			    kelas = COALESCE((SELECT string_agg(k.nama, ', ' ORDER BY k.nama) FROM kelas k WHERE k.walas_id = a.id), '')
		`)
		if err != nil {
			return err
		}

		return tx.Commit()
	}
	return fmt.Errorf("UpdateFields not supported on AstraDB")
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	if r.UseSupabase() {
		query := `DELETE FROM akun_admin WHERE id = $1`
		_, err := r.supabaseClient.DB.ExecContext(ctx, query, id)
		return err
	}
	return r.store.Delete(ctx, id)
}

func (r *Repo) ListPaged(ctx context.Context, role string, pageSize int, pageState string) ([]authmodel.Admin, string, error) {
	if r.UseSupabase() {
		var rows *sql.Rows
		var err error
		if role != "" {
			query := `
				SELECT a.id, a.nama, a.email, a.password, a.is_walas, a.kelas, a.role, a.created_at
				FROM akun_admin a
				WHERE a.role = $1
				ORDER BY a.nama ASC
			`
			rows, err = r.supabaseClient.DB.QueryContext(ctx, query, role)
		} else {
			query := `
				SELECT a.id, a.nama, a.email, a.password, a.is_walas, a.kelas, a.role, a.created_at
				FROM akun_admin a
				ORDER BY a.nama ASC
			`
			rows, err = r.supabaseClient.DB.QueryContext(ctx, query)
		}
		if err != nil {
			return nil, "", err
		}
		defer rows.Close()

		var list []authmodel.Admin
		for rows.Next() {
			var a authmodel.Admin
			var isWalas bool
			var createdAt time.Time
			if err := rows.Scan(&a.ID, &a.Nama, &a.Email, &a.Password, &isWalas, &a.Kelas, &a.Role, &createdAt); err != nil {
				return nil, "", err
			}
			a.IsWalas = isWalas
			a.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
			list = append(list, a)
		}
		return list, "", nil
	}
	return r.store.ListPaged(ctx, role, pageSize, pageState)
}
