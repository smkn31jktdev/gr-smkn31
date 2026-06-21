package rekap

// RekapBulanan
type RekapBulanan struct {
	ID               string  `json:"id" bson:"_id"`
	RekapKey         string  `json:"rekapKey" bson:"rekap_key"` // NIS_YYYY-MM
	NIS              string  `json:"nis" bson:"nis"`
	NamaSiswa        string  `json:"namaSiswa" bson:"nama"`
	Kelas            string  `json:"kelas" bson:"kelas"`
	BulanTahun       string  `json:"bulanTahun" bson:"bulan"`            // YYYY-MM
	Semester         string  `json:"semester,omitempty" bson:"semester"` // "ganjil"|"genap"
	TotalHadir       int     `json:"totalHadir" bson:"total_hadir"`
	TotalIzin        int     `json:"totalIzin" bson:"total_izin"`
	TotalSakit       int     `json:"totalSakit" bson:"total_sakit"`
	TotalTidakHadir  int     `json:"totalTidakHadir" bson:"total_tidak_hadir"`
	TotalMagang      int     `json:"totalMagang" bson:"total_magang"`
	TotalHariEfektif int     `json:"totalHariEfektif" bson:"total_hari_efektif"`
	PersentaseHadir  float64 `json:"persentaseHadir" bson:"persentase_hadir"`
	CreatedAt        string  `json:"createdAt,omitempty" bson:"created_at"`
	UpdatedAt        string  `json:"updatedAt,omitempty" bson:"updated_at"`
}

// KehadiranBulananSiswa
type KehadiranBulananSiswa struct {
	BulanTahun string                `json:"bulanTahun"`
	Kehadiran  []KehadiranHariItem   `json:"kehadiran"`
	Summary    KehadiranBulanSummary `json:"summary"`
}

// KehadiranHariItem
type KehadiranHariItem struct {
	Tanggal    string `json:"tanggal"`
	Status     string `json:"status"`
	WaktuAbsen string `json:"waktuAbsen"`
	Alasan     string `json:"alasan,omitempty"`
	FotoIzin   string `json:"fotoIzin,omitempty"`
}

// KehadiranBulanSummary
type KehadiranBulanSummary struct {
	TotalHadir       int     `json:"totalHadir"`
	TotalIzin        int     `json:"totalIzin"`
	TotalSakit       int     `json:"totalSakit"`
	TotalAlpa        int     `json:"totalAlpa"`
	TotalMagang      int     `json:"totalMagang"`
	TotalHariEfektif int     `json:"totalHariEfektif"`
	PersentaseHadir  float64 `json:"persentaseHadir"`
}

// RekapFilter
type RekapFilter struct {
	NIS        string
	Kelas      string
	BulanTahun string
	Semester   string
	Query      string
	Page       int
	Limit      int
}

// IncrementField
type IncrementField string

const (
	FieldHadir       IncrementField = "totalHadir"
	FieldIzin        IncrementField = "totalIzin"
	FieldSakit       IncrementField = "totalSakit"
	FieldTidakHadir  IncrementField = "totalTidakHadir"
	FieldMagang      IncrementField = "totalMagang"
	FieldHariEfektif IncrementField = "totalHariEfektif"
)

// StatusToField
func StatusToField(status string) IncrementField {
	switch status {
	case "hadir":
		return FieldHadir
	case "izin":
		return FieldIzin
	case "sakit":
		return FieldSakit
	case "magang":
		return FieldMagang
	default:
		return FieldTidakHadir
	}
}

// HitungSemester
func HitungSemester(bulan int) string {
	if bulan >= 7 {
		return "ganjil"
	}
	return "genap"
}

// RekapHarian
type RekapHarian struct {
	Tanggal     string `json:"tanggal"`
	TotalHadir  int    `json:"totalHadir"`
	TotalIzin   int    `json:"totalIzin"`
	TotalSakit  int    `json:"totalSakit"`
	TotalAlpa   int    `json:"totalAlpa"`
	TotalMagang int    `json:"totalMagang"`
	TotalSiswa  int    `json:"totalSiswa"`
}

// RekapMingguan
type RekapMingguan struct {
	MingguKe     int           `json:"mingguKe"`
	TanggalMulai string        `json:"tanggalMulai"`
	TanggalAkhir string        `json:"tanggalAkhir"`
	Hari         []RekapHarian `json:"hari"`
	Rangkuman    RekapHarian   `json:"rangkuman"`
}

// RingkasanSiswa
type RingkasanSiswa struct {
	NIS             string  `json:"nis"`
	NamaSiswa       string  `json:"namaSiswa"`
	Kelas           string  `json:"kelas"`
	BulanTahun      string  `json:"bulanTahun"`
	TotalHadir      int     `json:"totalHadir"`
	TotalIzin       int     `json:"totalIzin"`
	TotalSakit      int     `json:"totalSakit"`
	TotalTidakHadir int     `json:"totalTidakHadir"`
	TotalMagang     int     `json:"totalMagang"`
	PersentaseHadir float64 `json:"persentaseHadir"`
}

// Kelas-Jurusan
type KelasJurusanResponse struct {
	Kelas        []string `json:"kelas"`
	Jurusan      []string `json:"jurusan"`
	KelasLengkap []string `json:"kelasLengkap"`
}
