<script lang="ts">
  import { ChevronLeft, ChevronRight } from 'lucide-svelte';

  let {
    page = $bindable(),
    limit = $bindable(),
    totalStudentsCount,
    hasMore
  }: {
    page: number;
    limit: number;
    totalStudentsCount: number;
    hasMore: boolean;
  } = $props();
</script>

{#if limit !== -1}
  <div class="mt-4 pt-4 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-4 select-none">
    <!-- Left: Limit selector -->
    <div class="flex items-center gap-2">
      <span class="text-xs text-slate-400 font-bold">Tampilkan</span>
      <select 
        bind:value={limit} 
        class="bg-slate-50 border border-slate-200 text-slate-700 text-xs font-bold py-1.5 px-3 rounded-xl outline-none cursor-pointer focus:bg-white focus:border-slate-300 transition-colors"
      >
        <option value={50}>50</option>
        <option value={100}>100</option>
        <option value={200}>200</option>
        <option value={-1}>Seluruhnya</option>
      </select>
      <span class="text-xs text-slate-400 font-bold">per halaman</span>
    </div>

    <!-- Center: Text info -->
    <div class="text-xs text-slate-500 font-bold">
      Menampilkan {totalStudentsCount > 0 ? (page - 1) * limit + 1 : 0} - {Math.min(page * limit, totalStudentsCount)} dari {totalStudentsCount} siswa
    </div>

    <!-- Right: Prev/Next buttons -->
    <div class="flex items-center gap-2">
      <button 
        disabled={page <= 1}
        onclick={() => page = Math.max(1, page - 1)}
        class="px-3 py-1.5 bg-slate-50 border border-slate-200/50 hover:bg-slate-100 disabled:bg-slate-50 disabled:opacity-50 text-slate-655 font-bold text-xs rounded-xl transition-all cursor-pointer disabled:cursor-not-allowed border-none flex items-center gap-1.5"
      >
        <ChevronLeft class="w-3.5 h-3.5" />
        Sebelumnya
      </button>
      <button 
        disabled={!hasMore}
        onclick={() => page = page + 1}
        class="px-3 py-1.5 bg-[#00a294] text-white hover:bg-[#008f82] disabled:bg-slate-200 disabled:text-slate-400 font-bold text-xs rounded-xl transition-all cursor-pointer disabled:cursor-not-allowed border-none flex items-center gap-1.5"
      >
        Selanjutnya
        <ChevronRight class="w-3.5 h-3.5" />
      </button>
    </div>
  </div>
{:else}
  <!-- If limit is -1 (Seluruhnya) -->
  <div class="mt-4 pt-4 border-t border-slate-100 flex items-center justify-between select-none">
    <div class="flex items-center gap-2">
      <span class="text-xs text-slate-400 font-bold">Tampilkan</span>
      <select 
        bind:value={limit} 
        class="bg-slate-50 border border-slate-200 text-slate-700 text-xs font-bold py-1.5 px-3 rounded-xl outline-none cursor-pointer focus:bg-white focus:border-slate-300 transition-colors"
      >
        <option value={50}>50</option>
        <option value={100}>100</option>
        <option value={200}>200</option>
        <option value={-1}>Seluruhnya</option>
      </select>
      <span class="text-xs text-slate-400 font-bold">per halaman</span>
    </div>

    <div class="text-xs text-slate-500 font-bold">
      Menampilkan seluruh {totalStudentsCount} siswa
    </div>
    <div></div>
  </div>
{/if}
