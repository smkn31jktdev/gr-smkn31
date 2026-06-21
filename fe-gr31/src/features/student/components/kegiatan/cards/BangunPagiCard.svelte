<script lang="ts">
  import { Sun, Clock, Loader } from 'lucide-svelte';

  let { 
    waktu = $bindable(), 
    doa = $bindable(), 
    onsave, 
    loading = false 
  } = $props();
</script>

<div class="bg-white rounded-3xl border border-slate-100/90 p-5 shadow-[0_8px_30px_rgb(0,0,0,0.015)] flex flex-col justify-between h-full">
  <div>
    <!-- Header -->
    <div class="flex items-center gap-3 mb-5">
      <div class="w-9 h-9 rounded-full bg-amber-50 flex items-center justify-center text-amber-500">
        <Sun class="w-5 h-5" />
      </div>
      <h3 class="text-sm font-black text-slate-800 tracking-tight">Bangun Pagi</h3>
    </div>

    <!-- Inputs -->
    <div class="space-y-4">
      <div>
        <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1.5">Waktu Bangun</label>
        <div class="relative">
          <span class="absolute left-3.5 top-1/2 -translate-y-1/2 text-slate-400">
            <Clock class="w-4 h-4" />
          </span>
          <input 
            type="time" 
            bind:value={waktu} 
            disabled
            class="w-full pl-10 pr-4 py-2 border border-slate-100 rounded-xl text-xs font-bold text-slate-400 focus:outline-none bg-slate-50/50 cursor-not-allowed" 
          />
        </div>
      </div>

      <div>
        <label class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1.5">Membaca Doa?</label>
        <div class="grid grid-cols-2 gap-2">
          <button
            type="button"
            onclick={() => doa = true}
            class="py-2 px-4 rounded-xl text-xs font-extrabold border transition-all cursor-pointer {doa === true ? 'bg-amber-50 border-amber-200 text-amber-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
          >
            Ya
          </button>
          <button
            type="button"
            onclick={() => doa = false}
            class="py-2 px-4 rounded-xl text-xs font-extrabold border transition-all cursor-pointer {doa === false ? 'bg-slate-100 border-slate-200 text-slate-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
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
