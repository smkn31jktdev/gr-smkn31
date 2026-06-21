package kehadiran

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"be-gr31/internal/config"
	"be-gr31/internal/features/kehadiran/fetch/hadir"
	"be-gr31/internal/features/kehadiran/fetch/izin"
	"be-gr31/internal/features/kehadiran/monitoring"
	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	kehadiranmodel "be-gr31/internal/model/kehadiran"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/util"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrDuplicate    = kehadiranmodel.ErrDuplicate
	ErrOutOfTime    = kehadiranmodel.ErrOutOfTime
	ErrWeekend      = kehadiranmodel.ErrWeekend
	ErrTooFar       = kehadiranmodel.ErrTooFar
	ErrNotFound     = kehadiranmodel.ErrNotFound
	ErrUnauthorized = kehadiranmodel.ErrUnauthorized
	ErrBadRequest   = kehadiranmodel.ErrBadRequest
)

// Service
type Service struct {
	repo           *Repo
	rekapSvc       *RekapService
	cfg            *config.Config
	rdb            *redis.Client
	supabaseClient *supabase.Client
	studentDB      studentFinder
	monitoringSvc  *monitoring.Service
	adminStore     *astra.AdminStore
}

// studentFinder
type studentFinder interface {
	FindSiswaByNISN(ctx context.Context, nis string) (interface{}, error)
}

// simpleStudentFinder
type simpleStudentFinder struct {
	store interface {
		FindByNISN(ctx context.Context, nis string) (interface{}, error)
	}
}

// Membuat instance Service baru
func NewService(repo *Repo, rekapSvc *RekapService, cfg *config.Config, rdb *redis.Client, supabaseClient *supabase.Client, adminStore *astra.AdminStore) *Service {
	monitoringSvc := monitoring.NewService(repo, rekapSvc, rdb)
	return &Service{
		repo:           repo,
		rekapSvc:       rekapSvc,
		cfg:            cfg,
		rdb:            rdb,
		supabaseClient: supabaseClient,
		monitoringSvc:  monitoringSvc,
		adminStore:     adminStore,
	}
}

func (s *Service) resolveSiswaInfo(ctx context.Context, nis string) (string, string, error) {
	if s.supabaseClient != nil && s.supabaseClient.DB != nil {
		var nama, kelas string
		query := `SELECT nama, kelas FROM akun_siswa WHERE nis = $1`
		err := s.supabaseClient.DB.QueryRowContext(ctx, query, nis).Scan(&nama, &kelas)
		if err == nil {
			return nama, kelas, nil
		}
	}
	if s.rekapSvc != nil && s.rekapSvc.studentStore != nil {
		siswa, err := s.rekapSvc.studentStore.FindByNISN(ctx, nis)
		if err == nil && siswa != nil {
			return siswa.Nama, siswa.Kelas, nil
		}
	}
	return "", "", nil
}

