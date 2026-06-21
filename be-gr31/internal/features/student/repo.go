package student

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

// Repo storage untuk siswa
type Repo struct {
	studentStore   *astra.StudentStore
	adminStore     *astra.AdminStore
	supabaseClient *supabase.Client
}

// Membuat instance Repo baru
func NewRepo(client *astra.Client, studentStore *astra.StudentStore, supabaseClient *supabase.Client) *Repo {
	return &Repo{
		studentStore:   studentStore,
		adminStore:     astra.NewAdminStore(client),
		supabaseClient: supabaseClient,
	}
}

func (r *Repo) UseSupabase() bool {
	return r.supabaseClient != nil && r.supabaseClient.DB != nil
}

func (r *Repo) FindByNIS(ctx context.Context, nis string) (*authmodel.Siswa, error) {
	if r.UseSupabase() {
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
		err := r.supabaseClient.DB.QueryRowContext(ctx, query, nis).Scan(&s.ID, &s.NIS, &s.Nama, &s.Kelas, &waliID, &waliNama, &walasNama, &password, &createdAt)
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
	return r.studentStore.FindByNISN(ctx, nis)
}

func (r *Repo) FindByID(ctx context.Context, id string) (*authmodel.Siswa, error) {
	if r.UseSupabase() {
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
			WHERE s.id = $1
		`
		err := r.supabaseClient.DB.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.NIS, &s.Nama, &s.Kelas, &waliID, &waliNama, &walasNama, &password, &createdAt)
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
	return r.studentStore.FindByID(ctx, id)
}

func (r *Repo) resolveWaliID(ctx context.Context, walasName string) (sql.NullString, error) {
	if walasName == "" {
		return sql.NullString{Valid: false}, nil
	}

	// Fetch all admins to map their names
	rows, err := r.supabaseClient.DB.QueryContext(ctx, "SELECT id, nama FROM akun_admin")
	if err != nil {
		return sql.NullString{Valid: false}, err
	}
	defer rows.Close()

	adminMap := make(map[string]string)
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return sql.NullString{Valid: false}, err
		}
		adminMap[strings.ToLower(strings.TrimSpace(name))] = id
	}

	normalized := strings.ToLower(strings.TrimSpace(walasName))
	// 1. Direct match
	if id, ok := adminMap[normalized]; ok {
		return sql.NullString{String: id, Valid: true}, nil
	}

	// 2. Custom typo/spelling fallbacks
	var corrected string
	switch normalized {
	case "meythannisa salsabila":
		corrected = "meythannisa salsabilla"
	case "sriyani":
		corrected = "sriyani spd"
	case "syafira khairuninisa":
		corrected = "syafira khairunnisa"
	case "melki felix m":
		corrected = "melki felix mapan"
	}
	if corrected != "" {
		if id, ok := adminMap[corrected]; ok {
			return sql.NullString{String: id, Valid: true}, nil
		}
	}

	// 3. Prefix/substring fallback match
	for adminName, id := range adminMap {
		if len(normalized) >= 4 && strings.HasPrefix(adminName, normalized) {
			return sql.NullString{String: id, Valid: true}, nil
		}
		if len(adminName) >= 4 && strings.HasPrefix(normalized, adminName) {
			return sql.NullString{String: id, Valid: true}, nil
		}
	}

	return sql.NullString{Valid: false}, nil
}

func (r *Repo) Upsert(ctx context.Context, data *authmodel.Siswa) error {
	if r.UseSupabase() {
		waliID, err := r.resolveWaliID(ctx, data.Walas)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO akun_siswa (id, nis, nama, kelas, wali_id, password)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET
				nis = EXCLUDED.nis,
				nama = EXCLUDED.nama,
				kelas = EXCLUDED.kelas,
				wali_id = EXCLUDED.wali_id,
				password = CASE WHEN EXCLUDED.password = '' THEN akun_siswa.password ELSE EXCLUDED.password END
		`
		_, err = r.supabaseClient.DB.ExecContext(ctx, query, data.ID, data.NIS, data.Nama, data.Kelas, waliID, data.Password)
		return err
	}
	return r.studentStore.Upsert(ctx, data)
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	if r.UseSupabase() {
		query := `DELETE FROM akun_siswa WHERE id = $1`
		_, err := r.supabaseClient.DB.ExecContext(ctx, query, id)
		return err
	}
	return r.studentStore.Delete(ctx, id)
}

func (r *Repo) ListPaged(ctx context.Context, kelas string, pageSize int, pageState string) ([]authmodel.Siswa, string, error) {
	return r.studentStore.ListPaged(ctx, kelas, pageSize, pageState)
}

func (r *Repo) ListSupabase(ctx context.Context, filter authmodel.SiswaFilter, isGuruWali bool, adminID string) ([]authmodel.Siswa, int, error) {
	var total int
	var queryCount string
	var querySelect string
	var args []interface{}

	// Build WHERE clauses
	whereClauses := []string{"1=1"}
	paramIdx := 1

	if filter.Kelas != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("s.kelas = $%d", paramIdx))
		args = append(args, filter.Kelas)
		paramIdx++
	}

	if filter.Query != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(s.nama ILIKE $%d OR s.nis ILIKE $%d OR a.nama ILIKE $%d OR w.nama ILIKE $%d)", paramIdx, paramIdx, paramIdx, paramIdx))
		args = append(args, "%"+filter.Query+"%")
		paramIdx++
	}

	if isGuruWali && adminID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("s.wali_id = $%d", paramIdx))
		args = append(args, adminID)
		paramIdx++
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// Query Count
	queryCount = fmt.Sprintf(`
		SELECT COUNT(*)
		FROM akun_siswa s
		LEFT JOIN akun_admin a ON s.wali_id = a.id
		LEFT JOIN kelas k ON s.kelas = k.nama
		LEFT JOIN akun_admin w ON k.walas_id = w.id
		WHERE %s
	`, whereSQL)

	err := r.supabaseClient.DB.QueryRowContext(ctx, queryCount, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count students: %w", err)
	}

	// Query Select with Pagination
	limitSQL := ""
	if filter.Limit > 0 {
		limitSQL = fmt.Sprintf("LIMIT $%d OFFSET $%d", paramIdx, paramIdx+1)
		args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)
	}

	querySelect = fmt.Sprintf(`
		SELECT s.id, s.nis, s.nama, s.kelas, s.wali_id, 
		       a.nama AS wali_nama, 
		       w.nama AS walas_nama, 
		       s.created_at
		FROM akun_siswa s
		LEFT JOIN akun_admin a ON s.wali_id = a.id
		LEFT JOIN kelas k ON s.kelas = k.nama
		LEFT JOIN akun_admin w ON k.walas_id = w.id
		WHERE %s
		ORDER BY s.nama ASC
		%s
	`, whereSQL, limitSQL)

	rows, err := r.supabaseClient.DB.QueryContext(ctx, querySelect, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("select students: %w", err)
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
			return nil, 0, fmt.Errorf("scan student row: %w", err)
		}
		s.Walas = waliNama.String
		s.WaliKelas = walasNama.String
		s.CreatedAt = common.FlexTime(createdAt.Format(time.RFC3339))
		students = append(students, s)
	}

	return students, total, nil
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
				a.Role = "walas"
			} else {
				a.Role = "admin"
			}
		}
		return &a, nil
	}
	return r.adminStore.FindByID(ctx, id)
}
