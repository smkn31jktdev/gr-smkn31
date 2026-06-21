package evaluate

import g7model "be-gr31/internal/model/g7"

// EvalResult
type EvalResult struct {
	Skor int    `json:"skor"`
	Note string `json:"note"`
}

// EvalReport jurnal bulanan
type EvalReport struct {
	BangunPagi    EvalResult `json:"bangunPagi"`
	Beribadah     EvalResult `json:"beribadah"`
	Olahraga      EvalResult `json:"olahraga"`
	MakanSehat    EvalResult `json:"makanSehat"`
	GemarBelajar  EvalResult `json:"gemarBelajar"`
	Bermasyarakat EvalResult `json:"bermasyarakat"`
	TidurCepat    EvalResult `json:"tidurCepat"`
}

// Jurnals mengevaluasi slice jurnal bulanan menjadi EvalReport
func Jurnals(jurnals []g7model.G7) EvalReport {
	hari := len(jurnals)
	if hari == 0 {
		return EvalReport{
			BangunPagi:    EvalResult{0, "Belum ada data bangun pagi."},
			Beribadah:     EvalResult{0, "Belum ada data ibadah."},
			Olahraga:      EvalResult{0, "Belum ada data olahraga."},
			MakanSehat:    EvalResult{0, "Belum ada data makan sehat."},
			GemarBelajar:  EvalResult{0, "Belum ada data belajar."},
			Bermasyarakat: EvalResult{0, "Belum ada data kegiatan masyarakat."},
			TidurCepat:    EvalResult{0, "Belum ada data tidur."},
		}
	}
	return EvalReport{
		BangunPagi:    Bangun(jurnals),
		Beribadah:     Beribadah(jurnals),
		Olahraga:      Olahraga(jurnals, hari),
		MakanSehat:    Makan(jurnals, hari),
		GemarBelajar:  Belajar(jurnals, hari),
		Bermasyarakat: Masyarakat(jurnals),
		TidurCepat:    Tidur(jurnals),
	}
}

// ToSkorG7
func ToSkorG7(r EvalReport) g7model.SkorG7 {
	return g7model.SkorG7{
		BangunPagi:         r.BangunPagi.Skor,
		IbadahDoa:          r.Beribadah.Skor,
		IbadahSholatFajar:  r.Beribadah.Skor,
		IbadahSholat5Waktu: r.Beribadah.Skor,
		IbadahZikir:        r.Beribadah.Skor,
		IbadahDhuha:        r.Beribadah.Skor,
		IbadahRowatib:      r.Beribadah.Skor,
		IbadahTarawih:      0,
		IbadahPuasa:        0,
		IbadahZakat:        r.Beribadah.Skor,
		Olahraga:           r.Olahraga.Skor,
		MakanSehat:         r.MakanSehat.Skor,
		BelajarKitabSuci:   r.GemarBelajar.Skor,
		BelajarBukuUmum:    r.GemarBelajar.Skor,
		BelajarBukuMapel:   r.GemarBelajar.Skor,
		BelajarTugas:       r.GemarBelajar.Skor,
		Bermasyarakat:      r.Bermasyarakat.Skor,
		TidurCepat:         r.TidurCepat.Skor,
	}
}
