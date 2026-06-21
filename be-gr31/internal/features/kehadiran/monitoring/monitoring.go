package monitoring

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"be-gr31/internal/features/kehadiran/fetch/izin"
	"be-gr31/internal/features/kehadiran/fetch/magang"
	authmodel "be-gr31/internal/model/auth"
	kehadiranmodel "be-gr31/internal/model/kehadiran"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/util"

	"github.com/redis/go-redis/v9"
)

type Repo interface {
	ListPaged(ctx context.Context, filter kehadiranmodel.KehadiranFilter, limit int, state string) ([]kehadiranmodel.Kehadiran, string, error)
}

type RosterFetcher interface {
	FetchRoster(ctx context.Context, kelas string) []authmodel.Siswa
}

type Service struct {
	repo     Repo
	rekapSvc RosterFetcher
	rdb      *redis.Client
}

func NewService(repo Repo, rekapSvc RosterFetcher, rdb *redis.Client) *Service {
	return &Service{
		repo:     repo,
		rekapSvc: rekapSvc,
		rdb:      rdb,
	}
}

const cacheTTLKehadiran = 2 * time.Minute

func (s *Service) KehadiranCacheKey(f kehadiranmodel.KehadiranFilter) string {
	raw := fmt.Sprintf("kehadiran:%s:%s:%s:%s:%s:%s:%s",
		f.NIS, f.Kelas, f.Jurusan, f.Tanggal, f.BulanDari, f.BulanKe, f.Status)
	hash := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("cache:kehadiran:list:%x", hash[:8])
}

func (s *Service) InvalidateKehadiranListCache(ctx context.Context, tanggal string) {
	patterns := []string{
		"cache:kehadiran:list:*",
	}
	for _, pattern := range patterns {
		keys, err := s.rdb.Keys(ctx, pattern).Result()
		if err != nil {
			log.Printf("WARNING: gagal scan cache keys %s: %v", pattern, err)
			continue
		}
		if len(keys) > 0 {
			if err := s.rdb.Del(ctx, keys...).Err(); err != nil {
				log.Printf("WARNING: gagal invalidasi cache kehadiran: %v", err)
			}
		}
	}
}

// Status
func ClassifyStatus(status, alasan, fotoIzin string) string {
	if status == "hadir" {
		return status
	}
	if magang.IsMagang(alasan) {
		return "magang"
	}
	if status != "tidak_hadir" {
		return status
	}
	if alasan == "" && fotoIzin == "" {
		return "tidak_hadir"
	}
	if izin.IsSakit(alasan) {
		return "sakit"
	}
	if izin.IsIzin(alasan) {
		return "izin"
	}
	return "izin"
}

