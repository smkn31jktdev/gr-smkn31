package puasa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RamadanPeriod struct {
	HijriYear     int    `json:"hijri_year"`
	GregorianYear int    `json:"gregorian_year"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	EidDate       string `json:"eid_date"`
}

type PuasaConfig struct {
	Source  string          `json:"source"`
	Note    string          `json:"note"`
	Ramadan []RamadanPeriod `json:"ramadan"`
}

var (
	once    sync.Once
	config  *PuasaConfig
	loadErr error
)

func candidatePaths() []string {
	paths := []string{}
	if p := strings.TrimSpace(os.Getenv("PUASA_CALENDAR_PATH")); p != "" {
		paths = append(paths, p)
	}
	paths = append(paths,
		filepath.Join("data", "puasa.json"),
		filepath.Join(".", "data", "puasa.json"),
		filepath.Join("..", "data", "puasa.json"),
	)
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		paths = append(paths,
			filepath.Join(dir, "data", "puasa.json"),
			filepath.Join(dir, "..", "data", "puasa.json"),
			filepath.Join(dir, "..", "..", "data", "puasa.json"),
		)
	}
	return paths
}

// Load loads the puasa config from puasa.json dynamically
func Load() (*PuasaConfig, error) {
	once.Do(func() {
		var raw []byte
		var err error
		for _, p := range candidatePaths() {
			raw, err = os.ReadFile(p)
			if err == nil {
				break
			}
		}
		if err != nil {
			loadErr = fmt.Errorf("failed to read puasa.json: %w", err)
			return
		}
		var cfg PuasaConfig
		if err := json.Unmarshal(raw, &cfg); err != nil {
			loadErr = fmt.Errorf("failed to parse puasa.json: %w", err)
			return
		}
		config = &cfg
	})
	return config, loadErr
}

// ResetCache clears the loaded config cache, primarily for unit tests
func ResetCache() {
	once = sync.Once{}
	config = nil
	loadErr = nil
}

// IsRamadanDate checks if the given YYYY-MM-DD date falls within Ramadan
func IsRamadanDate(dateStr string) (bool, error) {
	cfg, err := Load()
	if err != nil {
		return false, err
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return false, fmt.Errorf("invalid date format YYYY-MM-DD: %w", err)
	}
	for _, p := range cfg.Ramadan {
		start, err1 := time.Parse("2006-01-02", p.StartDate)
		end, err2 := time.Parse("2006-01-02", p.EndDate)
		if err1 != nil || err2 != nil {
			continue
		}
		if !t.Before(start) && !t.After(end) {
			return true, nil
		}
	}
	return false, nil
}

// GetActiveRamadan retrieves the active Ramadan period details for a given date
func GetActiveRamadan(dateStr string) (*RamadanPeriod, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		return nil, fmt.Errorf("invalid date format YYYY-MM-DD: %w", err)
	}
	for _, p := range cfg.Ramadan {
		start, err1 := time.Parse("2006-01-02", p.StartDate)
		end, err2 := time.Parse("2006-01-02", p.EndDate)
		if err1 != nil || err2 != nil {
			continue
		}
		if !t.Before(start) && !t.After(end) {
			return &p, nil
		}
	}
	return nil, nil
}

// IsRamadanMonth checks if any day within the month (YYYY-MM) overlaps with Ramadan
func IsRamadanMonth(monthStr string) (bool, error) {
	cfg, err := Load()
	if err != nil {
		return false, err
	}
	monthStr = strings.TrimSpace(monthStr)
	if len(monthStr) != 7 || monthStr[4] != '-' {
		return false, fmt.Errorf("invalid month format, must be YYYY-MM")
	}
	monthStart, err := time.Parse("2006-01-02", monthStr+"-01")
	if err != nil {
		return false, fmt.Errorf("invalid month format: %w", err)
	}
	monthEnd := monthStart.AddDate(0, 1, -1)

	for _, p := range cfg.Ramadan {
		start, err1 := time.Parse("2006-01-02", p.StartDate)
		end, err2 := time.Parse("2006-01-02", p.EndDate)
		if err1 != nil || err2 != nil {
			continue
		}
		// Overlaps if start <= monthEnd && end >= monthStart
		if !start.After(monthEnd) && !end.Before(monthStart) {
			return true, nil
		}
	}
	return false, nil
}

// RamadanDaysInMonth returns the number of calendar days in monthStr (YYYY-MM) that are within Ramadan
func RamadanDaysInMonth(monthStr string) (int, error) {
	_, err := Load()
	if err != nil {
		return 0, err
	}
	monthStr = strings.TrimSpace(monthStr)
	if len(monthStr) != 7 || monthStr[4] != '-' {
		return 0, fmt.Errorf("invalid month format, must be YYYY-MM")
	}
	monthStart, err := time.Parse("2006-01-02", monthStr+"-01")
	if err != nil {
		return 0, fmt.Errorf("invalid month format: %w", err)
	}
	monthEnd := monthStart.AddDate(0, 1, -1)

	count := 0
	for d := monthStart; !d.After(monthEnd); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		isRamadan, _ := IsRamadanDate(dateStr)
		if isRamadan {
			count++
		}
	}
	return count, nil
}

// RegisterRoutes registers the puasa routes into Gin engine
func RegisterRoutes(r *gin.Engine) {
	group := r.Group("/v1/puasa")
	{
		group.GET("/check", handleCheck)
		group.GET("/calendar", handleCalendar)
	}
}

func handleCheck(c *gin.Context) {
	tanggal := c.Query("tanggal")
	if tanggal == "" {
		tanggal = time.Now().Format("2006-01-02")
	}

	isRamadan, err := IsRamadanDate(tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	period, _ := GetActiveRamadan(tanggal)

	ramadanDay := 0
	if isRamadan && period != nil {
		start, err1 := time.Parse("2006-01-02", period.StartDate)
		curr, err2 := time.Parse("2006-01-02", tanggal)
		if err1 == nil && err2 == nil {
			diff := curr.Sub(start)
			ramadanDay = int(diff.Hours()/24) + 1
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"tanggal":     tanggal,
		"is_ramadan":  isRamadan,
		"ramadan_day": ramadanDay,
		"period":      period,
	})
}

func handleCalendar(c *gin.Context) {
	cfg, err := Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}
