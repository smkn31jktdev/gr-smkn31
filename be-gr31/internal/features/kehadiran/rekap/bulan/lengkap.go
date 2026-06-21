package bulan

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"be-gr31/internal/kalender"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/util"
)

// GetRekapKelasLengkap
func (s *Service) GetRekapKelasLengkap(ctx context.Context, bulan, kelas string) (*rekapmodel.RekapKelasLengkap, error) {
	if _, err := time.Parse("2006-01", bulan); err != nil {
		return nil, fmt.Errorf("format bulan harus YYYY-MM")
	}

	hariEfektif := kalender.HariEfektif(bulan)

	cacheKey := "cache:rekap:lengkap:" + bulan + ":" + kelas
	var cached rekapmodel.RekapKelasLengkap
	if s.getCachedJSON(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	// Ambil seluruh data dari koleksi rekap_absensi untuk bulan tersebut.
	filter := rekapmodel.RekapFilter{
		Kelas:      kelas,
		BulanTahun: bulan,
	}
	fetcher := util.PagedFetcher[rekapmodel.RekapBulanan](func(ctx context.Context, size int, state string) ([]rekapmodel.RekapBulanan, string, error) {
		return s.store.ListPaged(ctx, filter, size, state)
	})
	rekapRecords, err := util.FetchAll(ctx, fetcher, util.MaxPageSize)
	if err != nil {
		return nil, err
	}

	// Ambil roster siswa agar siswa tanpa data tetap muncul
	roster := s.FetchRoster(ctx, kelas)

	// Counter per siswa, diinisialisasi dari roster
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

	for _, rec := range rekapRecords {
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

		c.hadir += rec.TotalHadir
		c.izin += rec.TotalIzin
		c.sakit += rec.TotalSakit
		c.alpa += rec.TotalTidakHadir
		c.magang += rec.TotalMagang
	}

	// Agregasi per kelas
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

		// Item per siswa: isi hari yang belum tercatat sebagai alpa.
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

	// 6. Susun ringkasan per kelas terurut + grand total.
	classNames := make([]string, 0, len(classes))
	for k := range classes {
		classNames = append(classNames, k)
	}
	sort.Strings(classNames)

	summaryByClass := make([]rekapmodel.RekapKelasSummary, 0, len(classNames))
	var grand rekapmodel.RekapRange
	grand.BulanTahun = bulan
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

	// Urutkan daftar siswa: kelas lalu nama.
	sort.Slice(students, func(i, j int) bool {
		if students[i].Kelas != students[j].Kelas {
			return students[i].Kelas < students[j].Kelas
		}
		return students[i].NamaSiswa < students[j].NamaSiswa
	})

	result := &rekapmodel.RekapKelasLengkap{
		BulanTahun:       bulan,
		HariEfektif:      hariEfektif,
		SummaryByClass:   summaryByClass,
		SummaryByStudent: students,
		SummaryRange:     grand,
	}
	s.setCachedJSON(ctx, cacheKey, result, 90*time.Second)
	return result, nil
}

// GetPersentaseKehadiranKelas
func (s *Service) GetPersentaseKehadiranKelas(ctx context.Context, bulan, kelas string) (*rekapmodel.RekapPersentaseKelas, error) {
	full, err := s.GetRekapKelasLengkap(ctx, bulan, kelas)
	if err != nil {
		return nil, err
	}
	return &rekapmodel.RekapPersentaseKelas{
		BulanTahun:     full.BulanTahun,
		HariEfektif:    full.HariEfektif,
		SummaryByClass: full.SummaryByClass,
		SummaryRange:   full.SummaryRange,
	}, nil
}

// classifyStatus
func classifyStatus(status, alasan string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "hadir":
		return "hadir"
	case "magang":
		return "magang"
	case "izin":
		return "izin"
	case "sakit":
		return "sakit"
	case "tidak_hadir":
		a := strings.ToLower(alasan)
		switch {
		case strings.Contains(a, "magang"):
			return "magang"
		case strings.Contains(a, "izin"):
			return "izin"
		case strings.Contains(a, "sakit"):
			return "sakit"
		default:
			return "alpa"
		}
	default:
		return "alpa"
	}
}

func normalizeKelas(kelas string) string {
	k := strings.TrimSpace(kelas)
	if k == "" {
		return "Tanpa Kelas"
	}
	return k
}

func clampPersen(v float64) float64 {
	if v > 100 {
		return 100
	}
	if v < 0 {
		return 0
	}
	return float64(int64(v*100+0.5)) / 100
}
