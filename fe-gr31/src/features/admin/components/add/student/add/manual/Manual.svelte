<script lang="ts">
  import { onMount } from 'svelte';
  import { UserPlus, Loader2, Lock } from 'lucide-svelte';
  import { createStudent, listAdmins } from '../../../../../logic/adminLogic';
  import { addToast } from '../../../../../../../stores/uiStore';
  import { currentUser } from '../../../../../../../stores/authStore';

  let { onSuccess }: { onSuccess: () => Promise<void> } = $props();

  let formSubmitting = $state(false);
  let formState = $state({
    nis: '',
    nama: '',
    kelas: '',
    walas: '',
    email: '',
    password: ''
  });

  let teachers = $state<any[]>([]);

  let isWalas = $derived(
    $currentUser?.role === 'walas' || 
    $currentUser?.role === 'guru_wali' || 
    $currentUser?.is_walas === true || 
    $currentUser?.isWalas === true
  );

  async function loadTeachers() {
    try {
      const res = await listAdmins(1, 100);
      teachers = (res.items || [])
        .filter(admin => !!admin.nama)
        .sort((a, b) => a.nama.localeCompare(b.nama));
    } catch (e) {
      console.error('Failed to load teachers:', e);
    }
  }

  onMount(async () => {
    if (isWalas) {
      formState.walas = $currentUser?.nama || '';
      formState.kelas = $currentUser?.kelas || '';
    } else {
      await loadTeachers();
    }
  });

  async function submitManualStudent() {
    if (!formState.nis || !formState.nama || !formState.kelas || !formState.password) {
      addToast('Semua field wajib diisi', 'warning');
      return;
    }
    
    formSubmitting = true;
    try {
      const payload = {
        nis: formState.nis,
        nama: formState.nama,
        kelas: formState.kelas,
        walas: formState.walas,
        agama: 'islam',
        email: formState.nis + '@student.smk31.sch.id',
        password: formState.password
      };
      
      const success = await createStudent(payload);
      if (success) {
        formState = {
          nis: '',
          nama: '',
          kelas: '',
          walas: '',
          email: '',
          password: ''
        };
        await onSuccess();
      }
    } catch (err) {
      console.error(err);
      addToast('Gagal menambahkan siswa', 'error');
    } finally {
      formSubmitting = false;
    }
  }
</script>

<div class="grid grid-cols-1 gap-6">
  <!-- Manual Form Card -->
  <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs text-left animate-fade-in">
    <div class="flex items-center gap-2.5 mb-4 border-b border-slate-55 pb-3">
      <div class="w-8 h-8 rounded-lg bg-slate-50 border border-slate-100 flex items-center justify-center text-slate-500">
        <UserPlus class="w-4 h-4" />
      </div>
      <div>
        <h3 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Informasi Siswa</h3>
      </div>
    </div>

    <form onsubmit={(e) => { e.preventDefault(); submitManualStudent(); }} class="space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <!-- NIS Input -->
        <div class="flex flex-col gap-1.5">
          <label for="nis-siswa" class="text-[11px] font-bold text-slate-700">NIS</label>
          <input 
            id="nis-siswa"
            type="text" 
            placeholder="Masukkan NIS"
            bind:value={formState.nis}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-300 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all font-mono"
            required
          />
        </div>

        <!-- Nama Input -->
        <div class="flex flex-col gap-1.5">
          <label for="nama-siswa" class="text-[11px] font-bold text-slate-700">Nama</label>
          <input 
            id="nama-siswa"
            type="text" 
            placeholder="Masukkan nama lengkap"
            bind:value={formState.nama}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-300 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all"
            required
          />
        </div>

        <!-- Kelas Input -->
        <div class="flex flex-col gap-1.5">
          <label for="kelas-siswa" class="text-[11px] font-bold text-slate-700">Kelas</label>
          <input 
            id="kelas-siswa"
            type="text" 
            placeholder="Contoh: 10A, 11B"
            bind:value={formState.kelas}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-300 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all"
            required
          />
        </div>

        <!-- Wali Kelas Input -->
        <div class="flex flex-col gap-1.5">
          <label for="walas-siswa" class="text-[11px] font-bold text-slate-700">Guru Wali</label>
          {#if isWalas}
            <div class="relative">
              <input 
                id="walas-siswa"
                type="text" 
                bind:value={formState.walas}
                class="w-full bg-slate-100 border border-slate-200 text-slate-500 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none cursor-not-allowed"
                disabled
              />
              <span class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400">
                <Lock class="w-3.5 h-3.5" />
              </span>
            </div>
          {:else}
            <select
              id="walas-siswa"
              bind:value={formState.walas}
              class="w-full bg-slate-50 border border-slate-100 focus:border-slate-300 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3 rounded-xl outline-none transition-all cursor-pointer"
            >
              <option value="">-- Pilih Guru Wali --</option>
              {#each teachers as teacher}
                <option value={teacher.nama}>{teacher.nama} ({teacher.kelas || 'Umum'})</option>
              {/each}
            </select>
          {/if}
        </div>

        <!-- Password Input -->
        <div class="flex flex-col gap-1.5 md:col-span-2">
          <label for="pass-siswa" class="text-[11px] font-bold text-slate-700">Password</label>
          <input 
            id="pass-siswa"
            type="password" 
            placeholder="Masukkan password"
            bind:value={formState.password}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-300 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all font-mono"
            required
          />
        </div>
      </div>

      <div class="flex justify-end pt-3 border-t border-slate-50">
        <!-- Submit Button -->
        <button 
          type="submit" 
          disabled={formSubmitting}
          class="flex items-center justify-center gap-1.5 px-6 py-2.5 bg-[#00a294] hover:bg-[#008f82] disabled:bg-slate-200 text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none"
        >
          {#if formSubmitting}
            <Loader2 class="w-3.5 h-3.5 animate-spin" />
            Menyimpan...
          {:else}
            <Lock class="w-3.5 h-3.5" />
            Simpan Siswa
          {/if}
        </button>
      </div>
    </form>
  </div>
</div>
