import { apiRequest, apiUpload } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { LombaKebersihan } from '../types/student.types';

// Lomba Kebersihan CRUD handlers
export async function createLomba(tanggal: string, foto: string[], catatan: string): Promise<boolean> {
  const { data, error } = await apiRequest('/v1/student/lomba', {
    method: 'POST',
    body: JSON.stringify({ tanggal, foto, catatan }),
  });

  if (error) {
    addToast(error, 'error');
    return false;
  }

  addToast('Foto kebersihan kelas berhasil dikirim!', 'success');
  return true;
}

export async function updateLomba(id: string, foto: string[], catatan: string): Promise<boolean> {
  const { data, error } = await apiRequest('/v1/student/lomba', {
    method: 'PUT',
    body: JSON.stringify({ id, foto, catatan }),
  });

  if (error) {
    addToast(error, 'error');
    return false;
  }

  addToast('Foto kebersihan kelas berhasil diperbarui!', 'success');
  return true;
}

export async function deleteLomba(id: string): Promise<boolean> {
  const { data, error } = await apiRequest(`/v1/student/lomba?id=${id}`, {
    method: 'DELETE',
  });

  if (error) {
    addToast(error, 'error');
    return false;
  }

  addToast('Foto kebersihan kelas berhasil dihapus', 'success');
  return true;
}

export async function listLombaSiswa(dari?: string, sampai?: string, page = 1, limit = 10) {
  const query = new URLSearchParams();
  if (dari) query.append('dari', dari);
  if (sampai) query.append('sampai', sampai);
  query.append('page', String(page));
  query.append('limit', String(limit));

  const { data, error } = await apiRequest<any>(`/v1/student/lomba?${query.toString()}`);
  if (error) {
    addToast('Gagal memuat daftar kebersihan kelas', 'error');
    return { items: [], total: 0, hasMore: false };
  }

  return {
    items: (data.items || []) as LombaKebersihan[],
    total: data.total || 0,
    hasMore: data.hasMore || false,
  };
}

export async function uploadLombaFile(file: File): Promise<string | null> {
  const formData = new FormData();
  formData.append('foto', file);

  const { data, error } = await apiUpload<{ url: string }>(
    '/v1/student/lomba/upload',
    formData
  );

  if (error || !data) {
    addToast(error ?? 'Gagal mengunggah foto kebersihan', 'error');
    return null;
  }

  return data.url;
}
