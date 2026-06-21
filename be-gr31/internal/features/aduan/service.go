package aduan

import (
	"context"
	"errors"
	"fmt"
	"time"

	aduanmodel "be-gr31/internal/model/aduan"
	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/storage/supabase"
	"be-gr31/internal/util"

	"github.com/google/uuid"
)

var (
	ErrNotFound     = errors.New("aduan tidak ditemukan")
	ErrUnauthorized = errors.New("tidak diizinkan mengakses aduan ini")
)

// Service
type Service struct {
	repo           *Repo
	studentStore   *astra.StudentStore
	adminStore     *astra.AdminStore
	supabaseClient *supabase.Client
}

func NewService(repo *Repo, studentStore *astra.StudentStore, adminStore *astra.AdminStore, supabaseClient *supabase.Client) *Service {
	return &Service{
		repo:           repo,
		studentStore:   studentStore,
		adminStore:     adminStore,
		supabaseClient: supabaseClient,
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
		Wali:      walas,
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

func (s *Service) findAdminByID(ctx context.Context, id string) (*authmodel.Admin, error) {
	if s.supabaseClient != nil && s.supabaseClient.DB != nil {
		var a authmodel.Admin
		var createdAt time.Time
		query := `
			SELECT id, nama, email, password, is_walas, kelas, role, created_at
			FROM akun_admin
			WHERE id = $1
		`
		var isWalas bool
		err := s.supabaseClient.DB.QueryRowContext(ctx, query, id).Scan(&a.ID, &a.Nama, &a.Email, &a.Password, &isWalas, &a.Kelas, &a.Role, &createdAt)
		if err != nil {
			return nil, err
		}
		a.IsWalas = isWalas
		if a.Role == "" {
			if isWalas {
				a.Role = "walas"
			} else {
				a.Role = "admin"
			}
		}
		return &a, nil
	}
	return s.adminStore.FindByID(ctx, id)
}

func (s *Service) applyWaliFilter(ctx context.Context, filter *aduanmodel.AduanFilter) {
	if filter.AdminRole == "" {
		return
	}
	isSpecialRole := filter.AdminRole == "super_admin" || filter.AdminRole == "guru_bk" || filter.AdminRole == "bk" || filter.AdminRole == "admin_bk"
	if !isSpecialRole && filter.AdminID != "" {
		admin, err := s.findAdminByID(ctx, filter.AdminID)
		if err == nil && admin != nil {
			filter.Wali = admin.Nama
		}
	}
}

// List aduan dengan pagination
func (s *Service) List(ctx context.Context, filter aduanmodel.AduanFilter) ([]aduanmodel.Aduan, bool, int, error) {
	s.applyWaliFilter(ctx, &filter)

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
	s.applyWaliFilter(ctx, &filter)

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

