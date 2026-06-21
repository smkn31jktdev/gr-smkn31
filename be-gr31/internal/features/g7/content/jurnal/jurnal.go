package jurnal

import (
	"context"
	"errors"
	"log"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/util"

	"be-gr31/internal/features/g7/content/wali"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound     = errors.New("data G7 tidak ditemukan")
	ErrFutureDate   = errors.New("tidak bisa input untuk tanggal mendatang")
	ErrNotTodayDate = errors.New("hanya diperbolehkan mengisi jurnal untuk hari ini")
)

type G7Store interface {
	FindByNISNTanggal(ctx context.Context, nisn, tanggal string) (*g7model.G7, error)
	Upsert(ctx context.Context, data *g7model.G7) error
	FindByID(ctx context.Context, id string) (*g7model.G7, error)
	Delete(ctx context.Context, id string) error
	ListPaged(ctx context.Context, filter g7model.G7Filter, pageSize int, pageState string) ([]g7model.G7, string, error)
}

type StudentFinder interface {
	FindStudentByNISN(ctx context.Context, nisn string) (*authmodel.Siswa, error)
}

type RosterFetcher interface {
	FetchRoster(ctx context.Context, kelas string) []authmodel.Siswa
}

type RekapService interface {
	UpdateRekapHariTercatat(ctx context.Context, nisn, bulanTahun string) error
}

type Service struct {
	store         G7Store
	studentFinder StudentFinder
	rdb           *redis.Client
	rosterFetcher RosterFetcher
	rekapSvc      RekapService
	waliSvc       *wali.Service
}

func NewService(
	store G7Store,
	studentFinder StudentFinder,
	rdb *redis.Client,
	rosterFetcher RosterFetcher,
	rekapSvc RekapService,
	waliSvc *wali.Service,
) *Service {
	return &Service{
		store:         store,
		studentFinder: studentFinder,
		rdb:           rdb,
		rosterFetcher: rosterFetcher,
		rekapSvc:      rekapSvc,
		waliSvc:       waliSvc,
	}
}

// Upsert saves or updates a student's daily G7 log
func (s *Service) Upsert(ctx context.Context, nisn string, req g7model.G7UpsertRequest) (*g7model.G7, error) {
	if err := util.ValidateDateFormat(req.Tanggal); err != nil {
		return nil, err
	}

	todayStr := util.NowJakarta().Format("2006-01-02")
	if req.Tanggal != todayStr {
		return nil, ErrNotTodayDate
	}

	if err := validateAktivitasLength(req); err != nil {
		return nil, err
	}
	namaSiswa, kelas := s.resolveSiswa(ctx, nisn)

	now := common.FlexTime(time.Now().Format(time.RFC3339))

	existing, err := s.store.FindByNISNTanggal(ctx, nisn, req.Tanggal)
	if err != nil {
		return nil, err
	}

	var data g7model.G7
	if existing != nil {
		data = *existing
		data.NamaSiswa = namaSiswa
		data.Kelas = kelas
		data.UpdatedAt = now
	} else {
		data = g7model.G7{
			ID:        uuid.New().String(),
			NISN:      nisn,
			NamaSiswa: namaSiswa,
			Kelas:     kelas,
			Tanggal:   req.Tanggal,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	data.Bangun = req.Bangun
	data.Ibadah = req.Ibadah
	data.Makan = req.Makan
	data.Olahraga = req.Olahraga
	data.Belajar = req.Belajar
	data.Bermasyarakat = req.Bermasyarakat
	data.Tidur = req.Tidur

	data.TotalDone = countDone(&data)

	if err := s.store.Upsert(ctx, &data); err != nil {
		return nil, err
	}

	bulanTahun := req.Tanggal[:7]
	cacheKey := "cache:g7:summary:" + nisn + ":" + bulanTahun
	if err := s.rdb.Del(ctx, cacheKey).Err(); err != nil {
		log.Printf("WARNING: gagal invalidasi cache G7 summary: %v", err)
	}

	// Update g7_rekap for this student and month
	if err := s.rekapSvc.UpdateRekapHariTercatat(ctx, nisn, bulanTahun); err != nil {
		log.Printf("WARNING: gagal update HariTercatat di g7_rekap: %v", err)
	}

	return &data, nil
}

// DashboardSiswa retrieves the G7 dashboard data for a student
func (s *Service) DashboardSiswa(ctx context.Context, nisn string) (*g7model.G7DashboardSiswa, error) {
	now := util.NowJakarta()
	tanggalHariIni := now.Format("2006-01-02")
	bulanTahun := util.BulanTahun(now)

	jurnalHariIni, err := s.store.FindByNISNTanggal(ctx, nisn, tanggalHariIni)
	if err != nil {
		return nil, err
	}

	progresHariIni := 0
	if jurnalHariIni != nil {
		progresHariIni = jurnalHariIni.TotalDone
	}

	filter := g7model.G7Filter{
		NISN:      nisn,
		BulanDari: bulanTahun,
		BulanKe:   bulanTahun,
	}
	fetcher := util.PagedFetcher[g7model.G7](func(ctx context.Context, size int, state string) ([]g7model.G7, string, error) {
		return s.store.ListPaged(ctx, filter, size, state)
	})
	semuaJurnal, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, err
	}

	totalDoneSum := 0
	for _, j := range semuaJurnal {
		totalDoneSum += j.TotalDone
	}
	hariTercatat := len(semuaJurnal)
	rataRata := 0.0
	if hariTercatat > 0 {
		rataRata = float64(totalDoneSum) / float64(hariTercatat)
	}

	return &g7model.G7DashboardSiswa{
		JurnalHariIni:  jurnalHariIni,
		ProgresHariIni: progresHariIni,
		RingkasanBulan: g7model.G7RingkasanBulan{
			BulanTahun:   bulanTahun,
			HariTercatat: hariTercatat,
			RataRataDone: rataRata,
			TotalDoneSum: totalDoneSum,
		},
	}, nil
}

// GetByTanggal retrieves a G7 daily log by student NISN and date
func (s *Service) GetByTanggal(ctx context.Context, nisn, tanggal string) (*g7model.G7, error) {
	data, err := s.store.FindByNISNTanggal(ctx, nisn, tanggal)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, ErrNotFound
	}
	return data, nil
}

