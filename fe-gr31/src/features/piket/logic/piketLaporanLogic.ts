import { writable, get } from 'svelte/store';
import {
	getRekapLengkap,
	getRekapSemester,
	getRekapBulananAdmin,
	getRekapMingguan,
	getKelas,
	listKehadiranAdmin,
	listStudents,
	getRekapSiswaDetail
} from '../../admin/logic/adminLogic';
import { addToast } from '../../../stores/uiStore';
import type {
	RekapSiswaItem,
	RekapKelasSummary,
	RekapRange,
	RekapSemesterItem,
	RekapBulanan
} from '../../admin/types/admin.types';
import type { Kehadiran } from '../../student/types/student.types';
import type { Siswa } from '../../auth/types/auth.types';

export const getBulanStr = () => {
	const d = new Date();
	return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
};

// Monday helper
export function getMonday(d: Date): Date {
	const day = d.getDay();
	const diff = d.getDate() - day + (day === 0 ? -6 : 1);
	const monday = new Date(d.setDate(diff));
	monday.setHours(0, 0, 0, 0);
	return monday;
}

export function formatDateStr(d: Date): string {
	const year = d.getFullYear();
	const month = String(d.getMonth() + 1).padStart(2, '0');
	const day = String(d.getDate()).padStart(2, '0');
	return `${year}-${month}-${day}`;
}

export function getCurrentSemester(): string {
	const now = new Date();
	const month = now.getMonth() + 1;
	const year = now.getFullYear();
	if (month >= 7) return `ganjil-${year}`;
	return `genap-${year}`;
}

// State stores
export const activeTab = writable<'kelas' | 'siswa' | 'lomba'>('kelas');
export const reportType = writable<'bulanan' | 'mingguan'>('bulanan');
export const loading = writable<boolean>(false);

// Filter stores
export const selectedKelas = writable<string>('');
export const selectedBulan = writable<string>(getBulanStr());
export const selectedWeekMonday = writable<string>(formatDateStr(getMonday(new Date())));

// Student Lookup stores
export const studentSearchQuery = writable<string>('');
export const foundStudents = writable<Siswa[]>([]);
export const selectedStudent = writable<Siswa | null>(null);

// Output data stores
export const items = writable<RekapSiswaItem[]>([]); // Laporan bulanan kelas (roster-join dari rekap-lengkap)
export const summaryByClass = writable<RekapKelasSummary[]>([]); // Ringkasan per kelas
export const summaryRange = writable<RekapRange | null>(null); // Grand total
export const weeklyClassItems = writable<RekapSiswaItem[]>([]); // Laporan mingguan (computed client-side)
export const rawWeeklyLogs = writable<Kehadiran[]>([]); // Raw logs for weekly chart
export const studentLogs = writable<Kehadiran[]>([]); // Daily breakdowns
export const studentSummary = writable({
	hadir: 0,
	izin: 0,
	sakit: 0,
	alpa: 0,
	magang: 0
});

// Monthly trend per-student (from backend rekap-bulanan, 6 bulan terakhir)
export const studentMonthlyTrend = writable<RekapBulanan[]>([]);

// ── Rekap Per Siswa (daftar terpaginasi) ─────────────────────────────────────
// Daftar rekap semua siswa untuk periode terpilih (bulanan/mingguan), diambil
// dari backend dengan pagination server-side agar fetch cepat.
export const studentRekapList = writable<RekapSiswaItem[]>([]);
export const studentRekapPage = writable<number>(1);
export const studentRekapLimit = writable<number>(50); // 50 | 100 | 0 (semua)
export const studentRekapTotal = writable<number>(0);
export const studentRekapHasMore = writable<boolean>(false);
export const studentRekapLoading = writable<boolean>(false);

// Semester trend data for monthly chart (from backend rekap-semester)
export const semesterTrend = writable<RekapSemesterItem[]>([]);

// Computed class-level stats derived from items (updated after loadData)
export const classTotalHariEfektif = writable<number>(0);

// Dynamic kelas list loaded from backend
export const availableKelas = writable<string[]>([]);

// Load kelas from backend
export async function loadKelas() {
	try {
		const kelas = await getKelas();
		if (kelas && kelas.length > 0) {
			availableKelas.set(kelas);
		}
	} catch (err) {
		console.error('Error loading kelas:', err);
	}
}

