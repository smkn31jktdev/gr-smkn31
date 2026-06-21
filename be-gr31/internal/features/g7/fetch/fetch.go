package fetch

import (
	"context"

	authmodel "be-gr31/internal/model/auth"
	g7model "be-gr31/internal/model/g7"
	"be-gr31/internal/storage/astra"
	"be-gr31/internal/util"
)

// StudentRosterLister defines the interface for student profile fetching
type StudentRosterLister interface {
	ListStudentsByKelas(ctx context.Context, kelas string, pageSize int, pageState string) ([]authmodel.Siswa, string, error)
	FindStudentByNISN(ctx context.Context, nisn string) (*authmodel.Siswa, error)
}

// Service handles fetching logic for G7
type Service struct {
	rekapStore   *astra.G7RekapStore
	g7Store      *astra.G7Store
	rosterLister StudentRosterLister
}

// NewService creates a new fetch Service instance
func NewService(rekapStore *astra.G7RekapStore, g7Store *astra.G7Store, rosterLister StudentRosterLister) *Service {
	return &Service{
		rekapStore:   rekapStore,
		g7Store:      g7Store,
		rosterLister: rosterLister,
	}
}

// GetRekap fetches a monthly rekap document by student NIS and month
func (s *Service) GetRekap(ctx context.Context, nisn, bulan string) (*g7model.G7Rekap, error) {
	return s.rekapStore.FindByNISNBulan(ctx, nisn, bulan)
}

// FindStudentByNISN fetches a student profile by NISN
func (s *Service) FindStudentByNISN(ctx context.Context, nisn string) (*authmodel.Siswa, error) {
	return s.rosterLister.FindStudentByNISN(ctx, nisn)
}

// FetchAllRekap retrieves all monthly recaps matching a filter
func (s *Service) FetchAllRekap(ctx context.Context, filter g7model.G7RekapFilter) ([]g7model.G7Rekap, error) {
	fetcher := util.PagedFetcher[g7model.G7Rekap](func(ctx context.Context, size int, state string) ([]g7model.G7Rekap, string, error) {
		return s.rekapStore.ListPaged(ctx, filter, size, state)
	})
	return util.FetchAll(ctx, fetcher, 100)
}

// CountHariTercatat counts the daily journal logs filled by a student in a month
func (s *Service) CountHariTercatat(ctx context.Context, nisn, bulan string) int {
	filter := g7model.G7Filter{NISN: nisn, BulanDari: bulan, BulanKe: bulan}
	fetcher := util.PagedFetcher[g7model.G7](func(ctx context.Context, size int, state string) ([]g7model.G7, string, error) {
		return s.g7Store.ListPaged(ctx, filter, size, state)
	})
	all, err := util.FetchAll(ctx, fetcher, 100)
	if err != nil {
		return 0
	}
	return len(all)
}

// CountSiswaKelas counts total students in a class
func (s *Service) CountSiswaKelas(ctx context.Context, kelas string) int {
	count := 0
	state := ""
	for {
		items, next, err := s.rosterLister.ListStudentsByKelas(ctx, kelas, 100, state)
		if err != nil {
			break
		}
		count += len(items)
		if next == "" || len(items) == 0 {
			break
		}
		state = next
	}
	return count
}

// FetchRoster retrieves the complete list of students in a class
func (s *Service) FetchRoster(ctx context.Context, kelas string) []authmodel.Siswa {
	all := make([]authmodel.Siswa, 0)
	state := ""
	for {
		items, next, err := s.rosterLister.ListStudentsByKelas(ctx, kelas, util.MaxPageSize, state)
		if err != nil {
			break
		}
		all = append(all, items...)
		if next == "" || len(items) == 0 {
			break
		}
		state = next
	}
	return all
}

// ListRekap fetches a paginated list of monthly recaps, applying guru wali filters if applicable
func (s *Service) ListRekap(ctx context.Context, filter g7model.G7RekapFilter, isGuruWali bool, asuhan []authmodel.Siswa) ([]g7model.G7Rekap, bool, int, error) {
	var allowedNISN map[string]bool
	if isGuruWali {
		allowedNISN = make(map[string]bool)
		for _, st := range asuhan {
			if st.NIS != "" {
				allowedNISN[st.NIS] = true
			}
		}
	}

	all, err := s.FetchAllRekap(ctx, filter)
	if err != nil {
		return nil, false, 0, err
	}

	if isGuruWali {
		filtered := make([]g7model.G7Rekap, 0)
		for _, r := range all {
			if allowedNISN[r.NISN] {
				filtered = append(filtered, r)
			}
		}
		all = filtered
	}

	if filter.Query != "" {
		all = util.FuzzyFilter(all, filter.Query, func(r g7model.G7Rekap) []string {
			return []string{r.NamaSiswa, r.Kelas}
		})
	}
	result, hasMore, total := util.PaginateSlice(all, filter.Page, filter.Limit)
	return result, hasMore, total, nil
}

// FetchJournalsPaged fetches daily journals with pagination
func (s *Service) FetchJournalsPaged(ctx context.Context, filter g7model.G7Filter, pageSize int, pageState string) ([]g7model.G7, string, error) {
	return s.g7Store.ListPaged(ctx, filter, pageSize, pageState)
}