// Create memproses absensi harian siswa
func (s *Service) Create(ctx context.Context, siswa SiswaInfo, req kehadiranmodel.AbsenRequest) (res *kehadiranmodel.Kehadiran, err error) {
	defer func() {
		if err != nil && req.FotoIzin != "" {
			s.cleanupUploadedFile(req.FotoIzin)
		}
	}()

	now := util.NowJakarta()
	tanggal := util.FormatTanggal(now)

	// Validasi hari sekolah
	if util.IsWeekend(now) {
		return nil, ErrWeekend
	}

	// Validasi jam absensi
	if err := util.ValidateAbsenTime(now, s.cfg.AbsensiStartHour, s.cfg.AbsensiStartMinute, s.cfg.AbsensiEndHour, s.cfg.AbsensiEndMinute); err != nil {
		return nil, ErrOutOfTime
	}

	// Validasi izin/sakit
	if err := izin.Validate(req.Status, req.Alasan, req.FotoIzin); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadRequest, err.Error())
	}

	// Verifikasi jarak + accuracy
	var jarak float64
	isQR := strings.HasPrefix(req.Tipe, "qr")
	if req.Status == "hadir" && !isQR {
		jarak, err = hadir.ValidateGPS(
			req.Koordinat, req.Status,
			s.cfg.SekolahLat, s.cfg.SekolahLng, s.cfg.SekolahRadiusMeter,
			req.Akurasi, s.cfg.MinGPSAccuracyMeter, s.cfg.MaxGPSAccuracyMeter,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrTooFar, err.Error())
		}
	} else {
		if req.Koordinat != nil {
			jarak = util.Haversine(req.Koordinat.Lat, req.Koordinat.Lng, s.cfg.SekolahLat, s.cfg.SekolahLng)
		}

		// Validasi token QR untuk absensi tipe QR
		if req.Status == "hadir" && isQR {
			token := ""
			if strings.HasPrefix(req.Tipe, "qr:") {
				token = strings.TrimPrefix(req.Tipe, "qr:")
			}
			expectedToken := fmt.Sprintf("SMKN31-ATTENDANCE-KEY-%s", tanggal)
			if token != expectedToken {
				return nil, fmt.Errorf("%w: QR Code tidak valid atau kedaluwarsa", ErrBadRequest)
			}
		}
	}

	// Cek duplikasi
	existing, err := s.repo.FindByNISNTanggal(ctx, siswa.NIS, tanggal)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicate
	}

	// Look up name and class if they are empty
	if siswa.Nama == "" || siswa.Kelas == "" {
		nama, kelas, err := s.resolveSiswaInfo(ctx, siswa.NIS)
		if err == nil && nama != "" {
			siswa.Nama = nama
			siswa.Kelas = kelas
		}
	}

	// Build model kehadiran
	classifiedStatus := ClassifyStatus(req.Status, req.Alasan, req.FotoIzin)
	data := &kehadiranmodel.Kehadiran{
		ID:         uuid.New().String(),
		NIS:        siswa.NIS,
		NamaSiswa:  siswa.Nama,
		Kelas:      siswa.Kelas,
		Tanggal:    tanggal,
		Hari:       util.NamaHari(now),
		Status:     classifiedStatus,
		WaktuAbsen: util.FormatWaktu(now),
		Alasan:     req.Alasan,
		Koordinat:  req.Koordinat,
		Jarak:      jarak,
		Akurasi:    req.Akurasi,
		DeviceInfo: req.DeviceInfo,
		FotoIzin:   req.FotoIzin,
		CreatedAt:  common.FlexTime(now.Format(time.RFC3339)),
		UpdatedAt:  common.FlexTime(now.Format(time.RFC3339)),
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return nil, err
	}

	// Dual-write rekap_absensi (WAJIB)
	bulanTahun := util.BulanTahun(now)
	if err := s.rekapSvc.IncrementRekap(ctx, siswa.NIS, siswa.Nama, siswa.Kelas, bulanTahun, classifiedStatus); err != nil {
		log.Printf("WARNING: dual-write rekap gagal untuk %s: %v", siswa.NIS, err)
	}

	// Invalidasi cache list kehadiran
	s.InvalidateKehadiranListCache(ctx, tanggal)

	return data, nil
}

// CreateByAdmin memproses absensi manual oleh admin
func (s *Service) CreateByAdmin(ctx context.Context, req kehadiranmodel.AdminAbsenRequest, siswaInfo SiswaInfo) (res *kehadiranmodel.Kehadiran, err error) {
	defer func() {
		if err != nil && req.FotoIzin != "" {
			s.cleanupUploadedFile(req.FotoIzin)
		}
	}()

	// Validasi format tanggal
	if err := util.ValidateDateFormat(req.Tanggal); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrBadRequest, err.Error())
	}

	// Validasi izin/sakit memerlukan alasan dan foto
	if req.Status == "izin" || req.Status == "sakit" {
		if req.Alasan == "" {
			return nil, fmt.Errorf("%w: alasan wajib diisi untuk status %s", ErrBadRequest, req.Status)
		}
	}

	// Cek duplikasi
	existing, err := s.repo.FindByNISNTanggal(ctx, req.NIS, req.Tanggal)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicate
	}

	// Parse tanggal untuk nama hari
	t, _ := util.ParseDate(req.Tanggal)
	now := common.FlexTime(time.Now().Format(time.RFC3339))

	// Look up name and class if they are empty
	if siswaInfo.Nama == "" || siswaInfo.Kelas == "" {
		nama, kelas, err := s.resolveSiswaInfo(ctx, req.NIS)
		if err == nil && nama != "" {
			siswaInfo.Nama = nama
			siswaInfo.Kelas = kelas
		}
	}

	classifiedStatus := ClassifyStatus(req.Status, req.Alasan, req.FotoIzin)
	data := &kehadiranmodel.Kehadiran{
		ID:         uuid.New().String(),
		NIS:        req.NIS,
		NamaSiswa:  siswaInfo.Nama,
		Kelas:      siswaInfo.Kelas,
		Tanggal:    req.Tanggal,
		Hari:       util.NamaHari(t),
		Status:     classifiedStatus,
		WaktuAbsen: util.FormatWaktu(t),
		Alasan:     req.Alasan,
		Koordinat:  req.Koordinat,
		Akurasi:    req.Akurasi,
		DeviceInfo: req.DeviceInfo,
		FotoIzin:   req.FotoIzin,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return nil, err
	}

	// Dual-write rekap
	bulanTahun := req.Tanggal[:7]
	if err := s.rekapSvc.IncrementRekap(ctx, req.NIS, siswaInfo.Nama, siswaInfo.Kelas, bulanTahun, classifiedStatus); err != nil {
		log.Printf("WARNING: dual-write rekap gagal: %v", err)
	}

	// Invalidasi cache list kehadiran
	s.InvalidateKehadiranListCache(ctx, req.Tanggal)

	return data, nil
}

