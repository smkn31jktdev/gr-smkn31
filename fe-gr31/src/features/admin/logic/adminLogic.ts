import { apiRequest, apiUpload } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type {
	G7Rekap,
	G7RekapStatistik,
	G7RekapKelasLengkap,
	G7SuggestResponse,
	RekapBulanan,
	RekapHarian,
	RekapKelasLengkap,
	RingkasanSiswa,
	RekapSemesterKelas,
	RekapMingguanKelas,
	KehadiranBulananSiswa
} from '../types/admin.types';
import type { Siswa, Admin } from '../../auth/types/auth.types';

// Kehadiran (Attendance) admin handlers
export async function getRekapHarian(tanggal: string, kelas?: string): Promise<RekapHarian | null> {
	const query = new URLSearchParams({ tanggal });
	if (kelas) query.append('kelas', kelas);

	const { data, error } = await apiRequest<RekapHarian>(
		`/v1/admin/rekap-harian?${query.toString()}`
	);
	if (error) {
		addToast('Gagal memuat rekap harian kehadiran', 'error');
		return null;
	}
	return data;
}

export async function getRingkasanKelas(kelas: string, bulan: string, page = 1, limit = 20) {
	const query = new URLSearchParams({ kelas, bulan, page: String(page), limit: String(limit) });
	const { data, error } = await apiRequest<any>(`/v1/admin/rekap-kelas?${query.toString()}`);
	if (error) {
		addToast('Gagal memuat rekap kelas', 'error');
		return { items: [], total: 0, hasMore: false };
	}
	return {
		items: (data.items || []) as RingkasanSiswa[],
		total: data.total || 0,
		hasMore: data.hasMore || false
	};
}

export async function getRekapSemester(
	kelas: string,
	semester: string
): Promise<RekapSemesterKelas | null> {
	const query = new URLSearchParams({ kelas, semester });
	const { data, error } = await apiRequest<RekapSemesterKelas>(
		`/v1/admin/rekap-semester?${query.toString()}`
	);
	if (error) {
		addToast('Gagal memuat rekap semester', 'error');
		return null;
	}
	return data;
}

export async function getRekapBulananAdmin(
	filters: { nis?: string; nisn?: string; kelas?: string; bulan?: string; semester?: string },
	page = 1,
	limit = 100
) {
	const query = new URLSearchParams({ page: String(page), limit: String(limit) });
	const nisVal = filters.nis || filters.nisn;
	if (nisVal) query.append('nisn', nisVal);
	if (filters.kelas) query.append('kelas', filters.kelas);
	if (filters.bulan) query.append('bulan', filters.bulan);
	if (filters.semester) query.append('semester', filters.semester);

	const { data, error } = await apiRequest<any>(`/v1/admin/rekap-bulanan?${query.toString()}`);
	if (error) return { items: [], total: 0, hasMore: false };
	return {
		items: (data.items || []) as RekapBulanan[],
		total: data.total || 0,
		hasMore: data.hasMore || false
	};
}

// getRekapLengkap
export async function getRekapLengkap(
	bulan: string,
	kelas?: string,
	page = 1,
	limit = 0
): Promise<RekapKelasLengkap | null> {
	const query = new URLSearchParams({ bulan, page: String(page), limit: String(limit) });
	if (kelas) query.append('kelas', kelas);

	const { data, error } = await apiRequest<RekapKelasLengkap>(
		`/v1/admin/rekap-lengkap?${query.toString()}`
	);
	if (error) {
		addToast('Gagal memuat rekap lengkap kehadiran', 'error');
		return null;
	}
	return data;
}

// getRekapMingguan
export async function getRekapMingguan(
	senin: string,
	kelas?: string,
	page = 1,
	limit = 0
): Promise<RekapMingguanKelas | null> {
	const query = new URLSearchParams({ senin, page: String(page), limit: String(limit) });
	if (kelas) query.append('kelas', kelas);

	const { data, error } = await apiRequest<RekapMingguanKelas>(
		`/v1/admin/rekap-mingguan?${query.toString()}`
	);
	if (error) {
		addToast('Gagal memuat rekap mingguan kehadiran', 'error');
		return null;
	}
	return data;
}

// getRekapSiswaDetail
export async function getRekapSiswaDetail(params: {
	nis?: string;
	nisn?: string;
	tipe: 'bulanan' | 'mingguan';
	bulan?: string;
	senin?: string;
}): Promise<KehadiranBulananSiswa | null> {
	const query = new URLSearchParams({
		nis: params.nis || params.nisn || '',
		tipe: params.tipe
	});
	if (params.bulan) query.append('bulan', params.bulan);
	if (params.senin) query.append('senin', params.senin);

	const { data, error } = await apiRequest<KehadiranBulananSiswa>(
		`/v1/admin/rekap-siswa-detail?${query.toString()}`
	);
	if (error) {
		addToast('Gagal memuat detail rekap siswa', 'error');
		return null;
	}
	return data;
}

