package g7

import (
	"be-gr31/internal/kalender/puasa"
	"errors"
	"fmt"
	"math"
)

const NilaiMaks = 90
var ErrSkorRange = errors.New("skor sub-indikator harus bernilai 0..5")
var ErrBadRequest = errors.New("request tidak valid")

type SkorG7 struct {
	// Seq 1 — Kebiasaan 1
	BangunPagi int `json:"bangunPagi" bson:"bangunPagi"`

	// Seq 2–10 — Kebiasaan 2
	IbadahDoa          int `json:"ibadahDoa" bson:"ibadahDoa"`
	IbadahSholatFajar  int `json:"ibadahSholatFajar" bson:"ibadahSholatFajar"`   // col E
	IbadahSholat5Waktu int `json:"ibadahSholat5Waktu" bson:"ibadahSholat5Waktu"` // col F
	IbadahZikir        int `json:"ibadahZikir" bson:"ibadahZikir"`               // col G
	IbadahDhuha        int `json:"ibadahDhuha" bson:"ibadahDhuha"`               // col H
	IbadahRowatib      int `json:"ibadahRowatib" bson:"ibadahRowatib"`           // col I
	IbadahTarawih      int `json:"ibadahTarawih" bson:"ibadahTarawih"`
	IbadahPuasa        int `json:"ibadahPuasa" bson:"ibadahPuasa"`
	IbadahZakat        int `json:"ibadahZakat" bson:"ibadahZakat"`

	// Seq 11 — Kebiasaan 3
	Olahraga int `json:"olahraga" bson:"olahraga"`

	// Seq 12 — Kebiasaan 4
	MakanSehat int `json:"makanSehat" bson:"makanSehat"`

	// Seq 13–16 — Kebiasaan 5
	BelajarKitabSuci int `json:"belajarKitabSuci" bson:"belajarKitabSuci"` // col O
	BelajarBukuUmum  int `json:"belajarBukuUmum" bson:"belajarBukuUmum"`
	BelajarBukuMapel int `json:"belajarBukuMapel" bson:"belajarBukuMapel"`
	BelajarTugas     int `json:"belajarTugas" bson:"belajarTugas"`

	// Seq 17 — Kebiasaan 6
	Bermasyarakat int `json:"bermasyarakat" bson:"bermasyarakat"`

	// Seq 18 — Kebiasaan 7
	TidurCepat int `json:"tidurCepat" bson:"tidurCepat"`
}

// indikatorKey
var indikatorKey = []string{
	"bangunPagi", "ibadahDoa", "ibadahSholatFajar", "ibadahSholat5Waktu",
	"ibadahZikir", "ibadahDhuha", "ibadahRowatib", "ibadahTarawih",
	"ibadahPuasa", "ibadahZakat", "olahraga", "makanSehat",
	"belajarKitabSuci", "belajarBukuUmum", "belajarBukuMapel", "belajarTugas",
	"bermasyarakat", "tidurCepat",
}

// ordered mengembalikan ke-18 skor dalam urutan kolom Excel C–T.
func (s SkorG7) ordered() []int {
	return []int{
		s.BangunPagi, s.IbadahDoa, s.IbadahSholatFajar,
		s.IbadahSholat5Waktu, s.IbadahZikir, s.IbadahDhuha,
		s.IbadahRowatib, s.IbadahTarawih, s.IbadahPuasa,
		s.IbadahZakat, s.Olahraga, s.MakanSehat,
		s.BelajarKitabSuci, s.BelajarBukuUmum, s.BelajarBukuMapel,
		s.BelajarTugas, s.Bermasyarakat, s.TidurCepat,
	}
}

// AsMap
func (s SkorG7) AsMap() map[string]int {
	ord := s.ordered()
	m := make(map[string]int, len(ord))
	for i, key := range indikatorKey {
		m[key] = ord[i]
	}
	return m
}

// ValidateSkor
func ValidateSkor(skor SkorG7) error {
	for i, v := range skor.ordered() {
		if v < 0 || v > 5 {
			return fmt.Errorf("%w (indikator %q = %d)", ErrSkorRange, indikatorKey[i], v)
		}
	}
	return nil
}

// HitungNilaiAkhir calculates scores and dynamically sets max points based on the indicators scored and Ramadan period
func HitungNilaiAkhir(skor SkorG7, bulanTahun string) (perolehan int, maks int, akhir float64, predikat string) {
	islamOnly := map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true, 7: true, 8: true}
	all := skor.ordered()

	// Detect if Islam indicators are used by checking if any of the Islam-only indicators are > 0.
	// Seq 3, 4, 5, 6, 7, 8, 9 are index 2, 3, 4, 5, 6, 7, 8.
	isIslam := false
	for i := 2; i <= 8; i++ {
		if all[i] > 0 {
			isIslam = true
			break
		}
	}

	isRamadanMonth, _ := puasa.IsRamadanMonth(bulanTahun)

	if isIslam {
		if isRamadanMonth {
			maks = 90
			for _, s := range all {
				perolehan += s
			}
		} else {
			maks = 75
			for i, s := range all {
				if i != 6 && i != 7 && i != 8 { // Skip index 6 (Rowatib), 7 (Tarawih), 8 (Puasa)
					perolehan += s
				}
			}
		}
	} else {
		maks = 55
		for i, s := range all {
			if !islamOnly[i] {
				perolehan += s
			}
		}
	}

	akhir = math.Round((float64(perolehan)/float64(maks))*100*100) / 100
	predikat = Predikat(akhir)
	return
}

// Predikat memetakan nilai akhir ke predikat.
func Predikat(akhir float64) string {
	switch {
	case akhir >= 90:
		return "Istimewa"
	case akhir >= 80:
		return "Sangat Baik"
	case akhir >= 70:
		return "Baik"
	case akhir >= 60:
		return "Cukup"
	default:
		return "Kurang"
	}
}
