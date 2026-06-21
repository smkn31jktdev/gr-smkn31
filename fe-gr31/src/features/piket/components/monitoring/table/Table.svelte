<script lang="ts">
  import { Loader2 } from 'lucide-svelte';
  import type { Kehadiran } from '../../../../student/types/student.types';
  import Magang from '../harian/magang/Magang.svelte';
  import Izin from '../harian/izin/Izin.svelte';
  import Hadir from '../harian/hadir/Hadir.svelte';

  let {
    loading,
    filteredLogs,
    onDelete,
    onOpenPermit
  }: {
    loading: boolean;
    filteredLogs: Kehadiran[];
    onDelete: (id: string) => void;
    onOpenPermit: (log: Kehadiran) => void;
  } = $props();
</script>

<div class="bg-white border border-slate-100/80 rounded-2xl rounded-b-none shadow-xs overflow-hidden">
  <div class="overflow-x-auto">
    <table class="w-full text-left border-collapse text-xs">
      <thead>
        <tr class="bg-slate-50/80 border-b border-slate-100 text-slate-400 font-bold text-[10px] uppercase tracking-wider">
          <th class="p-4 pl-6">Siswa</th>
          <th class="p-4">Kelas Lengkap</th>
          <th class="p-4">Status</th>
          <th class="p-4">Waktu</th>
          <th class="p-4">Keterangan</th>
          <th class="p-4 pr-6 text-center w-14"></th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          <tr>
            <td colspan="6" class="p-8 text-center text-slate-400">
              <div class="flex items-center justify-center gap-2">
                <Loader2 class="w-4 h-4 animate-spin text-slate-400" />
                <span class="font-medium">Memuat data log kehadiran...</span>
              </div>
            </td>
          </tr>
        {:else if filteredLogs.length === 0}
          <tr>
            <td colspan="6" class="p-8 text-center text-slate-400 font-medium">
              Tidak ada log kehadiran ditemukan
            </td>
          </tr>
        {:else}
          {#each filteredLogs as log}
            <tr class="border-b border-slate-50 hover:bg-slate-50/20 transition-colors">
              {#if log.status === 'magang'}
                <Magang 
                  {log} 
                  onDelete={onDelete} 
                  onOpenPermit={onOpenPermit} 
                />
              {:else if log.status === 'izin' || log.status === 'sakit'}
                <Izin 
                  {log} 
                  onDelete={onDelete} 
                  onOpenPermit={onOpenPermit} 
                />
              {:else}
                <Hadir 
                  {log} 
                  onDelete={onDelete} 
                  onOpenPermit={onOpenPermit}
                />
              {/if}
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
