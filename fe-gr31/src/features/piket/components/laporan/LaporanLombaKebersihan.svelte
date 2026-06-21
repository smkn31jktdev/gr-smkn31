<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    RefreshCw, 
    Calendar, 
    ChevronLeft, 
    ChevronRight,
    Loader2,
    X,
    ExternalLink,
    Camera,
    User,
    FileText
  } from 'lucide-svelte';
  
  import { currentUser } from '../../../../stores/authStore';
  
  import {
    logs,
    total,
    page,
    limit,
    loading,
    hasMore,
    kelasList,
    selectedKelas,
    selectedTanggal,
    selectedBulan,
    loadKelasList,
    loadLombaData,
    handleFilter,
    changePage
  } from '../../logic/piketLombaLogic';

  import { addToast } from '../../../../stores/uiStore';
  import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';
  import DatePicker from '../../../shared/components/DatePicker.svelte';

  // Generate Month list options (Indonesian month names)
  function getMonthOptions() {
    const options = [{ val: '', label: 'Semua Bulan' }];
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

  function formatTanggalIndo(tglStr: string) {
    if (!tglStr) return '';
    try {
      const d = new Date(tglStr);
      return d.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
    } catch {
      return tglStr;
    }
  }

  // Modals state for Image Lightbox
  let photoModalOpen = $state(false);
  let activePhotos = $state<string[]>([]);
  let activePhotoIndex = $state(0);
  let activeClassName = $state('');
  let activeDate = $state('');

  // Wali Kelas role check and class lock
  let role = $derived($currentUser?.role);
  let isWalas = $derived(
    role === 'walas' || 
    role === 'guru_wali' || 
    $currentUser?.is_walas === true || 
    $currentUser?.isWalas === true
  );
  let userKelas = $derived($currentUser?.kelas || '');

  function openPhotoModal(className: string, dateStr: string, photos: string[]) {
    activeClassName = className;
    activeDate = dateStr;
    activePhotos = photos;
    activePhotoIndex = 0;
    photoModalOpen = true;
  }

  onMount(() => {
    if (isWalas && userKelas) {
      selectedKelas.set(userKelas);
    }
    loadKelasList();
    loadLombaData();
  });
</script>

<div class="space-y-4">
  <!-- Filters Toolbar -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white border border-slate-100 rounded-2xl p-4 shadow-xs text-left">
    <div class="flex flex-1 flex-wrap gap-3">
      
      <!-- Class Selector -->
      {#if !isWalas}
        <div class="flex flex-col">
          <label for="filter-kelas" class="text-[9px] font-black uppercase text-slate-400 tracking-wider mb-1">Kelas</label>
          <div class="min-w-[130px]">
            <DropdownChoice
              options={[{ value: '', label: 'Semua Kelas' }, ...$kelasList.map(cls => ({ value: cls, label: cls }))]}
              bind:value={$selectedKelas}
              onchange={handleFilter}
              placeholder="Semua Kelas"
            />
          </div>
        </div>
      {:else}
        <div class="flex flex-col justify-end">
          <span class="text-[9px] font-black uppercase text-slate-400 tracking-wider mb-1">Kelas Diampu</span>
          <span class="inline-flex items-center px-3 py-2 bg-teal-50 border border-teal-100 text-teal-700 text-xs font-bold rounded-xl h-9">
            {userKelas}
          </span>
        </div>
      {/if}

      <!-- Month Selector -->
      <div class="flex flex-col">
        <label for="filter-bulan" class="text-[9px] font-black uppercase text-slate-400 tracking-wider mb-1">Periode Bulan</label>
        <div class="min-w-[140px]">
          <DropdownChoice
            options={monthsList.map(m => ({ value: m.val, label: m.label }))}
            bind:value={$selectedBulan}
            onchange={() => {
              selectedTanggal.set('');
              handleFilter();
            }}
            placeholder="Periode Bulan"
          />
        </div>
      </div>

      <!-- Specific Date Picker -->
      <div class="flex flex-col">
        <label for="filter-tanggal" class="text-[9px] font-black uppercase text-slate-400 tracking-wider mb-1">Pilih Tanggal Spesifik</label>
        <div class="w-48">
          <DatePicker
            bind:value={$selectedTanggal}
            onchange={() => {
              selectedBulan.set('');
              handleFilter();
            }}
            placeholder="Pilih Tanggal"
          />
        </div>
      </div>
    </div>

    <!-- Refresh Button -->
    <button 
      onclick={loadLombaData}
      disabled={$loading}
      class="flex items-center justify-center gap-1.5 px-4 py-2.5 bg-slate-50 hover:bg-slate-100 border border-slate-200 text-slate-600 font-bold text-xs rounded-xl transition-all cursor-pointer select-none self-end h-9"
    >
      <RefreshCw class="w-3.5 h-3.5 {$loading ? 'animate-spin' : ''}" />
      Refresh
    </button>
  </div>

  <!-- Content Section -->
  {#if $loading}
    <div class="flex flex-col items-center justify-center py-24 text-slate-450 bg-white border border-slate-100 rounded-2xl shadow-xs">
      <Loader2 class="w-8 h-8 animate-spin text-[#00a294] mb-3" />
      <p class="text-xs font-semibold">Memuat rekap laporan kebersihan kelas...</p>
    </div>
  {:else if $logs.length === 0}
    <div class="flex flex-col items-center justify-center py-24 text-slate-400 bg-white border border-slate-100 rounded-2xl shadow-xs">
      <Camera class="w-10 h-10 text-slate-200 mb-3" />
      <p class="text-xs font-bold">Tidak ada laporan kebersihan kelas terdeteksi</p>
      <p class="text-[10px] text-slate-450 mt-1">Gunakan filter di atas untuk melihat laporan kelas lainnya.</p>
    </div>
  {:else}
    <!-- Grid Card list -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5 text-left">
      {#each $logs as item}
        <div class="bg-white border border-slate-100 hover:border-slate-200/80 rounded-2xl p-5 shadow-xs transition-all flex flex-col justify-between group min-h-[200px]">
          <!-- Card Info -->
          <div class="space-y-3">
            <div class="flex items-start justify-between border-b border-slate-50 pb-2.5">
              <div class="space-y-0.5">
                <span class="px-2 py-0.5 rounded bg-teal-50 text-teal-700 text-[10px] font-black uppercase tracking-wider">
                  Kelas {item.kelas}
                </span>
                <p class="text-[10px] font-bold text-slate-400 pt-1">
                  {formatTanggalIndo(item.tanggal)}
                </p>
              </div>
              <span class="text-[9px] font-bold text-slate-400 flex items-center gap-1">
                <User class="w-3.5 h-3.5" /> {item.namaSiswa}
              </span>
            </div>

            <!-- Notes -->
            {#if item.catatan}
              <div class="p-3 rounded-xl bg-slate-50/50 border border-slate-100/60 flex items-start gap-2">
                <FileText class="w-3.5 h-3.5 text-slate-400 shrink-0 mt-0.5" />
                <p class="text-[11px] text-slate-600 font-medium leading-relaxed">{item.catatan}</p>
              </div>
            {/if}
          </div>

          <!-- Photo Previews and Button -->
          <div class="mt-4 space-y-2">
            {#if item.foto && item.foto.length > 0}
              <div class="grid grid-cols-4 gap-2">
                {#each item.foto.slice(0, 4) as photoUrl, idx}
                  <button 
                    onclick={() => openPhotoModal(item.kelas, item.tanggal, item.foto)}
                    class="aspect-video rounded-lg overflow-hidden bg-slate-50 border border-slate-100 hover:opacity-85 transition-opacity relative group/img cursor-pointer"
                  >
                    <img src={photoUrl} alt="Preview kebersihan" class="w-full h-full object-cover" />
                    {#if idx === 3 && item.foto.length > 4}
                      <div class="absolute inset-0 bg-black/55 flex items-center justify-center text-white font-bold text-xs">
                        +{item.foto.length - 4}
                      </div>
                    {/if}
                  </button>
                {/each}
              </div>

              <button 
                onclick={() => openPhotoModal(item.kelas, item.tanggal, item.foto)}
                class="w-full flex items-center justify-center gap-1.5 py-2 bg-slate-550 hover:bg-slate-600 text-white rounded-xl font-bold text-xs transition-all cursor-pointer mt-1"
              >
                <Camera class="w-3.5 h-3.5" />
                Lihat Foto Kebersihan ({item.foto.length})
              </button>
            {/if}
          </div>
        </div>
      {/each}
    </div>

    <!-- Pagination Footer -->
    {#if $total > 0}
      <div class="flex items-center justify-between border-t border-slate-100 pt-5 mt-4">
        <span class="text-xs font-semibold text-slate-400">
          Total {$total} Lomba Kebersihan Kelas
        </span>
        <div class="flex items-center gap-2">
          <button
            onclick={() => changePage(-1)}
            disabled={$page <= 1}
            class="flex items-center gap-1 px-3 py-1.5 border border-slate-200 bg-white text-slate-700 text-xs font-bold rounded-lg transition-all hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            <ChevronLeft class="w-3.5 h-3.5" /> Sebelumnya
          </button>
          <span class="text-xs font-bold text-slate-650 px-2 select-none">Halaman {$page}</span>
          <button
            onclick={() => changePage(1)}
            disabled={!$hasMore}
            class="flex items-center gap-1 px-3 py-1.5 border border-slate-200 bg-white text-slate-700 text-xs font-bold rounded-lg transition-all hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed cursor-pointer"
          >
            Berikutnya <ChevronRight class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    {/if}
  {/if}
</div>

<!-- Modal: Photo Lightbox Gallery -->
{#if photoModalOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4 animate-fade-in">
    <div class="relative w-full max-w-4xl bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h3 class="text-sm font-black text-slate-700 uppercase">Dokumentasi Kebersihan Kelas</h3>
          <p class="text-[10px] text-slate-450 font-bold mt-0.5">Kelas: {activeClassName} — {formatTanggalIndo(activeDate)}</p>
        </div>
        <button 
          onclick={() => photoModalOpen = false}
          class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-600 transition-colors border-none cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Main Photo Area -->
      <div class="flex-1 bg-slate-50 flex items-center justify-center relative p-6 min-h-[300px] overflow-hidden">
        <!-- Main Image -->
        <img 
          src={activePhotos[activePhotoIndex]} 
          alt="Foto Kebersihan Kelas" 
          class="max-w-full max-h-[50vh] sm:max-h-[60vh] object-contain rounded-xl shadow-xs" 
        />

        <!-- Prev button -->
        {#if activePhotos.length > 1}
          <button 
            onclick={() => activePhotoIndex = (activePhotoIndex - 1 + activePhotos.length) % activePhotos.length}
            class="absolute left-4 p-2.5 rounded-full bg-white/95 border border-slate-100 shadow-md text-slate-650 hover:text-slate-800 transition-all cursor-pointer border-none"
          >
            <ChevronLeft class="w-4 h-4" />
          </button>
          
          <!-- Next button -->
          <button 
            onclick={() => activePhotoIndex = (activePhotoIndex + 1) % activePhotos.length}
            class="absolute right-4 p-2.5 rounded-full bg-white/95 border border-slate-100 shadow-md text-slate-650 hover:text-slate-800 transition-all cursor-pointer border-none"
          >
            <ChevronRight class="w-4 h-4" />
          </button>
        {/if}
      </div>

      <!-- Indicator Footer -->
      {#if activePhotos.length > 1}
        <div class="px-6 py-3 border-t border-slate-100 flex items-center justify-center gap-1.5 bg-white shrink-0">
          {#each activePhotos as _, idx}
            <button 
              onclick={() => activePhotoIndex = idx}
              aria-label="Foto ke-{idx + 1}"
              class="w-2 h-2 rounded-full transition-all border-none cursor-pointer {activePhotoIndex === idx ? 'bg-slate-800 w-4' : 'bg-slate-200 hover:bg-slate-300'}"
            ></button>
          {/each}
        </div>
      {/if}
    </div>
  </div>
{/if}
