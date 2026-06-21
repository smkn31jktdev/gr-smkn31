<script lang="ts">
  import { onMount } from 'svelte';
  import { 
    listLombaSiswa, 
    createLomba, 
    updateLomba, 
    deleteLomba,
    uploadLombaFile
  } from '../../../logic/lombaLogic';
  import { getUploadUrl } from '../../../../../api/client';
  import { addToast } from '../../../../../stores/uiStore';
  import { currentUser } from '../../../../../stores/authStore';
  import Modal from '../../../../shared/components/Modal.svelte';
  import SubmitButton from '../../../../shared/components/SubmitButton.svelte';
  import type { LombaKebersihan } from '../../../types/student.types';
  import { 
    Calendar, 
    Camera, 
    UploadCloud, 
    Loader, 
    Trash2, 
    Plus, 
    User, 
    FileText, 
    CheckCircle2, 
    AlertCircle,
    Edit3,
    ExternalLink,
    RefreshCw,
    Brush
  } from 'lucide-svelte';

  const getTodayStr = () => new Date().toLocaleDateString('sv-SE');

  let listLombas = $state<LombaKebersihan[]>([]);
  let total = $state(0);
  let page = $state(1);
  let limit = $state(10);
  let loading = $state(false);

  // User state
  let currentNIS = $derived($currentUser?.nis || '');
  let currentKelas = $derived($currentUser?.kelas || '');

  function isSameWeek(dateStr1: string, dateStr2: string): boolean {
    if (!dateStr1 || !dateStr2) return false;
    const d1 = new Date(dateStr1);
    const d2 = new Date(dateStr2);
    
    // Get Monday of d1
    const day1 = d1.getDay();
    const diff1 = d1.getDate() - day1 + (day1 === 0 ? -6 : 1);
    const monday1 = new Date(d1.setDate(diff1));
    monday1.setHours(0,0,0,0);
    
    // Get Monday of d2
    const day2 = d2.getDay();
    const diff2 = d2.getDate() - day2 + (day2 === 0 ? -6 : 1);
    const monday2 = new Date(d2.setDate(diff2));
    monday2.setHours(0,0,0,0);
    
    return monday1.getTime() === monday2.getTime();
  }

  // Check if class has uploaded for this week
  let todayStr = getTodayStr();
  let thisWeekUpload = $derived(listLombas.find(item => isSameWeek(item.tanggal, todayStr)));

  // Modal forms
  let showModal = $state(false);
  let modalMode = $state<'create' | 'edit'>('create');
  let currentId = $state('');
  let tanggal = $state(getTodayStr());
  let catatan = $state('');
  let photos = $state<string[]>([]);
  let uploadLoading = $state(false);

  async function loadData() {
    loading = true;
    const res = await listLombaSiswa('', '', page, limit);
    listLombas = res.items;
    total = res.total;
    loading = false;
  }

  onMount(() => {
    loadData();
  });

  function openCreate() {
    modalMode = 'create';
    currentId = '';
    tanggal = getTodayStr();
    catatan = '';
    photos = [];
    showModal = true;
  }

  function openEdit(item: LombaKebersihan) {
    modalMode = 'edit';
    currentId = item.id;
    tanggal = item.tanggal;
    catatan = item.catatan || '';
    photos = [...(item.foto || [])];
    showModal = true;
  }

  async function handlePhotoUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      uploadLoading = true;
      addToast('Mengunggah foto...', 'info');
      try {
        const file = input.files[0];
        const url = await uploadLombaFile(file);
        if (url) {
          photos = [...photos, url];
        }
      } finally {
        uploadLoading = false;
      }
    }
  }

  function removePhoto(index: number) {
    photos = photos.filter((_, i) => i !== index);
  }

  async function handleSave(handlers: { resolve: () => void; reject: () => void }) {
    if (photos.length === 0) {
      addToast('Wajib mengunggah minimal 1 foto kebersihan', 'warning');
      handlers.reject();
      return;
    }

    if (modalMode === 'create') {
      const alreadyHasUpload = listLombas.some(item => isSameWeek(item.tanggal, tanggal));
      if (alreadyHasUpload) {
        addToast('Kelas Anda sudah mengunggah foto kebersihan untuk minggu tersebut', 'warning');
        handlers.reject();
        return;
      }
    }

    let success = false;
    if (modalMode === 'create') {
      success = await createLomba(tanggal, photos, catatan);
    } else {
      success = await updateLomba(currentId, photos, catatan);
    }

    if (success) {
      handlers.resolve();
      showModal = false;
      loadData();
    } else {
      handlers.reject();
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Apakah Anda yakin ingin menghapus data kebersihan ini?')) return;
    const success = await deleteLomba(id);
    if (success) loadData();
  }

  function formatTanggalIndo(tglStr: string) {
    if (!tglStr) return '';
    try {
      const d = new Date(tglStr);
      return d.toLocaleDateString('id-ID', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' });
    } catch {
      return tglStr;
    }
  }
</script>

<div class="space-y-6">
  
  <!-- Header Title -->
  <div class="text-left">
    <h2 class="text-xl font-black text-slate-800 tracking-tight font-display">Lomba Kebersihan Kelas</h2>
    <p class="text-xs text-slate-500 font-bold mt-0.5">
      Kelola berkas kebersihan untuk Kelas <span class="text-[#4db6ac]">{currentKelas}</span>
    </p>
  </div>

  <!-- Weekly Status Alert -->
  {#if thisWeekUpload}
    <div class="p-5 border border-teal-100 bg-teal-50/20 rounded-3xl flex items-start gap-4 text-left">
      <div class="w-10 h-10 rounded-full bg-teal-100/50 flex items-center justify-center text-teal-600 shrink-0">
        <CheckCircle2 class="w-5 h-5" />
      </div>
      <div>
        <h4 class="text-sm font-extrabold text-teal-900">Kebersihan Minggu Ini Sudah Dilaporkan</h4>
        <p class="text-xs text-teal-700/90 mt-0.5 leading-relaxed">
          Foto kebersihan kelas Anda sudah diunggah oleh <strong>{thisWeekUpload.namaSiswa}</strong> untuk minggu ini (pada {formatTanggalIndo(thisWeekUpload.tanggal)}). Terima kasih telah menjaga kebersihan kelas!
        </p>
      </div>
    </div>
  {:else}
    <div class="p-5 border border-amber-100 bg-amber-50/20 rounded-3xl flex items-start gap-4 text-left">
      <div class="w-10 h-10 rounded-full bg-amber-100/50 flex items-center justify-center text-amber-600 shrink-0">
        <AlertCircle class="w-5 h-5" />
      </div>
      <div>
        <h4 class="text-sm font-extrabold text-amber-900">Kebersihan Minggu Ini Belum Dilaporkan</h4>
        <p class="text-xs text-amber-700/90 mt-0.5 leading-relaxed">
          Silakan unggah foto kebersihan kelas Anda minggu ini untuk berpartisipasi dalam Lomba Kebersihan Kelas. 
          Hanya satu orang perwakilan saja yang perlu mengunggah foto sekali dalam seminggu.
        </p>
        <button onclick={openCreate} class="btn-primary text-xs py-2 px-5 mt-3 shadow-sm cursor-pointer inline-flex items-center gap-1.5 self-start">
          <Plus class="w-4 h-4" /> Upload Foto Minggu Ini
        </button>
      </div>
    </div>
  {/if}

  <!-- History Card -->
  <div class="bg-white rounded-3xl border border-slate-100/90 p-6 shadow-[0_8px_30px_rgb(0,0,0,0.015)] text-left">
    <div class="flex items-center justify-between mb-6 border-b border-slate-50 pb-4">
      <div>
        <h3 class="text-sm font-black text-slate-800 tracking-tight">Riwayat Pengiriman Kelas {currentKelas}</h3>
        <p class="text-[10px] text-slate-400 font-bold mt-0.5">Daftar unggahan foto kondisi kebersihan kelas mingguan</p>
      </div>
      <div class="flex items-center gap-3">
        <button 
          onclick={loadData} 
          disabled={loading}
          class="text-xs font-black text-[#4db6ac] hover:underline cursor-pointer inline-flex items-center gap-1.5 bg-transparent border-none"
        >
          <RefreshCw class="w-3.5 h-3.5 {loading ? 'animate-spin' : ''}" />
          Segarkan
        </button>
        {#if !thisWeekUpload}
          <button onclick={openCreate} class="btn-primary text-xs py-2 px-4 shadow-sm cursor-pointer inline-flex items-center gap-1">
            <Plus class="w-3.5 h-3.5" /> Upload Foto
          </button>
        {/if}
      </div>
    </div>

    {#if loading}
      <div class="p-12 text-center">
        <Loader class="w-6 h-6 animate-spin text-[#4db6ac] mx-auto mb-2" />
        <p class="text-xs text-slate-400 font-bold">Memuat riwayat kebersihan...</p>
      </div>
    {:else if listLombas.length === 0}
      <div class="p-16 text-center max-w-md mx-auto">
        <Brush class="w-12 h-12 text-slate-300 mx-auto mb-4" />
        <h4 class="text-sm font-bold text-slate-800">Belum ada riwayat foto</h4>
        <p class="text-xs text-slate-400 mt-1 mb-5 leading-relaxed">
          Kelas Anda belum pernah mengirimkan foto kebersihan. Jadilah perwakilan kelas pertama yang mengunggah foto kebersihan minggu ini!
        </p>
        <button onclick={openCreate} class="btn-primary text-xs py-2 px-6 cursor-pointer">Mulai Unggah</button>
      </div>
    {:else}
      <div class="space-y-4">
        {#each listLombas as item}
          <div class="p-5 border border-slate-100 rounded-3xl bg-slate-50/20 space-y-4">
            <!-- Header -->
            <div class="flex items-center justify-between border-b border-slate-100/50 pb-3">
              <div>
                <span class="text-xs font-black text-slate-700">{formatTanggalIndo(item.tanggal)}</span>
                <span class="text-[9px] font-bold text-[#4db6ac] mt-0.5 flex items-center gap-1">
                  <User class="w-3 h-3" /> Oleh: {item.namaSiswa}
                </span>
              </div>
              
              {#if item.nis === currentNIS}
                <div class="flex gap-2">
                  <button onclick={() => openEdit(item)} class="p-1.5 rounded-lg border border-slate-200 text-slate-500 hover:bg-slate-50 hover:text-slate-800 transition-colors cursor-pointer bg-white" title="Edit">
                    <Edit3 class="w-3.5 h-3.5" />
                  </button>
                  <button onclick={() => handleDelete(item.id)} class="p-1.5 rounded-lg border border-rose-100 text-rose-500 hover:bg-rose-50 hover:text-rose-800 transition-colors cursor-pointer bg-white" title="Hapus">
                    <Trash2 class="w-3.5 h-3.5" />
                  </button>
                </div>
              {/if}
            </div>

            <!-- Notes -->
            {#if item.catatan}
              <div class="p-3.5 rounded-2xl bg-white border border-slate-100/80 flex items-start gap-2.5">
                <FileText class="w-4 h-4 text-slate-400 shrink-0 mt-0.5" />
                <p class="text-xs text-slate-600 font-medium leading-relaxed">{item.catatan}</p>
              </div>
            {/if}

            <!-- Photo Grid -->
            {#if item.foto && item.foto.length > 0}
              <div class="space-y-2">
                <span class="text-[9px] font-black uppercase text-slate-400 tracking-wider">Foto Kebersihan ({item.foto.length})</span>
                <div class="grid grid-cols-4 sm:grid-cols-6 gap-3">
                  {#each item.foto as photoUrl}
                    <a href={getUploadUrl(photoUrl)} target="_blank" rel="noopener noreferrer" class="aspect-square rounded-2xl overflow-hidden bg-white border border-slate-100 block hover:opacity-85 transition-opacity relative group">
                      <img src={getUploadUrl(photoUrl)} alt="Foto Kebersihan" class="w-full h-full object-cover" />
                      <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity">
                        <ExternalLink class="w-4 h-4 text-white" />
                      </div>
                    </a>
                  {/each}
                </div>
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Modal Form -->
  <Modal show={showModal} title={modalMode === 'create' ? 'Unggah Foto Kebersihan Kelas' : 'Ubah Data Foto Kebersihan'} onclose={() => showModal = false}>
    <form class="space-y-5" onsubmit={(e) => e.preventDefault()}>
      
      <!-- Date -->
      <div>
        <label for="modal-tanggal" class="text-xxs font-black uppercase tracking-wider text-slate-400 mb-1.5 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5 text-[#4db6ac]" /> Tanggal Laporan
        </label>
        <input 
          type="date" 
          id="modal-tanggal"
          bind:value={tanggal} 
          class="w-full px-4 py-2 border border-slate-200 rounded-xl text-xs font-black text-slate-700 focus:outline-none focus:border-[#4db6ac] bg-slate-50/50" 
          disabled={modalMode === 'edit'} 
        />
      </div>

      <!-- Photo upload -->
      <div>
        <label class="text-xxs font-black uppercase tracking-wider text-slate-400 mb-2 flex items-center gap-1">
          <Camera class="w-3.5 h-3.5 text-[#4db6ac]" /> Foto Kondisi Kelas
        </label>
        
        <label class="group relative flex flex-col items-center justify-center border-2 border-dashed border-slate-200 hover:border-[#4db6ac] rounded-2xl bg-slate-50/50 hover:bg-slate-50/20 p-6 text-center transition-all duration-300 cursor-pointer">
          <UploadCloud class="w-8 h-8 text-slate-400 group-hover:text-[#4db6ac] group-hover:scale-105 transition-all mb-2" />
          <span class="text-xs font-bold text-slate-600">Klik untuk unggah foto kebersihan</span>
          <span class="text-[9px] text-slate-400 mt-1">PNG / JPG (Maks. 5MB)</span>
          
          <input
            type="file"
            accept="image/*"
            onchange={handlePhotoUpload}
            disabled={uploadLoading}
            class="hidden"
          />
          
          {#if uploadLoading}
            <div class="absolute inset-0 bg-white/95 rounded-2xl flex items-center justify-center text-slate-400">
              <Loader class="w-5 h-5 animate-spin text-[#4db6ac] mr-2" />
              <span class="text-xs font-bold">Mengunggah foto...</span>
            </div>
          {/if}
        </label>

        <!-- Thumbnails -->
        {#if photos.length > 0}
          <div class="grid grid-cols-4 gap-3 mt-3">
            {#each photos as photo, index}
              <div class="relative group aspect-square rounded-xl overflow-hidden border border-slate-100 bg-slate-50">
                <img src={getUploadUrl(photo)} alt="Bukti {index + 1}" class="w-full h-full object-cover" />
                <button
                  type="button"
                  onclick={() => removePhoto(index)}
                  class="absolute inset-0 bg-rose-600/90 text-white opacity-0 group-hover:opacity-100 flex flex-col items-center justify-center gap-1 transition-opacity cursor-pointer z-10"
                >
                  <Trash2 class="w-4 h-4" />
                  <span class="text-[8px] font-black uppercase">Hapus</span>
                </button>
              </div>
            {/each}
          </div>
        {/if}
      </div>

      <!-- Notes -->
      <div>
        <label for="modal-catatan" class="text-xxs font-black uppercase tracking-wider text-slate-400 mb-1.5 flex items-center gap-1">
          <FileText class="w-3.5 h-3.5 text-[#4db6ac]" /> Catatan / Keterangan (Opsional)
        </label>
        <textarea 
          id="modal-catatan"
          placeholder="Contoh: Kelas sudah disapu, dipel, sampah dikosongkan, meja tertata rapi..." 
          bind:value={catatan} 
          class="w-full px-4 py-2 border border-slate-200 rounded-xl text-xs font-medium text-slate-700 focus:outline-none focus:border-[#4db6ac] bg-slate-50/50 min-h-[80px]"
        ></textarea>
      </div>

      <div class="flex justify-end gap-3 pt-2">
        <button type="button" onclick={() => showModal = false} class="px-5 py-2.5 border border-slate-200 rounded-xl text-xs font-bold bg-white hover:bg-slate-50 text-slate-500 transition-colors cursor-pointer">Batal</button>
        <SubmitButton label="Simpan Berkas" loadingLabel="Menyimpan..." onclick={handleSave} />
      </div>
    </form>
  </Modal>
</div>
