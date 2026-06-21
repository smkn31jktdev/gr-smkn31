import { writable, get } from 'svelte/store';
import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { RekapHarian } from '../../admin/types/admin.types';

// State stores
export const loadingDashboard = writable<boolean>(false);
export const rekapHarian = writable<RekapHarian | null>(null);
export const currentTime = writable<string>('');
export const currentDateStr = writable<string>('');
export const countdownStr = writable<string>('00:00:00');

// Get today's YYYY-MM-DD date string based on local client timezone
export function getTodayDateString(): string {
  const d = new Date();
  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

export function formatIndonesianDate(d: Date): string {
  const days = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
  const months = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];
  return `${days[d.getDay()]}, ${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`;
}

export function updateClockAndCountdown() {
  const now = new Date();
  
  // Format Time: HH.MM.SS WIB
  const hours = String(now.getHours()).padStart(2, '0');
  const minutes = String(now.getMinutes()).padStart(2, '0');
  const seconds = String(now.getSeconds()).padStart(2, '0');
  currentTime.set(`${hours}.${minutes}.${seconds} WIB`);
  
  // Format Date
  currentDateStr.set(formatIndonesianDate(now));
  
  // Calculate Countdown to next midnight
  const midnight = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1, 0, 0, 0);
  const diffMs = midnight.getTime() - now.getTime();
  const diffSecs = Math.max(0, Math.floor(diffMs / 1000));
  
  const cHours = Math.floor(diffSecs / 3600);
  const cMinutes = Math.floor((diffSecs % 3600) / 60);
  const cSeconds = diffSecs % 60;
  
  countdownStr.set(
    [cHours, cMinutes, cSeconds]
      .map(v => String(v).padStart(2, '0'))
      .join(':')
  );
}

export async function loadPiketDashboardData() {
  loadingDashboard.set(true);
  try {
    const today = getTodayDateString();
    const { data, error } = await apiRequest<RekapHarian>(`/v1/admin/rekap-harian?tanggal=${today}`);
    if (error) {
      addToast(error || 'Gagal memuat rekap harian', 'error');
      return;
    }
    rekapHarian.set(data || null);
  } catch (err) {
    console.error('Error loading Piket dashboard:', err);
    addToast('Gagal memuat rekap harian', 'error');
  } finally {
    loadingDashboard.set(false);
  }
}
