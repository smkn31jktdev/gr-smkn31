package lomba

import (
	"context"
	"database/sql"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	lombamodel "be-gr31/internal/model/lomba"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
)

// Menangani operasi database untuk lomba kebersihan kelas.
type Repo struct {
	store          *astra.KebersihanStore
	studentStore   *astra.StudentStore
	adminStore     *astra.AdminStore
	supabaseClient *supabase.Client
}

// Membuat instance Repo baru
func NewRepo(store *astra.KebersihanStore, studentStore *astra.StudentStore, adminStore *astra.AdminStore, supabaseClient *supabase.Client) *Repo {
	return &Repo{
		store:          store,
		studentStore:   studentStore,
		adminStore:     adminStore,
		supabaseClient: supabaseClient,
	}
}

func (r *Repo) FindStudentByNISN(ctx context.Context, nisn string) (*authmodel.Siswa, error) {
	if r.supabaseClient != nil && r.supabaseClient.DB != nil {
		var s authmodel.Siswa
		var waliID sql.NullString
		var waliNama sql.NullString
		var walasNama sql.NullString
		var password string
		var createdAt time.Time
		query := `
			SELECT s.id, s.nis, s.nama, s.kelas, s.wali_id, 
			       a.nama AS wali_nama, 
			       w.nama AS walas_nama, 
			       s.password, s.created_at
			FROM akun_siswa s
			LEFT JOIN akun_admin a ON s.wali_id = a.id
			LEFT JOIN kelas k ON s.kelas = k.nama
			LEFT JOIN akun_admin w ON k.walas_id = w.id
			WHERE s.nis = $1
		`
		err := r.supabaseClient.DB.QueryRowContext(ctx, query, nisn).Scan(&s.ID, &s.NIS, &s.Nama, &s.Kelas, &waliID, &waliNama, &walasNama, &password, &createdAt)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		s.Walas = waliNama.String
		s.WaliKelas = walasNama.String
		s.Password = password
		s.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
		return &s, nil
	}
	return r.studentStore.FindByNISN(ctx, nisn)
}

func (r *Repo) Create(ctx context.Context, data *lombamodel.KebersihanKelas) error {
	return r.store.Create(ctx, data)
}

func (r *Repo) Update(ctx context.Context, data *lombamodel.KebersihanKelas) error {
	return r.store.Update(ctx, data)
}

func (r *Repo) FindByID(ctx context.Context, id string) (*lombamodel.KebersihanKelas, error) {
	return r.store.FindByID(ctx, id)
}

func (r *Repo) FindByKelasTanggal(ctx context.Context, kelas, tanggal string) (*lombamodel.KebersihanKelas, error) {
	return r.store.FindByKelasTanggal(ctx, kelas, tanggal)
}

func (r *Repo) FindByKelasDateRange(ctx context.Context, kelas, dariTgl, sampaiTgl string) ([]lombamodel.KebersihanKelas, error) {
	filter := lombamodel.KebersihanFilter{
		Kelas:     kelas,
		DariTgl:   dariTgl,
		SampaiTgl: sampaiTgl,
	}
	items, _, err := r.store.ListPaged(ctx, filter, 1, "")
	return items, err
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	return r.store.Delete(ctx, id)
}

func (r *Repo) ListPaged(ctx context.Context, filter lombamodel.KebersihanFilter, pageSize int, pageState string) ([]lombamodel.KebersihanKelas, string, error) {
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
