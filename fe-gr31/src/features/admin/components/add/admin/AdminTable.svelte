<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    UserPlus, 
    FolderDown, 
    RefreshCw,
    Users,
    ShieldCheck,
    GraduationCap,
    UserCheck
  } from 'lucide-svelte';
  
  import {
    adminsList,
    adminsLoading,
    loadAdmins
  } from '../../../logic/adminAdminsLogic';
  import Manual from './add/manual/Manual.svelte';
  import Sheets from './add/sheets/Sheets.svelte';
  import Table from './table/Table.svelte';

  let activeTab = $state<'manual' | 'sheets'>('manual');

  // Derived quick stats
  let totalAdmins = $derived($adminsList.length);
  let superAdminCount = $derived($adminsList.filter(a => a.role === 'super_admin').length);
  let guruBkCount = $derived($adminsList.filter(a => a.role === 'guru_bk' || a.role === 'admin_bk').length);
  let waliKelasCount = $derived($adminsList.filter(a => a.isWalas || a.role === 'walas' || a.role === 'guru_wali').length);
  let staffAdminCount = $derived($adminsList.filter(a => a.role === 'admin' || a.role === 'piket').length);

  onMount(() => {
    loadAdmins();
  });
</script>

<div class="space-y-6 select-none font-sans pb-10">
  
  <!-- Quick Stats Grid -->
  <div class="grid grid-cols-2 lg:grid-cols-5 gap-4">
    <!-- Stat 1: Total Admin/Guru -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
      <div class="space-y-1 text-left">
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Total Guru</span>
        <span class="text-lg font-extrabold text-slate-800">{totalAdmins}</span>
      </div>
      <div class="w-9 h-9 rounded-xl bg-slate-50 flex items-center justify-center text-slate-500 border border-slate-100">
        <Users class="w-4.5 h-4.5" />
      </div>
    </div>

    <!-- Stat 2: Super Admin -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
      <div class="space-y-1 text-left">
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Super Admin</span>
        <span class="text-lg font-extrabold text-slate-800">{superAdminCount}</span>
      </div>
      <div class="w-9 h-9 rounded-xl bg-slate-900 flex items-center justify-center text-white border border-slate-950">
        <ShieldCheck class="w-4.5 h-4.5" />
      </div>
    </div>

    <!-- Stat 3: Guru BK -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
      <div class="space-y-1 text-left">
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Guru BK</span>
        <span class="text-lg font-extrabold text-slate-800">{guruBkCount}</span>
      </div>
      <div class="w-9 h-9 rounded-xl bg-indigo-50/50 flex items-center justify-center text-indigo-600 border border-indigo-100/60">
        <Users class="w-4.5 h-4.5" />
      </div>
    </div>

    <!-- Stat 4: Wali Kelas -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
      <div class="space-y-1 text-left">
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Wali Kelas</span>
        <span class="text-lg font-extrabold text-slate-800">{waliKelasCount}</span>
      </div>
      <div class="w-9 h-9 rounded-xl bg-teal-50/50 flex items-center justify-center text-teal-655 border border-teal-100/60">
        <GraduationCap class="w-4.5 h-4.5" />
      </div>
    </div>

    <!-- Stat 5: Staf Admin -->
    <div class="bg-white p-4 rounded-2xl border border-slate-100/80 shadow-xs flex items-center justify-between hover:border-slate-200 transition-all duration-300">
      <div class="space-y-1 text-left">
        <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">Staf / Guru Non-Walas</span>
        <span class="text-lg font-extrabold text-slate-800">{staffAdminCount}</span>
      </div>
      <div class="w-9 h-9 rounded-xl bg-amber-50/50 flex items-center justify-center text-amber-655 border border-amber-100/60 font-mono">
        <UserCheck class="w-4.5 h-4.5" />
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
      onclick={() => { loadAdmins(); }}
      disabled={$adminsLoading}
      class="p-2 text-slate-400 hover:text-slate-655 hover:bg-slate-50 border border-transparent hover:border-slate-100 rounded-xl transition-all cursor-pointer flex items-center justify-center"
      title="Segarkan data"
    >
      <RefreshCw class="w-4 h-4 {$adminsLoading ? 'animate-spin text-slate-500' : ''}" />
    </button>
  </div>

  {#if activeTab === 'manual'}
    <Manual />
  {:else}
    <Sheets />
  {/if}

  <Table />

</div>
