package score

// ByMonthlyCount evaluates a count against standard monthly score thresholds
func ByMonthlyCount(n int) int {
	switch {
	case n > 24:
		return 5
	case n >= 15:
		return 4
	case n >= 9:
		return 3
	case n >= 4:
		return 2
	default:
		return 1
	}
}

// ByRatio evaluates the frequency ratio against standard score thresholds
func ByRatio(trueDays, total int) int {
	if total == 0 {
		return 1
	}
	r := float64(trueDays) / float64(total)
	switch {
	case r >= 0.85:
		return 5
	case r >= 0.60:
		return 4
	case r >= 0.35:
		return 3
	case r >= 0.15:
		return 2
	default:
		return 1
	}
}

// SholatFajar evaluates fajar prayer count against score thresholds
func SholatFajar(n int) int {
	switch {
	case n > 15:
		return 5
	case n >= 10:
		return 4
	case n >= 5:
		return 3
	case n >= 3:
		return 2
	default:
		return 1
	}
}

// Sholat5 evaluates 5 daily prayers count against score thresholds
func Sholat5(n int) int {
	switch {
	case n >= 28:
		return 5
	case n >= 26:
		return 4
	case n >= 25:
		return 3
	case n >= 20:
		return 2
	default:
		return 1
	}
}

// Zakat evaluates zakat/infaq sum against score thresholds
func Zakat(sum float64) int {
	switch {
	case sum > 75000:
		return 5
	case sum > 60000:
		return 4
	case sum > 50000:
		return 3
	case sum > 45000:
		return 2
	default:
		return 1
	}
}

// Makan evaluates healthy eating habits sum against score thresholds
func Makan(sum float64, hari int) int {
	if hari == 0 {
		return 1
	}
	a := sum / float64(hari)
	switch {
	case a >= 3.8:
		return 5
	case a >= 3.0:
		return 4
	case a >= 2.0:
		return 3
	case a >= 1.0:
		return 2
	default:
		return 1
	}
}

// Bermasyarakat evaluates distinct community activity counts against score thresholds
func Bermasyarakat(distinct int) int {
	switch {
	case distinct > 5:
		return 5
	case distinct > 4:
		return 4
	case distinct > 3:
		return 3
	case distinct > 2:
		return 2
	default:
		return 1
	}
}
