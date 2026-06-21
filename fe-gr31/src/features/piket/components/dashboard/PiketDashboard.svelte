<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { goto } from '$app/navigation';
  import { 
    User, 
    Users, 
    CheckCircle, 
    Clock, 
    Briefcase, 
    Calendar,
    QrCode, 
    Copy, 
    Download as DownloadIcon, 
    RefreshCw, 
    Loader2, 
    Activity,
    Info,
    AlertCircle
  } from 'lucide-svelte';
  import { currentUser } from '../../../../stores/authStore';
  import { addToast } from '../../../../stores/uiStore';
  import {
    loadingDashboard,
    rekapHarian,
    currentTime,
    currentDateStr,
    countdownStr,
    getTodayDateString,
    updateClockAndCountdown,
    loadPiketDashboardData
  } from '../../logic/piketDashboardLogic';

  let timerInterval: any;

  onMount(() => {
    loadPiketDashboardData();
    updateClockAndCountdown();
    // Update clock and countdown every second
    timerInterval = setInterval(updateClockAndCountdown, 1000);
  });

  onDestroy(() => {
    if (timerInterval) clearInterval(timerInterval);
  });

  // Derived counts for the stats display
  let totalSiswa = $derived($rekapHarian?.totalSiswa || 0);
  let totalHadir = $derived($rekapHarian?.totalHadir || 0);
  let totalIzinSakit = $derived(($rekapHarian?.totalIzin || 0) + ($rekapHarian?.totalSakit || 0));
  let totalMagang = $derived($rekapHarian?.totalMagang || 0);
  let totalBelumAbsen = $derived(
    Math.max(0, totalSiswa - totalHadir - totalIzinSakit - totalMagang)
  );

  function getQrData() {
    const today = getTodayDateString();
    return `SMKN31-ATTENDANCE-KEY-${today}`;
  }

  function handleCopyLink() {
    const qrData = getQrData();
    navigator.clipboard.writeText(qrData);
    addToast('Token QR berhasil disalin ke clipboard', 'success');
  }

  function handleDownload() {
    const today = getTodayDateString();
    const qrUrl = `https://api.qrserver.com/v1/create-qr-code/?size=350x350&data=SMKN31-ATTENDANCE-KEY-${today}`;
    
    const link = document.createElement('a');
    link.href = qrUrl;
    link.target = '_blank';
    link.download = `QR_Absen_${today}.png`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    addToast('Membuka QR Code di tab baru untuk diunduh', 'success');
  }
</script>

