<script lang="ts">
  import ProtectedRoute from '../../features/shared/components/ProtectedRoute.svelte';
  import StudentSidebar from '../../features/shared/components/sidebar/StudentSidebar.svelte';
  import { currentUser } from '../../stores/authStore';
  import { sidebarCollapsed } from '../../stores/uiStore';
  import { Menu, Bell } from 'lucide-svelte';

  let { children } = $props();

  let name = $derived($currentUser?.nama || 'Siswa');
  let nickname = $derived(name.split(' ')[0]);
  let initial = $derived(name.charAt(0).toUpperCase());

  // Track scroll position of the main viewport to blend/unblend header
  let scrolled = $state(false);

  function handleScroll(e: Event) {
    const target = e.currentTarget as HTMLElement;
    scrolled = target.scrollTop > 10;
  }
</script>

<ProtectedRoute allowedRoles={['siswa']}>
  <div class="flex min-h-screen bg-slate-50/50 font-sans">
    <!-- Sidebar component -->
    <StudentSidebar />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col h-screen overflow-hidden">
      <!-- Navbar top header bar -->
      <header class="h-16 flex items-center justify-between px-4 sm:px-8 shrink-0 select-none transition-all duration-300 z-30 {scrolled ? 'bg-white border-b border-slate-100 shadow-[0_4px_20px_rgba(0,0,0,0.015)]' : 'bg-transparent border-b border-transparent'}">
        <!-- Left Section: Hamburger menu icon wrapper -->
        <div class="flex items-center gap-4">
          <button 
            aria-label="Menu"
            onclick={() => sidebarCollapsed.update(v => !v)}
            class="p-2 border border-slate-100 hover:border-slate-200 rounded-xl hover:bg-slate-50 text-slate-400 hover:text-slate-600 transition-all cursor-pointer"
          >
            <Menu class="w-5 h-5" />
          </button>
        </div>

        <!-- Right Section: Bell Notifications & User Badge Dropdown -->
        <div class="flex items-center gap-5">
          <!-- Bell Alert indicator icon -->
          <button 
            aria-label="Notifikasi"
            class="p-1.5 rounded-full hover:bg-slate-50 text-slate-400 hover:text-slate-600 transition-colors relative cursor-pointer"
          >
            <Bell class="w-5 h-5" />
            <span class="absolute top-1.5 right-1.5 w-1.5 h-1.5 bg-rose-500 rounded-full border border-white"></span>
          </button>

          <!-- User drop-down element -->
          <div class="flex items-center gap-2.5 bg-slate-50 hover:bg-slate-100/80 px-3.5 py-1.5 rounded-full border border-slate-100 transition-colors cursor-pointer">
            <div class="w-7 h-7 rounded-full bg-[#4db6ac] text-white flex items-center justify-center font-black text-xs shadow-xxs">
              {initial}
            </div>
            <span class="text-xs font-black text-slate-600">{nickname}</span>
          </div>
        </div>
      </header>

      <!-- Page Content Viewport -->
      <main onscroll={handleScroll} class="flex-1 overflow-y-auto p-4 sm:p-8 bg-slate-50/30">
        {@render children()}
      </main>
    </div>
  </div>
</ProtectedRoute>
