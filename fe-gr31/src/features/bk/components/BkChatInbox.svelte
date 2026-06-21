<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    MessageSquare, 
    AlertCircle, 
    Clock, 
    CheckCircle, 
    ChevronRight, 
    User, 
    ArrowLeft, 
    RefreshCw, 
    Search,
    Loader2,
    Inbox,
    Download,
    Printer
  } from 'lucide-svelte';
  import SubmitButton from '../../shared/components/SubmitButton.svelte';
  import { BASE } from '../../../api/client';
  import {
    rooms,
    activeRoom,
    listLoading,
    replyText,
    selectedStatus,
    searchQuery,
    totalAduanCount,
    activeAduanCount,
    filteredRooms,
    loadRooms,
    handleSendReply,
    updateRoomStatus
  } from '../logic/bkChatLogic';

  onMount(() => {
    loadRooms();
  });

  async function handleDownloadCSV() {
    try {
      const token = localStorage.getItem('adminToken') || localStorage.getItem('studentToken');
      const statusVal = $selectedStatus;
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
      alert(err.message || 'Terjadi kesalahan saat mengunduh CSV');
    }
  }

  async function handlePrintRoom() {
    if (!$activeRoom) return;
    try {
      const token = localStorage.getItem('adminToken') || localStorage.getItem('studentToken');
      const url = `${BASE}/v1/admin/aduan/export/html?id=${$activeRoom.id}`;
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
      alert(err.message || 'Terjadi kesalahan saat memuat cetak aduan');
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

  function handleTabChange(status: string) {
    selectedStatus.set(status);
  }
</script>

<div class="h-[calc(100vh-8.5rem)] flex flex-col font-sans select-none">
  
  <!-- Header Bar -->
  <div class="flex items-center justify-between bg-white border border-slate-100 rounded-2xl p-4 mb-5 shadow-xs shrink-0">
    <div class="flex items-center">
      <a href="/bk" class="p-2 border border-slate-100 hover:border-slate-200 rounded-xl hover:bg-slate-50 text-slate-400 hover:text-slate-600 transition-all mr-3 flex items-center justify-center bg-white cursor-pointer">
        <ArrowLeft class="w-4 h-4" />
      </a>
      <div class="w-10 h-10 rounded-xl bg-slate-50 border border-slate-100 text-[#00a294] flex items-center justify-center mr-3 shrink-0">
        <MessageSquare class="w-5 h-5" />
      </div>
      <div class="text-left">
        <h2 class="text-xs font-black text-slate-800 uppercase tracking-tight leading-none">Layanan Aduan Siswa</h2>
        <span class="text-[9px] font-extrabold text-slate-400 uppercase tracking-wider mt-1 block">Bimbingan & Konseling</span>
      </div>
    </div>
    
    <div class="flex items-center gap-2">
      <button 
        onclick={handleDownloadCSV}
        disabled={$listLoading}
        class="flex items-center gap-1.5 px-3 py-1.5 text-[10px] font-bold text-[#00a294] hover:text-white border border-[#00a294]/20 hover:border-[#00a294] rounded-xl bg-[#00a294]/5 hover:bg-[#00a294] transition-all cursor-pointer shadow-xxs"
      >
        <Download class="w-3.5 h-3.5" />
        Unduh CSV
      </button>

      <button 
        onclick={loadRooms} 
        disabled={$listLoading}
        class="flex items-center gap-1.5 px-3 py-1.5 text-[10px] font-bold text-slate-600 hover:text-slate-900 border border-slate-150 rounded-xl bg-slate-50/50 hover:bg-slate-50 transition-colors cursor-pointer"
      >
        <RefreshCw class="w-3.5 h-3.5 {$listLoading ? 'animate-spin' : ''}" />
        Segarkan
      </button>
    </div>
  </div>

  <!-- Main Container -->
  <div class="flex-1 grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-0 overflow-hidden">
    
    <!-- Left Column: Search, Stats, Tabs & List (Col span 4) -->
    <div class="lg:col-span-4 bg-white border border-slate-100 rounded-2xl p-4 shadow-xs flex flex-col min-h-0 overflow-hidden">
      
      <!-- Top Overview Card -->
      <div class="bg-slate-50/50 border border-slate-100/60 rounded-xl p-3.5 mb-4 text-left">
        <div class="flex items-center gap-2">
          <div class="w-6 h-6 rounded-lg bg-[#00a294]/10 text-[#00a294] flex items-center justify-center shrink-0">
            <User class="w-3.5 h-3.5" />
          </div>
          <span class="text-[10px] font-extrabold text-[#00a294] uppercase tracking-wider">Info Konseling</span>
        </div>
        <p class="text-[10px] text-slate-400 leading-relaxed font-medium mt-1.5">
          Pantau dan tindaklanjuti laporan atau aduan siswa dengan sigap. Kerahasiaan diutamakan.
        </p>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-2 gap-3 mb-4 shrink-0">
        <div class="border border-slate-100 rounded-xl p-3 text-center bg-white shadow-xxs">
          <h4 class="text-sm font-extrabold text-slate-700 leading-none">{$activeAduanCount}</h4>
          <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block mt-1">Laporan Aktif</span>
        </div>
        <div class="border border-slate-100 rounded-xl p-3 text-center bg-white shadow-xxs">
          <h4 class="text-sm font-extrabold text-slate-700 leading-none">{$totalAduanCount}</h4>
          <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block mt-1">Total Aduan</span>
        </div>
      </div>

      <!-- Search Input -->
      <div class="relative mb-3.5 shrink-0">
        <Search class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input 
          type="text" 
          placeholder="Cari siswa atau kelas..."
          value={$searchQuery}
          oninput={(e) => searchQuery.set((e.target as HTMLInputElement).value)}
          class="w-full bg-slate-50/50 border border-slate-100 rounded-xl pl-9 pr-3 py-2 text-xs text-slate-700 placeholder-slate-400 outline-none focus:bg-white focus:border-slate-200 focus:ring-2 focus:ring-slate-150/50 transition-all font-sans"
        />
      </div>

      <!-- Navigation Tabs -->
      <div class="flex border border-slate-100 bg-slate-50/50 p-0.5 rounded-lg mb-3 shrink-0">
        <button 
          onclick={() => handleTabChange('')}
          class="flex-1 py-1 rounded-md text-[10px] font-bold transition-all border-none cursor-pointer {$selectedStatus === '' ? 'bg-white text-slate-800 shadow-xxs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
        >
          Semua
        </button>
        <button 
          onclick={() => handleTabChange('open')}
          class="flex-1 py-1 rounded-md text-[10px] font-bold transition-all border-none cursor-pointer {$selectedStatus === 'open' ? 'bg-white text-slate-800 shadow-xxs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
        >
          Baru
        </button>
        <button 
          onclick={() => handleTabChange('in_progress')}
          class="flex-1 py-1 rounded-md text-[10px] font-bold transition-all border-none cursor-pointer {$selectedStatus === 'in_progress' ? 'bg-white text-slate-800 shadow-xxs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
        >
          Proses
        </button>
        <button 
          onclick={() => handleTabChange('closed')}
          class="flex-1 py-1 rounded-md text-[10px] font-bold transition-all border-none cursor-pointer {$selectedStatus === 'closed' ? 'bg-white text-slate-800 shadow-xxs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
        >
          Selesai
        </button>
      </div>

      <!-- Scrollable Rooms List -->
      <div class="flex-1 overflow-y-auto pr-1 space-y-2 custom-scrollbar">
        {#if $listLoading && $rooms.length === 0}
          <div class="flex flex-col items-center justify-center py-12 text-slate-400">
            <Loader2 class="w-5 h-5 animate-spin text-slate-400 mb-2" />
            <p class="text-[10px] font-semibold">Memuat aduan...</p>
          </div>
        {:else if $filteredRooms.length === 0}
          <div class="flex flex-col items-center justify-center py-16 text-slate-400 text-center">
            <Inbox class="w-8 h-8 text-slate-350 mb-2" />
            <p class="text-[10px] font-bold">Tidak ada aduan ditemukan</p>
          </div>
        {:else}
          {#each $filteredRooms as room}
            <button
              onclick={() => activeRoom.set(room)}
              class="w-full text-left p-3 rounded-xl border transition-all duration-150 flex flex-col gap-1 cursor-pointer {$activeRoom?.id === room.id ? 'border-[#00a294] bg-[#00a294]/5' : 'border-slate-100 bg-white hover:bg-slate-50/55'}"
            >
              <div class="flex items-center justify-between">
                <span class="text-xs font-bold text-slate-700 truncate max-w-[140px] uppercase">{room.namaSiswa}</span>
                <span 
                  class="w-2 h-2 rounded-full shrink-0"
                  class:bg-amber-400={room.status === 'open'}
                  class:bg-blue-400={room.status === 'in_progress'}
                  class:bg-slate-300={room.status === 'closed'}
                  title={room.status}
                ></span>
              </div>
              
              <div class="flex justify-between items-center text-[10px]">
                <span class="text-slate-400 font-semibold">{room.kelas}</span>
                <span class="text-slate-400 font-medium">{formatDate(room.updatedAt)}</span>
              </div>
              
              <p class="text-[11px] text-slate-500 truncate mt-1 leading-normal">
                {room.messages[room.messages.length - 1]?.isi || 'Tidak ada pesan'}
              </p>
            </button>
          {/each}
        {/if}
      </div>

    </div>

    <!-- Right Column: Message thread detail (Col span 8) -->
    <div class="lg:col-span-8 bg-white border border-slate-100 rounded-2xl shadow-xs flex flex-col min-h-0 overflow-hidden font-sans">
      
      {#if $activeRoom}
        <!-- Header Detail -->
        <div class="p-4 border-b border-slate-100 bg-slate-50/20 flex items-center justify-between flex-wrap gap-3 shrink-0">
          <div class="text-left">
            <h4 class="text-xs font-extrabold text-slate-700 uppercase tracking-tight">
              {$activeRoom.namaSiswa} <span class="text-slate-450 mx-1">•</span> {$activeRoom.kelas}
            </h4>
            <p class="text-[9px] text-slate-400 font-medium mt-0.5">ID Tiket: {$activeRoom.id}</p>
          </div>
          
          <div class="flex items-center gap-2">
            <!-- Status Badge -->
            <span 
              class="px-2 py-0.5 rounded-md text-[9px] font-bold uppercase border"
              class:bg-amber-50={$activeRoom.status === 'open'}
              class:text-amber-700={$activeRoom.status === 'open'}
              class:border-amber-100={$activeRoom.status === 'open'}
              class:bg-blue-50={$activeRoom.status === 'in_progress'}
              class:text-blue-700={$activeRoom.status === 'in_progress'}
              class:border-blue-100={$activeRoom.status === 'in_progress'}
              class:bg-slate-50={$activeRoom.status === 'closed'}
              class:text-slate-650={$activeRoom.status === 'closed'}
              class:border-slate-150={$activeRoom.status === 'closed'}
            >
              {$activeRoom.status === 'open' ? 'Baru' : $activeRoom.status === 'in_progress' ? 'Diproses' : 'Selesai'}
            </span>

            <!-- Status Actions -->
            {#if $activeRoom.status !== 'closed'}
              {#if $activeRoom.status === 'open'}
                <button 
                  onclick={() => updateRoomStatus('in_progress')} 
                  class="px-2.5 py-1 text-[9px] font-extrabold uppercase bg-blue-50 hover:bg-blue-100 text-blue-600 border border-blue-150 rounded-lg transition-colors cursor-pointer"
                >
                  Proses
                </button>
              {/if}
              <button 
                onclick={() => updateRoomStatus('closed')} 
                class="px-2.5 py-1 text-[9px] font-extrabold uppercase bg-rose-50 hover:bg-rose-100 text-rose-600 border border-rose-150 rounded-lg transition-colors cursor-pointer"
              >
                Selesai
              </button>
            {:else}
              <button 
                onclick={() => updateRoomStatus('open')} 
                class="px-2.5 py-1 text-[9px] font-extrabold uppercase bg-amber-50 hover:bg-amber-100 text-amber-600 border border-amber-150 rounded-lg transition-colors cursor-pointer"
              >
                Buka Kembali
              </button>
            {/if}
            
            <button 
              onclick={handlePrintRoom}
              class="flex items-center gap-1 px-2.5 py-1 text-[9px] font-extrabold uppercase bg-emerald-50 hover:bg-emerald-100 text-emerald-650 border border-emerald-150 rounded-lg transition-colors cursor-pointer"
            >
              <Printer class="w-3 h-3" />
              Cetak Arsip
            </button>
          </div>
        </div>

        <!-- Thread Messages List -->
        <div class="flex-1 overflow-y-auto p-4 space-y-4 bg-slate-50/10 custom-scrollbar flex flex-col min-h-0">
          {#each $activeRoom.messages as msg}
            <div class="flex flex-col max-w-[70%]" class:self-end={msg.role === 'admin'} class:ml-auto={msg.role === 'admin'}>
              <span class="text-[9px] text-slate-400 font-bold mb-1 ml-1" class:text-right={msg.role === 'admin'}>
                {msg.role === 'admin' ? `Konselor (${$activeRoom.adminNama || 'Guru BK'})` : $activeRoom.namaSiswa}
              </span>
              <div 
                class="p-3 rounded-2xl text-xs leading-relaxed font-medium"
                class:bg-[#00a294]={msg.role === 'admin'}
                class:text-white={msg.role === 'admin'}
                class:rounded-tr-none={msg.role === 'admin'}
                class:bg-white={msg.role !== 'admin'}
                class:text-slate-700={msg.role !== 'admin'}
                class:border={msg.role !== 'admin'}
                class:border-slate-100={msg.role !== 'admin'}
                class:rounded-tl-none={msg.role !== 'admin'}
                class:shadow-xxs={msg.role !== 'admin'}
              >
                {msg.isi}
              </div>
              <span class="text-[9px] text-slate-400 mt-1 ml-1" class:text-right={msg.role === 'admin'}>
                {formatDate(msg.timestamp)}
              </span>
            </div>
          {/each}
        </div>

        <!-- Reply Input Bar -->
        {#if $activeRoom.status !== 'closed'}
          <div class="p-4 border-t border-slate-100 bg-white shrink-0">
            <form class="flex gap-2.5" onsubmit={(e) => e.preventDefault()}>
              <input
                type="text"
                placeholder="Tulis balasan konsultasi di sini..."
                value={$replyText}
                oninput={(e) => replyText.set((e.target as HTMLInputElement).value)}
                class="w-full bg-slate-50/50 border border-slate-100 rounded-xl px-4 py-2.5 text-xs text-slate-700 placeholder-slate-450 outline-none focus:bg-white focus:border-slate-200 focus:ring-2 focus:ring-slate-150/50 transition-all font-sans font-medium"
              />
              <SubmitButton
                label="Balas"
                loadingLabel="..."
                className="px-5 py-2.5 bg-[#00a294] hover:bg-[#008f83] rounded-xl text-white font-bold text-xs transition-all border-none cursor-pointer shrink-0"
                onclick={handleSendReply}
              />
            </form>
          </div>
        {/if}
      {:else}
        <!-- Thread Empty state -->
        <div class="flex-1 flex flex-col items-center justify-center p-6 text-center text-slate-400">
          <div class="w-14 h-14 rounded-2xl bg-slate-50 border border-slate-100/60 flex items-center justify-center mb-3">
            <Inbox class="w-6 h-6 text-slate-350" />
          </div>
          <h4 class="text-xs font-black text-slate-650 uppercase tracking-wider">Detail Percakapan</h4>
          <p class="text-[10px] text-slate-400 mt-1 max-w-[280px] leading-relaxed font-medium">Pilih salah satu tiket aduan di sebelah kiri untuk melihat pesan dan memberikan bimbingan/balasan.</p>
        </div>
      {/if}
    </div>

  </div>
</div>

<style>
  /* Custom scrollbar styling for a clean sleek feel */
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 99px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
  }
</style>
