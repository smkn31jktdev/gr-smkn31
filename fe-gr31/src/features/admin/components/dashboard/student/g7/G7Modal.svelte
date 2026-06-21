<script lang="ts">
  import { Loader2, Download } from 'lucide-svelte';
  import { 
    downloadStudentReportPDF, 
    downloadStudentSemesterReportPDF 
  } from '../../../../logic/adminDashboardLogic';
  import G7Summary from './summary/G7Summary.svelte';
  import G7Rating from './rating/G7Rating.svelte';

  let { 
    open = $bindable(), 
    selectedDetailStudent, 
    detailType, 
    detailLoading, 
    detailRekap, 
    detailEvaluate 
  }: {
    open: boolean,
    selectedDetailStudent: any,
    detailType: 'bulanan' | 'semester',
    detailLoading: boolean,
    detailRekap: any,
    detailEvaluate: any
  } = $props();

  let activeModalTab = $state<'rangkuman' | 'grafik'>('rangkuman');
</script>

{#if open && selectedDetailStudent}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm transition-all duration-300">
    <div class="bg-white w-full max-w-3xl rounded-3xl border border-slate-100 shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-200">
      <!-- Modal Header -->
      <div class="bg-slate-50/50 border-b border-slate-100 p-6 flex justify-between items-start">
        <div class="text-left">
          <h3 class="text-base font-black text-slate-700 font-display">
            {activeModalTab === 'rangkuman' ? 'Rangkuman Kebiasaan' : 'Grafik Rating Kebiasaan'}
          </h3>
          <p class="text-xs text-slate-400 font-bold mt-1 uppercase font-display">
            {selectedDetailStudent.nama} &bull; {detailRekap?.kelas || selectedDetailStudent.kelas || 'X Layanan Perbankan'}
          </p>
        </div>
        <button 
          onclick={() => { open = false; }}
          aria-label="Tutup"
          class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-all cursor-pointer border-none bg-transparent"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Tab Switcher -->
      <div class="px-6 flex border-b border-slate-100 bg-slate-50/20 shrink-0 select-none">
        <button 
          onclick={() => activeModalTab = 'rangkuman'}
          class="px-5 py-3 text-xs font-bold transition-all border-b-2 bg-transparent cursor-pointer {activeModalTab === 'rangkuman' ? 'border-[#00a294] text-[#00a294]' : 'border-transparent text-slate-400 hover:text-slate-600'}"
        >
          Rangkuman Kebiasaan
        </button>
        <button 
          onclick={() => activeModalTab = 'grafik'}
          class="px-5 py-3 text-xs font-bold transition-all border-b-2 bg-transparent cursor-pointer {activeModalTab === 'grafik' ? 'border-[#00a294] text-[#00a294]' : 'border-transparent text-slate-400 hover:text-slate-600'}"
        >
          Grafik Rating
        </button>
      </div>

      <!-- Modal Body -->
      <div class="p-6 max-h-[60vh] overflow-y-auto custom-scrollbar">
        {#if detailLoading}
          <div class="py-20 flex flex-col items-center justify-center text-slate-400 gap-3">
            <Loader2 class="w-8 h-8 animate-spin text-[#00a294]" />
            <span class="text-xs font-semibold">Memuat rincian indikator...</span>
          </div>
        {:else if (detailType === 'bulanan' && !detailEvaluate) || (detailType === 'semester' && !detailRekap)}
          <div class="py-20 text-center text-slate-400 text-xs font-semibold">
            Belum ada data jurnal G7 atau penilaian untuk periode ini.
          </div>
        {:else}
          {#if activeModalTab === 'rangkuman'}
            <G7Summary {detailType} {detailEvaluate} {detailRekap} />
          {:else}
            <G7Rating {detailType} {detailEvaluate} {detailRekap} />
          {/if}
        {/if}
      </div>

      <!-- Modal Footer -->
      <div class="bg-slate-50/50 border-t border-slate-100 p-5 flex justify-end gap-3 select-none">
        <button 
          onclick={() => { open = false; }}
          class="px-5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-700 text-xs font-bold rounded-xl transition-all cursor-pointer border-none"
        >
          Tutup
        </button>
        {#if !detailLoading}
          <button 
            onclick={() => {
              if (detailType === 'bulanan') {
                downloadStudentReportPDF(selectedDetailStudent.nis, selectedDetailStudent.nama, selectedDetailStudent.bulan);
              } else {
                downloadStudentSemesterReportPDF(
                  selectedDetailStudent.nis, 
                  selectedDetailStudent.nama, 
                  selectedDetailStudent.semester, 
                  detailRekap?.nilaiAkhir ?? 0, 
                  detailRekap?.predikat ?? 'Belum Dinilai'
                );
              }
              open = false;
            }}
            class="flex items-center gap-1.5 px-5 py-2.5 bg-[#00a294] hover:bg-[#008c80] text-white text-xs font-bold rounded-xl shadow-xs transition-all cursor-pointer border-none"
          >
            <Download class="w-3.5 h-3.5" />
            Download PDF
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 99px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
  }
</style>

