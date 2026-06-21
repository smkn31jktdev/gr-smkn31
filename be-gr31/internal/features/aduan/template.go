package aduan

const htmlTemplate = `<!DOCTYPE html>
<html lang="id">
<head>
	<meta charset="UTF-8"/>
	<title>Arsip Aduan Siswa - {{.ID}}</title>
	<style>
		@page { size: A4; margin: 15mm 15mm; }
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body { font-family: 'Segoe UI', Arial, sans-serif; font-size: 10pt; color: #1e293b; background: #fff; padding: 20px; }
		.container { width: 100%; max-width: 800px; margin: 0 auto; }
		
		/* Kop Surat / Header */
		.kop-surat { text-align: center; border-bottom: 3px double #000; padding-bottom: 12px; margin-bottom: 20px; }
		.kop-surat h1 { font-size: 14pt; font-weight: bold; text-transform: uppercase; margin-bottom: 2px; }
		.kop-surat h2 { font-size: 12pt; font-weight: bold; text-transform: uppercase; margin-bottom: 4px; color: #475569; }
		.kop-surat p { font-size: 8pt; color: #64748b; font-style: italic; }

		.title-doc { text-align: center; font-size: 12pt; font-weight: bold; text-transform: uppercase; margin: 15px 0; text-decoration: underline; }

		/* Info Card Block */
		.info-card {
			background-color: #f8fafc;
			border: 1px solid #cbd5e1;
			border-radius: 8px;
			padding: 16px 20px;
			margin-bottom: 24px;
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 12px 24px;
		}
		.info-group {
			display: flex;
			flex-direction: column;
			gap: 2px;
		}
		.info-label {
			font-size: 8pt;
			font-weight: bold;
			color: #475569;
			text-transform: uppercase;
			letter-spacing: 0.5px;
		}
		.info-value {
			font-size: 9.5pt;
			font-weight: 600;
			color: #0f172a;
		}
		.status-badge {
			display: inline-block;
			padding: 2px 6px;
			border-radius: 4px;
			font-size: 8pt;
			font-weight: bold;
			text-transform: uppercase;
		}
		.status-badge.pending { background-color: #fef3c7; color: #d97706; border: 1px solid #fde68a; }
		.status-badge.open { background-color: #ecfdf5; color: #059669; border: 1px solid #a7f3d0; }
		.status-badge.in_progress { background-color: #eff6ff; color: #2563eb; border: 1px solid #bfdbfe; }
		.status-badge.closed { background-color: #f1f5f9; color: #475569; border: 1px solid #cbd5e1; }

		.section-title { font-size: 10pt; font-weight: bold; text-transform: uppercase; border-bottom: 1px solid #cbd5e1; padding-bottom: 4px; margin-top: 20px; margin-bottom: 12px; color: #1e293b; }

		/* Chat Timeline */
		.chat-container { display: flex; flex-direction: column; gap: 12px; margin-bottom: 25px; }
		.chat-bubble { padding: 10px 14px; border-radius: 8px; font-size: 9.5pt; max-width: 85%; line-height: 1.5; word-wrap: break-word; }
		.chat-bubble.student { background-color: #f1f5f9; border: 1px solid #e2e8f0; align-self: flex-start; }
		.chat-bubble.admin { background-color: #f0fdfa; border: 1px solid #ccfbf1; align-self: flex-end; }
		.bubble-header { display: flex; justify-content: space-between; gap: 20px; font-size: 8pt; font-weight: bold; color: #64748b; margin-bottom: 4px; }
		.bubble-content { color: #0f172a; white-space: pre-wrap; font-weight: 500; }
		.bubble-time { font-size: 7.5pt; color: #94a3b8; text-align: right; margin-top: 4px; }

		/* Status History Table */
		.history-table { width: 100%; border-collapse: collapse; margin-bottom: 30px; font-size: 9pt; }
		.history-table th, .history-table td { border: 1px solid #cbd5e1; padding: 6px 8px; text-align: left; }
		.history-table th { background-color: #f8fafc; font-weight: bold; color: #334155; }
		
		/* Tanda Tangan */
		.ttd-section { margin-top: 40px; page-break-inside: avoid; }
		.ttd-date { text-align: right; font-size: 9.5pt; margin-bottom: 15px; }
		.ttd-table { width: 100%; }
		.ttd-table td { width: 50%; text-align: center; vertical-align: top; font-size: 9.5pt; }
		.ttd-space { height: 70px; }
		.ttd-name { font-weight: bold; text-decoration: underline; }
		.ttd-nip { font-size: 8.5pt; color: #475569; margin-top: 2px; }

		/* Floating Print Action (screen-only) */
		.print-action { position: fixed; bottom: 20px; right: 20px; z-index: 9999; }
		.print-btn { background-color: #0f172a; color: #ffffff; border: none; padding: 10px 16px; border-radius: 8px; font-size: 9pt; font-weight: bold; cursor: pointer; display: flex; align-items: center; gap: 6px; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); }
		.print-btn:hover { background-color: #1e293b; }

		@media print {
			body { padding: 0; -webkit-print-color-adjust: exact; print-color-adjust: exact; }
			.print-action { display: none; }
			.info-card { background-color: #f8fafc !important; border: 1px solid #cbd5e1 !important; }
			.status-badge.pending { background-color: #fef3c7 !important; color: #d97706 !important; }
			.status-badge.open { background-color: #ecfdf5 !important; color: #059669 !important; }
			.status-badge.in_progress { background-color: #eff6ff !important; color: #2563eb !important; }
			.status-badge.closed { background-color: #f1f5f9 !important; color: #475569 !important; }
			.chat-bubble.student { background-color: #f1f5f9 !important; border: 1px solid #e2e8f0 !important; }
			.chat-bubble.admin { background-color: #f0fdfa !important; border: 1px solid #ccfbf1 !important; }
		}
	</style>
</head>
<body onload="window.print()">
	<div class="print-action">
		<button class="print-btn" onclick="window.print()">
			Cetak Laporan / Simpan PDF
		</button>
	</div>

	<div class="container">
		<!-- Kop Surat -->
		<div class="kop-surat">
			<h1>Pemerintah Daerah Provinsi DKI Jakarta</h1>
			<h2>SMK Negeri 31 Jakarta</h2>
			<p>Jl. Kramat Raya No.42, Senen, Jakarta Pusat | Telp: (021) 390xxxx | Email: info@smkn31jakarta.sch.id</p>
		</div>

		<div class="title-doc">Arsip Aduan & Konseling Siswa</div>

		<!-- Informasi Siswa & Tiket -->
		<div class="info-card">
			<div class="info-group">
				<span class="info-label">ID Tiket</span>
				<span class="info-value">{{.ID}}</span>
			</div>
			<div class="info-group">
				<span class="info-label">Tanggal Pengaduan</span>
				<span class="info-value">{{formatDate .CreatedAt}}</span>
			</div>
			<div class="info-group">
				<span class="info-label">Nama Siswa</span>
				<span class="info-value" style="text-transform: uppercase;">{{.NamaSiswa}}</span>
			</div>
			<div class="info-group">
				<span class="info-label">Status Tiket</span>
				<span class="info-value">
					<span class="status-badge {{.Status}}">
						{{if eq .Status "open"}}BARU
						{{else if eq .Status "in_progress"}}DALAM PROSES
						{{else if eq .Status "closed"}}SELESAI
						{{else if eq .Status "pending"}}TERTUNDA
						{{else}}{{.Status}}
						{{end}}
					</span>
				</span>
			</div>
			<div class="info-group">
				<span class="info-label">NISN</span>
				<span class="info-value">{{.NISN}}</span>
			</div>
			<div class="info-group">
				<span class="info-label">Ditangani Oleh</span>
				<span class="info-value">{{if .AdminNama}}{{.AdminNama}}{{else}}-{{end}}</span>
			</div>
			<div class="info-group">
				<span class="info-label">Kelas / Guru Wali</span>
				<span class="info-value">{{.Kelas}} / {{if .Wali}}{{.Wali}}{{else}}-{{end}}</span>
			</div>
			<div class="info-group">
				<span class="info-label">Terakhir Diperbarui</span>
				<span class="info-value">{{formatDate .UpdatedAt}}</span>
			</div>
		</div>

		<!-- Riwayat Percakapan -->
		<div class="section-title">Riwayat Percakapan / Konseling</div>
		<div class="chat-container">
			{{range .Messages}}
			<div class="chat-bubble {{if eq .Role "admin"}}admin{{else}}student{{end}}">
				<div class="bubble-header">
					<span>{{if eq .Role "admin"}}Konselor / Guru BK{{else}}{{.From}}{{end}}</span>
					<span>{{if eq .Role "admin"}}[Staff/BK]{{else}}[Siswa]{{end}}</span>
				</div>
				<div class="bubble-content">{{.Isi}}</div>
				<div class="bubble-time">{{formatDate .Timestamp}}</div>
			</div>
			{{else}}
			<div style="text-align: center; color: #64748b; font-style: italic; font-size: 9.5pt; padding: 20px;">
				Tidak ada riwayat percakapan.
			</div>
			{{end}}
		</div>

		<!-- Riwayat Status -->
		{{if .StatusHistory}}
		<div class="section-title">Riwayat Log Status</div>
		<table class="history-table">
			<thead>
				<tr>
					<th style="width: 25%;">Tanggal & Waktu</th>
					<th style="width: 20%;">Status</th>
					<th style="width: 25%;">Diperbarui Oleh</th>
					<th style="width: 30%;">Catatan</th>
				</tr>
			</thead>
			<tbody>
				{{range .StatusHistory}}
				<tr>
					<td>{{formatDate .UpdatedAt}}</td>
					<td>
						{{if eq .Status "open"}}BARU
						{{else if eq .Status "in_progress"}}DALAM PROSES
						{{else if eq .Status "closed"}}SELESAI
						{{else}}{{.Status}}
						{{end}}
					</td>
					<td>{{.UpdatedBy}} ({{.Role}})</td>
					<td>{{if .Note}}{{.Note}}{{else}}-{{end}}</td>
				</tr>
				{{end}}
			</tbody>
		</table>
		{{end}}

		<!-- Tanda Tangan -->
		<div class="ttd-section">
			<table class="ttd-table">
				<tr>
					<td>
						<p>Siswa Yang Bersangkutan,</p>
						<div class="ttd-space" style="height: 92px;"></div>
						<p class="ttd-name">{{.NamaSiswa}}</p>
						<p class="ttd-nip">NISN. {{.NISN}}</p>
					</td>
					<td>
						<p style="margin-bottom: 2px;">Jakarta, {{currentDate}}</p>
						<p>Konselor / Guru BK,</p>
						<div class="ttd-space" style="height: 70px;"></div>
						<p class="ttd-name">{{if .AdminNama}}{{.AdminNama}}{{else}}Guru Bimbingan Konseling{{end}}</p>
						<p class="ttd-nip">Petugas Pelaksana</p>
					</td>
				</tr>
			</table>
		</div>
	</div>
</body>
</html>`
