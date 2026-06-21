<script lang="ts">
  import {
    Sunrise,
    Heart,
    Compass,
    Clock,
    Sparkles,
    Sun,
    Award,
    Moon,
    Calendar,
    Coins,
    Dumbbell,
    Utensils,
    Book,
    BookOpen,
    GraduationCap,
    ClipboardList,
    Users,
    BedDouble
  } from 'lucide-svelte';

  let { 
    distribusiPredikat = {},
    rataRataPerIndikator = {}
  } = $props<{
    distribusiPredikat?: Record<string, number>;
    rataRataPerIndikator?: Record<string, number>;
  }>();

  // Predicate list labels
  const predicates = [
    { key: 'Istimewa', color: 'bg-emerald-500' },
    { key: 'Sangat Baik', color: 'bg-teal-500' },
    { key: 'Baik', color: 'bg-blue-500' },
    { key: 'Cukup', color: 'bg-amber-500' },
    { key: 'Kurang', color: 'bg-rose-500' }
  ];

  const iconMap: Record<string, any> = {
    bangunPagi: Sunrise,
    ibadahDoa: Heart,
    ibadahSholatFajar: Compass,
    ibadahSholat5Waktu: Clock,
    ibadahZikir: Sparkles,
    ibadahDhuha: Sun,
    ibadahRowatib: Award,
    ibadahTarawih: Moon,
    ibadahPuasa: Calendar,
    ibadahZakat: Coins,
    olahraga: Dumbbell,
    makanSehat: Utensils,
    belajarKitabSuci: Book,
    belajarBukuUmum: BookOpen,
    belajarBukuMapel: GraduationCap,
    belajarTugas: ClipboardList,
    bermasyarakat: Users,
    tidurCepat: BedDouble
  };

  function getIndikatorLabel(key: string) {
    const labels: Record<string, string> = {
      bangunPagi: 'Bangun Pagi',
      ibadahDoa: 'Berdoa Harian',
      ibadahSholatFajar: 'Shalat Fajar (Muslim)',
      ibadahSholat5Waktu: 'Shalat 5 Waktu (Muslim)',
      ibadahZikir: 'Zikir Shalat (Muslim)',
      ibadahDhuha: 'Shalat Dhuha (Muslim)',
      ibadahRowatib: 'Shalat Rawatib (Muslim)',
      ibadahTarawih: 'Shalat Tarawih (Muslim)',
      ibadahPuasa: 'Puasa Ramadhan (Muslim)',
      ibadahZakat: 'Zakat/Infaq/Sedekah',
      olahraga: 'Berolahraga',
      makanSehat: 'Makan Sehat',
      belajarKitabSuci: 'Membaca Kitab Suci',
      belajarBukuUmum: 'Membaca Buku Umum',
      belajarBukuMapel: 'Membaca Buku Mapel',
      belajarTugas: 'Mengerjakan PR/Tugas',
      bermasyarakat: 'Bermasyarakat/Sosial',
      tidurCepat: 'Tidur Tepat Waktu'
    };
    return labels[key] || key;
  }

  let totalScored = $derived(
    (Object.values(distribusiPredikat) as number[]).reduce((sum, val) => sum + val, 0)
  );

  let indicatorsList = $derived(
    Object.entries(rataRataPerIndikator).map(([key, val]) => ({
      key,
      label: getIndikatorLabel(key),
      value: val as number
    }))
  );
</script>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
  <!-- Predicate Distribution Card -->
  <div class="card p-6">
    <div class="mb-4">
      <h3 class="text-sm font-bold text-foreground">Distribusi Predikat Nilai</h3>
      <p class="text-xs text-muted mt-0.5">Tingkatan predikat nilai bulanan siswa kelas ini</p>
    </div>

    {#if totalScored === 0}
      <p class="text-xs text-muted text-center py-12">Belum ada data nilai</p>
    {:else}
      <div class="space-y-4">
        {#each predicates as pred}
          {@const count = distribusiPredikat[pred.key] || 0}
          {@const percent = totalScored > 0 ? (count / totalScored) * 100 : 0}
          <div class="space-y-1">
            <div class="flex items-center justify-between text-xs font-semibold">
              <span class="text-foreground">{pred.key}</span>
              <span class="text-muted">{count} Siswa ({percent.toFixed(0)}%)</span>
            </div>
            <div class="w-full h-3.5 bg-gray-100 rounded-full overflow-hidden border border-border">
              <div 
                class="h-full rounded-full transition-all duration-500 {pred.color}" 
                style="width: {percent}%"
              ></div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Indicators Averages Card -->
  <div class="card p-6 flex flex-col h-[400px]">
    <div class="mb-4">
      <h3 class="text-sm font-bold text-foreground">Rata-rata per Sub-Indikator</h3>
      <p class="text-xs text-muted mt-0.5">Skor rata-rata siswa (skala 1–5)</p>
    </div>

    {#if indicatorsList.length === 0}
      <p class="text-xs text-muted text-center py-12 flex-1 flex items-center justify-center">Belum ada data penilaian jurnal</p>
    {:else}
      <div class="flex-1 overflow-y-auto pr-1 space-y-3.5">
        {#each indicatorsList as ind}
          {@const valPercent = (ind.value / 5) * 100}
          <div class="space-y-1">
            <div class="flex items-center justify-between text-xs font-semibold">
              <span class="text-foreground truncate max-w-[240px] flex items-center gap-1.5">
                {#if iconMap[ind.key]}
                  {@const IconComponent = iconMap[ind.key]}
                  <IconComponent class="w-3.5 h-3.5 text-slate-400 shrink-0" />
                {/if}
                {ind.label}
              </span>
              <span class="font-bold text-primary">{ind.value.toFixed(1)} <span class="text-muted font-normal">/ 5.0</span></span>
            </div>
            <div class="w-full h-2.5 bg-gray-100 rounded-full overflow-hidden border border-border">
              <!-- Color code based on average score: red < 3, yellow < 4, green >= 4 -->
              <div 
                class="h-full rounded-full transition-all duration-500" 
                class:bg-rose-500={ind.value < 3}
                class:bg-amber-500={ind.value >= 3 && ind.value < 4}
                class:bg-primary={ind.value >= 4}
                style="width: {valPercent}%"
              ></div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