// Available months for select option — dynamically keep a rolling 6-month window
export const availableMonths = (() => {
	const months = [];
	const labels = [
		'Januari',
		'Februari',
		'Maret',
		'April',
		'Mei',
		'Juni',
		'Juli',
		'Agustus',
		'September',
		'Oktober',
		'November',
		'Desember'
	];
	const now = new Date();
	for (let i = 0; i < 6; i++) {
		const d = new Date(now.getFullYear(), now.getMonth() - i, 1);
		const val = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
		months.push({ value: val, label: `${labels[d.getMonth()]} ${d.getFullYear()}` });
	}
	return months;
})();

export function changeWeek(offsetDays: number) {
	const currentMon = new Date(get(selectedWeekMonday));
	currentMon.setDate(currentMon.getDate() + offsetDays);
	selectedWeekMonday.set(formatDateStr(currentMon));
	studentRekapPage.set(1);
	loadData();
}

export async function searchForStudents() {
	const q = get(studentSearchQuery);
	const kelas = get(selectedKelas) || '';
	if (!q.trim()) {
		foundStudents.set([]);
		return;
	}
	try {
		const res = await listStudents(1, 40, q, kelas);
		foundStudents.set(res.items || []);
	} catch (err) {
		console.error('Error searching students:', err);
	}
}

export async function loadData() {
	const tab = get(activeTab);
	const type = get(reportType);
	const kelas = get(selectedKelas);

	loading.set(true);

	try {
		if (tab === 'kelas') {
			if (type === 'bulanan') {
				// ── Laporan Kelas Bulanan ──────────────────────────────────────────────
				// Sumber: rekap-lengkap (roster-join dari rekap_bulanan, paling akurat)
				// Siswa yang belum punya data tetap tampil dengan adaData=false
				const bulan = get(selectedBulan);
				const [rekapRes, semesterRes] = await Promise.all([
					getRekapLengkap(bulan, kelas),
					getRekapSemester(kelas, getCurrentSemester())
				]);

				if (rekapRes) {
					items.set(rekapRes.summaryByStudent || []);
					summaryByClass.set(rekapRes.summaryByClass || []);
					summaryRange.set(rekapRes.summaryRange || null);
					classTotalHariEfektif.set(rekapRes.hariEfektif || 0);
				} else {
					items.set([]);
					summaryByClass.set([]);
					summaryRange.set(null);
					classTotalHariEfektif.set(0);
				}

				semesterTrend.set(semesterRes?.bulan || []);
				rawWeeklyLogs.set([]);
			} else {
				// ── Laporan Kelas Mingguan ────────────────────────────────────────────
				// Sumber: backend /rekap-mingguan (agregasi data harian nyata dari
				// koleksi kehadiran). hariEfektif minggu dihitung server-side dari
				// kalender pendidikan — bukan perhitungan client-side.
				const mon = new Date(get(selectedWeekMonday));
				const seninStr = formatDateStr(getMonday(mon));

				const rekapRes = await getRekapMingguan(seninStr, kelas);

				if (rekapRes) {
					weeklyClassItems.set(rekapRes.summaryByStudent || []);
					items.set(rekapRes.summaryByStudent || []);
					summaryByClass.set(rekapRes.summaryByClass || []);
					summaryRange.set(rekapRes.summaryRange || null);
					classTotalHariEfektif.set(rekapRes.hariEfektif || 0);
				} else {
					weeklyClassItems.set([]);
					items.set([]);
					summaryByClass.set([]);
					summaryRange.set(null);
					classTotalHariEfektif.set(0);
				}

				// Fetch raw logs untuk chart harian mingguan (tetap dari kehadiran).
				const fri = new Date(getMonday(mon));
				fri.setDate(fri.getDate() + 4);
				const logsRes = await listKehadiranAdmin(
					{ kelas, dari: seninStr, sampai: formatDateStr(fri) },
					1,
					0
				);
				rawWeeklyLogs.set(logsRes.items || []);
			}
		} else {
			// ── Laporan Per Siswa ─────────────────────────────────────────────────
			const student = get(selectedStudent);
			if (!student) {
				await loadStudentRekapList();
				loading.set(false);
				return;
			}

			// Fetch log kehadiran + trend bulanan 6 bulan terakhir secara paralel
			const now = new Date();
			const sixMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 5, 1);
			const trendFrom = `${sixMonthsAgo.getFullYear()}-${String(sixMonthsAgo.getMonth() + 1).padStart(2, '0')}`;
			const trendTo = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;

			const [res, trendRes] = await Promise.all([
				getRekapSiswaDetail({
					nis: student.nis,
					tipe: type,
					bulan: type === 'bulanan' ? get(selectedBulan) : undefined,
					senin: type === 'mingguan' ? get(selectedWeekMonday) : undefined
				}),
				getRekapBulananAdmin({ nis: student.nis }, 1, 12)
			]);

			if (res) {
				const mappedLogs: Kehadiran[] = (res.kehadiran || []).map(
					(item) =>
						({
							tanggal: item.tanggal,
							status: item.status as any,
							waktuAbsen: item.waktuAbsen,
							alasan: item.alasan,
							fotoIzin: item.fotoIzin
						}) as Kehadiran
				);
				studentLogs.set(mappedLogs);

				studentSummary.set({
					hadir: res.summary.totalHadir,
					izin: res.summary.totalIzin,
					sakit: res.summary.totalSakit,
					alpa: res.summary.totalAlpa,
					magang: res.summary.totalMagang
				});
			} else {
				studentLogs.set([]);
				studentSummary.set({ hadir: 0, izin: 0, sakit: 0, alpa: 0, magang: 0 });
			}

			// Trend bulanan: ambil 6 bulan terakhir, sort asc
			const trendItems = (trendRes?.items || []) as RekapBulanan[];
			const filtered = trendItems
				.filter((r) => r.bulanTahun >= trendFrom && r.bulanTahun <= trendTo)
				.sort((a, b) => a.bulanTahun.localeCompare(b.bulanTahun));
			studentMonthlyTrend.set(filtered);
		}
	} catch (err) {
		console.error('Error loading report data:', err);
		addToast('Gagal memuat rekap laporan', 'error');
	} finally {
		loading.set(false);
	}
}

