<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    Search, 
    Loader2,
    RefreshCw,
    SlidersHorizontal,
    Plus,
    X
  } from 'lucide-svelte';
  import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';
  import DatePicker from '../../../shared/components/DatePicker.svelte';
  import Table from './table/Table.svelte';
  import Count from './table/count/Count.svelte';
  import {
    loading,
    kelasList,
    selectedKelas,
    selectedTanggal,
    searchQuery,
    selectedStatus,
    showModal,
    formState,
    loadKelasJurusan,
    loadData,
    handleFilter,
    openCreate,
    handleSave
  } from '../../logic/adminKehadiranSiswaLogic';

  import { uploadIzinFileAdmin } from '../../logic/adminLogic';
  import { getUploadUrl } from '../../../../api/client';
  import { addToast } from '../../../../stores/uiStore';

  onMount(() => {
    loadKelasJurusan();
    loadData();
  });

  let isSaving = $state(false);
  let selectedFile = $state<File | null>(null);
  let isUploading = $state(false);

  async function handleFileChange(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      selectedFile = input.files[0];
      isUploading = true;
      addToast('Mengunggah berkas izin...', 'info');
      try {
        const url = await uploadIzinFileAdmin(selectedFile, $formState.nis);
        if (url) {
          $formState.fotoIzin = url;
        } else {
          selectedFile = null;
        }
      } catch (err) {
        console.error(err);
        addToast('Gagal mengunggah file', 'error');
      } finally {
        isUploading = false;
      }
    }
  }

  async function onSave(e: Event) {
    e.preventDefault();
    isSaving = true;
    await handleSave({
      resolve: () => {
        isSaving = false;
        selectedFile = null;
        showModal.set(false);
      },
      reject: () => {
        isSaving = false;
      }
    });
  }

  function closeModal() {
    selectedFile = null;
    showModal.set(false);
  }
</script>

<div class="space-y-5 select-none font-sans pb-10 text-slate-700">
  <!-- Header Section -->
  <div class="flex items-start justify-between flex-wrap gap-4 border-b border-slate-100 pb-5">
    <div class="text-left space-y-1">
      <h1 class="text-xl font-bold tracking-tight text-slate-800">Presensi Kehadiran Siswa</h1>
      <p class="text-xs text-slate-400 font-medium font-sans">Kelola kehadiran harian siswa, izin/sakit, dan alpa (tanpa keterangan).</p>
    </div>
    
    <div class="flex items-center gap-3">
      <!-- Custom Date Picker -->
      <div class="w-48 text-left">
        <DatePicker
          bind:value={$selectedTanggal}
          onchange={handleFilter}
          placeholder="Pilih tanggal"
        />
      </div>

      <!-- Segarkan -->
      <button 
        onclick={loadData}
        disabled={$loading}
        class="flex items-center justify-center p-2.5 border border-slate-200 rounded-xl bg-white hover:bg-slate-50 text-slate-500 transition-all cursor-pointer shadow-xxs"
        title="Segarkan data"
      >
        <RefreshCw class="w-4 h-4 {$loading ? 'animate-spin' : ''}" />
      </button>
    </div>
  </div>

  <!-- Search and Counters Box -->
  <div class="flex flex-col md:flex-row gap-4 items-center justify-between">
    <!-- Search Bar -->
    <div class="relative w-full md:flex-1 group">
      <span class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
        <Search class="w-4 h-4 text-slate-400 group-focus-within:text-[#00a294] transition-colors" />
      </span>
      <input
        type="text"
        placeholder="Cari nama, kelas, atau nis..."
        bind:value={$searchQuery}
        class="w-full bg-slate-50/60 hover:bg-slate-50 border border-slate-200/50 hover:border-slate-300/80 focus:border-[#00a294]/50 focus:bg-white rounded-full pl-11 pr-4 py-2.5 text-xs text-slate-700 placeholder-slate-450 outline-none transition-all focus:shadow-xs"
        oninput={handleFilter}
      />
    </div>

    <Count />
  </div>

  <!-- Filters Bar -->
  <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xs flex flex-wrap items-center gap-4 text-left">
    <div class="flex items-center gap-2 text-slate-450 text-xs font-bold uppercase tracking-wider pr-1">
      <SlidersHorizontal class="w-4 h-4 text-slate-400" />
      <span>Filter</span>
    </div>

    <!-- Kelas select pill -->
    <div class="min-w-[130px] w-full sm:w-auto text-left">
      <DropdownChoice
        options={[{ value: '', label: 'Semua Kelas' }, ...$kelasList.map(k => ({ value: k, label: `Kelas ${k}` }))]}
        bind:value={$selectedKelas}
        onchange={handleFilter}
        placeholder="Semua Kelas"
      />
    </div>

    <!-- Status select pill -->
    <div class="min-w-[130px] w-full sm:w-auto text-left">
      <DropdownChoice
        options={[
          { value: '', label: 'Semua Status' },
          { value: 'hadir', label: 'Masuk' },
          { value: 'izin', label: 'Izin' },
          { value: 'sakit', label: 'Sakit' },
          { value: 'tidak_hadir', label: 'Tanpa Keterangan' }
        ]}
        bind:value={$selectedStatus}
        onchange={handleFilter}
        placeholder="Semua Status"
      />
    </div>
  </div>

  <Table />
