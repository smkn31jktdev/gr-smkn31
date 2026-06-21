package bukti

import (
	"context"
	"database/sql"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	sekolahmodel "be-gr31/internal/model/sekolah"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
)

// Repo bukti
type Repo struct {
	store          *astra.BuktiStore
	studentStore   *astra.StudentStore
	adminStore     *astra.AdminStore
	supabaseClient *supabase.Client
}

func NewRepo(store *astra.BuktiStore, studentStore *astra.StudentStore, adminStore *astra.AdminStore, supabaseClient *supabase.Client) *Repo {
	return &Repo{
		store:          store,
		studentStore:   studentStore,
		adminStore:     adminStore,
		supabaseClient: supabaseClient,
	}
}

func (r *Repo) Upsert(ctx context.Context, data *sekolahmodel.Bukti) error {
	return r.store.Upsert(ctx, data)
}

func (r *Repo) ListPaged(ctx context.Context, filter sekolahmodel.BuktiFilter, pageSize int, pageState string) ([]sekolahmodel.Bukti, string, error) {
	return r.store.ListPaged(ctx, filter, pageSize, pageState)
}

func (r *Repo) FindAdminByID(ctx context.Context, id string) (*authmodel.Admin, error) {
	if r.supabaseClient != nil && r.supabaseClient.DB != nil {
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
				a.Role = "walas"
			} else {
				a.Role = "admin"
			}
		}
		return &a, nil
	}
	return r.adminStore.FindByID(ctx, id)
}

func (r *Repo) ListStudentsByKelas(ctx context.Context, kelas string, pageSize int, pageState string) ([]authmodel.Siswa, string, error) {
	if r.supabaseClient != nil && r.supabaseClient.DB != nil {
		if pageState != "" {
			return nil, "", nil
		}

		query := `
			SELECT s.id, s.nis, s.nama, s.kelas, s.wali_id, 
			       a.nama AS wali_nama, 
			       w.nama AS walas_nama, 
			       s.created_at
			FROM akun_siswa s
			LEFT JOIN akun_admin a ON s.wali_id = a.id
			LEFT JOIN kelas k ON s.kelas = k.nama
			LEFT JOIN akun_admin w ON k.walas_id = w.id
		`
		var args []any
		if kelas != "" {
			query += " WHERE s.kelas = $1"
			args = append(args, kelas)
		}
		query += " ORDER BY s.nama ASC"

		rows, err := r.supabaseClient.DB.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, "", err
		}
		defer rows.Close()

		students := []authmodel.Siswa{}
		for rows.Next() {
			var s authmodel.Siswa
			var waliID sql.NullString
			var waliNama sql.NullString
			var walasNama sql.NullString
			var createdAt time.Time

			err := rows.Scan(&s.ID, &s.NIS, &s.Nama, &s.Kelas, &waliID, &waliNama, &walasNama, &createdAt)
			if err != nil {
				return nil, "", err
			}
			s.Walas = waliNama.String
			s.WaliKelas = walasNama.String
			s.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
			students = append(students, s)
		}
		return students, "", nil
	}
	return r.studentStore.ListPaged(ctx, kelas, pageSize, pageState)
}
