package bulan

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"sort"
	"strings"
	"time"

	"be-gr31/internal/kalender"
	authmodel "be-gr31/internal/model/auth"
	kehadiranmodel "be-gr31/internal/model/kehadiran"
	rekapmodel "be-gr31/internal/model/rekap"
	"be-gr31/internal/util"
)

// GetRekapHarian
func (s *Service) GetRekapHarian(ctx context.Context, tanggal, kelas string) (*rekapmodel.RekapHarian, error) {
	if s.kehadiranSvc == nil {
		return &rekapmodel.RekapHarian{Tanggal: tanggal}, nil
	}

	totalSiswa := 0
	if s.supabaseClient != nil && s.supabaseClient.DB != nil {
		var n int
		var query string
		var args []interface{}
		if kelas != "" {
			query = `SELECT COUNT(*) FROM akun_siswa WHERE kelas = $1`
			args = append(args, kelas)
		} else {
			query = `SELECT COUNT(*) FROM akun_siswa`
		}
		err := s.supabaseClient.DB.QueryRowContext(ctx, query, args...).Scan(&n)
		if err != nil {
			log.Printf("WARNING: gagal menghitung total siswa dari Supabase: %v", err)
		} else {
			totalSiswa = n
		}
	} else if s.studentStore != nil {
		n, err := s.studentStore.Count(ctx, kelas)
		if err != nil {
			log.Printf("WARNING: gagal menghitung total siswa: %v", err)
		} else {
			totalSiswa = n
		}
	}

	filter := kehadiranmodel.KehadiranFilter{
		Tanggal: tanggal,
		Kelas:   kelas,
	}
	fetcher := util.PagedFetcher[kehadiranmodel.Kehadiran](func(ctx context.Context, size int, state string) ([]kehadiranmodel.Kehadiran, string, error) {
		return s.kehadiranSvc.ListPaged(ctx, filter, size, state)
	})
	items, err := util.FetchAll(ctx, fetcher, util.MaxPageSize)
	if err != nil {
		return nil, err
	}

	result := &rekapmodel.RekapHarian{
		Tanggal:    tanggal,
		TotalSiswa: totalSiswa,
	}
	for i := range items {
		switch items[i].Status {
		case "hadir":
			result.TotalHadir++
		case "izin":
			result.TotalIzin++
		case "sakit":
			result.TotalSakit++
		case "magang":
			result.TotalMagang++
		default:
			result.TotalAlpa++
		}
	}

	if result.TotalSiswa == 0 {
		result.TotalSiswa = len(items)
	}

	return result, nil
}

// GetRingkasanKelas
func (s *Service) GetRingkasanKelas(ctx context.Context, kelas, bulanTahun string, page, limit int) ([]rekapmodel.RingkasanSiswa, bool, int, error) {
	filter := rekapmodel.RekapFilter{
		Kelas:      kelas,
		BulanTahun: bulanTahun,
		Page:       page,
		Limit:      limit,
	}

	rekaps, hasMore, total, err := s.ListRekap(ctx, filter)
	if err != nil {
		return nil, false, 0, err
	}

	result := make([]rekapmodel.RingkasanSiswa, 0, len(rekaps))
	for _, r := range rekaps {
		result = append(result, rekapmodel.RingkasanSiswa{
			NIS:             r.NIS,
			NamaSiswa:       r.NamaSiswa,
			Kelas:           r.Kelas,
			BulanTahun:      r.BulanTahun,
			TotalHadir:      r.TotalHadir,
			TotalIzin:       r.TotalIzin,
			TotalSakit:      r.TotalSakit,
			TotalTidakHadir: r.TotalTidakHadir,
			TotalMagang:     r.TotalMagang,
			PersentaseHadir: r.PersentaseHadir,
		})
	}
	return result, hasMore, total, nil
}

