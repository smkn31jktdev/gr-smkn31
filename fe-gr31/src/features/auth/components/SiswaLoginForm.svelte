<script lang="ts">
  import { onMount } from 'svelte';
  import { loginSiswa } from '../logic/authLogic';
  import { goto } from '$app/navigation';
  import SubmitButton from '../../shared/components/SubmitButton.svelte';

  let nis = $state('');
  let password = $state('');
  let rememberMe = $state(false);
  let errorMsg = $state('');
  let showPassword = $state(false);
  let submitBtnEl = $state<HTMLButtonElement>();

  onMount(() => {
    const remember = localStorage.getItem('siswa_remember') === 'true';
    if (remember) {
      nis = localStorage.getItem('siswa_nis') || '';
      password = localStorage.getItem('siswa_password') || '';
      rememberMe = true;
    }
  });

  async function handleSubmit(handlers: { resolve: () => void; reject: () => void }) {
    errorMsg = '';
    const success = await loginSiswa(nis, password);
    if (success) {
      if (rememberMe) {
        localStorage.setItem('siswa_remember', 'true');
        localStorage.setItem('siswa_nis', nis);
        localStorage.setItem('siswa_password', password);
      } else {
        localStorage.removeItem('siswa_remember');
        localStorage.removeItem('siswa_nis');
        localStorage.removeItem('siswa_password');
      }
      handlers.resolve();
      goto('/siswa', { replaceState: true });
    } else {
      errorMsg = 'Kombinasi NIS dan kata sandi salah. Silakan coba lagi.';
      handlers.reject();
    }
  }
</script>

<form
  class="space-y-6 w-full"
  onsubmit={(e) => e.preventDefault()}
  onkeydown={(e) => {
    if (e.key === 'Enter' && (e.target as HTMLElement).tagName !== 'BUTTON') {
      e.preventDefault();
      submitBtnEl?.click();
    }
  }}
>
  {#if errorMsg}
    <div class="p-4 rounded-2xl bg-rose-50 border border-rose-100 text-rose-600 text-xs font-bold leading-relaxed">
      {errorMsg}
    </div>
  {/if}

  <!-- NIS Field -->
  <div class="space-y-1.5 text-left">
    <label for="nis" class="block text-[11px] font-extrabold uppercase tracking-wider text-gray-500">
      NIS
    </label>
    <div class="relative">
      <!-- ID Card Icon -->
      <svg class="w-5 h-5 text-gray-400 absolute left-4 top-1/2 -translate-y-1/2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="4" width="18" height="16" rx="2" />
        <circle cx="9" cy="10" r="2" />
        <line x1="14" y1="9" x2="18" y2="9" />
        <line x1="14" y1="13" x2="18" y2="13" />
      </svg>
      <input
        type="text"
        id="nis"
        placeholder="Masukkan NIS Anda"
        bind:value={nis}
        class="w-full bg-[#f8fafc] border border-slate-200/80 rounded-2xl pl-12 pr-4 py-3.5 text-sm text-gray-800 placeholder-gray-400 outline-none focus:bg-white focus:border-[#0070f3] focus:ring-4 focus:ring-blue-100/40 transition-all duration-200"
        required
        maxlength="15"
      />
    </div>
  </div>

  <!-- Password Field -->
  <div class="space-y-1.5 text-left">
    <label for="password" class="block text-[11px] font-extrabold uppercase tracking-wider text-gray-500">
      Kata Sandi
    </label>
    <div class="relative">
      <!-- Padlock Icon -->
      <svg class="w-5 h-5 text-gray-400 absolute left-4 top-1/2 -translate-y-1/2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
        <path d="M7 11V7a5 5 0 0 1 10 0v4" />
      </svg>
      <input
        type={showPassword ? "text" : "password"}
        id="password"
        placeholder="Masukkan kata sandi"
        bind:value={password}
        class="w-full bg-[#f8fafc] border border-slate-200/80 rounded-2xl pl-12 pr-12 py-3.5 text-sm text-gray-800 placeholder-gray-400 outline-none focus:bg-white focus:border-[#0070f3] focus:ring-4 focus:ring-blue-100/40 transition-all duration-200"
        required
      />
      <!-- Eye Toggle Icon -->
      <button
        type="button"
        onclick={() => showPassword = !showPassword}
        class="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors focus:outline-none"
        aria-label={showPassword ? "Sembunyikan kata sandi" : "Tampilkan kata sandi"}
      >
        {#if showPassword}
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
            <line x1="1" y1="1" x2="23" y2="23" />
          </svg>
        {:else}
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
            <circle cx="12" cy="12" r="3" />
          </svg>
        {/if}
      </button>
    </div>
  </div>

  <!-- Remember Me Checkbox -->
  <div class="flex items-center text-left">
    <label class="flex items-center gap-2.5 cursor-pointer select-none group">
      <input 
        type="checkbox" 
        bind:checked={rememberMe} 
        class="w-4.5 h-4.5 rounded border-slate-300 text-[#0070f3] focus:ring-[#0070f3] cursor-pointer transition-colors" 
      />
      <span class="text-xs font-semibold text-slate-500 group-hover:text-slate-700 transition-colors">Ingat saya</span>
    </label>
  </div>

  <!-- Submit Button with Gradient Sky Blue/Blue Styling -->
  <div class="pt-2">
    <SubmitButton
      bind:el={submitBtnEl}
      label="Masuk"
      loadingLabel="Memproses Masuk..."
      className="w-full py-3.5 bg-linear-to-r! from-[#0070f3]! to-[#29b6f6]! hover:from-[#0060cb]! hover:to-[#02a5f4]! rounded-2xl! text-white font-bold text-sm shadow-[0_4px_15px_rgba(0,112,243,0.25)] hover:shadow-[0_6px_20px_rgba(0,112,243,0.35)] transition-all active:scale-[0.99] cursor-pointer border-none"
      onclick={handleSubmit}
    />
  </div>
</form>
