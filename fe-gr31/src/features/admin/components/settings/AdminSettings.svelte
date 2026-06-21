<script lang="ts">
  import { onMount } from 'svelte';
  import { currentUser } from '../../../../stores/authStore';
  import { addToast } from '../../../../stores/uiStore';
  import { 
    User, 
    Shield, 
    Lock, 
    Check, 
    Trash2, 
    AlertCircle, 
    Save,
    X
  } from 'lucide-svelte';
  import { 
    isSubmittingSettings, 
    updateAdminProfile 
  } from '../../logic/adminSettingsLogic';

  // Form states
  let nama = $state('');
  let email = $state('');
  let passwordSaatIni = $state('');
  let passwordBaru = $state('');
  let konfirmasiPassword = $state('');
  
  let fotoProfil = $state('');

  // Initialize form with current user details
  function initForm() {
    if ($currentUser) {
      nama = $currentUser.nama || '';
      email = $currentUser.email || '';
      fotoProfil = $currentUser.fotoProfil || '';
    }
    passwordSaatIni = '';
    passwordBaru = '';
    konfirmasiPassword = '';
  }

  onMount(() => {
    initForm();
  });

  function handleCancel() {
    initForm();
    addToast('Perubahan dibatalkan', 'info');
  }

  // Handle Save
  async function handleSave() {
    if (passwordSaatIni || passwordBaru || konfirmasiPassword) {
      if (!passwordSaatIni) {
        addToast('Kata sandi saat ini wajib diisi untuk mengubah kata sandi', 'warning');
        return;
      }
      if (!passwordBaru) {
        addToast('Kata sandi baru wajib diisi', 'warning');
        return;
      }
      if (passwordBaru.length < 8) {
        addToast('Kata sandi baru minimal harus 8 karakter', 'warning');
        return;
      }
      if (passwordBaru !== konfirmasiPassword) {
        addToast('Konfirmasi kata sandi baru tidak cocok', 'warning');
        return;
      }
    }

    const success = await updateAdminProfile(nama, email, fotoProfil);
    if (success) {
      passwordSaatIni = '';
      passwordBaru = '';
      konfirmasiPassword = '';
    }
  }

  function handleRemovePhoto() {
    fotoProfil = '';
    addToast('Foto profil dihapus', 'info');
  }
</script>

