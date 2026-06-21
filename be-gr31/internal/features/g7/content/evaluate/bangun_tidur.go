package evaluate

import (
	"fmt"

	g7model "be-gr31/internal/model/g7"
)

// Bangun evaluasi kebiasaan bangun pagi (indikator 1)
func Bangun(jurnals []g7model.G7) EvalResult {
	count0430 := 0
	count0500 := 0
	count0530 := 0
	hasData := false

	for _, j := range jurnals {
		if j.Bangun == nil || !j.Bangun.Done {
			continue
		}
		hasData = true
		w := j.Bangun.Waktu
		m := ParseKet(j.Bangun)
		berdoa := boolVal(m, "berdoa") || boolVal(m, "membacaDoa") || boolVal(m, "membacaDanBangunTidur")
		if berdoa {
			if w != "" && w <= "04:30" {
				count0430++
			} else if w != "" && w <= "05:00" {
				count0500++
			} else if w != "" && w <= "05:30" {
				count0530++
			}
		}
	}
	if !hasData {
		return EvalResult{0, "Belum ada data bangun pagi."}
	}
	n0430 := count0430
	n0500 := count0430 + count0500
	n0530 := count0430 + count0500 + count0530
	switch {
	case n0430 > 24:
		return EvalResult{5, "Istimewa: Lebih dari 24 kali bangun sebelum 04.30 dan berdoa."}
	case n0500 >= 15:
		return EvalResult{4, fmt.Sprintf("Sangat baik: %d kali bangun sebelum 05.00 dan berdoa.", n0500)}
	case n0500 >= 9:
		return EvalResult{3, fmt.Sprintf("Baik: %d kali bangun 05.00 dan berdoa.", n0500)}
	case n0500 >= 4:
		return EvalResult{2, fmt.Sprintf("Cukup baik: %d kali bangun 05.00 dan berdoa.", n0500)}
	default:
		return EvalResult{1, fmt.Sprintf("Kurang baik: Kurang dari 4 kali bangun sebelum 05.30 dan tidak berdoa. (%d kali)", n0530)}
	}
}

// Tidur evaluasi kebiasaan tidur cepat (indikator 7)
func Tidur(jurnals []g7model.G7) EvalResult {
	count := 0
	hasData := false
	for _, j := range jurnals {
		if j.Tidur == nil || !j.Tidur.Done {
			continue
		}
		hasData = true
		w := j.Tidur.Waktu
		m := ParseKet(j.Tidur)
		berdoa := boolVal(m, "berdoa") || boolVal(m, "membacaDoa") || boolVal(m, "membacaDanMasTidur")
		if berdoa && w != "" && w <= "22:00" {
			count++
		}
	}
	if !hasData {
		return EvalResult{0, "Belum ada data tidur."}
	}
	switch {
	case count > 24:
		return EvalResult{5, "Istimewa: Lebih dari 24 kali tidur sebelum jam 22.00 dan berdoa."}
	case count >= 15:
		return EvalResult{4, fmt.Sprintf("Sangat baik: %d kali tidur sebelum jam 22.00 dan berdoa.", count)}
	case count >= 9:
		return EvalResult{3, fmt.Sprintf("Baik: %d kali tidur sebelum jam 22.00 dan berdoa.", count)}
	case count >= 4:
		return EvalResult{2, fmt.Sprintf("Cukup baik: %d kali tidur sebelum jam 22.00 dan berdoa.", count)}
	default:
		return EvalResult{1, "Kurang baik: Kurang dari 4 kali tidur sebelum jam 22.00 dan berdoa."}
	}
}
