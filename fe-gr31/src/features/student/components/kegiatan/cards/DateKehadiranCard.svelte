<script lang="ts">
  import { Calendar, ShieldCheck, ShieldAlert, ArrowRight, Loader, Lock } from 'lucide-svelte';
  import type { Kehadiran } from '../../../types/student.types';
  import DatePicker from '../../../../shared/components/DatePicker.svelte';

  let { 
    tanggal = $bindable(), 
    maxDate = '',
    kehadiran = null as Kehadiran | null,
    loadingKehadiran = false,
    onchange,
    onabsen 
  } = $props();

  let formattedDate = $derived.by(() => {
    if (!tanggal) return '';
    try {
      const d = new Date(tanggal);
      return d.toLocaleDateString('id-ID', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
        year: 'numeric'
      });
    } catch {
      return tanggal;
    }
  });
</script>

<div class="bg-white rounded-3xl border border-slate-100/90 p-5 shadow-[0_8px_30px_rgb(0,0,0,0.015)] grid grid-cols-1 md:grid-cols-12 gap-6 items-center text-left">
  
  <!-- Date Selector (Cols 5) -->
  <div class="md:col-span-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4 md:border-r border-slate-100 md:pr-6">
    <div class="space-y-1">
      <div class="flex items-center gap-2.5">
        <div class="w-7 h-7 rounded-full bg-teal-50 flex items-center justify-center text-[#4db6ac]">
          <Calendar class="w-3.5 h-3.5" />
        </div>
        <span class="text-[10px] font-black uppercase tracking-wider text-slate-400 flex items-center gap-1">
          Tanggal Kegiatan
          <span class="text-[9px] font-extrabold text-amber-600 bg-amber-50 px-1.5 py-0.5 rounded border border-amber-200/50">Hari Ini</span>
        </span>
      </div>
      <div class="text-[10px] font-bold text-[#4db6ac]">
        {formattedDate}
      </div>
    </div>
    <div class="shrink-0 w-48">
      <DatePicker 
        bind:value={tanggal} 
        onchange={onchange}
        maxDate={maxDate}
        placeholder="Pilih Tanggal"
        disabled={true}
      />
    </div>
  </div>

  <!-- Kehadiran Status (Cols 7) -->
  <div class="md:col-span-7 flex flex-col sm:flex-row sm:items-center justify-between gap-4 pl-0 md:pl-2">
    <!-- Left Details -->
    <div class="flex items-center gap-3">
      <div class="w-8 h-8 rounded-full bg-slate-50 flex items-center justify-center text-slate-500 shrink-0">
        {#if kehadiran}
          <ShieldCheck class="w-4.5 h-4.5 text-[#4db6ac]" />
        {:else}
          <ShieldAlert class="w-4.5 h-4.5 text-amber-500" />
        {/if}
      </div>
      <div>
        <div class="flex items-center gap-2.5 mb-1">
          <span class="text-[10px] font-black uppercase tracking-wider text-slate-400">Status Absensi</span>
          
          {#if loadingKehadiran}
            <span class="px-2.5 py-0.5 bg-slate-50 rounded-full text-[8px] font-black text-slate-400">Loading...</span>
          {:else}
            <span class="px-2.5 py-0.5 rounded-full text-[8px] font-black uppercase shadow-xxs
              {kehadiran?.status === 'hadir' ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' : ''}
              {kehadiran?.status === 'izin' ? 'bg-sky-50 text-sky-600 border border-sky-100' : ''}
              {kehadiran?.status === 'sakit' ? 'bg-amber-50 text-amber-600 border border-amber-100' : ''}
              {kehadiran?.status === 'tidak_hadir' ? 'bg-rose-50 text-rose-600 border border-rose-100' : ''}
              {!kehadiran ? 'bg-amber-50 text-amber-700 border border-amber-200 animate-pulse' : ''}"
            >
              {kehadiran ? kehadiran.status.replace('_', ' ') : 'Belum Absen'}
            </span>
          {/if}
        </div>

        {#if loadingKehadiran}
          <p class="text-[10px] text-slate-400 font-bold">Memeriksa database...</p>
        {:else if kehadiran}
          <p class="text-[10px] text-slate-400 font-bold leading-normal">
            {kehadiran.status === 'hadir' ? 'Terverifikasi' : 'Tercatat'} 
            {#if kehadiran.waktuAbsen}
              • Pukul {kehadiran.waktuAbsen.substring(0, 8)} WIB
            {/if}
            {#if kehadiran.status === 'hadir'}
              • Jarak {kehadiran.jarak}m
            {/if}
          </p>
        {:else}
          <p class="text-[10px] text-slate-500 font-bold leading-normal">
            Silakan melakukan check-in kehadiran hari ini.
          </p>
        {/if}
      </div>
    </div>

    <!-- Action Button / Suffix -->
    <div class="shrink-0">
      {#if !kehadiran && !loadingKehadiran}
        <button 
          type="button" 
          onclick={onabsen}
          class="inline-flex items-center gap-1.5 px-4 py-2 border border-amber-200 hover:border-amber-400 bg-amber-50 hover:bg-amber-100 rounded-xl text-xs font-black text-amber-700 transition-all cursor-pointer shadow-xxs active:scale-97"
        >
          Absen Sekarang
          <ArrowRight class="w-3.5 h-3.5" />
        </button>
      {/if}
    </div>
  </div>

</div>