// List retrieves a list of G7 logs based on the filter
func (s *Service) List(ctx context.Context, filter g7model.G7Filter) ([]g7model.G7, bool, int, error) {
	if filter.NISN != "" {
		fetcher := util.PagedFetcher[g7model.G7](func(ctx context.Context, size int, state string) ([]g7model.G7, string, error) {
			return s.store.ListPaged(ctx, filter, size, state)
		})
		all, err := util.FetchAll(ctx, fetcher, 100)
		if err != nil {
			return nil, false, 0, err
		}
		result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
		return result, hasMore, total, nil
	}

	var targetStudents []authmodel.Siswa
	isGuruWali, asuhan, err := s.waliSvc.CheckIsGuruWali(ctx, filter.AdminID, filter.AdminRole)
	if err != nil {
		return nil, false, 0, err
	}

	if isGuruWali {
		targetStudents = asuhan
	} else {
		// Ambil seluruh siswa secara lengkap dengan auto-paging
		targetStudents = s.rosterFetcher.FetchRoster(ctx, "")
	}

	if len(targetStudents) == 0 {
		return []g7model.G7{}, false, 0, nil
	}

	// Ambil semua data G7 hari ini
	todayStr := util.NowJakarta().Format("2006-01-02")
	todayFilter := g7model.G7Filter{
		Tanggal: todayStr,
	}
	g7Fetcher := util.PagedFetcher[g7model.G7](func(ctx context.Context, size int, state string) ([]g7model.G7, string, error) {
		return s.store.ListPaged(ctx, todayFilter, size, state)
	})
	todayJurnals, err := util.FetchAll(ctx, g7Fetcher, 100)
	if err != nil {
		log.Printf("WARNING: gagal mengambil data jurnal hari ini: %v", err)
	}

	jurnalMap := make(map[string]g7model.G7)
	for _, g := range todayJurnals {
		jurnalMap[g.NISN] = g
	}

	// Gabungkan data profil siswa dengan data jurnal G7 hari ini
	allG7 := make([]g7model.G7, 0, len(targetStudents))
	for _, st := range targetStudents {
		if st.NIS == "" {
			continue
		}
		if g7data, found := jurnalMap[st.NIS]; found {
			allG7 = append(allG7, g7data)
		} else {
			// Buat dummy G7 jika belum mengisi hari ini
			allG7 = append(allG7, g7model.G7{
				ID:        "", // Penanda dummy
				NISN:      st.NIS,
				NamaSiswa: st.Nama,
				Kelas:     st.Kelas,
				Tanggal:   todayStr,
				TotalDone: 0,
			})
		}
	}

	// Terapkan filter query fuzzy jika ada pencarian
	if filter.Query != "" {
		allG7 = util.FuzzyFilter(allG7, filter.Query, func(g g7model.G7) []string {
			return []string{g.NamaSiswa, g.Kelas, g.NISN}
		})
	}

	result, hasMore, total := util.PaginateSlice(allG7, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}

// Delete removes a daily G7 log by ID
func (s *Service) Delete(ctx context.Context, id string) error {
	data, err := s.store.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNotFound
	}

	if err := s.store.Delete(ctx, id); err != nil {
		return err
	}

	bulanTahun := data.Tanggal[:7]
	cacheKey := "cache:g7:summary:" + data.NISN + ":" + bulanTahun
	_ = s.rdb.Del(ctx, cacheKey).Err()

	// Update g7_rekap count
	if err := s.rekapSvc.UpdateRekapHariTercatat(ctx, data.NISN, bulanTahun); err != nil {
		log.Printf("WARNING: gagal update HariTercatat setelah delete: %v", err)
	}

	return nil
}

