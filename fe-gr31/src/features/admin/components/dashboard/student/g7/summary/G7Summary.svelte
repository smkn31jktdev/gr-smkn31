<script lang="ts">
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
</script>

<div class="space-y-4">
  {#each HABITS as habit}
    {@const score = detailType === 'bulanan' ? (detailEvaluate?.[habit.key]?.skor ?? 0) : getSemesterHabitScore(habit.key, detailRekap)}
    {@const note = detailType === 'bulanan' ? (detailEvaluate?.[habit.key]?.note ?? 'Belum ada data') : `Rata-rata rating semester untuk kebiasaan ini.`}
    {@const label = getScoreLabel(score)}
    <div class="bg-slate-50/50 border border-slate-100 rounded-3xl p-5 flex items-start gap-4">
      <div class="w-10 h-10 rounded-full bg-white border border-slate-100 flex items-center justify-center font-bold text-slate-400 text-xs shrink-0 select-none">
        {habit.label.split('.')[0]}
      </div>
      <div class="flex-1 space-y-2 text-left">
        <h4 class="text-sm font-bold text-slate-700">{habit.label} : {habit.desc.split(': ')[1] || habit.desc}</h4>
        <div class="flex items-center gap-2">
          <span class="px-2.5 py-0.5 rounded-full text-[10px] font-bold bg-[#f3e8ff] text-indigo-600 border border-indigo-100">
            Nilai: {score} ({label})
          </span>
        </div>
        <div class="bg-white border border-slate-100/80 rounded-2xl p-3.5 text-xs text-slate-400 font-medium italic">
          "{note}"
        </div>
      </div>
    </div>
  {/each}
</div>
