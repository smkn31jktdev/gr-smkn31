<script lang="ts">
  import { Dumbbell, Clock, Loader } from 'lucide-svelte';
  import DropdownChoice from '../../../../shared/components/DropdownChoice.svelte';

  let { 
    aktivitas = $bindable(''), 
    durasi = $bindable(30), 
    onsave, 
    loading = false 
  } = $props();

  const options = [
    { value: 'Senam Sekolah', label: 'Senam Bersama Sekolah' },
    { value: 'Senam Kemasyarakatan', label: 'Senam Kemasyarakatan / Senam Pagi' },
    { value: 'Jalan Kaki ke Sekolah', label: 'Jalan Kaki ke Sekolah (Walk to School)' },
    { value: 'Bersepeda ke Sekolah', label: 'Bersepeda ke Sekolah (Ride to School)' },
    { value: 'Gym / Renang', label: 'Gym / Renang / Fitness' },
    { value: 'Running / Jogging', label: 'Running / Jogging' },
    { value: 'Olahraga Hobi (Basket/Futsal/Badminton)', label: 'Olahraga Hobi (Basket/Futsal/Badminton)' },
    { value: 'Lainnya', label: 'Aktivitas Fisik Lainnya' }
  ];
</script>

<div class="bg-white rounded-3xl border border-slate-100/90 p-5 shadow-[0_8px_30px_rgb(0,0,0,0.015)] flex flex-col justify-between h-full">
  <div>
    <!-- Header -->
    <div class="flex items-center gap-3 mb-5">
      <div class="w-9 h-9 rounded-full bg-[#e0f2f1] flex items-center justify-center text-[#0070f3]">
        <Dumbbell class="w-5 h-5 text-[#4db6ac]" />
      </div>
      <h3 class="text-sm font-black text-slate-800 tracking-tight">Olahraga</h3>
    </div>

    <!-- Inputs -->
    <div class="space-y-4">
      <div>
        <DropdownChoice
          label="Aktivitas Fisik"
          options={options}
          bind:value={aktivitas}
          placeholder="Pilih jenis olahraga..."
        />
      </div>

      <div>
        <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1.5">Durasi (Menit)</label>
        <div class="relative">
          <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400">
            <Clock class="w-4 h-4" />
          </span>
          <input 
            type="number" 
            min="1"
            bind:value={durasi} 
            class="w-full pl-10 pr-4 py-2 border border-slate-100 rounded-xl text-xs font-bold text-slate-600 focus:outline-none focus:border-[#4db6ac] transition-all bg-slate-50/50" 
          />
        </div>
      </div>
    </div>
  </div>

  <!-- Save button -->
  <button
    type="button"
    onclick={onsave}
    disabled={loading}
    class="w-full mt-6 bg-[#4db6ac] hover:bg-[#3ca59b] disabled:bg-slate-200 disabled:cursor-not-allowed text-white py-2.5 rounded-xl text-xs font-black transition-all shadow-xxs active:scale-[0.98] inline-flex items-center justify-center gap-1.5 cursor-pointer"
  >
    {#if loading}
      <Loader class="w-3.5 h-3.5 animate-spin" />
      Menyimpan...
    {:else}
      Simpan Data
    {/if}
  </button>
</div>
