import { apiRequest, apiUpload } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { Kehadiran, LatLng } from '../types/student.types';

// Types

export interface KehadiranHariItem {
	tanggal: string;
	status: 'hadir' | 'izin' | 'sakit' | 'tidak_hadir' | 'magang';
	waktuAbsen: string;
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

// RekapBulanan 
export interface RekapBulanan {
	id: string;
	rekapKey: string;
	nis: string;
	namaSiswa: string;
	kelas: string;
	bulanTahun: string;
	semester?: string;
	totalHadir: number;
	totalIzin: number;
	totalSakit: number;
	totalTidakHadir: number;
	totalMagang: number;
	totalHariEfektif: number;
	persentaseHadir: number;
}

// API calls

// Mengambil rekap kehadiran
export async function getRekapBulananSiswa(bulan: string): Promise<RekapBulanan | null> {
	const { data, error } = await apiRequest<RekapBulanan>(
		`/v1/student/kehadiran/rekap?bulan=${bulan}`
	);
	if (error) return null;
	return data;
}

// Mengambil rekap beberapa bulan
export async function getRekapRentangBulan(dari: string, sampai: string): Promise<RekapBulanan[]> {
	const { data, error } = await apiRequest<RekapBulanan[]>(
		`/v1/student/kehadiran/rekap?dari=${dari}&sampai=${sampai}`
	);
	if (error) return [];
	return data ?? [];
}

// Mengambil kehadiran bulanan siswa
export async function getKehadiranBulanan(bulan?: string): Promise<KehadiranBulananSiswa | null> {
	const query = new URLSearchParams();
	if (bulan) query.append('bulan', bulan);
	const { data, error } = await apiRequest<KehadiranBulananSiswa>(
		`/v1/student/kehadiran/bulanan?${query.toString()}`
	);
	if (error) return null;
	return data;
}

export async function uploadIzinFile(file: File): Promise<string | null> {
	const formData = new FormData();
	formData.append('foto', file);
	const { data, error } = await apiUpload<{ url: string }>(
		'/v1/student/kehadiran/upload-izin',
		formData
	);
	if (error || !data) {
		addToast(error ?? 'Gagal mengunggah file izin', 'error');
		return null;
	}
	addToast('Berkas bukti izin berhasil diunggah', 'success');
	return data.url;
}

const MIN_GPS_ACCURACY = 3;
const MAX_GPS_ACCURACY = 100;
const GPS_READING_INTERVAL_MS = 800;
const MAX_DRIFT_METERS = 500;

function gpsGetPosition(options: PositionOptions): Promise<GeolocationPosition> {
	return new Promise((resolve, reject) =>
		navigator.geolocation.getCurrentPosition(resolve, reject, options)
	);
}

function haversineMeters(lat1: number, lng1: number, lat2: number, lng2: number): number {
	const R = 6371000;
	const dLat = ((lat2 - lat1) * Math.PI) / 180;
	const dLng = ((lng2 - lng1) * Math.PI) / 180;
	const a =
		Math.sin(dLat / 2) ** 2 +
		Math.cos((lat1 * Math.PI) / 180) * Math.cos((lat2 * Math.PI) / 180) * Math.sin(dLng / 2) ** 2;
	return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
}

/**
 * Ambil lokasi GPS dengan validasi anti-fake GPS.
 * Strategi: dual-reading + accuracy check + drift check.
 * Return type identik dengan versi sebelumnya agar tidak break caller.
 */
export async function getGeolocation(): Promise<{ koordinat: LatLng; akurasi: number } | null> {
	if (typeof window === 'undefined' || !navigator.geolocation) {
		addToast('GPS tidak didukung di peramban ini', 'warning');
		return null;
	}

	const opts: PositionOptions = { enableHighAccuracy: true, timeout: 10000, maximumAge: 0 };

	try {
		// Reading pertama
		const r1 = await gpsGetPosition(opts);

		// Jeda singkat → GPS real akan menghasilkan koordinat sedikit berbeda (natural drift)
		await new Promise((res) => setTimeout(res, GPS_READING_INTERVAL_MS));

		// Reading kedua
		const r2 = await gpsGetPosition(opts);

		const c1 = r1.coords;
		const c2 = r2.coords;

		// Validasi 1: Accuracy terlalu sempurna → fake GPS
		if (c1.accuracy < MIN_GPS_ACCURACY || c2.accuracy < MIN_GPS_ACCURACY) {
			addToast(
				'GPS accuracy tidak realistis. Pastikan kamu tidak menggunakan aplikasi mock location.',
				'error'
			);
			return null;
		}

		// Validasi 2: Accuracy terlalu buruk → sinyal tidak memadai
		if (c1.accuracy >= MAX_GPS_ACCURACY || c2.accuracy >= MAX_GPS_ACCURACY) {
			addToast(
				'Sinyal GPS terlalu lemah. Pindah ke area terbuka dan coba lagi.',
				'error'
			);
			return null;
		}

		// Validasi 3: Koordinat identik sempurna → mock location (GPS real selalu micro-drift)
		// Catatan: Dilewati jika di localhost/development environment atau di PC/Laptop karena tidak memiliki sensor GPS fisik (menggunakan IP/Wi-Fi geoloc)
		const isDevMode = typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.DEV;
		const isDevelopment = isDevMode || (typeof window !== 'undefined' && 
			(window.location.hostname === 'localhost' || 
			 window.location.hostname === '127.0.0.1' || 
			 window.location.hostname.startsWith('192.168.')));
		const isPC = getBrowserDeviceInfo().platform === 'web';

		if (c1.latitude === c2.latitude && c1.longitude === c2.longitude) {
			if (!isDevelopment && !isPC) {
				addToast(
					'Lokasi GPS tidak bergerak sama sekali. Pastikan kamu tidak menggunakan aplikasi fake GPS.',
					'error'
				);
				return null;
			} else {
				console.warn('GPS micro-drift check bypassed on localhost/PC environment');
			}
		}

		// Validasi 4: Koordinat berubah terlalu jauh dalam < 1 detik → teleportasi / fake GPS
		const drift = haversineMeters(c1.latitude, c1.longitude, c2.latitude, c2.longitude);
		if (drift > MAX_DRIFT_METERS) {
			addToast(
				'Lokasi GPS berubah tidak wajar. Matikan aplikasi fake GPS dan coba lagi.',
				'error'
			);
			return null;
		}

		// Gunakan reading ke-2 (lebih settled) sebagai koordinat final
		return {
			koordinat: { lat: c2.latitude, lng: c2.longitude },
			akurasi: c2.accuracy
		};
	} catch (err: any) {
		const code = err?.code;
		if (code === 1 /* PERMISSION_DENIED */) {
			addToast('Izin lokasi ditolak. Harap izinkan akses GPS untuk absen.', 'error');
		} else if (code === 2 /* POSITION_UNAVAILABLE */) {
			addToast('Lokasi tidak tersedia. Pastikan GPS aktif dan kamu berada di area dengan sinyal.', 'error');
		} else if (code === 3 /* TIMEOUT */) {
			addToast('Waktu GPS habis. Pindah ke area terbuka dan coba lagi.', 'error');
		} else {
			addToast('Gagal mengakses lokasi GPS.', 'error');
		}
		return null;
	}
}

export function getSpecificIphoneModel(): string {
	if (typeof window === 'undefined') return 'iPhone';
	const w = Math.min(window.screen.width, window.screen.height);
	const h = Math.max(window.screen.width, window.screen.height);
	const dpr = window.devicePixelRatio || 1;

	// Check matching resolution (logical viewport size + device pixel ratio)
	if (w === 440 && h === 956 && dpr === 3) return 'iPhone 16 Pro Max';
	if (w === 402 && h === 874 && dpr === 3) return 'iPhone 16 Pro';
	if (w === 393 && h === 852 && dpr === 3) return 'iPhone 14 Pro / 15 Pro';
	if (w === 430 && h === 932 && dpr === 3) return 'iPhone 14 Pro Max / 15 Pro Max';
	if (w === 428 && h === 926 && dpr === 3) return 'iPhone 12/13 Pro Max / 15 Plus';
	if (w === 390 && h === 844 && dpr === 3) return 'iPhone 12/13/14/15';
	if (w === 360 && h === 780 && dpr === 3) return 'iPhone 12/13 mini';
	if (w === 414 && h === 896 && dpr === 3) return 'iPhone XS Max / 11 Pro Max';
	if (w === 414 && h === 896 && dpr === 2) return 'iPhone XR / 11';
	if (w === 375 && h === 812 && dpr === 3) return 'iPhone X / XS / 11 Pro';
	if (w === 414 && h === 736 && dpr === 3) return 'iPhone 6s/7/8 Plus';
	if (w === 375 && h === 667 && dpr === 2) return 'iPhone 6/7/8 / SE (2nd/3rd Gen)';
	if (w === 320 && h === 568 && dpr === 2) return 'iPhone 5/5s/5c / SE (1st Gen)';

	return 'iPhone';
}

export function getBrowserDeviceInfo(): { model: string; platform: string; osVersion: string; appVersion: string } {
	if (typeof window === 'undefined' || typeof navigator === 'undefined') {
		return { model: 'Unknown', platform: 'web', osVersion: 'Unknown', appVersion: '1.0.0' };
	}

	const ua = navigator.userAgent;
	let platform = 'web';
	let model = 'Web Browser';
	let osVersion = 'Unknown';

	// Detect iOS/iPhone/iPad
	if (/iPhone|iPad|iPod/.test(ua)) {
		platform = 'ios';
		if (/iPad/.test(ua)) {
			model = 'iPad';
		} else if (/iPod/.test(ua)) {
			model = 'iPod';
		} else {
			model = getSpecificIphoneModel();
		}
		
		const match = ua.match(/OS (\d+[._]\d+)/);
		if (match) {
			osVersion = 'iOS ' + match[1].replace(/_/g, '.');
		} else {
			osVersion = 'iOS';
		}
	}
	// Detect Android
	else if (/Android/.test(ua)) {
		platform = 'android';
		const osMatch = ua.match(/Android\s+([0-9\.]+)/);
		if (osMatch) {
			osVersion = 'Android ' + osMatch[1];
		} else {
			osVersion = 'Android';
		}

		// Try robust extraction of Android model name
		model = 'Android Device';
		const matchParentheses = ua.match(/\(([^)]+)\)/);
		if (matchParentheses && matchParentheses[1]) {
			const parts = matchParentheses[1].split(';').map(p => p.trim());
			const androidIndex = parts.findIndex(p => p.toLowerCase().includes('android'));
			if (androidIndex !== -1) {
				let possibleModel = '';
				for (let i = androidIndex + 1; i < parts.length; i++) {
					const p = parts[i];
					// Skip locales, 'U', 'wv' webview identifier, 'Mobile' etc.
					if (/^[a-z]{2}-[a-z]{2}$/i.test(p) || p === 'U' || p === 'wv' || p === 'Mobile') {
						continue;
					}
					possibleModel = p;
				}
				if (possibleModel) {
					model = possibleModel.split(/\bBuild\b/i)[0].trim();
				}
			}
		}
	}
	// Detect Windows PC
	else if (/Windows NT/.test(ua)) {
		model = 'PC / Laptop';
		platform = 'web';
		const match = ua.match(/Windows NT\s+([0-9\.]+)/);
		if (match) {
			const versionMap: Record<string, string> = {
				'10.0': '10/11',
				'6.3': '8.1',
				'6.2': '8',
				'6.1': '7'
			};
			osVersion = 'Windows ' + (versionMap[match[1]] || match[1]);
		} else {
			osVersion = 'Windows';
		}
	}
	// Detect Mac
	else if (/Macintosh/.test(ua)) {
		model = 'Mac';
		platform = 'web';
		const match = ua.match(/Mac OS X\s+([0-9\._]+)/);
		if (match) {
			osVersion = 'macOS ' + match[1].replace(/_/g, '.');
		} else {
			osVersion = 'macOS';
		}
	}

