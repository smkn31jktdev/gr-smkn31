package summary

import (
	"context"
	"errors"

	authmodel "be-gr31/internal/model/auth"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/features/g7/fetch"
	"be-gr31/internal/features/g7/content/evaluate"
	"be-gr31/internal/util"
)

// Service handles G7 monthly evaluation and reporting logic
type Service struct {
	fetchService *fetch.Service
}

// NewService creates a new summary Service instance
func NewService(fetchService *fetch.Service) *Service {
	return &Service{fetchService: fetchService}
}

// EvaluateJurnalBulanan mengevaluasi semua jurnal satu bulan → EvalReport
func (s *Service) EvaluateJurnalBulanan(ctx context.Context, nisn, bulan string) (*evaluate.EvalReport, error) {
	filter := g7model.G7Filter{NISN: nisn, BulanDari: bulan, BulanKe: bulan}
	fetcher := util.PagedFetcher[g7model.G7](func(ctx context.Context, size int, state string) ([]g7model.G7, string, error) {
		return s.fetchService.FetchJournalsPaged(ctx, filter, size, state)
	})
	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, err
	}
	report := evaluate.Jurnals(all)
	return &report, nil
}

// Statistik aggregates monthly G7 recaps to calculate class performance metrics and distribution of grades
func (s *Service) Statistik(ctx context.Context, bulan, kelas string, isGuruWali bool, asuhan []authmodel.Siswa, adminID, adminRole string) (*g7model.G7RekapStatistik, error) {
	var allowedNISN map[string]bool
	if isGuruWali {
		allowedNISN = make(map[string]bool)
		for _, st := range asuhan {
			if st.NIS != "" {
				allowedNISN[st.NIS] = true
			}
		}
	}

	all, err := s.fetchService.FetchAllRekap(ctx, g7model.G7RekapFilter{BulanTahun: bulan, Kelas: kelas})
	if err != nil {
		return nil, err
	}

	if isGuruWali {
		filtered := make([]g7model.G7Rekap, 0)
		for _, r := range all {
			if allowedNISN[r.NISN] {
				filtered = append(filtered, r)
			}
		}
		all = filtered
	}

	stat := &g7model.G7RekapStatistik{
		Kelas:      kelas,
		BulanTahun: bulan,
		DistribusiPredikat: map[string]int{
			"Istimewa": 0, "Sangat Baik": 0, "Baik": 0, "Cukup": 0, "Kurang": 0,
		},
		RataRataPerIndikator: map[string]float64{},
	}
	stat.SudahDinilai = len(all)

	if isGuruWali {
		stat.TotalSiswa = len(allowedNISN)
	} else if kelas != "" {
		stat.TotalSiswa = s.fetchService.CountSiswaKelas(ctx, kelas)
	}
	if stat.TotalSiswa < stat.SudahDinilai {
		stat.TotalSiswa = stat.SudahDinilai
	}
	stat.BelumDinilai = stat.TotalSiswa - stat.SudahDinilai

	if len(all) == 0 {
		return stat, nil
	}

	var sumNilai float64
	indikatorSum := map[string]int{}
	var tertinggi, terendah *g7model.G7RekapRingkas
	for i := range all {
		r := all[i]
		stat.DistribusiPredikat[r.Predikat]++
		sumNilai += r.NilaiAkhir
		for k, v := range r.Skor.AsMap() {
			indikatorSum[k] += v
		}
		ringkas := &g7model.G7RekapRingkas{NISN: r.NISN, Nama: r.NamaSiswa, NilaiAkhir: r.NilaiAkhir}
		if tertinggi == nil || r.NilaiAkhir > tertinggi.NilaiAkhir {
			tertinggi = ringkas
		}
		if terendah == nil || r.NilaiAkhir < terendah.NilaiAkhir {
			terendah = ringkas
		}
	}

	n := float64(len(all))
	stat.RataRataNilaiAkhir = round2(sumNilai / n)
	stat.NilaiTertinggi = tertinggi
	stat.NilaiTerendah = terendah
	for k, v := range indikatorSum {
		stat.RataRataPerIndikator[k] = round2(float64(v) / n)
	}
	return stat, nil
}

// BuildLaporanPDF merakit data untuk cetak laporan G7 siswa
func (s *Service) BuildLaporanPDF(ctx context.Context, nisn, bulan string) (*g7model.PDFLaporan, string, error) {
	student, sErr := s.fetchService.FindStudentByNISN(ctx, nisn)
	if sErr != nil {
		return nil, "", sErr
	}
	if student == nil {
		return nil, "", errors.New("siswa tidak ditemukan")
	}

	rekap, err := s.fetchService.GetRekap(ctx, nisn, bulan)
	if err != nil {
		return nil, "", err
	}
	if rekap == nil {
		report, rErr := s.EvaluateJurnalBulanan(ctx, nisn, bulan)
		if rErr != nil {
			return nil, "", rErr
		}

		hariTercatat := s.fetchService.CountHariTercatat(ctx, nisn, bulan)
		skor := evaluate.ToSkorG7(*report)
		
		perolehan, maks, akhir, predikat := g7model.HitungNilaiAkhir(skor, bulan)

		rekap = &g7model.G7Rekap{
			NISN:           nisn,
			NamaSiswa:      student.Nama,
			Kelas:          student.Kelas,
			BulanTahun:     bulan,
			HariTercatat:   hariTercatat,
			Skor:           skor,
			NilaiMaks:      maks,
			NilaiPerolehan: perolehan,
			NilaiAkhir:     akhir,
			Predikat:       predikat,
			Status:         "draft",
		}
	}

	if student.WaliKelas != "" {
		rekap.WaliKelas = student.WaliKelas
	}

	report, err := s.EvaluateJurnalBulanan(ctx, nisn, bulan)
	if err != nil {
		return nil, "", err
	}
	laporan := BuildPDFLaporan(rekap, *report)
	html, err := RenderHTMLLaporan(laporan)
	return laporan, html, err
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}
