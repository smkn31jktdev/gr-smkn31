package kalender

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Setenv("KALENDER_PENDIDIKAN_PATH", "../../data/kalender-pendidikan.json")
	os.Exit(m.Run())
}

func TestHariEfektifDariKalender(t *testing.T) {
	path := os.Getenv("KALENDER_PENDIDIKAN_PATH")
	if path == "" {
		path = "../../data/kalender-pendidikan.json"
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Gagal membaca file kalender: %v", err)
	}

	var kf kalenderFile
	if err := json.Unmarshal(raw, &kf); err != nil {
		t.Fatalf("Gagal unmarshal file kalender: %v", err)
	}

	for _, b := range kf.RincianBulan {
		mn := namaBulan[strings.ToLower(strings.TrimSpace(b.Bulan))]
		if mn == 0 || b.Tahun == 0 {
			continue
		}
		bulanKey := monthKey(b.Tahun, mn)

		t.Run(bulanKey, func(t *testing.T) {
			got := HariEfektif(bulanKey)
			if got != b.HariMasukEfektif {
				t.Errorf("HariEfektif(%q) = %d, want %d (dari JSON)", bulanKey, got, b.HariMasukEfektif)
			}
		})
	}
}

func TestHariEfektifFallbackWeekday(t *testing.T) {
	if got := HariEfektif("2026-07"); got != 23 {
		t.Errorf("HariEfektif(2026-07) fallback = %d, want 23", got)
	}
}

func TestIsLibur(t *testing.T) {
	path := os.Getenv("KALENDER_PENDIDIKAN_PATH")
	if path == "" {
		path = "../../data/kalender-pendidikan.json"
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Gagal membaca file kalender: %v", err)
	}

	var kf kalenderFile
	if err := json.Unmarshal(raw, &kf); err != nil {
		t.Fatalf("Gagal unmarshal file kalender: %v", err)
	}

	allLibur := make(map[string]bool)
	for _, b := range kf.RincianBulan {
		for _, d := range expandLibur(b.DaftarLibur) {
			allLibur[d] = true
			if !IsLibur(d) {
				t.Errorf("IsLibur(%q) = false, want true (dari JSON)", d)
			}
		}
	}

	bukanLibur := []string{"2026-06-02", "2026-06-16", "2026-06-18"}
	for _, d := range bukanLibur {
		if _, isLiburInJson := allLibur[d]; !isLiburInJson {
			if IsLibur(d) {
				t.Errorf("IsLibur(%q) = true, want false (tidak ada di JSON)", d)
			}
		}
	}
}

func TestIsHariEfektif(t *testing.T) {
	path := os.Getenv("KALENDER_PENDIDIKAN_PATH")
	if path == "" {
		path = "../../data/kalender-pendidikan.json"
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Gagal membaca file kalender: %v", err)
	}

	var kf kalenderFile
	if err := json.Unmarshal(raw, &kf); err != nil {
		t.Fatalf("Gagal unmarshal file kalender: %v", err)
	}

	allLibur := make(map[string]bool)
	for _, b := range kf.RincianBulan {
		for _, d := range expandLibur(b.DaftarLibur) {
			allLibur[d] = true
		}
	}

	for _, b := range kf.RincianBulan {
		mn := namaBulan[strings.ToLower(strings.TrimSpace(b.Bulan))]
		if mn == 0 || b.Tahun == 0 {
			continue
		}

		for day := 1; day <= 28; day++ {
			dateStr := fmt.Sprintf("%04d-%02d-%02d", b.Tahun, mn, day)
			parsedTime, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				continue
			}

			isWeekend := parsedTime.Weekday() == time.Saturday || parsedTime.Weekday() == time.Sunday
			_, isLibur := allLibur[dateStr]

			expectedEfektif := !isWeekend && !isLibur
			gotEfektif := IsHariEfektif(dateStr)

			if gotEfektif != expectedEfektif {
				t.Errorf("IsHariEfektif(%q) = %t, want %t (Weekend: %t, Libur: %t)", dateStr, gotEfektif, expectedEfektif, isWeekend, isLibur)
			}
		}
	}
}
