<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { currentUser, clearAuth } from '../../../../stores/authStore';
  import { goto } from '$app/navigation';
  import { addToast, sidebarCollapsed } from '../../../../stores/uiStore';
  import { 
    ChevronRight, 
    MessageSquare, 
    School, 
    Users, 
    Settings, 
    LogOut,
    Database,
    Eye,
    Laptop,
    PlusCircle,
    Trash2,
    ChevronDown,
    Home,
    Trophy
  } from 'lucide-svelte';
  import { pendingWalasAduan, loadPendingAduanForWalas } from '../../../../features/admin/logic/adminDashboardLogic';

  let role = $derived($currentUser?.role);
  let currentPath = $derived($page.url.pathname);
  let isCollapsed = $derived($sidebarCollapsed);

  onMount(() => {
    loadPendingAduanForWalas();
    if (window.innerWidth < 768) {
      sidebarCollapsed.set(true);
    }
  });

  let isSuperOrWalas = $derived(
    role === 'super_admin' || 
    role === 'walas' || 
    role === 'guru_wali' || 
    $currentUser?.is_walas === true || 
    $currentUser?.isWalas === true
  );

  let isPiketOrSuper = $derived(
    role === 'super_admin' || 
    role === 'piket' || 
    role === 'admin_piket' || 
    $currentUser?.is_piket === true ||
    $currentUser?.isPiket === true
  );

  // Dropdown states for Admin
  let monitoringOpen = $state(false);
  let aturUserOpen = $state(false);
  let aduanSiswaOpen = $state(false);

  $effect(() => {
    if (currentPath.includes('/admin/monitoring/')) {
      monitoringOpen = true;
    }
  });

  function handleLogout() {
    clearAuth();
    addToast('Anda berhasil logout', 'success');
    goto('/', { replaceState: true });
  }
</script>

