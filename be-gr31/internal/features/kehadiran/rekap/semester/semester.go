package semester

import (
	"context"
	"fmt"
	"time"

	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/util"
)

type RekapBulananLister interface {
	ListRekap(ctx context.Context, filter rekapmodel.RekapFilter) ([]rekapmodel.RekapBulanan, bool, int, error)
}

type Service struct {
	rekapLister RekapBulananLister
}

func NewService(rekapLister RekapBulananLister) *Service {
	return &Service{
		rekapLister: rekapLister,
	}
}

// GetRekapSemesterKelas mengambil tren kehadiran kelas selama satu semester.
func (s *Service) GetRekapSemesterKelas(ctx context.Context, kelas, semester string) (*rekapmodel.RekapSemesterKelas, error) {
	bulanList := bulanPerSemester(semester)

	bulanResult := make([]rekapmodel.RekapSemesterItem, 0, len(bulanList))
	for _, bulan := range bulanList {
		filter := rekapmodel.RekapFilter{
			Kelas:      kelas,
			BulanTahun: bulan,
			Page:       1,
			Limit:      util.MaxFetchAll,
		}
		rekaps, _, _, err := s.rekapLister.ListRekap(ctx, filter)
		if err != nil || len(rekaps) == 0 {
			continue
		}

		var h, i, sa, a, m, total int
		for _, r := range rekaps {
			h += r.TotalHadir
			i += r.TotalIzin
			sa += r.TotalSakit
			a += r.TotalTidakHadir
			m += r.TotalMagang
			total += r.TotalHariEfektif
		}
		persen := 0.0
		if total > 0 {
			persen = float64(h) / float64(total) * 100
		}
		bulanResult = append(bulanResult, rekapmodel.RekapSemesterItem{
			BulanTahun:      bulan,
			TotalHadir:      h,
			TotalIzin:       i,
			TotalSakit:      sa,
			TotalAlpa:       a,
			TotalMagang:     m,
			PersentaseHadir: persen,
		})
	}

	return &rekapmodel.RekapSemesterKelas{
		Kelas:    kelas,
		Semester: semester,
		Bulan:    bulanResult,
	}, nil
}

// bulanPerSemester mengembalikan daftar YYYY-MM untuk satu semester
func bulanPerSemester(semester string) []string {
	now := time.Now()
	tahun := now.Year()

	if len(semester) > 6 {
		if semester[:6] == "ganjil" {
			if len(semester) == 11 {
				fmt.Sscanf(semester[7:], "%d", &tahun)
			}
			from := fmt.Sprintf("%d-07", tahun)
			to := fmt.Sprintf("%d-12", tahun)
			return bulanDalamRentang(from, to)
		}
		if semester[:5] == "genap" {
			if len(semester) == 10 {
				fmt.Sscanf(semester[6:], "%d", &tahun)
			}
			from := fmt.Sprintf("%d-01", tahun)
			to := fmt.Sprintf("%d-06", tahun)
			return bulanDalamRentang(from, to)
		}
	}

	bulan := int(now.Month())
	if bulan >= 7 {
		from := fmt.Sprintf("%d-07", tahun)
		to := fmt.Sprintf("%d-12", tahun)
		return bulanDalamRentang(from, to)
	}
	from := fmt.Sprintf("%d-01", tahun)
	to := fmt.Sprintf("%d-06", tahun)
	return bulanDalamRentang(from, to)
}

func bulanDalamRentang(dari, sampai string) []string {
	if len(dari) < 7 || len(sampai) < 7 {
		return nil
	}
	var result []string
	current := dari[:7]
	end := sampai[:7]
	for current <= end {
		result = append(result, current)
		t, err := time.Parse("2006-01", current)
		if err != nil {
			break
		}
		current = t.AddDate(0, 1, 0).Format("2006-01")
	}
	return result
}
