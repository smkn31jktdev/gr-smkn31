<script lang="ts">
  import { onMount } from 'svelte';
  import { listG7Rekap, getKelas } from '../../logic/adminLogic';
  import type { G7Rekap } from '../../types/admin.types';
  import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';

  const getCurrentMonthStr = () => new Date().toISOString().substring(0, 7);

  let rekaps = $state<G7Rekap[]>([]);
  let total = $state(0);
  let page = $state(1);
  let limit = $state(10);
  let loading = $state(false);

  // Filters
  let kelasList = $state<string[]>(['X LP', 'XI LP', 'XII LP']);
  let selectedKelas = $state('');
  let selectedBulan = $state(getCurrentMonthStr());
  let searchQuery = $state('');
  let selectedStatus = $state('');

  async function loadData() {
    loading = true;
    const res = await listG7Rekap({
      kelas: selectedKelas,
      bulan: selectedBulan,
      q: searchQuery,
      status: selectedStatus
    }, page, limit);
    rekaps = res.items;
    total = res.total;
    loading = false;
  }

  onMount(async () => {
    const list = await getKelas();
    if (list && list.length > 0) {
      kelasList = list;
    }
    loadData();
  });

  function handleFilter() {
    page = 1;
    loadData();
  }
</script>

<div class="space-y-4">
  <!-- Filters Bar -->
  <div class="flex flex-col sm:flex-row gap-3 items-center justify-between">
    <div class="grid grid-cols-2 sm:flex gap-2 w-full sm:w-auto">
      <input
        type="text"
        placeholder="Cari nama siswa..."
        bind:value={searchQuery}
        class="input max-w-[200px]"
        oninput={handleFilter}
      />
      <div class="w-full sm:w-[130px] text-left">
        <DropdownChoice
          options={[
            { value: '', label: 'Semua Kelas' },
            ...kelasList.map(k => ({ value: k, label: k }))
          ]}
          bind:value={selectedKelas}
          onchange={handleFilter}
          placeholder="Semua Kelas"
        />
      </div>
      <input
        type="month"
        bind:value={selectedBulan}
        onchange={handleFilter}
        class="input max-w-[160px]"
      />
      <div class="w-full sm:w-[120px] text-left">
        <DropdownChoice
          options={[
            { value: '', label: 'Semua Status' },
            { value: 'draft', label: 'Draft' },
            { value: 'reviewed', label: 'Reviewed' },
            { value: 'final', label: 'Final' }
          ]}
          bind:value={selectedStatus}
          onchange={handleFilter}
          placeholder="Semua Status"
        />
      </div>
    </div>
  </div>

  <!-- Table Card -->
  <div class="card p-0 overflow-x-auto">
    <table class="w-full text-left border-collapse text-sm">
      <thead>
        <tr class="bg-gray-50 border-b border-border text-muted font-bold">
          <th class="p-4">NIS</th>
          <th class="p-4">Nama Siswa</th>
          <th class="p-4">Kelas</th>
          <th class="p-4">Bulan</th>
          <th class="p-4">Nilai Akhir</th>
          <th class="p-4">Predikat</th>
          <th class="p-4">Status</th>
          <th class="p-4 text-center">Aksi</th>
        </tr>
      </thead>
      <tbody>
        {#if loading}
          <tr>
            <td colspan="8" class="p-8 text-center text-muted">Memuat data rekap G7...</td>
          </tr>
        {:else if rekaps.length === 0}
          <tr>
            <td colspan="8" class="p-8 text-center text-muted font-medium">Tidak ada data rekap G7 ditemukan</td>
          </tr>
        {:else}
          {#each rekaps as rekap}
            <tr class="border-b border-border hover:bg-gray-50/50">
              <td class="p-4 font-semibold text-foreground">{rekap.nis}</td>
              <td class="p-4 font-bold text-foreground">{rekap.namaSiswa}</td>
              <td class="p-4 text-muted">{rekap.kelas}</td>
              <td class="p-4 text-muted font-medium">{rekap.bulanTahun}</td>
              <td class="p-4 font-bold" class:text-emerald-600={rekap.nilaiAkhir >= 80} class:text-rose-600={rekap.nilaiAkhir < 60}>
                {rekap.nilaiAkhir.toFixed(2)}
              </td>
              <td class="p-4 font-semibold text-foreground">{rekap.predikat || 'Kurang'}</td>
              <td class="p-4">
                <span 
                  class="inline-block px-2.5 py-0.5 rounded-full text-xxs font-bold uppercase"
                  class:bg-amber-100={rekap.status === 'draft'}
                  class:text-amber-800={rekap.status === 'draft'}
                  class:bg-blue-100={rekap.status === 'reviewed'}
                  class:text-blue-800={rekap.status === 'reviewed'}
                  class:bg-emerald-100={rekap.status === 'final'}
                  class:text-emerald-800={rekap.status === 'final'}
                >
                  {rekap.status}
                </span>
              </td>
              <td class="p-4 text-center">
                <a 
                  href={`/admin/g7/${rekap.nis}/${rekap.bulanTahun}`} 
                  class="btn-primary text-xxs py-1.5 px-3 shadow-none cursor-pointer"
                >
                  {rekap.status === 'final' ? 'Lihat Detail' : 'Beri Nilai'}
                </a>
              </td>
            </tr>
          {/each}
        {/if}
      </tbody>
    </table>
  </div>
</div>
