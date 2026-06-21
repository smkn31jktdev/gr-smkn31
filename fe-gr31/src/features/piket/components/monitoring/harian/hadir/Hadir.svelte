<script lang="ts">
  import { Trash2 } from 'lucide-svelte';
  import type { Kehadiran } from '../../../../../student/types/student.types';
  import Status from '../../../../../shared/components/Status.svelte';

  let { log, onDelete, onOpenPermit }: {
    log: Kehadiran;
    onDelete: (id: string) => void;
    onOpenPermit: (log: Kehadiran) => void;
  } = $props();

  const initial = $derived(log.namaSiswa ? log.namaSiswa.charAt(0).toUpperCase() : 'S');
</script>

<!-- Siswa (Avatar + Name & NISN) -->
<td class="p-4 pl-6">
  <button 
    onclick={() => onOpenPermit(log)}
    class="flex items-center gap-3 border-none bg-transparent p-0 cursor-pointer text-left focus:outline-none group/name"
    title="Klik untuk melihat detail absensi/perangkat"
  >
    <div class="w-8.5 h-8.5 rounded-xl bg-slate-50 text-slate-655 flex items-center justify-center font-bold text-xs border border-slate-100 shrink-0 group-hover/name:bg-slate-100 transition-colors">
      {initial}
    </div>
    <div class="text-left">
      <span class="font-bold text-slate-800 group-hover/name:text-[#00a294] uppercase tracking-wide block leading-tight transition-colors">{log.namaSiswa}</span>
      <span class="text-[10px] font-mono font-medium text-slate-400 block mt-0.5">{log.nis}</span>
      <span class="text-[9px] font-bold text-teal-600 group-hover/name:text-teal-700 block mt-0.5 transition-colors">Klik untuk detail</span>
    </div>
  </button>
</td>

<!-- Kelas Lengkap -->
<td class="p-4">
  <span class="font-semibold text-slate-600">{log.kelas || '-'}</span>
</td>

<!-- Status -->
<td class="p-4">
  <Status status={log.status} waktuAbsen={log.waktuAbsen} />
</td>

<!-- Waktu -->
<td class="p-4">
  <span class="font-semibold font-mono {log.status === 'hadir' && log.waktuAbsen && log.waktuAbsen.substring(0, 5) > '06:35' ? 'text-rose-600' : 'text-slate-550'}">
    {log.waktuAbsen ? log.waktuAbsen.substring(0, 5) : '-'}
  </span>
</td>

<!-- Keterangan -->
<td class="p-4 text-slate-550 font-medium max-w-[200px] truncate">
  {log.alasan || '-'}
</td>

<!-- Subtle Action Column -->
<td class="p-4 pr-6 text-center">
  <button 
    onclick={() => onDelete(log.id)} 
    class="inline-flex items-center justify-center p-1.5 text-slate-400 hover:text-rose-600 rounded-lg hover:bg-rose-50 transition-colors border-none bg-transparent cursor-pointer"
    title="Hapus data absensi"
  >
    <Trash2 class="w-4 h-4" />
  </button>
</td>

