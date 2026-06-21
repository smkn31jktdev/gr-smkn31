import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { G7Jurnal } from '../types/student.types';

// ── Tipe response dashboard dari backend ────────────────────────────────────

export interface G7RingkasanBulan {
	bulanTahun: string;
	hariTercatat: number;
	rataRataDone: number;
	totalDoneSum: number;
}

export interface G7DashboardSiswa {
	jurnalHariIni: G7Jurnal | null;
	progresHariIni: number; // 0–7, dihitung server-side
	ringkasanBulan: G7RingkasanBulan;
}

// ── API calls ────────────────────────────────────────────────────────────────

/**
 * Mengambil data dashboard G7 siswa yang sedang login:
 * jurnal hari ini + ringkasan bulan berjalan dalam satu request.
 * Backend: GET /v1/student/g7/dashboard
 * NISN diambil dari JWT siswa — tidak perlu dikirim dari frontend.
 */
export async function getDashboardG7(): Promise<G7DashboardSiswa | null> {
	const { data, error } = await apiRequest<G7DashboardSiswa>('/v1/student/g7/dashboard');
	if (error) return null;
	return data;
}

/**
 * Mengambil jurnal G7 siswa untuk tanggal tertentu.
 * Mengembalikan null jika belum ada jurnal (404) atau terjadi error.
 * Backend: GET /v1/student/g7/:tanggal
 */
export async function getJurnalForDate(tanggal: string): Promise<G7Jurnal | null> {
	const { data, error, status } = await apiRequest<G7Jurnal>(`/v1/student/g7/${tanggal}`);
	if (status === 404 || !data) return null;
	if (error) return null;
	return data;
}

/**
 * Submit atau update jurnal G7 siswa untuk tanggal tertentu.
 * NISN diambil dari JWT di backend — tidak perlu dikirim dari frontend.
 * Backend: POST /v1/student/g7
 */
export async function submitJurnal(
	tanggal: string,
	journalData: Partial<G7Jurnal>
): Promise<boolean> {
	const payload = {
		tanggal,
		bangun: journalData.bangun ?? { done: false, waktu: '', keterangan: '' },
		ibadah: journalData.ibadah ?? { done: false, waktu: '', keterangan: '' },
		makan: journalData.makan ?? { done: false, waktu: '', keterangan: '' },
		olahraga: journalData.olahraga ?? { done: false, waktu: '', keterangan: '' },
		belajar: journalData.belajar ?? { done: false, waktu: '', keterangan: '' },
		bermasyarakat: journalData.bermasyarakat ?? { done: false, waktu: '', keterangan: '' },
		tidur: journalData.tidur ?? { done: false, waktu: '', keterangan: '' }
	};

	const { error } = await apiRequest('/v1/student/g7', {
		method: 'POST',
		body: JSON.stringify(payload)
	});

	if (error) {
		addToast(error, 'error');
		return false;
	}

	addToast('Jurnal Harian G7 berhasil disimpan!', 'success');
	return true;
}

/**
 * Mengambil daftar jurnal G7 siswa dari collection kebiasaan_hebat.
 * Filter by bulan agar tidak load semua data.
 * Backend: GET /v1/student/g7?dari=YYYY-MM&sampai=YYYY-MM
 */
export async function listJurnalSiswa(dari?: string, sampai?: string, page = 1, limit = 31) {
	const now = new Date();
	const bulanDefault = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;

	const query = new URLSearchParams();
	// Hanya append jika tidak kosong — string kosong menyebabkan filter tidak benar di backend
	if (dari && dari.trim() !== '') query.append('dari', dari);
	else query.append('dari', bulanDefault);
	if (sampai && sampai.trim() !== '') query.append('sampai', sampai);
	else query.append('sampai', bulanDefault);
	query.append('page', String(page));
	query.append('limit', String(limit));

	const { data, error } = await apiRequest<{ items: G7Jurnal[]; total: number; hasMore: boolean }>(
		`/v1/student/g7?${query.toString()}`
	);
	if (error) {
		addToast('Gagal memuat riwayat jurnal G7', 'error');
		return { items: [] as G7Jurnal[], total: 0, hasMore: false };
	}

	return {
		items: data?.items ?? [],
		total: data?.total ?? 0,
		hasMore: data?.hasMore ?? false
	};
}
