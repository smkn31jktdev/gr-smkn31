import { writable, get } from 'svelte/store';
import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { Admin } from '../../auth/types/auth.types';

export const adminsList = writable<Admin[]>([]);
export const adminsLoading = writable<boolean>(false);

// Form state
export const formNama = writable<string>('');
export const formEmail = writable<string>('');
export const formRole = writable<string>('walas');
export const formKelas = writable<string>('');
export const formPassword = writable<string>('');
export const formSubmitting = writable<boolean>(false);

// Sheets state
export const sheetUrl = writable<string>('');
export const sheetLoading = writable<boolean>(false);
export const sheetAdminsList = writable<any[]>([]);
export const sheetSubmitting = writable<boolean>(false);

export async function loadAdmins() {
  adminsLoading.set(true);
  try {
    const { data, error } = await apiRequest<any>('/v1/admin/admins?limit=100');
    if (!error && data) {
      adminsList.set(data.items || []);
    }
  } catch (err) {
    console.error('Error loading admins:', err);
  } finally {
    adminsLoading.set(false);
  }
}

export async function submitManualAdmin(): Promise<boolean> {
  const nama = get(formNama);
  const email = get(formEmail);
  const role = get(formRole);
  const kelas = get(formKelas);
  const password = get(formPassword);

  if (!nama || !email || !password) {
    addToast('Semua field form wajib diisi', 'warning');
    return false;
  }

  formSubmitting.set(true);
  try {
    let apiRole = role;
    if (apiRole === 'guru_non_walas') {
      apiRole = 'admin';
    }

    const { error } = await apiRequest('/v1/admin/admins', {
      method: 'POST',
      body: JSON.stringify({ nama, email, role: apiRole, password, kelas })
    });

    if (error) {
      addToast(error, 'error');
      return false;
    }

    addToast('Admin/Guru berhasil ditambahkan', 'success');
    formNama.set('');
    formEmail.set('');
    formRole.set('walas');
    formPassword.set('');
    await loadAdmins();
    return true;
  } catch (err) {
    console.error(err);
    addToast('Gagal menambahkan admin', 'error');
    return false;
  } finally {
    formSubmitting.set(false);
  }
}

export async function removeAdmin(id: string): Promise<boolean> {
  if (!confirm('Apakah Anda yakin ingin menghapus data admin/guru ini?')) return false;

  try {
    const { error } = await apiRequest(`/v1/admin/admins/${id}`, {
      method: 'DELETE'
    });

    if (error) {
      addToast(error, 'error');
      return false;
    }

    addToast('Admin/Guru berhasil dihapus', 'success');
    await loadAdmins();
    return true;
  } catch (err) {
    console.error(err);
    addToast('Gagal menghapus admin', 'error');
    return false;
  }
}

export async function updateAdminFields(id: string, fields: { isWalas?: boolean; kelas?: string; role?: string }): Promise<boolean> {
  try {
    const { error } = await apiRequest(`/v1/admin/admins/${id}`, {
      method: 'PUT',
      body: JSON.stringify(fields)
    });

    if (error) {
      addToast(error, 'error');
      return false;
    }

    addToast('Hak akses/walas berhasil diperbarui', 'success');
    await loadAdmins();
    return true;
  } catch (err) {
    console.error(err);
    addToast('Gagal memperbarui data', 'error');
    return false;
  }
}


// Client-side Google Sheet parsing
export async function loadSheetData() {
  const url = get(sheetUrl).trim();
  if (!url) {
    addToast('URL Google Sheet wajib diisi', 'warning');
    return;
  }

  sheetLoading.set(true);
  try {
    // Convert sharing link to CSV export URL
    let csvUrl = url;
    const match = url.match(/\/d\/([a-zA-Z0-9-_]+)/);
    if (match) {
      const docId = match[1];
      csvUrl = `https://docs.google.com/spreadsheets/d/${docId}/export?format=csv`;
    } else {
      addToast('Format URL Google Sheet tidak valid', 'error');
      sheetLoading.set(false);
      return;
    }

    const res = await fetch(csvUrl);
    if (!res.ok) {
      throw new Error('Gagal mendownload data sheet. Pastikan sheet diatur publik (akses siapa saja dengan link).');
    }
    const csvText = await res.text();
    
    // Parse CSV
    const lines = csvText.split(/\r?\n/);
    if (lines.length < 2) {
      addToast('Data sheet kosong atau tidak valid', 'warning');
      sheetLoading.set(false);
      return;
    }

    const headers = lines[0].split(',').map(h => h.trim().toLowerCase());
    const parsedItems = [];

    for (let i = 1; i < lines.length; i++) {
      const line = lines[i].trim();
      if (!line) continue;
      
      const cols = line.split(',').map(c => c.trim().replace(/^"|"$/g, ''));
      
      const namaIdx = headers.indexOf('nama') !== -1 ? headers.indexOf('nama') : 0;
      const emailIdx = headers.indexOf('email') !== -1 ? headers.indexOf('email') : 1;
      const roleIdx = headers.indexOf('role') !== -1 ? headers.indexOf('role') : 2;
      const passIdx = headers.indexOf('password') !== -1 ? headers.indexOf('password') : 3;

      let roleVal = cols[roleIdx] || 'walas';
      if (roleVal === 'guru_non_walas') {
        roleVal = 'admin';
      }

      const item = {
        nama: cols[namaIdx] || '',
        email: cols[emailIdx] || '',
        role: roleVal,
        password: cols[passIdx] || 'changeme123'
      };

      if (item.nama && item.email) {
        parsedItems.push(item);
      }
    }

    sheetAdminsList.set(parsedItems);
    addToast(`Berhasil memuat ${parsedItems.length} data admin dari sheet`, 'success');
  } catch (err: any) {
    console.error(err);
    addToast(err.message || 'Gagal memproses Google Sheet', 'error');
  } finally {
    sheetLoading.set(false);
  }
}

export async function submitSheetAdmins(): Promise<boolean> {
  const items = get(sheetAdminsList);
  if (items.length === 0) {
    addToast('Tidak ada data admin untuk disimpan', 'warning');
    return false;
  }

  sheetSubmitting.set(true);
  try {
    const { data, error } = await apiRequest<any>('/v1/admin/tambah-admin/bulk', {
      method: 'POST',
      body: JSON.stringify(items)
    });

    if (error) {
      addToast(error, 'error');
      return false;
    }

    const berhasil = data?.berhasil || 0;
    const gagal = data?.gagal || 0;
    addToast(`Import selesai. Berhasil: ${berhasil}, Gagal: ${gagal}`, 'success');

    sheetAdminsList.set([]);
    sheetUrl.set('');
    await loadAdmins();
    return true;
  } catch (err) {
    console.error(err);
    addToast('Gagal mengimpor data bulk admin', 'error');
    return false;
  } finally {
    sheetSubmitting.set(false);
  }
}
