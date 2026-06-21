package aduan

import (
	"context"

	aduanmodel "be-gr31/internal/model/aduan"
	"be-gr31/internal/storage/mongodb"
)

// Repo menangani operasi storage untuk aduan.
type Repo struct {
	store *mongodb.AduanStore
}

func NewRepo(store *mongodb.AduanStore) *Repo {
	return &Repo{store: store}
}

func (r *Repo) Create(ctx context.Context, data *aduanmodel.Aduan) error {
	return r.store.Create(ctx, data)
}

func (r *Repo) Update(ctx context.Context, data *aduanmodel.Aduan) error {
	return r.store.Update(ctx, data)
}

func (r *Repo) FindByID(ctx context.Context, id string) (*aduanmodel.Aduan, error) {
	return r.store.FindByID(ctx, id)
}

func (r *Repo) ListPaged(ctx context.Context, filter aduanmodel.AduanFilter, pageSize int, pageState string) ([]aduanmodel.Aduan, string, error) {
	return r.store.ListPaged(ctx, filter, pageSize, pageState)
}
