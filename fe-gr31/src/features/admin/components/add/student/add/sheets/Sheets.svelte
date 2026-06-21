<script lang="ts">
  import { HelpCircle, RefreshCw, Loader2, Lock } from 'lucide-svelte';
  import { bulkImportStudents } from '../../../../../logic/adminLogic';
  import { addToast } from '../../../../../../../stores/uiStore';

  let { onSuccess }: { onSuccess: () => Promise<void> } = $props();

  let sheetUrl = $state('');
  let sheetLoading = $state(false);
  let sheetStudentsList = $state<any[]>([]);
  let sheetSubmitting = $state(false);

  // Client-side Google Sheets parser for students
  async function loadSheetData() {
    const url = sheetUrl.trim();
    if (!url) {
      addToast('URL Google Sheet wajib diisi', 'warning');
      return;
    }

    sheetLoading = true;
    try {
      let csvUrl = url;
      const match = url.match(/\/d\/([a-zA-Z0-9-_]+)/);
      if (match) {
        const docId = match[1];
        csvUrl = `https://docs.google.com/spreadsheets/d/${docId}/export?format=csv`;
      } else {
        addToast('Format URL Google Sheet tidak valid', 'error');
        sheetLoading = false;
        return;
      }

      const res = await fetch(csvUrl);
      if (!res.ok) {
        throw new Error('Gagal mendownload data sheet. Pastikan sheet diatur publik (akses siapa saja dengan link).');
      }
      const csvText = await res.text();
      
      const lines = csvText.split(/\r?\n/);
      if (lines.length < 2) {
        addToast('Data sheet kosong atau tidak valid', 'warning');
        sheetLoading = false;
        return;
      }

      const headers = lines[0].split(',').map(h => h.trim().toLowerCase());
      const parsedItems = [];

      for (let i = 1; i < lines.length; i++) {
        const line = lines[i].trim();
        if (!line) continue;
        
        const cols = line.split(',').map(c => c.trim().replace(/^"|"$/g, ''));
        
        let nisIdx = headers.indexOf('nis');
        if (nisIdx === -1) nisIdx = headers.indexOf('nisn');
        if (nisIdx === -1) nisIdx = 0;
        
        const namaIdx = headers.indexOf('nama') !== -1 ? headers.indexOf('nama') : 1;
        const kelasIdx = headers.indexOf('kelas') !== -1 ? headers.indexOf('kelas') : 2;
        
        // Find index of wali kelas / walas / guru wali / guru_wali
        let walasIdx = headers.indexOf('walas');
        if (walasIdx === -1) walasIdx = headers.indexOf('guru wali');
        if (walasIdx === -1) walasIdx = headers.indexOf('guru_wali');
        if (walasIdx === -1) walasIdx = headers.indexOf('wali kelas');
        if (walasIdx === -1) walasIdx = 3;
        
        const passIdx = headers.indexOf('password') !== -1 ? headers.indexOf('password') : 4;

        const nisVal = cols[nisIdx] || '';
        const item = {
          nis: nisVal,
          nama: cols[namaIdx] || '',
          kelas: cols[kelasIdx] || 'X LP',
          walas: cols[walasIdx] || '',
          agama: 'islam',
          email: nisVal ? nisVal + '@student.smk31.sch.id' : '',
          password: cols[passIdx] || 'changeme123'
        };

        if (item.nis && item.nama) {
          parsedItems.push(item);
        }
      }

      sheetStudentsList = parsedItems;
      addToast(`Berhasil memuat ${parsedItems.length} data siswa dari sheet`, 'success');
    } catch (err: any) {
      console.error(err);
      addToast(err.message || 'Gagal memproses Google Sheet', 'error');
    } finally {
      sheetLoading = false;
    }
  }

  async function submitSheetStudents() {
    if (sheetStudentsList.length === 0) {
      addToast('Tidak ada data siswa untuk disimpan', 'warning');
      return;
    }

    sheetSubmitting = true;
    try {
      const success = await bulkImportStudents(sheetStudentsList);
      if (success) {
        sheetStudentsList = [];
        sheetUrl = '';
        await onSuccess();
      }
    } catch (err) {
      console.error(err);
      addToast('Gagal mengimpor data bulk siswa', 'error');
    } finally {
      sheetSubmitting = false;
    }
  }
</script>

