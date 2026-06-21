package bulan

import (
	"context"
	"sort"
	"strings"
	"time"

	authmodel "be-gr31/internal/model/auth"
	kehadiranmodel "be-gr31/internal/model/kehadiran"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/util"
)

// GetRekapMingguanKelas
func (s *Service) GetRekapMingguanKelas(ctx context.Context, kelas, senin string) (*rekapmodel.RekapMingguanKelas, error) {
	mon := parseMondayOrCurrent(senin)
	fri := mon.AddDate(0, 0, 4)

	seninStr := mon.Format("2006-01-02")
	jumatStr := fri.Format("2006-01-02")

	cacheKey := "cache:rekap:mingguan:" + seninStr + ":" + kelas
	var cached rekapmodel.RekapMingguanKelas
	if s.getCachedJSON(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	// Hari efektif minggu ini = hari kerja Senin–Jumat (tanpa mengecualikan libur).
	hariEfektif := 0
	for d := mon; !d.After(fri); d = d.AddDate(0, 0, 1) {
		if isWeekday(d) {
			hariEfektif++
		}
	}

	records := s.fetchKehadiranRentang(ctx, kelas, seninStr, jumatStr)
	roster := s.FetchRoster(ctx, kelas)

	summaryByClass, students, grand := aggregateFromKehadiran(records, roster, hariEfektif)

	result := &rekapmodel.RekapMingguanKelas{
		Senin:            seninStr,
		Jumat:            jumatStr,
		HariEfektif:      hariEfektif,
		SummaryByClass:   summaryByClass,
		SummaryByStudent: students,
		SummaryRange:     grand,
	}
	s.setCachedJSON(ctx, cacheKey, result, 90*time.Second)
	return result, nil
}

// parseMondayOrCurrent mengembalikan tanggal Senin. Jika input valid, ia
// dinormalkan ke Senin pada minggu yang sama; jika kosong/invalid, dipakai Senin
// minggu berjalan.
func parseMondayOrCurrent(senin string) time.Time {
	loc := util.JakartaLoc()
	var base time.Time
	if t, err := time.ParseInLocation("2006-01-02", strings.TrimSpace(senin), loc); err == nil {
		base = t
	} else {
		base = util.NowJakarta()
	}
	// Normalkan ke Senin minggu tersebut.
	wd := int(base.Weekday()) // Minggu=0..Sabtu=6
	if wd == 0 {
		wd = 7 // perlakukan Minggu sebagai akhir minggu
	}
	monday := base.AddDate(0, 0, -(wd - 1))
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, loc)
}

// isWeekday mengembalikan true jika hari adalah Senin–Jumat.
func isWeekday(t time.Time) bool {
	wd := t.Weekday()
	return wd != time.Saturday && wd != time.Sunday
}

// isWeekdayStr sama dengan isWeekday namun menerima tanggal string YYYY-MM-DD.
func isWeekdayStr(tanggal string) bool {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(tanggal))
	if err != nil {
		return false
	}
	return isWeekday(t)
}

// fetchKehadiranRentang mengambil seluruh catatan kehadiran dalam rentang tanggal
// (inklusif) untuk kelas tertentu (semua kelas bila kosong), dengan auto-paging.
func (s *Service) fetchKehadiranRentang(ctx context.Context, kelas, dari, sampai string) []kehadiranmodel.Kehadiran {
	if s.kehadiranSvc == nil {
		return nil
	}
	filter := kehadiranmodel.KehadiranFilter{
		Kelas:         kelas,
		TanggalDari:   dari,
		TanggalSampai: sampai,
	}
	fetcher := util.PagedFetcher[kehadiranmodel.Kehadiran](func(ctx context.Context, size int, state string) ([]kehadiranmodel.Kehadiran, string, error) {
		return s.kehadiranSvc.ListPaged(ctx, filter, size, state)
	})
	all, err := util.FetchAll(ctx, fetcher, util.MaxPageSize)
	if err != nil {
		return nil
	}
	return all
}