// Summary calculates the monthly summary for students
func (s *Service) Summary(ctx context.Context, bulanTahun, kelas, adminID, adminRole string) ([]g7model.G7Summary, error) {
	isGuruWali, asuhan, err := s.waliSvc.CheckIsGuruWali(ctx, adminID, adminRole)
	if err != nil {
		return nil, err
	}

	var allowedNISN map[string]bool
	if isGuruWali {
		allowedNISN = make(map[string]bool)
		for _, st := range asuhan {
			if st.NIS != "" {
				allowedNISN[st.NIS] = true
			}
		}
	}

	filter := g7model.G7Filter{
		Kelas:     kelas,
		BulanDari: bulanTahun,
		BulanKe:   bulanTahun,
	}

	fetcher := util.PagedFetcher[g7model.G7](func(ctx context.Context, size int, state string) ([]g7model.G7, string, error) {
		return s.store.ListPaged(ctx, filter, size, state)
	})

	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, err
	}

	type aggData struct {
		NamaSiswa    string
		Kelas        string
		TotalDoneSum int
		HariCount    int
	}
	agg := make(map[string]*aggData)

	for _, g := range all {
		if isGuruWali && !allowedNISN[g.NISN] {
			continue
		}
		if _, ok := agg[g.NISN]; !ok {
			agg[g.NISN] = &aggData{NamaSiswa: g.NamaSiswa, Kelas: g.Kelas}
		}
		agg[g.NISN].TotalDoneSum += g.TotalDone
		agg[g.NISN].HariCount++
	}

	summaries := make([]g7model.G7Summary, 0, len(agg))
	for nisn, data := range agg {
		var rataRata float64
		if data.HariCount > 0 {
			rataRata = float64(data.TotalDoneSum) / float64(data.HariCount)
		}
		summaries = append(summaries, g7model.G7Summary{
			NISN:         nisn,
			NamaSiswa:    data.NamaSiswa,
			Kelas:        data.Kelas,
			BulanTahun:   bulanTahun,
			RataRataDone: rataRata,
			HariTercatat: data.HariCount,
		})
	}

	return summaries, nil
}

func (s *Service) resolveSiswa(ctx context.Context, nisn string) (nama, kelas string) {
	siswa, err := s.studentFinder.FindStudentByNISN(ctx, nisn)
	if err != nil || siswa == nil {
		log.Printf("WARNING: profil siswa NISN=%s tidak ditemukan: %v", nisn, err)
		return "", ""
	}
	return siswa.Nama, siswa.Kelas
}

func countDone(g *g7model.G7) int {
	count := 0
	aktivitas := []*g7model.Aktivitas{
		g.Bangun, g.Ibadah, g.Makan, g.Olahraga,
		g.Belajar, g.Bermasyarakat, g.Tidur,
	}
	for _, a := range aktivitas {
		if a != nil && a.Done {
			count++
		}
	}
	return count
}

func validateAktivitasLength(req g7model.G7UpsertRequest) error {
	checkAkt := func(a *g7model.Aktivitas, name string) error {
		if a == nil {
			return nil
		}
		if len([]rune(a.Keterangan)) > 255 {
			return errors.New("keterangan " + name + " maksimal 255 karakter")
		}
		if a.Waktu != "" {
			if err := util.ValidateTimeFormat(a.Waktu); err != nil {
				return errors.New("format waktu " + name + " harus HH:mm")
			}
		}
		return nil
	}

	sections := map[string]*g7model.Aktivitas{
		"bangun": req.Bangun, "ibadah": req.Ibadah, "makan": req.Makan,
		"olahraga": req.Olahraga, "belajar": req.Belajar,
		"bermasyarakat": req.Bermasyarakat, "tidur": req.Tidur,
	}
	for name, a := range sections {
		if err := checkAkt(a, name); err != nil {
			return err
		}
	}
	return nil
}
