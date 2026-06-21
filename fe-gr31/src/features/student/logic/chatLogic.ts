import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { Aduan } from '../types/student.types';

export async function createAduan(isi: string): Promise<Aduan | null> {
  if (!isi.trim()) {
    addToast('Pesan konsultasi tidak boleh kosong', 'warning');
    return null;
  }

  const { data, error } = await apiRequest<Aduan>('/v1/student/aduan', {
    method: 'POST',
    body: JSON.stringify({ isi }),
  });

  if (error || !data) {
    addToast(error ?? 'Gagal mengirim aduan', 'error');
    return null;
  }

  addToast('Aduan/konsultasi berhasil dibuat!', 'success');
  return data;
}

export async function listAduanSiswa(status?: 'open' | 'in_progress' | 'closed', page = 1, limit = 10) {
  const query = new URLSearchParams();
  if (status) query.append('status', status);
  query.append('page', String(page));
  query.append('limit', String(limit));

  const { data, error } = await apiRequest<any>(`/v1/student/aduan?${query.toString()}`);
  if (error) {
    addToast('Gagal memuat daftar aduan/konsultasi', 'error');
    return { items: [], total: 0, hasMore: false };
  }

  return {
    items: (data.items || []) as Aduan[],
    total: data.total || 0,
    hasMore: data.hasMore || false,
  };
}
