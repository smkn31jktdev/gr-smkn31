<script lang="ts">
  import { Calendar, Clock } from 'lucide-svelte';
  import { formatIndonesianDateStr } from './absensi';

  let { log } = $props<{ log: any }>();
</script>

<div class="grid grid-cols-2 gap-3">
  <div class="bg-white border border-slate-100 rounded-xl p-3 flex items-center gap-2.5 text-left shadow-xxs">
    <Calendar class="w-4 h-4 text-slate-400 shrink-0" />
    <div class="min-w-0">
      <span class="text-[9px] font-bold text-slate-400 block uppercase">Tanggal</span>
      <span class="text-[10px] font-extrabold text-slate-750 block truncate mt-0.5">
        {formatIndonesianDateStr(log.tanggal)}
      </span>
    </div>
  </div>
  <div class="bg-white border border-slate-100 rounded-xl p-3 flex items-center gap-2.5 text-left shadow-xxs">
    <Clock class="w-4 h-4 text-slate-400 shrink-0" />
    <div class="min-w-0">
      <span class="text-[9px] font-bold text-slate-400 block uppercase">Waktu Absen</span>
      <span class="text-[10px] font-extrabold text-slate-750 block mt-0.5">
        {log.waktuAbsen ? `${log.waktuAbsen.substring(0, 5)} WIB` : '-'}
      </span>
    </div>
  </div>
</div>
