<script lang="ts">
  import { BookOpen, Loader } from 'lucide-svelte';
  import DropdownChoice from '../../../../shared/components/DropdownChoice.svelte';

  let { 
    belajarMandiri = $bindable(false),
    kitabSuci = $bindable(false),
    bukuUmum = $bindable(false),
    bukuMapel = $bindable(false),
    tugasPR = $bindable(false),
    onsave, 
    loading = false 
  } = $props();

  const options = [
    { value: 'kitabSuci', label: 'Membaca Kitab Suci' },
    { value: 'bukuUmum', label: 'Membaca Buku Bacaan / Novel' },
    { value: 'bukuMapel', label: 'Membaca Buku Pelajaran' },
    { value: 'tugasPR', label: 'Mengerjakan Tugas / PR' }
  ];

  let selectedTopic = $state('');

  // Sync from props (parent) to local state
  $effect(() => {
    if (kitabSuci) {
      selectedTopic = 'kitabSuci';
    } else if (bukuUmum) {
      selectedTopic = 'bukuUmum';
    } else if (bukuMapel) {
      selectedTopic = 'bukuMapel';
    } else if (tugasPR) {
      selectedTopic = 'tugasPR';
    } else {
      selectedTopic = '';
    }
  });

  // Sync from local state to props
  function handleTopicChange(val: string) {
    selectedTopic = val;
    kitabSuci = val === 'kitabSuci';
    bukuUmum = val === 'bukuUmum';
    bukuMapel = val === 'bukuMapel';
    tugasPR = val === 'tugasPR';
  }
</script>

<div class="bg-white rounded-3xl border border-slate-100/90 p-5 shadow-[0_8px_30px_rgb(0,0,0,0.015)] flex flex-col justify-between h-full">
  <div>
    <!-- Header -->
    <div class="flex items-center gap-3 mb-5">
      <div class="w-9 h-9 rounded-full bg-violet-50 flex items-center justify-center text-violet-500">
        <BookOpen class="w-5 h-5" />
      </div>
      <h3 class="text-sm font-black text-slate-800 tracking-tight">Belajar Mandiri</h3>
    </div>

    <!-- Inputs -->
    <div class="space-y-4">
      <div>
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-2.5">Belajar mandiri hari ini?</span>
        <div class="grid grid-cols-2 gap-2">
          <button
            type="button"
            onclick={() => belajarMandiri = true}
            class="py-2 px-4 rounded-xl text-xs font-extrabold border transition-all cursor-pointer {belajarMandiri === true ? 'bg-violet-50 border-violet-200 text-violet-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
          >
            Ya
          </button>
          <button
            type="button"
            onclick={() => { belajarMandiri = false; kitabSuci = false; bukuUmum = false; bukuMapel = false; tugasPR = false; }}
            class="py-2 px-4 rounded-xl text-xs font-extrabold border transition-all cursor-pointer {belajarMandiri === false ? 'bg-slate-100 border-slate-200 text-slate-700 shadow-xxs' : 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'}"
          >
            Tidak
          </button>
        </div>
      </div>

      {#if belajarMandiri}
        <!-- Dropdown for learning indicators -->
        <div class="pt-2.5 border-t border-slate-50">
          <DropdownChoice
            label="Aktivitas Belajar"
            options={options}
            bind:value={selectedTopic}
            onchange={handleTopicChange}
            placeholder="Pilih aktivitas belajar..."
          />
        </div>
      {/if}
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
