import { writable, get } from 'svelte/store';
import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { Aduan } from '../../student/types/student.types';

// State stores
export const loadingDashboard = writable<boolean>(false);
export const aduanList = writable<Aduan[]>([]);
export const totalAduan = writable<number>(0);
export const openAduan = writable<number>(0);
export const inProgressAduan = writable<number>(0);
export const closedAduan = writable<number>(0);
export const recentActiveAduan = writable<Aduan[]>([]);

export async function loadBkDashboardData() {
  loadingDashboard.set(true);
  try {
    const { data, error } = await apiRequest<any>('/v1/admin/aduan');
    if (error) {
      addToast(error || 'Gagal memuat data dashboard BK', 'error');
      return;
    }

    const items = (data?.items || []) as Aduan[];
    aduanList.set(items);

    // Compute stats
    totalAduan.set(items.length);
    openAduan.set(items.filter(a => a.status === 'open').length);
    inProgressAduan.set(items.filter(a => a.status === 'in_progress').length);
    closedAduan.set(items.filter(a => a.status === 'closed').length);

    // Filter and sort for the 5 most recent active aduan (open or in_progress first, sorted by updatedAt)
    const recent = [...items]
      .filter(a => a.status !== 'closed')
      .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
      .slice(0, 5);
      
    recentActiveAduan.set(recent);
  } catch (err) {
    console.error('Error loading BK dashboard data:', err);
    addToast('Gagal memuat data dashboard BK', 'error');
  } finally {
    loadingDashboard.set(false);
  }
}
