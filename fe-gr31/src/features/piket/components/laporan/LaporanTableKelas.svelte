<script lang="ts">
  import type { RekapSiswaItem, RekapKelasSummary } from '../../../admin/types/admin.types';

  let { activeList, summaryByClass = [], selectedKelas = '' } = $props<{
    activeList: RekapSiswaItem[];
    summaryByClass?: RekapKelasSummary[];
    selectedKelas?: string;
  }>();

  function getKehadiranColor(rate: number): string {
    if (rate >= 80) return 'text-emerald-600';
    if (rate >= 70) return 'text-amber-600';
    return 'text-rose-600';
  }

  function getBarColor(rate: number): string {
    if (rate >= 80) return 'bg-emerald-500';
    if (rate >= 70) return 'bg-amber-500';
    return 'bg-rose-500';
  }
</script>

<!-- Rekap Per Kelas Table -->
<div class="bg-white border border-slate-100/80 rounded-2xl shadow-xs overflow-hidden">
  <div class="p-4 border-b border-slate-100 text-left flex items-center justify-between">
    <h3 class="text-xs font-bold text-slate-800 uppercase tracking-wider">
      {selectedKelas === '' ? 'Tabel Kehadiran Per Kelas' : 'Tabel Rekapitulasi Data Kehadiran'}
    </h3>
    <span class="text-[10px] text-slate-400 font-medium">
      {selectedKelas === '' ? `${summaryByClass.length} Kelas` : `${activeList.length} Siswa`}
    </span>
  </div>
  <div class="overflow-x-auto">
    {#if selectedKelas === ''}
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-slate-50/80 border-b border-slate-100 text-slate-400 font-bold text-[10px] uppercase tracking-wider">
            <th class="p-4 pl-6 w-12 text-center">No</th>
            <th class="p-4">Kelas</th>
            <th class="p-4 text-center w-16">Siswa</th>
            <th class="p-4 text-center w-20 text-emerald-600">Hadir</th>
            <th class="p-4 text-center w-20 text-amber-600">Sakit</th>
            <th class="p-4 text-center w-20 text-blue-600">Izin</th>
            <th class="p-4 text-center w-20 text-rose-600 font-bold">Alpha</th>
            <th class="p-4 text-center w-20 text-slate-500">Magang</th>
            <th class="p-4 pr-6 text-center w-40">Tingkat Kehadiran</th>
          </tr>
        </thead>
        <tbody>
          {#if summaryByClass.length === 0}
            <tr>
              <td colspan="9" class="p-8 text-center text-slate-400 font-medium">
                Tidak ada data rekap kelas ditemukan
              </td>
            </tr>
          {:else}
            {#each summaryByClass as item, idx}
              <tr class="border-b border-slate-50 hover:bg-slate-50/20 transition-colors">
                <td class="p-4 pl-6 text-center font-bold text-slate-400">{idx + 1}</td>
                <td class="p-4 font-bold text-slate-800">{item.kelas}</td>
                <td class="p-4 text-center font-semibold text-slate-550">{item.totalSiswa}</td>
                <td class="p-4 text-center font-bold text-emerald-600">{item.totalHadir}</td>
                <td class="p-4 text-center font-semibold text-amber-600">{item.totalSakit}</td>
                <td class="p-4 text-center font-semibold text-blue-600">{item.totalIzin}</td>
                <td class="p-4 text-center font-bold text-rose-600">{item.totalAlpa}</td>
                <td class="p-4 text-center font-semibold text-slate-500">{item.totalMagang}</td>
                <td class="p-4 pr-6">
                  <div class="flex items-center justify-center gap-3">
                    <span class="font-bold w-12 text-right {getKehadiranColor(item.tingkatKehadiran)}">
                      {item.tingkatKehadiran.toFixed(1)}%
                    </span>
                    <div class="w-20 h-2 bg-slate-100 rounded-full overflow-hidden">
                      <div class="h-full rounded-full {getBarColor(item.tingkatKehadiran)}" style="width: {item.tingkatKehadiran}%"></div>
                    </div>
                  </div>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    {:else}
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-slate-50/80 border-b border-slate-100 text-slate-400 font-bold text-[10px] uppercase tracking-wider">
            <th class="p-4 pl-6 w-12 text-center">No</th>
            <th class="p-4">NIS</th>
            <th class="p-4">Nama Siswa</th>
            <th class="p-4">Kelas</th>
            <th class="p-4 text-center w-16">Hadir</th>
            <th class="p-4 text-center w-16">Izin</th>
            <th class="p-4 text-center w-16">Sakit</th>
            <th class="p-4 text-center w-16">Alpa</th>
            <th class="p-4 text-center w-16">Magang</th>
            <th class="p-4 pr-6 text-center w-24">Kehadiran</th>
          </tr>
        </thead>
        <tbody>
          {#if activeList.length === 0}
            <tr>
              <td colspan="10" class="p-8 text-center text-slate-400 font-medium">
                Tidak ada data rekap laporan ditemukan untuk filter ini
              </td>
            </tr>
          {:else}
            {#each activeList as item, idx}
              <tr class="border-b border-slate-50 hover:bg-slate-50/20 transition-colors {!item.adaData ? 'opacity-60' : ''}">
                <td class="p-4 pl-6 text-center font-bold text-slate-400">{idx + 1}</td>
                <td class="p-4 text-slate-550 font-mono font-medium">{item.nis}</td>
                <td class="p-4">
                  <span class="font-bold text-slate-800 uppercase tracking-wide">{item.namaSiswa}</span>
                  {#if !item.adaData}
                    <span class="ml-1.5 text-[9px] font-bold text-slate-400 bg-slate-100 px-1.5 py-0.5 rounded-md">Belum ada data</span>
                  {/if}
                </td>
                <td class="p-4 text-slate-550 font-semibold">{item.kelas}</td>
                <td class="p-4 text-center font-bold text-emerald-600">{item.totalHadir}</td>
                <td class="p-4 text-center font-semibold text-slate-550">{item.totalIzin}</td>
                <td class="p-4 text-center font-semibold text-slate-550">{item.totalSakit}</td>
                <td class="p-4 text-center font-bold text-rose-600">{item.totalAlpa}</td>
                <td class="p-4 text-center font-semibold text-slate-550">{item.totalMagang}</td>
                <td class="p-4 pr-6 text-center font-bold {getKehadiranColor(item.tingkatKehadiran)}">
                  {item.tingkatKehadiran.toFixed(1)}%
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    {/if}
  </div>
</div>
