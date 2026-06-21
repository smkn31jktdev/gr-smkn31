package sekolah

import "be-gr31/internal/model/common"

// Kegiatan siswa.
type Kegiatan struct {
	ID        string                 `json:"id" bson:"_id"`
	NISN      string                 `json:"nisn" bson:"nisn"`
	NamaSiswa string                 `json:"namaSiswa" bson:"namaSiswa"`
	Kelas     string                 `json:"kelas" bson:"kelas"`
	Tanggal   string                 `json:"tanggal" bson:"tanggal"`
	Section   string                 `json:"section" bson:"section"`
	Payload   map[string]interface{} `json:"payload" bson:"payload"`
	CreatedAt common.FlexTime        `json:"createdAt" bson:"createdAt"`
	UpdatedAt common.FlexTime        `json:"updatedAt" bson:"updatedAt"`
}

// KegiatanCreateRequest
type KegiatanCreateRequest struct {
	Tanggal string                 `json:"tanggal" binding:"required"`
	Section string                 `json:"section" binding:"required"`
	Payload map[string]interface{} `json:"payload" binding:"required"`
}

// KegiatanUpdateRequest
type KegiatanUpdateRequest struct {
	ID      string                 `json:"id" binding:"required"`
	Payload map[string]interface{} `json:"payload" binding:"required"`
}

// KegiatanDeleteRequest
type KegiatanDeleteRequest struct {
	ID string `json:"id" binding:"required"`
}

// KegiatanFilter
type KegiatanFilter struct {
	NISN      string
	Kelas     string
	Section   string
	Tanggal   string
	DariTgl   string
	SampaiTgl string
	Query     string
	Page      int
	Limit     int
}