// loadStudentRekapList memuat daftar rekap per-siswa (terpaginasi) untuk periode
// terpilih. Sumber: backend rekap-lengkap (bulanan) / rekap-mingguan (mingguan)
// dengan pagination server-side. Dipakai pada tab "Laporan Per Siswa".
export async function loadStudentRekapList() {
	const type = get(reportType);
	const kelas = get(selectedKelas);
	const page = get(studentRekapPage);
	const limit = get(studentRekapLimit);

	studentRekapLoading.set(true);
	try {
		if (type === 'bulanan') {
			const bulan = get(selectedBulan);
			const res = await getRekapLengkap(bulan, kelas, page, limit);
			if (res) {
				studentRekapList.set(res.summaryByStudent || []);
				studentRekapTotal.set(res.totalSiswa || 0);
				studentRekapHasMore.set(res.hasMore || false);
			} else {
				studentRekapList.set([]);
				studentRekapTotal.set(0);
				studentRekapHasMore.set(false);
			}
		} else {
			const mon = new Date(get(selectedWeekMonday));
			const seninStr = formatDateStr(getMonday(mon));
			const res = await getRekapMingguan(seninStr, kelas, page, limit);
			if (res) {
				studentRekapList.set(res.summaryByStudent || []);
				studentRekapTotal.set(res.totalSiswa || 0);
				studentRekapHasMore.set(res.hasMore || false);
			} else {
				studentRekapList.set([]);
				studentRekapTotal.set(0);
				studentRekapHasMore.set(false);
			}
		}
	} catch (err) {
		console.error('Error loading student rekap list:', err);
		addToast('Gagal memuat daftar rekap siswa', 'error');
		studentRekapList.set([]);
		studentRekapTotal.set(0);
		studentRekapHasMore.set(false);
	} finally {
		studentRekapLoading.set(false);
	}
}

// changeStudentRekapPage menggeser halaman daftar rekap siswa lalu memuat ulang.
export function changeStudentRekapPage(delta: number) {
	const next = get(studentRekapPage) + delta;
	if (next < 1) return;
	studentRekapPage.set(next);
	loadStudentRekapList();
}

// setStudentRekapLimit mengganti ukuran halaman (50 | 100 | 0=semua) dan reset ke hal. 1.
export function setStudentRekapLimit(limit: number) {
	studentRekapLimit.set(limit);
	studentRekapPage.set(1);
	loadStudentRekapList();
}

// resetStudentRekapPage mengembalikan ke halaman 1 (dipakai saat ganti periode/kelas).
export function resetStudentRekapPage() {
	studentRekapPage.set(1);
}
