import { writable, get } from 'svelte/store';
import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import { getKelas } from '../../admin/logic/adminLogic';
import type { LombaKebersihan } from '../../student/types/student.types';

export const logs = writable<LombaKebersihan[]>([]);
export const total = writable<number>(0);
export const page = writable<number>(1);
export const limit = writable<number>(10);
export const loading = writable<boolean>(false);
export const hasMore = writable<boolean>(false);

export const kelasList = writable<string[]>([]);

export const selectedKelas = writable<string>('');
export const selectedTanggal = writable<string>('');
export const selectedBulan = writable<string>(''); // YYYY-MM format

export async function loadKelasList() {
  try {
    const data = await getKelas();
    if (data) {
      kelasList.set(data || []);
    }
  } catch (err) {
    console.error('Error loading kelas:', err);
  }
}

export async function loadLombaData() {
  loading.set(true);
  try {
    const query = new URLSearchParams();
    const kelas = get(selectedKelas);
    const tanggal = get(selectedTanggal);
    const bulan = get(selectedBulan);

    if (kelas) query.append('kelas', kelas);
    if (tanggal) query.append('tanggal', tanggal);
    if (bulan) {
      const [year, month] = bulan.split('-');
      const lastDay = new Date(Number(year), Number(month), 0).getDate();
      query.append('dari', `${bulan}-01`);
      query.append('sampai', `${bulan}-${String(lastDay).padStart(2, '0')}`);
    }
    
    query.append('page', String(get(page)));
    query.append('limit', String(get(limit)));

    const { data, error } = await apiRequest<any>(`/v1/admin/lomba?${query.toString()}`);
    if (error) {
      addToast('Gagal memuat data laporan kebersihan', 'error');
      logs.set([]);
      total.set(0);
      hasMore.set(false);
    } else {
      logs.set(data.items || []);
      total.set(data.total || 0);
      hasMore.set(data.hasMore || false);
    }
  } catch (err) {
    console.error('Error loading lomba data:', err);
    addToast('Gagal memuat data laporan kebersihan', 'error');
  } finally {
    loading.set(false);
  }
}

export function handleFilter() {
  page.set(1);
  loadLombaData();
}

export function changePage(delta: number) {
  const nextPage = get(page) + delta;
  if (nextPage < 1) return;
  page.set(nextPage);
  loadLombaData();
}
