package admin

import (
	"context"
	"errors"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	"be-gr31/internal/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound  = errors.New("admin tidak ditemukan")
	ErrDuplicate = errors.New("email sudah terdaftar")
)

// Menangani business logic manajemen admin
type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// Create menambahkan admin baru
func (s *Service) Create(ctx context.Context, req authmodel.AdminCreateRequest) (*authmodel.Admin, error) {
	existing, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicate
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := common.FlexTime(time.Now().Format(time.RFC3339))
	role := req.Role
	isWalas := req.IsWalas
	if isWalas {
		role = "walas"
	} else if role == "walas" {
		isWalas = true
	} else if role == "" {
		role = "admin"
	}

	data := &authmodel.Admin{
		ID:        uuid.New().String(),
		Nama:      req.Nama,
		Email:     req.Email,
		Role:      role,
		IsWalas:   isWalas,
		Kelas:     req.Kelas,
		Password:  string(hash),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Upsert(ctx, data); err != nil {
		return nil, err
	}

	data.Password = ""
	return data, nil
}

// Update memperbarui data admin.
func (s *Service) Update(ctx context.Context, id string, req authmodel.AdminUpdateRequest) (*authmodel.Admin, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	if req.IsWalas != nil {
		existing.IsWalas = *req.IsWalas
		if existing.IsWalas {
			existing.Role = "walas"
		} else if existing.Role == "walas" || existing.Role == "guru_wali" {
			existing.Role = "admin"
		}
	}
	if req.Kelas != nil {
		existing.Kelas = *req.Kelas
	}
	if req.Role != nil {
		existing.Role = *req.Role
		if existing.Role == "walas" || existing.Role == "guru_wali" {
			existing.IsWalas = true
		} else {
			existing.IsWalas = false
		}
	}

	if err := s.repo.UpdateFields(ctx, id, existing.IsWalas, existing.Kelas, existing.Role); err != nil {
		return nil, err
	}

	return existing, nil
}

// Menghapus admin
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// List mengambil daftar admin
func (s *Service) List(ctx context.Context, filter authmodel.AdminFilter) ([]authmodel.Admin, bool, int, error) {
	candidates, err := util.FetchAll(ctx, func(ctx context.Context, size int, state string) ([]authmodel.Admin, string, error) {
		return s.repo.ListPaged(ctx, filter.Role, size, state)
	}, 100)
	if err != nil {
		return nil, false, 0, err
	}

	if filter.Query != "" {
		candidates = util.FuzzyFilter(candidates, filter.Query, func(a authmodel.Admin) []string {
			return []string{a.Nama, a.Email}
		})
	}

	result, hasMore, total := util.PaginateSlice(candidates, filter.Page, filter.Limit)
	for i := range result {
		result[i].Password = ""
	}
	return result, hasMore, total, nil
}

// Mengimpor banyak admin sekaligus
func (s *Service) BulkCreate(ctx context.Context, items []authmodel.AdminBulkItem) (int, []string, error) {
	success := 0
	var errs []string

	for _, item := range items {
		isWalas := false
		if item.Role == "walas" || item.Role == "guru_wali" {
			isWalas = true
		}
		req := authmodel.AdminCreateRequest{
			Nama:     item.Nama,
			Email:    item.Email,
			Role:     item.Role,
			IsWalas:  isWalas,
			Password: item.Password,
		}
		if req.Password == "" {
			req.Password = "changeme123"
		}
		if _, err := s.Create(ctx, req); err != nil {
			errs = append(errs, item.Email+": "+err.Error())
		} else {
			success++
		}
	}
	return success, errs, nil
}
