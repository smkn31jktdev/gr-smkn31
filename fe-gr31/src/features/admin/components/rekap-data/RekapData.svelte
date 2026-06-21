<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Download, 
    RefreshCw, 
    FileSpreadsheet, 
    Database, 
    Calendar, 
    ChevronRight, 
    Loader2,
    Users,
    TrendingUp,
    AlertCircle,
    FileText
  } from 'lucide-svelte';
  
  import { 
    loadRekapSummary, 
    loadRekapPreview, 
    fetchAllRekapForExport, 
    handleExportG7RekapCSV 
  } from '../../logic/rekapDataLogic';
  import { getG7RekapKelas } from '../../logic/adminLogic';
  import type { G7Rekap, G7RekapKelasLengkap, G7RekapSiswaItem } from '../../types/admin.types';
  import { addToast } from '../../../../stores/uiStore';
  import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';
  import { apiRequest } from '../../../../api/client';

  // State
  let selectedBulan = $state('2026-06');
  let selectedKelas = $state('');
  
  let terpilihCount = $state(0);
  let totalArsipCount = $state(0);
  let previewItems = $state<G7Rekap[]>([]);
  let loading = $state(false);
  let exporting = $state(false);
  let ramadanPeriods = $state<{ start_date: string; end_date: string }[]>([]);

  let isRamadan = $derived.by(() => {
    const targetBulan = selectedBulan; // YYYY-MM
    return ramadanPeriods.some(p => {
      const startMonth = p.start_date.substring(0, 7);
      const endMonth = p.end_date.substring(0, 7);
      return targetBulan === startMonth || targetBulan === endMonth;
    });
  });

  // Rekap kelas state (endpoint baru)
  let rekapKelasLoading = $state(false);
  let rekapKelasData = $state<G7RekapKelasLengkap | null>(null);

  // Active tab
  let activeView = $state<'export' | 'rekap-kelas'>('rekap-kelas');

  function getMonthOptions() {
    const options = [];
    const monthNames = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    const anchorYear = 2026;
    const anchorMonth = 5;
    for (let i = -12; i <= 3; i++) {
      const d = new Date(anchorYear, anchorMonth + i, 1);
      const val = d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0');
      const label = `${monthNames[d.getMonth()]} ${d.getFullYear()}`;
      options.push({ val, label });
    }
    return options.reverse();
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

  async function loadData() {
    loading = true;
    const summary = await loadRekapSummary(selectedBulan, selectedKelas);
    terpilihCount = summary.terpilihCount;
    totalArsipCount = summary.totalArsipCount;
    previewItems = await loadRekapPreview(selectedBulan, selectedKelas);
    loading = false;
  }

  async function loadRekapKelas() {
    rekapKelasLoading = true;
    const res = await getG7RekapKelas(selectedBulan, selectedKelas);
    rekapKelasData = res;
    rekapKelasLoading = false;
  }

  function handleRefresh() {
    if (activeView === 'rekap-kelas') {
      loadRekapKelas();
    } else {
      loadData();
    }
    addToast('Data rekap berhasil diperbarui', 'success');
  }

  async function handleExport() {
    if (exporting) return;
    exporting = true;
    try {
      const items = await fetchAllRekapForExport(selectedBulan, selectedKelas);
      if (items.length === 0) {
        addToast('Tidak ada data rekap pada periode terpilih', 'warning');
      } else {
        handleExportG7RekapCSV(items, selectedBulan, selectedKelas, isRamadan);
        addToast('Laporan G7 berhasil diexport', 'success');
      }
    } catch (e) {
      console.error(e);
      addToast('Gagal melakukan export CSV', 'error');
    } finally {
      exporting = false;
    }
  }

  function getPredikatClass(predikat: string) {
    switch (predikat) {
      case 'Istimewa': return 'bg-emerald-50 text-emerald-600 border border-emerald-100';
      case 'Sangat Baik': return 'bg-blue-50 text-blue-600 border border-blue-100';
      case 'Baik': return 'bg-cyan-50 text-cyan-600 border border-cyan-100';
      case 'Cukup': return 'bg-amber-50 text-amber-600 border border-amber-100';
      case 'Kurang': return 'bg-rose-50 text-rose-600 border border-rose-100';
      default: return 'bg-slate-50 text-slate-500 border border-slate-100';
    }
  }

  $effect(() => {
    if (activeView === 'rekap-kelas') {
      loadRekapKelas();
    } else {
      loadData();
    }
  });

  onMount(async () => {
    loadRekapKelas();
    try {
      const res = await apiRequest<{ ramadan: { start_date: string; end_date: string }[] }>('/v1/puasa/calendar');
      if (res.data && res.data.ramadan) {
        ramadanPeriods = res.data.ramadan;
      }
    } catch (e) {
      console.error('Error fetching calendar:', e);
    }
  });
</script>

<div class="space-y-6 select-none font-sans max-w-[1600px] mx-auto p-1">
  <!-- Header -->
  <div class="flex items-center justify-between flex-wrap gap-4">
    <div>
      <h2 class="text-xl font-extrabold tracking-tight text-slate-800 font-display uppercase">Rekap Data G7</h2>
      <p class="text-xs text-slate-400 mt-0.5 font-medium">Tinjau rekap bulanan 7 kebiasaan dan download laporan spreadsheet kelas</p>
    </div>
    <div class="flex items-center gap-3">
      <button
        onclick={handleExport}
        disabled={exporting}
        class="flex items-center gap-2 bg-[#00a294] hover:bg-[#008c80] disabled:bg-slate-200 text-white font-bold text-xs py-2.5 px-4 rounded-xl shadow-xs transition-all cursor-pointer border-none"
      >
        {#if exporting}
          <Loader2 class="w-4 h-4 animate-spin" />
        {:else}
          <Download class="w-4 h-4" />
        {/if}
        Export CSV
      </button>
      <button
        onclick={handleRefresh}
        disabled={loading || rekapKelasLoading}
        class="flex items-center gap-2 bg-slate-50 hover:bg-slate-100 border border-slate-200/50 text-slate-600 font-bold text-xs py-2.5 px-4 rounded-xl transition-all cursor-pointer border-none"
      >
        <RefreshCw class="w-4 h-4 {(loading || rekapKelasLoading) ? 'animate-spin' : ''}" />
        Segarkan
      </button>
    </div>
  </div>

  <!-- Period & Class Filters (shared) -->
  <div class="bg-white rounded-2xl border border-slate-100 shadow-xs p-5">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div class="flex flex-col gap-1.5 text-left">
        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest">BULAN LAPORAN</label>
        <DropdownChoice
          options={monthsList.map(m => ({ value: m.val, label: m.label }))}
          bind:value={selectedBulan}
          placeholder="Pilih Bulan"
        />
      </div>
      <!-- Tab switcher -->
      <div class="flex flex-col gap-1.5">
        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest">TAMPILAN</label>
        <div class="flex border border-slate-200 rounded-xl p-1 bg-slate-50/50">
          <button
            onclick={() => activeView = 'rekap-kelas'}
            class="flex-1 py-2 rounded-lg text-xs font-bold transition-all border-none cursor-pointer {activeView === 'rekap-kelas' ? 'bg-[#00a294] text-white shadow-xs' : 'bg-transparent text-slate-450 hover:text-slate-700'}"
          >
            Rekap Kelas
          </button>
          <button
            onclick={() => activeView = 'export'}
            class="flex-1 py-2 rounded-lg text-xs font-bold transition-all border-none cursor-pointer {activeView === 'export' ? 'bg-[#00a294] text-white shadow-xs' : 'bg-transparent text-slate-450 hover:text-slate-700'}"
          >
            Export Data
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- â”€â”€ VIEW 1: REKAP KELAS (endpoint baru g7/rekap-kelas) â”€â”€ -->
  {#if activeView === 'rekap-kelas'}
    {#if rekapKelasLoading}
      <div class="bg-white rounded-2xl border border-slate-100 p-16 flex flex-col items-center justify-center gap-3">
        <Loader2 class="w-8 h-8 animate-spin text-slate-300" />
        <p class="text-xs text-slate-400 font-semibold">Memuat rekap nilai kelas...</p>
      </div>
    {:else if rekapKelasData}
      <!-- Tabel lengkap per siswa (roster-join) -->
      <div class="bg-white rounded-3xl border border-slate-100 shadow-[0_8px_30px_rgb(0,0,0,0.015)] overflow-hidden">
        <div class="p-6 border-b border-slate-100 flex items-center gap-3">
          <div class="p-2 bg-slate-50 border border-slate-100 rounded-xl text-slate-500">
            <FileText class="w-5 h-5" />
          </div>
          <h3 class="text-sm font-black text-slate-800 font-display">Preview Data</h3>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="border-b border-slate-100">
                <th class="py-4 px-6 text-[10px] font-black text-slate-400 uppercase tracking-wider">NIS</th>
                <th class="py-4 px-6 text-[10px] font-black text-slate-400 uppercase tracking-wider">Nama Siswa</th>
                <th class="py-4 px-6 text-[10px] font-black text-slate-400 uppercase tracking-wider">Kelas</th>
                <th class="py-4 px-6 text-[10px] font-black text-slate-400 uppercase tracking-wider">Periode</th>
              </tr>
            </thead>
            <tbody>
              {#each rekapKelasData.siswa as item}
                <tr class="border-b border-slate-50/60 hover:bg-slate-50/30 transition-colors">
                  <td class="py-4 px-6 text-xs font-semibold text-slate-500 font-mono">{item.nis}</td>
                  <td class="py-4 px-6 text-xs font-bold text-slate-700">
                    <a href="/admin/g7/{item.nis}/{rekapKelasData?.bulanTahun}" class="font-bold text-slate-800 uppercase tracking-wide hover:text-[#00a294] hover:underline">
                      {item.namaSiswa}
                    </a>
                  </td>
                  <td class="py-4 px-6 text-xs font-semibold text-slate-500">{item.kelas || rekapKelasData.kelas}</td>
                  <td class="py-4 px-6 text-xs">
                    <span class="px-2.5 py-1 rounded-md text-[10px] font-bold bg-blue-50 text-blue-600 border border-blue-100">
                      {formatMonthYear(rekapKelasData.bulanTahun)}
                    </span>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    {:else}
      <div class="bg-white rounded-2xl border border-slate-100 p-12 flex flex-col items-center justify-center gap-3 text-center">
        <AlertCircle class="w-8 h-8 text-slate-300" />
        <p class="text-xs text-slate-400 font-semibold">Pilih kelas dan bulan untuk melihat rekap nilai G7.</p>
      </div>
    {/if}

  <!-- â”€â”€ VIEW 2: EXPORT DATA â”€â”€ -->
  {:else}
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-stretch">
      <div class="lg:col-span-2 bg-white rounded-3xl border border-slate-100 shadow-[0_8px_30px_rgb(0,0,0,0.015)] p-6 flex flex-col justify-between">
        <div>
          <div class="flex items-center gap-2.5 mb-5 border-b border-slate-50 pb-4">
            <div class="p-2 bg-emerald-50 rounded-xl text-emerald-600">
              <Calendar class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-sm font-black text-slate-700">Export CSV</h3>
              <p class="text-[10px] text-slate-400 font-medium">Download data rekap G7 semua siswa dalam format spreadsheet</p>
            </div>
          </div>
        </div>
        <div class="flex items-center gap-3 mt-4">
          <button
            onclick={handleExport}
            disabled={exporting}
            class="flex-1 md:flex-initial flex items-center justify-center gap-2 bg-[#00a294] hover:bg-[#008c80] disabled:bg-slate-200 text-white font-bold text-xs py-3 px-6 rounded-xl shadow-md hover:shadow-lg transition-all cursor-pointer select-none"
          >
            {#if exporting}
              <Loader2 class="w-4 h-4 animate-spin" />
              Memproses...
            {:else}
              <Download class="w-4 h-4" />
              Download CSV
            {/if}
          </button>
        </div>
      </div>

      <div class="bg-white rounded-3xl border border-slate-100 shadow-[0_8px_30px_rgb(0,0,0,0.015)] p-6 flex flex-col gap-5">
        <div class="flex items-center gap-2 mb-2 border-b border-slate-50 pb-4">
          <div class="p-2 bg-blue-50 rounded-xl text-blue-600">
            <Database class="w-5 h-5" />
          </div>
          <div>
            <h3 class="text-sm font-black text-slate-700">Ringkasan Data</h3>
            <p class="text-[10px] text-slate-400 font-medium">Metrik cakupan laporan di database</p>
          </div>
        </div>
        <div class="bg-slate-50/60 border border-slate-100 p-5 rounded-2xl relative overflow-hidden">
          <div class="flex flex-col">
            <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">TERPILIH</span>
            {#if loading}
              <div class="h-8 w-24 bg-slate-200 animate-pulse rounded-lg mt-2 mb-1"></div>
            {:else}
              <span class="text-2xl font-black text-slate-700 mt-1">{terpilihCount}</span>
            {/if}
            <span class="text-[10px] text-slate-400 font-bold mt-1">Data siswa bulan {formatMonthYear(selectedBulan)}</span>
          </div>
          <ChevronRight class="absolute right-5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-300" />
        </div>
        <div class="bg-slate-50/60 border border-slate-100 p-5 rounded-2xl relative overflow-hidden">
          <div class="flex flex-col">
            <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest">TOTAL ARSIP</span>
            {#if loading}
              <div class="h-8 w-24 bg-slate-200 animate-pulse rounded-lg mt-2 mb-1"></div>
            {:else}
              <span class="text-2xl font-black text-slate-700 mt-1">{totalArsipCount}</span>
            {/if}
            <span class="text-[10px] text-slate-400 font-bold mt-1">Total semua periode</span>
          </div>
          <ChevronRight class="absolute right-5 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-300" />
        </div>
      </div>
    </div>

    <!-- Preview Table -->
    <div class="bg-white rounded-3xl border border-slate-100 shadow-[0_8px_30px_rgb(0,0,0,0.015)] p-6">
      <div class="flex items-center gap-2.5 mb-5 border-b border-slate-50 pb-4">
        <div class="p-2 bg-violet-50 rounded-xl text-violet-600">
          <FileSpreadsheet class="w-5 h-5" />
        </div>
        <div>
          <h3 class="text-sm font-black text-slate-700">Preview Data (5 Teratas)</h3>
          <p class="text-[10px] text-slate-400 font-medium">Tampilan cepat 5 data siswa untuk validasi sebelum export</p>
        </div>
      </div>
      {#if loading}
        <div class="space-y-3">
          {#each Array(5) as _}
            <div class="h-12 bg-slate-50/40 rounded-xl border border-slate-100/50 animate-pulse"></div>
          {/each}
        </div>
      {:else}
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="border-b border-slate-100">
                <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">NIS</th>
                <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Nama Siswa</th>
                <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Kelas</th>
                <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Periode</th>
                <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Nilai Akhir</th>
                <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Predikat</th>
              </tr>
            </thead>
            <tbody>
              {#if previewItems.length === 0}
                <tr>
                  <td colspan="6" class="py-12 text-center text-xs text-slate-400 font-medium">
                    Tidak ada data rekap untuk periode ini.
                  </td>
                </tr>
              {:else}
                {#each previewItems as item}
                  <tr class="border-b border-slate-50/60 hover:bg-slate-50/30 transition-colors">
                    <td class="py-4 px-4 text-xs font-semibold text-slate-500 font-mono">{item.nis}</td>
                    <td class="py-4 px-4 text-xs font-bold text-slate-700">{item.namaSiswa}</td>
                    <td class="py-4 px-4 text-xs font-semibold text-slate-500">{item.kelas}</td>
                    <td class="py-4 px-4 text-xs font-semibold text-slate-400">{formatMonthYear(item.bulanTahun)}</td>
                    <td class="py-4 px-4 text-xs font-black text-slate-700">{item.nilaiAkhir.toFixed(2)}</td>
                    <td class="py-4 px-4">
                      <span class="inline-block px-2.5 py-1 text-[9px] font-extrabold rounded-full uppercase tracking-wider {getPredikatClass(item.predikat)}">
                        {item.predikat}
                      </span>
                    </td>
                  </tr>
                {/each}
              {/if}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  {/if}
</div>