// Mengambil daftar kehadiran dengan filter + fuzzy + pagination.
func (s *Service) List(ctx context.Context, filter kehadiranmodel.KehadiranFilter) ([]kehadiranmodel.Kehadiran, bool, int, error) {
	cacheKey := s.KehadiranCacheKey(filter)
	var cachedAll []kehadiranmodel.Kehadiran
	if cached, err := s.rdb.Get(ctx, cacheKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &cachedAll); err == nil {
			all := cachedAll
			if filter.Query != "" && filter.Tanggal == "" {
				all = util.FuzzyFilter(all, filter.Query, func(k kehadiranmodel.Kehadiran) []string {
					return []string{k.NamaSiswa, k.NIS, k.Kelas}
				})
			}
			result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
			return result, hasMore, total, nil
		}
	}

	var all []kehadiranmodel.Kehadiran

	// Case 1: Daily Monitoring (filter.Tanggal is specified)
	if filter.Tanggal != "" {
		// 1. Fetch all students roster
		var allStudents []authmodel.Siswa
		if s.rekapSvc != nil {
			allStudents = s.rekapSvc.FetchRoster(ctx, "")
		}

		// 2. Filter students roster based on grade, major, NISN, search query
		var filteredStudents []authmodel.Siswa
		for _, student := range allStudents {
			if filter.NIS != "" && student.NIS != filter.NIS {
				continue
			}

			// Kelas filter
			if filter.Kelas != "" {
				if strings.Contains(filter.Kelas, " ") {
					if student.Kelas != filter.Kelas {
						continue
					}
				} else {
					parts := strings.Fields(student.Kelas)
					if len(parts) == 0 || parts[0] != filter.Kelas {
						continue
					}
				}
			}

			// Jurusan filter
			if filter.Jurusan != "" {
				parts := strings.Fields(student.Kelas)
				if len(parts) <= 1 {
					continue
				}
				jur := strings.Join(parts[1:], " ")
				if !strings.Contains(strings.ToLower(jur), strings.ToLower(filter.Jurusan)) {
					continue
				}
			}

			// Fuzzy search query (q) filter
			if filter.Query != "" {
				q := strings.ToLower(filter.Query)
				matchName := strings.Contains(strings.ToLower(student.Nama), q)
				matchNisn := strings.Contains(strings.ToLower(student.NIS), q)
				matchKelas := strings.Contains(strings.ToLower(student.Kelas), q)
				if !matchName && !matchNisn && !matchKelas {
					continue
				}
			}

			filteredStudents = append(filteredStudents, student)
		}

		// 3. Fetch all attendance records for this date
		repoFilter := kehadiranmodel.KehadiranFilter{
			Tanggal: filter.Tanggal,
		}
		fetcher := util.PagedFetcher[kehadiranmodel.Kehadiran](func(ctx context.Context, size int, state string) ([]kehadiranmodel.Kehadiran, string, error) {
			return s.repo.ListPaged(ctx, repoFilter, size, state)
		})
		attendanceRecords, err := util.FetchAll(ctx, fetcher, 100)
		if err != nil {
			return nil, false, 0, err
		}

		// Index attendance by NIS
		attendanceByNIS := make(map[string]kehadiranmodel.Kehadiran)
		for _, rec := range attendanceRecords {
			attendanceByNIS[rec.NIS] = rec
		}

		// 4. Join roster and attendance
		var resultRecords []kehadiranmodel.Kehadiran
		for _, student := range filteredStudents {
			rec, found := attendanceByNIS[student.NIS]
			if found {
				rec.NamaSiswa = student.Nama
				rec.Kelas = student.Kelas
				rec.Status = ClassifyStatus(rec.Status, rec.Alasan, rec.FotoIzin)
				if filter.Status != "" {
					if filter.Status == "izin_sakit" {
						if rec.Status != "izin" && rec.Status != "sakit" {
							continue
						}
					} else {
						if rec.Status != filter.Status {
							continue
						}
					}
				}
				resultRecords = append(resultRecords, rec)
			} else {
				if filter.Status != "" && filter.Status != "belum" {
					continue
				}
				resultRecords = append(resultRecords, kehadiranmodel.Kehadiran{
					ID:        "belum_" + student.NIS,
					NIS:       student.NIS,
					NamaSiswa: student.Nama,
					Kelas:     student.Kelas,
					Tanggal:   filter.Tanggal,
					Status:    "belum",
				})
			}
		}

		// Sort: by class, then by student name
		sort.Slice(resultRecords, func(i, j int) bool {
			if resultRecords[i].Kelas != resultRecords[j].Kelas {
				return resultRecords[i].Kelas < resultRecords[j].Kelas
			}
			return resultRecords[i].NamaSiswa < resultRecords[j].NamaSiswa
		})

		all = resultRecords

	} else {
		// Case 2: Monthly/Range Query (no single date)
		isExactClass := false
		if filter.Kelas != "" && strings.Contains(filter.Kelas, " ") {
			isExactClass = true
		}

		memoryFilter := false
		gradeFilter := ""
		jurusanFilter := ""

		repoFilter := filter

		if filter.Kelas != "" && !strings.Contains(filter.Kelas, " ") {
			gradeFilter = filter.Kelas
			repoFilter.Kelas = ""
			memoryFilter = true
		}
		if filter.Jurusan != "" {
			jurusanFilter = filter.Jurusan
			if !isExactClass {
				repoFilter.Kelas = ""
				memoryFilter = true
			}
		}
		if filter.Status == "izin_sakit" {
			repoFilter.Status = ""
			memoryFilter = true
		}

		fetcher := util.PagedFetcher[kehadiranmodel.Kehadiran](func(ctx context.Context, size int, state string) ([]kehadiranmodel.Kehadiran, string, error) {
			return s.repo.ListPaged(ctx, repoFilter, size, state)
		})

		var err error
		all, err = util.FetchAll(ctx, fetcher, 100)
		if err != nil {
			return nil, false, 0, err
		}

		var allStudents []authmodel.Siswa
		if s.rekapSvc != nil {
			allStudents = s.rekapSvc.FetchRoster(ctx, "")
		}
		studentMap := make(map[string]authmodel.Siswa)
		for _, student := range allStudents {
			studentMap[student.NIS] = student
		}

		for i := range all {
			all[i].Status = ClassifyStatus(all[i].Status, all[i].Alasan, all[i].FotoIzin)
			if student, found := studentMap[all[i].NIS]; found {
				if all[i].NamaSiswa == "" {
					all[i].NamaSiswa = student.Nama
				}
				if all[i].Kelas == "" || all[i].Kelas == "Tanpa Kelas" {
					all[i].Kelas = student.Kelas
				}
			}
		}

		if memoryFilter {
			filtered := make([]kehadiranmodel.Kehadiran, 0)
			for _, k := range all {
				matchGrade := true
				matchJurusan := true
				matchStatus := true

				if gradeFilter != "" {
					parts := strings.Fields(k.Kelas)
					if len(parts) == 0 || parts[0] != gradeFilter {
						matchGrade = false
					}
				}
				if jurusanFilter != "" {
					parts := strings.Fields(k.Kelas)
					if len(parts) <= 1 {
						matchJurusan = false
					} else {
						jur := strings.Join(parts[1:], " ")
						if !strings.Contains(strings.ToLower(jur), strings.ToLower(jurusanFilter)) {
							matchJurusan = false
						}
					}
				}
				if filter.Status == "izin_sakit" {
					if k.Status != "izin" && k.Status != "sakit" {
						matchStatus = false
					}
				}

				if matchGrade && matchJurusan && matchStatus {
					filtered = append(filtered, k)
				}
			}
			all = filtered
		}

		// Fuzzy search by name/NIS/class
		if filter.Query != "" {
			all = util.FuzzyFilter(all, filter.Query, func(k kehadiranmodel.Kehadiran) []string {
				return []string{k.NamaSiswa, k.NIS, k.Kelas}
			})
		}
	}

	// Cache the result (for 2 mins)
	if data, err := json.Marshal(all); err == nil {
		if setErr := s.rdb.Set(ctx, cacheKey, data, cacheTTLKehadiran).Err(); setErr != nil {
			log.Printf("WARNING: gagal cache kehadiran list: %v", setErr)
		}
	}

	result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}