// aggregateFromKehadiran mengagregasi catatan kehadiran harian menjadi ringkasan
// per-kelas, daftar per-siswa, dan grand-total. Dipakai oleh rekap mingguan.
//
// Aturan:
//   - Deduplikasi per (nis, tanggal); abaikan akhir pekan (Sabtu/Minggu).
//     Hari libur kalender TIDAK dikecualikan untuk rekap mingguan.
//   - Ringkasan kelas memakai jumlah TERCATAT (tanpa mengisi alpa).
//   - Daftar per-siswa mengisi hari yang belum tercatat sebagai alpa.
//   - %Kehadiran = (hadir + magang) / (totalSiswa × hariEfektif) × 100.
func aggregateFromKehadiran(records []kehadiranmodel.Kehadiran, roster []authmodel.Siswa, hariEfektif int) ([]rekapmodel.RekapKelasSummary, []rekapmodel.RekapSiswaItem, rekapmodel.RekapRange) {
	uniqueByStudentDate := make(map[string]kehadiranmodel.Kehadiran, len(records))
	for _, rec := range records {
		nis := strings.TrimSpace(rec.NIS)
		tanggal := strings.TrimSpace(rec.Tanggal)
		if nis == "" || tanggal == "" {
			continue
		}
		if !isWeekdayStr(tanggal) {
			continue
		}
		uniqueByStudentDate[nis+"_"+tanggal] = rec
	}

	type studentCounter struct {
		nama    string
		kelas   string
		hadir   int
		izin    int
		sakit   int
		alpa    int
		magang  int
		adaData bool
	}
	counters := map[string]*studentCounter{}
	for _, st := range roster {
		nis := strings.TrimSpace(st.NIS)
		if nis == "" {
			continue
		}
		counters[nis] = &studentCounter{
			nama:  strings.TrimSpace(st.Nama),
			kelas: normalizeKelas(st.Kelas),
		}
	}

	for _, rec := range uniqueByStudentDate {
		nis := strings.TrimSpace(rec.NIS)
		c, ok := counters[nis]
		if !ok {
			c = &studentCounter{
				nama:  strings.TrimSpace(rec.NamaSiswa),
				kelas: normalizeKelas(rec.Kelas),
			}
			counters[nis] = c
		}
		if c.kelas == "Tanpa Kelas" {
			if k := normalizeKelas(rec.Kelas); k != "Tanpa Kelas" {
				c.kelas = k
			}
		}
		if c.nama == "" {
			c.nama = strings.TrimSpace(rec.NamaSiswa)
		}
		c.adaData = true

		switch classifyStatus(rec.Status, rec.Alasan) {
		case "hadir":
			c.hadir++
		case "magang":
			c.magang++
		case "izin":
			c.izin++
		case "sakit":
			c.sakit++
		default:
			c.alpa++
		}
	}

	type classAgg struct {
		totalSiswa map[string]struct{}
		hadirSiswa map[string]struct{}
		hadir      int
		izin       int
		sakit      int
		alpa       int
		magang     int
	}
	classes := map[string]*classAgg{}
	ensureClass := func(k string) *classAgg {
		ca, ok := classes[k]
		if !ok {
			ca = &classAgg{
				totalSiswa: map[string]struct{}{},
				hadirSiswa: map[string]struct{}{},
			}
			classes[k] = ca
		}
		return ca
	}

	students := make([]rekapmodel.RekapSiswaItem, 0, len(counters))
	nisList := make([]string, 0, len(counters))
	for nis := range counters {
		nisList = append(nisList, nis)
	}
	sort.Strings(nisList)

	for _, nis := range nisList {
		c := counters[nis]
		kls := c.kelas
		if kls == "" {
			kls = "Tanpa Kelas"
		}

		ca := ensureClass(kls)
		ca.totalSiswa[nis] = struct{}{}
		ca.hadir += c.hadir
		ca.izin += c.izin
		ca.sakit += c.sakit
		ca.alpa += c.alpa
		ca.magang += c.magang
		if c.hadir > 0 || c.magang > 0 {
			ca.hadirSiswa[nis] = struct{}{}
		}

		alpaDisplay := c.alpa
		counted := c.hadir + c.izin + c.sakit + c.alpa + c.magang
		if missing := hariEfektif - counted; missing > 0 {
			alpaDisplay += missing
		}

		efektifHadir := c.hadir + c.magang
		tingkat := 0.0
		if hariEfektif > 0 {
			tingkat = clampPersen(float64(efektifHadir) / float64(hariEfektif) * 100)
		}

		students = append(students, rekapmodel.RekapSiswaItem{
			NIS:              nis,
			NamaSiswa:        c.nama,
			Kelas:            kls,
			TotalHadir:       c.hadir,
			TotalIzin:        c.izin,
			TotalSakit:       c.sakit,
			TotalAlpa:        alpaDisplay,
			TotalMagang:      c.magang,
			HariEfektif:      hariEfektif,
			TingkatKehadiran: tingkat,
			AdaData:          c.adaData,
		})
	}

	classNames := make([]string, 0, len(classes))
	for k := range classes {
		classNames = append(classNames, k)
	}
	sort.Strings(classNames)

	summaryByClass := make([]rekapmodel.RekapKelasSummary, 0, len(classNames))
	var grand rekapmodel.RekapRange
	grand.HariEfektif = hariEfektif
	uniqueHadirAll := map[string]struct{}{}
	totalSiswaAll := map[string]struct{}{}

	for _, k := range classNames {
		ca := classes[k]
		totalSiswa := len(ca.totalSiswa)
		siswaUnikHadir := len(ca.hadirSiswa)

		tingkat := 0.0
		expected := totalSiswa * hariEfektif
		if expected > 0 {
			tingkat = clampPersen(float64(ca.hadir+ca.magang) / float64(expected) * 100)
		}

		summaryByClass = append(summaryByClass, rekapmodel.RekapKelasSummary{
			Kelas:            k,
			TotalSiswa:       totalSiswa,
			HariEfektif:      hariEfektif,
			TotalHadir:       ca.hadir,
			TotalIzin:        ca.izin,
			TotalSakit:       ca.sakit,
			TotalAlpa:        ca.alpa,
			TotalMagang:      ca.magang,
			SiswaUnikHadir:   siswaUnikHadir,
			TingkatKehadiran: tingkat,
		})

		grand.TotalHadir += ca.hadir
		grand.TotalIzin += ca.izin
		grand.TotalSakit += ca.sakit
		grand.TotalAlpa += ca.alpa
		grand.TotalMagang += ca.magang
		for nis := range ca.totalSiswa {
			totalSiswaAll[nis] = struct{}{}
		}
		for nis := range ca.hadirSiswa {
			uniqueHadirAll[nis] = struct{}{}
		}
	}

	grand.TotalSiswa = len(totalSiswaAll)
	grand.SiswaUnikHadir = len(uniqueHadirAll)
	if expectedAll := grand.TotalSiswa * hariEfektif; expectedAll > 0 {
		grand.TingkatKehadiran = clampPersen(float64(grand.TotalHadir+grand.TotalMagang) / float64(expectedAll) * 100)
	}

	sort.Slice(students, func(i, j int) bool {
		if students[i].Kelas != students[j].Kelas {
			return students[i].Kelas < students[j].Kelas
		}
		return students[i].NamaSiswa < students[j].NamaSiswa
	})

	return summaryByClass, students, grand
}
