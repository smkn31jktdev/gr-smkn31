package kehadiran

import (
	"context"

	kehadiranmodel "be-gr31/internal/model/kehadiran"
)

type KehadiranStore interface {
	Create(ctx context.Context, data *kehadiranmodel.Kehadiran) error
	FindByID(ctx context.Context, id string) (*kehadiranmodel.Kehadiran, error)
	FindByNISNTanggal(ctx context.Context, nisn, tanggal string) (*kehadiranmodel.Kehadiran, error)
	ListPaged(ctx context.Context, filter kehadiranmodel.KehadiranFilter, pageSize int, pageState string) ([]kehadiranmodel.Kehadiran, string, error)
	Update(ctx context.Context, id string, fields map[string]any) error
	Delete(ctx context.Context, id string) error
}

// Menangani operasi storage untuk kehadiran
type Repo struct {
	store KehadiranStore
}

// Membuat instance Repo baru
func NewRepo(store KehadiranStore) *Repo {
	return &Repo{store: store}
}

// Menyimpan data kehadiran baru
func (r *Repo) Create(ctx context.Context, data *kehadiranmodel.Kehadiran) error {
	return r.store.Create(ctx, data)
}

// Mencari kehadiran berdasarkan NISN + tanggal
func (r *Repo) FindByNISNTanggal(ctx context.Context, nisn, tanggal string) (*kehadiranmodel.Kehadiran, error) {
	return r.store.FindByNISNTanggal(ctx, nisn, tanggal)
}

// Mengambil kehadiran berdasarkan ID
func (r *Repo) FindByID(ctx context.Context, id string) (*kehadiranmodel.Kehadiran, error) {
	return r.store.FindByID(ctx, id)
}

// Menghapus kehadiran berdasarkan ID
func (r *Repo) Delete(ctx context.Context, id string) error {
	return r.store.Delete(ctx, id)
}

// Mengambil kehadiran dengan pagination cursor
func (r *Repo) ListPaged(ctx context.Context, filter kehadiranmodel.KehadiranFilter, pageSize int, pageState string) ([]kehadiranmodel.Kehadiran, string, error) {
	return r.store.ListPaged(ctx, filter, pageSize, pageState)
}

// Mengupdate field tertentu pada dokumen kehadiran
func (r *Repo) Update(ctx context.Context, id string, fields map[string]any) error {
	return r.store.Update(ctx, id, fields)
}
