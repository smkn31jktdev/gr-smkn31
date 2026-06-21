<script lang="ts">
  import { 
    CheckCircle, 
    Clock, 
    HelpCircle, 
    AlertCircle, 
    User, 
    ExternalLink 
  } from 'lucide-svelte';
  import type { Kehadiran } from '../../../student/types/student.types';
  import type { RekapBulanan } from '../../../admin/types/admin.types';
  import { getUploadUrl } from '../../../../api/client';

  // Svelte 5 Props destructuring
  let { studentSummary, studentLogs, monthlyTrend = [], reportType = 'bulanan' } = $props<{
    studentSummary: {
      hadir: number;
      izin: number;
      sakit: number;
      alpa: number;
      magang: number;
    };
    studentLogs: Kehadiran[];
    monthlyTrend?: RekapBulanan[];
    reportType?: string;
  }>();

  // Indonesian date formatter helper
  function formatIndonesianDateStr(dateStr: string): string {
    if (!dateStr) return '';
    const parts = dateStr.split('-');
    if (parts.length !== 3) return dateStr;
    const d = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]));
    const days = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
    const months = [
      'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
      'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
    ];
    return `${days[d.getDay()]}, ${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`;
  }

  const monthLabels: Record<string, string> = {
    '01': 'Jan', '02': 'Feb', '03': 'Mar', '04': 'Apr',
    '05': 'Mei', '06': 'Jun', '07': 'Jul', '08': 'Agu',
    '09': 'Sep', '10': 'Okt', '11': 'Nov', '12': 'Des'
  };

  function getStatusBadgeClass(status: string, waktuAbsen?: string) {
    if (status === 'hadir' && waktuAbsen && waktuAbsen.substring(0, 5) > '06:35') {
      return 'bg-rose-50 text-rose-700 border-rose-100/40';
    }
    switch(status) {
      case 'hadir': return 'bg-teal-50 text-[#00a294] border-teal-100/30';
      case 'izin': return 'bg-slate-50 text-slate-650 border-slate-200/50';
      case 'sakit': return 'bg-amber-50/80 text-amber-700 border-amber-100/40';
      case 'magang': return 'bg-violet-50 text-violet-700 border-violet-100/40';
      default: return 'bg-rose-50 text-rose-700 border-rose-100/40';
    }
  }

  function getStatusLabel(status: string, waktuAbsen?: string): string {
    if (status === 'hadir' && waktuAbsen && waktuAbsen.substring(0, 5) > '06:35') return 'terlambat';
    if (status === 'tidak_hadir') return 'alpa';
    return status;
  }

  // Chart helpers
  const chartH = 60;
  const barGap = 6;
  let chartItems = $derived(monthlyTrend.slice(-6));
  let maxVal = $derived(Math.max(...chartItems.map((r: RekapBulanan) => r.persentaseHadir || 0), 50));
  let barW = $derived(chartItems.length > 0 ? Math.max(20, Math.floor((280 - barGap * (chartItems.length - 1)) / chartItems.length)) : 36);
  let svgW = $derived(chartItems.length > 0 ? barW * chartItems.length + barGap * (chartItems.length - 1) : 280);
</script>

