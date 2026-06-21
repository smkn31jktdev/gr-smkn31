<script lang="ts">
  import { onMount } from 'svelte';
  import { Loader2, Edit3, Trash2, Lock, UserPlus } from 'lucide-svelte';
  import { deleteStudent, updateStudent, listAdmins } from '../../../../logic/adminLogic';
  import { addToast } from '../../../../../../stores/uiStore';
  import { currentUser } from '../../../../../../stores/authStore';
  import Modal from '../../../../../shared/components/Modal.svelte';
  import DropdownChoice from '../../../../../shared/components/DropdownChoice.svelte';
  import SearchBar from '../../../../../shared/components/SearchBar.svelte';
  import Pagination from '../pagination/Pagination.svelte';
  import type { Siswa } from '../../../../../auth/types/auth.types';

  let {
    students,
    loading,
    onRefresh,
    filterQuery = $bindable(),
    selectedKelas = $bindable(),
    totalStudentsCount,
    hasMore,
    page = $bindable(),
    limit = $bindable()
  }: {
    students: Siswa[];
    loading: boolean;
    onRefresh: () => Promise<void>;
    filterQuery: string;
    selectedKelas: string;
    totalStudentsCount: number;
    hasMore: boolean;
    page: number;
    limit: number;
  } = $props();

  let deletingStudentId = $state<string | null>(null);

  // Edit Modal control
  let showEditModal = $state(false);
  let editingStudentId = $state('');
  let editFormState = $state({
    nis: '',
    nama: '',
    kelas: '',
    walas: '',
    email: '',
    password: ''
  });
  let editSubmitting = $state(false);

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
    if (!isWalas) {
      await loadTeachers();
    }
  });

  function openEdit(student: Siswa) {
    editingStudentId = student.id;
    editFormState = {
      nis: student.nis,
      nama: student.nama,
      kelas: student.kelas,
      walas: student.walas || '',
      email: student.email || '',
      password: ''
    };
    showEditModal = true;
  }

  async function handleUpdateStudent() {
    if (!editFormState.nama || !editFormState.kelas) {
      addToast('Nama dan Kelas wajib diisi', 'warning');
      return;
    }
    editSubmitting = true;
    try {
      const payload: Partial<Siswa> & { password?: string } = {
        nis: editFormState.nis,
        nama: editFormState.nama,
        kelas: editFormState.kelas,
        walas: editFormState.walas
      };
      if (editFormState.password) {
        payload.password = editFormState.password;
      }
      const success = await updateStudent(editingStudentId, payload);
      if (success) {
        showEditModal = false;
        await onRefresh();
      }
    } catch (err) {
      console.error(err);
      addToast('Gagal memperbarui data siswa', 'error');
    } finally {
      editSubmitting = false;
    }
  }
</script>

