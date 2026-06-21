import { writable, get } from 'svelte/store';
import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import { listStudents } from './adminLogic';
import type { Siswa } from '../../auth/types/auth.types';

export const studentsList = writable<Siswa[]>([]);
export const loadingStudents = writable<boolean>(false);
export const deletingData = writable<boolean>(false);

// Muat data
export async function fetchStudentsForDelete() {
  loadingStudents.set(true);
  try {
    const res = await listStudents(1, 1000, '', '');
    studentsList.set(res.items || []);
  } catch (err) {
    console.error('Error fetching students for delete:', err);
    addToast('Gagal memuat daftar siswa', 'error');
  } finally {
    loadingStudents.set(false);
  }
}

export async function deleteStudentActivities(nis: string, rural: string): Promise<boolean> {
  const bulan = rural;
  if (!nis) {
    addToast('Pilih siswa terlebih dahulu', 'warning');
    return false;
  }
  if (!bulan) {
    addToast('Pilih bulan terlebih dahulu', 'warning');
    return false;
  }

  deletingData.set(true);
  try {
    // 1. Daily journals for student & month
    const { data: listData, error: listError } = await apiRequest<any>(
      `/v1/admin/g7?nis=${nis}&dari=${bulan}&sampai=${bulan}&limit=100`
    );

    if (listError) {
      addToast(listError || 'Gagal mengambil data kegiatan siswa', 'error');
      return false;
    }

    const items = listData?.items || [];
    
    // 2. Delete each journal one by one
    let deletedJournalsCount = 0;
    for (const item of items) {
      const { error: delError } = await apiRequest('/v1/admin/g7', {
        method: 'DELETE',
        body: JSON.stringify({ id: item.id })
      });
      if (!delError) {
        deletedJournalsCount++;
      }
    }

    // 3. Check and delete rekap if exists
    let deletedRekap = false;
    const { data: rekapData, error: rekapError } = await apiRequest<any>(
      `/v1/admin/g7/rekap/${nis}/${bulan}`
    );

    if (!rekapError && rekapData?.id) {
      const { error: delRekapError } = await apiRequest('/v1/admin/g7/rekap', {
        method: 'DELETE',
        body: JSON.stringify({ id: rekapData.id })
      });
      if (!delRekapError) {
        deletedRekap = true;
      }
    }

    if (deletedJournalsCount > 0 || deletedRekap) {
      addToast(
        `Berhasil menghapus ${deletedJournalsCount} catatan kegiatan${deletedRekap ? ' dan rekap bulanan' : ''} secara permanen`,
        'success'
      );
      return true;
    } else {
      addToast('Tidak ada data kegiatan ditemukan untuk siswa dan periode ini', 'info');
      return false;
    }
  } catch (err) {
    console.error('Error deleting student activities:', err);
    addToast('Gagal menghapus data kegiatan siswa', 'error');
    return false;
  } finally {
    deletingData.set(false);
  }
}
