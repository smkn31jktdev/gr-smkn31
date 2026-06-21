package aduan

import (
	"context"
	"errors"
	"fmt"
	"time"

	aduanmodel "be-gr31/internal/model/aduan"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/util"

	"github.com/google/uuid"
)

var (
	ErrNotFound     = errors.New("aduan tidak ditemukan")
	ErrUnauthorized = errors.New("tidak diizinkan mengakses aduan ini")
)

// Service
type Service struct {
	repo         *Repo
	studentStore *astra.StudentStore
}

func NewService(repo *Repo, studentStore *astra.StudentStore) *Service {
	return &Service{
		repo:         repo,
		studentStore: studentStore,
	}
}

// Membuat tiket aduan baru
func (s *Service) Create(ctx context.Context, nisn, isi string) (*aduanmodel.Aduan, error) {
	student, err := s.studentStore.FindByNISN(ctx, nisn)
	if err != nil {
		return nil, fmt.Errorf("lookup student: %w", err)
	}

	namaSiswa := nisn
	kelas := ""
	walas := ""
	if student != nil {
		namaSiswa = student.Nama
		kelas = student.Kelas
		walas = student.Walas
	}

	now := util.NowJakarta()
	id := fmt.Sprintf("ADU-%s-%s", now.Format("20060102"), uuid.New().String()[:8])

	data := &aduanmodel.Aduan{
		ID:        id,
		NISN:      nisn,
		NamaSiswa: namaSiswa,
		Kelas:     kelas,
		Messages: []aduanmodel.Message{
			{
				ID:        uuid.New().String(),
				From:      namaSiswa,
				Role:      "student",
				Isi:       isi,
				Timestamp: now.Format(time.RFC3339),
			},
		},
		Status:    "pending",
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
		Walas:     walas,
		StatusHistory: []aduanmodel.StatusHistory{
			{
				Status:    "pending",
				UpdatedBy: namaSiswa,
				Role:      "student",
				UpdatedAt: now.Format(time.RFC3339),
				Note:      "",
			},
		},
	}

	if err := s.repo.Create(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}

// Respons menambahkan pesan balasan dari admin
func (s *Service) Respond(ctx context.Context, adminNama, aduanID, isi string) (*aduanmodel.Aduan, error) {
	aduan, err := s.repo.FindByID(ctx, aduanID)
	if err != nil {
		return nil, err
	}
	if aduan == nil {
		return nil, ErrNotFound
	}

	now := util.NowJakarta()
	aduan.Messages = append(aduan.Messages, aduanmodel.Message{
		ID:        uuid.New().String(),
		From:      adminNama,
		Role:      "admin",
		Isi:       isi,
		Timestamp: now.Format(time.RFC3339),
	})
	aduan.AdminNama = adminNama

	prevStatus := aduan.Status
	if aduan.Status == "open" || aduan.Status == "pending" {
		aduan.Status = "in_progress"
	}
	aduan.UpdatedAt = now.Format(time.RFC3339)

	if prevStatus != aduan.Status {
		aduan.StatusHistory = append(aduan.StatusHistory, aduanmodel.StatusHistory{
			Status:    aduan.Status,
			UpdatedBy: adminNama,
			Role:      "admin",
			UpdatedAt: now.Format(time.RFC3339),
			Note:      "Admin replied, auto-progressed",
		})
	}

	if err := s.repo.Update(ctx, aduan); err != nil {
		return nil, err
	}
	return aduan, nil
}

// Update status
func (s *Service) UpdateStatus(ctx context.Context, aduanID, status, updatedBy, role string) (*aduanmodel.Aduan, error) {
	aduan, err := s.repo.FindByID(ctx, aduanID)
	if err != nil {
		return nil, err
	}
	if aduan == nil {
		return nil, ErrNotFound
	}

	now := util.NowJakarta()
	aduan.Status = status
	aduan.UpdatedAt = now.Format(time.RFC3339)
	aduan.StatusHistory = append(aduan.StatusHistory, aduanmodel.StatusHistory{
		Status:    status,
		UpdatedBy: updatedBy,
		Role:      role,
		UpdatedAt: now.Format(time.RFC3339),
		Note:      "",
	})

	if err := s.repo.Update(ctx, aduan); err != nil {
		return nil, err
	}
	return aduan, nil
}

// List aduan dengan pagination
func (s *Service) List(ctx context.Context, filter aduanmodel.AduanFilter) ([]aduanmodel.Aduan, bool, int, error) {
	fetcher := util.PagedFetcher[aduanmodel.Aduan](func(ctx context.Context, size int, state string) ([]aduanmodel.Aduan, string, error) {
		return s.repo.ListPaged(ctx, filter, size, state)
	})

	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, false, 0, err
	}

	result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}

// ListAll fetches all complaints matching the filter without pagination limits
func (s *Service) ListAll(ctx context.Context, filter aduanmodel.AduanFilter) ([]aduanmodel.Aduan, error) {
	fetcher := util.PagedFetcher[aduanmodel.Aduan](func(ctx context.Context, size int, state string) ([]aduanmodel.Aduan, string, error) {
		return s.repo.ListPaged(ctx, filter, size, state)
	})

	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// GetRoom detail aduan
func (s *Service) GetRoom(ctx context.Context, aduanID string) (*aduanmodel.Aduan, error) {
	aduan, err := s.repo.FindByID(ctx, aduanID)
	if err != nil {
		return nil, err
	}
	if aduan == nil {
		return nil, ErrNotFound
	}
	return aduan, nil
}