</div>

<!-- Modal Form Catat Absensi -->
{#if $showModal}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-xs">
    <div 
      class="relative w-full max-w-md bg-white rounded-3xl shadow-xl border border-slate-100 overflow-hidden transform transition-all text-left"
    >
      <!-- Modal Header -->
      <div class="flex items-center justify-between px-6 py-4.5 border-b border-slate-100">
        <h3 class="text-base font-bold text-slate-800">Catat Absensi Siswa</h3>
        <button 
          onclick={closeModal}
          class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-all cursor-pointer border-none bg-transparent"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Modal Body -->
      <form onsubmit={onSave} class="p-6 space-y-4.5 font-sans">
        <!-- NIS Input -->
        <div class="space-y-1.5">
          <label for="nis" class="block text-[10px] font-black text-slate-400 uppercase tracking-widest">NIS SISWA</label>
          <input
            type="text"
            id="nis"
            placeholder="Masukkan NIS siswa..."
            bind:value={$formState.nis}
            class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-xs outline-none focus:bg-white focus:border-[#00a294] transition-all"
            required
          />
        </div>

        <!-- Status Select -->
        <div class="space-y-1.5">
          <label for="status" class="block text-[10px] font-black text-slate-400 uppercase tracking-widest">STATUS KEHADIRAN</label>
          <div class="w-full text-left">
            <DropdownChoice
              options={[
                { value: 'hadir', label: 'Masuk' },
                { value: 'izin', label: 'Izin' },
                { value: 'sakit', label: 'Sakit' },
                { value: 'magang', label: 'Magang' },
                { value: 'tidak_hadir', label: 'Tanpa Keterangan' }
              ]}
              bind:value={$formState.status}
              placeholder="Pilih status"
            />
          </div>
        </div>

        <!-- Tanggal Input -->
        <div class="space-y-1.5 text-left">
          <!-- svelte-ignore a11y_label_has_associated_control -->
          <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest">TANGGAL ABSEN</label>
          <DatePicker
            bind:value={$formState.tanggal}
            placeholder="Pilih tanggal"
          />
        </div>

        <!-- Bukti Foto (Tampil Jika Masuk, Izin, atau Sakit) -->
        {#if $formState.status === 'hadir' || $formState.status === 'izin' || $formState.status === 'sakit'}
          <div class="space-y-1.5 text-left">
            <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest">UNGGAH BUKTI FOTO / SURAT BUKTI (OPSIONAL, MAKS 5MB)</label>
            <input
              type="file"
              accept="image/*,.pdf"
              onchange={handleFileChange}
              disabled={isUploading}
              class="w-full text-xs text-slate-500 file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-xs file:font-semibold file:bg-[#00a294]/10 file:text-[#00a294] hover:file:bg-[#00a294]/20 file:cursor-pointer"
            />
            {#if $formState.fotoIzin}
              <div class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-800 text-xs font-semibold flex items-center justify-between mt-2">
                <span>✓ Berkas bukti berhasil diunggah</span>
                <a href={getUploadUrl($formState.fotoIzin)} target="_blank" class="underline text-[#00a294] font-bold">Lihat Berkas</a>
              </div>
            {/if}
          </div>
        {/if}

        <!-- Alasan Textarea -->
        <div class="space-y-1.5">
          <label for="alasan" class="block text-[10px] font-black text-slate-400 uppercase tracking-widest">KETERANGAN / ALASAN (OPSIONAL)</label>
          <textarea
            id="alasan"
            placeholder="Tulis alasan jika sakit atau izin..."
            bind:value={$formState.alasan}
            rows="3"
            class="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-xs outline-none focus:bg-white focus:border-[#00a294] transition-all resize-none"
          ></textarea>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-100">
          <button
            type="button"
            onclick={closeModal}
            class="px-4.5 py-2.5 border border-slate-200 hover:bg-slate-50 text-slate-500 rounded-xl font-bold text-xs transition-all cursor-pointer bg-white"
          >
            Batal
          </button>
          <button
            type="submit"
            disabled={isSaving}
            class="flex items-center gap-1.5 px-4.5 py-2.5 bg-[#00a294] hover:bg-[#008c80] text-white rounded-xl font-bold text-xs shadow-xs transition-all cursor-pointer disabled:opacity-50 border-none"
          >
            {#if isSaving}
              <Loader2 class="w-3.5 h-3.5 animate-spin" />
              Menyimpan...
            {:else}
              Simpan Absensi
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
