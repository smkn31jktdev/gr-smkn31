<script lang="ts">
  import { User, CheckCircle, Clock, AlertCircle } from 'lucide-svelte';
  import { logs } from '../../../../logic/adminKehadiranSiswaLogic';

  // Calculate counts based on today's logs
  let totalCount = $derived($logs.length);
  let hadirCount = $derived($logs.filter(l => l.status === 'hadir').length);
  let izinSakitCount = $derived($logs.filter(l => l.status === 'izin' || l.status === 'sakit').length);
  let alpaCount = $derived($logs.filter(l => l.status === 'tidak_hadir').length);
</script>

<div class="flex items-center gap-4.5 shrink-0 bg-white border border-slate-200/50 rounded-full px-5 py-2.5 shadow-xxs">
  <!-- Total count -->
  <div class="flex items-center gap-1.5 text-slate-500 font-semibold text-xs" title="Total Siswa">
    <User class="w-4 h-4 text-slate-400" />
    <span>{totalCount}</span>
  </div>
  <!-- Hadir count -->
  <div class="flex items-center gap-1.5 text-teal-600 font-semibold text-xs" title="Masuk">
    <CheckCircle class="w-4 h-4 text-[#00a294]" />
    <span>{hadirCount}</span>
  </div>
  <!-- Izin/Sakit count -->
  <div class="flex items-center gap-1.5 text-amber-600 font-semibold text-xs" title="Izin / Sakit">
    <Clock class="w-4 h-4 text-amber-500" />
    <span>{izinSakitCount}</span>
  </div>
  <!-- Alpa count -->
  <div class="flex items-center gap-1.5 text-rose-600 font-semibold text-xs" title="Tanpa Keterangan">
    <AlertCircle class="w-4 h-4 text-rose-500" />
    <span>{alpaCount}</span>
  </div>
</div>
