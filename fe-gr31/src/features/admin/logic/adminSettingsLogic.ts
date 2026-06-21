import { writable, get } from 'svelte/store';
import { currentUser } from '../../../stores/authStore';
import { addToast } from '../../../stores/uiStore';

export const isSubmittingSettings = writable<boolean>(false);

export async function updateAdminProfile(nama: string, email: string, fotoProfil: string): Promise<boolean> {
  if (!nama.trim() || !email.trim()) {
    addToast('Nama Lengkap dan Email wajib diisi', 'warning');
    return false;
  }

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (!emailRegex.test(email)) {
    addToast('Format email tidak valid', 'warning');
    return false;
  }

  isSubmittingSettings.set(true);
  try {
    await new Promise((resolve) => setTimeout(resolve, 800));

    const current = get(currentUser);
    const updatedUser = {
      ...current,
      nama: nama.trim(),
      email: email.trim(),
      fotoProfil: fotoProfil
    };

    currentUser.set(updatedUser);
    localStorage.setItem('currentUser', JSON.stringify(updatedUser));
    
    addToast('Pengaturan akun berhasil disimpan!', 'success');
    return true;
  } catch (err) {
    console.error('Error updating admin profile:', err);
    addToast('Gagal menyimpan pengaturan akun', 'error');
    return false;
  } finally {
    isSubmittingSettings.set(false);
  }
}