<div class="space-y-6 select-none font-sans pb-10 text-slate-700">
  <!-- Welcome Title Header -->
  <div class="flex items-center justify-between border-b border-slate-100 pb-5">
    <div class="space-y-1 text-left">
      <h1 class="text-lg font-bold tracking-tight text-slate-800 uppercase">
        Dashboard Guru Piket
      </h1>
      <p class="text-xs text-slate-400 font-medium">
        Monitoring absensi dan kehadiran harian siswa secara real-time.
      </p>
    </div>

    <button 
      onclick={loadPiketDashboardData}
      disabled={$loadingDashboard}
      class="flex items-center gap-1.5 px-3.5 py-2 text-xs font-semibold text-slate-600 hover:text-slate-900 border border-slate-200 rounded-xl bg-white hover:bg-slate-50 transition-all cursor-pointer shadow-xs active:scale-98 disabled:opacity-50"
    >
      <RefreshCw class="w-3.5 h-3.5 {$loadingDashboard ? 'animate-spin' : ''}" />
      Segarkan
    </button>
  </div>

  <!-- Stats Grid Row -->
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
    <!-- Guru Piket Profile Card -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xs flex items-center gap-3.5">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-450 border border-slate-100 flex items-center justify-center shrink-0">
        <User class="w-4.5 h-4.5" />
      </div>
      <div class="text-left truncate">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Petugas Piket</p>
        <h3 class="text-xs font-bold text-slate-700 mt-0.5 truncate max-w-[125px]" title={$currentUser?.nama}>
          {$currentUser?.nama || 'Guru Piket'}
        </h3>
      </div>
    </div>

    <!-- Total Students Card -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xs flex items-center gap-3.5 hover:border-slate-200 transition-all">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-500 border border-slate-100 flex items-center justify-center shrink-0">
        <Users class="w-4.5 h-4.5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Total Siswa</p>
        <h3 class="text-sm font-bold text-slate-800 mt-0.5">
          {#if $loadingDashboard}
            <Loader2 class="w-3 h-3 animate-spin text-slate-400" />
          {:else}
            {totalSiswa} <span class="text-[10px] font-normal text-slate-400">Siswa</span>
          {/if}
        </h3>
      </div>
    </div>

    <!-- Present Today Card -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xs flex items-center gap-3.5 hover:border-slate-200 transition-all">
      <div class="w-10 h-10 rounded-xl bg-teal-50/50 text-[#00a294] border border-teal-100/30 flex items-center justify-center shrink-0">
        <CheckCircle class="w-4.5 h-4.5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-450">Hadir</p>
        <h3 class="text-sm font-bold text-slate-800 mt-0.5">
          {#if $loadingDashboard}
            <Loader2 class="w-3 h-3 animate-spin text-slate-400" />
          {:else}
            {totalHadir} <span class="text-[10px] font-normal text-slate-450">Siswa</span>
          {/if}
        </h3>
      </div>
    </div>

    <!-- Sick / Leave Card -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xs flex items-center gap-3.5 hover:border-slate-200 transition-all">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-500 border border-slate-100 flex items-center justify-center shrink-0">
        <Calendar class="w-4.5 h-4.5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Izin / Sakit</p>
        <h3 class="text-sm font-bold text-slate-800 mt-0.5">
          {#if $loadingDashboard}
            <Loader2 class="w-3 h-3 animate-spin text-slate-400" />
          {:else}
            {totalIzinSakit} <span class="text-[10px] font-normal text-slate-400">Siswa</span>
          {/if}
        </h3>
      </div>
    </div>

    <!-- Internship Card -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xs flex items-center gap-3.5 hover:border-slate-200 transition-all">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-550 border border-slate-100 flex items-center justify-center shrink-0">
        <Briefcase class="w-4.5 h-4.5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Magang</p>
        <h3 class="text-sm font-bold text-slate-800 mt-0.5">
          {#if $loadingDashboard}
            <Loader2 class="w-3 h-3 animate-spin text-slate-400" />
          {:else}
            {totalMagang} <span class="text-[10px] font-normal text-slate-400">Siswa</span>
          {/if}
        </h3>
      </div>
    </div>
  </div>

  <!-- Monitoring Absensi Siswa CTA Panel -->
  <div class="bg-white border border-slate-100/80 rounded-2xl p-5 shadow-xs flex flex-col md:flex-row md:items-center justify-between gap-6 text-left">
    <div class="space-y-2 max-w-2xl">
      <h3 class="text-xs font-bold text-slate-800 uppercase tracking-wider">Rekapitulasi Kehadiran Hari Ini</h3>
      <p class="text-[11px] text-slate-450 leading-relaxed font-medium">
        Seluruh log kehadiran siswa dikumpulkan secara otomatis melalui scan QR Code mandiri oleh siswa. Klik tombol di samping untuk memantau detail data absensi, mengubah status kehadiran, atau menginput absen siswa secara manual.
      </p>
      
      <!-- Metrics row badges -->
      <div class="flex flex-wrap gap-2 pt-1.5">
        <span class="px-2.5 py-1 rounded-lg text-[9px] font-bold uppercase bg-teal-50/60 text-[#00a294] border border-teal-100/30">
          {totalHadir} Hadir
        </span>
        <span class="px-2.5 py-1 rounded-lg text-[9px] font-bold uppercase bg-slate-50 text-slate-600 border border-slate-200/50">
          {totalIzinSakit} Izin/Sakit
        </span>
        <span class="px-2.5 py-1 rounded-lg text-[9px] font-bold uppercase bg-slate-50 text-slate-600 border border-slate-200/50">
          {totalMagang} Magang
        </span>
        <span class="px-2.5 py-1 rounded-lg text-[9px] font-bold uppercase bg-slate-100/50 text-slate-500 border border-slate-200/30">
          {totalBelumAbsen} Belum Absen
        </span>
      </div>
    </div>
    
    <button 
      onclick={() => goto('/piket/monitoring')}
      class="flex items-center justify-center gap-1.5 px-5 py-3 bg-[#00a294] hover:bg-[#008f83] text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none shrink-0"
    >
      <Activity class="w-4 h-4" />
      Buka Monitoring
    </button>
  </div>

  <!-- QR Kehadiran Harian Card -->
  <div class="bg-white border border-slate-100/80 rounded-2xl p-6 shadow-xs text-left">
    <div class="flex items-center gap-2 border-b border-slate-100 pb-4 mb-6">
      <QrCode class="w-5 h-5 text-[#00a294]" />
      <div>
        <h3 class="text-xs font-bold text-slate-800 uppercase tracking-wider">QR Code Kehadiran Harian</h3>
        <p class="text-[10px] text-slate-400 font-medium">QR Code ini berubah otomatis setiap hari demi keamanan absensi siswa</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-center">
      <!-- Left: QR Code Image -->
      <div class="lg:col-span-5 flex justify-center border border-slate-100 rounded-2xl p-5 bg-slate-50/30 max-w-[280px] mx-auto shadow-inner">
        <div class="bg-white p-3 rounded-xl border border-slate-100 shadow-xs">
          <img 
            src="https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=SMKN31-ATTENDANCE-KEY-{getTodayDateString()}" 
            alt="Kehadiran QR Code" 
            class="w-full aspect-square max-w-[190px] object-contain"
          />
        </div>
      </div>

      <!-- Right: Date, Time & Reset countdown -->
      <div class="lg:col-span-7 space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
          <!-- Calendar Date -->
          <div class="border border-slate-100 rounded-xl p-3.5 flex items-center gap-3 bg-white shadow-xs">
            <div class="w-9 h-9 rounded-lg bg-[#00a294]/5 text-[#00a294] flex items-center justify-center shrink-0 border border-teal-100/20">
              <Calendar class="w-4 h-4" />
            </div>
            <div class="truncate">
              <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Tanggal</span>
              <span class="text-xs font-bold text-slate-700 block mt-0.5">{$currentDateStr}</span>
            </div>
          </div>

          <!-- Live Clock -->
          <div class="border border-slate-100 rounded-xl p-3.5 flex items-center gap-3 bg-white shadow-xs">
            <div class="w-9 h-9 rounded-lg bg-slate-50 text-slate-500 flex items-center justify-center shrink-0 border border-slate-100">
              <Clock class="w-4 h-4" />
            </div>
            <div>
              <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Waktu Sekarang</span>
              <span class="text-xs font-bold text-slate-700 block mt-0.5 font-mono">{$currentTime}</span>
            </div>
          </div>
        </div>

        <!-- Countdown box -->
        <div class="bg-slate-50/60 border border-slate-150 rounded-xl p-4 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <AlertCircle class="w-4 h-4 text-slate-400" />
            <span class="text-[10px] font-bold text-slate-555 uppercase tracking-wider">Masa Berlaku Token Sisa</span>
          </div>
          <span class="text-sm font-bold text-slate-700 font-mono tracking-widest bg-white px-3 py-1 border border-slate-150 rounded-lg shadow-xxs">{$countdownStr}</span>
        </div>

        <!-- Actions Buttons -->
        <div class="flex flex-wrap gap-3 pt-2">
          <button 
            onclick={handleCopyLink}
            class="flex items-center justify-center gap-2 px-5 py-2.5 border border-slate-200 hover:border-slate-300 text-slate-600 bg-white hover:bg-slate-50 rounded-xl font-bold text-xs transition-all cursor-pointer shadow-xs active:scale-98"
          >
            <Copy class="w-4 h-4 text-slate-450" />
            Salin Kode Token
          </button>
          
          <button 
            onclick={handleDownload}
            class="flex items-center justify-center gap-2 px-5 py-2.5 bg-[#00a294] hover:bg-[#008f83] text-white rounded-xl font-bold text-xs transition-all cursor-pointer shadow-xs border-none active:scale-98"
          >
            <DownloadIcon class="w-4 h-4 text-white" />
            Unduh Gambar QR
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
