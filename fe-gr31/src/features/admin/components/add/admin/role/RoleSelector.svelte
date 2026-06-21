<script lang="ts">
  import { updateAdminFields } from '../../../../logic/adminAdminsLogic';
  import DropdownChoice from '../../../../../shared/components/DropdownChoice.svelte';

  interface Props {
    admin: any;
    hasAnyWalas: boolean;
    localToggledOnAdminIds: Record<string, boolean>;
  }

  let {
    admin,
    hasAnyWalas,
    localToggledOnAdminIds = $bindable()
  }: Props = $props();

  let showPopover = $state(false);
  let tempSelectedClasses = $state<string[]>([]);

  const classList = [
    "X Akuntansi", "X Animasi", "X Bisnis Ritel", "X DKV", "X Layanan Perbankan", "X Manajemen Perkantoran",
    "XI Akuntansi", "XI Animasi", "XI Bisnis Ritel", "XI DKV", "XI Layanan Perbankan", "XI Manajemen Perkantoran",
    "XII Akuntansi", "XII Animasi", "XII Bisnis Ritel", "XII DKV", "XII Layanan Perbankan", "XII Manajemen Perkantoran"
  ];
</script>

<td class="py-3 px-4">
  {#if admin.role === 'super_admin'}
    <span class="px-2 py-0.5 rounded-lg text-[9px] font-extrabold uppercase bg-slate-900 text-white border border-slate-950">
      Super Admin
    </span>
  {:else}
    <div class="flex items-center gap-3">
      <!-- Toggle switch -->
      <label class="relative inline-flex items-center cursor-pointer select-none">
        <input 
          type="checkbox" 
          checked={admin.isWalas || !!localToggledOnAdminIds[admin.id]}
          onchange={async (e) => {
            const checked = e.currentTarget.checked;
            if (checked) {
              localToggledOnAdminIds[admin.id] = true;
            } else {
              delete localToggledOnAdminIds[admin.id];
              await updateAdminFields(admin.id, { isWalas: false, role: 'admin', kelas: '' });
            }
          }}
          class="sr-only peer"
        />
        <div class="w-9 h-5 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#00a294]"></div>
      </label>
      
      {#if admin.isWalas || localToggledOnAdminIds[admin.id]}
        <span class="px-2 py-0.5 rounded bg-teal-50 text-teal-700 border border-teal-100 text-[9px] font-bold uppercase shrink-0">
          Wali Kelas
        </span>
      {:else}
        <select
          value={admin.role}
          onchange={async (e) => {
            const nextRole = e.currentTarget.value;
            await updateAdminFields(admin.id, { role: nextRole, isWalas: false });
          }}
          class="bg-slate-50 border border-slate-100 text-slate-600 text-[10px] font-bold py-0.5 px-1.5 rounded-lg outline-none transition-all focus:border-slate-300 focus:bg-white cursor-pointer"
        >
          <option value="admin">Staf / Guru Non-Walas</option>
          <option value="piket">Admin Piket</option>
          <option value="admin_bk">Admin BK</option>
          <option value="super_admin">Super Admin</option>
        </select>
      {/if}
    </div>
  {/if}
</td>

{#if hasAnyWalas}
  <td class="py-3 px-4">
    {#if admin.isWalas || localToggledOnAdminIds[admin.id]}
      <div class="relative inline-block text-left">
         <button
          onclick={() => {
            if (showPopover) {
              showPopover = false;
            } else {
              showPopover = true;
              tempSelectedClasses = admin.kelas ? admin.kelas.split(', ').filter(Boolean) : [];
            }
          }}
          class="bg-white border border-slate-200 text-slate-700 text-[11px] font-bold py-1.5 px-3 rounded-xl outline-none hover:border-slate-300 transition-all cursor-pointer flex items-center justify-between gap-2 min-w-[130px]"
        >
          <span class="truncate max-w-[150px]">{admin.kelas || 'Pilih Kelas'}</span>
          <span class="text-[9px] text-slate-400">▼</span>
        </button>
        
        {#if showPopover}
          <!-- Popover Menu -->
          <div class="absolute z-50 mt-1 w-52 rounded-xl bg-white border border-slate-200 shadow-lg text-left animate-fade-in origin-top-left left-0">
            <!-- Header / Info -->
            <div class="px-3 py-1.5 border-b border-slate-100 bg-slate-50/50 rounded-t-xl flex justify-between items-center">
              <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider">Pilih Kelas Wali</span>
              <button 
                onclick={() => showPopover = false}
                class="text-slate-400 hover:text-slate-655 bg-transparent border-none cursor-pointer p-0.5 text-xs"
              >
                ✕
              </button>
            </div>
            
            <!-- Scrollable Checklist -->
            <div class="max-h-52 overflow-y-auto custom-scrollbar py-1">
              {#each classList as kl}
                <label class="flex items-center gap-2 px-3 py-1.5 hover:bg-slate-50 cursor-pointer text-xs font-semibold text-slate-700 select-none">
                  <input
                    type="checkbox"
                    checked={tempSelectedClasses.includes(kl)}
                    onchange={(e) => {
                      const checked = e.currentTarget.checked;
                      if (checked) {
                        if (!tempSelectedClasses.includes(kl)) {
                          tempSelectedClasses = [...tempSelectedClasses, kl];
                        }
                      } else {
                        tempSelectedClasses = tempSelectedClasses.filter(c => c !== kl);
                      }
                      tempSelectedClasses.sort();
                    }}
                    class="w-3.5 h-3.5 rounded border-slate-300 text-[#00a294] focus:ring-[#00a294]"
                  />
                  <span>{kl}</span>
                </label>
              {/each}
            </div>
            
            <!-- Footer Selesai -->
            <div class="border-t border-slate-100 p-1.5 bg-slate-50 rounded-b-xl flex justify-end gap-1.5">
              <button
                onclick={() => showPopover = false}
                class="px-2.5 py-1 bg-slate-200 hover:bg-slate-300 text-slate-750 text-[10px] font-bold rounded-lg cursor-pointer border-none"
              >
                Batal
              </button>
              <button
                onclick={async () => {
                  const nextKelas = tempSelectedClasses.join(', ');
                  const success = await updateAdminFields(admin.id, { isWalas: true, role: 'walas', kelas: nextKelas });
                  if (success) {
                    if (nextKelas !== '') {
                      delete localToggledOnAdminIds[admin.id];
                    }
                    showPopover = false;
                  }
                }}
                class="px-3 py-1 bg-slate-800 hover:bg-slate-900 text-white text-[10px] font-bold rounded-lg cursor-pointer border-none"
              >
                Simpan
              </button>
            </div>
          </div>
        {/if}
      </div>
    {:else}
      <span class="text-slate-300 font-semibold">-</span>
    {/if}
  </td>
{/if}
