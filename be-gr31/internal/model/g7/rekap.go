package g7

import "be-gr31/internal/model/common"

// Status rekap bulanan G7
const (
	StatusDraft    = "draft"
	StatusReviewed = "reviewed"
	StatusFinal    = "final"
)

type G7Rekap struct {
	ID           string `json:"id" bson:"_id"`
	NISN         string `json:"nis" bson:"nis"`
	NamaSiswa    string `json:"namaSiswa" bson:"namaSiswa"`
	Kelas        string `json:"kelas" bson:"kelas"`
	BulanTahun   string `json:"bulanTahun" bson:"bulanTahun"`
	HariTercatat int    `json:"hariTercatat" bson:"hariTercatat"`

	Skor SkorG7 `json:"skor" bson:"skor"`

	NilaiMaks      int     `json:"nilaiMaks" bson:"nilaiMaks"`
	NilaiPerolehan int     `json:"nilaiPerolehan" bson:"nilaiPerolehan"`
	NilaiAkhir     float64 `json:"nilaiAkhir" bson:"nilaiAkhir"`
	Predikat       string  `json:"predikat" bson:"predikat"`

	WaliKelas string `json:"waliKelas" bson:"waliKelas"`
	GuruBK    string `json:"guruBK" bson:"guruBK"`

	Status       string          `json:"status" bson:"status"`
	TanggalFinal string          `json:"tanggalFinal,omitempty" bson:"tanggalFinal,omitempty"`
	CreatedAt    common.FlexTime `json:"createdAt" bson:"createdAt"`
	UpdatedAt    common.FlexTime `json:"updatedAt" bson:"updatedAt"`
}

// RekapKey
func (r G7Rekap) RekapKey() string { return r.NISN + "_" + r.BulanTahun }
func (r G7Rekap) JumlahAssessor() int {
	n := 0
	if r.WaliKelas != "" {
		n++
	}
	if r.GuruBK != "" {
		n++
	}
	return n
}

// G7RekapUpsertRequest
type G7RekapUpsertRequest struct {
	NISN       string `json:"nis" binding:"required"`
	BulanTahun string `json:"bulanTahun" binding:"required"`
	Skor       SkorG7 `json:"skor"`
	WaliKelas  string `json:"waliKelas"`
	GuruBK     string `json:"guruBK"`
	Status     string `json:"status"`
}

// G7RekapFilter
type G7RekapFilter struct {
	NISN       string
	Kelas      string
	BulanTahun string
	Predikat   string
	Status     string
	Query      string
	Page       int
	Limit      int
	AdminID    string
	AdminRole  string
}

// G7RekapRingkas
type G7RekapRingkas struct {
	NISN       string  `json:"nis"`
	Nama       string  `json:"nama"`
	NilaiAkhir float64 `json:"nilaiAkhir"`
}

// G7RekapStatistik
type G7RekapStatistik struct {
	Kelas                string             `json:"kelas,omitempty"`
	BulanTahun           string             `json:"bulanTahun"`
	TotalSiswa           int                `json:"totalSiswa"`
	SudahDinilai         int                `json:"sudahDinilai"`
	BelumDinilai         int                `json:"belumDinilai"`
	RataRataNilaiAkhir   float64            `json:"rataRataNilaiAkhir"`
	DistribusiPredikat   map[string]int     `json:"distribusiPredikat"`
	NilaiTertinggi       *G7RekapRingkas    `json:"nilaiTertinggi,omitempty"`
	NilaiTerendah        *G7RekapRingkas    `json:"nilaiTerendah,omitempty"`
	RataRataPerIndikator map[string]float64 `json:"rataRataPerIndikator"`
}

// G7SuggestResponse
type G7SuggestResponse struct {
	NISN         string            `json:"nis"`
	BulanTahun   string            `json:"bulanTahun"`
	Skor         SkorG7            `json:"skor"`
	Catatan      map[string]string `json:"catatan"`
	HariTercatat int               `json:"hariTercatat"`
	IsAdvisory   bool              `json:"isAdvisory"`
}

// DeleteG7RekapRequest
type DeleteG7RekapRequest struct {
	ID string `json:"id" binding:"required"`
}

// ── Rekap Kelas Lengkap (roster-join) ─────────────────────────────────────────

// G7RekapSiswaItem baris rekap nilai G7 satu siswa.
// Siswa yang belum dinilai tetap muncul dengan sudahDinilai=false.
type G7RekapSiswaItem struct {
	NISN           string  `json:"nis"`
	NamaSiswa      string  `json:"namaSiswa"`
	Kelas          string  `json:"kelas"`
	NilaiPerolehan int     `json:"nilaiPerolehan"`
	NilaiMaks      int     `json:"nilaiMaks"`
	NilaiAkhir     float64 `json:"nilaiAkhir"`
	Predikat       string  `json:"predikat"`
	Status         string  `json:"status"`
	HariTercatat   int     `json:"hariTercatat"`
	SudahDinilai   bool    `json:"sudahDinilai"`
}

// G7RekapKelasLengkap response GET /v1/admin/g7/rekap-kelas.
// Menggabungkan rekap nilai G7 (g7_rekap) dengan roster siswa, plus statistik
// distribusi predikat — semua dalam satu panggilan untuk memudahkan rekap admin.
type G7RekapKelasLengkap struct {
	Kelas        string             `json:"kelas"`
	BulanTahun   string             `json:"bulanTahun"`
	TotalSiswa   int                `json:"totalSiswa"`
	SudahDinilai int                `json:"sudahDinilai"`
	BelumDinilai int                `json:"belumDinilai"`
	Statistik    *G7RekapStatistik  `json:"statistik"`
	Siswa        []G7RekapSiswaItem `json:"siswa"`
}

// G7SemesterStudentItem baris rekap nilai G7 satu siswa untuk satu semester.
type G7SemesterStudentItem struct {
	NISN        string             `json:"nis"`
	NamaSiswa   string             `json:"namaSiswa"`
	Kelas       string             `json:"kelas"`
	NilaiAkhir  float64            `json:"nilaiAkhir"`
	Predikat    string             `json:"predikat"`
	MonthsCount int                `json:"monthsCount"`
	Skor        map[string]float64 `json:"skor"`
}

// G7SemesterKelasResponse response GET /v1/admin/g7/rekap-semester.
type G7SemesterKelasResponse struct {
	Semester string                  `json:"semester"`
	Kelas    string                  `json:"kelas,omitempty"`
	Siswa    []G7SemesterStudentItem `json:"siswa"`
}

