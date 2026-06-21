package izin

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrBadRequest = errors.New("request tidak valid")
)

// Validate izin/sakit requirements
func Validate(status, alasan, fotoIzin string) error {
	if status == "izin" || status == "sakit" {
		if alasan == "" {
			return fmt.Errorf("alasan wajib diisi untuk status %s", status)
		}
		if fotoIzin == "" {
			return fmt.Errorf("foto izin wajib ada untuk status %s", status)
		}
	}
	return nil
}

// IsSakit checks if the reason points to sick keywords
func IsSakit(alasan string) bool {
	alasanLower := strings.ToLower(alasan)
	keywordsSakit := []string{"sakit", "demam", "dokter", "pusing", "pilek", "batuk", "gigi", "rawat", "klinik", "puskesmas", "opname", "rs", "rumah sakit"}
	for _, kw := range keywordsSakit {
		if strings.Contains(alasanLower, kw) {
			return true
		}
	}
	return false
}

// IsIzin checks if the reason points to permit keywords
func IsIzin(alasan string) bool {
	alasanLower := strings.ToLower(alasan)
	keywordsIzin := []string{"izin", "ijin", "acara", "keperluan", "nikah", "melayat", "keluarga", "kondangan", "pulang", "pergi", "halangan"}
	for _, kw := range keywordsIzin {
		if strings.Contains(alasanLower, kw) {
			return true
		}
	}
	return false
}
