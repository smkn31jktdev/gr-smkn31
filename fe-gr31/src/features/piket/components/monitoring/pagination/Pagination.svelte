<script lang="ts">
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';
  import DropdownChoice from '../../../../shared/components/DropdownChoice.svelte';

  let {
    limit = $bindable(),
    page = $bindable(),
    total,
    hasMore,
    loading,
    onFilter,
    onLoadData
  }: {
    limit: number;
    page: number;
    total: number;
    hasMore: boolean;
    loading: boolean;
    onFilter: () => void;
    onLoadData: () => void;
  } = $props();
</script>

<div class="flex flex-col sm:flex-row items-center justify-between gap-4 bg-white border border-slate-100/80 border-t-0 rounded-b-2xl px-6 py-4 shadow-xs">
  <div class="flex items-center gap-2 text-slate-500 text-xs">
    <span>Tampilkan</span>
    <div class="w-24 text-left">
      <DropdownChoice
        options={[
          { value: 50, label: '50' },
          { value: 100, label: '100' },
          { value: 0, label: 'Semua' }
        ]}
        bind:value={limit}
        onchange={onFilter}
        placeholder="Tampilkan"
      />
    </div>
    <span>data per halaman (Total: {total} siswa)</span>
  </div>

  {#if limit > 0}
    <div class="flex items-center gap-2">
      <button 
        onclick={() => { if (page > 1) { page--; onLoadData(); } }}
        disabled={page === 1 || loading}
        class="p-2 border border-slate-200 rounded-xl bg-white hover:bg-slate-50 shadow-xxs transition-all cursor-pointer disabled:opacity-50 flex items-center justify-center border-none"
        aria-label="Sebelumnya"
      >
        <ChevronLeft class="w-4 h-4 {page === 1 || loading ? 'text-slate-300' : 'text-black'}" />
      </button>
      <span class="text-xs font-bold text-slate-700 px-2">Halaman {page}</span>
      <button 
        onclick={() => { if (hasMore) { page++; onLoadData(); } }}
        disabled={!hasMore || loading}
        class="p-2 border border-slate-200 rounded-xl bg-white hover:bg-slate-50 shadow-xxs transition-all cursor-pointer disabled:opacity-50 flex items-center justify-center border-none"
        aria-label="Selanjutnya"
      >
        <ChevronRight class="w-4 h-4 {!hasMore || loading ? 'text-slate-300' : 'text-black'}" />
      </button>
    </div>
  {/if}
</div>
