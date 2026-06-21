import { writable, get } from 'svelte/store';
import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';

export interface Bukti {
  id: string;
  nis: string;
  namaSiswa: string;
  kelas: string;
  bulan: string;
  foto: string[];
  linkYT: string[];
}

// State
export const loadingBukti = writable<boolean>(false);
export const searchNIS = writable<string>('');
export const searchKelas = writable<string>('');
export const searchBulan = writable<string>('');
export const buktiList = writable<Bukti[]>([]);

export async function loadBuktiData() {
  loadingBukti.set(true);
  try {
    const query = new URLSearchParams();
    const nisVal = get(searchNIS);
    const kelasVal = get(searchKelas);
    const bulanVal = get(searchBulan);

    if (nisVal) query.append('nisn', nisVal);
    if (kelasVal) query.append('kelas', kelasVal);
    if (bulanVal) query.append('bulan', bulanVal);

    const { data, error } = await apiRequest<any>(`/v1/admin/bukti?${query.toString()}`);
    if (error) {
      addToast(error || 'Gagal memuat bukti kegiatan', 'error');
      return;
    }
    
    const items = (data?.items || []) as Bukti[];
    buktiList.set(items);
  } catch (err) {
    console.error('Error loading bukti:', err);
    addToast('Gagal memuat bukti kegiatan', 'error');
  } finally {
    loadingBukti.set(false);
  }
}
