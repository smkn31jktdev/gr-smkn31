package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var reDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var reTime = regexp.MustCompile(`^\d{2}:\d{2}$`)
var reBulan = regexp.MustCompile(`^\d{4}-\d{2}$`)

// ValidateDateFormat memvalidasi format YYYY-MM-DD.
func ValidateDateFormat(s string) error {
	if !reDate.MatchString(s) {
		return errors.New("format tanggal harus YYYY-MM-DD")
	}
	_, err := time.Parse("2006-01-02", s)
	if err != nil {
		return errors.New("tanggal tidak valid")
	}
	return nil
}

// ValidateTimeFormat memvalidasi format HH:mm.
func ValidateTimeFormat(s string) error {
	if s == "" {
		return nil
	}
	if !reTime.MatchString(s) {
		return errors.New("format waktu harus HH:mm")
	}
	return nil
}

// ValidateBulanFormat memvalidasi format YYYY-MM.
func ValidateBulanFormat(s string) error {
	if !reBulan.MatchString(s) {
		return errors.New("format bulan harus YYYY-MM")
	}
	return nil
}

// ValidateNotFutureDate memvalidasi bahwa tanggal tidak di masa depan (WIB).
func ValidateNotFutureDate(tanggal string) error {
	t, err := time.ParseInLocation("2006-01-02", tanggal, JakartaLoc())
	if err != nil {
		return errors.New("tanggal tidak valid")
	}
	now := NowJakarta()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, JakartaLoc())
	if t.After(today) {
		return errors.New("tidak bisa input untuk tanggal mendatang")
	}
	return nil
}

// ValidateMaxLength memvalidasi panjang string maksimum.
func ValidateMaxLength(s string, max int, field string) error {
	if len([]rune(s)) > max {
		return errors.New(field + " maksimal " + strconv.Itoa(max) + " karakter")
	}
	return nil
}

// SanitizeFilename membersihkan nama file dari karakter berbahaya.
func SanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, "/", "")
	name = strings.ReplaceAll(name, "\\", "")
	name = strings.ReplaceAll(name, "..", "")
	return strings.TrimSpace(name)
}
