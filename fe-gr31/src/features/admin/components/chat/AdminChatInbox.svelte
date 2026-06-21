<script lang="ts">
  import { onMount } from 'svelte';
  import { apiRequest, BASE } from '../../../../api/client';
  import { addToast } from '../../../../stores/uiStore';
  import SubmitButton from '../../../shared/components/SubmitButton.svelte';
  import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';
  import type { Aduan } from '../../../student/types/student.types';
  import { RotateCw, Inbox, MessageSquare, Download, Printer } from 'lucide-svelte';

  let rooms = $state<Aduan[]>([]);
  let activeRoom = $state<Aduan | null>(null);
  let replyText = $state('');
  let listLoading = $state(false);
  let selectedStatus = $state('');

  async function handleDownloadCSV() {
    try {
      const token = localStorage.getItem('adminToken') || localStorage.getItem('studentToken');
      const statusVal = selectedStatus;
      const url = `${BASE}/v1/admin/aduan/export/csv` + (statusVal ? `?status=${statusVal}` : '');
      const response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      if (!response.ok) {
        throw new Error('Gagal mengunduh CSV');
      }
      const blob = await response.blob();
      const blobUrl = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = blobUrl;
      a.download = `arsip-aduan-siswa${statusVal ? '-' + statusVal : ''}.csv`;
      document.body.appendChild(a);
      a.click();
      a.remove();
      window.URL.revokeObjectURL(blobUrl);
    } catch (err: any) {
      addToast(err.message || 'Terjadi kesalahan saat mengunduh CSV', 'error');
    }
  }

  async function handlePrintRoom() {
    if (!activeRoom) return;
    try {
      const token = localStorage.getItem('adminToken') || localStorage.getItem('studentToken');
      const url = `${BASE}/v1/admin/aduan/export/html?id=${activeRoom.id}`;
      const response = await fetch(url, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });
      if (!response.ok) {
        throw new Error('Gagal mencetak aduan');
      }
      const html = await response.text();
      const blob = new Blob([html], { type: 'text/html' });
      const blobUrl = window.URL.createObjectURL(blob);
      window.open(blobUrl, '_blank');
    } catch (err: any) {
      addToast(err.message || 'Terjadi kesalahan saat memuat cetak aduan', 'error');
    }
  }

  async function loadRooms() {
    listLoading = true;
    const q = new URLSearchParams();
    if (selectedStatus) q.append('status', selectedStatus);
    
    const { data, error } = await apiRequest<any>(`/v1/admin/aduan?${q.toString()}`);
    if (error) {
      addToast('Gagal memuat inbox konsultasi', 'error');
    } else {
      rooms = data.items || [];
      if (activeRoom) {
        const refreshed = rooms.find(r => r.id === activeRoom?.id);
        if (refreshed) activeRoom = refreshed;
      }
    }
    listLoading = false;
  }

  onMount(() => {
    loadRooms();
  });

  async function handleSendReply(handlers: { resolve: () => void; reject: () => void }) {
    if (!activeRoom || !replyText.trim()) {
      handlers.reject();
      return;
    }

    const { error } = await apiRequest('/v1/admin/aduan/respond', {
      method: 'POST',
      body: JSON.stringify({ aduanId: activeRoom.id, isi: replyText }),
    });

    if (error) {
      addToast(error, 'error');
      handlers.reject();
    } else {
      replyText = '';
      handlers.resolve();
      addToast('Pesan balasan terkirim', 'success');
      await loadRooms();
    }
  }

  async function updateStatus(status: 'open' | 'in_progress' | 'closed') {
    if (!activeRoom) return;
    
    const { error } = await apiRequest('/v1/admin/aduan/status', {
      method: 'POST',
      body: JSON.stringify({ aduanId: activeRoom.id, status }),
    });

    if (error) {
      addToast(error, 'error');
    } else {
      addToast(`Status aduan diperbarui menjadi ${status}`, 'success');
      await loadRooms();
    }
  }

  function formatDate(isoStr: string) {
    if (!isoStr) return '';
    const date = new Date(isoStr);
    return date.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<div class="card p-0 grid grid-cols-1 md:grid-cols-3 min-h-[550px] max-h-[650px] overflow-hidden">
  
  <!-- Left Side: Chat Inbox list -->
  <div class="border-r border-border flex flex-col h-full overflow-hidden">
    <div class="p-4 border-b border-border flex items-center justify-between bg-gray-50/50 gap-2">
      <div class="w-full max-w-[140px] text-left">
        <DropdownChoice
          options={[
            { value: '', label: 'Semua Status' },
            { value: 'open', label: 'Open' },
            { value: 'in_progress', label: 'In Progress' },
            { value: 'closed', label: 'Closed' }
          ]}
          bind:value={selectedStatus}
          onchange={loadRooms}
          placeholder="Semua Status"
        />
      </div>
      
      <div class="flex items-center gap-1.5 shrink-0">
        <button 
          onclick={handleDownloadCSV} 
          class="p-1.5 border rounded-lg hover:bg-gray-100 text-emerald-600 border-emerald-200 hover:border-emerald-300 transition-colors text-xs flex items-center gap-1 cursor-pointer font-bold"
        >
          <Download class="w-3.5 h-3.5" /> Unduh CSV
        </button>

        <button 
          onclick={loadRooms} 
          class="p-1.5 border rounded-lg hover:bg-gray-100 text-muted transition-colors text-xs flex items-center gap-1 cursor-pointer font-semibold"
        >
          <RotateCw class="w-3.5 h-3.5" /> Segarkan
        </button>
      </div>
    </div>

    <!-- Rooms Inbox list -->
    <div class="flex-1 overflow-y-auto p-3 space-y-2">
      {#if listLoading && rooms.length === 0}
        <p class="text-xs text-muted text-center py-8">Memuat chat inbox...</p>
      {:else if rooms.length === 0}
        <div class="text-center py-12 text-muted px-4">
          <Inbox class="w-10 h-10 text-slate-300 mx-auto mb-2" />
          <p class="text-xs font-semibold mt-2">Inbox kosong.</p>
        </div>
      {:else}
        {#each rooms as room}
          <button
            onclick={() => activeRoom = room}
            class="w-full text-left p-3.5 rounded-xl border transition-all duration-150 flex flex-col gap-1 cursor-pointer {activeRoom?.id === room.id ? 'border-primary bg-primary-light/30' : 'border-border bg-white hover:bg-gray-50'}"
          >
            <div class="flex items-center justify-between">
              <span class="text-xs font-bold text-foreground truncate max-w-[150px]">{room.namaSiswa} ({room.kelas})</span>
              <span 
                class="inline-block px-1.5 py-0.5 rounded-full text-[10px] font-bold uppercase"
                class:bg-amber-100={room.status === 'open'}
                class:text-amber-800={room.status === 'open'}
                class:bg-blue-100={room.status === 'in_progress'}
                class:text-blue-800={room.status === 'in_progress'}
                class:bg-gray-100={room.status === 'closed'}
                class:text-gray-800={room.status === 'closed'}
              >
                {room.status}
              </span>
            </div>
            
            <p class="text-xs text-muted truncate leading-relaxed">
              {room.messages[room.messages.length - 1]?.isi || 'Tidak ada pesan'}
            </p>
            
            <span class="text-[10px] text-muted block text-right mt-1 font-medium">
              {formatDate(room.updatedAt)}
            </span>
          </button>
        {/each}
      {/if}
    </div>
  </div>

  <!-- Right Side: Message Thread -->
  <div class="col-span-2 flex flex-col h-full overflow-hidden bg-gray-50/30">
    {#if activeRoom}
      <!-- Header -->
      <div class="p-4 border-b border-border bg-surface flex items-center justify-between flex-wrap gap-2">
        <div>
          <h4 class="text-sm font-bold text-foreground">
            {activeRoom.namaSiswa} ({activeRoom.kelas})
          </h4>
          <p class="text-[10px] text-muted font-medium mt-0.5">Tiket ID: {activeRoom.id}</p>
        </div>
        
        <!-- Status actions -->
        <div class="flex items-center gap-1.5">
          {#if activeRoom.status !== 'closed'}
            {#if activeRoom.status === 'open'}
              <button onclick={() => updateStatus('in_progress')} class="px-2.5 py-1 text-xxs font-bold bg-blue-50 text-blue-600 hover:bg-blue-100 border border-blue-200 rounded-lg transition-colors cursor-pointer">
                Proses
              </button>
            {/if}
            <button onclick={() => updateStatus('closed')} class="px-2.5 py-1 text-xxs font-bold bg-rose-50 text-rose-600 hover:bg-rose-100 border border-rose-200 rounded-lg transition-colors cursor-pointer">
              Selesai & Tutup
            </button>
          {:else}
            <button onclick={() => updateStatus('open')} class="px-2.5 py-1 text-xxs font-bold bg-amber-50 text-amber-600 hover:bg-amber-100 border border-amber-200 rounded-lg transition-colors cursor-pointer">
              Buka Kembali
            </button>
          {/if}

          <button 
            onclick={handlePrintRoom}
            class="px-2.5 py-1 text-xxs font-bold bg-emerald-50 text-emerald-600 hover:bg-emerald-100 border border-emerald-200 rounded-lg transition-colors cursor-pointer flex items-center gap-1 uppercase"
          >
            <Printer class="w-3 h-3" /> Cetak Arsip
          </button>
        </div>
      </div>

      <!-- Messages Thread -->
      <div class="flex-1 overflow-y-auto p-4 space-y-3.5">
        {#each activeRoom.messages as msg}
          <div class="flex flex-col max-w-[75%]" class:self-end={msg.role === 'admin'} class:ml-auto={msg.role === 'admin'}>
            <span class="text-xxs text-muted font-bold mb-1 ml-1" class:text-right={msg.role === 'admin'}>
              {msg.role === 'admin' ? `Staf/Guru (${activeRoom.adminNama || 'Konselor'})` : activeRoom.namaSiswa}
            </span>
            <div 
              class="p-3.5 rounded-2xl text-sm leading-relaxed"
              class:bg-blue-600={msg.role === 'admin'}
              class:text-white={msg.role === 'admin'}
              class:rounded-tr-none={msg.role === 'admin'}
              class:bg-surface={msg.role !== 'admin'}
              class:text-foreground={msg.role !== 'admin'}
              class:border={msg.role !== 'admin'}
              class:border-border={msg.role !== 'admin'}
              class:rounded-tl-none={msg.role !== 'admin'}
            >
              {msg.isi}
            </div>
            <span class="text-xxs text-muted mt-1 ml-1" class:text-right={msg.role === 'admin'}>
              {formatDate(msg.timestamp)}
            </span>
          </div>
        {/each}
      </div>

      <!-- Input Area (if not closed) -->
      {#if activeRoom.status !== 'closed'}
        <div class="p-4 border-t border-border bg-surface">
          <form class="flex gap-2" onsubmit={(e) => e.preventDefault()}>
            <input
              type="text"
              placeholder="Tulis balasan konsultasi di sini..."
              bind:value={replyText}
              class="input py-2 flex-1"
            />
            <SubmitButton
              label="Balas"
              loadingLabel="..."
              className="px-6"
              onclick={handleSendReply}
            />
          </form>
        </div>
      {/if}
    {:else}
      <!-- Empty state -->
      <div class="flex-1 flex flex-col items-center justify-center p-6 text-center">
        <MessageSquare class="w-12 h-12 text-slate-300 mx-auto mb-4" />
        <h4 class="text-base font-bold text-foreground mt-3">Detail Percakapan</h4>
        <p class="text-xs text-muted mt-1 max-w-[280px] mx-auto">Pilih salah satu tiket aduan di sebelah kiri untuk melihat pesan dan memberikan bimbingan/balasan.</p>
      </div>
    {/if}
  </div>
</div>