<div class="grid grid-cols-1 gap-6 animate-fade-in">
  <!-- Sheets Control Card -->
  <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs text-left space-y-5">
    <!-- Panduan Import -->
    <div class="p-4 rounded-xl bg-slate-50 border border-slate-100/60 text-slate-500 text-xs leading-relaxed space-y-2">
      <div class="flex items-center gap-1.5 font-bold text-slate-700">
        <HelpCircle class="w-4 h-4 text-[#00a294]" />
        <span>Panduan Import Google Sheets Siswa</span>
      </div>
      <ol class="list-decimal pl-4 space-y-1.5 font-semibold text-[11px] text-slate-500">
        <li>Atur Google Sheet Anda agar dapat diakses publik (Akses Link: <span class="text-[#00a294]">Anyone with the link / Siapa saja dengan link</span>).</li>
        <li>Pastikan kolom baris pertama diisi dengan header persis berikut: <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">NIS</span> (atau <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">NISN</span>), <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Nama</span>, <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Kelas</span>, <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Walas</span> (atau <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Guru Wali</span>), <span class="font-mono bg-white px-1.5 py-0.5 border border-slate-150 rounded text-slate-600">Password</span>.</li>
      </ol>
    </div>

    <!-- URL Input & Fetch Controls -->
    <div class="flex flex-col gap-2">
      <label for="sheet-url" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Tautan / URL Spreadsheet Google</label>
      <div class="flex flex-col sm:flex-row gap-3">
        <input 
          id="sheet-url"
          type="text" 
          placeholder="Contoh: https://docs.google.com/spreadsheets/d/.../edit?usp=sharing"
          bind:value={sheetUrl}
          class="flex-1 bg-slate-50 border border-slate-100 focus:border-slate-350 focus:bg-white text-slate-755 text-xs font-semibold py-2.5 px-3.5 rounded-xl outline-none transition-all"
        />
        
        <div class="flex gap-2 shrink-0">
          <!-- Load Button -->
          <button 
            onclick={loadSheetData}
            disabled={sheetLoading}
            class="flex items-center justify-center gap-1.5 px-4 py-2.5 bg-slate-50 hover:bg-slate-100 disabled:bg-slate-100 text-slate-600 border border-slate-200/50 rounded-xl font-bold text-xs transition-all cursor-pointer"
          >
            {#if sheetLoading}
              <Loader2 class="w-3.5 h-3.5 animate-spin" />
              Mengunduh...
            {:else}
              <RefreshCw class="w-3.5 h-3.5" />
              Muat Data Sheet
            {/if}
          </button>

          <!-- Save All Button -->
          {#if sheetStudentsList.length > 0}
            <button 
              onclick={submitSheetStudents}
              disabled={sheetSubmitting}
              class="flex items-center justify-center gap-1.5 px-5 py-2.5 bg-slate-800 hover:bg-slate-900 disabled:bg-slate-250 text-white rounded-xl font-bold text-xs shadow-xs transition-all active:scale-98 cursor-pointer border-none"
            >
              {#if sheetSubmitting}
                <Loader2 class="w-3.5 h-3.5 animate-spin" />
                Menyimpan...
              {:else}
                <Lock class="w-3.5 h-3.5" />
                Simpan Semua ({sheetStudentsList.length})
              {/if}
            </button>
          {/if}
        </div>
      </div>
    </div>
  </div>

  <!-- Preview Table Card -->
  {#if sheetStudentsList.length > 0}
    <div class="bg-white rounded-2xl border border-slate-100 p-6 shadow-xs min-h-[250px] flex flex-col text-left">
      <!-- Header -->
      <div class="border-b border-slate-50 pb-3 mb-4 flex justify-between items-center">
        <div>
          <h3 class="text-xs font-bold text-slate-700 uppercase tracking-wider">Preview Baris Google Sheet</h3>
          <p class="text-[10px] text-slate-400 font-semibold mt-0.5">Tinjau list siswa sebelum disimpan ke database</p>
        </div>
        <span class="px-2.5 py-0.5 rounded-lg bg-slate-50 border border-slate-100 text-[10px] font-extrabold text-slate-500 font-mono">
          {sheetStudentsList.length} Baris
        </span>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100">
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">No</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">NIS</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Nama Lengkap</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Kelas</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Guru Wali</th>
              <th class="py-2.5 px-4 text-[9px] font-bold text-slate-400 uppercase tracking-wider">Password</th>
            </tr>
          </thead>
          <tbody>
            {#each sheetStudentsList as stud, idx}
              <tr class="border-b border-slate-50 hover:bg-slate-50/20 transition-colors">
                <td class="py-2.5 px-4 text-xs font-bold text-slate-400 font-mono">{idx + 1}</td>
                <td class="py-2.5 px-4 text-xs font-bold text-slate-700 font-mono">{stud.nis}</td>
                <td class="py-2.5 px-4 text-xs font-bold text-slate-700 uppercase">{stud.nama}</td>
                <td class="py-2.5 px-4 text-xs font-bold text-slate-500">{stud.kelas}</td>
                <td class="py-2.5 px-4 text-xs font-medium text-slate-655">{stud.walas || '-'}</td>
                <td class="py-2.5 px-4 text-xs font-medium text-slate-400 font-mono">{stud.password}</td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  {/if}
</div>
