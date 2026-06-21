package g7

// PDFRowIndikator satu baris dalam tabel PDF laporan G7
type PDFRowIndikator struct {
	No         int    `json:"no"`
	Indikator  string `json:"indikator"`
	Skor       int    `json:"skor"`       // 1–5 (0 = belum diisi)
	Keterangan string `json:"keterangan"` // note dari evaluasi
}

// PDFLaporan data lengkap untuk render PDF laporan G7
type PDFLaporan struct {
	NamaSiswa    string            `json:"namaSiswa"`
	NIS          string            `json:"nis"`
	Kelas        string            `json:"kelas"`
	BulanTahun   string            `json:"bulanTahun"`
	TahunAjaran  string            `json:"tahunAjaran"`
	Indikator    []PDFRowIndikator `json:"indikator"`
	WaliKelas    string            `json:"waliKelas"`
	OrangTua     string            `json:"orangTua"`
	NilaiAkhir   float64           `json:"nilaiAkhir"`
	Predikat     string            `json:"predikat"`
	TanggalCetak string            `json:"tanggalCetak"`
}
