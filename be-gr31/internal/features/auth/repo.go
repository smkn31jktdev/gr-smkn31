package auth

import (
	"context"
	"database/sql"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
)

type Repo struct {
	students       *astra.StudentStore
	admins         *astra.AdminStore
	supabaseClient *supabase.Client
}

func NewRepo(_ *astra.Client, students *astra.StudentStore, admins *astra.AdminStore, supabaseClient *supabase.Client) *Repo {
	return &Repo{students: students, admins: admins, supabaseClient: supabaseClient}
}

func (r *Repo) UseSupabase() bool {
	return r.supabaseClient != nil && r.supabaseClient.DB != nil
}

func (r *Repo) FindSiswaByNIS(ctx context.Context, nis string) (*authmodel.Siswa, error) {
	if r.UseSupabase() {
		var s authmodel.Siswa
		var waliID sql.NullString
		var walasNama sql.NullString
		var password string
		var createdAt time.Time
		query := `
			SELECT s.id, s.nis, s.nama, s.kelas, s.wali_id, a.nama AS walas_nama, s.password, s.created_at
			FROM akun_siswa s
			LEFT JOIN akun_admin a ON s.wali_id = a.id
			WHERE s.nis = $1
		`
		err := r.supabaseClient.DB.QueryRowContext(ctx, query, nis).Scan(&s.ID, &s.NIS, &s.Nama, &s.Kelas, &waliID, &walasNama, &password, &createdAt)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		s.Walas = walasNama.String
		s.Password = password
		s.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
		return &s, nil
	}
	return r.students.FindByNISN(ctx, nis)
}

func (r *Repo) FindSiswaByID(ctx context.Context, id string) (*authmodel.Siswa, error) {
	if r.UseSupabase() {
		var s authmodel.Siswa
		var waliID sql.NullString
		var walasNama sql.NullString
		var password string
		var createdAt time.Time
		query := `
			SELECT s.id, s.nis, s.nama, s.kelas, s.wali_id, a.nama AS walas_nama, s.password, s.created_at
			FROM akun_siswa s
			LEFT JOIN akun_admin a ON s.wali_id = a.id
			WHERE s.id = $1
		`
		err := r.supabaseClient.DB.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.NIS, &s.Nama, &s.Kelas, &waliID, &walasNama, &password, &createdAt)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		s.Walas = walasNama.String
		s.Password = password
		s.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
		return &s, nil
	}
	return r.students.FindByID(ctx, id)
}

func (r *Repo) FindAdminByEmail(ctx context.Context, email string) (*authmodel.Admin, error) {
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
		if a.Role == "" {
			if isWalas {
				a.Role = "guru_wali"
			} else {
				a.Role = "admin"
			}
		}
		return &a, nil
	}
	return r.admins.FindByEmail(ctx, email)
}

func (r *Repo) FindAdminByID(ctx context.Context, id string) (*authmodel.Admin, error) {
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
		if a.Role == "" {
			if isWalas {
				a.Role = "guru_wali"
			} else {
				a.Role = "admin"
			}
		}
		return &a, nil
	}
	return r.admins.FindByID(ctx, id)
}

func (r *Repo) UpdateAdminRole(ctx context.Context, id string, role string) error {
	if r.UseSupabase() {
		query := `UPDATE akun_admin SET role = $1 WHERE id = $2`
		_, err := r.supabaseClient.DB.ExecContext(ctx, query, role, id)
		return err
	}
	return nil
}
