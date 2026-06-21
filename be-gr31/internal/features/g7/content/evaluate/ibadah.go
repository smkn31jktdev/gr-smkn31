package evaluate

import (
	"be-gr31/internal/kalender/puasa"
	g7model "be-gr31/internal/model/g7"
)

// Beribadah evaluates the Beribadah habit (indicator 2) based on the updated AGENTS.md daily score rules.
func Beribadah(jurnals []g7model.G7) EvalResult {
	totalDailyScore := 0.0
	validDays := 0
	hasData := false

	for _, j := range jurnals {
		if j.Ibadah == nil || !j.Ibadah.Done {
			continue
		}
		hasData = true
		validDays++
		m := ParseKet(j.Ibadah)

		// 1. Ibadah Wajib lengkap = +1.0 poin
		// Wajib: Sholat Fajar & Sholat 5 waktu berjamaah
		wajibScore := 0.0
		if boolVal(m, "sholatFajar") {
			wajibScore += 0.5
		}
		if boolVal(m, "sholat5Waktu") {
			wajibScore += 0.5
		}

		// 2. Ibadah Sunnah lengkap = +1.0 poin (proporsional)
		// Sunnah: Dhuha, Rawatib, Zikir
		sunnahCount := 0
		if boolVal(m, "sholatDhuha") {
			sunnahCount++
		}
		if boolVal(m, "rowatib") || boolVal(m, "sholatSunah") {
			sunnahCount++
		}
		if boolVal(m, "zikir") {
			sunnahCount++
		}

		isRam, _ := puasa.IsRamadanDate(j.Tanggal)
		divisor := 2.0
		if isRam {
			divisor = 3.0
		}
		sunnahScore := float64(sunnahCount) / divisor

		// 3. Berdoa untuk diri sendiri dan orang tua = +1.0 poin
		doaScore := 0.0
		if boolVal(m, "berdoa") {
			doaScore = 1.0
		}

		// 4. Sedekah/infaq diisi = +1.0 poin
		sedekahScore := 0.0
		infaqVal := floatVal(m, "infaq")
		zakatVal := floatVal(m, "zakat")
		if infaqVal > 0 || zakatVal > 0 {
			sedekahScore = 1.0
		} else {
			if boolVal(m, "sedekah") || boolVal(m, "infaqBool") {
				sedekahScore = 1.0
			}
		}

		dailyScore := wajibScore + sunnahScore + doaScore + sedekahScore
		totalDailyScore += dailyScore
	}

	if !hasData || validDays == 0 {
		return EvalResult{0, "Belum ada data ibadah."}
	}

	avgScore := totalDailyScore / float64(validDays)

	var score int
	var note string
	switch {
	case avgScore >= 3.5:
		score = 5
		note = "Selalu melaksanakan ibadah wajib dan sunnah secara konsisten."
	case avgScore >= 2.5:
		score = 4
		note = "Sering melaksanakan ibadah wajib dan sebagian ibadah sunnah."
	case avgScore >= 1.5:
		score = 3
		note = "Melaksanakan ibadah wajib dengan cukup rutin, namun ibadah sunnah masih belum konsisten."
	case avgScore >= 0.5:
		score = 2
		note = "Melaksanakan sebagian ibadah wajib dan sunnah, namun frekuensinya tidak teratur dan masih perlu ditingkatkan."
	default:
		score = 1
		note = "Jarang melaksanakan ibadah wajib maupun sunnah."
	}

	return EvalResult{score, note}
}
