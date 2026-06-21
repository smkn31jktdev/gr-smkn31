<script lang="ts">
  import { User } from 'lucide-svelte';

  let { 
    name = '', 
    nis = '', 
    kelas = '' 
  } = $props();
</script>

<div class="bg-white rounded-3xl border border-slate-100/90 p-5 shadow-[0_8px_30px_rgb(0,0,0,0.015)] flex flex-col md:flex-row md:items-center justify-between gap-6">
  <div class="flex items-center gap-3">
    <div class="w-9 h-9 rounded-full bg-teal-50 flex items-center justify-center text-[#4db6ac]">
      <User class="w-5 h-5" />
    </div>
    <span class="text-[10px] font-black uppercase tracking-wider text-[#4db6ac]">Data Siswa</span>
  </div>

  <div class="flex flex-wrap items-center gap-x-8 gap-y-3">
    <!-- Student Name -->
    <div>
      <span class="block text-[8px] font-black uppercase tracking-widest text-slate-400">Nama Siswa</span>
      <span class="text-xs font-black text-slate-700 tracking-tight">{name}</span>
    </div>

    <!-- NIS -->
    <div>
      <span class="block text-[8px] font-black uppercase tracking-widest text-slate-400">NIS</span>
      <span class="text-xs font-black text-slate-700 tracking-tight">{nis}</span>
    </div>

    <!-- Kelas -->
    <div>
      <span class="block text-[8px] font-black uppercase tracking-widest text-slate-400 mb-0.5">Kelas</span>
      <span class="inline-flex px-3.5 py-1 border border-teal-200 bg-teal-50/50 rounded-full text-[10px] font-extrabold text-[#4db6ac] shadow-xxs">
        {kelas}
      </span>
    </div>
  </div>
</div>
