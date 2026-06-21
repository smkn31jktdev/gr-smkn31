package util

import (
	"be-gr31/internal/kalender"
	"fmt"
	"time"
)

// Lokasi timezone Jakarta (WIB, UTC+7).
var jakartaLoc *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Fallback jika timezone data tidak tersedia
		loc = time.FixedZone("WIB", 7*60*60)
	}
	jakartaLoc = loc
}

// NowJakarta mengembalikan waktu saat ini dalam timezone Jakarta.
func NowJakarta() time.Time {
	return time.Now().In(jakartaLoc)
}

// JakartaLoc mengembalikan lokasi timezone Jakarta.
func JakartaLoc() *time.Location {
	return jakartaLoc
}

// ValidateAbsenTime memvalidasi apakah waktu berada dalam rentang absensi.
func ValidateAbsenTime(t time.Time, startHour, startMin, endHour, endMin int) error {
	h, m, _ := t.Clock()
	totalMin := h*60 + m
	startTotal := startHour*60 + startMin
	endTotal := endHour*60 + endMin

	if totalMin < startTotal || totalMin > endTotal {
		return fmt.Errorf("di luar jam absensi (%02d:%02d–%02d:%02d WIB)",
			startHour, startMin, endHour, endMin)
	}
	return nil
}

// FormatTanggal memformat time.Time menjadi string YYYY-MM-DD.
func FormatTanggal(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatWaktu memformat time.Time menjadi HH:mm:ss.
func FormatWaktu(t time.Time) string {
	return t.Format("15:04:05")
}

// FormatDateTime memformat time.Time menjadi RFC3339.
func FormatDateTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// NamaHari mengembalikan nama hari dalam Bahasa Indonesia.
func NamaHari(t time.Time) string {
	days := map[time.Weekday]string{
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
		time.Sunday:    "Minggu",
	}
	return days[t.Weekday()]
}

// Mengembalikan true jika hari adalah Sabtu atau Minggu
func IsWeekend(t time.Time) bool {
	wd := t.Weekday()
	return wd == time.Saturday || wd == time.Sunday
}

// Mengembalikan string format YYYY-MM dari Time
func BulanTahun(t time.Time) string {
	return t.Format("2006-01")
}

// Parse string YYYY-MM-DD menjadi Time dalam timezone Jakarta
func ParseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, jakartaLoc)
}

// Mengembalikan tanggal Senin (awal minggu) dan Minggu (akhir minggu)
func GetWeekRange(t time.Time) (string, string) {
	wd := t.Weekday()
	offset := int(wd) - 1
	if wd == time.Sunday {
		offset = 6
	}
	monday := t.AddDate(0, 0, -offset)
	sunday := monday.AddDate(0, 0, 6)
	return monday.Format("2006-01-02"), sunday.Format("2006-01-02")
}

// Hari Efektif
func CountHariEfektifBulan(bulan string) int {
	return kalender.HariEfektif(bulan)
}

// Rentang hari Efektif
func CountHariEfektifRentang(dari, sampai string) (int, error) {
	start, err := time.Parse("2006-01-02", dari)
	if err != nil {
		return 0, fmt.Errorf("format dari harus YYYY-MM-DD")
	}
	end, err := time.Parse("2006-01-02", sampai)
	if err != nil {
		return 0, fmt.Errorf("format sampai harus YYYY-MM-DD")
	}
	if end.Before(start) {
		return 0, fmt.Errorf("tanggal sampai tidak boleh sebelum tanggal dari")
	}
	return countWeekdays(start, end), nil
}

// countWeekdays menghitung hari Senin–Jumat dalam rentang inklusif.
func countWeekdays(start, end time.Time) int {
	count := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		wd := d.Weekday()
		if wd != time.Saturday && wd != time.Sunday {
			count++
		}
	}
	return count
}
