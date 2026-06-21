package kehadiran

import (
	"errors"

	"be-gr31/internal/model/common"
)

// DeviceInfo menyimpan informasi perangkat siswa saat melakukan absensi.
type DeviceInfo struct {
	Model      string `json:"model,omitempty" bson:"model,omitempty"`           // cth: Xiaomi 23001927F
	Platform   string `json:"platform,omitempty" bson:"platform,omitempty"`     // android | ios | web
	OsVersion  string `json:"osVersion,omitempty" bson:"osVersion,omitempty"`   // cth: Android 13
	AppVersion string `json:"appVersion,omitempty" bson:"appVersion,omitempty"` // cth: 1.4.2
}

// Kehadiran Siswa
type Kehadiran struct {
	ID         string          `json:"id" bson:"_id"`
	NIS        string          `json:"nis" bson:"nis"`
	NamaSiswa  string          `json:"namaSiswa" bson:"namaSiswa"`
	Kelas      string          `json:"kelas" bson:"kelas"`
	Tanggal    string          `json:"tanggal" bson:"tanggal"`
	Hari       string          `json:"hari" bson:"hari"`
	Status     string          `json:"status" bson:"status"`
	WaktuAbsen string          `json:"waktuAbsen" bson:"waktuAbsen"`
	Alasan     string          `json:"alasan,omitempty" bson:"alasan,omitempty"`
	Koordinat  *LatLng         `json:"koordinat,omitempty" bson:"koordinat,omitempty"`
	Jarak      float64         `json:"jarak" bson:"jarak"`
	Akurasi    float64         `json:"akurasi" bson:"akurasi"`
	DeviceInfo *DeviceInfo     `json:"deviceInfo,omitempty" bson:"deviceInfo,omitempty"`
	FotoIzin   string          `json:"fotoIzin,omitempty" bson:"fotoIzin,omitempty"`
	CreatedAt  common.FlexTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt  common.FlexTime `json:"updatedAt" bson:"updatedAt"`
}

// LatLng
type LatLng struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// AbsenRequest adalah body request absensi siswa.
type AbsenRequest struct {
	Status     string      `json:"status" binding:"required,oneof=hadir tidak_hadir izin sakit magang"`
	Alasan     string      `json:"alasan"`
	Koordinat  *LatLng     `json:"koordinat"`
	Akurasi    float64     `json:"akurasi"`
	FotoIzin   string      `json:"fotoIzin"`
	Tipe       string      `json:"tipe"`
	DeviceInfo *DeviceInfo `json:"deviceInfo"`
}

// AdminAbsenRequest
type AdminAbsenRequest struct {
	NIS        string      `json:"nis" binding:"required"`
	Status     string      `json:"status" binding:"required,oneof=hadir tidak_hadir izin sakit magang"`
	Tanggal    string      `json:"tanggal" binding:"required"`
	Alasan     string      `json:"alasan"`
	Koordinat  *LatLng     `json:"koordinat"`
	Akurasi    float64     `json:"akurasi"`
	FotoIzin   string      `json:"fotoIzin"`
	DeviceInfo *DeviceInfo `json:"deviceInfo"`
}

// DeleteRequest
type DeleteRequest struct {
	ID string `json:"id" binding:"required"`
}

// KehadiranFilter
type KehadiranFilter struct {
	NIS           string
	Kelas         string
	Jurusan       string
	Tanggal       string
	TanggalDari   string
	TanggalSampai string
	BulanDari     string
	BulanKe       string
	Status        string
	Query         string
	Page          int
	Limit         int
	AdminID       string
	AdminRole     string
}

var (
	ErrDuplicate    = errors.New("absensi hari ini sudah ada")
	ErrOutOfTime    = errors.New("di luar jam absensi")
	ErrWeekend      = errors.New("tidak ada sekolah pada hari Sabtu dan Minggu")
	ErrTooFar       = errors.New("lokasi terlalu jauh dari sekolah")
	ErrNotFound     = errors.New("data tidak ditemukan")
	ErrUnauthorized = errors.New("tidak diizinkan")
	ErrBadRequest   = errors.New("request tidak valid")
)

// SiswaInfo
type SiswaInfo struct {
	NIS   string
	Nama  string
	Kelas string
}
