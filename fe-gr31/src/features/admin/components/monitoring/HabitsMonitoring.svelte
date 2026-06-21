<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Search, 
    Calendar, 
    RefreshCw, 
    ArrowLeft, 
    ArrowRight, 
    Sun,
    Heart, 
    Flame, 
    BookOpen, 
    Users2, 
    Compass, 
    Moon, 
    Sparkles,
    Loader2 
  } from 'lucide-svelte';
  
  import { 
    listG7DailyJournals, 
    getHabitMetadata 
  } from '../../logic/monitoringLogic';
  import Content from './content/Content.svelte';
  import type { G7Jurnal } from '../../../student/types/student.types';
  import { addToast } from '../../../../stores/uiStore';

  // Props
  let { kegiatan } = $props<{ kegiatan: string }>();
  const getTodayStr = () => {
    const d = new Date();
    return d.toLocaleDateString('sv-SE');
  };

  let selectedKelas = $state('');
  let searchQuery = $state('');
  let currentPage = $state(1);
  let limit = 20;

  let items = $state<G7Jurnal[]>([]);
  let total = $state(0);
  let hasMore = $state(false);
  let loading = $state(false);

  let isModalOpen = $state(false);
  let selectedStudent = $state<G7Jurnal | null>(null);
  let loadingHistory = $state(false);
  let historyItems = $state<G7Jurnal[]>([]);

  async function openDetail(item: G7Jurnal) {
    selectedStudent = item;
    isModalOpen = true;
    loadingHistory = true;
    historyItems = [];
    try {
      const res = await listG7DailyJournals({
        nis: item.nis,
        limit: 100
      });
      historyItems = res.items.sort((a, b) => b.tanggal.localeCompare(a.tanggal));
    } catch (e) {
      console.error(e);
      addToast('Gagal memuat riwayat harian siswa', 'error');
    } finally {
      loadingHistory = false;
    }
  }

  function closeModal() {
    isModalOpen = false;
    selectedStudent = null;
    historyItems = [];
  }

  let meta = $derived(getHabitMetadata(kegiatan));

  async function loadData() {
    loading = true;
    try {
      const res = await listG7DailyJournals({
        kelas: selectedKelas,
        q: searchQuery,
        page: currentPage,
        limit: limit
      });
      items = res.items;
      total = res.total;
      hasMore = res.hasMore;
    } catch (e) {
      console.error(e);
      addToast('Gagal memuat data monitoring', 'error');
    } finally {
      loading = false;
    }
  }

  function handleFilterChange() {
    currentPage = 1;
    loadData();
  }

  function handleRefresh() {
    loadData();
    addToast('Data berhasil diperbarui', 'success');
  }

  function handlePageChange(direction: 'next' | 'prev') {
    if (direction === 'next' && hasMore) {
      currentPage += 1;
      loadData();
    } else if (direction === 'prev' && currentPage > 1) {
      currentPage -= 1;
      loadData();
    }
  }

  $effect(() => {
    const _k = kegiatan;
    currentPage = 1;
    loadData();
  });

  onMount(() => {
    loadData();
  });
</script>