// getG7RekapKelas
export async function getG7RekapKelas(
	bulan: string,
	kelas: string
): Promise<G7RekapKelasLengkap | null> {
	const query = new URLSearchParams({ bulan, kelas });

	const { data, error } = await apiRequest<G7RekapKelasLengkap>(
		`/v1/admin/g7/rekap-kelas?${query.toString()}`
	);
	if (error) {
		addToast('Gagal memuat rekap kelas G7', 'error');
		return null;
	}
	return data;
}

export async function submitAbsenAdmin(params: {
	nis?: string;
	nisn?: string;
	status: 'hadir' | 'tidak_hadir' | 'izin' | 'sakit' | 'magang';
	tanggal: string;
	alasan?: string;
	fotoIzin?: string;
}): Promise<boolean> {
	const bodyPayload = {
		...params,
		nisn: params.nis || params.nisn
	};
	const { data, error } = await apiRequest('/v1/admin/kehadiran', {
		method: 'POST',
		body: JSON.stringify(bodyPayload)
	});

	if (error) {
		addToast(error, 'error');
		return false;
	}

	addToast('Data absensi siswa berhasil direkam oleh Admin', 'success');
	return true;
}

export async function listKehadiranAdmin(filters: Record<string, string>, page = 1, limit = 20) {
	const query = new URLSearchParams();
	Object.entries(filters).forEach(([k, v]) => {
		if (v) query.append(k, v);
	});
	query.append('page', String(page));
	query.append('limit', String(limit));

	const { data, error } = await apiRequest<any>(`/v1/admin/kehadiran?${query.toString()}`);
	if (error) return { items: [], total: 0, hasMore: false };
	return { items: data.items || [], total: data.total || 0, hasMore: data.hasMore || false };
}

export async function deleteKehadiran(id: string): Promise<boolean> {
	const { error } = await apiRequest('/v1/admin/kehadiran', {
		method: 'DELETE',
		body: JSON.stringify({ id })
	});

	if (error) {
		addToast(error, 'error');
		return false;
	}

	addToast('Data kehadiran berhasil dihapus', 'success');
	return true;
}

export async function uploadIzinFileAdmin(file: File, nis: string): Promise<string | null> {
	const formData = new FormData();
	formData.append('foto', file);
	const { data, error } = await apiUpload<{ url: string }>(
		`/v1/admin/kehadiran/upload-izin?nis=${nis || 'unknown'}`,
		formData
	);
	if (error || !data) {
		addToast(error ?? 'Gagal mengunggah file izin', 'error');
		return null;
	}
	addToast('Berkas bukti izin berhasil diunggah', 'success');
	return data.url;
}

// G7 Grading & Rekap admin handlers
export async function listG7Rekap(filters: Record<string, string>, page = 1, limit = 20) {
	const query = new URLSearchParams();
	Object.entries(filters).forEach(([k, v]) => {
		if (v) query.append(k, v);
	});
	query.append('page', String(page));
	query.append('limit', String(limit));

	const { data, error } = await apiRequest<any>(`/v1/admin/g7/rekap?${query.toString()}`);
	if (error) return { items: [], total: 0, hasMore: false };
	return {
		items: (data.items || []) as G7Rekap[],
		total: data.total || 0,
		hasMore: data.hasMore || false
	};
}

export async function getG7RekapDetail(nis: string, bulan: string): Promise<G7Rekap | null> {
	const { data, error } = await apiRequest<G7Rekap>(`/v1/admin/g7/rekap/${nis}/${bulan}`);
	if (error) return null;
	return data;
}

export async function getG7EvaluateDetail(nis: string, bulan: string): Promise<any | null> {
	const { data, error } = await apiRequest<any>(`/v1/admin/g7/evaluate/${nis}/${bulan}`);
	if (error) return null;
	return data;
}

export async function getG7Suggest(nis: string, bulan: string): Promise<G7SuggestResponse | null> {
	const { data, error } = await apiRequest<G7SuggestResponse>(
		`/v1/admin/g7/suggest/${nis}/${bulan}`
	);
	if (error) return null;
	return data;
}

export async function saveG7Rekap(payload: Partial<G7Rekap>): Promise<boolean> {
	const { data, error } = await apiRequest('/v1/admin/g7/rekap', {
		method: 'POST',
		body: JSON.stringify(payload)
	});

	if (error) {
		addToast(error, 'error');
		return false;
	}

	addToast('Penilaian rekap G7 berhasil disimpan!', 'success');
	return true;
}

