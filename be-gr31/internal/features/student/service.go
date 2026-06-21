package student

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	"be-gr31/internal/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNotFound  = errors.New("siswa tidak ditemukan")
	ErrDuplicate = errors.New("NIS sudah terdaftar")
)

// Service
type Service struct {
	repo *Repo
	rdb  cacheInvalidator
}

// cacheInvalidator
type cacheInvalidator interface {
	InvalidateStudentCache(ctx context.Context, kelas string) error
}

// simpleRedisInvalidator
type simpleRedisInvalidator struct{}

func (s *simpleRedisInvalidator) InvalidateStudentCache(ctx context.Context, kelas string) error {
	return nil
}

// NewService
func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
		rdb:  &simpleRedisInvalidator{},
	}
}

// Mnambahkan siswa baru
func (s *Service) Create(ctx context.Context, req authmodel.SiswaBulkItem) (*authmodel.Siswa, error) {
	existing, err := s.repo.FindByNIS(ctx, req.NIS)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrDuplicate
	}

	// Hash password
	password := req.Password
	if password == "" {
		password = "123456"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := common.FlexTime(time.Now().Format(time.RFC3339))
	data := &authmodel.Siswa{
		ID:        uuid.New().String(),
		NIS:       req.NIS,
		Nama:      req.Nama,
		Kelas:     req.Kelas,
		Walas:     req.Walas,
		Password:  string(hash),
		IsOnline:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Upsert(ctx, data); err != nil {
		return nil, err
	}

	data.Password = ""
	return data, nil
}

// Memperbarui data siswa
func (s *Service) Update(ctx context.Context, id string, req authmodel.SiswaUpdateRequest) (*authmodel.Siswa, error) {
	siswa, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if siswa == nil {
		return nil, ErrNotFound
	}

	oldKelas := siswa.Kelas
	if req.Nama != "" {
		siswa.Nama = req.Nama
	}
	if req.Kelas != "" {
		siswa.Kelas = req.Kelas
	}
	if req.Walas != "" {
		siswa.Walas = req.Walas
	}
	siswa.UpdatedAt = common.FlexTime(time.Now().Format(time.RFC3339))

	if err := s.repo.Upsert(ctx, siswa); err != nil {
		return nil, err
	}

	// Invalidasi cache
	s.rdb.InvalidateStudentCache(ctx, oldKelas)
	if req.Kelas != "" && req.Kelas != oldKelas {
		s.rdb.InvalidateStudentCache(ctx, req.Kelas)
	}

	siswa.Password = ""
	return siswa, nil
}

// Menghapus siswa
func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// List dengan filter + fuzzy + pagination
func (s *Service) List(ctx context.Context, filter authmodel.SiswaFilter) ([]authmodel.Siswa, bool, int, error) {
	if s.repo.UseSupabase() {
		isGuruWali, err := s.isGuruWaliRole(ctx, filter.AdminID, filter.AdminRole)
		if err != nil {
			return nil, false, 0, err
		}

		students, total, err := s.repo.ListSupabase(ctx, filter, isGuruWali, filter.AdminID)
		if err != nil {
			return nil, false, 0, err
		}

		hasMore := false
		if filter.Limit > 0 {
			hasMore = (filter.Page * filter.Limit) < total
		}
		return students, hasMore, total, nil
	}

	isGuruWali, asuhan, err := s.checkIsGuruWali(ctx, filter.AdminID, filter.AdminRole)
	if err != nil {
		return nil, false, 0, err
	}

	var candidates []authmodel.Siswa
	if isGuruWali {
		candidates = asuhan
		if filter.Kelas != "" {
			var filtered []authmodel.Siswa
			for _, st := range candidates {
				if st.Kelas == filter.Kelas {
					filtered = append(filtered, st)
				}
			}
			candidates = filtered
		}
	} else {
		candidates, err = util.FetchAll(ctx, func(ctx context.Context, size int, state string) ([]authmodel.Siswa, string, error) {
			return s.repo.ListPaged(ctx, filter.Kelas, size, state)
		}, 100)
		if err != nil {
			return nil, false, 0, err
		}
	}

	if filter.Query != "" {
		candidates = util.FuzzyFilter(candidates, filter.Query, func(s authmodel.Siswa) []string {
			return []string{s.Nama, s.Kelas, s.NIS}
		})
	}

	result, hasMore, total := util.PaginateSlice(candidates, filter.Page, filter.Limit)
	// Hapus password dari semua result
	for i := range result {
		result[i].Password = ""
	}
	return result, hasMore, total, nil
}

// Mengimpor banyak siswa sekaligus
func (s *Service) BulkCreate(ctx context.Context, items []authmodel.SiswaBulkItem) (int, []string, error) {
	success := 0
	var errs []string

	for _, item := range items {
		if _, err := s.Create(ctx, item); err != nil {
			errs = append(errs, item.NIS+": "+err.Error())
		} else {
			success++
		}
	}
	return success, errs, nil
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func (s *Service) checkIsGuruWali(ctx context.Context, adminID, adminRole string) (bool, []authmodel.Siswa, error) {
	role := strings.ToLower(strings.TrimSpace(adminRole))
	if role == "super_admin" {
		return false, nil, nil
	}

	if adminID == "" {
		return true, []authmodel.Siswa{}, nil
	}

	admin, err := s.repo.FindAdminByID(ctx, adminID)
	if err != nil || admin == nil {
		log.Printf("WARNING: admin tidak ditemukan atau error saat lookup: ID=%s, err=%v", adminID, err)
		return true, []authmodel.Siswa{}, nil
	}

	// Ambil seluruh siswa dengan auto-paging
	var students []authmodel.Siswa
	state := ""
	for {
		items, next, err := s.repo.ListPaged(ctx, "", 100, state)
		if err != nil {
			return false, nil, err
		}
		students = append(students, items...)
		if next == "" || len(items) == 0 {
			break
		}
		state = next
	}

	normalizedAdminName := normalizeSpace(admin.Nama)
	var asuhan []authmodel.Siswa
	for _, st := range students {
		if strings.EqualFold(normalizeSpace(st.Walas), normalizedAdminName) {
			asuhan = append(asuhan, st)
		}
	}

	return true, asuhan, nil
}

func (s *Service) isGuruWaliRole(ctx context.Context, adminID, adminRole string) (bool, error) {
	role := strings.ToLower(strings.TrimSpace(adminRole))
	if role == "super_admin" {
		return false, nil
	}
	return true, nil
}
