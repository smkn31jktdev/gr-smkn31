<script lang="ts">
  import { 
    Trash2, 
    Loader2, 
    ChevronLeft, 
    ChevronRight 
  } from 'lucide-svelte';
  import DropdownChoice from '../../../../shared/components/DropdownChoice.svelte';
  import {
    logs,
    loading,
    page,
    limit,
    total,
    hasMore,
    handleFilter,
    handleDelete,
    loadData,
    openCreateForStudent
  } from '../../../logic/adminKehadiranSiswaLogic';

  // Indonesian date formatter
  function formatIndonesianDateStr(dateStr: string): string {
    if (!dateStr) return '';
    const parts = dateStr.split('-');
    if (parts.length !== 3) return dateStr;
    const d = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]));
    const days = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
    const months = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    return `${days[d.getDay()]}, ${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`;
  }

  function getStatusBadgeClass(status: string) {
    switch(status) {
      case 'hadir': return 'bg-teal-50 text-[#00a294] border-teal-100/30';
      case 'izin': return 'bg-slate-50 text-slate-650 border-slate-200/50';
      case 'sakit': return 'bg-amber-50/80 text-amber-700 border-amber-100/40';
      case 'magang': return 'bg-indigo-50/80 text-indigo-750 border-indigo-100/40';
      case 'belum': return 'bg-slate-100 text-slate-400 border-slate-200/50';
      default: return 'bg-rose-50 text-rose-700 border-rose-100/40';
    }
  }
</script>

<!-- Logs Table -->
<div class="bg-white border border-slate-100/80 rounded-2xl rounded-b-none shadow-xs overflow-hidden">
  <div class="overflow-x-auto">
    <table class="w-full text-left border-collapse text-xs">
      <thead>
        <tr class="bg-slate-50/80 border-b border-slate-100 text-slate-400 font-bold text-[10px] uppercase tracking-wider">
          <th class="p-4 pl-6">Siswa</th>
          <th class="p-4">Kelas</th>
          <th class="p-4">Status</th>
          <th class="p-4">Tanggal</th>
          <th class="p-4">Keterangan</th>
          <th class="p-4 pr-6 text-center w-14"></th>
        </tr>
      </thead>
      <tbody>
        {#if $loading}
          <tr>
            <td colspan="6" class="p-8 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <Loader2 class="w-4 h-4 animate-spin text-slate-400" />
                <span class="font-medium">Memuat data log kehadiran...</span>
              </div>
            </td>
          </tr>
        {:else if $logs.length === 0}
          <tr>
            <td colspan="6" class="p-8 text-center text-slate-400 font-medium">
              Tidak ada log kehadiran ditemukan
            </td>
          </tr>
        {:else}
          {#each $logs as log}
            {@const initial = log.namaSiswa ? log.namaSiswa.charAt(0).toUpperCase() : 'S'}
            <tr 
              onclick={() => openCreateForStudent(log.nis, (log.status as any) !== 'belum' ? log.status : 'hadir', log.fotoIzin || '', log.alasan || '')}
              class="border-b border-slate-50 hover:bg-slate-50/20 transition-colors cursor-pointer"
            >
              <!-- Siswa -->
              <td class="p-4 pl-6">
                <div class="flex items-center gap-3">
                  <div class="w-8.5 h-8.5 rounded-xl bg-slate-50 text-slate-655 flex items-center justify-center font-bold text-xs border border-slate-100 shrink-0">
                    {initial}
                  </div>
                  <div class="text-left">
                    <span class="font-bold text-slate-800 uppercase tracking-wide block leading-tight">{log.namaSiswa}</span>
                    <span class="text-[10px] font-mono font-medium text-slate-400 block mt-0.5">{log.nis}</span>
                  </div>
                </div>
              </td>

              <!-- Kelas -->
              <td class="p-4">
                <span class="font-semibold text-slate-600">{log.kelas || '-'}</span>
              </td>

              <!-- Status -->
              <td class="p-4">
                <span class="inline-block px-2.5 py-0.5 rounded-md border text-[9px] font-bold uppercase tracking-wider {getStatusBadgeClass(log.status)}">
                  {log.status === 'tidak_hadir' ? 'tanpa keterangan' : log.status}
                </span>
              </td>

              <!-- Tanggal -->
              <td class="p-4">
                <span class="font-semibold text-slate-550">
                  {formatIndonesianDateStr(log.tanggal)}
                </span>
              </td>

              <!-- Keterangan -->
              <td class="p-4 text-slate-550 font-medium max-w-[200px] truncate">
                {log.alasan || '-'}
              </td>

              <!-- Delete Action -->
              <td class="p-4 pr-6 text-center">
                <button 
                  onclick={(e) => { e.stopPropagation(); handleDelete(log.id); }} 
                  class="inline-flex items-center justify-center p-1.5 text-slate-400 hover:text-rose-600 rounded-lg hover:bg-rose-50 transition-colors border-none bg-transparent cursor-pointer"
                  title="Hapus data absensi"
                >
                  <Trash2 class="w-4 h-4" />
                </button>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>

<!-- Pagination Controls -->
<div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white border border-slate-100/80 border-t-0 rounded-b-2xl px-6 py-4 shadow-xs">
  <div class="flex items-center gap-2 text-slate-500 text-xs">
    <span>Tampilkan</span>
    <div class="w-24 text-left">
      <DropdownChoice
        options={[
          { value: 10, label: '10' },
          { value: 50, label: '50' },
          { value: 100, label: '100' }
        ]}
        bind:value={$limit}
        onchange={handleFilter}
        placeholder="Tampilkan"
      />
    </div>
    <span>data per halaman (Total: {$total} records)</span>
  </div>

  <div class="flex items-center gap-2">
    <button 
      onclick={() => { if ($page > 1) { page.update(p => p - 1); loadData(); } }}
      disabled={$page === 1 || $loading}
      class="p-2 border border-slate-200 rounded-xl bg-white hover:bg-slate-50 shadow-xxs transition-all cursor-pointer disabled:opacity-50 flex items-center justify-center"
      aria-label="Sebelumnya"
    >
      <ChevronLeft class="w-4 h-4 {$page === 1 || $loading ? 'text-slate-300' : 'text-black'}" />
    </button>
    <span class="text-xs font-bold text-slate-700 px-2">Halaman {$page}</span>
    <button 
      onclick={() => { if ($hasMore) { page.update(p => p + 1); loadData(); } }}
      disabled={!$hasMore || $loading}
      class="p-2 border border-slate-200 rounded-xl bg-white hover:bg-slate-50 shadow-xxs transition-all cursor-pointer disabled:opacity-50 flex items-center justify-center"
      aria-label="Selanjutnya"
    >
      <ChevronRight class="w-4 h-4 {!$hasMore || $loading ? 'text-slate-300' : 'text-black'}" />
    </button>
  </div>
</div>