<div class="space-y-5">
  <!-- Summary Cards Row -->
  <div class="grid grid-cols-2 sm:grid-cols-5 gap-4">
    <!-- Hadir -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xxs text-left flex items-center gap-3">
      <div class="w-9 h-9 rounded-lg bg-teal-50/50 text-[#00a294] flex items-center justify-center shrink-0 border border-teal-100/20">
        <CheckCircle class="w-4 h-4" />
      </div>
      <div>
        <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Hadir</span>
        <span class="text-sm font-bold text-slate-800 block mt-0.5">{studentSummary.hadir} <span class="text-[10px] font-normal text-slate-400">Hari</span></span>
      </div>
    </div>

    <!-- Izin -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xxs text-left flex items-center gap-3">
      <div class="w-9 h-9 rounded-lg bg-slate-50 text-slate-500 flex items-center justify-center shrink-0 border border-slate-100">
        <Clock class="w-4 h-4" />
      </div>
      <div>
        <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Izin</span>
        <span class="text-sm font-bold text-slate-800 block mt-0.5">{studentSummary.izin} <span class="text-[10px] font-normal text-slate-400">Hari</span></span>
      </div>
    </div>

    <!-- Sakit -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xxs text-left flex items-center gap-3">
      <div class="w-9 h-9 rounded-lg bg-slate-50 text-slate-500 flex items-center justify-center shrink-0 border border-slate-100">
        <HelpCircle class="w-4 h-4" />
      </div>
      <div>
        <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Sakit</span>
        <span class="text-sm font-bold text-slate-800 block mt-0.5">{studentSummary.sakit} <span class="text-[10px] font-normal text-slate-400">Hari</span></span>
      </div>
    </div>

    <!-- Alpa -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xxs text-left flex items-center gap-3">
      <div class="w-9 h-9 rounded-lg bg-rose-50/50 text-rose-600 flex items-center justify-center shrink-0 border border-rose-100/20">
        <AlertCircle class="w-4 h-4" />
      </div>
      <div>
        <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Alpa</span>
        <span class="text-sm font-bold text-rose-600 block mt-0.5">{studentSummary.alpa} <span class="text-[10px] font-normal text-slate-400">Hari</span></span>
      </div>
    </div>

    <!-- Magang -->
    <div class="bg-white border border-slate-100/80 rounded-2xl p-4 shadow-xxs text-left flex items-center gap-3">
      <div class="w-9 h-9 rounded-lg bg-violet-50 text-violet-600 flex items-center justify-center shrink-0 border border-violet-100/20">
        <User class="w-4 h-4" />
      </div>
      <div>
        <span class="text-[9px] font-bold text-slate-400 uppercase tracking-wider block">Magang</span>
        <span class="text-sm font-bold text-slate-800 block mt-0.5">{studentSummary.magang} <span class="text-[10px] font-normal text-slate-400">Hari</span></span>
      </div>
    </div>
  </div>

  <!-- Monthly Trend Chart — only show if we have data -->
  {#if monthlyTrend.length > 0}
    <div class="bg-white border border-slate-100/80 rounded-2xl shadow-xs overflow-hidden">
      <div class="p-4 border-b border-slate-100 flex items-center justify-between text-left">
        <div>
          <h3 class="text-xs font-bold text-slate-800 uppercase tracking-wider">Tren Kehadiran Bulanan</h3>
          <p class="text-[10px] text-slate-400 font-medium mt-0.5">Persentase kehadiran 6 bulan terakhir</p>
        </div>
        <div class="text-right">
          <span class="text-lg font-bold text-[#00a294]">
            {monthlyTrend.length > 0 ? (monthlyTrend[monthlyTrend.length - 1].persentaseHadir ?? 0).toFixed(1) : 0}%
          </span>
          <span class="text-[10px] text-slate-400 font-medium block">Bulan ini</span>
        </div>
      </div>
      <div class="p-5 flex items-end justify-center overflow-x-auto">
        <svg width={svgW} height={chartH + 28} viewBox="0 0 {svgW} {chartH + 28}" class="overflow-visible">
          {#each chartItems as item, idx}
            {@const x = idx * (barW + barGap)}
            {@const pct = item.persentaseHadir ?? 0}
            {@const bh = Math.max(4, Math.round((pct / maxVal) * chartH))}
            {@const by = chartH - bh}
            {@const label = monthLabels[item.bulanTahun.slice(5, 7)] ?? item.bulanTahun.slice(5, 7)}
            {@const isLast = idx === chartItems.length - 1}
            <g>
              <!-- Bar background -->
              <rect
                x={x} y={0} width={barW} height={chartH}
                rx="4" fill="#f8fafc" opacity="0.6"
              />
              <!-- Bar fill -->
              <rect
                x={x} y={by} width={barW} height={bh}
                rx="4"
                fill={isLast ? '#00a294' : '#a7f3d0'}
              />
              <!-- Percentage label -->
              {#if pct > 0}
                <text
                  x={x + barW / 2} y={by - 4}
                  text-anchor="middle" font-size="8" font-weight="600"
                  fill={isLast ? '#00a294' : '#6b7280'}
                >{pct.toFixed(0)}%</text>
              {/if}
              <!-- Month label -->
              <text
                x={x + barW / 2} y={chartH + 18}
                text-anchor="middle" font-size="8" font-weight="600"
                fill={isLast ? '#00a294' : '#9ca3af'}
              >{label}</text>
            </g>
          {/each}
        </svg>
      </div>
    </div>
  {/if}

  <!-- Daily Breakdown List Table -->
  <div class="bg-white border border-slate-100/80 rounded-2xl shadow-xs overflow-hidden">
    <div class="p-4 border-b border-slate-100 text-left">
      <h3 class="text-xs font-bold text-slate-850 uppercase tracking-wider">Breakdown Log Kehadiran Harian</h3>
    </div>
    <div class="overflow-x-auto">
      <table class="w-full text-left border-collapse text-xs">
        <thead>
          <tr class="bg-slate-50/80 border-b border-slate-100 text-slate-400 font-bold text-[10px] uppercase tracking-wider">
            <th class="p-4 pl-6">Hari & Tanggal</th>
            <th class="p-4">Status</th>
            <th class="p-4">Waktu Absen</th>
            <th class="p-4 pr-6">Alasan / Catatan Lampiran</th>
          </tr>
        </thead>
        <tbody>
          {#if studentLogs.length === 0}
            <tr>
              <td colspan="4" class="p-8 text-center text-slate-400 font-medium">
                Tidak ada catatan log kehadiran untuk siswa ini pada periode terpilih.
              </td>
            </tr>
          {:else}
            {#each studentLogs as log}
              <tr class="border-b border-slate-50 hover:bg-slate-50/20 transition-colors">
                <td class="p-4 pl-6 font-semibold text-slate-700">{formatIndonesianDateStr(log.tanggal)}</td>
                <td class="p-4">
                  <span class="inline-block px-2.5 py-0.5 rounded-md border text-[9px] font-bold uppercase tracking-wider {getStatusBadgeClass(log.status, log.waktuAbsen)}">
                    {getStatusLabel(log.status, log.waktuAbsen)}
                  </span>
                </td>
                <td class="p-4 font-semibold font-mono">
                  <span class={log.status === 'hadir' && log.waktuAbsen && log.waktuAbsen.substring(0, 5) > '06:35' ? 'text-rose-600' : 'text-slate-500'}>
                    {log.waktuAbsen ? log.waktuAbsen.substring(0, 5) : '-'}
                  </span>
                </td>
                <td class="p-4 pr-6 text-slate-500 font-medium max-w-sm truncate">
                  {#if log.fotoIzin}
                    <a 
                      href={getUploadUrl(log.fotoIzin)} 
                      target="_blank" 
                      class="inline-flex items-center gap-1 text-[#00a294] hover:text-[#008f83] font-bold hover:underline"
                    >
                      Lihat Lampiran
                      <ExternalLink class="w-3.5 h-3.5" />
                    </a>
                  {:else}
                    {log.alasan || '-'}
                  {/if}
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  </div>
</div>
