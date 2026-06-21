package kegiatan

import (
	"context"

	sekolahmodel "be-gr31/internal/model/sekolah"
	"be-gr31/internal/storage/astra"
)

// Repo menangani operasi storage untuk kegiatan.
type Repo struct {
	store *astra.KegiatanStore
}

// NewRepo membuat instance Repo baru.
func NewRepo(store *astra.KegiatanStore) *Repo {
	return &Repo{store: store}
}

func (r *Repo) Create(ctx context.Context, data *sekolahmodel.Kegiatan) error {
	return r.store.Create(ctx, data)
}

func (r *Repo) Update(ctx context.Context, data *sekolahmodel.Kegiatan) error {
	return r.store.Update(ctx, data)
}

func (r *Repo) FindByID(ctx context.Context, id string) (*sekolahmodel.Kegiatan, error) {
	return r.store.FindByID(ctx, id)
}

func (r *Repo) Delete(ctx context.Context, id string) error {
	return r.store.Delete(ctx, id)
}

func (r *Repo) ListPaged(ctx context.Context, filter sekolahmodel.KegiatanFilter, pageSize int, pageState string) ([]sekolahmodel.Kegiatan, string, error) {
	return r.store.ListPaged(ctx, filter, pageSize, pageState)
}