<div class="space-y-6 select-none font-sans max-w-[1600px] mx-auto p-1">
  <!-- Top bar message -->
  {#if loading}
    <div class="bg-slate-50 border border-slate-100 rounded-2xl px-5 py-3.5 flex items-center gap-2 text-xs font-semibold text-slate-500 animate-pulse transition-all">
      <Loader2 class="w-4 h-4 animate-spin text-[#00a294]" />
      Memuat data monitoring kegiatan...
    </div>
  {/if}

  <!-- Header Habit Card (Premium Design) -->
  <div class="bg-white rounded-3xl border border-slate-100 shadow-[0_8px_30px_rgb(0,0,0,0.012)] p-6 flex items-center gap-4.5 relative overflow-hidden transition-all duration-300 hover:scale-[1.002]">
    <!-- Icon matching active habit -->
    <div class="p-4 rounded-2xl shadow-xxs
      {kegiatan === 'bangun-pagi' ? 'bg-amber-50 text-amber-500 border border-amber-100' : ''}
      {kegiatan === 'beribadah' ? 'bg-emerald-50 text-emerald-500 border border-emerald-100' : ''}
      {kegiatan === 'makan-sehat' ? 'bg-rose-50 text-rose-500 border border-rose-100' : ''}
      {kegiatan === 'olahraga' ? 'bg-cyan-50 text-cyan-500 border border-cyan-100' : ''}
      {kegiatan === 'belajar' ? 'bg-violet-50 text-violet-500 border border-violet-100' : ''}
      {kegiatan === 'bermasyarakat' ? 'bg-blue-50 text-blue-500 border border-blue-100' : ''}
      {kegiatan === 'tidur-cukup' ? 'bg-indigo-50 text-indigo-500 border border-indigo-100' : ''}
    ">
      {#if kegiatan === 'bangun-pagi'}
        <Sun class="w-7 h-7" />
      {:else if kegiatan === 'beribadah'}
        <Heart class="w-7 h-7" />
      {:else if kegiatan === 'makan-sehat'}
        <Flame class="w-7 h-7" />
      {:else if kegiatan === 'olahraga'}
        <Compass class="w-7 h-7" />
      {:else if kegiatan === 'belajar'}
        <BookOpen class="w-7 h-7" />
      {:else if kegiatan === 'bermasyarakat'}
        <Users2 class="w-7 h-7" />
      {:else}
        <Moon class="w-7 h-7" />
      {/if}
    </div>

    <div>
      <h2 class="text-lg font-black text-slate-700">{meta.title}</h2>
      <p class="text-xs text-slate-400 font-bold mt-0.5">{meta.subtitle}</p>
    </div>
  </div>

  <!-- Filter & Search Toolbar (Mockup Styled) -->
  <div class="bg-white rounded-3xl border border-slate-100 shadow-[0_8px_30px_rgb(0,0,0,0.012)] p-6">
    <div class="flex flex-col lg:flex-row items-center justify-between gap-4 mb-6">
      <div class="flex flex-col md:flex-row items-center gap-3 w-full lg:w-auto">
        <!-- 1. Search Bar -->
        <div class="relative w-full md:w-[320px]">
          <Search class="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
          <input
            type="text"
            placeholder="Cari siswa berdasarkan nama atau NIS..."
            bind:value={searchQuery}
            oninput={handleFilterChange}
            class="w-full pl-10 pr-4 py-2.5 bg-slate-50 border border-slate-100 hover:bg-slate-100/30 focus:border-[#00a294] focus:bg-white text-slate-700 text-xs font-semibold rounded-xl outline-none transition-all"
          />
        </div>

        <!-- 3. Class selector removed -->
      </div>

      <!-- Total & Refresh info -->
      <div class="flex items-center gap-4 w-full lg:w-auto justify-between lg:justify-end">
        <span class="inline-block bg-slate-100 text-slate-600 font-extrabold text-[10px] uppercase tracking-wider py-1.5 px-3 rounded-full">
          Total: {total} Siswa
        </span>

        <button
          onclick={handleRefresh}
          disabled={loading}
          class="flex items-center gap-2 bg-slate-50 hover:bg-slate-100 border border-slate-200/50 text-slate-600 font-bold text-xs py-2 px-4 rounded-xl transition-all cursor-pointer select-none"
        >
          <RefreshCw class="w-3.5 h-3.5 {loading ? 'animate-spin' : ''}" />
          Refresh
        </button>
      </div>
    </div>

    <!-- Dynamic Preview Data Table -->
    {#if loading}
      <!-- Skeleton layout -->
      <div class="space-y-3">
        <div class="h-9 bg-slate-50 rounded-xl border border-slate-100 animate-pulse"></div>
        {#each Array(5) as _}
          <div class="h-14 bg-slate-50/40 rounded-xl border border-slate-100/50 animate-pulse"></div>
        {/each}
      </div>
    {:else}
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100">
              <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">NIS</th>
              <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Nama Siswa</th>
              <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider text-center">Status Jurnal</th>
            </tr>
          </thead>
          <tbody>
            {#if items.length === 0}
              <tr>
                <td colspan="3" class="py-24 text-center">
                  <div class="flex flex-col items-center justify-center text-slate-400 gap-3">
                    <div class="p-3 bg-slate-50 border border-slate-100 rounded-2xl">
                      <Sun class="w-8 h-8 text-slate-300" />
                    </div>
                    <span class="text-xs font-semibold">Tidak ada data siswa ditemukan.</span>
                  </div>
                </td>
              </tr>
            {:else}
              {#each items as item}
                {@const isDone = item.totalDone > 0}
                <tr 
                  onclick={() => openDetail(item)}
                  class="border-b border-slate-50/60 hover:bg-slate-50/50 transition-all duration-200 cursor-pointer"
                >
                  <td class="py-4.5 px-4 text-xs font-semibold text-[#00a294] font-mono hover:underline">{item.nis}</td>
                  <td class="py-4.5 px-4 text-xs font-black text-slate-700">
                    <span class="hover:text-[#00a294] hover:underline transition-colors">{item.namaSiswa}</span>
                  </td>
                  
                  <!-- Status badging -->
                  <td class="py-4.5 px-4 text-center">
                    <span class="inline-block px-3 py-1 text-[9px] font-extrabold rounded-full uppercase tracking-wider
                      {isDone ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' : 'bg-rose-50 text-rose-500 border border-rose-100'}
                    ">
                      {isDone ? `${item.totalDone}/7 Selesai` : 'Belum Mengisi'}
                    </span>
                  </td>
                </tr>
              {/each}
            {/if}
          </tbody>
        </table>
      </div>

      <!-- Pagination Footer Controls -->
      {#if total > limit}
        <div class="flex items-center justify-between border-t border-slate-100 pt-5 mt-5">
          <span class="text-[10px] text-slate-400 font-bold">
            Halaman {currentPage} (Menampilkan {items.length} dari {total} siswa)
          </span>

          <div class="flex items-center gap-2">
            <button
              onclick={() => handlePageChange('prev')}
              disabled={currentPage === 1}
              class="p-2 bg-slate-50 hover:bg-slate-100 disabled:opacity-40 text-slate-600 rounded-lg border border-slate-200/50 cursor-pointer select-none transition-all"
            >
              <ArrowLeft class="w-4 h-4" />
            </button>
            <button
              onclick={() => handlePageChange('next')}
              disabled={!hasMore}
              class="p-2 bg-slate-50 hover:bg-slate-100 disabled:opacity-40 text-slate-600 rounded-lg border border-slate-200/50 cursor-pointer select-none transition-all"
            >
              <ArrowRight class="w-4 h-4" />
            </button>
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>

<!-- Modal Detail G7 (Premium Glassmorphic Popup) -->
{#if isModalOpen && selectedStudent}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm transition-all duration-300">
    <div class="bg-white w-full max-w-3xl rounded-3xl border border-slate-100 shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      <!-- Modal Header -->
      <div class="bg-slate-50/50 border-b border-slate-100 p-6 flex justify-between items-start">
        <div>
          <span class="text-[9px] font-black uppercase tracking-wider text-slate-400 bg-slate-100 py-1 px-2.5 rounded-md">
            Rekap Jurnal &bull; {meta.title}
          </span>
          <h3 class="text-base font-black text-slate-700 mt-2">{selectedStudent.namaSiswa}</h3>
          <p class="text-xs text-slate-400 font-bold mt-1">
            NIS: {selectedStudent.nis} &bull; Kelas: {selectedStudent.kelas}
          </p>
        </div>
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button 
          onclick={closeModal}
          class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-all cursor-pointer"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Modal Body (History Table) -->
      <Content {loadingHistory} {historyItems} {meta} {kegiatan} />

      <!-- Modal Footer -->
      <div class="bg-slate-50/50 border-t border-slate-100 p-5 flex justify-end">
        <button 
          onclick={closeModal}
          class="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-bold rounded-xl transition-all cursor-pointer"
        >
          Tutup Detail
        </button>
      </div>
    </div>
  </div>
{/if}
