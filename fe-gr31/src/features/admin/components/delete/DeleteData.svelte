<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Trash2, 
    Loader2, 
    AlertTriangle, 
    Calendar, 
    User,
    X,
    Check
  } from 'lucide-svelte';
  import { 
    studentsList, 
    loadingStudents, 
    deletingData, 
    fetchStudentsForDelete, 
    deleteStudentActivities 
  } from '../../logic/deleteDataLogic';
  import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';

  let selectedNis = $state('');
  let selectedBulan = $state('');
  let showConfirm = $state(false);

  // Month list for selection
  const months = [
    { value: '2025-07', label: 'Juli 2025' },
    { value: '2025-08', label: 'Agustus 2025' },
    { value: '2025-09', label: 'September 2025' },
    { value: '2025-10', label: 'Oktober 2025' },
    { value: '2025-11', label: 'November 2025' },
    { value: '2025-12', label: 'Desember 2025' },
    { value: '2026-01', label: 'Januari 2026' },
    { value: '2026-02', label: 'Februari 2026' },
    { value: '2026-03', label: 'Maret 2026' },
    { value: '2026-04', label: 'April 2026' },
    { value: '2026-05', label: 'Mei 2026' },
    { value: '2026-06', label: 'Juni 2026' }
  ];

  onMount(() => {
    fetchStudentsForDelete();
  });

  // Reset confirm state when choices change
  $effect(() => {
    if (selectedNis || selectedBulan) {
      showConfirm = false;
    }
  });

  async function handleDelete() {
    const success = await deleteStudentActivities(selectedNis, selectedBulan);
    if (success) {
      selectedNis = '';
      selectedBulan = '';
      showConfirm = false;
    }
  }
</script>

<div class="space-y-6 select-none font-sans pb-10 text-left animate-fade-in">
  
  <!-- Header Title -->
  <div>
    <h2 class="text-xl font-extrabold tracking-tight text-slate-800 font-sans uppercase">Hapus Kegiatan Siswa</h2>
    <p class="text-xs text-slate-400 font-semibold mt-0.5">Pilih siswa dan periode bulan untuk menghapus seluruh jurnal harian dan rekap dari database secara permanen.</p>
  </div>

  <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs space-y-6">
    
    <!-- Selectors Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Student Selector -->
      <div class="flex flex-col gap-2">
        <label for="select-student" class="text-[10px] font-bold text-slate-450 uppercase tracking-widest flex items-center gap-1.5">
          <User class="w-3.5 h-3.5 text-slate-400" />
          Pilih Siswa
        </label>
        
        {#if $loadingStudents}
          <div class="flex items-center gap-2 py-2.5 px-3.5 bg-slate-50 border border-slate-100 rounded-xl text-slate-400 text-xs">
            <Loader2 class="w-3.5 h-3.5 animate-spin" />
            <span class="font-semibold">Memuat daftar siswa...</span>
          </div>
        {:else}
          <DropdownChoice
            options={$studentsList.map(student => ({
              value: student.nis,
              label: `${student.nama.toUpperCase()} (NIS: ${student.nis}) - KELAS ${student.kelas}`
            }))}
            bind:value={selectedNis}
            searchable={true}
            placeholder="Pilih Siswa"
          />
        {/if}
      </div>

      <!-- Month Selector -->
      <div class="flex flex-col gap-2">
        <label for="select-month" class="text-[10px] font-bold text-slate-450 uppercase tracking-widest flex items-center gap-1.5">
          <Calendar class="w-3.5 h-3.5 text-slate-400" />
          Pilih Bulan
        </label>
        
        <DropdownChoice
          options={months.map(m => ({ value: m.value, label: m.label }))}
          bind:value={selectedBulan}
          disabled={!selectedNis}
          placeholder="Pilih Bulan"
        />
      </div>
    </div>

    <!-- Danger Zone Panel -->
    <div class="p-5 bg-rose-50/50 border border-rose-100 rounded-2xl flex flex-col sm:flex-row items-start gap-4">
      <div class="w-10 h-10 rounded-xl bg-rose-100 flex items-center justify-center text-rose-600 shrink-0">
        <AlertTriangle class="w-5 h-5" />
      </div>
      
      <div class="flex-1 space-y-3.5">
        <div class="space-y-1">
          <h4 class="text-xs font-bold text-rose-800 uppercase tracking-wider">Area Berbahaya</h4>
          <p class="text-[11px] text-rose-600 font-semibold leading-relaxed">
            Tindakan ini akan <span class="font-bold text-rose-700">menghapus permanen</span> semua catatan kegiatan jurnal harian dan rekap penilaian bulanan untuk siswa dan periode bulan yang dipilih. Data yang telah dihapus tidak dapat dipulihkan kembali.
          </p>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          {#if !showConfirm}
            <button 
              type="button" 
              onclick={() => {
                if (!selectedNis || !selectedBulan) return;
                showConfirm = true;
              }}
              disabled={!selectedNis || !selectedBulan || $deletingData}
              class="flex items-center gap-1.5 px-5 py-2.5 bg-rose-500 hover:bg-rose-600 disabled:bg-rose-200 text-white rounded-xl font-bold text-xs shadow-xs transition-all border-none cursor-pointer"
            >
              <Trash2 class="w-4 h-4" />
              Hapus Data Kegiatan
            </button>
          {:else}
            <!-- Svelte Inline Confirmation Flow -->
            <div class="flex items-center gap-2 bg-rose-100 border border-rose-200 p-1.5 rounded-xl animate-fade-in">
              <span class="text-[10px] font-extrabold text-rose-800 uppercase tracking-wider px-2">
                Konfirmasi Hapus Permanen?
              </span>
              
              <button 
                type="button" 
                onclick={handleDelete}
                disabled={$deletingData}
                class="flex items-center gap-1 px-3 py-1.5 bg-rose-600 hover:bg-rose-700 text-white rounded-lg font-bold text-[10px] uppercase border-none cursor-pointer transition-colors"
              >
                {#if $deletingData}
                  <Loader2 class="w-3 h-3 animate-spin" />
                  Menghapus...
                {:else}
                  <Check class="w-3 h-3" />
                  Ya, Hapus
                {/if}
              </button>

              <button 
                type="button" 
                onclick={() => showConfirm = false}
                disabled={$deletingData}
                class="flex items-center gap-1 px-3 py-1.5 bg-white hover:bg-slate-100 text-slate-650 border border-slate-200 rounded-lg font-bold text-[10px] uppercase cursor-pointer transition-colors"
              >
                <X class="w-3 h-3" />
                Batal
              </button>
            </div>
          {/if}
        </div>
      </div>
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
