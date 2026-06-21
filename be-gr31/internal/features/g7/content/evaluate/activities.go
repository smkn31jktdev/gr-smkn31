package evaluate

import (
	"fmt"
	"regexp"
	"strings"

	g7model "be-gr31/internal/model/g7"
)

var durasiRegex = regexp.MustCompile(`durasi:\s*(\d+)`)

// Olahraga
func Olahraga(jurnals []g7model.G7, hari int) EvalResult {
	active := 0
	for _, j := range jurnals {
		if j.Olahraga != nil && j.Olahraga.Done {
			durasi := 0
			// Try to parse from keterangan
			ketLower := strings.ToLower(j.Olahraga.Keterangan)
			matches := durasiRegex.FindStringSubmatch(ketLower)
			if len(matches) > 1 {
				fmt.Sscanf(matches[1], "%d", &durasi)
			} else {
				// Fallback to parsing from waktu
				fmt.Sscanf(j.Olahraga.Waktu, "%d", &durasi)
			}

			// Olahraga is valid only if duration is >= 30 minutes
			if durasi >= 30 {
				active++
			}
		}
	}
	if active == 0 {
		return EvalResult{0, "Belum ada data olahraga."}
	}
	r := float64(active) / float64(hari)
	switch {
	case r >= 0.90:
		return EvalResult{5, fmt.Sprintf("Istimewa: Selalu berolahraga (rasio olahraga valid %d%%).", int(r*100))}
	case r >= 0.70:
		return EvalResult{4, fmt.Sprintf("Sangat baik: Sering berolahraga (rasio olahraga valid %d%%).", int(r*100))}
	case r >= 0.50:
		return EvalResult{3, fmt.Sprintf("Baik: Cukup sering berolahraga (rasio olahraga valid %d%%).", int(r*100))}
	case r >= 0.30:
		return EvalResult{2, fmt.Sprintf("Cukup baik: Kadang-kadang berolahraga (rasio olahraga valid %d%%).", int(r*100))}
	default:
		return EvalResult{1, fmt.Sprintf("Kurang baik: Jarang berolahraga (rasio olahraga valid %d%%).", int(r*100))}
	}
}

// Makan
func Makan(jurnals []g7model.G7, hari int) EvalResult {
	healthyDays := 0
	for _, j := range jurnals {
		if j.Makan == nil {
			continue
		}
		m := ParseKet(j.Makan)
		if strVal(m, "makanUtama") != "" || strVal(m, "jenisMakanan") != "" ||
			strVal(m, "laukPauk") != "" || strVal(m, "jenisLaukSayur") != "" ||
			boolVal(m, "sayurBuah") || boolVal(m, "makanSayurAtauBuah") ||
			boolVal(m, "susuSuplemen") || boolVal(m, "minumSuplemen") {
			healthyDays++
		}
	}
	if healthyDays == 0 {
		return EvalResult{0, "Belum ada data makan sehat."}
	}
	r := float64(healthyDays) / float64(hari)
	switch {
	case r >= 0.90:
		return EvalResult{5, fmt.Sprintf("Istimewa: Sangat konsisten mengonsumsi makanan sehat dan bergizi (rasio %d%%).", int(r*100))}
	case r >= 0.75:
		return EvalResult{4, fmt.Sprintf("Sangat baik: Rutin mengonsumsi makanan sehat dan bergizi (rasio %d%%).", int(r*100))}
	case r >= 0.50:
		return EvalResult{3, fmt.Sprintf("Baik: Cukup sering mengonsumsi makanan sehat dan bergizi (rasio %d%%).", int(r*100))}
	case r >= 0.25:
		return EvalResult{2, fmt.Sprintf("Cukup baik: Kadang-kadang mengonsumsi makanan sehat dan bergizi (rasio %d%%).", int(r*100))}
	default:
		return EvalResult{1, fmt.Sprintf("Kurang baik: Jarang mengonsumsi makanan sehat dan bergizi (rasio %d%%).", int(r*100))}
	}
}

// Belajar
func Belajar(jurnals []g7model.G7, hari int) EvalResult {
	done := 0
	for _, j := range jurnals {
		if j.Belajar != nil && j.Belajar.Done {
			done++
		}
	}
	if done == 0 {
		return EvalResult{0, "Belum ada data belajar."}
	}
	r := float64(done) / float64(hari)
	switch {
	case r >= 0.90:
		return EvalResult{5, fmt.Sprintf("Istimewa: Sangat aktif belajar mandiri secara rutin (rasio %d%%).", int(r*100))}
	case r >= 0.75:
		return EvalResult{4, fmt.Sprintf("Sangat baik: Sering belajar mandiri secara rutin (rasio %d%%).", int(r*100))}
	case r >= 0.50:
		return EvalResult{3, fmt.Sprintf("Baik: Cukup sering belajar mandiri (rasio %d%%).", int(r*100))}
	case r >= 0.25:
		return EvalResult{2, fmt.Sprintf("Cukup baik: Kadang-kadang belajar mandiri (rasio %d%%).", int(r*100))}
	default:
		return EvalResult{1, fmt.Sprintf("Kurang baik: Jarang belajar mandiri (rasio %d%%).", int(r*100))}
	}
}

// Masyarakat
func Masyarakat(jurnals []g7model.G7) EvalResult {
	participated := 0
	for _, j := range jurnals {
		if j.Bermasyarakat != nil && j.Bermasyarakat.Done {
			participated++
		}
	}
	if participated == 0 {
		return EvalResult{0, "Belum ada data kegiatan masyarakat."}
	}
	switch {
	case participated > 5:
		return EvalResult{5, fmt.Sprintf("Istimewa: Sangat aktif bermasyarakat (melakukan %d kegiatan).", participated)}
	case participated == 5:
		return EvalResult{4, "Sangat baik: Aktif bermasyarakat (melakukan 5 kegiatan)."}
	case participated == 4:
		return EvalResult{3, "Baik: Cukup aktif bermasyarakat (melakukan 4 kegiatan)."}
	case participated == 3:
		return EvalResult{2, "Cukup baik: Kadang-kadang mengikuti kegiatan bermasyarakat (melakukan 3 kegiatan)."}
	default:
		return EvalResult{1, fmt.Sprintf("Kurang baik: Jarang mengikuti kegiatan bermasyarakat (hanya %d kegiatan).", participated)}
	}
}
