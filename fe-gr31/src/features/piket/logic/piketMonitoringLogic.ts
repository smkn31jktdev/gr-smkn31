import { writable, get } from 'svelte/store';
import { listKehadiranAdmin, submitAbsenAdmin, deleteKehadiran, getKelas } from '../../admin/logic/adminLogic';
import { addToast } from '../../../stores/uiStore';
import type { Kehadiran } from '../../student/types/student.types';

export const getTodayStr = () => new Date().toLocaleDateString('sv-SE');

// State stores
export const logs = writable<Kehadiran[]>([]);
export const total = writable<number>(0);
export const page = writable<number>(1);
export const limit = writable<number>(50);
export const loading = writable<boolean>(false);
export const hasMore = writable<boolean>(false);

// Dynamic dropdown lists from backend
export const kelasList = writable<string[]>([]);

// Filter stores
export const selectedKelas = writable<string>('');
export const selectedTanggal = writable<string>(getTodayStr());
export const searchQuery = writable<string>('');
export const selectedStatus = writable<string>('');
export const urutkanWaktu = writable<string>('default');
export const urutkanNama = writable<string>('default');

// Modal / Form state stores
export const showModal = writable<boolean>(false);
export const formState = writable({
  nis: '',
  status: 'hadir' as 'hadir' | 'tidak_hadir' | 'izin' | 'sakit',
  tanggal: getTodayStr(),
  alasan: ''
});

export async function loadKelasJurusan() {
  try {
    const data = await getKelas();
    if (data) {
      kelasList.set(data || []);
    }
  } catch (err) {
    console.error('Error loading kelas:', err);
  }
}

export async function loadData() {
  loading.set(true);
  try {
    const filters = {
      kelas: get(selectedKelas),
      tanggal: get(selectedTanggal),
      q: get(searchQuery),
      status: get(selectedStatus)
    };
    
    const currentPage = get(page);
    const currentLimit = get(limit);

    const res = await listKehadiranAdmin(filters, currentPage, currentLimit);
    logs.set(res.items || []);
    total.set(res.total || 0);
    hasMore.set(res.hasMore || false);
  } catch (err) {
    console.error('Error loading attendance logs:', err);
    addToast('Gagal memuat log kehadiran', 'error');
  } finally {
    loading.set(false);
  }
}


export function handleFilter() {
  page.set(1);
  loadData();
}

export function openCreate() {
  formState.set({
    nis: '',
    status: 'hadir',
    tanggal: getTodayStr(),
    alasan: ''
  });
  showModal.set(true);
}

export async function handleSave(handlers: { resolve: () => void; reject: () => void }) {
  const currentForm = get(formState);
  if (!currentForm.nis.trim()) {
    addToast('NIS siswa harus diisi', 'warning');
    handlers.reject();
    return;
  }

  const success = await submitAbsenAdmin(currentForm);
  if (success) {
    handlers.resolve();
    showModal.set(false);
    loadData();
  } else {
    handlers.reject();
  }
}

export async function handleDelete(id: string) {
  if (!confirm('Apakah Anda yakin ingin menghapus data absensi ini?')) return;
  const success = await deleteKehadiran(id);
  if (success) {
    loadData();
  }
}
