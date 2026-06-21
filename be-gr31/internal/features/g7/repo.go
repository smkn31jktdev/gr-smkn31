package g7

import (
	"context"
	"database/sql"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
)

// Repo
type Repo struct {
	store          *astra.G7Store
	rekapStore     *astra.G7RekapStore
	studentStore   *astra.StudentStore
	adminStore     *astra.AdminStore
	supabaseClient *supabase.Client
}

// NewRepo membuat instance Repo baru.
func NewRepo(store *astra.G7Store, rekapStore *astra.G7RekapStore, studentStore *astra.StudentStore, adminStore *astra.AdminStore, supabaseClient *supabase.Client) *Repo {
	return &Repo{
		store:          store,
		rekapStore:     rekapStore,
		studentStore:   studentStore,
		adminStore:     adminStore,
		supabaseClient: supabaseClient,
	}
}

// Jurnal harian
func (r *Repo) Upsert(ctx context.Context, data *g7model.G7) error {
	return r.store.Upsert(ctx, data)
}
func (r *Repo) FindByNISNTanggal(ctx context.Context, nisn, tanggal string) (*g7model.G7, error) {
	return r.store.FindByNISNTanggal(ctx, nisn, tanggal)
}

// Mengambil G7 berdasarkan ID
func (r *Repo) FindByID(ctx context.Context, id string) (*g7model.G7, error) {
	return r.store.FindByID(ctx, id)
}

// Menghapus G7 berdasarkan ID
func (r *Repo) Delete(ctx context.Context, id string) error {
	return r.store.Delete(ctx, id)
}

// Mengambil G7 dengan pagination cursor
func (r *Repo) ListPaged(ctx context.Context, filter g7model.G7Filter, pageSize int, pageState string) ([]g7model.G7, string, error) {
	return r.store.ListPaged(ctx, filter, pageSize, pageState)
}

// Rekap bulanan
func (r *Repo) UpsertRekap(ctx context.Context, data *g7model.G7Rekap) error {
	return r.rekapStore.Upsert(ctx, data)
}

// FindRekapByNISBulan mencari rekap berdasarkan NIS + bulanTahun.
func (r *Repo) FindRekapByNISBulan(ctx context.Context, nis, bulan string) (*g7model.G7Rekap, error) {
	return r.rekapStore.FindByNISNBulan(ctx, nis, bulan)
}

// FindRekapByID mengambil rekap berdasarkan ID.
func (r *Repo) FindRekapByID(ctx context.Context, id string) (*g7model.G7Rekap, error) {
	return r.rekapStore.FindByID(ctx, id)
}

// ListRekapPaged mengambil rekap dengan pagination cursor.
func (r *Repo) ListRekapPaged(ctx context.Context, filter g7model.G7RekapFilter, pageSize int, pageState string) ([]g7model.G7Rekap, string, error) {
	return r.rekapStore.ListPaged(ctx, filter, pageSize, pageState)
}

// DeleteRekap menghapus rekap berdasarkan ID.
func (r *Repo) DeleteRekap(ctx context.Context, id string) error {
	return r.rekapStore.Delete(ctx, id)
}

// Siswa
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

// ListStudentsByKelas mengambil siswa pada satu kelas
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

func (r *Repo) ListStudentsByWalas(ctx context.Context, walas string) ([]authmodel.Siswa, error) {
	return r.studentStore.ListByWalas(ctx, walas)
}
