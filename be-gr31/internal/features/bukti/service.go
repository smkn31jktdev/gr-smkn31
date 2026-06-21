package bukti

import (
	"context"
	"errors"
	"strings"
	"time"

	authmodel "be-gr31/internal/model/auth"
	"be-gr31/internal/model/common"
	sekolahmodel "be-gr31/internal/model/sekolah"
	"be-gr31/internal/util"

	"github.com/google/uuid"
)

var ErrBadRequest = errors.New("request tidak valid")
var ErrNotCurrentMonth = errors.New("hanya dapat mengunggah bukti untuk bulan laporan saat ini")

// Service menangani business logic bukti
type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

// Create menyimpan atau update bukti per bulan
func (s *Service) Create(ctx context.Context, nis, namaSiswa, kelas string, req sekolahmodel.BuktiCreateRequest) (*sekolahmodel.Bukti, error) {
	if err := util.ValidateBulanFormat(req.Bulan); err != nil {
		return nil, ErrBadRequest
	}

	currentMonth := util.BulanTahun(util.NowJakarta())
	if req.Bulan != currentMonth {
		return nil, ErrNotCurrentMonth
	}

	now := common.FlexTime(time.Now().Format(time.RFC3339))
	data := &sekolahmodel.Bukti{
		ID:        uuid.New().String(),
		NIS:       nis,
		NamaSiswa: namaSiswa,
		Kelas:     kelas,
		Bulan:     req.Bulan,
		Foto:      req.Foto,
		LinkYT:    req.LinkYT,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if data.Foto == nil {
		data.Foto = []string{}
	}
	if data.LinkYT == nil {
		data.LinkYT = []string{}
	}

	if err := s.repo.Upsert(ctx, data); err != nil {
		return nil, err
	}
	return data, nil
}

// List mengambil daftar bukti dengan filter + pagination
func (s *Service) List(ctx context.Context, filter sekolahmodel.BuktiFilter) ([]sekolahmodel.Bukti, bool, int, error) {
	var allowedNISN map[string]bool
	isGuruWali := filter.AdminRole != "" && filter.AdminRole != "super_admin"

	if isGuruWali {
		admin, err := s.repo.FindAdminByID(ctx, filter.AdminID)
		if err == nil && admin != nil {
			allowedNISN = make(map[string]bool)
			// Fetch all students
			var students []authmodel.Siswa
			state := ""
			for {
				items, next, err := s.repo.ListStudentsByKelas(ctx, "", 100, state)
				if err != nil {
					break
				}
				students = append(students, items...)
				if next == "" || len(items) == 0 {
					break
				}
				state = next
			}

			normalizedAdminName := normalizeSpace(admin.Nama)
			for _, st := range students {
				if strings.EqualFold(normalizeSpace(st.Walas), normalizedAdminName) {
					if st.NIS != "" {
						allowedNISN[st.NIS] = true
					}
				}
			}
		}
	}

	fetcher := util.PagedFetcher[sekolahmodel.Bukti](func(ctx context.Context, size int, state string) ([]sekolahmodel.Bukti, string, error) {
		return s.repo.ListPaged(ctx, filter, size, state)
	})

	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return nil, false, 0, err
	}

	if isGuruWali {
		filtered := make([]sekolahmodel.Bukti, 0)
		for _, b := range all {
			if allowedNISN != nil && allowedNISN[b.NIS] {
				filtered = append(filtered, b)
			}
		}
		all = filtered
	}

	result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}
