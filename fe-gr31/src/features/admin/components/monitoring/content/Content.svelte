<script lang="ts">
  import { Loader2, Sun } from 'lucide-svelte';
  import { parseHabitDetails } from '../../../logic/monitoringLogic';
  import type { G7Jurnal } from '../../../../student/types/student.types';

  let { 
    loadingHistory, 
    historyItems, 
    meta, 
    kegiatan 
  } = $props<{
    loadingHistory: boolean;
    historyItems: G7Jurnal[];
    meta: { title: string; dbKey: string };
    kegiatan: string;
  }>();
</script>

<div class="p-6 max-h-[60vh] overflow-y-auto">
  {#if loadingHistory}
    <div class="py-20 flex flex-col items-center justify-center text-slate-400 gap-3">
      <Loader2 class="w-8 h-8 animate-spin text-[#00a294]" />
      <span class="text-xs font-semibold">Memuat riwayat harian...</span>
    </div>
  {:else if historyItems.length === 0}
    <div class="py-20 text-center">
      <div class="flex flex-col items-center justify-center text-slate-400 gap-3">
        <div class="p-3 bg-slate-50 border border-slate-100 rounded-2xl">
          <Sun class="w-8 h-8 text-slate-300" />
        </div>
        <p class="text-xs font-bold text-slate-400">Siswa belum memiliki riwayat pengisian pembiasaan G7.</p>
      </div>
    </div>
  {:else}
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="border-b border-slate-100">
            <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Tanggal</th>
            <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider">Rincian Pembiasaan</th>
            <th class="py-3 px-4 text-[10px] font-black text-slate-400 uppercase tracking-wider text-center">Status</th>
          </tr>
        </thead>
        <tbody>
          {#each historyItems as hist}
            {@const act = parseHabitDetails(hist, meta.dbKey)}
            <tr class="border-b border-slate-50/60 hover:bg-slate-50/30 transition-all duration-200">
              <td class="py-4 px-4 text-xs font-bold text-slate-500 font-mono">{hist.tanggal}</td>
              <td class="py-4 px-4 text-xs text-slate-600 font-semibold">
                {#if !act.done}
                  <span class="text-slate-300 font-medium">Belum mengisi</span>
                {:else}
                  {#if kegiatan === 'bangun-pagi'}
                    <div class="grid grid-cols-2 gap-x-6 gap-y-1 max-w-sm">
                      <div>Waktu Bangun: <span class="font-black text-slate-700">{act.waktu}</span></div>
                      <div>Membaca Doa: <span class="font-black {act.doa === 'Ya' ? 'text-emerald-500' : 'text-rose-500'}">{act.doa}</span></div>
                    </div>
                  {:else if kegiatan === 'beribadah'}
                    <div class="space-y-1.5 max-w-md">
                      <div>Pembiasaan: <span class="font-black text-slate-700">{act.display}</span></div>
                      {#if act.infaq && act.infaq !== '-'}
                        <div>Infaq/Zakat: <span class="font-black text-emerald-600 bg-emerald-50 px-2 py-0.5 rounded-md">{act.infaq}</span></div>
                      {/if}
                    </div>
                  {:else if kegiatan === 'makan-sehat'}
                    <div class="grid grid-cols-2 gap-x-6 gap-y-1.5 max-w-md">
                      <div>Makanan Utama: <span class="font-black text-slate-700">{act.utama}</span></div>
                      <div>Lauk Pauk: <span class="font-black text-slate-700">{act.lauk}</span></div>
                      <div>Sayur & Buah: <span class="font-black text-slate-700">{act.sayurBuah}</span></div>
                      <div>Susu & Suplemen: <span class="font-black text-slate-700">{act.susuSuplemen}</span></div>
                    </div>
                  {:else if kegiatan === 'olahraga'}
                    <div class="grid grid-cols-2 gap-x-6 gap-y-1 max-w-sm">
                      <div>Aktivitas Olahraga: <span class="font-black text-slate-700">{act.aktivitas}</span></div>
                      <div>Durasi Latihan: <span class="font-black text-slate-700">{act.durasi}</span></div>
                    </div>
                  {:else if kegiatan === 'belajar'}
                    <div class="space-y-1 max-w-md">
                      <div>Aktivitas Belajar: <span class="font-black text-slate-700">{act.display}</span></div>
                    </div>
                  {:else if kegiatan === 'bermasyarakat'}
                    <div class="grid grid-cols-2 gap-x-6 gap-y-1.5 max-w-md">
                      <div class="col-span-2">Aktivitas Sosial: <span class="font-black text-slate-700">{act.kegiatan}</span></div>
                      <div>Lokasi Kegiatan: <span class="font-black text-slate-700">{act.lokasi}</span></div>
                      <div>Diketahui RT/Ortu: <span class="font-black text-slate-700">{act.diketahuiOT}</span></div>
                    </div>
                  {:else if kegiatan === 'tidur-cukup'}
                    <div class="grid grid-cols-2 gap-x-6 gap-y-1 max-w-sm">
                      <div>Waktu Tidur: <span class="font-black text-slate-700">{act.waktu}</span></div>
                      <div>Membaca Doa: <span class="font-black {act.doa === 'Ya' ? 'text-emerald-500' : 'text-rose-500'}">{act.doa}</span></div>
                    </div>
                  {/if}
                {/if}
              </td>
              <td class="py-4 px-4 text-center">
                <span class="inline-block px-2.5 py-0.5 text-[9px] font-extrabold rounded-full uppercase tracking-wider
                  {act.done ? 'bg-emerald-50 text-emerald-600 border border-emerald-100' : 'bg-slate-50 text-slate-400 border border-slate-100'}
                ">
                  {act.done ? 'Selesai' : 'Belum'}
                </span>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
