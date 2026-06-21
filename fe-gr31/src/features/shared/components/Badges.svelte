<script lang="ts">
  import { User, CheckCircle, Clock, AlertCircle } from 'lucide-svelte';

  let {
    total = 0,
    hadir = 0,
    izinSakit = 0,
    alpa = 0
  }: {
    total?: number;
    hadir?: number;
    izinSakit?: number;
    alpa?: number;
  } = $props();
</script>

<div class="flex items-center gap-4.5 shrink-0 bg-white border border-slate-200/50 rounded-full px-5 py-2.5 shadow-xxs">
  <!-- Total count -->
  <div class="flex items-center gap-1.5 text-slate-500 font-semibold text-xs" title="Total Kehadiran Hari Ini">
    <User class="w-4 h-4 text-slate-400" />
    <span>{total}</span>
  </div>
  <!-- Hadir count -->
  <div class="flex items-center gap-1.5 text-teal-600 font-semibold text-xs" title="Hadir">
    <CheckCircle class="w-4 h-4 text-[#00a294]" />
    <span>{hadir}</span>
  </div>
  <!-- Izin/Sakit count -->
  <div class="flex items-center gap-1.5 text-amber-600 font-semibold text-xs" title="Izin / Sakit">
    <Clock class="w-4 h-4 text-amber-500" />
    <span>{izinSakit}</span>
  </div>
  <!-- Alpa count -->
  <div class="flex items-center gap-1.5 text-rose-600 font-semibold text-xs" title="Alpa (Tanpa Keterangan)">
    <AlertCircle class="w-4 h-4 text-rose-500" />
    <span>{alpa}</span>
  </div>
</div>
