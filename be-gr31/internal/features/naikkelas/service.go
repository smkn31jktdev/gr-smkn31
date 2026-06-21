package naikkelas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/util"
)

// HasilNaikKelas merangkum hasil proses naik kelas satu tahun ajaran.
type HasilNaikKelas struct {
	TahunAjaran   string   `json:"tahunAjaran"` // cth: "2026/2027"
	TotalDiproses int      `json:"totalDiproses"`
	TotalBerhasil int      `json:"totalBerhasil"`
	TotalGagal    int      `json:"totalGagal"`
	Detail        []string `json:"detail,omitempty"` // error per-siswa jika ada
	DurasiMs      int64    `json:"durasiMs"`
	ExecutedAt    string   `json:"executedAt"`
}

// Service naik kelas — support Supabase (PostgreSQL) dan AstraDB.
type Service struct {
	astraClient    *astra.Client
	studentStore   *astra.StudentStore
	supabaseClient *supabase.Client
}

func NewService(astraClient *astra.Client, studentStore *astra.StudentStore, supabaseClient *supabase.Client) *Service {
	return &Service{
		astraClient:    astraClient,
		studentStore:   studentStore,
		supabaseClient: supabaseClient,
	}
}

func useSupabase(sc *supabase.Client) bool {
	return sc != nil && sc.DB != nil
}

// naikkanTingkat
func naikkanTingkat(kelas string) (string, bool) {
	parts := strings.Fields(strings.TrimSpace(kelas))
	if len(parts) == 0 {
		return kelas, false
	}
	switch strings.ToUpper(parts[0]) {
	case "X":
		parts[0] = "XI"
		return strings.Join(parts, " "), true
	case "XI":
		parts[0] = "XII"
		return strings.Join(parts, " "), true
	default:
		// XII atau format tidak dikenal → tidak diproses
		return kelas, false
	}
}

// TahunAjaranSekarang
func TahunAjaranSekarang(t time.Time) string {
	year := t.Year()
	if t.Month() >= time.July {
		return fmt.Sprintf("%d/%d", year, year+1)
	}
	return fmt.Sprintf("%d/%d", year-1, year)
}

