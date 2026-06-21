package lomba

import (
	"context"
	"errors"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	lombamodel "be-gr31/internal/model/lomba"
	"be-gr31/internal/util"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("data kebersihan tidak ditemukan")
	ErrUnauthorized  = errors.New("tidak diizinkan mengubah data kebersihan ini")
	ErrAlreadyExists = errors.New("kelas Anda sudah mengirimkan foto kebersihan untuk minggu ini")
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetStudentInfo(ctx context.Context, nisn string) (*authmodel.Siswa, error) {
	return s.repo.FindStudentByNISN(ctx, nisn)
}

func (s *Service) Create(ctx context.Context, nisn string, req lombamodel.KebersihanCreateRequest) (*lombamodel.KebersihanKelas, error) {
	if err := util.ValidateDateFormat(req.Tanggal); err != nil {
		return nil, err
	}

	// Fetch student profile to get actual Nama and Kelas
	student, err := s.repo.FindStudentByNISN(ctx, nisn)
	if err != nil {
		return nil, err
	}
	if student == nil {
		return nil, errors.New("siswa tidak ditemukan")
	}

	// Parse date to check the week range
	t, err := util.ParseDate(req.Tanggal)
	if err != nil {
		return nil, err
	}
	mondayStr, sundayStr := util.GetWeekRange(t)

	// Enforce: only one upload per class per week
	existing, err := s.repo.FindByKelasDateRange(ctx, student.Kelas, mondayStr, sundayStr)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, ErrAlreadyExists
	}

	now := common.FlexTime(time.Now().Format(time.RFC3339))
	data := &lombamodel.KebersihanKelas{
		ID:        uuid.New().String(),
		Kelas:     student.Kelas,
		NISN:      nisn,
		NamaSiswa: student.Nama,
		Tanggal:   req.Tanggal,
		Foto:      req.Foto,
		Catatan:   req.Catatan,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Service) Update(ctx context.Context, nisn string, req lombamodel.KebersihanUpdateRequest) (*lombamodel.KebersihanKelas, error) {
	existing, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	// Only the uploader can edit their submission
	if existing.NISN != nisn {
		return nil, ErrUnauthorized
	}

	existing.Foto = req.Foto
	existing.Catatan = req.Catatan
	existing.UpdatedAt = common.FlexTime(time.Now().Format(time.RFC3339))

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) Delete(ctx context.Context, nisn, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}

	// Only the uploader can delete their submission
	if existing.NISN != nisn {
		return ErrUnauthorized
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context, filter lombamodel.KebersihanFilter) ([]lombamodel.KebersihanKelas, bool, int, error) {
	isGuruWali := filter.AdminRole != "" && filter.AdminRole != "super_admin"
	if isGuruWali {
		admin, err := s.repo.FindAdminByID(ctx, filter.AdminID)
		if err == nil && admin != nil {
			if admin.IsWalas || admin.Role == "walas" || admin.Role == "guru_wali" {
				filter.Kelas = admin.Kelas
			}
		} else {
			return nil, false, 0, nil
		}
	}

	fetcher := util.PagedFetcher[lombamodel.KebersihanKelas](func(ctx context.Context, size int, state string) ([]lombamodel.KebersihanKelas, string, error) {
		return s.repo.ListPaged(ctx, filter, size, state)
	})

	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, false, 0, err
	}

	result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}
