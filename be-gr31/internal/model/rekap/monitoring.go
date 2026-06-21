package rekap

// SummaryPerSiswa
type SummaryPerSiswa struct {
	NIS             string  `json:"nis"`
	NamaSiswa       string  `json:"namaSiswa"`
	Kelas           string  `json:"kelas"`
	TotalHadir      int     `json:"totalHadir"`
	TotalIzin       int     `json:"totalIzin"`
	TotalSakit      int     `json:"totalSakit"`
	TotalAlpa       int     `json:"totalAlpa"`
	TotalMagang     int     `json:"totalMagang"`
	PersentaseHadir float64 `json:"persentaseHadir"`
}

// ListKehadiranAdminResponse
type ListKehadiranAdminResponse struct {
	Total            int               `json:"total"`
	Page             int               `json:"page"`
	Limit            int               `json:"limit"`
	HasMore          bool              `json:"hasMore"`
	SummaryByStudent []SummaryPerSiswa `json:"summaryByStudent"`
}

// RekapSemesterItem
type RekapSemesterItem struct {
	BulanTahun      string  `json:"bulanTahun"`
	TotalHadir      int     `json:"totalHadir"`
	TotalIzin       int     `json:"totalIzin"`
	TotalSakit      int     `json:"totalSakit"`
	TotalAlpa       int     `json:"totalAlpa"`
	TotalMagang     int     `json:"totalMagang"`
	PersentaseHadir float64 `json:"persentaseHadir"`
}

// RekapSemesterKelas
type RekapSemesterKelas struct {
	Kelas    string              `json:"kelas"`
	Semester string              `json:"semester"`
	Bulan    []RekapSemesterItem `json:"bulan"`
}

// UpdateKehadiranRequest
type UpdateKehadiranRequest struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=hadir tidak_hadir izin sakit magang"`
	Alasan string `json:"alasan"`
}

// Rekap Kelas Lengkap
// RekapSiswaItem
type RekapSiswaItem struct {
	NIS              string  `json:"nis"`
	NamaSiswa        string  `json:"namaSiswa"`
	Kelas            string  `json:"kelas"`
	TotalHadir       int     `json:"totalHadir"`
	TotalIzin        int     `json:"totalIzin"`
	TotalSakit       int     `json:"totalSakit"`
	TotalAlpa        int     `json:"totalAlpa"`
	TotalMagang      int     `json:"totalMagang"`
	HariEfektif      int     `json:"hariEfektif"`
	TingkatKehadiran float64 `json:"tingkatKehadiran"`
	AdaData          bool    `json:"adaData"`
}

// RekapKelasSummary
type RekapKelasSummary struct {
	Kelas            string  `json:"kelas"`
	TotalSiswa       int     `json:"totalSiswa"`
	HariEfektif      int     `json:"hariEfektif"`
	TotalHadir       int     `json:"totalHadir"`
	TotalIzin        int     `json:"totalIzin"`
	TotalSakit       int     `json:"totalSakit"`
	TotalAlpa        int     `json:"totalAlpa"`
	TotalMagang      int     `json:"totalMagang"`
	SiswaUnikHadir   int     `json:"siswaUnikHadir"`
	TingkatKehadiran float64 `json:"tingkatKehadiran"`
}

// RekapRange
type RekapRange struct {
	BulanTahun       string  `json:"bulanTahun"`
	HariEfektif      int     `json:"hariEfektif"`
	TotalSiswa       int     `json:"totalSiswa"`
	TotalHadir       int     `json:"totalHadir"`
	TotalIzin        int     `json:"totalIzin"`
	TotalSakit       int     `json:"totalSakit"`
	TotalAlpa        int     `json:"totalAlpa"`
	TotalMagang      int     `json:"totalMagang"`
	SiswaUnikHadir   int     `json:"siswaUnikHadir"`
	TingkatKehadiran float64 `json:"tingkatKehadiran"`
}

// RekapKelasLengkap
type RekapKelasLengkap struct {
	BulanTahun       string              `json:"bulanTahun"`
	HariEfektif      int                 `json:"hariEfektif"`
	SummaryByClass   []RekapKelasSummary `json:"summaryByClass"`
	SummaryByStudent []RekapSiswaItem    `json:"summaryByStudent"`
	SummaryRange     RekapRange          `json:"summaryRange"`

	TotalSiswa int  `json:"totalSiswa"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	HasMore    bool `json:"hasMore"`
}

// RekapPersentaseKelas
type RekapPersentaseKelas struct {
	BulanTahun     string              `json:"bulanTahun"`
	HariEfektif    int                 `json:"hariEfektif"`
	SummaryByClass []RekapKelasSummary `json:"summaryByClass"`
	SummaryRange   RekapRange          `json:"summaryRange"`
}

// RekapMingguanKelas
type RekapMingguanKelas struct {
	Senin            string              `json:"senin"`
	Jumat            string              `json:"jumat"`
	HariEfektif      int                 `json:"hariEfektif"`
	SummaryByClass   []RekapKelasSummary `json:"summaryByClass"`
	SummaryByStudent []RekapSiswaItem    `json:"summaryByStudent"`
	SummaryRange     RekapRange          `json:"summaryRange"`

	TotalSiswa int  `json:"totalSiswa"`
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	HasMore    bool `json:"hasMore"`
}
