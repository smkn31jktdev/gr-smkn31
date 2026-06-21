package kegiatan

import (
	"context"
	"errors"
	"time"

	"be-gr31/internal/model/common"
	sekolahmodel "be-gr31/internal/model/sekolah"
	"be-gr31/internal/util"

	"github.com/google/uuid"
)

var (
	ErrNotFound    = errors.New("kegiatan tidak ditemukan")
	ErrUnauthorized = errors.New("tidak diizinkan mengubah kegiatan ini")
)

// Service menangani business logic kegiatan.
type Service struct {
	repo *Repo
}

// NewService membuat instance Service baru.
func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// Create membuat kegiatan baru.
func (s *Service) Create(ctx context.Context, nisn, namaSiswa, kelas string, req sekolahmodel.KegiatanCreateRequest) (*sekolahmodel.Kegiatan, error) {
	if err := util.ValidateDateFormat(req.Tanggal); err != nil {
		return nil, err
	}

	now := common.FlexTime(time.Now().Format(time.RFC3339))
	data := &sekolahmodel.Kegiatan{
		ID:        uuid.New().String(),
		NISN:      nisn,
		NamaSiswa: namaSiswa,
		Kelas:     kelas,
		Tanggal:   req.Tanggal,
		Section:   req.Section,
		Payload:   req.Payload,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}

// Update memperbarui kegiatan.
func (s *Service) Update(ctx context.Context, nisn string, req sekolahmodel.KegiatanUpdateRequest) (*sekolahmodel.Kegiatan, error) {
	existing, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}
	// Pastikan hanya pemilik yang bisa update
	if existing.NISN != nisn {
		return nil, ErrUnauthorized
	}

	existing.Payload = req.Payload
	existing.UpdatedAt = common.FlexTime(time.Now().Format(time.RFC3339))

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// Delete menghapus kegiatan.
func (s *Service) Delete(ctx context.Context, nisn, id string) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	if existing.NISN != nisn {
		return ErrUnauthorized
	}
	return s.repo.Delete(ctx, id)
}

// List mengambil daftar kegiatan dengan filter + pagination.
func (s *Service) List(ctx context.Context, filter sekolahmodel.KegiatanFilter) ([]sekolahmodel.Kegiatan, bool, int, error) {
	fetcher := util.PagedFetcher[sekolahmodel.Kegiatan](func(ctx context.Context, size int, state string) ([]sekolahmodel.Kegiatan, string, error) {
		return s.repo.ListPaged(ctx, filter, size, state)
	})

	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, false, 0, err
	}

	if filter.Query != "" {
		all = util.FuzzyFilter(all, filter.Query, func(k sekolahmodel.Kegiatan) []string {
			return []string{k.NamaSiswa, k.Kelas}
		})
	}

	result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}
