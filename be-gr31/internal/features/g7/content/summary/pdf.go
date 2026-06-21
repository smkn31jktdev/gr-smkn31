package summary

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"time"

	"be-gr31/internal/features/g7/content/evaluate"
	g7model "be-gr31/internal/model/g7"
)

const indikatorNames = `Bangun Pagi : Siswa biasa bangun pagi sebelum jam 05.00
Beribadah : Siswa biasa beribadah
Berolah Raga : Siswa melakukan kebiasaan berolahraga
Makan Sehat Bergizi : Siswa memiliki kebiasaan makan sehat bergizi
Gemar Belajar : Siswa memiliki kebiasaan gemar belajar
Baik di Masyarakat : Siswa memiliki perilaku baik di masyarakat
Tidur Lebih Cepat : Siswa biasa tidur/istirahat malam sebelum jam 22.00`

// BuildPDFLaporan merakit PDFLaporan dari rekap G7 dan hasil evaluasi jurnal
func BuildPDFLaporan(rekap *g7model.G7Rekap, report evaluate.EvalReport) *g7model.PDFLaporan {
	names := strings.Split(indikatorNames, "\n")

	evals := []evaluate.EvalResult{
		report.BangunPagi,
		report.Beribadah,
		report.Olahraga,
		report.MakanSehat,
		report.GemarBelajar,
		report.Bermasyarakat,
		report.TidurCepat,
	}

	rows := make([]g7model.PDFRowIndikator, 0, 7)
	for i, name := range names {
		ev := evals[i]
		ket := ev.Note
		if ev.Skor == 0 {
			ket = "Belum diisi"
		}
		rows = append(rows, g7model.PDFRowIndikator{
			No:         i + 1,
			Indikator:  name,
			Skor:       ev.Skor,
			Keterangan: ket,
		})
	}

	bulan := rekap.BulanTahun
	tahunStr := ""
	if len(bulan) >= 4 {
		tahunStr = bulan[:4]
	}

	return &g7model.PDFLaporan{
		NamaSiswa:    rekap.NamaSiswa,
		NIS:          rekap.NISN,
		Kelas:        rekap.Kelas,
		BulanTahun:   bulan,
		TahunAjaran:  tahunStr,
		Indikator:    rows,
		WaliKelas:    rekap.WaliKelas,
		OrangTua:     "",
		NilaiAkhir:   rekap.NilaiAkhir,
		Predikat:     rekap.Predikat,
		TanggalCetak: time.Now().Format("2 January 2006"),
	}
}