// GetKehadiranBulananSiswa
func (s *Service) GetKehadiranBulananSiswa(ctx context.Context, nis, bulanTahun string) (*rekapmodel.KehadiranBulananSiswa, error) {
	filter := kehadiranmodel.KehadiranFilter{
		NIS:       nis,
		BulanDari: bulanTahun,
		BulanKe:   bulanTahun,
	}
	semuaKehadiran, _, err := s.kehadiranSvc.ListPaged(ctx, filter, 500, "")
	if err != nil {
		return nil, err
	}

	items := make([]rekapmodel.KehadiranHariItem, 0, len(semuaKehadiran))
	for _, k := range semuaKehadiran {
		items = append(items, rekapmodel.KehadiranHariItem{
			Tanggal:    k.Tanggal,
			Status:     k.Status,
			WaktuAbsen: k.WaktuAbsen,
			Alasan:     k.Alasan,
			FotoIzin:   k.FotoIzin,
		})
	}

	summary := rekapmodel.KehadiranBulanSummary{}
	rekap, err := s.store.FindByKey(ctx, nis+"_"+bulanTahun)
	if err == nil && rekap != nil {
		summary.TotalHadir = rekap.TotalHadir
		summary.TotalIzin = rekap.TotalIzin
		summary.TotalSakit = rekap.TotalSakit
		summary.TotalAlpa = rekap.TotalTidakHadir
		summary.TotalMagang = rekap.TotalMagang
		summary.TotalHariEfektif = rekap.TotalHariEfektif
		summary.PersentaseHadir = rekap.PersentaseHadir
	} else {
		hariEfektif := kalender.HariEfektif(bulanTahun)
		summary.TotalHariEfektif = hariEfektif
		for _, k := range semuaKehadiran {
			switch k.Status {
			case "hadir":
				summary.TotalHadir++
			case "izin":
				summary.TotalIzin++
			case "sakit":
				summary.TotalSakit++
			case "magang":
				summary.TotalMagang++
			default:
				summary.TotalAlpa++
			}
		}

		counted := summary.TotalHadir + summary.TotalIzin + summary.TotalSakit + summary.TotalMagang + summary.TotalAlpa
		if missing := hariEfektif - counted; missing > 0 {
			summary.TotalAlpa += missing
		}

		if summary.TotalHariEfektif > 0 {
			persen := float64(summary.TotalHadir+summary.TotalMagang) / float64(summary.TotalHariEfektif) * 100
			if persen > 100 {
				persen = 100
			} else if persen < 0 {
				persen = 0
			}
			summary.PersentaseHadir = math.Round(persen*10) / 10
		}
	}

	return &rekapmodel.KehadiranBulananSiswa{
		BulanTahun: bulanTahun,
		Kehadiran:  items,
		Summary:    summary,
	}, nil
}

// GetRingkasanSiswa mengambil rekap per-bulan satu siswa dalam rentang dari–sampai
func (s *Service) GetRingkasanSiswa(ctx context.Context, nis, dari, sampai string) ([]rekapmodel.RekapBulanan, error) {
	bulanList := BulanDalamRentang(dari, sampai)

	result := make([]rekapmodel.RekapBulanan, 0, len(bulanList))
	for _, bulan := range bulanList {
		r, err := s.store.FindByKey(ctx, nis+"_"+bulan)
		if err != nil || r == nil {
			continue
		}
		result = append(result, *r)
	}
	return result, nil
}

// FetchRoster mengambil seluruh siswa pada kelas (auto-paging)
func (s *Service) FetchRoster(ctx context.Context, kelas string) []authmodel.Siswa {
	if s.supabaseClient != nil && s.supabaseClient.DB != nil {
		var query string
		var args []interface{}
		if kelas != "" {
			query = `SELECT nis, nama, kelas FROM akun_siswa WHERE kelas = $1 ORDER BY nama ASC`
			args = append(args, kelas)
		} else {
			query = `SELECT nis, nama, kelas FROM akun_siswa ORDER BY nama ASC`
		}
		rows, err := s.supabaseClient.DB.QueryContext(ctx, query, args...)
		if err != nil {
			log.Printf("WARNING: gagal memuat roster kelas %q dari Supabase: %v", kelas, err)
			return nil
		}
		defer rows.Close()

		var students []authmodel.Siswa
		for rows.Next() {
			var st authmodel.Siswa
			if err := rows.Scan(&st.NIS, &st.Nama, &st.Kelas); err != nil {
				log.Printf("WARNING: gagal scan row student: %v", err)
				continue
			}
			students = append(students, st)
		}
		return students
	}

	if s.studentStore == nil {
		return nil
	}
	fetcher := util.PagedFetcher[authmodel.Siswa](func(ctx context.Context, size int, state string) ([]authmodel.Siswa, string, error) {
		return s.studentStore.ListPaged(ctx, kelas, size, state)
	})
	all, err := util.FetchAll(ctx, fetcher, util.MaxPageSize)
	if err != nil {
		log.Printf("WARNING: gagal memuat roster kelas %q: %v", kelas, err)
		return nil
	}
	return all
}

