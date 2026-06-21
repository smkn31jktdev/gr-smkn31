import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';

// Bukti upload handlers
export async function createBukti(bulan: string, foto: string[], linkYT: string[]): Promise<boolean> {
  const { data, error } = await apiRequest('/v1/student/bukti', {
    method: 'POST',
    body: JSON.stringify({ bulan, foto, linkYT }),
  });

  if (error) {
    addToast(error, 'error');
    return false;
  }

  addToast('Bukti kegiatan berhasil dikirim!', 'success');
  return true;
}

export async function listBuktiSiswa(bulan?: string, page = 1, limit = 10) {
  const query = new URLSearchParams();
  if (bulan) query.append('bulan', bulan);
  query.append('page', String(page));
  query.append('limit', String(limit));

  const { data, error } = await apiRequest<any>(`/v1/student/bukti?${query.toString()}`);
  if (error) {
    addToast('Gagal memuat riwayat upload bukti', 'error');
    return { items: [], total: 0, hasMore: false };
  }

  return {
    items: data.items || [],
    total: data.total || 0,
    hasMore: data.hasMore || false,
  };
}
