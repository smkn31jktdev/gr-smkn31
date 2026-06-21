import { writable, get, derived } from 'svelte/store';
import { apiRequest } from '../../../api/client';
import { addToast } from '../../../stores/uiStore';
import type { Aduan } from '../../student/types/student.types';

// State stores
export const rooms = writable<Aduan[]>([]);
export const activeRoom = writable<Aduan | null>(null);
export const listLoading = writable<boolean>(false);
export const replyText = writable<string>('');
export const selectedStatus = writable<string>('');
export const searchQuery = writable<string>('');

// Metrics
export const totalAduanCount = derived(rooms, ($rooms) => $rooms.length);
export const activeAduanCount = derived(rooms, ($rooms) => 
  $rooms.filter(r => r.status === 'open' || r.status === 'in_progress').length
);

// Filtered rooms
export const filteredRooms = derived(
  [rooms, selectedStatus, searchQuery],
  ([$rooms, $status, $search]) => {
    let result = [...$rooms];
    
    // Filter by status tab
    if ($status) {
      result = result.filter(r => r.status === $status);
    }
    
    // Filter by search query (student name or class)
    if ($search.trim()) {
      const q = $search.toLowerCase();
      result = result.filter(
        r => r.namaSiswa.toLowerCase().includes(q) || r.kelas.toLowerCase().includes(q)
      );
    }
    
    // Sort by updatedAt descending
    return result.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime());
  }
);

export async function loadRooms() {
  listLoading.set(true);
  try {
    const { data, error } = await apiRequest<any>('/v1/admin/aduan');
    if (error) {
      addToast(error || 'Gagal memuat inbox aduan', 'error');
      return;
    }

    const items = (data?.items || []) as Aduan[];
    rooms.set(items);

    // Refresh active room data if one is currently selected
    const active = get(activeRoom);
    if (active) {
      const refreshed = items.find(r => r.id === active.id);
      if (refreshed) {
        activeRoom.set(refreshed);
      }
    }
  } catch (err) {
    console.error('Error loading BK rooms:', err);
    addToast('Gagal memuat inbox aduan', 'error');
  } finally {
    listLoading.set(false);
  }
}

export async function handleSendReply(handlers: { resolve: () => void; reject: () => void }) {
  const active = get(activeRoom);
  const text = get(replyText);

  if (!active || !text.trim()) {
    handlers.reject();
    return;
  }

  try {
    const { error } = await apiRequest('/v1/admin/aduan/respond', {
      method: 'POST',
      body: JSON.stringify({ aduanId: active.id, isi: text }),
    });

    if (error) {
      addToast(error, 'error');
      handlers.reject();
    } else {
      replyText.set('');
      handlers.resolve();
      addToast('Pesan balasan terkirim', 'success');
      await loadRooms();
    }
  } catch (err) {
    console.error('Error sending BK reply:', err);
    addToast('Gagal mengirim balasan', 'error');
    handlers.reject();
  }
}

export async function updateRoomStatus(status: 'open' | 'in_progress' | 'closed') {
  const active = get(activeRoom);
  if (!active) return;

  try {
    const { error } = await apiRequest('/v1/admin/aduan/status', {
      method: 'POST',
      body: JSON.stringify({ aduanId: active.id, status }),
    });

    if (error) {
      addToast(error, 'error');
    } else {
      addToast(`Status aduan diperbarui`, 'success');
      await loadRooms();
    }
  } catch (err) {
    console.error('Error updating status:', err);
    addToast('Gagal memperbarui status', 'error');
  }
}