// RenderHTMLLaporan menghasilkan HTML string siap print/PDF sesuai layout gambar
func RenderHTMLLaporan(data *g7model.PDFLaporan) (string, error) {
	tmpl, err := template.New("laporan").Funcs(template.FuncMap{
		"seq": func(start, end int) []int {
			r := make([]int, 0, end-start+1)
			for i := start; i <= end; i++ {
				r = append(r, i)
			}
			return r
		},
		"centang": func(skor, target int) string {
			if skor == target {
				return "V"
			}
			return ""
		},
		"bulanIndo": bulanIndonesia,
	}).Parse(htmlTemplate)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func bulanIndonesia(bulanTahun string) string {
	if len(bulanTahun) < 7 {
		return bulanTahun
	}
	bulanMap := map[string]string{
		"01": "Januari", "02": "Februari", "03": "Maret", "04": "April",
		"05": "Mei", "06": "Juni", "07": "Juli", "08": "Agustus",
		"09": "September", "10": "Oktober", "11": "November", "12": "Desember",
	}
	mm := bulanTahun[5:7]
	yyyy := bulanTahun[:4]
	return bulanMap[mm] + " " + yyyy
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8"/>
<title>Laporan 7 Kebiasaan - {{.NamaSiswa}}</title>
<style>
  @page { size: A4; margin: 18mm 14mm; }
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: Arial, sans-serif; font-size: 10pt; color: #000; background: #fff; }
  .page { width: 100%; }

  /* Header */
  .header { text-align: center; margin-bottom: 14px; }
  .header h1 { font-size: 12pt; font-weight: bold; text-decoration: underline; color: #003399; }
  .header h2 { font-size: 11pt; font-weight: bold; color: #003399; }
  .header h3 { font-size: 11pt; font-weight: bold; color: #003399; }

  /* Info siswa */
  .info-siswa { margin-bottom: 14px; }
  .info-siswa table td { padding: 2px 4px; font-size: 10pt; font-weight: bold; }

  /* Tabel nilai */
  .tbl-nilai { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
  .tbl-nilai th, .tbl-nilai td { border: 1px solid #444; padding: 5px 6px; vertical-align: middle; text-align: center; }
  .tbl-nilai .col-no { width: 30px; }
  .tbl-nilai .col-ind { width: 34%; color: #003399; font-size: 9.5pt; }
  .tbl-nilai td.col-ind { text-align: left; }
  .tbl-nilai th.col-ind { text-align: center !important; }
  .tbl-nilai .col-skor { width: 56px; }
  .tbl-nilai .col-ket { font-size: 9pt; }
  .tbl-nilai td.col-ket { text-align: left; }
  .tbl-nilai th.col-ket { text-align: center !important; }
  .tbl-nilai thead th { background: #f5f5dc; font-weight: bold; color: #003399; font-size: 9pt; }
  .tbl-nilai .th-skor-grp { font-size: 9pt; }
  .tbl-nilai .centang { font-weight: bold; }
  .row-alt { background: #fff; }

  /* Keterangan belum diisi */
  .belum { color: #888; font-style: italic; }

  /* Tanda tangan */
  .ttd-section { margin-top: 10px; }
  .ttd-kota { text-align: right; margin-bottom: 20px; font-size: 10pt; color: #003399; }
  .ttd-table { width: 100%; }
  .ttd-table td { width: 33%; vertical-align: top; text-align: center; font-size: 10pt; }
  .ttd-nama { margin-top: 80px; border-top: 1px solid #333; padding-top: 4px; font-weight: bold; color: #000; font-size: 10pt; }
  .ttd-role { font-size: 10pt; color: #333; }

  @media print {
    body { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
  }
</style>
</head>
<body>
<div class="page">
  <!-- Header -->
  <div class="header">
    <h1>LAPORAN PROSES 7 KEBIASAAN BAIK ANAK INDONESIA HEBAT</h1>
    <h2>SMK NEGERI 31 JAKARTA</h2>
    <h3>TAHUN {{.TahunAjaran}}</h3>
  </div>

  <!-- Info Siswa -->
  <div class="info-siswa">
    <table>
      <tr><td>Nama Siswa</td><td>: <strong>{{.NamaSiswa}}</strong></td></tr>
      <tr><td>NIS / NISN</td><td>: <strong>{{.NIS}}</strong></td></tr>
      <tr><td>Kelas</td><td>: <strong>{{.Kelas}}</strong></td></tr>
    </table>
  </div>

  <!-- Tabel Penilaian -->
  <table class="tbl-nilai">
    <thead>
      <tr>
        <th class="col-no" rowspan="2">NO</th>
        <th class="col-ind" rowspan="2" style="text-align: center;">INDIKATOR</th>
        <th colspan="5" class="th-skor-grp">
          <table width="100%" style="border:none;border-collapse:collapse;">
            <tr>
              <td style="border:none;width:20%;text-align:center;font-size:8.5pt;">Kurang<br/>Baik</td>
              <td style="border:none;width:20%;text-align:center;font-size:8.5pt;">Cukup<br/>Baik</td>
              <td style="border:none;width:20%;text-align:center;font-size:8.5pt;">Baik</td>
              <td style="border:none;width:20%;text-align:center;font-size:8.5pt;">Sangat<br/>Baik</td>
              <td style="border:none;width:20%;text-align:center;font-size:8.5pt;">Istimewa</td>
            </tr>
          </table>
        </th>
        <th class="col-ket" rowspan="2" style="text-align: center;">KETERANGAN</th>
      </tr>
      <tr>
        <th class="col-skor">1</th>
        <th class="col-skor">2</th>
        <th class="col-skor">3</th>
        <th class="col-skor">4</th>
        <th class="col-skor">5</th>
      </tr>
    </thead>
    <tbody>
      {{range .Indikator}}
      <tr class="row-alt">
        <td>{{.No}}</td>
        <td class="col-ind">{{.Indikator}}</td>
        <td class="centang">{{centang .Skor 1}}</td>
        <td class="centang">{{centang .Skor 2}}</td>
        <td class="centang">{{centang .Skor 3}}</td>
        <td class="centang">{{centang .Skor 4}}</td>
        <td class="centang">{{centang .Skor 5}}</td>
        <td class="col-ket {{if eq .Skor 0}}belum{{end}}">{{.Keterangan}}</td>
      </tr>
      {{end}}
    </tbody>
  </table>

  <!-- Tanda Tangan -->
  <div class="ttd-section">
    <div class="ttd-kota">Jakarta, {{bulanIndo .BulanTahun}}</div>
    <table class="ttd-table">
      <tr>
        <td>
          <div class="ttd-role">Guru Wali</div>
          <div class="ttd-nama">{{.WaliKelas}}</div>
        </td>
        <td>
          <div class="ttd-role">Orang Tua / Wali</div>
          <div class="ttd-nama">Orang Tua / Wali</div>
        </td>
        <td>
          <div class="ttd-role">Siswa</div>
          <div class="ttd-nama">{{.NamaSiswa}}</div>
        </td>
      </tr>
    </table>
  </div>
</div>
</body>
</html>`
