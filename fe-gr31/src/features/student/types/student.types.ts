export interface Aktivitas {
	done: boolean;
	waktu?: string;
	keterangan?: string;
}

export interface G7Jurnal {
	id: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	tanggal: string;
	bangun?: Aktivitas;
	ibadah?: Aktivitas;
	makan?: Aktivitas;
	olahraga?: Aktivitas;
	belajar?: Aktivitas;
	bermasyarakat?: Aktivitas;
	tidur?: Aktivitas;
	totalDone: number;
}

export interface LatLng {
	lat: number;
	lng: number;
}

export interface Kehadiran {
	id: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	tanggal: string;
	hari: string;
	status: 'hadir' | 'tidak_hadir' | 'izin' | 'sakit' | 'magang';
	waktuAbsen: string;
	alasan?: string;
	koordinat?: LatLng;
	jarak: number;
	akurasi: number;
	fotoIzin?: string;
}

export interface Kegiatan {
	id: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	tanggal: string;
	section: string;
	payload: Record<string, any>;
	createdAt?: string;
	updatedAt?: string;
}

export interface Bukti {
	id: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	bulan: string; // YYYY-MM
	foto: string[];
	linkYT: string[];
}

export interface Message {
	role: 'student' | 'admin';
	isi: string;
	timestamp: string;
}

export interface Aduan {
	id: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	messages: Message[];
	status: 'open' | 'in_progress' | 'closed';
	adminNama?: string;
	wali?: string;
	createdAt: string;
	updatedAt: string;
}

export interface LombaKebersihan {
	id: string;
	kelas: string;
	nis: string;
	namaSiswa: string;
	tanggal: string;
	foto: string[];
	catatan: string;
	createdAt?: string;
	updatedAt?: string;
}

