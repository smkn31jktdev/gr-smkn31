import type { Siswa, Admin } from '../../auth/types/auth.types';

export interface SkorG7 {
	bangunPagi: number;
	ibadahDoa: number;
	ibadahSholatFajar: number;
	ibadahSholat5Waktu: number;
	ibadahZikir: number;
	ibadahDhuha: number;
	ibadahRowatib: number;
	ibadahTarawih: number;
	ibadahPuasa: number;
	ibadahZakat: number;
	olahraga: number;
	makanSehat: number;
	belajarKitabSuci: number;
	belajarBukuUmum: number;
	belajarBukuMapel: number;
	belajarTugas: number;
	bermasyarakat: number;
	tidurCepat: number;
}

export interface G7Rekap {
	id: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	agama?: string;
	bulanTahun: string; // YYYY-MM
	hariTercatat: number;
	skor: SkorG7;
	nilaiMaks: number;
	nilaiPerolehan: number;
	nilaiAkhir: number;
	predikat: string;
	waliKelas: string;
	guruBK: string;
	status: 'draft' | 'reviewed' | 'final';
	tanggalFinal?: string;
	createdAt?: string;
	updatedAt?: string;
}

export interface G7RekapRingkas {
	nis: string;
	nama: string;
	nilaiAkhir: number;
}

export interface G7RekapStatistik {
	kelas?: string;
	bulanTahun: string;
	totalSiswa: number;
	sudahDinilai: number;
	belumDinilai: number;
	rataRataNilaiAkhir: number;
	distribusiPredikat: Record<string, number>;
	nilaiTertinggi?: G7RekapRingkas;
	nilaiTerendah?: G7RekapRingkas;
	rataRataPerIndikator: Record<string, number>;
}

export interface G7SuggestResponse {
	nis: string;
	bulanTahun: string;
	skor: SkorG7;
	catatan: Record<string, string>;
	hariTercatat: number;
	isAdvisory: boolean;
}

export interface RekapBulanan {
	id: string;
	rekapKey: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	bulanTahun: string;
	semester?: 'ganjil' | 'genap';
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalTidakHadir: number;
	totalMagang: number;
	totalHariEfektif: number;
	persentaseHadir: number;
}

export interface RekapHarian {
	tanggal: string;
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalAlpa: number;
	totalMagang: number;
	totalSiswa: number;
}

export interface RingkasanSiswa {
	nis: string;
	namaSiswa: string;
	kelas: string;
	bulanTahun: string;
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalTidakHadir: number;
	totalMagang: number;
	persentaseHadir: number;
}

export interface RekapSemesterItem {
	bulanTahun: string;
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalAlpa: number;
	totalMagang: number;
	persentaseHadir: number;
}

export interface RekapSemesterKelas {
	kelas: string;
	semester: string;
	bulan: RekapSemesterItem[];
}

// Rekap Lengkap Kehadiran

export interface RekapSiswaItem {
	nis: string;
	namaSiswa: string;
	kelas: string;
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalAlpa: number;
	totalMagang: number;
	hariEfektif: number;
	tingkatKehadiran: number;
	adaData: boolean;
}

// RekapKelasSummary ringkasan agregat satu kelas.
export interface RekapKelasSummary {
	kelas: string;
	totalSiswa: number;
	hariEfektif: number;
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalAlpa: number;
	totalMagang: number;
	siswaUnikHadir: number;
	tingkatKehadiran: number;
}

// RekapRange grand-total seluruh kelas dalam periode
export interface RekapRange {
	bulanTahun: string;
	hariEfektif: number;
	totalSiswa: number;
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalAlpa: number;
	totalMagang: number;
	siswaUnikHadir: number;
	tingkatKehadiran: number;
}

// RekapKelasLengkap
export interface RekapKelasLengkap {
	bulanTahun: string;
	hariEfektif: number;
	summaryByClass: RekapKelasSummary[];
	summaryByStudent: RekapSiswaItem[];
	summaryRange: RekapRange;
	totalSiswa?: number;
	page?: number;
	limit?: number;
	hasMore?: boolean;
}

// RekapMingguanKelas — rekap kehadiran satu minggu (Senin–Jumat) dari data harian nyata.
export interface RekapMingguanKelas {
	senin: string; // YYYY-MM-DD
	jumat: string; // YYYY-MM-DD
	hariEfektif: number;
	summaryByClass: RekapKelasSummary[];
	summaryByStudent: RekapSiswaItem[];
	summaryRange: RekapRange;
	totalSiswa?: number;
	page?: number;
	limit?: number;
	hasMore?: boolean;
}

export interface KehadiranHariItem {
	tanggal: string;
	status: string;
	waktuAbsen: string;
	alasan?: string;
	fotoIzin?: string;
}

export interface KehadiranBulanSummary {
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalAlpa: number;
	totalMagang: number;
	totalHariEfektif: number;
	persentaseHadir: number;
}

export interface KehadiranBulananSiswa {
	bulanTahun: string;
	kehadiran: KehadiranHariItem[];
	summary: KehadiranBulanSummary;
}

// Rekap Kelas G7

export interface G7RekapSiswaItem {
	nis: string;
	namaSiswa: string;
	kelas: string;
	nilaiPerolehan: number;
	nilaiMaks: number;
	nilaiAkhir: number;
	predikat: string;
	status: string;
	hariTercatat: number;
	sudahDinilai: boolean;
}

// G7RekapKelasLengkap
export interface G7RekapKelasLengkap {
	kelas: string;
	bulanTahun: string;
	totalSiswa: number;
	sudahDinilai: number;
	belumDinilai: number;
	statistik: G7RekapStatistik | null;
	siswa: G7RekapSiswaItem[];
}
