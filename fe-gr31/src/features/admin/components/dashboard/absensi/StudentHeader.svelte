<script lang="ts">
  let { log } = $props<{ log: any }>();
</script>

<div class="flex items-center gap-4 bg-white border border-slate-100 rounded-2xl p-4 shadow-xxs">
  <div class="w-12 h-12 rounded-xl bg-teal-50 text-teal-600 border border-teal-100/50 flex items-center justify-center font-black text-base shrink-0 select-none uppercase">
    {log.namaSiswa ? log.namaSiswa.charAt(0).toUpperCase() : '?'}
  </div>
  <div class="text-left flex-1 min-w-0">
    <h4 class="text-sm font-black text-slate-850 uppercase truncate leading-snug" title={log.namaSiswa}>
      {log.namaSiswa}
    </h4>
    <p class="text-[10px] font-bold text-slate-400 mt-0.5">
      NIS: {log.nis} <span class="mx-1 text-slate-200">•</span> Kelas {log.kelas || '-'}
    </p>
  </div>
  
  <!-- Status Badge -->
  <div class="shrink-0">
    {#if log.status === 'hadir'}
      <span class="px-2.5 py-1 text-[9px] font-black uppercase rounded-lg bg-teal-50 text-[#00a294] border border-teal-150/20">
        Hadir
      </span>
    {:else if log.status === 'izin'}
      <span class="px-2.5 py-1 text-[9px] font-black uppercase rounded-lg bg-amber-50 text-amber-600 border border-amber-150/20">
        Izin
      </span>
    {:else if log.status === 'sakit'}
      <span class="px-2.5 py-1 text-[9px] font-black uppercase rounded-lg bg-orange-50 text-orange-600 border border-orange-150/20">
        Sakit
      </span>
    {:else if log.status === 'magang'}
      <span class="px-2.5 py-1 text-[9px] font-black uppercase rounded-lg bg-blue-50 text-blue-600 border border-blue-150/20">
        Magang
      </span>
    {:else if log.status === 'tidak_hadir'}
      <span class="px-2.5 py-1 text-[9px] font-black uppercase rounded-lg bg-red-50 text-red-600 border border-red-150/20">
        Alpa
      </span>
    {:else}
      <span class="px-2.5 py-1 text-[9px] font-black uppercase rounded-lg bg-slate-50 text-slate-500 border border-slate-200">
        Belum Absen
      </span>
    {/if}
  </div>
</div>
