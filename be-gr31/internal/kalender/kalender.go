package kalender

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Libur struct {
	Tanggal        string `json:"tanggal,omitempty"`
	TanggalMulai   string `json:"tanggal_mulai,omitempty"`
	TanggalSelesai string `json:"tanggal_selesai,omitempty"`
	Keterangan     string `json:"keterangan"`
}

// BulanInfo
type BulanInfo struct {
	Bulan             string  `json:"bulan"`
	Tahun             int     `json:"tahun"`
	HariKerjaTersedia int     `json:"hari_kerja_tersedia"`
	HariMasukEfektif  int     `json:"hari_masuk_efektif"`
	DaftarLibur       []Libur `json:"daftar_libur"`
	KegiatanAkademik  string  `json:"kegiatan_akademik"`
	BulanTahun        string  `json:"bulanTahun"`
}

type kalenderFile struct {
	Wilayah      string      `json:"wilayah"`
	TahunAjaran  string      `json:"tahun_ajaran"`
	Semester     string      `json:"semester"`
	RincianBulan []BulanInfo `json:"rincian_bulan"`
}

var namaBulan = map[string]int{
	"januari": 1, "februari": 2, "maret": 3, "april": 4,
	"mei": 5, "juni": 6, "juli": 7, "agustus": 8,
	"september": 9, "oktober": 10, "november": 11, "desember": 12,
}

type store struct {
	byBulan  map[string]BulanInfo
	liburSet map[string]struct{}
	loaded   bool
}

var (
	once   sync.Once
	cached *store
)

func candidatePaths() []string {
	paths := []string{}
	if p := strings.TrimSpace(os.Getenv("KALENDER_PENDIDIKAN_PATH")); p != "" {
		paths = append(paths, p)
	}
	paths = append(paths,
		filepath.Join("data", "kalender-pendidikan.json"),
		filepath.Join(".", "data", "kalender-pendidikan.json"),
	)
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		paths = append(paths,
			filepath.Join(dir, "data", "kalender-pendidikan.json"),
			filepath.Join(dir, "..", "data", "kalender-pendidikan.json"),
		)
	}
	return paths
}

func load() *store {
	once.Do(func() {
		cached = &store{
			byBulan:  map[string]BulanInfo{},
			liburSet: map[string]struct{}{},
		}

		var raw []byte
		var err error
		for _, p := range candidatePaths() {
			raw, err = os.ReadFile(p)
			if err == nil {
				break
			}
		}
		if err != nil || len(raw) == 0 {
			return
		}

		var kf kalenderFile
		if err := json.Unmarshal(raw, &kf); err != nil {
			return
		}

		for _, b := range kf.RincianBulan {
			mn := namaBulan[strings.ToLower(strings.TrimSpace(b.Bulan))]
			if mn == 0 || b.Tahun == 0 {
				continue
			}
			key := monthKey(b.Tahun, mn)
			b.BulanTahun = key
			cached.byBulan[key] = b

			for _, d := range expandLibur(b.DaftarLibur) {
				cached.liburSet[d] = struct{}{}
			}
		}
		cached.loaded = true
	})
	return cached
}

func monthKey(tahun, bulan int) string {
	return time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC).Format("2006-01")
}

func expandLibur(list []Libur) []string {
	out := []string{}
	for _, l := range list {
		if t := strings.TrimSpace(l.Tanggal); t != "" {
			out = append(out, t)
			continue
		}
		mulai := strings.TrimSpace(l.TanggalMulai)
		selesai := strings.TrimSpace(l.TanggalSelesai)
		if mulai == "" || selesai == "" {
			continue
		}
		start, err1 := time.Parse("2006-01-02", mulai)
		end, err2 := time.Parse("2006-01-02", selesai)
		if err1 != nil || err2 != nil || end.Before(start) {
			continue
		}
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			out = append(out, d.Format("2006-01-02"))
		}
	}
	return out
}

// Info mengembalikan
func Info(bulan string) (BulanInfo, bool) {
	s := load()
	b, ok := s.byBulan[normalizeBulan(bulan)]
	return b, ok
}

// HariEfektif
func HariEfektif(bulan string) int {
	if b, ok := Info(bulan); ok && b.HariMasukEfektif > 0 {
		return b.HariMasukEfektif
	}
	return countWeekdaysInMonth(bulan)
}

// DaftarLibur
func DaftarLibur(bulan string) []Libur {
	if b, ok := Info(bulan); ok {
		return b.DaftarLibur
	}
	return nil
}

// IsLibur
func IsLibur(tanggal string) bool {
	s := load()
	_, ok := s.liburSet[strings.TrimSpace(tanggal)]
	return ok
}

// IsHariEfektif
func IsHariEfektif(tanggal string) bool {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(tanggal))
	if err != nil {
		return false
	}
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return false
	}
	return !IsLibur(tanggal)
}

func normalizeBulan(bulan string) string {
	bulan = strings.TrimSpace(bulan)
	if len(bulan) >= 7 {
		return bulan[:7]
	}
	return bulan
}

func countWeekdaysInMonth(bulan string) int {
	t, err := time.Parse("2006-01", normalizeBulan(bulan))
	if err != nil {
		return 0
	}
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			count++
		}
	}
	return count
}