<div class="space-y-6 select-none font-sans pb-10 text-left animate-fade-in">
  
  <!-- Header Title -->
  <div>
    <h2 class="text-xl font-extrabold tracking-tight text-slate-800 font-sans uppercase">Pengaturan Akun</h2>
    <p class="text-xs text-slate-400 font-semibold mt-0.5">Kelola informasi profil, email, dan keamanan akun Anda.</p>
  </div>

  <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs space-y-6">
    
    <!-- Profile Card Panel -->
    <div class="flex flex-col sm:flex-row items-center gap-5 p-5 bg-slate-50/50 rounded-2xl border border-slate-100/60">
      <!-- Avatar container -->
      <div class="relative">
        {#if fotoProfil}
          <img src={fotoProfil} alt="Profil" class="w-16 h-16 rounded-full object-cover border-2 border-[#00a294]/25 shadow-xs" />
        {:else}
          <div class="w-16 h-16 rounded-full bg-[#00a294]/10 border-2 border-[#00a294]/20 flex items-center justify-center text-[#00a294] font-extrabold text-xl shadow-xs">
            {nama ? nama.charAt(0).toUpperCase() : 'A'}
          </div>
        {/if}
      </div>
      
      <!-- User info & badges -->
      <div class="flex-1 text-center sm:text-left space-y-1.5">
        <h3 class="text-sm font-extrabold text-slate-800 uppercase tracking-wide">{nama || 'Nama Admin'}</h3>
        <p class="text-[11px] font-mono text-slate-400 font-semibold">{email || 'email@domain.com'}</p>
        
        <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2 pt-0.5">
          <!-- Verified Badge -->
          <span class="flex items-center gap-1.5 px-3 py-1 rounded-full bg-emerald-50 text-emerald-600 border border-emerald-100 text-[9px] font-bold uppercase tracking-wider">
            <Check class="w-3.5 h-3.5" />
            Akun Terverifikasi
          </span>
          
          <!-- Hapus Foto Button -->
          <button 
            type="button" 
            onclick={handleRemovePhoto}
            class="flex items-center gap-1.5 px-3 py-1 rounded-full bg-rose-50 hover:bg-rose-100 text-rose-600 border border-rose-100 text-[9px] font-bold uppercase tracking-wider cursor-pointer transition-colors border-none"
          >
            <Trash2 class="w-3.5 h-3.5" />
            Hapus Foto
          </button>
        </div>
      </div>
    </div>

    <!-- Informasi Pribadi Section -->
    <div class="space-y-4">
      <div class="flex items-center gap-2 border-b border-slate-50 pb-2">
        <User class="w-4 h-4 text-slate-400" />
        <h4 class="text-[11px] font-bold text-slate-700 uppercase tracking-widest">Informasi Pribadi</h4>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <!-- Full Name input -->
        <div class="flex flex-col gap-1.5">
          <label for="fullName" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Nama Lengkap</label>
          <input 
            id="fullName"
            type="text" 
            bind:value={nama}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all"
            placeholder="Masukkan nama lengkap Anda"
          />
        </div>

        <!-- Email input -->
        <div class="flex flex-col gap-1.5">
          <label for="emailAddr" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Email</label>
          <input 
            id="emailAddr"
            type="email" 
            bind:value={email}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all font-mono"
            placeholder="Masukkan alamat email Anda"
          />
        </div>
      </div>
    </div>

    <!-- Keamanan Section -->
    <div class="space-y-4 pt-2">
      <div class="flex items-center gap-2 border-b border-slate-50 pb-2">
        <Shield class="w-4 h-4 text-slate-400" />
        <h4 class="text-[11px] font-bold text-slate-700 uppercase tracking-widest">Keamanan</h4>
      </div>

      <!-- Warning Alert Banner -->
      <div class="flex items-start gap-2.5 p-3.5 bg-amber-50/50 border border-amber-100/60 rounded-xl text-left">
        <AlertCircle class="w-4.5 h-4.5 text-amber-600 shrink-0 mt-0.5" />
        <p class="text-[10px] text-amber-700 font-semibold leading-relaxed">
          Kosongkan kolom di bawah jika Anda tidak ingin mengubah kata sandi saat ini.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <!-- Current Password -->
        <div class="flex flex-col gap-1.5">
          <label for="currPass" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Kata Sandi Saat Ini</label>
          <div class="relative">
            <Lock class="w-3.5 h-3.5 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
            <input 
              id="currPass"
              type="password" 
              bind:value={passwordSaatIni}
              class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 pl-9 pr-3.5 rounded-xl outline-none transition-all font-mono"
              placeholder="••••••••"
            />
          </div>
        </div>

        <!-- New Password -->
        <div class="flex flex-col gap-1.5">
          <label for="newPass" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Kata Sandi Baru</label>
          <input 
            id="newPass"
            type="password" 
            bind:value={passwordBaru}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all font-mono"
            placeholder="••••••••"
          />
        </div>

        <!-- Confirm Password -->
        <div class="flex flex-col gap-1.5">
          <label for="confirmPass" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Konfirmasi Kata Sandi</label>
          <input 
            id="confirmPass"
            type="password" 
            bind:value={konfirmasiPassword}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all font-mono"
            placeholder="••••••••"
          />
        </div>
      </div>
    </div>

    <!-- Action Buttons footer -->
    <div class="flex items-center justify-end gap-3 pt-4 border-t border-slate-50">
      <button 
        type="button" 
        onclick={handleCancel}
        class="flex items-center gap-1.5 px-5 py-2.5 bg-transparent hover:bg-slate-50 text-slate-500 rounded-xl font-bold text-xs border border-slate-200 cursor-pointer transition-all"
      >
        <X class="w-3.5 h-3.5" />
        Batalkan
      </button>

      <button 
        type="button" 
        onclick={handleSave}
        disabled={$isSubmittingSettings}
        class="flex items-center justify-center gap-1.5 px-6 py-2.5 bg-[#00a294] hover:bg-[#008f82] disabled:bg-slate-200 text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none"
      >
        {#if $isSubmittingSettings}
          <div class="w-3.5 h-3.5 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
          Menyimpan...
        {:else}
          <Save class="w-3.5 h-3.5" />
          Simpan Perubahan
        {/if}
      </button>
    </div>

  </div>
</div>

<style>
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(2px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .animate-fade-in {
    animation: fadeIn 0.2s ease-out forwards;
  }
</style>