<aside class="bg-white border-r border-slate-100 flex flex-col justify-between shadow-[2px_0_8px_rgba(0,0,0,0.01)] shrink-0 h-screen select-none font-sans transition-all duration-300 md:sticky md:top-0 max-md:fixed max-md:top-0 max-md:left-0 max-md:z-50 max-md:h-screen max-md:w-64 max-md:shadow-xl {isCollapsed ? 'md:w-20 max-md:-translate-x-full' : 'md:w-64 max-md:translate-x-0'}">
  <div class="flex flex-col sidebar-scrollbar flex-1">
    <div class="{isCollapsed ? 'p-3' : 'p-6'} flex flex-col items-center justify-center border-b border-slate-100 transition-all duration-300">
      <img src="/assets/img/7kaih.png" alt="7 KAIH Logo" class="{isCollapsed ? 'h-9' : 'h-20'} w-auto object-contain transition-all duration-300 hover:scale-105" />
    </div>

    <div class="{isCollapsed ? 'p-2' : 'p-4'} transition-all duration-300">
      {#if !isCollapsed}
        <span class="block px-4 mb-3 text-[10px] font-extrabold uppercase tracking-widest text-slate-400">
          ADMIN MENU
        </span>
      {/if}
      
      <nav class="flex flex-col gap-1.5">
        <a
          href="/admin"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#00a294]={currentPath === '/admin'}
          class:text-white={currentPath === '/admin'}
          class:text-slate-500={currentPath !== '/admin'}
          class:hover:bg-slate-50={currentPath !== '/admin'}
        >
          <div class="flex items-center gap-3">
            <Home class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Dashboard
            {/if}
          </div>
          {#if currentPath === '/admin' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>

        <a
          href="/admin/g7"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#00a294]={currentPath === '/admin/g7'}
          class:text-white={currentPath === '/admin/g7'}
          class:text-slate-500={currentPath !== '/admin/g7'}
          class:hover:bg-slate-50={currentPath !== '/admin/g7'}
        >
          <div class="flex items-center gap-3">
            <Database class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Rekap Data
            {/if}
          </div>
          {#if currentPath === '/admin/g7' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>

        <div class="flex flex-col">
          <button
            type="button"
            onclick={() => !isCollapsed && (monitoringOpen = !monitoringOpen)}
            class="w-full flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 text-slate-500 hover:bg-slate-50 border-none cursor-pointer text-left"
          >
            <div class="flex items-center gap-3">
              <Laptop class="w-5 h-5 transition-transform duration-200" />
              {#if !isCollapsed}
                Monitoring
              {/if}
            </div>
            {#if !isCollapsed}
              <ChevronDown class="w-4 h-4 transition-transform duration-200 {monitoringOpen ? 'rotate-180' : ''}" />
            {/if}
          </button>
          {#if monitoringOpen && !isCollapsed}
            <div class="pl-9 pr-2 py-1 flex flex-col gap-1 border-l border-slate-100/60 ml-6 mt-1">
              <a href="/admin/monitoring/bangun-pagi" class="px-3 py-2 text-[11px] font-bold rounded-xl transition-all {currentPath === '/admin/monitoring/bangun-pagi' ? 'text-[#00a294] bg-slate-50/50' : 'text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50'}">
                Bangun Pagi
              </a>
              <a href="/admin/monitoring/beribadah" class="px-3 py-2 text-[11px] font-bold rounded-xl transition-all {currentPath === '/admin/monitoring/beribadah' ? 'text-[#00a294] bg-slate-50/50' : 'text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50'}">
                Beribadah
              </a>
              <a href="/admin/monitoring/makan-sehat" class="px-3 py-2 text-[11px] font-bold rounded-xl transition-all {currentPath === '/admin/monitoring/makan-sehat' ? 'text-[#00a294] bg-slate-50/50' : 'text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50'}">
                Makan Sehat
              </a>
              <a href="/admin/monitoring/olahraga" class="px-3 py-2 text-[11px] font-bold rounded-xl transition-all {currentPath === '/admin/monitoring/olahraga' ? 'text-[#00a294] bg-slate-50/50' : 'text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50'}">
                Olahraga
              </a>
              <a href="/admin/monitoring/belajar" class="px-3 py-2 text-[11px] font-bold rounded-xl transition-all {currentPath === '/admin/monitoring/belajar' ? 'text-[#00a294] bg-slate-50/50' : 'text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50'}">
                Belajar
              </a>
              <a href="/admin/monitoring/bermasyarakat" class="px-3 py-2 text-[11px] font-bold rounded-xl transition-all {currentPath === '/admin/monitoring/bermasyarakat' ? 'text-[#00a294] bg-slate-50/50' : 'text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50'}">
                Bermasyarakat
              </a>
              <a href="/admin/monitoring/tidur-cukup" class="px-3 py-2 text-[11px] font-bold rounded-xl transition-all {currentPath === '/admin/monitoring/tidur-cukup' ? 'text-[#00a294] bg-slate-50/50' : 'text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50'}">
                Tidur Cukup
              </a>
            </div>
          {/if}
        </div>

        <!-- 4. Tambah Absensi -->
        {#if isPiketOrSuper}
          <a
            href="/admin/kehadiran"
            class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
            class:bg-[#00a294]={currentPath === '/admin/kehadiran'}
            class:text-white={currentPath === '/admin/kehadiran'}
            class:text-slate-500={currentPath !== '/admin/kehadiran'}
            class:hover:bg-slate-50={currentPath !== '/admin/kehadiran'}
          >
            <div class="flex items-center gap-3">
              <PlusCircle class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
              {#if !isCollapsed}
                Tambah Absensi
              {/if}
            </div>
            {#if currentPath === '/admin/kehadiran' && !isCollapsed}
              <ChevronRight class="w-4 h-4 text-white" />
            {/if}
          </a>
        {/if}

        <!-- 5. Bukti Kegiatan -->
        <a
          href="/admin/sekolah"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#00a294]={currentPath === '/admin/sekolah'}
          class:text-white={currentPath === '/admin/sekolah'}
          class:text-slate-500={currentPath !== '/admin/sekolah'}
          class:hover:bg-slate-50={currentPath !== '/admin/sekolah'}
        >
          <div class="flex items-center gap-3">
            <School class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Bukti Kegiatan
            {/if}
          </div>
          {#if currentPath === '/admin/sekolah' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>

        {#if isSuperOrWalas}
          <!-- Lomba Kebersihan -->
          <a
            href="/admin/lomba"
            class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
            class:bg-[#00a294]={currentPath === '/admin/lomba'}
            class:text-white={currentPath === '/admin/lomba'}
            class:text-slate-500={currentPath !== '/admin/lomba'}
            class:hover:bg-slate-50={currentPath !== '/admin/lomba'}
          >
            <div class="flex items-center gap-3">
              <Trophy class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
              {#if !isCollapsed}
                Lomba Kebersihan
              {/if}
            </div>
            {#if currentPath === '/admin/lomba' && !isCollapsed}
              <ChevronRight class="w-4 h-4 text-white" />
            {/if}
          </a>
        {/if}

        <!-- 6. Atur User (Collapsible Dropdown) -->
        <div class="flex flex-col">
          <button
            type="button"
            onclick={() => !isCollapsed && (aturUserOpen = !aturUserOpen)}
            class="w-full flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 text-slate-500 hover:bg-slate-50 border-none cursor-pointer text-left"
          >
            <div class="flex items-center gap-3">
              <Users class="w-5 h-5 transition-transform duration-200" />
              {#if !isCollapsed}
                Atur User
              {/if}
            </div>
            {#if !isCollapsed}
              <ChevronDown class="w-4 h-4 transition-transform duration-200 {aturUserOpen ? 'rotate-180' : ''}" />
            {/if}
          </button>
          {#if aturUserOpen && !isCollapsed}
            <div class="pl-9 pr-2 py-1 flex flex-col gap-1 border-l border-slate-100 ml-6 mt-1">
              {#if role === 'super_admin'}
                <a href="/admin/admins" class="px-3 py-2 text-[11px] font-bold text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50 rounded-xl transition-all">
                  Kelola Admin
                </a>
              {/if}
              <a href="/admin/students" class="px-3 py-2 text-[11px] font-bold text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50 rounded-xl transition-all">
                Kelola Siswa
              </a>
            </div>
          {/if}
        </div>

        <!-- 7. Aduan Siswa (Collapsible Dropdown) -->
        <div class="flex flex-col">
          <button
            type="button"
            onclick={() => !isCollapsed && (aduanSiswaOpen = !aduanSiswaOpen)}
            class="w-full flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 text-slate-500 hover:bg-slate-50 border-none cursor-pointer text-left"
          >
            <div class="flex items-center gap-3 relative">
              <MessageSquare class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
              {#if isCollapsed && isSuperOrWalas && $pendingWalasAduan.length > 0}
                <span class="absolute top-0.5 right-0.5 block h-2 w-2 rounded-full bg-rose-500 ring-2 ring-white animate-pulse"></span>
              {/if}
              {#if !isCollapsed}
                <span class="flex items-center gap-2">
                  Aduan Siswa
                  {#if isSuperOrWalas && $pendingWalasAduan.length > 0}
                    <span class="px-1.5 py-0.5 rounded-full bg-rose-500 text-white text-[9px] font-black leading-none animate-pulse shrink-0">
                      {$pendingWalasAduan.length}
                    </span>
                  {/if}
                </span>
              {/if}
            </div>
            {#if !isCollapsed}
              <ChevronDown class="w-4 h-4 transition-transform duration-200 {aduanSiswaOpen ? 'rotate-180' : ''}" />
            {/if}
          </button>
          {#if aduanSiswaOpen && !isCollapsed}
            <div class="pl-9 pr-2 py-1 flex flex-col gap-1 border-l border-slate-100 ml-6 mt-1">
              <a href="/admin/chat" class="px-3 py-2 text-[11px] font-bold text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50 rounded-xl transition-all flex items-center justify-between">
                <span>Chat Aduan</span>
                {#if isSuperOrWalas && $pendingWalasAduan.length > 0}
                  <span class="px-1.5 py-0.5 rounded-full bg-rose-500 text-white text-[9px] font-black leading-none animate-pulse shrink-0">
                    {$pendingWalasAduan.length}
                  </span>
                {/if}
              </a>
              <a href="/admin/chat" class="px-3 py-2 text-[11px] font-bold text-slate-500 hover:text-[#00a294] hover:bg-slate-50/50 rounded-xl transition-all">
                Monitoring Aduan
              </a>
            </div>
          {/if}
        </div>

        <!-- 8. Pengaturan Akun -->
        <a
          href="/admin/settings"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#00a294]={currentPath === '/admin/settings'}
          class:text-white={currentPath === '/admin/settings'}
          class:text-slate-500={currentPath !== '/admin/settings'}
          class:hover:bg-slate-50={currentPath !== '/admin/settings'}
        >
          <div class="flex items-center gap-3">
            <Settings class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Pengaturan Akun
            {/if}
          </div>
          {#if currentPath === '/admin/settings' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>

        <!-- 9. Hapus Data -->
        <a
          href="/admin/delete-data"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group {currentPath === '/admin/delete-data' ? 'bg-rose-600 text-white' : 'text-slate-500 hover:bg-rose-50/50'}"
        >
          <div class="flex items-center gap-3">
            <Trash2 class="w-5 h-5 transition-transform duration-200 group-hover:scale-105 {currentPath === '/admin/delete-data' ? 'text-white' : 'text-rose-500'}" />
            {#if !isCollapsed}
              <span>Hapus Data</span>
            {/if}
          </div>
          {#if currentPath === '/admin/delete-data' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>
      </nav>
    </div>
  </div>

  <!-- Logout Button footer -->
  <div class="p-4 border-t border-slate-100 bg-slate-50/50">
    <button
      onclick={handleLogout}
      class="w-full flex items-center justify-center gap-2 px-4 py-2.5 text-xs font-bold text-rose-600 bg-rose-50 hover:bg-rose-100 rounded-xl transition-all duration-300 cursor-pointer"
    >
      <LogOut class="w-4 h-4" />
      {#if !isCollapsed}
        Keluar Aplikasi
      {/if}
    </button>
  </div>
</aside>

{#if !isCollapsed}
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    role="presentation"
    class="fixed inset-0 bg-slate-900/40 backdrop-blur-xs z-40 transition-opacity duration-300 md:hidden"
    onclick={() => sidebarCollapsed.set(true)}
  ></div>
{/if}
