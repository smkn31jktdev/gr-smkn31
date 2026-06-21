<script lang="ts">
  import { Utensils, Loader } from 'lucide-svelte';
  import DropdownChoice from '../../../../shared/components/DropdownChoice.svelte';

  let { 
    utama = $bindable(''), 
    lauk = $bindable(''), 
    sayurBuah = $bindable(false), 
    susuSuplemen = $bindable(false), 
    onsave, 
    loading = false 
  } = $props();

  const options = [
    { value: 'Makan Sahur', label: 'Makan Sahur' },
    { value: 'Sarapan', label: 'Sarapan' },
    { value: 'Makan Siang', label: 'Makan Siang' },
    { value: 'Makan Malam', label: 'Makan Malam' }
  ];
</script>

<div class="bg-white rounded-3xl border border-slate-100/90 p-5 shadow-[0_8px_30px_rgb(0,0,0,0.015)] flex flex-col justify-between h-full">
  <div>
    <!-- Header -->
    <div class="flex items-center gap-3 mb-5">
      <div class="w-9 h-9 rounded-full bg-teal-50 flex items-center justify-center text-teal-500">
        <Utensils class="w-5 h-5" />
      </div>
      <h3 class="text-sm font-black text-slate-800 tracking-tight">Makan Sehat</h3>
    </div>

    <!-- Inputs -->
    <div class="space-y-4">
      <div>
        <DropdownChoice
          label="Makanan Utama"
          options={options}
          bind:value={utama}
          placeholder="Pilih jenis makan..."
        />
      </div>

      <div>
        <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1.5">Lauk Pauk</label>
        <input 
          type="text" 
          placeholder="Ayam, Tahu, Sayur..." 
          bind:value={lauk} 
          class="w-full px-3.5 py-2 border border-slate-100 rounded-xl text-xs font-bold text-slate-600 focus:outline-none focus:border-[#4db6ac] transition-all bg-slate-50/50" 
        />
      </div>

      <!-- Sayur / Buah (Ya/Tidak) -->
      <div class="flex items-center justify-between py-1">
        <span class="text-xs font-bold text-slate-600">Sayur / Buah?</span>
        <div class="flex gap-1.5">
          <button
            type="button"
            onclick={() => sayurBuah = true}
            class="py-1 px-3.5 rounded-lg text-xxs font-extrabold border transition-all cursor-pointer {sayurBuah === true ? 'bg-teal-50 border-teal-200 text-teal-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
          >
            Ya
          </button>
          <button
            type="button"
            onclick={() => sayurBuah = false}
            class="py-1 px-3.5 rounded-lg text-xxs font-extrabold border transition-all cursor-pointer {sayurBuah === false ? 'bg-slate-100 border-slate-200 text-slate-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
          >
            Tidak
          </button>
        </div>
      </div>

      <!-- Susu / Suplemen (Ya/Tidak) -->
      <div class="flex items-center justify-between py-1">
        <span class="text-xs font-bold text-slate-600">Susu / Suplemen?</span>
        <div class="flex gap-1.5">
          <button
            type="button"
            onclick={() => susuSuplemen = true}
            class="py-1 px-3.5 rounded-lg text-xxs font-extrabold border transition-all cursor-pointer {susuSuplemen === true ? 'bg-teal-50 border-teal-200 text-teal-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
          >
            Ya
          </button>
          <button
            type="button"
            onclick={() => susuSuplemen = false}
            class="py-1 px-3.5 rounded-lg text-xxs font-extrabold border transition-all cursor-pointer {susuSuplemen === false ? 'bg-slate-100 border-slate-200 text-slate-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
          >
            Tidak
          </button>
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
