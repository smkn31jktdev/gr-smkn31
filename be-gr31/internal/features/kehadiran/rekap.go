package kehadiran

import (
	"context"

	"be-gr31/internal/features/kehadiran/rekap/bulan"
	"be-gr31/internal/features/kehadiran/rekap/semester"
	authmodel "be-gr31/internal/model/auth"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"

	"github.com/redis/go-redis/v9"
)

type RekapService struct {
	store          *astra.RekapStore
	kehadiranSvc   KehadiranStore
	studentStore   *astra.StudentStore
	rdb            *redis.Client
	supabaseClient *supabase.Client
	bulanSvc       *bulan.Service
	semesterSvc    *semester.Service
}

// NewRekapService membuat RekapService baru.
func NewRekapService(store *astra.RekapStore, rdb *redis.Client, supabaseClient *supabase.Client) *RekapService {
	s := &RekapService{
		store:          store,
		rdb:            rdb,
		supabaseClient: supabaseClient,
	}
	s.initSubServices()
	return s
}

func (s *RekapService) initSubServices() {
	s.bulanSvc = bulan.NewService(s.store, s.kehadiranSvc, s.studentStore, s.rdb, s.supabaseClient)
	s.semesterSvc = semester.NewService(s.bulanSvc)
}

// SetKehadiranStore menyuntikkan KehadiranStore untuk query rekap harian
func (s *RekapService) SetKehadiranStore(ks KehadiranStore) {
	s.kehadiranSvc = ks
	s.initSubServices()
}

// SetStudentStore menyuntikkan StudentStore untuk join roster siswa pada rekap
func (s *RekapService) SetStudentStore(ss *astra.StudentStore) {
	s.studentStore = ss
	s.initSubServices()
}

// IncrementRekap melakukan increment counter setelah absensi baru dibuat
func (s *RekapService) IncrementRekap(ctx context.Context, nis, namaSiswa, kelas, bulanTahun, status string) error {
	return s.bulanSvc.IncrementRekap(ctx, nis, namaSiswa, kelas, bulanTahun, status)
}

// DecrementRekap melakukan decrement setelah absensi dihapus
func (s *RekapService) DecrementRekap(ctx context.Context, nis, namaSiswa, kelas, bulanTahun, status string) error {
	return s.bulanSvc.DecrementRekap(ctx, nis, namaSiswa, kelas, bulanTahun, status)
}

// Mengambil rekap bulanan satu siswa
func (s *RekapService) GetRekap(ctx context.Context, nis, bulanTahun string) (*rekapmodel.RekapBulanan, error) {
	return s.bulanSvc.GetRekap(ctx, nis, bulanTahun)
}

// Mengambil daftar rekap bulanan dengan filter + fuzzy + pagination
func (s *RekapService) ListRekap(ctx context.Context, filter rekapmodel.RekapFilter) ([]rekapmodel.RekapBulanan, bool, int, error) {
	return s.bulanSvc.ListRekap(ctx, filter)
}

// Mengambil ringkasan kehadiran untuk 1 hari
func (s *RekapService) GetRekapHarian(ctx context.Context, tanggal, kelas string) (*rekapmodel.RekapHarian, error) {
	return s.bulanSvc.GetRekapHarian(ctx, tanggal, kelas)
}

// Mengambil ringkasan semua siswa dalam 1 kelas untuk bulan tertentu
func (s *RekapService) GetRingkasanKelas(ctx context.Context, kelas, bulanTahun string, page, limit int) ([]rekapmodel.RingkasanSiswa, bool, int, error) {
	return s.bulanSvc.GetRingkasanKelas(ctx, kelas, bulanTahun, page, limit)
}

// Membangun rekap bulanan lengkap untuk admin
func (s *RekapService) GetRekapKelasLengkap(ctx context.Context, bulan, kelas string) (*rekapmodel.RekapKelasLengkap, error) {
	return s.bulanSvc.GetRekapKelasLengkap(ctx, bulan, kelas)
}

// Mengambil ringkasan persentase kehadiran per-kelas (ringan)
func (s *RekapService) GetPersentaseKehadiranKelas(ctx context.Context, bulan, kelas string) (*rekapmodel.RekapPersentaseKelas, error) {
	return s.bulanSvc.GetPersentaseKehadiranKelas(ctx, bulan, kelas)
}

// Membangun rekap kehadiran satu minggu dari data harian nyata
func (s *RekapService) GetRekapMingguanKelas(ctx context.Context, kelas, senin string) (*rekapmodel.RekapMingguanKelas, error) {
	return s.bulanSvc.GetRekapMingguanKelas(ctx, kelas, senin)
}

// Mengambil seluruh siswa pada kelas (auto-paging)
func (s *RekapService) FetchRoster(ctx context.Context, kelas string) []authmodel.Siswa {
	return s.bulanSvc.FetchRoster(ctx, kelas)
}

// Mengambil data kehadiran bulanan lengkap untuk satu siswa
func (s *RekapService) GetKehadiranBulananSiswa(ctx context.Context, nis, bulanTahun string) (*rekapmodel.KehadiranBulananSiswa, error) {
	return s.bulanSvc.GetKehadiranBulananSiswa(ctx, nis, bulanTahun)
}

// Mengambil data kehadiran mingguan lengkap untuk satu siswa
func (s *RekapService) GetKehadiranMingguanSiswa(ctx context.Context, nis, senin string) (*rekapmodel.KehadiranBulananSiswa, error) {
	return s.bulanSvc.GetKehadiranMingguanSiswa(ctx, nis, senin)
}

// Mengambil rekap per-bulan satu siswa dalam rentang dari–sampai
func (s *RekapService) GetRingkasanSiswa(ctx context.Context, nis, dari, sampai string) ([]rekapmodel.RekapBulanan, error) {
	return s.bulanSvc.GetRingkasanSiswa(ctx, nis, dari, sampai)
}

// Mengambil tren kehadiran kelas selama satu semester
func (s *RekapService) GetRekapSemesterKelas(ctx context.Context, kelas, semester string) (*rekapmodel.RekapSemesterKelas, error) {
	return s.semesterSvc.GetRekapSemesterKelas(ctx, kelas, semester)
}

// Mengambil seluruh unique kelas dari students database
func (s *RekapService) GetKelas(ctx context.Context) ([]string, error) {
	return s.bulanSvc.GetKelas(ctx)
}

// Mengambil seluruh unique kelas dan jurusan dari students database
func (s *RekapService) GetKelasJurusan(ctx context.Context) (*rekapmodel.KelasJurusanResponse, error) {
	return s.bulanSvc.GetKelasJurusan(ctx)
}
