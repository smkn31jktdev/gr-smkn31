<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Search, 
    RefreshCw, 
    Image as ImageIcon, 
    Video as VideoIcon, 
    Loader2
  } from 'lucide-svelte';
  import Photos from './media/photos/Photos.svelte';
  import Videos from './media/videos/Videos.svelte';
  
  import {
    loadingBukti,
    searchNIS,
    searchKelas,
    searchBulan,
    buktiList,
    loadBuktiData
  } from '../../logic/buktiLogic';
  import type { Bukti } from '../../logic/buktiLogic';
  import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';

  // Available classes
  const classesList = [
    { value: '', label: 'Semua Kelas' },
    { value: 'X LP', label: 'X LP' },
    { value: 'XI LP', label: 'XI LP' },
    { value: 'XII LP', label: 'XII LP' }
  ];

  // Generate Month list options (Indonesian month names)
  function getMonthOptions() {
    const options = [{ val: '', label: 'Semua Periode' }];
    const monthNames = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    
    const anchorYear = 2026;
    const anchorMonth = 5; // June (0-indexed)
    
    for (let i = -12; i <= 3; i++) {
      const d = new Date(anchorYear, anchorMonth + i, 1);
      const val = d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0');
      const label = `${monthNames[d.getMonth()]} ${d.getFullYear()}`;
      options.push({ val, label });
    }
    return options;
  }

  const monthsList = getMonthOptions();

  function formatMonthYear(val: string) {
    if (!val) return '';
    const [year, month] = val.split('-');
    const monthNames = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    const mIndex = parseInt(month, 10) - 1;
    return `${monthNames[mIndex] || month} ${year}`;
  }

  // Modals state
  let photoModalOpen = $state(false);
  let activePhotos = $state<string[]>([]);
  let activeStudentName = $state('');

  let videoModalOpen = $state(false);
  let activeVideos = $state<string[]>([]);

  function openPhotoModal(studentName: string, photos: string[]) {
    activeStudentName = studentName;
    activePhotos = photos;
    photoModalOpen = true;
  }

  function openVideoModal(studentName: string, videos: string[]) {
    activeStudentName = studentName;
    activeVideos = videos;
    videoModalOpen = true;
  }

  // Debounced search trigger
  let debounceTimer: any;
  function handleSearchInput(e: Event) {
    clearTimeout(debounceTimer);
    const target = e.target as HTMLInputElement;
    searchNIS.set(target.value);
    debounceTimer = setTimeout(() => {
      loadBuktiData();
    }, 300);
  }

  onMount(() => {
    loadBuktiData();
  });
</script>

