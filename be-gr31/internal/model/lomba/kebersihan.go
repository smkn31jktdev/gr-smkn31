package lomba

import "be-gr31/internal/model/common"

// KebersihanKelas merepresentasikan data unggahan kebersihan kelas.
type KebersihanKelas struct {
	ID        string          `json:"id" bson:"_id"`
	Kelas     string          `json:"kelas" bson:"kelas"`
	NISN      string          `json:"nisn" bson:"nisn"`
	NamaSiswa string          `json:"namaSiswa" bson:"namaSiswa"`
	Tanggal   string          `json:"tanggal" bson:"tanggal"`
	Foto      []string        `json:"foto" bson:"foto"`
	Catatan   string          `json:"catatan" bson:"catatan"`
	CreatedAt common.FlexTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt common.FlexTime `json:"updatedAt" bson:"updatedAt"`
}

// KebersihanCreateRequest adalah body request unggah data kebersihan kelas baru.
type KebersihanCreateRequest struct {
	Tanggal string   `json:"tanggal" binding:"required"`
	Foto    []string `json:"foto" binding:"required,min=1"`
	Catatan string   `json:"catatan"`
}

// KebersihanUpdateRequest adalah body request perbarui data kebersihan kelas.
type KebersihanUpdateRequest struct {
	ID      string   `json:"id" binding:"required"`
	Foto    []string `json:"foto" binding:"required,min=1"`
	Catatan string   `json:"catatan"`
}

// KebersihanFilter adalah filter untuk list data kebersihan kelas.
type KebersihanFilter struct {
	Kelas     string
	Tanggal   string
	DariTgl   string
	SampaiTgl string
	Page      int
	Limit     int
	AdminID   string
	AdminRole string
}