<!-- Current Students List Card -->
<div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs flex flex-col min-h-[300px] text-left">
  <!-- Card Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-slate-55 pb-3 mb-4">
    <div>
      <h3 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Daftar Siswa Sekolah</h3>
      <p class="text-[10px] text-slate-400 font-semibold mt-0.5">Kelola data profile dan akun siswa yang aktif di sekolah</p>
    </div>
    
    <!-- Search & Filters in list -->
    <div class="flex flex-wrap gap-2 w-full sm:w-auto items-center">
      <!-- Filter Kelas Dropdown -->
      <div class="min-w-[130px] text-left">
        <DropdownChoice
          options={[
            { value: '', label: 'Semua Kelas' },
            { value: 'X LP', label: 'X LP' },
            { value: 'XI LP', label: 'XI LP' },
            { value: 'XII LP', label: 'XII LP' }
          ]}
          bind:value={selectedKelas}
          placeholder="Semua Kelas"
        />
      </div>

      <!-- Search Bar -->
      <SearchBar 
        bind:value={filterQuery}
        placeholder="Cari nama atau NIS..."
        class="w-full sm:w-60"
        size="sm"
      />
    </div>
  </div>

  <!-- Students Table -->
  {#if loading}
    <div class="flex-1 flex flex-col items-center justify-center py-12 text-slate-400">
      <Loader2 class="w-6 h-6 animate-spin mb-2" />
      <span class="text-[11px] font-bold">Memuat data siswa...</span>
    </div>
  {:else if students.length === 0}
    <div class="flex-1 flex flex-col items-center justify-center py-12 text-slate-400">
      <UserPlus class="w-8 h-8 text-slate-200 mb-2" />
      <span class="text-[11px] font-bold">Tidak ada siswa ditemukan</span>
    </div>
  {:else}
    <div class="overflow-x-auto custom-scrollbar">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-slate-100">
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">NIS</th>
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Nama Lengkap</th>
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Kelas</th>
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Wali Kelas (Walas)</th>
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Guru Wali (Akademik)</th>
            <th class="py-3 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider text-center w-40">Aksi</th>
          </tr>
        </thead>
        <tbody>
          {#each students as student}
            <tr class="border-b border-slate-50 hover:bg-slate-50/30 transition-colors">
              <td class="py-3 px-4 text-xs font-semibold text-slate-500 font-mono">{student.nis}</td>
              <td class="py-3 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-xl bg-slate-50 border border-slate-100 text-slate-500 flex items-center justify-center font-extrabold text-[11px] shrink-0 select-none">
                    {student.nama.charAt(0).toUpperCase()}
                  </div>
                  <span class="text-xs font-bold text-slate-755 uppercase tracking-wide">{student.nama}</span>
                </div>
              </td>
              <td class="py-3 px-4">
                <span class="px-2 py-0.5 rounded bg-slate-900 border border-slate-950 text-white text-[9.5px] font-extrabold uppercase">
                  {student.kelas}
                </span>
              </td>
              <td class="py-3 px-4 text-xs font-semibold text-[#00a294]">{student.waliKelas || '-'}</td>
              <td class="py-3 px-4 text-xs font-medium text-slate-600">{student.walas || '-'}</td>
              <td class="py-3 px-4">
                <div class="flex items-center justify-center h-8 gap-2">
                  {#if deletingStudentId === student.id}
                    <div class="flex items-center gap-1.5 bg-rose-50 border border-rose-100 px-2 py-1 rounded-lg animate-fade-in">
                      <span class="text-[9px] font-extrabold text-rose-600 uppercase tracking-wider">Yakin?</span>
                      <button
                        onclick={async () => {
                          const success = await deleteStudent(student.id);
                          if (success) {
                            deletingStudentId = null;
                            await onRefresh();
                          }
                        }}
                        class="px-2 py-0.5 bg-rose-500 hover:bg-rose-600 text-white rounded-md text-[9px] font-bold border-none cursor-pointer transition-colors"
                      >
                        Hapus
                      </button>
                      <button
                        onclick={() => deletingStudentId = null}
                        class="px-2 py-0.5 bg-white hover:bg-slate-100 text-slate-500 border border-slate-200 rounded-md text-[9px] font-bold cursor-pointer transition-colors"
                      >
                        Batal
                      </button>
                    </div>
                  {:else}
                    <button 
                      onclick={() => openEdit(student)}
                      class="p-1.5 text-slate-400 hover:text-slate-655 hover:bg-slate-50 rounded-lg transition-all border-none cursor-pointer flex items-center justify-center"
                      title="Edit Siswa"
                    >
                      <Edit3 class="w-3.5 h-3.5" />
                    </button>
                    <button 
                      onclick={() => deletingStudentId = student.id}
                      class="p-1.5 text-slate-400 hover:text-rose-55 hover:bg-rose-50 rounded-lg transition-all border-none cursor-pointer flex items-center justify-center"
                      title="Hapus Siswa"
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

  <!-- Pagination Controls -->
  <Pagination
    bind:page
    bind:limit
    {totalStudentsCount}
    {hasMore}
  />
</div>

<!-- Edit Student Modal (Restyled, Premium) -->
<Modal show={showEditModal} title="Edit Informasi Siswa" onclose={() => showEditModal = false}>
  <form class="space-y-4 text-left p-2" onsubmit={(e) => { e.preventDefault(); handleUpdateStudent(); }}>
    <div class="grid grid-cols-2 gap-4">
      <div class="flex flex-col gap-1.5">
        <label for="edit-nis" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Nomor NIS (Tetap)</label>
        <input 
          id="edit-nis"
          type="text" 
          bind:value={editFormState.nis} 
          class="w-full bg-slate-100/80 border border-slate-200 text-slate-455 text-xs font-semibold py-2 px-3 rounded-xl outline-none font-mono" 
          disabled 
        />
      </div>
      <div class="flex flex-col gap-1.5">
        <label for="edit-nama" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Nama Lengkap Siswa *</label>
        <input 
          id="edit-nama"
          type="text" 
          bind:value={editFormState.nama} 
          class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-705 text-xs font-semibold py-2 px-3 rounded-xl outline-none transition-all" 
          required 
        />
      </div>
    </div>
    
    <div class="grid grid-cols-2 gap-4">
      <div class="flex flex-col gap-1.5">
        <label for="edit-kelas" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Kelas Siswa *</label>
        <input 
          id="edit-kelas"
          type="text" 
          bind:value={editFormState.kelas} 
          class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-705 text-xs font-semibold py-2 px-3 rounded-xl outline-none transition-all" 
          required 
        />
      </div>
      <div class="flex flex-col gap-1.5">
        <label for="edit-walas" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Nama Guru Wali (Akademik)</label>
        {#if isWalas}
          <div class="relative">
            <input 
              id="edit-walas"
              type="text" 
              bind:value={editFormState.walas} 
              class="w-full bg-slate-100 border border-slate-200 text-slate-500 text-xs font-semibold py-2 px-3 rounded-xl outline-none cursor-not-allowed" 
              disabled
            />
            <span class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400">
              <Lock class="w-3.5 h-3.5" />
            </span>
          </div>
        {:else}
          <select
            id="edit-walas"
            bind:value={editFormState.walas}
            class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-705 text-xs font-semibold py-2 px-3 rounded-xl outline-none transition-all cursor-pointer"
          >
            <option value="">-- Pilih Guru Wali --</option>
            {#each teachers as teacher}
              <option value={teacher.nama}>{teacher.nama} ({teacher.kelas || 'Umum'})</option>
            {/each}
          </select>
        {/if}
      </div>
    </div>

    <div class="flex flex-col gap-1.5">
      <label for="edit-pass" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Ubah Kata Sandi (Opsional)</label>
      <input 
        id="edit-pass"
        type="password" 
        bind:value={editFormState.password} 
        placeholder="Biarkan kosong jika tidak diubah" 
        class="w-full bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2 px-3 rounded-xl outline-none transition-all font-mono" 
      />
    </div>

    <div class="flex justify-end gap-3 pt-4 border-t border-slate-50 mt-4">
      <button 
        type="button" 
        onclick={() => showEditModal = false} 
        class="px-4 py-2 border rounded-xl text-xs font-bold bg-white hover:bg-slate-50 text-slate-500 border-slate-200 transition-colors cursor-pointer"
      >
        Batal
      </button>
      
      <button 
        type="submit" 
        disabled={editSubmitting}
        class="flex items-center justify-center gap-1.5 px-5 py-2 bg-slate-800 hover:bg-slate-900 disabled:bg-slate-200 text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none"
      >
        {#if editSubmitting}
          <Loader2 class="w-3.5 h-3.5 animate-spin" />
          Menyimpan...
        {:else}
          <Lock class="w-3.5 h-3.5" />
          Simpan Perubahan
        {/if}
      </button>
    </div>
  </form>
</Modal>

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