<div class="space-y-6 select-none font-sans max-w-[1600px] mx-auto p-1 pb-10">
  <!-- Title Header Block -->
  <div>
    <h2 class="text-xl font-bold tracking-tight text-slate-800 font-display uppercase">Bukti Kegiatan</h2>
    <p class="text-xs text-slate-400 mt-0.5 font-medium">Tinjau dokumentasi foto dan video kegiatan yang diunggah oleh siswa setiap bulannya.</p>
  </div>

  <!-- Filters Toolbar -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white border border-slate-100 rounded-2xl p-4 shadow-xs">
    <div class="flex flex-1 flex-col sm:flex-row gap-3">
      <!-- Search Input -->
      <div class="relative flex-1 max-w-md">
        <Search class="w-4 h-4 text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input 
          type="text" 
          placeholder="Filter NIS..."
          value={$searchNIS}
          oninput={handleSearchInput}
          class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-700 text-xs font-bold py-2.5 pl-10 pr-4 rounded-xl outline-none transition-all"
        />
      </div>

      <!-- Class Selector -->
      <div class="min-w-[130px] text-left">
        <DropdownChoice
          options={classesList.map(cls => ({ value: cls.value, label: cls.label }))}
          value={$searchKelas}
          onchange={(val) => {
            searchKelas.set(val);
            loadBuktiData();
          }}
          placeholder="Semua Kelas"
        />
      </div>

      <!-- Month Selector -->
      <div class="min-w-[140px] text-left">
        <DropdownChoice
          options={monthsList.map(m => ({ value: m.val, label: m.label }))}
          value={$searchBulan}
          onchange={(val) => {
            searchBulan.set(val);
            loadBuktiData();
          }}
          placeholder="Semua Periode"
        />
      </div>
    </div>

    <!-- Refresh Button -->
    <button 
      onclick={loadBuktiData}
      disabled={$loadingBukti}
      class="flex items-center justify-center gap-1.5 px-4 py-2.5 bg-slate-50 hover:bg-slate-100 border border-slate-200/50 text-slate-600 font-bold text-xs rounded-xl transition-all cursor-pointer select-none"
    >
      <RefreshCw class="w-3.5 h-3.5 {$loadingBukti ? 'animate-spin' : ''}" />
      Refresh
    </button>
  </div>

  <!-- Content Section -->
  {#if $loadingBukti}
    <div class="flex flex-col items-center justify-center py-32 text-slate-400 bg-white border border-slate-100 rounded-3xl shadow-xs">
      <Loader2 class="w-8 h-8 animate-spin text-slate-350 mb-3" />
      <p class="text-xs font-semibold">Memuat bukti kegiatan...</p>
    </div>
  {:else if $buktiList.length === 0}
    <div class="flex flex-col items-center justify-center py-32 text-slate-400 bg-white border border-slate-100 rounded-3xl shadow-xs">
      <ImageIcon class="w-10 h-10 text-slate-200 mb-3" />
      <p class="text-xs font-bold">Tidak ada bukti laporan terdeteksi</p>
    </div>
  {:else}
    <!-- Grid Card list -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      {#each $buktiList as item}
        <div class="bg-white border border-slate-100 hover:border-slate-200/80 rounded-2xl p-5 shadow-xs transition-all flex flex-col justify-between min-h-[160px] group">
          <!-- Card Header Info -->
          <div>
            <div class="flex items-start justify-between gap-3">
              <div class="flex items-center gap-3">
                <!-- User Profile Initials Badge -->
                <div class="w-9 h-9 rounded-lg bg-slate-50 border border-slate-100 text-slate-500 flex items-center justify-center font-bold text-xs shrink-0 select-none">
                  {item.namaSiswa.charAt(0).toUpperCase()}
                </div>
                
                <div class="space-y-0.5">
                  <h4 class="text-xs font-bold text-slate-700 uppercase group-hover:text-slate-900 transition-colors truncate max-w-[170px] sm:max-w-none">
                    {item.namaSiswa}
                  </h4>
                  <p class="text-[10px] font-semibold text-slate-400 uppercase tracking-wide">
                    {item.kelas}
                  </p>
                </div>
              </div>
              
              <!-- Month/Period -->
              <span class="text-[10px] font-bold text-slate-400 shrink-0">
                {formatMonthYear(item.bulan)}
              </span>
            </div>
          </div>

          <!-- Card Actions (Buttons) -->
          <div class="grid grid-cols-1 gap-2.5 mt-5">
            <!-- Lihat Foto Button -->
            {#if item.foto && item.foto.length > 0}
              <button 
                onclick={() => openPhotoModal(item.namaSiswa, item.foto)}
                class="w-full flex items-center justify-center gap-1.5 py-2.5 bg-slate-50 hover:bg-slate-100 border border-slate-200/50 rounded-xl text-slate-650 font-bold text-xs transition-all cursor-pointer"
              >
                <ImageIcon class="w-3.5 h-3.5" />
                Lihat Foto ({item.foto.length})
              </button>
            {/if}

            <!-- Lihat Video Button -->
            {#if item.linkYT && item.linkYT.length > 0}
              <button 
                onclick={() => openVideoModal(item.namaSiswa, item.linkYT)}
                class="w-full flex items-center justify-center gap-1.5 py-2.5 bg-rose-50 hover:bg-rose-100/60 border border-rose-100/50 rounded-xl text-rose-600 font-bold text-xs transition-all cursor-pointer"
              >
                <VideoIcon class="w-3.5 h-3.5" />
                Lihat Video ({item.linkYT.length})
              </button>
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

<Photos 
  bind:open={photoModalOpen} 
  photos={activePhotos} 
  studentName={activeStudentName} 
/>

<Videos 
  bind:open={videoModalOpen} 
  videos={activeVideos} 
  studentName={activeStudentName} 
/>

<style>
  /* Custom scrollbar styling for a clean sleek feel */
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 99px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
  }
</style>
