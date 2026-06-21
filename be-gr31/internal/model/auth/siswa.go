package auth

import (
	"be-gr31/internal/model/common"
	"encoding/json"
)

// Siswa merepresentasikan data siswa di database.
type Siswa struct {
	ID        string          `json:"id" bson:"_id"`
	NIS       string          `json:"nis" bson:"nis"`
	Nama      string          `json:"nama" bson:"nama"`
	Kelas     string          `json:"kelas" bson:"kelas"`
	Walas     string          `json:"walas" bson:"walas"`
	WaliKelas string          `json:"waliKelas" bson:"wali_kelas"`
	Password  string          `json:"password,omitempty" bson:"password"`
	IsOnline  bool            `json:"isOnline" bson:"isOnline"`
	CreatedAt common.FlexTime `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt common.FlexTime `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type SiswaAlias Siswa

// UnmarshalJSON melayani unmarshal custom untuk mendukung legacy key "nisn" dari AstraDB.
func (s *Siswa) UnmarshalJSON(data []byte) error {
	type tempSiswa struct {
		SiswaAlias
		NISN string `json:"nisn"`
	}
	var tmp tempSiswa
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	*s = Siswa(tmp.SiswaAlias)
	if s.NIS == "" && tmp.NISN != "" {
		s.NIS = tmp.NISN
	}
	return nil
}

// Request login siswa
type SiswaLoginRequest struct {
	NIS      string `json:"nis" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Respon Sukses
type SiswaLoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Siswa        Siswa  `json:"siswa"`
}

// Request refresh token
type SiswaRefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Request update data siswa
type SiswaUpdateRequest struct {
	Nama  string `json:"nama"`
	Kelas string `json:"kelas"`
	Walas string `json:"walas"`
}

// Parameter filter untuk list siswa
type SiswaFilter struct {
	Kelas     string
	Query     string
	Page      int
	Limit     int
	AdminID   string
	AdminRole string
}

// Satu item dalam import bulk siswa
type SiswaBulkItem struct {
	NIS      string `json:"nis" binding:"required"`
	Nama     string `json:"nama" binding:"required"`
	Kelas    string `json:"kelas" binding:"required"`
	Walas    string `json:"walas"`
	Password string `json:"password"` 
}
