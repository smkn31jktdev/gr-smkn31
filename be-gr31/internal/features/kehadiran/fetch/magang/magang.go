package magang

import (
	"strings"
)

// IsMagang checks if the reason points to internship/PKL keywords
func IsMagang(alasan string) bool {
	alasanLower := strings.ToLower(alasan)
	keywordsMagang := []string{"magang", "pkl", "prakerin", "kerja praktek"}
	for _, kw := range keywordsMagang {
		if strings.Contains(alasanLower, kw) {
			return true
		}
	}
	return false
}
