package semester

import (
	"context"
	"errors"
	"math"
	"sort"
	"strings"

	authmodel "be-gr31/internal/model/auth"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/features/g7/fetch"
)

// Service handles G7 semester recap operations
type Service struct {
	fetchService *fetch.Service
}

// NewService creates a new Service instance
func NewService(fetchService *fetch.Service) *Service {
	return &Service{fetchService: fetchService}
}

// GetRekapSemesterKelas aggregates and returns G7 semester averages per student
func (s *Service) GetRekapSemesterKelas(ctx context.Context, semester, kelas string, isGuruWali bool, asuhan []authmodel.Siswa) ([]g7model.G7SemesterStudentItem, error) {
	bulanList := bulanPerSemester(semester)
	if len(bulanList) == 0 {
		return nil, errors.New("semester format tidak valid")
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

	type studentAgg struct {
		NISN            string
		NamaSiswa       string
		Kelas           string
		TotalNilaiAkhir float64
		MonthsCount     int
		Grades          map[string][]int
	}

	studentMap := make(map[string]*studentAgg)

	for _, m := range bulanList {
		filter := g7model.G7RekapFilter{
			BulanTahun: m,
			Kelas:      kelas,
		}

		rekaps, err := s.fetchService.FetchAllRekap(ctx, filter)
		if err != nil {
			continue
		}

		for _, r := range rekaps {
			if isGuruWali && !allowedNISN[r.NISN] {
				continue
			}

			sa, ok := studentMap[r.NISN]
			if !ok {
				sa = &studentAgg{
					NISN:      r.NISN,
					NamaSiswa: r.NamaSiswa,
					Kelas:     r.Kelas,
					Grades:    make(map[string][]int),
				}
				studentMap[r.NISN] = sa
			}

			sa.TotalNilaiAkhir += r.NilaiAkhir
			sa.MonthsCount++

			for key, val := range r.Skor.AsMap() {
				sa.Grades[key] = append(sa.Grades[key], val)
			}
		}
	}

	var list []g7model.G7SemesterStudentItem
	for _, sAgg := range studentMap {
		avgNilai := 0.0
		if sAgg.MonthsCount > 0 {
			avgNilai = math.Round((sAgg.TotalNilaiAkhir/float64(sAgg.MonthsCount))*100) / 100
		}

		averagedSkor := make(map[string]float64)
		for key, arr := range sAgg.Grades {
			sum := 0
			for _, v := range arr {
				sum += v
			}
			avgVal := 0.0
			if len(arr) > 0 {
				avgVal = math.Round((float64(sum)/float64(len(arr)))*10) / 10
			}
			averagedSkor[key] = avgVal
		}

		list = append(list, g7model.G7SemesterStudentItem{
			NISN:        sAgg.NISN,
			NamaSiswa:   sAgg.NamaSiswa,
			Kelas:       sAgg.Kelas,
			NilaiAkhir:  avgNilai,
			Predikat:    g7model.Predikat(avgNilai),
			MonthsCount: sAgg.MonthsCount,
			Skor:        averagedSkor,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].Kelas != list[j].Kelas {
			return list[i].Kelas < list[j].Kelas
		}
		return list[i].NamaSiswa < list[j].NamaSiswa
	})

	return list, nil
}

func bulanPerSemester(semester string) []string {
	if strings.HasSuffix(semester, "-genap") {
		parts := strings.Split(semester, "-")
		if len(parts) > 0 {
			years := strings.Split(parts[0], "/")
			if len(years) > 1 {
				year := years[1]
				return []string{
					year + "-01", year + "-02", year + "-03",
					year + "-04", year + "-05", year + "-06",
				}
			}
		}
	} else if strings.HasSuffix(semester, "-ganjil") {
		parts := strings.Split(semester, "-")
		if len(parts) > 0 {
			years := strings.Split(parts[0], "/")
			if len(years) > 0 {
				year := years[0]
				return []string{
					year + "-07", year + "-08", year + "-09",
					year + "-10", year + "-11", year + "-12",
				}
			}
		}
	}
	return nil
}
