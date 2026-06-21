<script lang="ts">
  import { 
    User, 
    Camera, 
    Link as LinkIcon, 
    Trash2, 
    Plus, 
    UploadCloud, 
    Calendar, 
    Loader,
    Lock
  } from 'lucide-svelte';
  import { createBukti } from '../../logic/buktiLogic';
  import { uploadIzinFile } from '../../logic/kehadiranLogic';
  import { addToast } from '../../../../stores/uiStore';
  import { currentUser } from '../../../../stores/authStore';

  const getCurrentMonthStr = () => new Date().toISOString().substring(0, 7);

  let { onsuccess } = $props();

  // State
  let bulan = $state(getCurrentMonthStr());
  let linkYTInput = $state('');
  let linksYT = $state<string[]>([]);
  let photos = $state<string[]>([]);
  let uploadLoading = $state(false);
  let submitting = $state(false);

  // Derived user details
  let name = $derived($currentUser?.nama || 'Siswa');
  let kelas = $derived($currentUser?.kelas || 'Kelas');
  let nis = $derived($currentUser?.nis || '-');

  // Format month to Indonesian: "Juni 2026"
  let formattedMonth = $derived.by(() => {
    if (!bulan) return '';
    try {
      const [year, month] = bulan.split('-');
      const d = new Date(Number(year), Number(month) - 1, 1);
      return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' });
    } catch {
      return bulan;
    }
  });

  function addYTLink() {
    const trimmed = linkYTInput.trim();
    if (!trimmed) return;
    if (!trimmed.includes('youtube.com') && !trimmed.includes('youtu.be') && !trimmed.includes('instagram.com') && !trimmed.includes('tiktok.com') && !trimmed.includes('drive.google.com')) {
      addToast('Tautan media / konten tidak valid. Harap gunakan link YouTube, Instagram, TikTok, dll.', 'warning');
      return;
    }
    if (linksYT.includes(trimmed)) {
      addToast('Tautan tersebut sudah ditambahkan', 'info');
      return;
    }
    linksYT = [...linksYT, trimmed];
    linkYTInput = '';
  }

  function removeYTLink(index: number) {
    linksYT = linksYT.filter((_, i) => i !== index);
  }

  async function handlePhotoUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files.length > 0) {
      uploadLoading = true;
      addToast('Mengunggah foto...', 'info');
      try {
        const file = input.files[0];
        const url = await uploadIzinFile(file);
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

  function handleReset() {
    photos = [];
    linksYT = [];
    linkYTInput = '';
    bulan = getCurrentMonthStr();
  }

  async function handleSubmit() {
    if (photos.length === 0 && linksYT.length === 0) {
      addToast('Harap unggah minimal 1 foto atau lampirkan 1 tautan media/video', 'error');
      return;
    }

    submitting = true;
    const success = await createBukti(bulan, photos, linksYT);
    submitting = false;

    if (success) {
      handleReset();
      if (onsuccess) onsuccess();
    }
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start text-left">
  
  <!-- Left Side: Informasi Siswa (Col 5) -->
  <div class="lg:col-span-5 bg-white rounded-3xl border border-slate-100/90 p-6 shadow-[0_8px_30px_rgb(0,0,0,0.015)] space-y-6">
    <div class="flex items-center gap-3 border-b border-slate-50 pb-4">
      <div class="w-9 h-9 rounded-full bg-teal-50 flex items-center justify-center text-[#4db6ac]">
        <User class="w-5 h-5" />
      </div>
      <h3 class="text-sm font-black text-slate-800 tracking-tight">Informasi Siswa</h3>
    </div>

    <!-- Student Details -->
    <div class="space-y-4">
      <div>
        <span class="block text-[8px] font-black uppercase tracking-widest text-slate-400">Nama Lengkap</span>
        <span class="text-xs font-black text-slate-700 tracking-tight mt-0.5 block">{name}</span>
      </div>

      <div>
        <span class="block text-[8px] font-black uppercase tracking-widest text-slate-400">NIS Siswa</span>
        <span class="text-xs font-black text-slate-700 tracking-tight mt-0.5 block">{nis}</span>
      </div>

      <div>
        <span class="block text-[8px] font-black uppercase tracking-widest text-slate-400 mb-1">Kelas Akademik</span>
        <span class="inline-flex px-3.5 py-1 border border-teal-200 bg-teal-50/50 rounded-full text-[10px] font-extrabold text-[#4db6ac] shadow-xxs">
          {kelas}
        </span>
      </div>
    </div>

    <!-- Month Picker Card Section -->
    <div class="pt-5 border-t border-slate-50 space-y-3">
      <label for="bulan" class="flex items-center gap-2 text-[10px] font-black uppercase tracking-wider text-slate-400">
        <Calendar class="w-3.5 h-3.5 text-amber-500" />
        Bulan Laporan
        <span class="text-[9px] font-extrabold text-amber-600 bg-amber-50 px-1.5 py-0.5 rounded border border-amber-200/50 normal-case ml-auto">Bulan Berjalan</span>
      </label>
      <div class="relative font-sans">
        <input
          type="month"
          id="bulan"
          bind:value={bulan}
          disabled={true}
          class="w-full px-4 py-2 border border-slate-100 rounded-xl text-xs font-black text-slate-400 focus:outline-none bg-slate-50 cursor-not-allowed"
        />
        <p class="text-[10px] font-bold text-[#4db6ac] mt-1.5 flex items-center gap-1">
          <Calendar class="w-3.5 h-3.5" />
          Bukti untuk: {formattedMonth}
        </p>
      </div>
    </div>
  </div>

  <!-- Right Side: Upload and Video Link forms (Col 7) -->
  <div class="lg:col-span-7 bg-white rounded-3xl border border-slate-100/90 p-6 shadow-[0_8px_30px_rgb(0,0,0,0.015)] space-y-6">
    
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-slate-50 pb-4">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-full bg-teal-50 flex items-center justify-center text-[#4db6ac]">
          <Camera class="w-5 h-5" />
        </div>
        <h3 class="text-sm font-black text-slate-800 tracking-tight">Foto Dokumentasi</h3>
      </div>
      <span class="px-2.5 py-1 bg-teal-50 text-teal-700 text-[8px] font-black rounded-full shadow-xxs">
        {photos.length} FOTO
      </span>
    </div>

    <!-- Dropzone -->
    <div class="space-y-4">
      <label class="group relative flex flex-col items-center justify-center border-2 border-dashed border-slate-100 hover:border-[#4db6ac] rounded-2xl bg-slate-50/55 hover:bg-slate-50/20 p-8 text-center transition-all duration-300 cursor-pointer">
        <UploadCloud class="w-8 h-8 text-slate-400 group-hover:text-[#4db6ac] group-hover:scale-105 transition-all mb-3" />
        <span class="text-xs font-black text-slate-700 tracking-tight">Klik untuk unggah foto kegiatan</span>
        <span class="text-[10px] text-slate-400 font-bold mt-1">PNG / JPG (Maks. 5MB)</span>
        
        <input
          type="file"
          accept="image/*"
          onchange={handlePhotoUpload}
          disabled={uploadLoading || submitting}
          class="hidden"
        />
        
        {#if uploadLoading}
          <div class="absolute inset-0 bg-white/80 rounded-2xl flex items-center justify-center text-slate-400">
            <Loader class="w-5 h-5 animate-spin text-[#4db6ac] mr-2" />
            <span class="text-xs font-bold">Mengunggah...</span>
          </div>
        {/if}
      </label>

      <!-- Thumbnails Grid -->
      {#if photos.length > 0}
        <div class="grid grid-cols-4 sm:grid-cols-5 gap-3">
          {#each photos as photo, index}
            <div class="relative group aspect-square rounded-xl overflow-hidden border border-slate-100 bg-slate-50">
              <img src={photo} alt="Bukti {index + 1}" class="w-full h-full object-cover" />
              <!-- Hover delete button -->
              <button
                type="button"
                onclick={() => removePhoto(index)}
                class="absolute inset-0 bg-rose-600/90 text-white opacity-0 group-hover:opacity-100 flex flex-col items-center justify-center gap-1.5 transition-opacity cursor-pointer z-10"
              >
                <Trash2 class="w-4 h-4" />
                <span class="text-[9px] font-black uppercase">Hapus</span>
              </button>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <!-- YouTube Links / Content Section -->
    <div class="pt-6 border-t border-slate-50 space-y-4">
      <div class="flex items-center gap-2.5">
        <div class="w-7 h-7 rounded-full bg-teal-50 flex items-center justify-center text-[#4db6ac]">
          <LinkIcon class="w-4 h-4" />
        </div>
        <span class="text-[10px] font-black uppercase tracking-wider text-slate-400">Link Video / Konten</span>
      </div>
      
      <p class="text-[10px] text-slate-400 font-bold">
        Tautkan URL video/dokumentasi di Instagram, YouTube, TikTok, Google Drive, atau Facebook.
      </p>

      <div class="flex gap-2">
        <input
          type="url"
          placeholder="https://instagram.com/... atau platform lain"
          bind:value={linkYTInput}
          class="flex-1 px-4 py-2 border border-slate-100 rounded-xl text-xs font-bold text-slate-600 focus:outline-none focus:border-[#4db6ac] bg-slate-50/50"
        />
        <button
          type="button"
          onclick={addYTLink}
          class="px-4 py-2 bg-slate-50 hover:bg-slate-100 border border-slate-100 text-xs font-black rounded-xl text-slate-600 transition-colors cursor-pointer inline-flex items-center gap-1 shrink-0"
        >
          <Plus class="w-3.5 h-3.5" />
          Tambah
        </button>
      </div>

      <!-- Links List -->
      {#if linksYT.length > 0}
        <ul class="space-y-2">
          {#each linksYT as link, index}
            <li class="flex items-center justify-between p-3 bg-slate-50 border border-slate-100 rounded-2xl text-xs">
              <span class="truncate text-slate-500 font-bold pr-4 flex-1">{link}</span>
              <button
                type="button"
                onclick={() => removeYTLink(index)}
                class="text-rose-600 hover:text-rose-800 transition-colors cursor-pointer shrink-0"
                aria-label="Hapus tautan"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </div>

    <!-- Save & Reset buttons -->
    <div class="flex justify-end gap-3 pt-4 border-t border-slate-50">
      <button
        type="button"
        onclick={handleReset}
        disabled={submitting}
        class="px-5 py-2.5 border border-slate-100 hover:border-slate-200 text-xs font-black text-slate-500 rounded-xl hover:bg-slate-50 transition-colors cursor-pointer disabled:opacity-50"
      >
        Reset
      </button>

      <button
        type="button"
        onclick={handleSubmit}
        disabled={submitting || uploadLoading}
        class="px-8 py-2.5 bg-[#4db6ac] hover:bg-[#3ca59b] disabled:bg-slate-200 disabled:cursor-not-allowed text-white rounded-xl text-xs font-black transition-all shadow-xxs active:scale-[0.98] inline-flex items-center justify-center gap-1.5 cursor-pointer shrink-0"
      >
        {#if submitting}
          <Loader class="w-3.5 h-3.5 animate-spin" />
          Mengirim...
        {:else}
          Kirim Bukti
        {/if}
      </button>
    </div>
  </div>
</div>
