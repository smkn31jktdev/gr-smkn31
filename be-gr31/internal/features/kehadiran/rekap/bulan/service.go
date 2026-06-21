package bulan

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	kehadiranmodel "be-gr31/internal/model/kehadiran"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/util"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type KehadiranStore interface {
	ListPaged(ctx context.Context, filter kehadiranmodel.KehadiranFilter, pageSize int, pageState string) ([]kehadiranmodel.Kehadiran, string, error)
}

type Service struct {
	store          *astra.RekapStore
	kehadiranSvc   KehadiranStore
	studentStore   *astra.StudentStore
	rdb            *redis.Client
	supabaseClient *supabase.Client
}

func NewService(store *astra.RekapStore, kehadiranSvc KehadiranStore, studentStore *astra.StudentStore, rdb *redis.Client, supabaseClient *supabase.Client) *Service {
	return &Service{
		store:          store,
		kehadiranSvc:   kehadiranSvc,
		studentStore:   studentStore,
		rdb:            rdb,
		supabaseClient: supabaseClient,
	}
}

func (s *Service) SetKehadiranStore(ks KehadiranStore) {
	s.kehadiranSvc = ks
}

func (s *Service) SetStudentStore(ss *astra.StudentStore) {
	s.studentStore = ss
}

// IncrementRekap melakukan increment counter setelah absensi baru dibuat.
func (s *Service) IncrementRekap(ctx context.Context, nis, namaSiswa, kelas, bulanTahun, status string) error {
	return s.mutatRekap(ctx, nis, namaSiswa, kelas, bulanTahun, status, 1)
}

// DecrementRekap melakukan decrement setelah absensi dihapus.
func (s *Service) DecrementRekap(ctx context.Context, nis, namaSiswa, kelas, bulanTahun, status string) error {
	return s.mutatRekap(ctx, nis, namaSiswa, kelas, bulanTahun, status, -1)
}

