import { apiRequest } from '../../../api/client';
import { setAuth, clearAuth } from '../../../stores/authStore';
import { addToast } from '../../../stores/uiStore';
import type { SiswaLoginResponse, AdminLoginResponse } from '../types/auth.types';

export async function loginSiswa(nis: string, password: string): Promise<boolean> {
  if (!nis || !password) {
    addToast('NIS dan password harus diisi', 'error');
    return false;
  }

  const { data, error } = await apiRequest<SiswaLoginResponse>('/v1/student/login', {
    method: 'POST',
    body: JSON.stringify({ nis, password }),
  });

  if (error || !data) {
    addToast(error ?? 'Login gagal. Cek NIS dan password Anda', 'error');
    return false;
  }

  // Set the user in store (student has role 'student' inside token, we'll assign it here for client compatibility)
  const user = {
    ...data.siswa,
    role: 'student',
  };
  setAuth(user, data.accessToken);
  addToast('Login siswa berhasil!', 'success');
  return true;
}

export async function loginAdmin(email: string, password: string): Promise<boolean> {
  if (!email || !password) {
    addToast('Email dan password harus diisi', 'error');
    return false;
  }

  const { data, error } = await apiRequest<AdminLoginResponse>('/v1/admin/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });

  if (error || !data) {
    addToast(error ?? 'Login gagal. Cek email dan password Anda', 'error');
    return false;
  }

  setAuth(data.admin, data.accessToken);
  addToast(`Login berhasil! Selamat datang, ${data.admin.nama}`, 'success');
  return true;
}

export function logout() {
  clearAuth();
  addToast('Anda telah keluar', 'success');
}
