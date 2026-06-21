<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    User, 
    Users, 
    Mail, 
    Loader2,
    RefreshCw
  } from 'lucide-svelte';
  
  import {
    dashboardLoading,
    activeAdminName,
    activeAdminEmail,
    totalSiswaCount,
    studentSearchQuery,
    studentsList,
    totalStudentsFiltered,
    loadDashboardData,
    loadSemesterData,
    pendingWalasAduan
  } from '../../logic/adminDashboardLogic';
  
  import { getG7RekapDetail, getG7EvaluateDetail, listKehadiranAdmin } from '../../logic/adminLogic';
  import { getTodayStr } from '../../logic/adminKehadiranSiswaLogic';
  import { addToast } from '../../../../stores/uiStore';
  import { currentUser } from '../../../../stores/authStore';

  // Subcomponents
  import ManualBook from './book/ManualBook.svelte';
  import LaporanBulanan from './student/bulan/LaporanBulanan.svelte';
  import LaporanSemester from './student/semester/LaporanSemester.svelte';
  import G7Modal from './student/g7/G7Modal.svelte';
  import SearchBar from '../../../shared/components/SearchBar.svelte';
  import AttendanceDetailModal from './AttendanceDetailModal.svelte';

  // Detail Modal State
  let selectedDetailStudent = $state<any>(null);
  let detailModalOpen = $state(false);
  let detailLoading = $state(false);
  let detailRekap = $state<any>(null);
  let detailEvaluate = $state<any>(null);
  let detailType = $state<'bulanan' | 'semester'>('bulanan');

  // Attendance Detail Modal State
  let attendanceModalOpen = $state(false);
  let selectedAttendanceLog = $state<any>(null);

  function openAttendanceDetail(log: any) {
    selectedAttendanceLog = log;
    attendanceModalOpen = true;
  }

  async function showMonthlyDetail(nis: string, nama: string, bulan: string) {
    selectedDetailStudent = { nis, nama, bulan };
    detailType = 'bulanan';
    detailModalOpen = true;
    detailLoading = true;
    detailRekap = null;
    detailEvaluate = null;
    
    try {
      const [rekapRes, evaluateRes] = await Promise.all([
        getG7RekapDetail(nis, bulan),
        getG7EvaluateDetail(nis, bulan)
      ]);
      detailRekap = rekapRes;
      detailEvaluate = evaluateRes;
    } catch (e) {
      console.error(e);
      addToast('Gagal memuat detail rekap G7', 'error');
    } finally {
      detailLoading = false;
    }
  }

  function showSemesterDetail(report: any, semester: string) {
    selectedDetailStudent = { 
      nis: report.nis, 
      nama: report.namaSiswa, 
      semester 
    };
    detailType = 'semester';
    detailModalOpen = true;
    detailRekap = report;
    detailEvaluate = null;
    detailLoading = false;
  }

  let activeTab = $state<'bulanan' | 'semester'>('bulanan');

  function switchTab(tab: 'bulanan' | 'semester') {
    activeTab = tab;
    if (tab === 'bulanan') {
      loadDashboardData();
    } else {
      loadSemesterData();
    }
  }

  // Debounced search for Student List
  let studentDebounceTimer: any;
  function handleStudentSearch(e: Event) {
    clearTimeout(studentDebounceTimer);
    const target = e.target as HTMLInputElement;
    studentSearchQuery.set(target.value);
    studentDebounceTimer = setTimeout(() => {
      loadDashboardData();
    }, 300);
  }

  let isWalas = $derived(
    $currentUser?.role === 'walas' || 
    $currentUser?.role === 'guru_wali' || 
    $currentUser?.is_walas === true || 
    $currentUser?.isWalas === true
  );

  let walasKelas = $derived($currentUser?.kelas || '');

  let attendanceLogs = $state<any[]>([]);
  let attendanceLoading = $state(false);
  let attendanceSearchQuery = $state('');

  async function loadTodayAttendance() {
    if (!isWalas || !walasKelas) return;
    attendanceLoading = true;
    try {
      const res = await listKehadiranAdmin({
        kelas: walasKelas,
        tanggal: getTodayStr()
      }, 1, 100);
      attendanceLogs = res.items || [];
    } catch (e) {
      console.error(e);
      addToast('Gagal memuat absensi harian', 'error');
    } finally {
      attendanceLoading = false;
    }
  }

  let filteredAttendance = $derived(
    attendanceSearchQuery.trim() === ''
      ? attendanceLogs
      : attendanceLogs.filter(item => 
          item.namaSiswa.toLowerCase().includes(attendanceSearchQuery.toLowerCase()) ||
          item.nis.toLowerCase().includes(attendanceSearchQuery.toLowerCase())
        )
  );

  let statsTotal = $derived(attendanceLogs.length);
  let statsHadir = $derived(attendanceLogs.filter(i => i.status === 'hadir').length);
  let statsIzin = $derived(attendanceLogs.filter(i => i.status === 'izin').length);
  let statsSakit = $derived(attendanceLogs.filter(i => i.status === 'sakit').length);
  let statsMagang = $derived(attendanceLogs.filter(i => i.status === 'magang').length);
  let statsAlpa = $derived(attendanceLogs.filter(i => i.status === 'tidak_hadir').length);
  let statsBelum = $derived(attendanceLogs.filter(i => i.status === 'belum').length);

  onMount(() => {
    loadDashboardData();
    loadSemesterData();
    if (isWalas) {
      loadTodayAttendance();
    }
  });
