package rekap

import (
	"context"

	"be-gr31/internal/features/g7/fetch"
	"be-gr31/internal/features/g7/content/rekap/bulan"
	"be-gr31/internal/features/g7/content/rekap/semester"
	"be-gr31/internal/features/g7/content/summary"
	authmodel "be-gr31/internal/model/auth"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/storage/astra"
)

// Re-export errors for package compatibility
var (
	ErrNotFound      = bulan.ErrNotFound
	ErrRekapFinal    = bulan.ErrRekapFinal
	ErrFinalizeRule  = bulan.ErrFinalizeRule
	ErrInvalidStatus = bulan.ErrInvalidStatus
	ErrSiswaNotFound = bulan.ErrSiswaNotFound
)

// Service coordinates
type Service struct {
	bulanSvc    *bulan.Service
	semesterSvc *semester.Service
}

// NewService
func NewService(
	rekapStore *astra.G7RekapStore,
	fetchService *fetch.Service,
	summaryService *summary.Service,
	rosterLister fetch.StudentRosterLister,
) *Service {
	return &Service{
		bulanSvc:    bulan.NewService(rekapStore, fetchService, summaryService, rosterLister),
		semesterSvc: semester.NewService(fetchService),
	}
}

// UpsertRekap delegates G7 score updates to the monthly rekap service
func (s *Service) UpsertRekap(ctx context.Context, req g7model.G7RekapUpsertRequest) (*g7model.G7Rekap, error) {
	return s.bulanSvc.UpsertRekap(ctx, req)
}

// DeleteRekap delegates G7 rekap deletion to the monthly rekap service
func (s *Service) DeleteRekap(ctx context.Context, id string) error {
	return s.bulanSvc.DeleteRekap(ctx, id)
}

// RekapKelasLengkap delegates roster aggregation to the monthly rekap service
func (s *Service) RekapKelasLengkap(ctx context.Context, bulanVal, kelas string, isGuruWali bool, asuhan []authmodel.Siswa, adminID, adminRole string) (*g7model.G7RekapKelasLengkap, error) {
	return s.bulanSvc.RekapKelasLengkap(ctx, bulanVal, kelas, isGuruWali, asuhan, adminID, adminRole)
}

// GetRekapSemesterKelas delegates semester aggregation to the semester rekap service
func (s *Service) GetRekapSemesterKelas(ctx context.Context, sem, kelas string, isGuruWali bool, asuhan []authmodel.Siswa) ([]g7model.G7SemesterStudentItem, error) {
	return s.semesterSvc.GetRekapSemesterKelas(ctx, sem, kelas, isGuruWali, asuhan)
}

// UpdateRekapHariTercatat delegates daily logs count updates to the monthly rekap service
func (s *Service) UpdateRekapHariTercatat(ctx context.Context, nisn, bulanTahun string) error {
	return s.bulanSvc.UpdateRekapHariTercatat(ctx, nisn, bulanTahun)
}
