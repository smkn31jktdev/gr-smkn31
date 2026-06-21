package g7

import "be-gr31/internal/model/common"

// Coloumn Pencatatan G7
type G7 struct {
	ID            string          `json:"id" bson:"_id"`
	NISN          string          `json:"nis" bson:"nis"`
	NamaSiswa     string          `json:"namaSiswa" bson:"namaSiswa"`
	Kelas         string          `json:"kelas" bson:"kelas"`
	Tanggal       string          `json:"tanggal" bson:"tanggal"`
	Bangun        *Aktivitas      `json:"bangun" bson:"bangun"`
	Ibadah        *Aktivitas      `json:"ibadah" bson:"ibadah"`
	Makan         *Aktivitas      `json:"makan" bson:"makan"`
	Olahraga      *Aktivitas      `json:"olahraga" bson:"olahraga"`
	Belajar       *Aktivitas      `json:"belajar" bson:"belajar"`
	Bermasyarakat *Aktivitas      `json:"bermasyarakat" bson:"bermasyarakat"`
	Tidur         *Aktivitas      `json:"tidur" bson:"tidur"`
	TotalDone     int             `json:"totalDone" bson:"totalDone"`
	CreatedAt     common.FlexTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt     common.FlexTime `json:"updatedAt" bson:"updatedAt"`
}

// Aktivitas 
type Aktivitas struct {
	Done       bool   `json:"done" bson:"done"`
	Waktu      string `json:"waktu,omitempty" bson:"waktu,omitempty"`
	Keterangan string `json:"keterangan,omitempty" bson:"keterangan,omitempty"`
}

// ValidSections 
var ValidSections = []string{
	"bangun", "ibadah", "makan", "olahraga",
	"belajar", "bermasyarakat", "tidur",
}

// G7UpsertRequest
type G7UpsertRequest struct {
	Tanggal       string     `json:"tanggal" binding:"required"`
	Bangun        *Aktivitas `json:"bangun"`
	Ibadah        *Aktivitas `json:"ibadah"`
	Makan         *Aktivitas `json:"makan"`
	Olahraga      *Aktivitas `json:"olahraga"`
	Belajar       *Aktivitas `json:"belajar"`
	Bermasyarakat *Aktivitas `json:"bermasyarakat"`
	Tidur         *Aktivitas `json:"tidur"`
}

// G7Filter 
type G7Filter struct {
	NISN      string
	Kelas     string
	Tanggal   string
	BulanDari string
	BulanKe   string
	Query     string
	Page      int
	Limit     int
	AdminID   string
	AdminRole string
}

// G7Summary 
type G7Summary struct {
	NISN         string  `json:"nis"`
	NamaSiswa    string  `json:"namaSiswa"`
	Kelas        string  `json:"kelas"`
	BulanTahun   string  `json:"bulanTahun"`
	RataRataDone float64 `json:"rataRataDone"`
	HariTercatat int     `json:"hariTercatat"`
}

// endpoint GET
type G7DashboardSiswa struct {
	JurnalHariIni *G7 `json:"jurnalHariIni"`
	ProgresHariIni int `json:"progresHariIni"`

	RingkasanBulan G7RingkasanBulan `json:"ringkasanBulan"`
}

// G7RingkasanBulan 
type G7RingkasanBulan struct {
	BulanTahun   string  `json:"bulanTahun"`
	HariTercatat int     `json:"hariTercatat"`
	RataRataDone float64 `json:"rataRataDone"`
	TotalDoneSum int     `json:"totalDoneSum"`
}

// DeleteG7Request
type DeleteG7Request struct {
	ID string `json:"id" binding:"required"`
}