</script>

<div class="space-y-6 select-none font-sans pb-10">
  <!-- Title Header -->
  <div class="space-y-1 text-left">
    <h1 class="text-xl font-bold tracking-tight text-slate-800 uppercase leading-none">
      Dashboard Overview
    </h1>
    <p class="text-xs text-slate-400 font-medium">
      Pantau aktivitas sistem sekolah dan kelola data siswa dengan mudah dalam satu tampilan.
    </p>
  </div>

  <!-- Pending Aduan Ping Notification for Walas -->
  {#if isWalas && $pendingWalasAduan.length > 0}
    <div class="relative overflow-hidden bg-amber-50/80 border border-amber-200/60 rounded-2xl p-4.5 text-left shadow-xs flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 animate-fade-in">
      <!-- Glow background decoration -->
      <div class="absolute -right-10 -top-10 w-32 h-32 bg-amber-200/20 rounded-full blur-2xl pointer-events-none"></div>
      
      <div class="flex items-start gap-3.5 z-10">
        <!-- Pulse Ping Icon -->
        <div class="relative flex items-center justify-center w-10 h-10 rounded-xl bg-amber-100/80 text-amber-600 shrink-0 border border-amber-200/30">
          <span class="absolute top-1 right-1 flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-rose-500 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-rose-500"></span>
          </span>
          <svg class="w-5 h-5 animate-bounce" style="animation-duration: 2.5s;" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
          </svg>
        </div>
        
        <div class="space-y-1">
          <h4 class="text-xs font-extrabold text-amber-900 uppercase tracking-tight">Tindak Lanjut Aduan Siswa</h4>
          <p class="text-[11px] text-amber-700 font-medium leading-relaxed">
            Halo Guru Wali! Terdapat <strong>{$pendingWalasAduan.length} aduan aktif</strong> dari siswa di kelas Anda yang dibuat dalam 7 hari terakhir dan membutuhkan tanggapan Anda segera.
          </p>
        </div>
      </div>
      
      <a 
        href="/admin/chat"
        class="shrink-0 z-10 px-4 py-2 bg-amber-600 hover:bg-amber-700 text-white rounded-xl text-[11px] font-bold shadow-[0_2px_8px_rgba(217,119,6,0.2)] transition-all flex items-center justify-center gap-1.5 cursor-pointer border-none"
      >
        Tanggapi Siswa <span class="text-xs">→</span>
      </a>
    </div>
  {/if}

  <!-- Stats Grid -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
    <!-- Active Admin Card -->
    <div class="flex items-center gap-4 bg-white border border-slate-100 rounded-2xl p-4 shadow-xs">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-400 border border-slate-100 flex items-center justify-center shrink-0">
        <User class="w-4 h-4" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Admin Aktif</p>
        <h3 class="text-xs font-bold text-slate-700 mt-0.5">{$activeAdminName}</h3>
      </div>
    </div>

    <!-- Total Siswa Card -->
    <div class="flex items-center gap-4 bg-white border border-slate-100 rounded-2xl p-4 shadow-xs">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-400 border border-slate-100 flex items-center justify-center shrink-0">
        <Users class="w-4 h-4" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Total Siswa</p>
        <h3 class="text-xs font-bold text-slate-700 mt-0.5">{$totalSiswaCount} Siswa</h3>
      </div>
    </div>

    <!-- Registered Email Card -->
    <div class="flex items-center gap-4 bg-white border border-slate-100 rounded-2xl p-4 shadow-xs">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-400 border border-slate-100 flex items-center justify-center shrink-0">
        <Mail class="w-4 h-4" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Email Terdaftar</p>
        <h3 class="text-xs font-bold text-slate-700 mt-0.5 truncate max-w-[180px]">{$activeAdminEmail}</h3>
      </div>
    </div>
  </div>

  <!-- Manual Book Section -->
  <ManualBook />

  <!-- Two Column Content Panel -->
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
    
    <!-- Left Column: Reports with Tabs (Col span 7) -->
    <div class="lg:col-span-7 bg-white border border-slate-100 rounded-2xl p-5 shadow-xs flex flex-col min-h-[520px]">
      {#snippet tabSwitcherSnippet()}
        <div class="flex border border-slate-100 bg-slate-50/50 p-1 rounded-xl shrink-0">
          <button 
            onclick={() => switchTab('bulanan')}
            class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all border-none cursor-pointer {activeTab === 'bulanan' ? 'bg-white text-slate-800 shadow-xs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
          >
            Rekap Bulanan
          </button>
          <button 
            onclick={() => switchTab('semester')}
            class="px-4 py-1.5 rounded-lg text-xs font-bold transition-all border-none cursor-pointer {activeTab === 'semester' ? 'bg-white text-slate-800 shadow-xs' : 'text-slate-400 hover:text-slate-600 bg-transparent'}"
          >
            Rekap Semester
          </button>
        </div>
      {/snippet}

      {#if activeTab === 'bulanan'}
        <LaporanBulanan 
          onShowDetail={showMonthlyDetail}
          tabSwitcher={tabSwitcherSnippet}
        />
      {:else}
        <LaporanSemester 
          onShowDetail={showSemesterDetail}
          tabSwitcher={tabSwitcherSnippet}
        />
      {/if}
    </div>

    <!-- Right Column: List Siswa (Col span 5) -->
    {#if isWalas}
      <div class="lg:col-span-5 bg-white border border-slate-100 rounded-2xl p-5 shadow-xs flex flex-col min-h-[520px]">
        <!-- Panel Header -->
        <div class="flex items-center justify-between border-b border-slate-100 pb-4">
          <div class="text-left">
            <h2 class="text-xs font-bold text-slate-700 uppercase tracking-tight">Presensi Harian Kelas</h2>
            <p class="text-[10px] text-slate-400 font-medium mt-0.5">Kehadiran hari ini: {walasKelas}</p>
          </div>
          
          <button 
            onclick={loadTodayAttendance}
            disabled={attendanceLoading}
            class="flex items-center justify-center p-1.5 border border-slate-200 rounded-lg bg-white hover:bg-slate-50 text-slate-500 transition-all cursor-pointer shadow-xxs"
            title="Segarkan data"
          >
            <RefreshCw class="w-3.5 h-3.5 {attendanceLoading ? 'animate-spin' : ''}" />
          </button>
        </div>

        <!-- Quick Metrics Grid -->
        <div class="grid grid-cols-3 gap-2 mt-4 text-left">
          <!-- Hadir Card -->
          <div class="bg-teal-50/40 border border-teal-100/30 rounded-xl p-2.5">
            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Hadir</span>
            <span class="text-xs font-black text-slate-800 mt-1 block">{statsHadir} / {statsTotal}</span>
          </div>
          <!-- Izin/Sakit Card -->
          <div class="bg-amber-50/40 border border-amber-100/30 rounded-xl p-2.5">
            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Izin/Sakit</span>
            <span class="text-xs font-black text-slate-800 mt-1 block">{statsIzin + statsSakit}</span>
          </div>
          <!-- Belum Absen Card -->
          <div class="bg-slate-50/60 border border-slate-150 rounded-xl p-2.5">
            <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Belum Absen</span>
            <span class="text-xs font-black text-slate-800 mt-1 block">{statsBelum + statsAlpa}</span>
          </div>
        </div>

        <!-- Search Input -->
        <SearchBar 
          bind:value={attendanceSearchQuery}
          placeholder="Cari nama atau NIS siswa..."
          class="my-4 shrink-0"
        />

        <!-- Scrollable Attendance List -->
        {#if attendanceLoading}
          <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
            <Loader2 class="w-6 h-6 animate-spin text-slate-400 mb-2" />
            <p class="text-[11px] font-semibold">Memuat absensi...</p>
          </div>
        {:else if filteredAttendance.length === 0}
          <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
            <Users class="w-8 h-8 text-slate-300 mb-2" />
            <p class="text-[11px] font-bold">Tidak ada siswa ditemukan</p>
          </div>
        {:else}
          <div class="flex-1 overflow-y-auto max-h-[320px] pr-1 space-y-2 custom-scrollbar">
            {#each filteredAttendance as student}
              <button
                type="button"
                onclick={() => openAttendanceDetail(student)}
                class="w-full flex items-center justify-between p-3 bg-white hover:bg-slate-50/80 border border-slate-100 hover:border-slate-200 rounded-xl transition-all duration-150 cursor-pointer text-left focus:outline-none group/student shadow-xxs"
              >
                <div class="flex items-center gap-3 text-left">
                  <div class="w-8 h-8 rounded-lg bg-slate-50 border border-slate-100 text-slate-500 flex items-center justify-center font-bold text-[10px] shrink-0 select-none group-hover/student:bg-slate-100 transition-colors">
                    {student.namaSiswa ? student.namaSiswa.charAt(0).toUpperCase() : '?'}
                  </div>
                  
                  <div class="space-y-0.5">
                    <h4 class="text-xs font-bold text-slate-700 uppercase truncate max-w-[150px] group-hover/student:text-[#00a294] transition-colors" title={student.namaSiswa}>
                      {student.namaSiswa}
                    </h4>
                    <p class="text-[9px] font-medium text-slate-400">
                      NIS: {student.nis} {#if student.waktuAbsen}<span class="mx-1 text-slate-200">•</span> {student.waktuAbsen.substring(0, 5)}{/if} {#if student.deviceInfo?.model || student.device?.model}<span class="mx-1 text-slate-200">•</span> <span class="text-slate-500 font-semibold">{student.deviceInfo?.model || student.device?.model}</span>{/if}
                    </p>
                  </div>
                </div>
                
                <!-- Attendance Status Badge -->
                {#if student.status === 'hadir'}
                  <span class="px-2 py-0.5 text-[9px] font-bold uppercase rounded-md bg-teal-50 text-[#00a294] border border-teal-150/20">
                    Hadir
                  </span>
                {:else if student.status === 'izin'}
                  <span class="px-2 py-0.5 text-[9px] font-bold uppercase rounded-md bg-amber-50 text-amber-600 border border-amber-150/20">
                    Izin
                  </span>
                {:else if student.status === 'sakit'}
                  <span class="px-2 py-0.5 text-[9px] font-bold uppercase rounded-md bg-orange-50 text-orange-600 border border-orange-150/20">
                    Sakit
                  </span>
                {:else if student.status === 'magang'}
                  <span class="px-2 py-0.5 text-[9px] font-bold uppercase rounded-md bg-blue-50 text-blue-600 border border-blue-150/20">
                    Magang
                  </span>
                {:else if student.status === 'tidak_hadir'}
                  <span class="px-2 py-0.5 text-[9px] font-bold uppercase rounded-md bg-red-50 text-red-650 border border-red-150/20">
                    Alpa
                  </span>
                {:else}
                  <span class="px-2 py-0.5 text-[9px] font-bold uppercase rounded-md bg-slate-50 text-slate-450 border border-slate-200/60">
                    Belum
                  </span>
                {/if}
              </button>
            {/each}
          </div>
        {/if}
      </div>
    {:else}
      <div class="lg:col-span-5 bg-white border border-slate-100 rounded-2xl p-5 shadow-xs flex flex-col min-h-[520px]">
        <!-- Panel Header -->
        <div class="flex items-center justify-between border-b border-slate-100 pb-4">
          <div class="text-left">
            <h2 class="text-xs font-bold text-slate-700 uppercase tracking-tight">Daftar Siswa</h2>
            <p class="text-[10px] text-slate-400 font-medium mt-0.5">Daftar lengkap siswa aktif sekolah</p>
          </div>
          
          <!-- Total count badge -->
          <span class="px-2 py-0.5 rounded-md text-[9px] font-bold uppercase bg-slate-50 text-slate-500 border border-slate-100">
            {$totalStudentsFiltered} siswa
          </span>
        </div>

        <!-- Search Input -->
        <SearchBar 
          bind:value={$studentSearchQuery}
          placeholder="Cari nama, kelas, NIS, atau wali..."
          oninput={handleStudentSearch}
          class="my-4 shrink-0"
        />

        <!-- Scrollable Student List -->
        {#if $dashboardLoading}
          <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
            <Loader2 class="w-6 h-6 animate-spin text-slate-400 mb-2" />
            <p class="text-[11px] font-semibold">Memuat daftar...</p>
          </div>
        {:else if $studentsList.length === 0}
          <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
            <Users class="w-8 h-8 text-slate-300 mb-2" />
            <p class="text-[11px] font-bold">Tidak ada siswa ditemukan</p>
          </div>
        {:else}
          <div class="flex-1 overflow-y-auto max-h-[380px] pr-1 space-y-2 custom-scrollbar">
            {#each $studentsList as student}
              <div class="flex items-center justify-between p-3 bg-white border border-slate-50 rounded-xl transition-colors hover:bg-slate-50/20">
                <div class="flex items-center gap-3 text-left">
                  <!-- Standardized Initials Badge -->
                  <div class="w-8 h-8 rounded-lg bg-slate-50 border border-slate-100 text-slate-500 flex items-center justify-center font-bold text-[10px] shrink-0 select-none">
                    {student.nama.charAt(0).toUpperCase()}
                  </div>
                  
                  <div class="space-y-0.5">
                    <h4 class="text-xs font-bold text-slate-700 uppercase truncate max-w-[150px]">
                      {student.nama}
                    </h4>
                    <p class="text-[10px] font-medium text-slate-400 truncate max-w-[150px]">
                      {student.kelas} <span class="mx-1 text-slate-250">•</span> {student.walas || 'Guru Wali'}
                    </p>
                  </div>
                </div>
                
                <!-- Indicator status dot (online/offline check) -->
                <span class="w-1.5 h-1.5 rounded-full {student.isOnline ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.5)]' : 'bg-slate-300'}"></span>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    {/if}

  </div>
</div>

<!-- Modal Detail Indikator G7 -->
<G7Modal 
  bind:open={detailModalOpen}
  selectedDetailStudent={selectedDetailStudent}
  detailType={detailType}
  detailLoading={detailLoading}
  detailRekap={detailRekap}
  detailEvaluate={detailEvaluate}
/>

<!-- Modal Detail Absensi Siswa (Walas Dashboard) -->
<AttendanceDetailModal
  bind:show={attendanceModalOpen}
  log={selectedAttendanceLog}
  onclose={() => (attendanceModalOpen = false)}
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