func (s *Service) mutatRekap(ctx context.Context, nis, namaSiswa, kelas, bulanTahun, status string, delta int) error {
	rekapKey := nis + "_" + bulanTahun
	now := time.Now().Format(time.RFC3339)

	// Hitung semester dari bulan (format YYYY-MM)
	semester := "genap"
	if len(bulanTahun) >= 7 {
		bulanInt, _ := strconv.Atoi(bulanTahun[5:7])
		semester = rekapmodel.HitungSemester(bulanInt)
	}

	// Cek apakah dokumen rekap sudah ada
	existing, err := s.store.FindByKey(ctx, rekapKey)
	if err != nil {
		return err
	}

	if existing == nil {
		newRekap := rekapmodel.RekapBulanan{
			ID:         uuid.New().String(),
			RekapKey:   rekapKey,
			NIS:        nis,
			NamaSiswa:  namaSiswa,
			Kelas:      kelas,
			BulanTahun: bulanTahun,
			Semester:   semester,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := s.store.Upsert(ctx, &newRekap); err != nil {
			return err
		}
	}

	incMap := buildIncMap(status, delta)

	// Atomic increment via $inc
	if err := s.store.IncrementCounters(ctx, rekapKey, incMap, now); err != nil {
		return err
	}

	s.invalidateCache(ctx, nis, bulanTahun)
	return nil
}

func buildIncMap(status string, delta int) map[string]int {
	m := map[string]int{
		"total_hari_efektif": delta,
	}
	switch status {
	case "hadir":
		m["total_hadir"] = delta
	case "izin":
		m["total_izin"] = delta
	case "sakit":
		m["total_sakit"] = delta
	case "magang":
		m["total_magang"] = delta
	default: // tidak_hadir / alpa
		m["total_tidak_hadir"] = delta
	}
	return m
}

func (s *Service) GetRekap(ctx context.Context, nis, bulanTahun string) (*rekapmodel.RekapBulanan, error) {
	return s.store.FindByKey(ctx, nis+"_"+bulanTahun)
}

const cacheTTLRekap = 2 * time.Minute

func (s *Service) ListRekap(ctx context.Context, filter rekapmodel.RekapFilter) ([]rekapmodel.RekapBulanan, bool, int, error) {
	var all []rekapmodel.RekapBulanan

	cacheKey := s.rekapCacheKey(filter)
	if cached, err := s.rdb.Get(ctx, cacheKey).Bytes(); err == nil {
		if err := json.Unmarshal(cached, &all); err == nil {
			if filter.Query != "" {
				all = util.FuzzyFilter(all, filter.Query, func(r rekapmodel.RekapBulanan) []string {
					return []string{r.NamaSiswa, r.NIS, r.Kelas}
				})
			}
			result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
			return result, hasMore, total, nil
		}
	}

	fetcher := util.PagedFetcher[rekapmodel.RekapBulanan](func(ctx context.Context, size int, state string) ([]rekapmodel.RekapBulanan, string, error) {
		return s.store.ListPaged(ctx, filter, size, state)
	})

	var err error
	all, err = util.FetchAll(ctx, fetcher, util.MaxPageSize)
	if err != nil {
		return nil, false, 0, err
	}

	if data, err := json.Marshal(all); err == nil {
		if setErr := s.rdb.Set(ctx, cacheKey, data, cacheTTLRekap).Err(); setErr != nil {
			log.Printf("WARNING: gagal cache rekap list: %v", setErr)
		}
	}

	if filter.Query != "" {
		all = util.FuzzyFilter(all, filter.Query, func(r rekapmodel.RekapBulanan) []string {
			return []string{r.NamaSiswa, r.NIS, r.Kelas}
		})
	}

	result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}

func (s *Service) rekapCacheKey(f rekapmodel.RekapFilter) string {
	raw := fmt.Sprintf("rekap:%s:%s:%s:%s",
		f.NIS, f.Kelas, f.BulanTahun, f.Semester)
	hash := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("cache:rekap:list:%x", hash[:8])
}

func (s *Service) invalidateCache(ctx context.Context, nis, bulanTahun string) {
	keys := []string{
		"cache:rekap:" + nis + ":" + bulanTahun,
		"cache:rekap:kelas:" + bulanTahun,
	}
	for _, k := range keys {
		if err := s.rdb.Del(ctx, k).Err(); err != nil {
			log.Printf("WARNING: gagal invalidasi cache %s: %v", k, err)
		}
	}

	listKeys, err := s.rdb.Keys(ctx, "cache:rekap:list:*").Result()
	if err == nil && len(listKeys) > 0 {
		if err := s.rdb.Del(ctx, listKeys...).Err(); err != nil {
			log.Printf("WARNING: gagal invalidasi cache rekap list: %v", err)
		}
	}

	// Invalidasi cache agregasi rekap lengkap & mingguan (semua kelas/periode).
	for _, pattern := range []string{"cache:rekap:lengkap:*", "cache:rekap:mingguan:*"} {
		if ks, err := s.rdb.Keys(ctx, pattern).Result(); err == nil && len(ks) > 0 {
			if err := s.rdb.Del(ctx, ks...).Err(); err != nil {
				log.Printf("WARNING: gagal invalidasi cache %s: %v", pattern, err)
			}
		}
	}
}

// getCachedJSON membaca nilai JSON dari Redis ke dst. Mengembalikan true bila hit.
func (s *Service) getCachedJSON(ctx context.Context, key string, dst any) bool {
	if s.rdb == nil {
		return false
	}
	b, err := s.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(b, dst) == nil
}

// setCachedJSON menyimpan nilai sebagai JSON ke Redis dengan TTL.
func (s *Service) setCachedJSON(ctx context.Context, key string, val any, ttl time.Duration) {
	if s.rdb == nil {
		return
	}
	if b, err := json.Marshal(val); err == nil {
		if err := s.rdb.Set(ctx, key, b, ttl).Err(); err != nil {
			log.Printf("WARNING: gagal set cache %s: %v", key, err)
		}
	}
}

func BulanDalamRentang(dari, sampai string) []string {
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
