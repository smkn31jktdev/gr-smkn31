<script lang="ts">
  let { 
    totalSiswa = 0,
    sudahHadir = 0,
    sedangIzin = 0,
    sedangSakit = 0,
    tanpaKeterangan = 0,
    g7RataRata = 0
  } = $props<{
    totalSiswa?: number;
    sudahHadir?: number;
    sedangIzin?: number;
    sedangSakit?: number;
    tanpaKeterangan?: number;
    g7RataRata?: number;
  }>();

  // Percentage calculations
  let hadirPercent = $derived(totalSiswa > 0 ? (sudahHadir / totalSiswa) * 100 : 0);
</script>

<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
  <!-- Total Students -->
  <div class="card p-5 bg-linear-to-br from-indigo-500 to-indigo-600 text-white border-none shadow-md">
    <span class="text-xs font-bold uppercase tracking-wider text-indigo-100">Total Siswa</span>
    <div class="flex items-baseline gap-2 mt-2">
      <span class="text-3xl font-black">{totalSiswa}</span>
      <span class="text-xs font-semibold text-indigo-200">Siswa Aktif</span>
    </div>
  </div>

  <!-- Attendance rate -->
  <div class="card p-5 bg-linear-to-br from-emerald-500 to-emerald-600 text-white border-none shadow-md">
    <span class="text-xs font-bold uppercase tracking-wider text-emerald-100">Tingkat Kehadiran</span>
    <div class="flex items-baseline gap-2 mt-2">
      <span class="text-3xl font-black">{hadirPercent.toFixed(1)}%</span>
      <span class="text-xs font-semibold text-emerald-200">{sudahHadir} dari {totalSiswa} siswa</span>
    </div>
  </div>

  <!-- Permits & Sickness -->
  <div class="card p-5 bg-linear-to-br from-amber-500 to-amber-600 text-white border-none shadow-md">
    <span class="text-xs font-bold uppercase tracking-wider text-amber-100">Izin & Sakit Hari Ini</span>
    <div class="flex items-baseline gap-2 mt-2">
      <span class="text-3xl font-black">{sedangIzin + sedangSakit}</span>
      <span class="text-xs font-semibold text-amber-200">{sedangIzin} Izin | {sedangSakit} Sakit</span>
    </div>
  </div>

  <!-- G7 average done -->
  <div class="card p-5 bg-linear-to-br from-purple-500 to-purple-600 text-white border-none shadow-md">
    <span class="text-xs font-bold uppercase tracking-wider text-purple-100">Rerata G7 Harian</span>
    <div class="flex items-baseline gap-2 mt-2">
      <span class="text-3xl font-black">{g7RataRata.toFixed(1)}</span>
      <span class="text-xs font-semibold text-purple-200">dari 7.0 Kebiasaan Baik</span>
    </div>
  </div>
</div>
