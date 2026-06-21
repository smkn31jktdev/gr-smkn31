<script lang="ts">
  import { UserPlus, Loader2, Lock } from 'lucide-svelte';
  import {
    formNama,
    formEmail,
    formRole,
    formKelas,
    formPassword,
    formSubmitting,
    submitManualAdmin
  } from '../../../../../logic/adminAdminsLogic';
  import DropdownChoice from '../../../../../../shared/components/DropdownChoice.svelte';

  let isGuru = $state(false);
  let tipeGuru = $state<'wali' | 'non_wali'>('wali');

  $effect(() => {
    if (isGuru) {
      if (tipeGuru === 'wali') {
        $formRole = 'walas';
        if (!$formKelas) {
          $formKelas = '';
        }
      } else {
        $formRole = 'admin';
        $formKelas = '';
      }
    } else {
      $formKelas = '';
      if ($formRole === 'guru_wali') {
        $formRole = 'admin';
      }
    }
  });
</script>

<div class="grid grid-cols-1 gap-6">
  <!-- Manual Form Card -->
  <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs text-left animate-fade-in">
    <div class="flex items-center gap-2.5 mb-4 border-b border-slate-50 pb-3">
      <div class="w-8 h-8 rounded-lg bg-slate-50 border border-slate-100 flex items-center justify-center text-slate-500">
        <UserPlus class="w-4 h-4" />
      </div>
      <div>
        <h3 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Registrasi Guru / Admin Baru</h3>
        <p class="text-[10px] text-slate-400 font-semibold mt-0.5">Buat akun guru kelas, bimbingan konseling, atau administrator baru secara langsung.</p>
      </div>
    </div>

    <form onsubmit={(e) => { e.preventDefault(); submitManualAdmin(); }} class="space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Nama Input -->
        <div class="flex flex-col gap-1.5">
          <label for="nama-guru" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Nama Lengkap Guru</label>
          <input 
            id="nama-guru"
            type="text" 
            placeholder="Masukkan nama guru"
            bind:value={$formNama}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-300 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all"
            required
          />
        </div>

        <!-- Email Input -->
        <div class="flex flex-col gap-1.5">
          <label for="email-guru" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Alamat Email Resmi</label>
          <input 
            id="email-guru"
            type="email" 
            placeholder="Masukkan email"
            bind:value={$formEmail}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all font-mono"
            required
          />
        </div>

        <!-- Password Input -->
        <div class="flex flex-col gap-1.5">
          <label for="pass-guru" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Kata Sandi Akun</label>
          <input 
            id="pass-guru"
            type="password" 
            placeholder="Masukkan password"
            bind:value={$formPassword}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all font-mono"
            required
          />
        </div>
      </div>

      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pt-3 border-t border-slate-55">
        <!-- Toggle isGuru -->
        <div class="flex flex-wrap items-center gap-6">
          <div class="flex items-center gap-2.5">
            <span class="text-[11px] font-bold text-slate-700">Akun Wali Kelas?</span>
            <label class="relative inline-flex items-center cursor-pointer select-none">
              <input 
                type="checkbox" 
                bind:checked={isGuru}
                class="sr-only peer"
              />
              <div class="w-9 h-5 bg-slate-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-[#00a294]"></div>
            </label>
          </div>

          <!-- Conditional Class Selector field -->
          {#if isGuru}
            <div class="flex items-center gap-2 animate-fade-in min-w-[180px] text-left">
              <span class="text-[11px] font-bold text-slate-700 shrink-0">Kelas Wali:</span>
              <DropdownChoice
                options={[
                  "X Akuntansi", "X Animasi", "X Bisnis Ritel", "X DKV", "X Layanan Perbankan", "X Manajemen Perkantoran",
                  "XI Akuntansi", "XI Animasi", "XI Bisnis Ritel", "XI DKV", "XI Layanan Perbankan", "XI Manajemen Perkantoran",
                  "XII Akuntansi", "XII Animasi", "XII Bisnis Ritel", "XII DKV", "XII Layanan Perbankan", "XII Manajemen Perkantoran"
                ]}
                bind:value={$formKelas}
                placeholder="Pilih Kelas"
              />
            </div>
          {:else}
            <div class="flex items-center gap-2.5 min-w-[180px] text-left">
              <span class="text-[11px] font-bold text-slate-700 shrink-0">Tingkat Akses:</span>
              <DropdownChoice
                options={[
                  { value: 'admin', label: 'Staf Admin' },
                  { value: 'piket', label: 'Admin Piket' },
                  { value: 'admin_bk', label: 'Admin BK' },
                  { value: 'super_admin', label: 'Super Admin' }
                ]}
                bind:value={$formRole}
                placeholder="Pilih Akses"
              />
            </div>
          {/if}
        </div>

        <!-- Submit Button -->
        <button 
          type="submit" 
          disabled={$formSubmitting}
          class="flex items-center justify-center gap-1.5 px-6 py-2.5 bg-slate-800 hover:bg-slate-900 disabled:bg-slate-200 text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none"
        >
          {#if $formSubmitting}
            <Loader2 class="w-3.5 h-3.5 animate-spin" />
            Menyimpan...
          {:else}
            <Lock class="w-3.5 h-3.5" />
            Simpan Admin Baru
          {/if}
        </button>
      </div>
    </form>
  </div>
</div>