// SudahDieksekusi
func (s *Service) SudahDieksekusi(ctx context.Context, tahunAjaran string) (bool, error) {
	if !useSupabase(s.supabaseClient) {
		return false, nil
	}
	var count int
	err := s.supabaseClient.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM naik_kelas_log WHERE tahun_ajaran = $1`, tahunAjaran,
	).Scan(&count)
	if err != nil {
		return false, nil
	}
	return count > 0, nil
}

// catatLog
func (s *Service) catatLog(ctx context.Context, hasil *HasilNaikKelas) {
	if !useSupabase(s.supabaseClient) {
		return
	}
	// Buat tabel jika belum ada
	_, _ = s.supabaseClient.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS naik_kelas_log (
			id             SERIAL PRIMARY KEY,
			tahun_ajaran   TEXT NOT NULL UNIQUE,
			total_diproses INT  DEFAULT 0,
			total_berhasil INT  DEFAULT 0,
			total_gagal    INT  DEFAULT 0,
			executed_at    TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	_, err := s.supabaseClient.DB.ExecContext(ctx, `
		INSERT INTO naik_kelas_log (tahun_ajaran, total_diproses, total_berhasil, total_gagal, executed_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (tahun_ajaran) DO NOTHING
	`, hasil.TahunAjaran, hasil.TotalDiproses, hasil.TotalBerhasil, hasil.TotalGagal, hasil.ExecutedAt)
	if err != nil {
		log.Printf("[naik-kelas] WARNING: gagal mencatat log: %v", err)
	}
}

// JalankanNaikKelas
func (s *Service) JalankanNaikKelas(ctx context.Context, tahunAjaran string, force bool) (*HasilNaikKelas, error) {
	mulai := time.Now()

	if !force {
		sudah, err := s.SudahDieksekusi(ctx, tahunAjaran)
		if err != nil {
			return nil, fmt.Errorf("cek log eksekusi: %w", err)
		}
		if sudah {
			return nil, fmt.Errorf(
				"naik kelas tahun ajaran %s sudah pernah dieksekusi; gunakan ?force=true untuk override",
				tahunAjaran,
			)
		}
	}

	hasil := &HasilNaikKelas{
		TahunAjaran: tahunAjaran,
		ExecutedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	if useSupabase(s.supabaseClient) {
		if err := s.jalankanSupabase(ctx, hasil); err != nil {
			return nil, err
		}
	} else {
		if err := s.jalankanAstra(ctx, hasil); err != nil {
			return nil, err
		}
	}

	hasil.DurasiMs = time.Since(mulai).Milliseconds()
	s.catatLog(ctx, hasil)

	log.Printf("[naik-kelas] Selesai — %d berhasil, %d gagal dari %d siswa (durasi: %dms)",
		hasil.TotalBerhasil, hasil.TotalGagal, hasil.TotalDiproses, hasil.DurasiMs)

	return hasil, nil
}

// Batch UPDATE
func (s *Service) jalankanSupabase(ctx context.Context, hasil *HasilNaikKelas) error {
	db := s.supabaseClient.DB

	var total int
	err := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM akun_siswa
		WHERE kelas ~ '^X(\s|$)' OR kelas ~ '^XI(\s|$)'
	`).Scan(&total)
	if err != nil {
		return fmt.Errorf("hitung siswa naik kelas: %w", err)
	}
	hasil.TotalDiproses = total

	resXII, err := db.ExecContext(ctx, `
		UPDATE akun_siswa
		SET kelas = REGEXP_REPLACE(kelas, '^XI(\s|$)', 'XII\1')
		WHERE kelas ~ '^XI(\s|$)'
	`)
	if err != nil {
		hasil.Detail = append(hasil.Detail, fmt.Sprintf("UPDATE XI→XII gagal: %v", err))
		hasil.TotalGagal = total
		return nil
	}
	rowsXII, _ := resXII.RowsAffected()

	resXI, err := db.ExecContext(ctx, `
		UPDATE akun_siswa
		SET kelas = REGEXP_REPLACE(kelas, '^X(\s|$)', 'XI\1')
		WHERE kelas ~ '^X(\s|$)'
	`)
	if err != nil {
		hasil.Detail = append(hasil.Detail, fmt.Sprintf("UPDATE X→XI gagal: %v", err))
		hasil.TotalBerhasil = int(rowsXII)
		hasil.TotalGagal = total - int(rowsXII)
		return nil
	}
	rowsXI, _ := resXI.RowsAffected()

	hasil.TotalBerhasil = int(rowsXII + rowsXI)
	hasil.TotalGagal = total - hasil.TotalBerhasil
	return nil
}

// jalankanAstra
func (s *Service) jalankanAstra(ctx context.Context, hasil *HasilNaikKelas) error {
	allSiswa, err := util.FetchAll(ctx,
		util.PagedFetcher[authmodel.Siswa](func(ctx context.Context, size int, state string) ([]authmodel.Siswa, string, error) {
			return s.studentStore.ListPaged(ctx, "", size, state)
		}),
		100,
	)
	if err != nil {
		return fmt.Errorf("fetch all siswa dari AstraDB: %w", err)
	}

	for _, siswa := range allSiswa {
		kelasLama := siswa.Kelas
		kelasBaru, perluNaik := naikkanTingkat(kelasLama)
		if !perluNaik {
			continue
		}
		hasil.TotalDiproses++

		docURL := s.astraClient.CollectionDocumentURL("akun_siswa", siswa.ID)
		patch := map[string]interface{}{"kelas": kelasBaru}

		if err := s.astraClient.Patch(ctx, docURL, patch); err != nil {
			hasil.TotalGagal++
			hasil.Detail = append(hasil.Detail,
				fmt.Sprintf("NIS=%s kelas=%s→%s err: %v", siswa.NIS, kelasLama, kelasBaru, err))
		} else {
			hasil.TotalBerhasil++
		}
	}

	return nil
}
