<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { currentUser, clearAuth } from '../../../../stores/authStore';
  import { goto } from '$app/navigation';
  import { addToast, sidebarCollapsed } from '../../../../stores/uiStore';
  import { 
    Home, 
    ClipboardCheck, 
    Image, 
    ChevronRight, 
    LogOut,
    Trophy
  } from 'lucide-svelte';

  let name = $derived($currentUser?.nama || $currentUser?.email || 'Siswa');
  let currentPath = $derived($page.url.pathname);
  let isCollapsed = $derived($sidebarCollapsed);

  onMount(() => {
    if (window.innerWidth < 768) {
      sidebarCollapsed.set(true);
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
    <!-- Header Branding Logo -->
    <div class="{isCollapsed ? 'p-3' : 'p-6'} flex flex-col items-center justify-center border-b border-slate-100 transition-all duration-300">
      <img src="/assets/img/7kaih.png" alt="7 KAIH Logo" class="{isCollapsed ? 'h-9' : 'h-20'} w-auto object-contain transition-all duration-300 hover:scale-105" />
    </div>

    <!-- Navigation links -->
    <div class="{isCollapsed ? 'p-2' : 'p-4'} transition-all duration-300">
      {#if !isCollapsed}
        <span class="block px-4 mb-3 text-[10px] font-extrabold uppercase tracking-widest text-slate-400">
          Menu Utama
        </span>
      {/if}
      
      <nav class="flex flex-col gap-1.5">
        <a
          href="/siswa"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#4db6ac]={currentPath === '/siswa'}
          class:text-white={currentPath === '/siswa'}
          class:text-slate-500={currentPath !== '/siswa'}
          class:hover:bg-slate-50={currentPath !== '/siswa'}
        >
          <div class="flex items-center gap-3">
            <Home class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Dashboard
            {/if}
          </div>
          {#if currentPath === '/siswa' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>

        <a
          href="/siswa/kegiatan"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#4db6ac]={currentPath === '/siswa/kegiatan'}
          class:text-white={currentPath === '/siswa/kegiatan'}
          class:text-slate-500={currentPath !== '/siswa/kegiatan'}
          class:hover:bg-slate-50={currentPath !== '/siswa/kegiatan'}
        >
          <div class="flex items-center gap-3">
            <ClipboardCheck class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Input Kegiatan
            {/if}
          </div>
          {#if currentPath === '/siswa/kegiatan' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>

        <a
          href="/siswa/bukti"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#4db6ac]={currentPath === '/siswa/bukti'}
          class:text-white={currentPath === '/siswa/bukti'}
          class:text-slate-500={currentPath !== '/siswa/bukti'}
          class:hover:bg-slate-50={currentPath !== '/siswa/bukti'}
        >
          <div class="flex items-center gap-3">
            <Image class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Dokumentasi
            {/if}
          </div>
          {#if currentPath === '/siswa/bukti' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>

        <a
          href="/siswa/lomba"
          class="flex items-center {isCollapsed ? 'justify-center p-3' : 'justify-between px-4 py-3'} rounded-2xl text-xs font-bold transition-all duration-200 group"
          class:bg-[#4db6ac]={currentPath === '/siswa/lomba'}
          class:text-white={currentPath === '/siswa/lomba'}
          class:text-slate-500={currentPath !== '/siswa/lomba'}
          class:hover:bg-slate-50={currentPath !== '/siswa/lomba'}
        >
          <div class="flex items-center gap-3">
            <Trophy class="w-5 h-5 transition-transform duration-200 group-hover:scale-105" />
            {#if !isCollapsed}
              Lomba Kebersihan
            {/if}
          </div>
          {#if currentPath === '/siswa/lomba' && !isCollapsed}
            <ChevronRight class="w-4 h-4 text-white" />
          {/if}
        </a>
      </nav>
    </div>
  </div>

  <!-- Help Banner at the Bottom -->
  {#if !isCollapsed}
    <div class="m-4 p-4 bg-[#e0f2f1]/30 border border-[#4db6ac]/20 rounded-2xl text-left shadow-xxs">
      <h4 class="text-xs font-bold text-[#0070f3] mb-1">Butuh Bantuan?</h4>
      <p class="text-[10px] text-slate-500 font-medium leading-relaxed">
        Hubungi wali siswa jika mengalami kendala, atau lihat Manual Book GR31.
      </p>
    </div>
  {/if}

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