	return {
		model,
		platform,
		osVersion,
		appVersion: '1.0.0'
	};
}

export async function getBrowserDeviceInfoAsync(): Promise<{ model: string; platform: string; osVersion: string; appVersion: string }> {
	const info = getBrowserDeviceInfo();

	if (typeof navigator !== 'undefined' && (navigator as any).userAgentData) {
		try {
			const hints = await (navigator as any).userAgentData.getHighEntropyValues(['model', 'platformVersion']);
			if (hints.model) {
				info.model = hints.model;
			}
			if (hints.platformVersion) {
				if (info.platform === 'android') {
					info.osVersion = 'Android ' + hints.platformVersion;
				}
			}
		} catch (e) {
			console.warn('Error getting UA client hints:', e);
		}
	}

	return info;
}

export async function submitKehadiran(params: {
	status: 'hadir' | 'tidak_hadir' | 'izin' | 'sakit' | 'magang';
	alasan?: string;
	koordinat?: LatLng;
	akurasi?: number;
	fotoIzin?: string;
	tipe?: string;
}): Promise<boolean> {
	const deviceInfo = await getBrowserDeviceInfoAsync();
	const payload = {
		...params,
		deviceInfo
	};

	const { error } = await apiRequest('/v1/student/kehadiran', {
		method: 'POST',
		body: JSON.stringify(payload)
	});
	if (error) {
		addToast(error, 'error');
		return false;
	}
	addToast('Absensi berhasil direkam!', 'success');
	
	setTimeout(() => {
		if (typeof window !== 'undefined') {
			window.location.reload();
		}
	}, 1500);

	return true;
}

export async function listKehadiranSiswa(dari?: string, sampai?: string, page = 1, limit = 10) {
	const query = new URLSearchParams();
	if (dari) query.append('dari', dari);
	if (sampai) query.append('sampai', sampai);
	query.append('page', String(page));
	query.append('limit', String(limit));
	const { data, error } = await apiRequest<any>(`/v1/student/kehadiran?${query.toString()}`);
	if (error) {
		addToast('Gagal mengambil data riwayat kehadiran', 'error');
		return { items: [], total: 0, hasMore: false };
	}
	return {
		items: data.items || [],
		total: data.total || 0,
		hasMore: data.hasMore || false
	};
}
