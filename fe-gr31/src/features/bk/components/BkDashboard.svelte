<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { 
    MessageSquare, 
    AlertCircle, 
    Clock, 
    CheckCircle, 
    ChevronRight, 
    User, 
    ArrowRight, 
    Loader2, 
    Info 
  } from 'lucide-svelte';
  import { currentUser } from '../../../stores/authStore';
  import {
    loadingDashboard,
    totalAduan,
    openAduan,
    inProgressAduan,
    closedAduan,
    recentActiveAduan,
    loadBkDashboardData
  } from '../logic/bkDashboardLogic';

  onMount(() => {
    loadBkDashboardData();
  });

  function formatDate(isoStr: string) {
    if (!isoStr) return '';
    const date = new Date(isoStr);
    return date.toLocaleDateString('id-ID', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<div class="space-y-6 select-none font-sans pb-10">
  <!-- Welcome Banner -->
  <div class="space-y-1 text-left">
    <h1 class="text-xl font-bold tracking-tight text-slate-800 uppercase leading-none">
      Beranda Guru BK
    </h1>
    <p class="text-xs text-slate-400 font-medium">
      Selamat datang di halaman bimbingan dan konseling. Anda bisa memantau dan menindaklanjuti keluhan / aduan siswa di sini.
    </p>
  </div>

  <!-- Stats Grid -->
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
    <!-- Total Aduan Card -->
    <div class="bg-white border border-slate-100 rounded-2xl p-4 shadow-xs flex items-center gap-4">
      <div class="w-10 h-10 rounded-xl bg-slate-50 text-[#00a294] border border-slate-100 flex items-center justify-center shrink-0">
        <MessageSquare class="w-5 h-5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Total Aduan</p>
        <h3 class="text-base font-extrabold text-slate-700 mt-0.5">{$totalAduan} Tiket</h3>
      </div>
    </div>

    <!-- Baru (Open) Card -->
    <div class="bg-white border border-slate-100 rounded-2xl p-4 shadow-xs flex items-center gap-4">
      <div class="w-10 h-10 rounded-xl bg-amber-50 text-amber-500 border border-amber-100/60 flex items-center justify-center shrink-0">
        <AlertCircle class="w-5 h-5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Aduan Baru</p>
        <h3 class="text-base font-extrabold text-slate-700 mt-0.5">{$openAduan} Tiket</h3>
      </div>
    </div>

    <!-- Sedang Diproses Card -->
    <div class="bg-white border border-slate-100 rounded-2xl p-4 shadow-xs flex items-center gap-4">
      <div class="w-10 h-10 rounded-xl bg-blue-50 text-blue-500 border border-blue-100/60 flex items-center justify-center shrink-0">
        <Clock class="w-5 h-5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Sedang Diproses</p>
        <h3 class="text-base font-extrabold text-slate-700 mt-0.5">{$inProgressAduan} Tiket</h3>
      </div>
    </div>

    <!-- Selesai Card -->
    <div class="bg-white border border-slate-100 rounded-2xl p-4 shadow-xs flex items-center gap-4">
      <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-500 border border-emerald-100/60 flex items-center justify-center shrink-0">
        <CheckCircle class="w-5 h-5" />
      </div>
      <div class="text-left">
        <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Telah Selesai</p>
        <h3 class="text-base font-extrabold text-slate-700 mt-0.5">{$closedAduan} Tiket</h3>
      </div>
    </div>
  </div>

  <!-- Two Column Content Panel -->
  <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
    
    <!-- Left Column: Recent Active Complaints (Col span 7) -->
    <div class="lg:col-span-7 bg-white border border-slate-100 rounded-2xl p-5 shadow-xs flex flex-col min-h-[420px]">
      <div class="flex items-center justify-between border-b border-slate-100 pb-4">
        <div class="text-left">
          <h2 class="text-xs font-bold text-slate-700 uppercase tracking-tight">Aduan Perlu Tindak Lanjut</h2>
          <p class="text-[10px] text-slate-400 font-medium mt-0.5">Daftar aduan aktif terbaru dari siswa</p>
        </div>
        
        <button 
          onclick={loadBkDashboardData}
          class="px-2.5 py-1.5 text-[10px] font-bold text-slate-600 hover:text-slate-900 border border-slate-150 rounded-lg bg-slate-50/50 hover:bg-slate-50 transition-colors cursor-pointer"
        >
          Segarkan
        </button>
      </div>

      <!-- Scrollable List -->
      {#if $loadingDashboard}
        <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
          <Loader2 class="w-6 h-6 animate-spin text-slate-400 mb-2" />
          <p class="text-[11px] font-semibold">Memuat data aduan...</p>
        </div>
      {:else}
        {#if $recentActiveAduan.length === 0}
          <div class="flex-1 flex flex-col items-center justify-center py-20 text-slate-400">
            <CheckCircle class="w-8 h-8 text-slate-300 mb-2" />
            <p class="text-[11px] font-bold">Semua aduan telah diselesaikan</p>
          </div>
        {:else}
          <div class="flex-1 overflow-y-auto max-h-[320px] pr-1 space-y-2 mt-4 custom-scrollbar">
            {#each $recentActiveAduan as aduan}
              <button 
                onclick={() => goto('/admin/chat')}
                class="w-full flex items-center justify-between p-3.5 bg-white hover:bg-slate-50/50 border border-slate-100 hover:border-slate-200/80 rounded-xl transition-all group cursor-pointer text-left"
              >
                <div class="flex items-center gap-3">
                  <!-- Standardized Initials Badge -->
                  <div class="w-9 h-9 rounded-lg bg-slate-50 border border-slate-100 text-slate-500 flex items-center justify-center font-bold text-xs shrink-0 select-none">
                    {aduan.namaSiswa.charAt(0).toUpperCase()}
                  </div>
                  
                  <div class="space-y-0.5 max-w-[200px] sm:max-w-[320px]">
                    <h4 class="text-xs font-bold text-slate-700 uppercase group-hover:text-slate-900 transition-colors truncate">
                      {aduan.namaSiswa}
                    </h4>
                    <p class="text-[10px] font-medium text-slate-400">
                      Kelas {aduan.kelas} <span class="mx-1 text-slate-200">•</span> {formatDate(aduan.updatedAt)}
                    </p>
                    {#if aduan.messages && aduan.messages.length > 0}
                      <p class="text-[10px] text-slate-500 truncate mt-1">
                        {aduan.messages[aduan.messages.length - 1].isi}
                      </p>
                    {/if}
                  </div>
                </div>
                
                <div class="flex items-center gap-2">
                  <span 
                    class="px-2 py-0.5 rounded-md text-[9px] font-bold uppercase border"
                    class:bg-amber-50={aduan.status === 'open'}
                    class:text-amber-700={aduan.status === 'open'}
                    class:border-amber-100={aduan.status === 'open'}
                    class:bg-blue-50={aduan.status === 'in_progress'}
                    class:text-blue-700={aduan.status === 'in_progress'}
                    class:border-blue-100={aduan.status === 'in_progress'}
                  >
                    {aduan.status === 'open' ? 'Baru' : 'Diproses'}
                  </span>
                  <ChevronRight class="w-4 h-4 text-slate-300 group-hover:text-slate-400 transition-all" />
                </div>
              </button>
            {/each}
          </div>
        {/if}
      {/if}
    </div>

    <!-- Right Column: Account & CTAs (Col span 5) -->
    <div class="lg:col-span-5 space-y-6">
      
      <!-- Account Profile Box -->
      <div class="bg-white border border-slate-100 rounded-2xl p-5 shadow-xs">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl bg-slate-50 border border-slate-100 text-[#00a294] flex items-center justify-center shrink-0">
            <User class="w-6 h-6" />
          </div>
          <div class="text-left">
            <p class="text-[9px] font-bold uppercase tracking-wider text-slate-400">Profil Akun</p>
            <h3 class="text-sm font-bold text-slate-700 mt-0.5">{$currentUser?.nama || 'Admin BK SMKN 31 Jakarta'}</h3>
            <p class="text-[10px] text-slate-400 font-medium mt-0.5">{$currentUser?.email || 'bk@smkn31jkt.id'}</p>
          </div>
        </div>
      </div>

      <!-- Teal CTA Box (Pantau Chat) -->
      <button 
        onclick={() => goto('/admin/chat')}
        class="w-full bg-[#00a294] hover:bg-[#008f83] text-white border-none rounded-2xl p-5 shadow-xs flex items-center justify-between text-left cursor-pointer transition-all duration-200 group active:scale-98"
      >
        <div class="space-y-1">
          <p class="text-[9px] font-bold uppercase tracking-wider text-teal-100">Aduan Siswa</p>
          <h3 class="text-sm font-extrabold text-white">Pantau Chat</h3>
          <p class="text-[10px] text-teal-50/80 font-medium mt-0.5">Tanggapi aduan dan keluhan siswa secara langsung</p>
        </div>
        <div class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center text-white shrink-0 group-hover:translate-x-1 transition-transform">
          <ArrowRight class="w-4 h-4" />
        </div>
      </button>

      <!-- System Information Box -->
      <div class="bg-white border border-slate-100 rounded-2xl p-5 shadow-xs text-left">
        <div class="flex items-center gap-2 border-b border-slate-100 pb-3">
          <Info class="w-4 h-4 text-slate-400" />
          <h3 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Informasi Sistem</h3>
        </div>
        <div class="mt-3 space-y-2.5">
          <p class="text-[11px] text-slate-500 leading-relaxed">
            Sebagai Guru BK / Piket, tugas utama Anda adalah menanggapi laporan atau keluhan yang dikirimkan oleh siswa, baik secara langsung maupun yang diteruskan oleh guru wali mereka.
          </p>
          <p class="text-[11px] text-slate-500 leading-relaxed font-bold border-l-2 border-[#00a294] pl-2 bg-slate-50/50 py-1.5 rounded-r-md">
            Pastikan Anda mengecek menu Aduan Siswa (Chat) secara berkala agar setiap permasalahan dapat segera ditindaklanjuti.
          </p>
        </div>
      </div>

    </div>

  </div>
</div>

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