func (s *Service) findAdminByID(ctx context.Context, id string) (*authmodel.Admin, error) {
	if s.supabaseClient != nil && s.supabaseClient.DB != nil {
		var a authmodel.Admin
		var createdAt time.Time
		query := `
			SELECT id, nama, email, password, is_walas, kelas, role, created_at
			FROM akun_admin
			WHERE id = $1
		`
		var isWalas bool
		err := s.supabaseClient.DB.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Nama, &a.Email, &a.Password, &isWalas, &a.Kelas, &a.Role, &createdAt)
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
	return s.adminStore.FindByID(ctx, id)
}

// List mengambil daftar kehadiran dengan filter + fuzzy + pagination.
func (s *Service) List(ctx context.Context, filter kehadiranmodel.KehadiranFilter) ([]kehadiranmodel.Kehadiran, bool, int, error) {
	isGuruWali := filter.AdminRole != "" && filter.AdminRole != "super_admin"
	if isGuruWali {
		admin, err := s.findAdminByID(ctx, filter.AdminID)
		if err == nil && admin != nil {
			if admin.IsWalas || admin.Role == "walas" || admin.Role == "guru_wali" {
				filter.Kelas = admin.Kelas
			}
		}
	}
	return s.monitoringSvc.List(ctx, filter)
}

// InvalidateKehadiranListCache
func (s *Service) InvalidateKehadiranListCache(ctx context.Context, tanggal string) {
	s.monitoringSvc.InvalidateKehadiranListCache(ctx, tanggal)
}

// Delete menghapus kehadiran dan decrement rekap
func (s *Service) Delete(ctx context.Context, id string) error {
	data, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// Decrement rekap
	bulanTahun := data.Tanggal[:7]
	if err := s.rekapSvc.DecrementRekap(ctx, data.NIS, data.NamaSiswa, data.Kelas, bulanTahun, data.Status); err != nil {
		log.Printf("WARNING: decrement rekap gagal: %v", err)
	}

	// Invalidasi cache list kehadiran
	s.InvalidateKehadiranListCache(ctx, data.Tanggal)

	return nil
}

// SiswaInfo adalah data minimal siswa yang diperlukan untuk membuat kehadiran.
type SiswaInfo = kehadiranmodel.SiswaInfo

// Mengubah status kehadiran oleh admin (override)
func (s *Service) UpdateByAdmin(ctx context.Context, id, statusBaru, alasan string) (*kehadiranmodel.Kehadiran, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	classifiedStatus := ClassifyStatus(statusBaru, alasan, existing.FotoIzin)

	fields := map[string]any{
		"status": classifiedStatus,
		"alasan": alasan,
	}
	if err := s.repo.Update(ctx, id, fields); err != nil {
		return nil, err
	}

	// Dual-write: swap counter lama → baru
	bulanTahun := existing.Tanggal[:7]
	if existing.Status != classifiedStatus {
		if err := s.rekapSvc.DecrementRekap(ctx, existing.NIS, existing.NamaSiswa, existing.Kelas, bulanTahun, existing.Status); err != nil {
			log.Printf("WARNING: decrement rekap lama gagal: %v", err)
		}
		if err := s.rekapSvc.IncrementRekap(ctx, existing.NIS, existing.NamaSiswa, existing.Kelas, bulanTahun, classifiedStatus); err != nil {
			log.Printf("WARNING: increment rekap baru gagal: %v", err)
		}
	}

	// Invalidasi cache list kehadiran
	s.InvalidateKehadiranListCache(ctx, existing.Tanggal)

	existing.Status = classifiedStatus
	existing.Alasan = alasan
	return existing, nil
}

// GetByID mengambil satu record kehadiran berdasarkan ID (untuk admin detail view).
func (s *Service) GetByID(ctx context.Context, id string) (*kehadiranmodel.Kehadiran, error) {
	data, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrNotFound
	}
	return data, nil
}

// SekolahRadiusMeter mengembalikan radius sekolah dari konfigurasi.
func (s *Service) SekolahRadiusMeter() float64 {
	return s.cfg.SekolahRadiusMeter
}

// Fungsi tidak_hadir secara dinamis berdasarkan alasan dan lampiran
func ClassifyStatus(status, alasan, fotoIzin string) string {
	return monitoring.ClassifyStatus(status, alasan, fotoIzin)
}

func (s *Service) cleanupUploadedFile(fotoPath string) {
	if fotoPath == "" {
		return
	}
	cleanPath := strings.TrimPrefix(fotoPath, "/uploads")
	filePath := filepath.Join(s.cfg.UploadDir, cleanPath)
	if err := os.Remove(filePath); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("WARNING: gagal menghapus file terunggah %s: %v", filePath, err)
		}
	} else {
		log.Printf("INFO: berhasil menghapus file terunggah yang tidak jadi digunakan: %s", filePath)
	}
}
