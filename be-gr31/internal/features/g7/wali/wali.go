package wali

import (
	"context"
	"errors"
	"log"
	"strings"

	authmodel "be-gr31/internal/model/auth"
)

type AdminFinder interface {
	FindAdminByID(ctx context.Context, id string) (*authmodel.Admin, error)
}

type RosterFetcher interface {
	FetchRoster(ctx context.Context, kelas string) []authmodel.Siswa
}

type Service struct {
	adminFinder   AdminFinder
	rosterFetcher RosterFetcher
}

func NewService(adminFinder AdminFinder, rosterFetcher RosterFetcher) *Service {
	return &Service{
		adminFinder:   adminFinder,
		rosterFetcher: rosterFetcher,
	}
}

func normalizeSpace(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

// Cek Guru Wali
func (s *Service) CheckIsGuruWali(ctx context.Context, adminID, adminRole string) (bool, []authmodel.Siswa, error) {
	role := strings.ToLower(strings.TrimSpace(adminRole))
	if role == "super_admin" {
		return false, nil, nil
	}

	if adminID == "" {
		return true, []authmodel.Siswa{}, nil
	}

	admin, err := s.adminFinder.FindAdminByID(ctx, adminID)
	if err != nil || admin == nil {
		log.Printf("WARNING: admin tidak ditemukan atau error saat lookup: ID=%s, err=%v", adminID, err)
		return true, []authmodel.Siswa{}, nil
	}

	// Ambil seluruh siswa lintas semua halaman
	students := s.rosterFetcher.FetchRoster(ctx, "")

	normalizedAdminName := normalizeSpace(admin.Nama)
	var asuhan []authmodel.Siswa
	for _, st := range students {
		if strings.EqualFold(normalizeSpace(st.Walas), normalizedAdminName) {
			asuhan = append(asuhan, st)
		}
	}

	return true, asuhan, nil
}

// Verifikasi Admin
func (s *Service) VerifyAdminAccessToStudent(ctx context.Context, adminID, adminRole, studentNIS string) error {
	role := strings.ToLower(strings.TrimSpace(adminRole))
	if role == "super_admin" {
		return nil
	}

	isGuruWali, asuhan, err := s.CheckIsGuruWali(ctx, adminID, adminRole)
	if err != nil {
		return err
	}

	if isGuruWali {
		for _, st := range asuhan {
			if st.NIS == studentNIS {
				return nil
			}
		}
	}

	return errors.New("forbidden: anda tidak memiliki akses ke data siswa ini")
}
