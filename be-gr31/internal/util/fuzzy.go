package util

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// FuzzyThreshold
const FuzzyThreshold = 0.35

// Trigram memecah string menjadi set trigram
func Trigram(s string) map[string]bool {
	s = normalizeFuzzy(s)
	trigrams := make(map[string]bool)
	runes := []rune(s)
	for i := 0; i <= len(runes)-3; i++ {
		trigrams[string(runes[i:i+3])] = true
	}
	return trigrams
}

// TrigramSimilarity menghitung kesamaan dua string berdasarkan trigram (0.0–1.0)
func TrigramSimilarity(a, b string) float64 {
	ta := Trigram(a)
	tb := Trigram(b)
	if len(ta) == 0 && len(tb) == 0 {
		return 1.0
	}
	if len(ta) == 0 || len(tb) == 0 {
		return 0.0
	}
	intersection := 0
	for t := range ta {
		if tb[t] {
			intersection++
		}
	}
	return float64(intersection) / float64(len(ta)+len(tb)-intersection)
}

// JaroSimilarity menghitung Jaro similarity (0.0–1.0)
func JaroSimilarity(a, b string) float64 {
	a = normalizeFuzzy(a)
	b = normalizeFuzzy(b)
	if a == b {
		return 1.0
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	matchDist := int(math.Max(float64(len(a)), float64(len(b)))/2) - 1
	if matchDist < 0 {
		matchDist = 0
	}

	aMatches := make([]bool, len(a))
	bMatches := make([]bool, len(b))
	matches := 0
	transpositions := 0

	for i := range a {
		start := int(math.Max(0, float64(i-matchDist)))
		end := int(math.Min(float64(i+matchDist+1), float64(len(b))))
		for j := start; j < end; j++ {
			if bMatches[j] || a[i] != b[j] {
				continue
			}
			aMatches[i] = true
			bMatches[j] = true
			matches++
			break
		}
	}

	if matches == 0 {
		return 0.0
	}

	k := 0
	for i := range a {
		if !aMatches[i] {
			continue
		}
		for !bMatches[k] {
			k++
		}
		if a[i] != b[k] {
			transpositions++
		}
		k++
	}

	return (float64(matches)/float64(len(a)) +
		float64(matches)/float64(len(b)) +
		float64(matches-transpositions/2)/float64(matches)) / 3
}

// JaroWinkler menambahkan prefix bonus ke Jaro similarity
func JaroWinkler(a, b string) float64 {
	jaro := JaroSimilarity(a, b)
	prefix := 0
	maxPfx := int(math.Min(4, math.Min(float64(len(a)), float64(len(b)))))
	for i := 0; i < maxPfx; i++ {
		if a[i] == b[i] {
			prefix++
		} else {
			break
		}
	}
	return jaro + float64(prefix)*0.1*(1-jaro)
}

// FuzzyScore menggabungkan trigram + Jaro-Winkler dengan bobot
func FuzzyScore(query, target string) float64 {
	const (
		weightTrigram     = 0.45
		weightJaroWinkler = 0.55
	)
	tScore := TrigramSimilarity(query, target)
	jScore := JaroWinkler(query, target)
	return weightTrigram*tScore + weightJaroWinkler*jScore
}

// FuzzyMatch
func FuzzyMatch(query, target string) bool {
	return FuzzyScore(query, target) >= FuzzyThreshold
}

// FuzzyResult adalah hasil satu item dengan score.
type FuzzyResult[T any] struct {
	Item  T
	Score float64
}

// FuzzyFilter
func FuzzyFilter[T any](items []T, query string, fieldSelector func(T) []string) []T {
	if query == "" {
		return items
	}
	query = normalizeFuzzy(query)

	results := make([]FuzzyResult[T], 0, len(items))
	for _, item := range items {
		fields := fieldSelector(item)
		bestScore := 0.0
		for _, field := range fields {
			if s := FuzzyScore(query, field); s > bestScore {
				bestScore = s
			}
		}
		if bestScore >= FuzzyThreshold {
			results = append(results, FuzzyResult[T]{Item: item, Score: bestScore})
		}
	}

	// Sort descending by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	out := make([]T, len(results))
	for i, r := range results {
		out[i] = r.Item
	}
	return out
}

// normalizeFuzzy
func normalizeFuzzy(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
