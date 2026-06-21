<script lang="ts">
  import { HelpCircle, RefreshCw, Loader2, Lock } from 'lucide-svelte';
  import {
    sheetUrl,
    sheetLoading,
    sheetAdminsList,
    sheetSubmitting,
    loadSheetData,
    submitSheetAdmins
  } from '../../../../../logic/adminAdminsLogic';

  // Local role labels helper
  function getRoleLabel(role: string) {
    switch (role) {
      case 'super_admin': return 'Super Admin';
      case 'guru_bk':
      case 'admin_bk': return 'Admin BK';
      case 'guru_wali':
      case 'walas': return 'Wali Kelas';
      case 'piket': return 'Admin Piket';
      case 'admin': return 'Staf / Guru Non-Walas';
      default: return role;
    }
  }
</script>

<div class="grid grid-cols-1 gap-6 animate-fade-in">
  
  <!-- Sheets Control Card -->
  <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs text-left space-y-5">
    <!-- Panduan Import -->
    <div class="p-4 rounded-xl bg-slate-50 border border-slate-100/60 text-slate-500 text-xs leading-relaxed space-y-2">
      <div class="flex items-center gap-1.5 font-bold text-slate-700">
        <HelpCircle class="w-4 h-4 text-[#00a294]" />
        <span>Panduan Import Google Sheets</span>
      </div>
      <ol class="list-decimal pl-4 space-y-1.5 font-semibold text-[11px] text-slate-500">
        <li>Atur Google Sheet Anda agar dapat diakses publik (Akses Link: <span class="text-[#00a294]">Anyone with the link / Siapa saja dengan link</span>).</li>
        <li>Pastikan kolom baris pertama diisi dengan header persis berikut: <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Nama</span>, <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Email</span>, <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Role</span>, <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Password</span>.</li>
        <li>Isi nilai kolom <span class="font-bold">Role</span> dengan: <span class="font-mono bg-white px-1 border border-slate-150 rounded">super_admin</span>, <span class="font-mono bg-white px-1 border border-slate-150 rounded">guru_bk</span>, <span class="font-mono bg-white px-1 border border-slate-150 rounded">guru_wali</span>, atau <span class="font-mono bg-white px-1 border border-slate-150 rounded">admin</span>.</li>
      </ol>
    </div>

    <!-- URL Input & Fetch Controls -->
    <div class="flex flex-col gap-2">
      <label for="sheet-url" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Tautan / URL Spreadsheet Google</label>
      <div class="flex flex-col sm:flex-row gap-3">
        <input 
          id="sheet-url"
          type="text" 
          placeholder="Contoh: https://docs.google.com/spreadsheets/d/.../edit?usp=sharing"
          bind:value={$sheetUrl}
          class="flex-1 bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all"
        />
        
        <div class="flex gap-2 shrink-0">
          <!-- Load Button -->
          <button 
            onclick={loadSheetData}
            disabled={$sheetLoading}
            class="flex items-center justify-center gap-1.5 px-4 py-2.5 bg-slate-50 hover:bg-slate-100 disabled:bg-slate-100 text-slate-600 border border-slate-200/50 rounded-xl font-bold text-xs transition-all cursor-pointer"
          >
            {#if $sheetLoading}
              <Loader2 class="w-3.5 h-3.5 animate-spin" />
              Mengunduh...
            {:else}
              <RefreshCw class="w-3.5 h-3.5" />
              Muat Data Sheet
            {/if}
          </button>

          <!-- Save All Button -->
          {#if $sheetAdminsList.length > 0}
            <button 
              onclick={submitSheetAdmins}
              disabled={$sheetSubmitting}
              class="flex items-center justify-center gap-1.5 px-5 py-2.5 bg-slate-800 hover:bg-slate-900 disabled:bg-slate-250 text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none"
            >
              {#if $sheetSubmitting}
                <Loader2 class="w-3.5 h-3.5 animate-spin" />
                Menyimpan...
              {:else}
                <Lock class="w-3.5 h-3.5" />
                Simpan Semua ({$sheetAdminsList.length})
              {/if}
            </button>
          {/if}
        </div>
      </div>
    </div>
  </div>

  <!-- Preview Table Card -->
  {#if $sheetAdminsList.length > 0}
    <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs min-h-[250px] flex flex-col text-left">
      <!-- Header -->
      <div class="border-b border-slate-50 pb-3 mb-4 flex justify-between items-center">
        <div>
          <h3 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Preview Baris Google Sheet</h3>
          <p class="text-[10px] text-slate-400 font-semibold mt-0.5">Tinjau list akun sebelum disimpan permanen ke database</p>
        </div>
        <span class="px-2.5 py-0.5 rounded-lg bg-slate-50 border border-slate-100 text-[10px] font-extrabold text-slate-500 font-mono">
          {$sheetAdminsList.length} Baris
        </span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100">
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">No</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Nama Lengkap</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Alamat Email</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Role Akses</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Kata Sandi</th>
            </tr>
          </thead>
          <tbody>
            {#each $sheetAdminsList as admin, idx}
              <tr class="border-b border-slate-50 hover:bg-slate-50/20 transition-colors">
                <td class="py-2.5 px-4 text-xs font-bold text-slate-400 font-mono">{idx + 1}</td>
                <td class="py-2.5 px-4 text-xs font-bold text-slate-700 uppercase">{admin.nama}</td>
                <td class="py-2.5 px-4 text-xs font-medium text-slate-500 font-mono">{admin.email}</td>
                <td class="py-2.5 px-4">
                  <span class="px-2 py-0.5 rounded bg-slate-50 text-slate-500 border border-slate-100 text-[9px] font-bold uppercase">
                    {getRoleLabel(admin.role)}
                  </span>
                </td>
                <td class="py-2.5 px-4 text-xs font-medium text-slate-400 font-mono">{admin.password}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}

</div>