// BuildSummaryByStudent
func BuildSummaryByStudent(items []kehadiranmodel.Kehadiran) []rekapmodel.SummaryPerSiswa {
	type agg struct {
		nama          string
		kelas         string
		h, i, s, a, m int
	}
	grouped := map[string]*agg{}
	for _, k := range items {
		entry, ok := grouped[k.NIS]
		if !ok {
			entry = &agg{nama: k.NamaSiswa, kelas: k.Kelas}
			grouped[k.NIS] = entry
		}
		switch k.Status {
		case "hadir":
			entry.h++
		case "izin":
			entry.i++
		case "sakit":
			entry.s++
		case "magang":
			entry.m++
		default:
			entry.a++
		}
	}
	result := make([]rekapmodel.SummaryPerSiswa, 0, len(grouped))
	for nis, e := range grouped {
		total := e.h + e.i + e.s + e.a + e.m
		persen := 0.0
		if total > 0 {
			persen = float64(e.h) / float64(total) * 100
		}
		result = append(result, rekapmodel.SummaryPerSiswa{
			NIS:             nis,
			NamaSiswa:       e.nama,
			Kelas:           e.kelas,
			TotalHadir:      e.h,
			TotalIzin:       e.i,
			TotalSakit:      e.s,
			TotalAlpa:       e.a,
			TotalMagang:     e.m,
			PersentaseHadir: persen,
		})
	}
	return result
}
