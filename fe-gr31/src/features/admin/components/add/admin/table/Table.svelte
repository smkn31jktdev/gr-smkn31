<script lang="ts">
  import { Loader2, Trash2, UserPlus } from 'lucide-svelte';
  import {
    adminsList,
    adminsLoading,
    removeAdmin
  } from '../../../../logic/adminAdminsLogic';
  import RoleSelector from '../role/RoleSelector.svelte';
  import SearchBar from '../../../../../shared/components/SearchBar.svelte';

  let filterQuery = $state('');
  let deletingAdminId = $state<string | null>(null);
  let localToggledOnAdminIds = $state<Record<string, boolean>>({});

  let hasAnyWalas = $derived(
    $adminsList.some(a => a.isWalas) || 
    Object.values(localToggledOnAdminIds).some(Boolean)
  );

  // Filtered admins list based on search query
  let filteredAdmins = $derived(
    $adminsList.filter(admin => 
      admin.nama.toLowerCase().includes(filterQuery.toLowerCase()) ||
      admin.email.toLowerCase().includes(filterQuery.toLowerCase()) ||
      admin.role.toLowerCase().includes(filterQuery.toLowerCase())
    )
  );
</script>

<!-- Current Admins List Card -->
<div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs flex flex-col min-h-[300px] text-left">
  <!-- Card Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-55 pb-3 mb-4">
    <div>
      <h3 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Daftar Guru & Administrator</h3>
      <p class="text-[10px] text-slate-400 font-semibold mt-0.5">Kelola akun guru atau staf pembimbing yang aktif di aplikasi</p>
    </div>
    
    <!-- Search in list -->
    <SearchBar 
      bind:value={filterQuery}
      placeholder="Cari nama atau role..."
      class="w-full sm:w-64"
      size="sm"
    />
  </div>

  <!-- Admins Table -->
  {#if $adminsLoading}
    <div class="flex-1 flex flex-col items-center justify-center py-12 text-slate-400">
      <Loader2 class="w-6 h-6 animate-spin mb-2" />
      <span class="text-[11px] font-bold">Memuat data guru...</span>
    </div>
  {:else if filteredAdmins.length === 0}
    <div class="flex-1 flex flex-col items-center justify-center py-12 text-slate-400">
      <UserPlus class="w-8 h-8 text-slate-200 mb-2" />
      <span class="text-[11px] font-bold">Tidak ada staf guru ditemukan</span>
    </div>
  {:else}
    <div class="overflow-x-auto custom-scrollbar">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-slate-100">
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Nama Lengkap</th>
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Alamat Email</th>
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Wali Kelas</th>
            {#if hasAnyWalas}
              <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Kelas Wali</th>
            {/if}
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider text-center w-40">Aksi</th>
          </tr>
        </thead>
        <tbody>
          {#each filteredAdmins as admin}
            <tr class="border-b border-slate-50 hover:bg-slate-50/30 transition-colors">
              <td class="py-3 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-slate-50 border border-slate-100 text-slate-500 flex items-center justify-center font-extrabold text-[11px] shrink-0 select-none">
                    {admin.nama.charAt(0).toUpperCase()}
                  </div>
                  <span class="text-xs font-bold text-slate-755 uppercase tracking-wide">{admin.nama}</span>
                </div>
              </td>
              <td class="py-3 px-4 text-xs font-medium text-slate-500 font-mono">{admin.email}</td>
              
              <!-- Role toggling/selector cells -->
              <RoleSelector 
                {admin} 
                {hasAnyWalas} 
                bind:localToggledOnAdminIds 
              />

              <td class="py-3 px-4">
                <div class="flex items-center justify-center h-8">
                  {#if deletingAdminId === admin.id}
                    <div class="flex items-center gap-1.5 bg-rose-50 border border-rose-100 px-2 py-1 rounded-lg animate-fade-in">
                      <span class="text-[9px] font-extrabold text-rose-600 uppercase tracking-wider">Yakin?</span>
                      <button
                        onclick={async () => {
                          const success = await removeAdmin(admin.id);
                          if (success) deletingAdminId = null;
                        }}
                        class="px-2 py-0.5 bg-rose-500 hover:bg-rose-600 text-white rounded-md text-[9px] font-bold border-none cursor-pointer transition-colors"
                      >
                        Hapus
                      </button>
                      <button
                        onclick={() => deletingAdminId = null}
                        class="px-2 py-0.5 bg-white hover:bg-slate-100 text-slate-500 border border-slate-200 rounded-md text-[9px] font-bold cursor-pointer transition-colors"
                      >
                        Batal
                      </button>
                    </div>
                  {:else}
                    <button 
                      onclick={() => deletingAdminId = admin.id}
                      class="p-1.5 text-slate-400 hover:text-rose-500 hover:bg-rose-50 rounded-lg transition-all border-none cursor-pointer flex items-center justify-center"
                      title="Hapus Admin"
                    >
                      <Trash2 class="w-3.5 h-3.5" />
                    </button>
                  {/if}
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>

<style>
  /* Custom scrollbar styling for a clean sleek feel */
  .custom-scrollbar::-webkit-scrollbar {
    height: 4px;
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #e2e8f0;
    border-radius: 99px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #cbd5e1;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(2px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .animate-fade-in {
    animation: fadeIn 0.2s ease-out forwards;
  }
</style>
