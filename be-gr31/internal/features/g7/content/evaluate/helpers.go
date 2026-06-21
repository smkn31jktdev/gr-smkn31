package evaluate

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	g7model "be-gr31/internal/model/g7"
)

var infaqRegex = regexp.MustCompile(`(?:infaq|zakat|rp)[\s:]*(?:rp)?[\s:]*(\d+)`)

func ParseKet(a *g7model.Aktivitas) map[string]any {
	if a == nil || a.Keterangan == "" {
		return map[string]any{}
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(a.Keterangan), &m); err == nil {
		return m
	}

	m = make(map[string]any)
	ketLower := strings.ToLower(a.Keterangan)

	// 1. Check for berdoa/doa (Bangun/Tidur)
	if strings.Contains(ketLower, "doa: ya") || strings.Contains(ketLower, "doa: true") || strings.Contains(ketLower, "doa:ya") {
		m["berdoa"] = true
		m["membacaDoa"] = true
		m["membacaDanBangunTidur"] = true
		m["membacaDanMasTidur"] = true
	} else if strings.Contains(ketLower, "doa: tidak") || strings.Contains(ketLower, "doa: false") || strings.Contains(ketLower, "doa:tidak") {
		m["berdoa"] = false
		m["membacaDoa"] = false
		m["membacaDanBangunTidur"] = false
		m["membacaDanMasTidur"] = false
	}

	// 2. Check for sholat 5 waktu (Ibadah)
	if strings.Contains(ketLower, "sholat 5 waktu") || strings.Contains(ketLower, "5 waktu") {
		m["sholat5Waktu"] = true
	}

	// 3. Check for sholat fajar / qoblal subuh
	if strings.Contains(ketLower, "sholat fajar") || strings.Contains(ketLower, "qoblal subuh") || strings.Contains(ketLower, "fajar") {
		m["sholatFajar"] = true
	}

	// 4. Check for zikir
	if strings.Contains(ketLower, "zikir") {
		m["zikir"] = true
	}

	// 5. Check for sholat dhuha
	if strings.Contains(ketLower, "dhuha") {
		m["sholatDhuha"] = true
	}

	// 6. Check for rowatib / rawatib
	if strings.Contains(ketLower, "rowatib") || strings.Contains(ketLower, "rawatib") {
		m["rowatib"] = true
		m["sholatSunah"] = true
	}

	// 7. Parse infaq/zakat nominals
	matches := infaqRegex.FindStringSubmatch(ketLower)
	if len(matches) > 1 {
		var val float64
		if _, err := fmt.Sscanf(matches[1], "%f", &val); err == nil {
			m["infaq"] = val
			m["zakat"] = val
		}
	}

	// 7b. Parse puasa
	if strings.Contains(ketLower, "puasa: ya") || strings.Contains(ketLower, "puasa: true") || strings.Contains(ketLower, "puasa:ya") {
		m["puasa"] = true
	} else if strings.Contains(ketLower, "puasa: tidak") || strings.Contains(ketLower, "puasa: false") || strings.Contains(ketLower, "puasa:tidak") {
		m["puasa"] = false
	}

	// 7c. Parse tarawih
	tarawihRegex := regexp.MustCompile(`tarawih:\s*(\d+)`)
	tMatches := tarawihRegex.FindStringSubmatch(ketLower)
	if len(tMatches) > 1 {
		var val int
		if _, err := fmt.Sscanf(tMatches[1], "%d", &val); err == nil {
			m["tarawihRokaat"] = val
		}
	}

	// 8. Makan parsing
	if strings.Contains(ketLower, "makan utama") || strings.Contains(ketLower, "makanan utama") {
		m["makanUtama"] = "Makan"
	}
	if strings.Contains(ketLower, "lauk") {
		m["laukPauk"] = "Lauk"
		m["jenisLaukSayur"] = "Lauk"
	}
	if strings.Contains(ketLower, "sayur/buah: ya") || strings.Contains(ketLower, "sayur/buah:ya") || strings.Contains(ketLower, "sayur: ya") || strings.Contains(ketLower, "sayur:ya") || strings.Contains(ketLower, "buah: ya") || strings.Contains(ketLower, "buah:ya") {
		m["sayurBuah"] = true
		m["makanSayurAtauBuah"] = true
	}
	if strings.Contains(ketLower, "susu/suplemen: ya") || strings.Contains(ketLower, "susu/suplemen:ya") || strings.Contains(ketLower, "susu: ya") || strings.Contains(ketLower, "susu:ya") || strings.Contains(ketLower, "suplemen: ya") || strings.Contains(ketLower, "suplemen:ya") {
		m["susuSuplemen"] = true
		m["minumSuplemen"] = true
	}

	return m
}

func boolVal(m map[string]any, key string) bool {
	b, _ := m[key].(bool)
	return b
}

func strVal(m map[string]any, key string) string {
	s, _ := m[key].(string)
	return s
}

func floatVal(m map[string]any, key string) float64 {
	switch t := m[key].(type) {
	case float64:
		return t
	case int:
		return float64(t)
	}
	return 0
}

