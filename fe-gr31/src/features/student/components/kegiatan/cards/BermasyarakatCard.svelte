<script lang="ts">
  import { Users, Clock, Loader } from 'lucide-svelte';
  import DropdownChoice from '../../../../shared/components/DropdownChoice.svelte';
  import ClockTimer from '../../../../shared/components/ClockTimer.svelte';

  let { 
    kegiatan = $bindable(''), 
    lokasi = $bindable(''), 
    waktu = $bindable(''), 
    diketahuiOT = $bindable(false), 
    onsave, 
    loading = false 
  } = $props();

  const options = [
    { value: 'Membersihkan Tempat Ibadah', label: 'Membersihkan Tempat Ibadah' },
    { value: 'Membersihkan Got / Jalan', label: 'Membersihkan Got / Jalan Lingkungan' },
    { value: 'Merawat Tanaman', label: 'Merawat Tanaman / Penghijauan Umum' },
    { value: 'Petugas Ibadah', label: 'Petugas Ibadah (Muazin/Koster/dll)' },
    { value: 'Khotib / Penceramah', label: 'Khotib / Penceramah / Kultum' },
    { value: 'Mengajar Ngaji / Taklim', label: 'Mengajar Ngaji / Taklim / Belajar Bersama' },
    { value: 'Aktivitas Sosial Lainnya', label: 'Aktivitas Sosial Lainnya' }
  ];
</script>

<div class="bg-white rounded-3xl border border-slate-100/90 p-6 shadow-[0_8px_30px_rgb(0,0,0,0.015)] space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-3">
    <div class="w-9 h-9 rounded-full bg-indigo-50 flex items-center justify-center text-indigo-500">
      <Users class="w-5 h-5" />
    </div>
    <h3 class="text-sm font-black text-slate-800 tracking-tight">Bermasyarakat</h3>
  </div>

  <!-- Inputs Grid -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-5 items-end">
    <div>
      <DropdownChoice
        label="Jenis Kegiatan"
        options={options}
        bind:value={kegiatan}
        placeholder="Pilih jenis kegiatan..."
      />
    </div>

    <div>
      <label for="tempat-kegiatan" class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1.5">Tempat Kegiatan</label>
      <input 
        id="tempat-kegiatan"
        type="text" 
        placeholder="Masjid, Balai Warga, dll" 
        bind:value={lokasi} 
        class="w-full px-3.5 py-2 border border-slate-100 rounded-xl text-xs font-bold text-slate-600 focus:outline-none focus:border-[#4db6ac] transition-all bg-slate-50/50" 
      />
    </div>

    <div>
      <ClockTimer
        mode="timepicker"
        label="Waktu Pelaksanaan"
        bind:value={waktu}
        placeholder="--:--"
      />
    </div>
  </div>

  <!-- Verification Statement Checkbox -->
  <label class="flex items-start gap-3 p-4 bg-slate-50/50 border border-slate-100 rounded-2xl cursor-pointer hover:bg-slate-50 transition-colors select-none">
    <input 
      type="checkbox" 
      bind:checked={diketahuiOT} 
      class="w-4.5 h-4.5 rounded border-slate-200 text-[#4db6ac] focus:ring-[#4db6ac] mt-0.5" 
    />
    <div class="flex-1">
      <span class="text-xs font-bold text-slate-600 leading-snug">
        Saya menyatakan kegiatan ini diketahui oleh Orang Tua / Wali / RT setempat
      </span>
      <p class="text-[9px] text-slate-400 font-semibold mt-0.5">
        (Paraf Orang Tua/Wali diperlukan untuk verifikasi fisik jurnal bulanan)
      </p>
    </div>
  </label>

  <!-- Save Button at the Bottom -->
  <div class="flex justify-end">
    <button
      type="button"
      onclick={onsave}
      disabled={loading}
      class="px-8 py-2.5 bg-[#4db6ac] hover:bg-[#3ca59b] disabled:bg-slate-200 disabled:cursor-not-allowed text-white rounded-xl text-xs font-black transition-all shadow-xxs active:scale-[0.98] inline-flex items-center justify-center gap-1.5 cursor-pointer"
    >
      {#if loading}
        <Loader class="w-3.5 h-3.5 animate-spin" />
        Menyimpan...
      {:else}
        Simpan Data
      {/if}
    </button>
  </div>
</div>
