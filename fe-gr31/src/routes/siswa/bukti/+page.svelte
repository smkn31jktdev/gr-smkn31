<script lang="ts">
  import { onMount } from 'svelte';
  import BuktiUploadForm from '../../../features/student/components/bukti/BuktiUploadForm.svelte';
  import { listBuktiSiswa } from '../../../features/student/logic/buktiLogic';
  import type { Bukti } from '../../../features/student/types/student.types';
  import { RefreshCw, Video, Image as ImageIcon, ExternalLink } from 'lucide-svelte';

  let listBuktis = $state<Bukti[]>([]);
  let loadingHistory = $state(false);

  async function loadHistory() {
    loadingHistory = true;
    const res = await listBuktiSiswa();
    listBuktis = res.items;
    loadingHistory = false;
  }

  function formatBulan(bulanStr: string) {
    if (!bulanStr) return '';
    try {
      const [year, month] = bulanStr.split('-');
      const d = new Date(Number(year), Number(month) - 1, 1);
      return d.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' });
    } catch {
      return bulanStr;
    }
  }

  onMount(() => {
    loadHistory();
  });
</script>

<div class="space-y-6">
  <!-- Header Title -->
  <div class="text-left">
    <h2 class="text-xl font-black text-slate-800 tracking-tight font-display">Dokumentasi Kegiatan</h2>
    <p class="text-xs text-slate-500 font-bold mt-0.5">Unggah bukti foto dan tautan video kegiatan bulanan Anda untuk divalidasi oleh Wali Kelas.</p>
  </div>

  <!-- Upload Form -->
  <BuktiUploadForm onsuccess={loadHistory} />

  <!-- History Card -->
  <div class="bg-white rounded-3xl border border-slate-100/90 p-6 shadow-[0_8px_30px_rgb(0,0,0,0.015)] text-left">
    <div class="flex items-center justify-between mb-6 border-b border-slate-50 pb-4">
      <div>
        <h3 class="text-sm font-black text-slate-800 tracking-tight">Riwayat Unggahan Bukti</h3>
        <p class="text-[10px] text-slate-400 font-bold mt-0.5">Daftar laporan berkas per bulan yang telah diunggah</p>
      </div>
      <button 
        onclick={loadHistory} 
        disabled={loadingHistory}
        class="text-xs font-black text-[#4db6ac] hover:underline cursor-pointer inline-flex items-center gap-1.5 bg-transparent border-none"
      >
        <RefreshCw class="w-3.5 h-3.5 {loadingHistory ? 'animate-spin' : ''}" />
        Segarkan
      </button>
    </div>

    {#if loadingHistory}
      <p class="text-xs text-slate-400 text-center py-12 font-medium">Memuat riwayat unggahan...</p>
    {:else if listBuktis.length === 0}
      <p class="text-xs text-slate-400 text-center py-12 font-medium">Belum ada bukti kegiatan yang diunggah.</p>
    {:else}
      <div class="space-y-4">
        {#each listBuktis as item}
          <div class="p-5 border border-slate-100 rounded-3xl bg-slate-50/20 space-y-4">
            <div class="flex items-center justify-between border-b border-slate-50 pb-3">
              <span class="text-xs font-black text-slate-700">Laporan Bulan: {formatBulan(item.bulan)}</span>
            </div>

            <!-- YouTube Video Links -->
            {#if item.linkYT && item.linkYT.length > 0}
              <div class="space-y-2">
                <h4 class="text-[10px] font-black uppercase tracking-wider text-slate-400 flex items-center gap-1.5">
                  <Video class="w-3.5 h-3.5 text-[#4db6ac]" />
                  Tautan Video / Konten
                </h4>
                <ul class="space-y-1.5">
                  {#each item.linkYT as link}
                    <li class="text-xs text-[#4db6ac] hover:underline truncate inline-flex items-center gap-1.5">
                      <ExternalLink class="w-3 h-3 shrink-0" />
                      <a href={link} target="_blank" rel="noopener noreferrer">{link}</a>
                    </li>
                  {/each}
                </ul>
              </div>
            {/if}

            <!-- Uploaded Photos Grid -->
            {#if item.foto && item.foto.length > 0}
              <div class="space-y-2">
                <h4 class="text-[10px] font-black uppercase tracking-wider text-slate-400 flex items-center gap-1.5">
                  <ImageIcon class="w-3.5 h-3.5 text-[#4db6ac]" />
                  Foto Kegiatan
                </h4>
                <div class="grid grid-cols-4 sm:grid-cols-6 gap-3">
                  {#each item.foto as imgUrl}
                    <a 
                      href={imgUrl} 
                      target="_blank" 
                      rel="noopener noreferrer" 
                      class="aspect-square rounded-2xl overflow-hidden border border-slate-100 bg-white block hover:opacity-85 transition-opacity"
                    >
                      <img src={imgUrl} alt="Bukti Foto" class="w-full h-full object-cover" />
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
</div>
