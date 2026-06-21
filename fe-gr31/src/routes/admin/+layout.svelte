<script lang="ts">
  import ProtectedRoute from '../../features/shared/components/ProtectedRoute.svelte';
  import AdminSidebar from '../../features/shared/components/sidebar/AdminSidebar.svelte';
  import { sidebarCollapsed } from '../../stores/uiStore';
  import { currentUser } from '../../stores/authStore';
  import { Menu } from 'lucide-svelte';

  let { children } = $props();

  let name = $derived($currentUser?.nama || 'Super Admin 31');
  let roleLabel = $derived($currentUser?.role === 'super_admin' ? 'Administrator' : 'Guru / Staf');
  let initial = $derived(name.charAt(0).toUpperCase());
</script>

<ProtectedRoute allowedRoles={['admin']}>
  <div class="flex min-h-screen bg-slate-50/50 font-sans">
    <!-- Sidebar -->
    <AdminSidebar />

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col h-screen overflow-hidden">
      <!-- Top header bar (navbar) -->
      <header class="h-16 border-b border-slate-100 bg-white flex items-center justify-between px-4 sm:px-8 shrink-0 select-none z-10">
        <div class="flex items-center gap-4">
          <button 
            aria-label="Menu"
            onclick={() => sidebarCollapsed.update(v => !v)}
            class="p-2 border border-slate-100 hover:border-slate-200 rounded-xl hover:bg-slate-50 text-slate-400 hover:text-slate-600 transition-all cursor-pointer bg-white"
          >
            <Menu class="w-5 h-5" />
          </button>
        </div>

        <!-- Right side: Admin profile info matching image -->
        <div class="flex items-center gap-3">
          <div class="text-right hidden sm:block">
            <h4 class="text-xs font-black text-slate-700 leading-none">{name}</h4>
            <span class="text-[9px] font-extrabold text-slate-400 uppercase tracking-wider mt-1 block">{roleLabel}</span>
          </div>
          
          <!-- Avatar icon -->
          <div class="w-9 h-9 rounded-2xl bg-blue-50 text-[#0070f3] flex items-center justify-center font-black text-xs border border-blue-100 shadow-xxs">
            {initial}
          </div>
        </div>
      </header>

      <!-- Page Content Viewport -->
      <main class="flex-1 overflow-y-auto p-4 sm:p-8 bg-slate-50/20">
        {@render children()}
      </main>
    </div>
  </div>
</ProtectedRoute>