// GetKelas extracts all unique kelas from the students database
func (s *Service) GetKelas(ctx context.Context) ([]string, error) {
	cacheKey := "cache:kelas_list"
	if cached, err := s.rdb.Get(ctx, cacheKey).Bytes(); err == nil {
		var cachedResponse []string
		if err := json.Unmarshal(cached, &cachedResponse); err == nil {
			return cachedResponse, nil
		}
	}

	students := s.FetchRoster(ctx, "")

	kelasMap := make(map[string]bool)
	for _, student := range students {
		raw := strings.TrimSpace(student.Kelas)
		if raw != "" {
			kelasMap[raw] = true
		}
	}

	kelas := make([]string, 0, len(kelasMap))
	for k := range kelasMap {
		kelas = append(kelas, k)
	}
	sort.Strings(kelas)

	if data, err := json.Marshal(kelas); err == nil {
		if setErr := s.rdb.Set(ctx, cacheKey, data, 10*time.Minute).Err(); setErr != nil {
			log.Printf("WARNING: gagal cache kelas_list: %v", setErr)
		}
	}

	return kelas, nil
}

// GetKelasJurusan
func (s *Service) GetKelasJurusan(ctx context.Context) (*rekapmodel.KelasJurusanResponse, error) {
	kelas, err := s.GetKelas(ctx)
	if err != nil {
		return nil, err
	}
	return &rekapmodel.KelasJurusanResponse{
		KelasLengkap: kelas,
	}, nil
}

// GetKehadiranMingguanSiswa mengambil data kehadiran mingguan lengkap (Senin-Jumat) untuk satu siswa.
func (s *Service) GetKehadiranMingguanSiswa(ctx context.Context, nis, senin string) (*rekapmodel.KehadiranBulananSiswa, error) {
	mon := parseMondayOrCurrent(senin)
	fri := mon.AddDate(0, 0, 4)

	seninStr := mon.Format("2006-01-02")
	jumatStr := fri.Format("2006-01-02")

	filter := kehadiranmodel.KehadiranFilter{
		NIS:       nis,
		BulanDari: seninStr,
		BulanKe:   jumatStr,
	}
	semuaKehadiran, _, err := s.kehadiranSvc.ListPaged(ctx, filter, 50, "")
	if err != nil {
		return nil, err
	}

	// Buat map tanggal untuk pencarian cepat
	hadirMap := make(map[string]kehadiranmodel.Kehadiran)
	for _, k := range semuaKehadiran {
		hadirMap[k.Tanggal] = k
	}

	items := make([]rekapmodel.KehadiranHariItem, 0, 5)
	summary := rekapmodel.KehadiranBulanSummary{
		TotalHariEfektif: 5,
	}

	// Loop Senin s.d. Jumat untuk memastikan semua hari terisi
	for d := mon; !d.After(fri); d = d.AddDate(0, 0, 1) {
		tglStr := d.Format("2006-01-02")
		if k, found := hadirMap[tglStr]; found {
			items = append(items, rekapmodel.KehadiranHariItem{
				Tanggal:    k.Tanggal,
				Status:     k.Status,
				WaktuAbsen: k.WaktuAbsen,
				Alasan:     k.Alasan,
				FotoIzin:   k.FotoIzin,
			})
			switch classifyStatus(k.Status, k.Alasan) {
			case "hadir":
				summary.TotalHadir++
			case "izin":
				summary.TotalIzin++
			case "sakit":
				summary.TotalSakit++
			case "magang":
				summary.TotalMagang++
			default:
				summary.TotalAlpa++
			}
		} else {
			// Jika belum ada absen di hari sekolah, dianggap Alpa
			items = append(items, rekapmodel.KehadiranHariItem{
				Tanggal:    tglStr,
				Status:     "tidak_hadir",
				WaktuAbsen: "",
			})
			summary.TotalAlpa++
		}
	}

	if summary.TotalHariEfektif > 0 {
		persen := float64(summary.TotalHadir+summary.TotalMagang) / float64(summary.TotalHariEfektif) * 100
		summary.PersentaseHadir = math.Round(persen*10) / 10
	}

	// Sort items by Tanggal asc
	sort.Slice(items, func(i, j int) bool {
		return items[i].Tanggal < items[j].Tanggal
	})

	return &rekapmodel.KehadiranBulananSiswa{
		BulanTahun: seninStr + "_" + jumatStr,
		Kehadiran:  items,
		Summary:    summary,
	}, nil
}

