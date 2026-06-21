package bulan

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/features/g7/fetch"
	"be-gr31/internal/features/g7/content/summary"
	"be-gr31/internal/util"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("data G7 tidak ditemukan")
	ErrRekapFinal    = errors.New("rekap sudah final dan tidak bisa diubah")
	ErrFinalizeRule  = errors.New("rekap hanya bisa difinalisasi jika status reviewed dan terisi 2 penilai")
	ErrInvalidStatus = errors.New("status tidak valid (gunakan draft|reviewed|final)")
	ErrSiswaNotFound = errors.New("siswa tidak ditemukan")
)

// Service handles G7 monthly recap operations
type Service struct {
	rekapStore     *astra.G7RekapStore
	fetchService   *fetch.Service
	summaryService *summary.Service
	rosterLister   fetch.StudentRosterLister
}

// NewService creates a new Service instance
func NewService(
	rekapStore *astra.G7RekapStore,
	fetchService *fetch.Service,
	summaryService *summary.Service,
	rosterLister fetch.StudentRosterLister,
) *Service {
	return &Service{
		rekapStore:     rekapStore,
		fetchService:   fetchService,
		summaryService: summaryService,
		rosterLister:   rosterLister,
	}
}

// UpsertRekap saves or updates monthly G7 score assessments
func (s *Service) UpsertRekap(ctx context.Context, req g7model.G7RekapUpsertRequest) (*g7model.G7Rekap, error) {
	if _, err := time.Parse("2006-01", req.BulanTahun); err != nil {
		return nil, fmt.Errorf("request tidak valid: bulanTahun harus YYYY-MM")
	}
	if !validRekapStatus(req.Status) {
		return nil, ErrInvalidStatus
	}

	siswa, err := s.rosterLister.FindStudentByNISN(ctx, req.NISN)
	if err != nil {
		return nil, err
	}
	if siswa == nil {
		return nil, ErrSiswaNotFound
	}

	if err := g7model.ValidateSkor(req.Skor); err != nil {
		return nil, err
	}

	existing, err := s.rekapStore.FindByNISNBulan(ctx, req.NISN, req.BulanTahun)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Status == g7model.StatusFinal {
		return nil, ErrRekapFinal
	}

	now := common.FlexTime(time.Now().Format(time.RFC3339))
	var data g7model.G7Rekap
	if existing != nil {
		data = *existing
	} else {
		data = g7model.G7Rekap{ID: uuid.New().String(), CreatedAt: now}
	}

	data.NISN = req.NISN
	data.NamaSiswa = siswa.Nama
	data.Kelas = siswa.Kelas
	data.BulanTahun = req.BulanTahun
	data.Skor = req.Skor

	if req.WaliKelas != "" {
		data.WaliKelas = req.WaliKelas
	}
	if req.GuruBK != "" {
		data.GuruBK = req.GuruBK
	}

	perolehan, maks, akhir, predikat := g7model.HitungNilaiAkhir(req.Skor, req.BulanTahun)
	data.NilaiMaks = maks
	data.NilaiPerolehan = perolehan
	data.NilaiAkhir = akhir
	data.Predikat = predikat
	data.HariTercatat = s.fetchService.CountHariTercatat(ctx, req.NISN, req.BulanTahun)

	target := req.Status
	if target == "" {
		if existing != nil {
			target = existing.Status
		} else {
			target = g7model.StatusDraft
		}
	}
	if target == "" {
		target = g7model.StatusDraft
	}
	if target == g7model.StatusFinal {
		curStatus := ""
		if existing != nil {
			curStatus = existing.Status
		}
		if curStatus != g7model.StatusReviewed || data.JumlahAssessor() < 2 {
			return nil, ErrFinalizeRule
		}
		data.TanggalFinal = util.NowJakarta().Format("2006-01-02")
	}
	data.Status = target
	data.UpdatedAt = now

	if err := s.rekapStore.Upsert(ctx, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// DeleteRekap deletes a monthly G7 recap by ID
func (s *Service) DeleteRekap(ctx context.Context, id string) error {
	data, err := s.rekapStore.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if data == nil {
		return ErrNotFound
	}
	if data.Status == g7model.StatusFinal {
		return ErrRekapFinal
	}
	return s.rekapStore.Delete(ctx, id)
}

// RekapKelasLengkap returns G7 score recaps for an entire class roster joined with summary statistics
func (s *Service) RekapKelasLengkap(ctx context.Context, bulan, kelas string, isGuruWali bool, asuhan []authmodel.Siswa, adminID, adminRole string) (*g7model.G7RekapKelasLengkap, error) {
	if _, err := time.Parse("2006-01", bulan); err != nil {
		return nil, fmt.Errorf("request tidak valid: bulanTahun harus YYYY-MM")
	}

	var roster []authmodel.Siswa
	var allowedNISN map[string]bool
	if isGuruWali {
		roster = asuhan
		allowedNISN = make(map[string]bool)
		for _, st := range roster {
			if st.NIS != "" {
				allowedNISN[st.NIS] = true
			}
		}
	} else {
		roster = s.fetchService.FetchRoster(ctx, kelas)
	}

	rekaps, err := s.fetchService.FetchAllRekap(ctx, g7model.G7RekapFilter{BulanTahun: bulan, Kelas: kelas})
	if err != nil {
		return nil, err
	}

	if isGuruWali {
		filteredRekaps := make([]g7model.G7Rekap, 0)
		for _, r := range rekaps {
			if allowedNISN[r.NISN] {
				filteredRekaps = append(filteredRekaps, r)
			}
		}
		rekaps = filteredRekaps
	}

	rekapByNISN := make(map[string]g7model.G7Rekap, len(rekaps))
	for _, r := range rekaps {
		rekapByNISN[r.NISN] = r
	}

	items := make([]g7model.G7RekapSiswaItem, 0)
	seen := map[string]struct{}{}

	addItem := func(nisn, nama, kls string, r *g7model.G7Rekap) {
		if nisn == "" {
			return
		}
		if _, dup := seen[nisn]; dup {
			return
		}
		seen[nisn] = struct{}{}

		item := g7model.G7RekapSiswaItem{
			NISN:      nisn,
			NamaSiswa: nama,
			Kelas:     kls,
			NilaiMaks: g7model.NilaiMaks,
		}
		if r != nil {
			item.NilaiPerolehan = r.NilaiPerolehan
			item.NilaiMaks = r.NilaiMaks
			item.NilaiAkhir = r.NilaiAkhir
			item.Predikat = r.Predikat
			item.Status = r.Status
			item.HariTercatat = r.HariTercatat
			item.SudahDinilai = true
		} else {
			item.Status = "belum_dinilai"
			item.Predikat = "-"
		}
		items = append(items, item)
	}

	for _, st := range roster {
		var rp *g7model.G7Rekap
		if r, ok := rekapByNISN[st.NIS]; ok {
			rp = &r
		}
		addItem(st.NIS, st.Nama, st.Kelas, rp)
	}
	for i := range rekaps {
		r := rekaps[i]
		addItem(r.NISN, r.NamaSiswa, r.Kelas, &r)
	}

	stat, err := s.summaryService.Statistik(ctx, bulan, kelas, isGuruWali, asuhan, adminID, adminRole)
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].Kelas != items[j].Kelas {
			return items[i].Kelas < items[j].Kelas
		}
		return items[i].NamaSiswa < items[j].NamaSiswa
	})

	sudah := len(rekaps)
	total := len(items)
	if stat != nil && stat.TotalSiswa > total {
		total = stat.TotalSiswa
	}
	belum := total - sudah
	if belum < 0 {
		belum = 0
	}

	return &g7model.G7RekapKelasLengkap{
		Kelas:        kelas,
		BulanTahun:   bulan,
		TotalSiswa:   total,
		SudahDinilai: sudah,
		BelumDinilai: belum,
		Statistik:    stat,
		Siswa:        items,
	}, nil
}

