package sekolah

import "be-gr31/internal/model/common"

// Bukti
type Bukti struct {
	ID        string          `json:"id" bson:"_id"`
	NIS       string          `json:"nis" bson:"nis"`
	NamaSiswa string          `json:"namaSiswa" bson:"namaSiswa"`
	Kelas     string          `json:"kelas" bson:"kelas"`
	Bulan     string          `json:"bulan" bson:"bulan"`
	Foto      []string        `json:"foto" bson:"foto"`
	LinkYT    []string        `json:"linkYT" bson:"linkYT"`
	CreatedAt common.FlexTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt common.FlexTime `json:"updatedAt" bson:"updatedAt"`
}

// BuktiCreateRequest
type BuktiCreateRequest struct {
	Bulan  string   `json:"bulan" binding:"required"`
	Foto   []string `json:"foto"`
	LinkYT []string `json:"linkYT"`
}

// BuktiFilter
type BuktiFilter struct {
	NIS       string
	Kelas     string
	Bulan     string
	Page      int
	Limit     int
	AdminID   string
	AdminRole string
}

