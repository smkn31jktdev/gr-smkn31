<script lang="ts">
  import { Loader2, FileText, ChevronRight, Download } from 'lucide-svelte';
  import { 
    selectedMonth, 
    laporanSearchQuery, 
    G7ReportsList, 
    G7ReportsCount, 
    dashboardLoading, 
    loadDashboardData, 
    downloadAllReportsPDF 
  } from '../../../../logic/adminDashboardLogic';
  import DropdownChoice from '../../../../../shared/components/DropdownChoice.svelte';
  import SearchBar from '../../../../../shared/components/SearchBar.svelte';

  let { onShowDetail, tabSwitcher }: {
    onShowDetail: (nis: string, nama: string, bulanTahun: string) => void,
    tabSwitcher: import('svelte').Snippet
  } = $props();

  // Generate available months dynamically (last 12 months in reverse chronological order)
  const availableMonths = (() => {
    const monthNames = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    const months = [];
    const now = new Date();
    for (let i = 0; i < 12; i++) {
      const d = new Date(now.getFullYear(), now.getMonth() - i, 1);
      const val = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
      const label = `${monthNames[d.getMonth()]} ${d.getFullYear()}`;
      months.push({ value: val, label });
    }
    return months;
  })();

  // Debounced search for Laporan Bulanan
  let debounceTimer: any;
  function handleSearch(e: Event) {
    clearTimeout(debounceTimer);
    const target = e.target as HTMLInputElement;
    laporanSearchQuery.set(target.value);
    debounceTimer = setTimeout(() => {
      loadDashboardData();
    }, 300);
  }
</script>

<!-- Tabs Navigation Header -->
<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-100 pb-4">
  {@render tabSwitcher()}
  
  <div class="flex items-center gap-2 shrink-0">
    <!-- Bulk Download for Monthly -->
    <button 
      onclick={() => downloadAllReportsPDF($selectedMonth, $G7ReportsCount)}
      class="flex items-center gap-1 px-3 py-1.5 bg-slate-50 hover:bg-slate-100 text-slate-600 border border-slate-100 rounded-lg font-bold text-[10px] transition-all cursor-pointer border-none"
    >
      <Download class="w-3 h-3" />
      Download Semua ({$G7ReportsCount})
    </button>
    
    <!-- Month Dropdown Selector -->
    <div class="w-28 text-left">
      <DropdownChoice
        options={availableMonths.map(m => ({ value: m.value, label: m.label }))}
        value={$selectedMonth}
        onchange={(val) => {
          selectedMonth.set(val);
          loadDashboardData();
        }}
        placeholder="Pilih Bulan"
      />
    </div>
  </div>
</div>

<!-- Search Input -->
<SearchBar 
  bind:value={$laporanSearchQuery}
  placeholder="Cari siswa di laporan bulanan..."
  oninput={handleSearch}
  class="my-4 shrink-0"
/>


<!-- Scrollable List -->
{#if $dashboardLoading}
  <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
    <Loader2 class="w-6 h-6 animate-spin text-slate-400 mb-2" />
    <p class="text-[11px] font-semibold">Memuat data...</p>
  </div>
{:else if $G7ReportsList.length === 0}
  <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
    <FileText class="w-8 h-8 text-slate-300 mb-2" />
    <p class="text-[11px] font-bold">Tidak ada laporan bulan ini</p>
  </div>
{:else}
  <div class="flex-1 overflow-y-auto max-h-[380px] pr-1 space-y-2 custom-scrollbar">
    {#each $G7ReportsList as report}
      <button 
        onclick={() => onShowDetail(report.nis, report.namaSiswa, report.bulanTahun)}
        class="w-full flex items-center justify-between p-3 bg-white hover:bg-slate-50/50 border border-slate-100 hover:border-slate-200/80 rounded-xl transition-all group cursor-pointer text-left"
      >
        <div class="flex items-center gap-3">
          <!-- Standardized Initials Badge -->
          <div class="w-9 h-9 rounded-lg bg-slate-50 border border-slate-100 text-slate-500 flex items-center justify-center font-bold text-xs shrink-0 select-none">
            {report.namaSiswa.charAt(0).toUpperCase()}
          </div>
          
          <div class="space-y-0.5">
            <h4 class="text-xs font-bold text-slate-700 uppercase group-hover:text-slate-900 transition-colors">
              {report.namaSiswa}
            </h4>
            <p class="text-[10px] font-medium text-slate-400">
              {report.kelas} <span class="mx-1 text-slate-250">•</span> NIS: {report.nis}
            </p>
          </div>
        </div>
        
        <div class="flex items-center gap-2">
          {#if report.predikat}
            <span class="px-2 py-0.5 rounded-md text-[9px] font-bold uppercase bg-slate-50 text-slate-500 border border-slate-100">
              {report.predikat}
            </span>
          {/if}
          <ChevronRight class="w-3.5 h-3.5 text-slate-300 group-hover:text-slate-400 transition-all" />
        </div>
      </button>
    {/each}
  </div>
{/if}

<style>
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