// UpdateRekapHariTercatat updates the journal logs counter in the G7 monthly rekap
func (s *Service) UpdateRekapHariTercatat(ctx context.Context, nisn, bulanTahun string) error {
	count := s.fetchService.CountHariTercatat(ctx, nisn, bulanTahun)

	existing, err := s.rekapStore.FindByNISNBulan(ctx, nisn, bulanTahun)
	if err != nil {
		return err
	}

	now := common.FlexTime(time.Now().Format(time.RFC3339))
	if existing != nil {
		if existing.Status == g7model.StatusFinal {
			return nil
		}
		existing.HariTercatat = count
		existing.UpdatedAt = now
		return s.rekapStore.Upsert(ctx, existing)
	}

	siswa, err := s.rosterLister.FindStudentByNISN(ctx, nisn)
	if err != nil {
		return err
	}
	if siswa == nil {
		return ErrSiswaNotFound
	}

	newRekap := &g7model.G7Rekap{
		ID:           uuid.New().String(),
		NISN:         nisn,
		NamaSiswa:    siswa.Nama,
		Kelas:        siswa.Kelas,
		BulanTahun:   bulanTahun,
		HariTercatat: count,
		Status:       g7model.StatusDraft,
		NilaiMaks:    g7model.NilaiMaks,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.rekapStore.Upsert(ctx, newRekap)
}

func validRekapStatus(s string) bool {
	switch s {
	case "", g7model.StatusDraft, g7model.StatusReviewed, g7model.StatusFinal:
		return true
	default:
		return false
	}
}

func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}
