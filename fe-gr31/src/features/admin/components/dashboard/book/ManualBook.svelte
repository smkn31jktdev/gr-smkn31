<script lang="ts">
  import { BookOpen, Download } from 'lucide-svelte';
  import { downloadManualBookPDF } from '../../../logic/adminDashboardLogic';
</script>

<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white border border-slate-100 rounded-2xl p-5 shadow-xs">
  <div class="flex items-center gap-3.5 text-left shrink-0">
    <div class="w-10 h-10 rounded-xl bg-slate-50 text-slate-400 border border-slate-100 flex items-center justify-center shrink-0">
      <BookOpen class="w-4 h-4" />
    </div>
    <div>
      <h3 class="text-xs font-bold text-slate-700">Manual Book</h3>
      <p class="text-[11px] text-slate-400 font-medium mt-0.5">Unduh panduan penggunaan dashboard untuk Guru Wali</p>
    </div>
  </div>
  <button 
    onclick={downloadManualBookPDF}
    class="flex items-center justify-center gap-1.5 px-4 py-2 bg-slate-800 hover:bg-slate-900 text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none"
  >
    <Download class="w-3.5 h-3.5" />
    Unduh PDF
  </button>
</div>
