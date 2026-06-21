package security

import (
	"strings"
	"unicode"
)
// Judol Keywords
var judolKeywords = []string{
	"slot", "s1ot", "sl0t", "s|ot",
	"togel", "toge1", "t0gel", "toto", "t0t0", "4d", "3d", "2d",
	"judi", "jud1", "j4di",
	"gacor", "g4cor", "gac0r",
	"maxwin", "m4xwin", "max-win",
	"scatter", "sc4tter",
	"jackpot", "j4ckpot",
	"rtp", "r.t.p",
	"spin", "sp1n",
	"pragmatic", "pr4gmatic",
	"pgsoft", "pg soft",
	"habanero",
	"joker123", "joker-123",
	"idn", "idnplay",
	"sbobet", "sb0bet",
	"paris", "p4ris", "pari5", "pa4ris", "p4r1s",
	"daftar", "deposit", "withdraw", "wd ",
	"slotgacor", "s1otg4c0r",
	"bonus", "b0nus",
	"promo", "pr0mo",
	"link alternatif",
	"login slot",
	"agen slot", "agen togel",
	"bandar slot", "bandar togel",
	"situs slot", "situs togel", "situs judi",
}

var judolDomains = []string{
	"slot", "togel", "poker", "casino", "betting",
	"betwin", "maxwin", "gacor", "jackpot", "spin", "paris",
}

func normalizeText(s string) string {
	var b strings.Builder
	prev := ' '
	for _, r := range strings.ToLower(s) {
		if r == '.' || r == '-' || r == '_' || r == '*' || r == '|' {
			continue
		}
		if unicode.IsSpace(r) {
			if prev != ' ' {
				b.WriteRune(' ')
			}
			prev = ' '
			continue
		}
		b.WriteRune(r)
		prev = r
	}
	return strings.TrimSpace(b.String())
}

// ContainsJudol
func ContainsJudol(text string) (bool, string) {
	norm := normalizeText(text)
	for _, kw := range judolKeywords {
		if strings.Contains(norm, kw) {
			return true, kw
		}
	}
	return false, ""
}

// ContainsJudolURL
func ContainsJudolURL(rawURL string) bool {
	lower := strings.ToLower(rawURL)
	for _, d := range judolDomains {
		if strings.Contains(lower, d) {
			return true
		}
	}
	return false
}
