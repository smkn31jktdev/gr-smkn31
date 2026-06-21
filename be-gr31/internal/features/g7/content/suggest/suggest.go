package suggest

import (
	"context"
	"fmt"
	"time"

	"be-gr31/internal/features/g7/content/evaluate"
	"be-gr31/internal/features/g7/content/score"
	"be-gr31/internal/features/g7/fetch"
	"be-gr31/internal/kalender/puasa"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/util"
)

// Service handles G7 score suggestion logic
type Service struct {
	fetchService *fetch.Service
}

// NewService creates a new Service instance
func NewService(fetchService *fetch.Service) *Service {
	return &Service{fetchService: fetchService}
}

// Suggest evaluates journals to generate scoring recommendations
func (s *Service) Suggest(ctx context.Context, nisn, bulan string) (*g7model.G7SuggestResponse, error) {
	if _, err := time.Parse("2006-01", bulan); err != nil {
		return nil, fmt.Errorf("bulan harus YYYY-MM")
	}

	filter := g7model.G7Filter{NISN: nisn, BulanDari: bulan, BulanKe: bulan}
	fetcher := util.PagedFetcher[g7model.G7](func(ctx context.Context, size int, state string) ([]g7model.G7, string, error) {
		return s.fetchService.FetchJournalsPaged(ctx, filter, size, state)
	})
	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, err
	}
	hari := len(all)

	counts := aggregateCounts(all)

	isRamadan, _ := puasa.IsRamadanMonth(bulan)

	suggestedRowatib := 0
	if isRamadan {
		suggestedRowatib = score.ByRatio(counts.rowatib, hari)
	}

	suggestedTarawih := 0
	if isRamadan {
		suggestedTarawih = scoreTarawih(counts.tarawihSum)
	}

	suggestedPuasa := 0
	if isRamadan {
		ramadanDays, _ := puasa.RamadanDaysInMonth(bulan)
		missed := ramadanDays - counts.puasaCount
		suggestedPuasa = scorePuasa(missed)
	}

	skor := g7model.SkorG7{
		BangunPagi:         score.ByMonthlyCount(counts.bangun),
		IbadahDoa:          score.ByMonthlyCount(counts.doa),
		IbadahSholatFajar:  score.SholatFajar(counts.fajar),
		IbadahSholat5Waktu: score.Sholat5(counts.sholat5),
		IbadahZikir:        score.ByRatio(counts.zikir, hari),
		IbadahDhuha:        score.ByRatio(counts.dhuha, hari),
		IbadahRowatib:      suggestedRowatib,
		IbadahTarawih:      suggestedTarawih,
		IbadahPuasa:        suggestedPuasa,
		IbadahZakat:        score.Zakat(counts.zakatSum),
		Olahraga:           score.ByRatio(counts.olahraga, hari),
		MakanSehat:         score.Makan(counts.makanSum, hari),
		BelajarKitabSuci:   score.ByRatio(counts.belajar, hari),
		BelajarBukuUmum:    score.ByRatio(counts.belajar, hari),
		BelajarBukuMapel:   score.ByRatio(counts.belajar, hari),
		BelajarTugas:       score.ByRatio(counts.belajar, hari),
		Bermasyarakat:      score.Bermasyarakat(counts.masyarakatDistinct),
		TidurCepat:         score.ByMonthlyCount(counts.tidur),
	}

	var catatanRowatib string
	var catatanTarawih string
	var catatanPuasa string

	if isRamadan {
		catatanRowatib = fmt.Sprintf("Rowatib tercatat %d kali → skor %d", counts.rowatib, skor.IbadahRowatib)
		catatanTarawih = fmt.Sprintf("Total tarawih %d rokaat → skor %d", counts.tarawihSum, skor.IbadahTarawih)
		ramadanDays, _ := puasa.RamadanDaysInMonth(bulan)
		catatanPuasa = fmt.Sprintf("Puasa %d dari %d hari Ramadhan → skor %d", counts.puasaCount, ramadanDays, skor.IbadahPuasa)
	} else {
		catatanRowatib = "Luar Ramadhan → skor 5 (tidak dinilai)"
		catatanTarawih = "Luar Ramadhan → skor 5 (tidak dinilai)"
		catatanPuasa = "Luar Ramadhan → skor 5 (tidak dinilai)"
	}

	catatan := map[string]string{
		"bangunPagi":    fmt.Sprintf("Bangun ≤05:00 sebanyak %d kali → skor %d", counts.bangun, skor.BangunPagi),
		"tidurCepat":    fmt.Sprintf("Tidur ≤22:00 sebanyak %d kali → skor %d", counts.tidur, skor.TidurCepat),
		"ibadahDoa":     fmt.Sprintf("Berdoa tercatat %d kali → skor %d", counts.doa, skor.IbadahDoa),
		"ibadahZakat":   fmt.Sprintf("Total infaq/sodaqoh Rp%.0f → skor %d", counts.zakatSum, skor.IbadahZakat),
		"makanSehat":    fmt.Sprintf("Rata-rata %.1f jenis makanan/hari → skor %d", safeAvg(counts.makanSum, hari), skor.MakanSehat),
		"bermasyarakat": fmt.Sprintf("%d jenis kegiatan berbeda → skor %d", counts.masyarakatDistinct, skor.Bermasyarakat),
		"belajar":       "Jurnal harian belum memisahkan 4 sub-indikator belajar; skor disamakan dari frekuensi belajar mandiri — mohon disesuaikan manual.",
		"ibadahRowatib": catatanRowatib,
		"ibadahTarawih": catatanTarawih,
		"ibadahPuasa":   catatanPuasa,
	}

	return &g7model.G7SuggestResponse{
		NISN:         nisn,
		BulanTahun:   bulan,
		Skor:         skor,
		Catatan:      catatan,
		HariTercatat: hari,
		IsAdvisory:   true,
	}, nil
}

