<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Users,
    School,
    RefreshCw,
    UserPlus,
    FolderDown
  } from 'lucide-svelte';
  import { listStudents } from '../../../logic/adminLogic';
  import Manual from './add/manual/Manual.svelte';
  import Sheets from './add/sheets/Sheets.svelte';
  import Table from './table/Table.svelte';
  import ClassStats from './class/ClassStats.svelte';
  import type { Siswa } from '../../../../auth/types/auth.types';

  let allStudents = $state<Siswa[]>([]);
  let loading = $state(false);
  let activeTab = $state<'manual' | 'sheets'>('manual');
  let filterQuery = $state('');
  let selectedKelas = $state('');

  // Pagination states
  let page = $state(1);
  let limit = $state(50);
  let totalStudentsCount = $state(0);
  let hasMore = $state(false);

  // Derived quick stats
  let totalStudents = $derived(totalStudentsCount);
  let uniqueKelasCount = $derived(new Set(allStudents.map(s => s.kelas?.trim().toUpperCase()).filter(Boolean)).size);

  // Keep track of last search filters to detect changes and reset page
  let lastQuery = '';
  let lastKelas = '';
  let lastLimit = 50;

  async function loadData() {
    loading = true;
    try {
      const limitParam = limit === -1 ? 'all' : limit;
      const res = await listStudents(page, limitParam as any, filterQuery, selectedKelas);
      allStudents = res.items || [];
      totalStudentsCount = res.total || allStudents.length;
      hasMore = res.hasMore || false;
    } catch (err) {
      console.error('Error loading students:', err);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    if (filterQuery !== lastQuery || selectedKelas !== lastKelas || limit !== lastLimit) {
      lastQuery = filterQuery;
      lastKelas = selectedKelas;
      lastLimit = limit;
      page = 1;
      return;
    }
    loadData();
  });
</script>

<div class="space-y-6 select-none font-sans pb-10">

  <!-- Quick Stats Grid -->
  <div class="grid grid-cols-2 lg:grid-cols-5 gap-4">
    <!-- Stat 1: Total Siswa -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
      <div class="space-y-1 text-left">
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Total Siswa</span>
        <span class="text-lg font-extrabold text-slate-800">{totalStudents}</span>
      </div>
      <div class="w-9 h-9 rounded-xl bg-slate-50 flex items-center justify-center text-slate-500 border border-slate-100">
        <Users class="w-4.5 h-4.5" />
      </div>
    </div>

    <ClassStats students={allStudents} />

    <!-- Stat 5: Banyak Kelas -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
      <div class="space-y-1 text-left">
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Banyak Kelas</span>
        <span class="text-lg font-extrabold text-slate-800">{uniqueKelasCount}</span>
      </div>
      <div class="w-9 h-9 rounded-xl bg-slate-50 flex items-center justify-center text-slate-500 border border-slate-100">
        <School class="w-4.5 h-4.5" />
      </div>
    </div>
  </div>

  <!-- Tab selectors & Refresh Button -->
  <div class="flex items-center justify-between gap-4 border-b border-slate-100 pb-2">
    <div class="flex gap-1.5 p-1 bg-slate-100/60 rounded-xl w-fit shrink-0">
      <button 
        onclick={() => activeTab = 'manual'}
        class="px-4 py-1.5 text-xs font-bold rounded-lg transition-all duration-200 flex items-center gap-2 border-none cursor-pointer {activeTab === 'manual' ? 'bg-white text-slate-800 shadow-xs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
      >
        <UserPlus class="w-3.5 h-3.5" />
        Registrasi Manual
      </button>
      <button 
        onclick={() => activeTab = 'sheets'}
        class="px-4 py-1.5 text-xs font-bold rounded-lg transition-all duration-200 flex items-center gap-2 border-none cursor-pointer {activeTab === 'sheets' ? 'bg-white text-slate-800 shadow-xs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
      >
        <FolderDown class="w-3.5 h-3.5" />
        Impor Google Sheet
      </button>
    </div>

    <button
      onclick={() => { loadData(); }}
      disabled={loading}
      class="p-2 text-slate-400 hover:text-slate-650 hover:bg-slate-50 border border-transparent hover:border-slate-100 rounded-xl transition-all cursor-pointer flex items-center justify-center"
      title="Segarkan data siswa"
    >
      <RefreshCw class="w-4 h-4 {loading ? 'animate-spin text-slate-500' : ''}" />
    </button>
  </div>

  {#if activeTab === 'manual'}
    <Manual onSuccess={loadData} />
  {:else}
    <Sheets onSuccess={loadData} />
  {/if}

  <Table
    students={allStudents}
    {loading}
    onRefresh={loadData}
    bind:filterQuery
    bind:selectedKelas
    {totalStudentsCount}
    {hasMore}
    bind:page
    bind:limit
  />

</div>
