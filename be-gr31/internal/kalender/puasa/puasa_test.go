package puasa

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	os.Setenv("PUASA_CALENDAR_PATH", "../../../data/puasa.json")
	os.Exit(m.Run())
}

func TestLoad(t *testing.T) {
	ResetCache()
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Gagal me-load kalender puasa: %v", err)
	}
	if cfg == nil {
		t.Fatal("Config nil setelah load")
	}
	if len(cfg.Ramadan) == 0 {
		t.Error("Konfigurasi Ramadan kosong")
	}
}

func TestIsRamadanDate(t *testing.T) {
	ResetCache()

	tests := []struct {
		dateStr string
		want    bool
		wantErr bool
	}{
		// Ramadan 1447 (2026-02-19 to 2026-03-20)
		{"2026-02-18", false, false},
		{"2026-02-19", true, false},
		{"2026-03-01", true, false},
		{"2026-03-20", true, false},
		{"2026-03-21", false, false}, // Eid

		// Ramadan 1448 (2027-02-08 to 2027-03-09)
		{"2027-02-07", false, false},
		{"2027-02-08", true, false},
		{"2027-03-09", true, false},
		{"2027-03-10", false, false}, // Eid

		// Invalid format
		{"2026-03-99", false, true},
		{"invalid-date", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.dateStr, func(t *testing.T) {
			got, err := IsRamadanDate(tt.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsRamadanDate(%q) error = %v, wantErr %v", tt.dateStr, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsRamadanDate(%q) = %v, want %v", tt.dateStr, got, tt.want)
			}
		})
	}
}

func TestGetActiveRamadan(t *testing.T) {
	ResetCache()

	// Ramadan 1447 (2026-02-19 to 2026-03-20)
	period, err := GetActiveRamadan("2026-03-01")
	if err != nil {
		t.Fatalf("GetActiveRamadan failed: %v", err)
	}
	if period == nil {
		t.Fatal("GetActiveRamadan returned nil, wanted Ramadan period")
	}
	if period.HijriYear != 1447 {
		t.Errorf("expected HijriYear 1447, got %d", period.HijriYear)
	}

	// Outside Ramadan
	period, err = GetActiveRamadan("2026-06-01")
	if err != nil {
		t.Fatalf("GetActiveRamadan failed: %v", err)
	}
	if period != nil {
		t.Errorf("expected nil period outside Ramadan, got %v", period)
	}
}

func TestIsRamadanMonth(t *testing.T) {
	ResetCache()

	tests := []struct {
		monthStr string
		want     bool
		wantErr  bool
	}{
		{"2026-02", true, false}, // Overlaps: Ramadan start 2026-02-19
		{"2026-03", true, false}, // Overlaps: Ramadan end 2026-03-20
		{"2026-04", false, false},
		{"2026-06", false, false},
		{"2026-02-01", false, true}, // Invalid format YYYY-MM-DD for month check
		{"invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.monthStr, func(t *testing.T) {
			got, err := IsRamadanMonth(tt.monthStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsRamadanMonth(%q) error = %v, wantErr %v", tt.monthStr, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsRamadanMonth(%q) = %v, want %v", tt.monthStr, got, tt.want)
			}
		})
	}
}

func TestRoutes(t *testing.T) {
	ResetCache()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterRoutes(r)

	t.Run("Check Ramadan Date true", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/puasa/check?tanggal=2026-03-01", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.Code)
		}

		var body map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to parse response body: %v", err)
		}

		if body["is_ramadan"] != true {
			t.Errorf("expected is_ramadan true, got %v", body["is_ramadan"])
		}

		period, ok := body["period"].(map[string]interface{})
		if !ok || period == nil {
			t.Fatal("expected period details in response")
		}
		if period["hijri_year"] != float64(1447) {
			t.Errorf("expected Hijri year 1447, got %v", period["hijri_year"])
		}
	})

	t.Run("Check Ramadan Date false", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/puasa/check?tanggal=2026-06-01", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.Code)
		}

		var body map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to parse response body: %v", err)
		}

		if body["is_ramadan"] != false {
			t.Errorf("expected is_ramadan false, got %v", body["is_ramadan"])
		}
	})

	t.Run("Calendar", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/v1/puasa/calendar", nil)
		resp := httptest.NewRecorder()
		r.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.Code)
		}

		var body map[string]interface{}
		if err := json.Unmarshal(resp.Body.Bytes(), &body); err != nil {
			t.Fatalf("failed to parse response body: %v", err)
		}

		if body["ramadan"] == nil {
			t.Error("expected ramadan list in config response")
		}
	})
}

func TestRamadanDaysInMonth(t *testing.T) {
	ResetCache()

	tests := []struct {
		monthStr string
		want     int
		wantErr  bool
	}{
		{"2026-02", 10, false}, // Feb 19 to Feb 28 = 10 days
		{"2026-03", 20, false}, // Mar 1 to Mar 20 = 20 days
		{"2026-04", 0, false},
		{"2026-02-01", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.monthStr, func(t *testing.T) {
			got, err := RamadanDaysInMonth(tt.monthStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("RamadanDaysInMonth(%q) error = %v, wantErr %v", tt.monthStr, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RamadanDaysInMonth(%q) = %d, want %d", tt.monthStr, got, tt.want)
			}
		})
	}
}


