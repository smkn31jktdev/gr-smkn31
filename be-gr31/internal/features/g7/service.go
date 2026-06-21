package g7

import (
	"context"
	"errors"
	"fmt"
	"strings"

	authmodel "be-gr31/internal/model/auth"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/features/g7/fetch"
	"be-gr31/internal/features/g7/content/jurnal"
	"be-gr31/internal/features/g7/content/rekap"
	"be-gr31/internal/features/g7/content/suggest"
	"be-gr31/internal/features/g7/content/summary"
	"be-gr31/internal/features/g7/content/wali"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound      = rekap.ErrNotFound
	ErrRekapFinal    = rekap.ErrRekapFinal
	ErrFinalizeRule  = rekap.ErrFinalizeRule
	ErrInvalidStatus = rekap.ErrInvalidStatus
	ErrSiswaNotFound = rekap.ErrSiswaNotFound
	ErrFutureDate    = errors.New("tidak bisa input untuk tanggal mendatang")
	ErrBadRequest    = g7model.ErrBadRequest
)

var ErrSkorRange = g7model.ErrSkorRange

type Service struct {
	repo           *Repo
	rdb            *redis.Client
	fetchService   *fetch.Service
	rekapService   *rekap.Service
	summaryService *summary.Service
	suggestService *suggest.Service
	waliService    *wali.Service
	jurnalService  *jurnal.Service
}

func NewService(repo *Repo, rdb *redis.Client) *Service {
	fetchSvc := fetch.NewService(repo.rekapStore, repo.store, repo)
	summarySvc := summary.NewService(fetchSvc)
	rekapSvc := rekap.NewService(repo.rekapStore, fetchSvc, summarySvc, repo)
	suggestSvc := suggest.NewService(fetchSvc)
	waliSvc := wali.NewService(repo, fetchSvc)
	jurnalSvc := jurnal.NewService(repo, repo, rdb, fetchSvc, rekapSvc, waliSvc)

	return &Service{
		repo:           repo,
		rdb:            rdb,
		fetchService:   fetchSvc,
		summaryService: summarySvc,
		rekapService:   rekapSvc,
		suggestService: suggestSvc,
		waliService:    waliSvc,
		jurnalService:  jurnalSvc,
	}
}

func (s *Service) FindG7ByID(ctx context.Context, id string) (*g7model.G7, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) FindRekapByID(ctx context.Context, id string) (*g7model.G7Rekap, error) {
	return s.repo.FindRekapByID(ctx, id)
}

func (s *Service) checkIsGuruWali(ctx context.Context, adminID, adminRole string) (bool, []authmodel.Siswa, error) {
	return s.waliService.CheckIsGuruWali(ctx, adminID, adminRole)
}

func (s *Service) VerifyAdminAccessToStudent(ctx context.Context, adminID, adminRole, studentNIS string) error {
	return s.waliService.VerifyAdminAccessToStudent(ctx, adminID, adminRole, studentNIS)
}

func (s *Service) fetchRoster(ctx context.Context, kelas string) []authmodel.Siswa {
	return s.fetchService.FetchRoster(ctx, kelas)
}

func (s *Service) UpsertRekap(ctx context.Context, req g7model.G7RekapUpsertRequest) (*g7model.G7Rekap, error) {
	return s.rekapService.UpsertRekap(ctx, req)
}

func (s *Service) GetRekap(ctx context.Context, nisn, bulan string) (*g7model.G7Rekap, error) {
	return s.fetchService.GetRekap(ctx, nisn, bulan)
}

func (s *Service) ListRekap(ctx context.Context, filter g7model.G7RekapFilter) ([]g7model.G7Rekap, bool, int, error) {
	isGuruWali, asuhan, err := s.checkIsGuruWali(ctx, filter.AdminID, filter.AdminRole)
	if err != nil {
		return nil, false, 0, err
	}
	return s.fetchService.ListRekap(ctx, filter, isGuruWali, asuhan)
}

func (s *Service) DeleteRekap(ctx context.Context, id string) error {
	return s.rekapService.DeleteRekap(ctx, id)
}

func (s *Service) Statistik(ctx context.Context, bulan, kelas, adminID, adminRole string) (*g7model.G7RekapStatistik, error) {
	isGuruWali, asuhan, err := s.checkIsGuruWali(ctx, adminID, adminRole)
	if err != nil {
		return nil, err
	}
	return s.summaryService.Statistik(ctx, bulan, kelas, isGuruWali, asuhan, adminID, adminRole)
}

func (s *Service) RekapKelasLengkap(ctx context.Context, bulan, kelas, adminID, adminRole string) (*g7model.G7RekapKelasLengkap, error) {
	isGuruWali, asuhan, err := s.checkIsGuruWali(ctx, adminID, adminRole)
	if err != nil {
		return nil, err
	}
	return s.rekapService.RekapKelasLengkap(ctx, bulan, kelas, isGuruWali, asuhan, adminID, adminRole)
}

func (s *Service) EvaluateJurnalBulanan(ctx context.Context, nisn, bulan string) (*EvalReport, error) {
	return s.summaryService.EvaluateJurnalBulanan(ctx, nisn, bulan)
}

func (s *Service) BuildLaporanPDF(ctx context.Context, nisn, bulan string) (*g7model.PDFLaporan, string, error) {
	return s.summaryService.BuildLaporanPDF(ctx, nisn, bulan)
}

func (s *Service) UpdateRekapHariTercatat(ctx context.Context, nisn, bulanTahun string) error {
	return s.rekapService.UpdateRekapHariTercatat(ctx, nisn, bulanTahun)
}

func (s *Service) GetRekapSemesterKelas(ctx context.Context, semester, kelas, adminID, adminRole string) ([]g7model.G7SemesterStudentItem, error) {
	isGuruWali, asuhan, err := s.checkIsGuruWali(ctx, adminID, adminRole)
	if err != nil {
		return nil, err
	}
	return s.rekapService.GetRekapSemesterKelas(ctx, semester, kelas, isGuruWali, asuhan)
}

func (s *Service) Suggest(ctx context.Context, nisn, bulan string) (*g7model.G7SuggestResponse, error) {
	resp, err := s.suggestService.Suggest(ctx, nisn, bulan)
	if err != nil {
		if strings.Contains(err.Error(), "bulan harus YYYY-MM") {
			return nil, fmt.Errorf("%w: %v", ErrBadRequest, err)
		}
		return nil, err
	}
	return resp, nil
}

func (s *Service) Upsert(ctx context.Context, nisn string, req g7model.G7UpsertRequest) (*g7model.G7, error) {
	return s.jurnalService.Upsert(ctx, nisn, req)
}

func (s *Service) DashboardSiswa(ctx context.Context, nisn string) (*g7model.G7DashboardSiswa, error) {
	return s.jurnalService.DashboardSiswa(ctx, nisn)
}

func (s *Service) GetByTanggal(ctx context.Context, nisn, tanggal string) (*g7model.G7, error) {
	return s.jurnalService.GetByTanggal(ctx, nisn, tanggal)
}

func (s *Service) List(ctx context.Context, filter g7model.G7Filter) ([]g7model.G7, bool, int, error) {
	return s.jurnalService.List(ctx, filter)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.jurnalService.Delete(ctx, id)
}

func (s *Service) Summary(ctx context.Context, bulanTahun, kelas, adminID, adminRole string) ([]g7model.G7Summary, error) {
	return s.jurnalService.Summary(ctx, bulanTahun, kelas, adminID, adminRole)
}