type suggestCounts struct {
	bangun, tidur, doa                    int
	fajar, sholat5, zikir, dhuha, rowatib int
	olahraga, belajar                     int
	zakatSum, makanSum                    float64
	masyarakatDistinct                    int
	tarawihSum, puasaCount                int
}

func aggregateCounts(all []g7model.G7) suggestCounts {
	jenisMasyarakat := map[string]bool{}
	var c suggestCounts

	for i := range all {
		g := all[i]
		if g.Bangun != nil && g.Bangun.Done && sugTimeLE(g.Bangun.Waktu, "05:00") {
			c.bangun++
		}
		if g.Tidur != nil && g.Tidur.Done && sugTimeLE(g.Tidur.Waktu, "22:00") {
			c.tidur++
		}
		if g.Ibadah != nil {
			ki := sugParseKet(g.Ibadah)
			if sugBool(ki, "berdoa") {
				c.doa++
			}
			if sugBool(ki, "sholatFajar") {
				c.fajar++
			}
			if sugBool(ki, "sholat5Waktu") {
				c.sholat5++
			}
			if sugBool(ki, "zikir") {
				c.zikir++
			}
			if sugBool(ki, "sholatDhuha") {
				c.dhuha++
			}
			if sugBool(ki, "rowatib") || sugBool(ki, "sholatSunah") {
				c.rowatib++
			}
			if sugBool(ki, "puasa") {
				c.puasaCount++
			}
			c.tarawihSum += int(sugFloat(ki, "tarawihRokaat"))
			c.zakatSum += sugFloat(ki, "infaq")
		}
		if g.Olahraga != nil && g.Olahraga.Done {
			c.olahraga++
		}
		if g.Belajar != nil && g.Belajar.Done {
			c.belajar++
		}
		if g.Makan != nil {
			km := sugParseKet(g.Makan)
			n := 0
			if sugStr(km, "makanUtama") != "" {
				n++
			}
			if sugStr(km, "laukPauk") != "" {
				n++
			}
			if sugBool(km, "sayurBuah") {
				n++
			}
			if sugBool(km, "susuSuplemen") {
				n++
			}
			c.makanSum += float64(n)
		}
		if g.Bermasyarakat != nil && g.Bermasyarakat.Done {
			if j := sugStr(sugParseKet(g.Bermasyarakat), "jenis"); j != "" {
				jenisMasyarakat[j] = true
			}
		}
	}
	c.masyarakatDistinct = len(jenisMasyarakat)
	return c
}

func sugParseKet(a *g7model.Aktivitas) map[string]any {
	return evaluate.ParseKet(a)
}

func sugBool(m map[string]any, key string) bool  { b, _ := m[key].(bool); return b }
func sugStr(m map[string]any, key string) string { s, _ := m[key].(string); return s }
func sugFloat(m map[string]any, key string) float64 {
	switch t := m[key].(type) {
	case float64:
		return t
	case int:
		return float64(t)
	}
	return 0
}
func sugTimeLE(a, b string) bool { return a != "" && a <= b }
func safeAvg(sum float64, n int) float64 {
	if n == 0 {
		return 0
	}
	return sum / float64(n)
}

func scoreTarawih(rokaat int) int {
	switch {
	case rokaat > 330:
		return 5
	case rokaat >= 308:
		return 4
	case rokaat >= 275:
		return 3
	case rokaat >= 231:
		return 2
	default:
		return 1
	}
}

func scorePuasa(missed int) int {
	if missed <= 0 {
		return 5
	}
	switch {
	case missed <= 2:
		return 3
	case missed <= 4:
		return 2
	default:
		return 1
	}
}
