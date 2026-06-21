<script lang="ts">
  import { Download } from 'lucide-svelte';
  import { addToast } from '../../../../../../../stores/uiStore';

  let { detailType, detailEvaluate, detailRekap }: {
    detailType: 'bulanan' | 'semester',
    detailEvaluate: any,
    detailRekap: any
  } = $props();

  const HABITS = [
    { key: 'bangunPagi', label: '1. Bangun Pagi', desc: 'Siswa biasa bangun pagi sebelum jam 04.30' },
    { key: 'beribadah', label: '2. Beribadah', desc: 'Siswa biasa beribadah' },
    { key: 'olahraga', label: '3. Berolah raga', desc: 'Siswa melakukan kebiasaan berolahraga' },
    { key: 'makanSehat', label: '4. Makan sehat bergizi', desc: 'Siswa memiliki kebiasaan makan sehat bergizi' },
    { key: 'gemarBelajar', label: '5. Gemar belajar', desc: 'Siswa memiliki kebiasaan belajar' },
    { key: 'bermasyarakat', label: '6. Baik di masyarakat', desc: 'Siswa memiliki perilaku di masyarakat' },
    { key: 'tidurCepat', label: '7. Tidur lebih cepat', desc: 'Siswa biasa tidur / istirahat malam jam 22.00' }
  ];

  function getSemesterHabitScore(habitKey: string, rekap: any): number {
    if (!rekap || !rekap.skor) return 0;
    if (habitKey === 'bangunPagi') return rekap.skor.bangunPagi ?? 0;
    if (habitKey === 'beribadah') {
      const keys = ['ibadahDoa', 'ibadahSholatFajar', 'ibadahSholat5Waktu', 'ibadahZikir', 'ibadahDhuha', 'ibadahRowatib', 'ibadahZakat'];
      let sum = 0, count = 0;
      keys.forEach(k => {
        if (rekap.skor[k] !== undefined && rekap.skor[k] > 0) {
          sum += rekap.skor[k];
          count++;
        }
      });
      return count > 0 ? Number((sum / count).toFixed(1)) : 0;
    }
    if (habitKey === 'olahraga') return rekap.skor.olahraga ?? 0;
    if (habitKey === 'makanSehat') return rekap.skor.makanSehat ?? 0;
    if (habitKey === 'gemarBelajar') {
      const keys = ['belajarKitabSuci', 'belajarBukuUmum', 'belajarBukuMapel', 'belajarTugas'];
      let sum = 0, count = 0;
      keys.forEach(k => {
        if (rekap.skor[k] !== undefined && rekap.skor[k] > 0) {
          sum += rekap.skor[k];
          count++;
        }
      });
      return count > 0 ? Number((sum / count).toFixed(1)) : 0;
    }
    if (habitKey === 'bermasyarakat') return rekap.skor.bermasyarakat ?? 0;
    if (habitKey === 'tidurCepat') return rekap.skor.tidurCepat ?? 0;
    return 0;
  }

  function getScoreLabel(score: number | undefined) {
    if (score === undefined || score === 0) return 'Dilewati';
    if (score >= 4.5) return 'Istimewa';
    if (score >= 3.5) return 'Sangat Baik';
    if (score >= 2.5) return 'Baik';
    if (score >= 1.5) return 'Cukup';
    return 'Kurang';
  }

  let scores = $derived.by(() => {
    if (detailType === 'bulanan') {
      return HABITS.map(h => detailEvaluate?.[h.key]?.skor ?? 0);
    } else {
      return HABITS.map(h => getSemesterHabitScore(h.key, detailRekap));
    }
  });

  let stats = $derived.by(() => {
    if (scores.length === 0) return { avg: 0, max: 0, min: 0, dom: 0, domLabel: '—', distPct: [] };
    const validScores = scores.filter(s => s > 0);
    if (validScores.length === 0) return { avg: 0, max: 0, min: 0, dom: 0, domLabel: '—', distPct: [
      { count: 0, pct: 0 },
      { count: 0, pct: 0 },
      { count: 0, pct: 0 },
      { count: 0, pct: 0 },
      { count: 0, pct: 0 }
    ] };

    const sum = validScores.reduce((a, b) => a + b, 0);
    const avg = Number((sum / validScores.length).toFixed(1));
    const max = Math.max(...validScores);
    const min = Math.min(...validScores);

    const counts: Record<number, number> = {};
    validScores.forEach(s => { 
      const sInt = Math.round(s);
      counts[sInt] = (counts[sInt] || 0) + 1; 
    });
    let dom = Math.round(validScores[0]);
    let maxCount = 0;
    Object.entries(counts).forEach(([s, count]) => {
      if (count > maxCount) {
        maxCount = count;
        dom = Number(s);
      }
    });
    const domLabel = getScoreLabel(dom);

    const dist = [0, 0, 0, 0, 0];
    validScores.forEach(s => {
      const idx = Math.max(1, Math.min(5, Math.round(s))) - 1;
      dist[idx]++;
    });

    const distPct = dist.map(count => ({
      count,
      pct: Math.round((count / validScores.length) * 100)
    }));

    return { avg, max, min, dom, domLabel, distPct };
  });

  let chartPoints = $derived.by(() => {
    const width = 600;
    const height = 180;
    const paddingLeft = 70;
    const paddingRight = 30;
    const paddingTop = 20;
    const paddingBottom = 40;

    const chartW = width - paddingLeft - paddingRight;
    const chartH = height - paddingTop - paddingBottom;

    return scores.map((s, i) => {
      const x = paddingLeft + (i * chartW) / 6;
      const sClamped = Math.max(1, Math.min(5, s || 1));
      const y = paddingTop + chartH - ((sClamped - 1) * chartH) / 4;
      return { x, y, score: s };
    });
  });

  let linePath = $derived.by(() => {
    if (chartPoints.length === 0) return '';
    return chartPoints.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x} ${p.y}`).join(' ');
  });

  let smoothPath = $derived.by(() => {
    if (chartPoints.length === 0) return '';
    let path = `M ${chartPoints[0].x} ${chartPoints[0].y}`;
    for (let i = 0; i < chartPoints.length - 1; i++) {
      const p0 = chartPoints[i];
      const p1 = chartPoints[i + 1];
      const cpX1 = p0.x + (p1.x - p0.x) / 2;
      const cpY1 = p0.y;
      const cpX2 = p0.x + (p1.x - p0.x) / 2;
      const cpY2 = p1.y;
      path += ` C ${cpX1} ${cpY1}, ${cpX2} ${cpY2}, ${p1.x} ${p1.y}`;
    }
    return path;
  });

  let areaPath = $derived.by(() => {
    if (chartPoints.length === 0) return '';
    const first = chartPoints[0];
    const last = chartPoints[chartPoints.length - 1];
    const bottomY = 180 - 40;
    return `${smoothPath} L ${last.x} ${bottomY} L ${first.x} ${bottomY} Z`;
  });
</script>

<div>
  <!-- Header inside tab -->
  <div class="flex items-center justify-between mb-4">
    <div class="text-left">
      <h4 class="text-sm font-bold text-slate-700">Grafik Rating Kebiasaan</h4>
      <p class="text-[10px] text-slate-400 font-bold">Perbandingan rating per indikator</p>
    </div>
    <button
      onclick={() => {
        addToast('Grafik berhasil diexport', 'success');
      }}
      class="flex items-center gap-1.5 px-3 py-1.5 bg-slate-50 hover:bg-slate-100 border border-slate-200/50 text-slate-600 font-bold text-[10px] rounded-lg transition-all cursor-pointer select-none"
    >
      <Download class="w-3 h-3" />
      Unduh Grafik
    </button>
  </div>

  <!-- SVG Line Chart with Grid lines and marker dots -->
  <div class="bg-white border border-slate-100 rounded-3xl p-4 flex justify-center">
    <svg width="600" height="180" viewBox="0 0 600 180" class="overflow-visible">
      <defs>
        <linearGradient id="chart-grad" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="#00a294" stop-opacity="0.15" />
          <stop offset="100%" stop-color="#00a294" stop-opacity="0.00" />
        </linearGradient>
      </defs>

      <!-- Horizontal Grid lines -->
      {#each [1, 2, 3, 4, 5] as lvl}
        {@const y = 20 + 120 - ((lvl - 1) * 120) / 4}
        <line x1="70" y1={y} x2="570" y2={y} stroke="#f1f5f9" stroke-width="1" />
        <text x="60" y={y + 3} text-anchor="end" class="text-[9.5px] font-bold fill-slate-400">
          {lvl} - {lvl === 5 ? 'Istimewa' : lvl === 4 ? 'Sangat Baik' : lvl === 3 ? 'Baik' : lvl === 2 ? 'Cukup' : 'Kurang'}
        </text>
      {/each}

      <!-- Gradient Area & Smooth Bezier Line -->
      {#if linePath}
        <path d={areaPath} fill="url(#chart-grad)" />
        <path d={smoothPath} fill="none" stroke="#00a294" stroke-width="3" stroke-linecap="round" />
      {/if}

      <!-- Data Markers and Rotated labels -->
      {#each chartPoints as p, i}
        <!-- Marker Dot -->
        {@const markerColor = p.score >= 4.5 ? '#10b981' : p.score >= 3.5 ? '#06b6d4' : p.score >= 2.5 ? '#f59e0b' : p.score >= 1.5 ? '#ea580c' : '#ef4444'}
        <circle cx={p.x} cy={p.y} r="5" fill="#fff" stroke={markerColor} stroke-width="2.5" />
        <circle cx={p.x} cy={p.y} r="2" fill={markerColor} />

        <!-- Bottom label rotated slightly -->
        <text x={p.x} y={152} text-anchor="end" transform="rotate(-15 {p.x} 152)" class="text-[8.5px] font-bold fill-slate-500">
          {HABITS[i].label.split('. ')[1]}
        </text>
      {/each}
    </svg>
  </div>

  <!-- Statistics Cards -->
  <div class="grid grid-cols-4 gap-4 mt-6">
    <div class="bg-slate-50/50 border border-slate-100 p-4.5 rounded-2xl text-center">
      <span class="text-[9px] font-black text-slate-400 uppercase tracking-wider font-display">Rata-rata</span>
      <h4 class="text-xl font-black text-slate-700 mt-1 font-mono">{stats.avg}</h4>
    </div>
    <div class="bg-slate-50/50 border border-slate-100 p-4.5 rounded-2xl text-center">
      <span class="text-[9px] font-black text-slate-400 uppercase tracking-wider font-display">Tertinggi</span>
      <h4 class="text-xl font-black text-emerald-500 mt-1 font-mono">{stats.max}</h4>
    </div>
    <div class="bg-slate-50/50 border border-slate-100 p-4.5 rounded-2xl text-center">
      <span class="text-[9px] font-black text-slate-400 uppercase tracking-wider font-display">Terendah</span>
      <h4 class="text-xl font-black text-rose-500 mt-1 font-mono">{stats.min}</h4>
    </div>
    <div class="bg-slate-50/50 border border-slate-100 p-4.5 rounded-2xl text-center">
      <span class="text-[9px] font-black text-slate-400 uppercase tracking-wider font-display">Dominan</span>
      <h4 class="text-xl font-black text-slate-700 mt-1 font-mono">{stats.dom}</h4>
      <p class="text-[8px] font-bold text-slate-400 mt-0.5 uppercase">{stats.domLabel}</p>
    </div>
  </div>

  <!-- Distribusi Rating Section -->
  <div class="mt-6 text-left border-t border-slate-100 pt-5">
    <h4 class="text-xs font-black text-slate-400 uppercase tracking-wider mb-4">Distribusi Rating</h4>
    <div class="space-y-3">
      {#each [{ idx: 0, label: '1. Kurang Baik', colorClass: 'bg-rose-500' }, 
              { idx: 1, label: '2. Cukup Baik', colorClass: 'bg-orange-500' }, 
              { idx: 2, label: '3. Baik', colorClass: 'bg-blue-500' }, 
              { idx: 3, label: '4. Sangat Baik', colorClass: 'bg-indigo-500' }, 
              { idx: 4, label: '5. Istimewa', colorClass: 'bg-emerald-500' }] as rating}
        {@const item = stats.distPct[rating.idx]}
        <div class="flex items-center gap-4 text-xs font-semibold">
          <span class="w-24 text-slate-500 select-none text-[11px] font-bold">{rating.label}</span>
          <div class="flex-1 h-2 bg-slate-100 rounded-full overflow-hidden">
            <div class="h-full rounded-full transition-all duration-500 {rating.colorClass}" style="width: {item?.pct ?? 0}%"></div>
          </div>
          <span class="w-16 text-right text-slate-600 font-mono text-[11px]">{item?.count ?? 0} ({item?.pct ?? 0}%)</span>
        </div>
      {/each}
    </div>
  </div>
</div>