export async function getG7Statistik(
	bulan: string,
	kelas?: string
): Promise<G7RekapStatistik | null> {
	const query = new URLSearchParams({ bulan });
	if (kelas) query.append('kelas', kelas);

	const { data, error } = await apiRequest<G7RekapStatistik>(
		`/v1/admin/g7/statistik?${query.toString()}`
	);
	if (error) return null;
	return data;
}

export async function getG7RekapSemester(
	semester: string,
	kelas?: string
): Promise<any[] | null> {
	const query = new URLSearchParams({ semester });
	if (kelas) query.append('kelas', kelas);

	const { data, error } = await apiRequest<any[]>(
		`/v1/admin/g7/rekap-semester?${query.toString()}`
	);
	if (error) return null;
	return data;
}

// Student & Admin CRUD
export async function listStudents(page = 1, limit = 20, queryStr = '', kelas = '') {
	const q = new URLSearchParams({ page: String(page), limit: String(limit) });
	if (queryStr) q.append('q', queryStr);
	if (kelas) q.append('kelas', kelas);

	const { data, error } = await apiRequest<any>(`/v1/admin/students?${q.toString()}`);
	if (error) return { items: [], total: 0, hasMore: false };
	return {
		items: (data.items || []) as Siswa[],
		total: data.total || 0,
		hasMore: data.hasMore || false
	};
}

export async function createStudent(student: Partial<Siswa>): Promise<boolean> {
	const { error } = await apiRequest('/v1/admin/students', {
		method: 'POST',
		body: JSON.stringify(student)
	});
	if (error) {
		addToast(error, 'error');
		return false;
	}
	addToast('Siswa berhasil ditambahkan', 'success');
	return true;
}

export async function updateStudent(id: string, student: Partial<Siswa>): Promise<boolean> {
	const { error } = await apiRequest(`/v1/admin/students/${id}`, {
		method: 'PUT',
		body: JSON.stringify(student)
	});
	if (error) {
		addToast(error, 'error');
		return false;
	}
	addToast('Siswa berhasil diperbarui', 'success');
	return true;
}

export async function deleteStudent(id: string): Promise<boolean> {
	const { error } = await apiRequest(`/v1/admin/students/${id}`, { method: 'DELETE' });
	if (error) {
		addToast(error, 'error');
		return false;
	}
	addToast('Siswa berhasil dihapus', 'success');
	return true;
}

export async function bulkImportStudents(items: any[]): Promise<boolean> {
	const { error } = await apiRequest('/v1/admin/tambah-siswa/bulk', {
		method: 'POST',
		body: JSON.stringify({ items })
	});
	if (error) {
		addToast(error, 'error');
		return false;
	}
	addToast('Bulk import siswa berhasil!', 'success');
	return true;
}

export async function listAdmins(page = 1, limit = 20) {
	const q = new URLSearchParams({ page: String(page), limit: String(limit) });
	const { data, error } = await apiRequest<any>(`/v1/admin/admins?${q.toString()}`);
	if (error) return { items: [], total: 0, hasMore: false };
	return {
		items: (data.items || []) as Admin[],
		total: data.total || 0,
		hasMore: data.hasMore || false
	};
}

export async function createAdmin(admin: any): Promise<boolean> {
	const { error } = await apiRequest('/v1/admin/admins', {
		method: 'POST',
		body: JSON.stringify(admin)
	});
	if (error) {
		addToast(error, 'error');
		return false;
	}
	addToast('Admin/Guru berhasil ditambahkan', 'success');
	return true;
}

export async function deleteAdmin(id: string): Promise<boolean> {
	const { error } = await apiRequest(`/v1/admin/admins/${id}`, { method: 'DELETE' });
	if (error) {
		addToast(error, 'error');
		return false;
	}
	addToast('Admin/Guru berhasil dihapus', 'success');
	return true;
}

export async function getKelasJurusan(): Promise<{
	kelas: string[];
	jurusan: string[];
	kelasLengkap: string[];
} | null> {
	const { data, error } = await apiRequest<{
		kelas: string[];
		jurusan: string[];
		kelasLengkap: string[];
	}>('/v1/admin/kelas-jurusan');
	if (error) {
		addToast('Gagal memuat data kelas dan jurusan', 'error');
		return null;
	}
	return data;
}

export async function getKelas(): Promise<string[] | null> {
	const { data, error } = await apiRequest<string[]>('/v1/admin/kelas');
	if (error) {
		addToast('Gagal memuat data kelas', 'error');
		return null;
	}
	return data;
}
