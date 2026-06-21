<script lang="ts">
  import DatePicker from '../../features/shared/components/DatePicker.svelte';
  import ClockTimer from '../../features/shared/components/ClockTimer.svelte';
  import DropdownChoice, { type DropdownOption } from '../../features/shared/components/DropdownChoice.svelte';
  import { Shield, BookOpen, GraduationCap, ArrowLeft, Calendar, Info, Clock, Check } from 'lucide-svelte';

  // State bindings for DatePicker
  let selectedDate = $state('2026-06-11');
  const highlightedDates = ['2026-06-01', '2026-06-05', '2026-06-11', '2026-06-15', '2026-06-25'];

  // State bindings for TimePicker
  let selectedTime = $state('14:15');

  // State bindings for Dropdown Choice (Single)
  let selectedAssessor = $state('bk');
  const assessorOptions: DropdownOption[] = [
    { 
      value: 'walas', 
      label: 'Budi Santoso, S.Pd', 
      description: 'Wali Kelas X LP', 
      badge: 'Wali Kelas', 
      badgeColor: 'bg-[#e3f2fd] text-[#0d47a1]', 
      icon: GraduationCap 
    },
    { 
      value: 'bk', 
      label: 'Siti Rahma, M.Psi', 
      description: 'Guru Bimbingan Konseling', 
      badge: 'Guru BK', 
      badgeColor: 'bg-emerald-100 text-emerald-800', 
      icon: Shield 
    },
    { 
      value: 'agama', 
      label: 'Ustadz Ahmad Fauzi', 
      description: 'Guru Agama Islam', 
      badge: 'Guru Agama', 
      badgeColor: 'bg-amber-100 text-amber-800', 
      icon: BookOpen 
    }
  ];

  // State bindings for Dropdown Choice (Multiple)
  let selectedStudents = $state<string[]>(['nisn001']);
  const studentOptions: DropdownOption[] = [
    { 
      value: 'nisn001', 
      label: 'Kirana Akbar', 
      description: 'NISN: 0012345678', 
      avatar: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=80&fit=crop&q=80', 
      badge: 'Siswa X LP' 
    },
    { 
      value: 'nisn002', 
      label: 'Kirani', 
      description: 'NISN: 0012345679', 
      avatar: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=80&fit=crop&q=80', 
      badge: 'Siswi X LP' 
    },
    { 
      value: 'nisn003', 
      label: 'Rian Hidayat', 
      description: 'NISN: 0012345680', 
      avatar: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=80&fit=crop&q=80', 
      badge: 'Siswa X LP' 
    }
  ];

  // Event handlers / log records
  let timerLogs = $state<string[]>([]);
  let stopwatchLaps = $state<string[]>([]);

  function handleTimerComplete() {
    timerLogs = [`Timer selesai pada ${new Date().toLocaleTimeString()}`, ...timerLogs];
  }

  function handleStopwatchLap(lapTime: string, laps: string[]) {
    stopwatchLaps = laps;
  }
</script>

<svelte:head>
  <title>Demo Komponen Premium GR31</title>
</svelte:head>

<div class="relative min-h-screen bg-gray-50/50 pb-20 select-none overflow-x-hidden">
  <!-- Dynamic soft lighting glow overlays -->
  <div class="pointer-events-none absolute top-0 left-0 h-[600px] w-[600px] -translate-x-1/3 -translate-y-1/3 rounded-full bg-blue-400/5 blur-[120px]"></div>
  <div class="pointer-events-none absolute right-0 top-1/3 h-[500px] w-[500px] translate-x-1/3 rounded-full bg-teal-400/5 blur-[120px]"></div>
  <div class="pointer-events-none absolute left-1/3 bottom-0 h-[600px] w-[600px] translate-y-1/3 rounded-full bg-amber-400/5 blur-[120px]"></div>

  <!-- Header Banner -->
  <header class="border-b border-gray-100 bg-white/80 backdrop-blur-md sticky top-0 z-40">
    <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <a href="/" class="p-2 hover:bg-gray-100 rounded-xl transition-all text-gray-500 hover:text-gray-700">
          <ArrowLeft class="w-5 h-5" />
        </a>
        <div>
          <h1 class="text-sm font-black text-gray-800 tracking-wide uppercase">Demo Komponen Premium</h1>
          <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-0.5">Svelte 5 Runes & Tailwind CSS v4</p>
        </div>
      </div>
      <span class="inline-flex items-center rounded-full bg-[#e3f2fd] px-3.5 py-1 text-[9px] font-black tracking-widest text-[#0d47a1] uppercase">
        SMK Negeri 31 Jakarta
      </span>
    </div>
  </header>

  <!-- Content Grid -->
  <main class="max-w-7xl mx-auto px-6 mt-8">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      
      <!-- LEFT COLUMN: Calendar & Choice (lg:col-span-8) -->
      <div class="lg:col-span-8 flex flex-col gap-8">
        
        <!-- CARD 1: DatePicker & Calendar -->
        <section class="card bg-white p-6 rounded-2xl border border-gray-200/60 shadow-sm">
          <div class="flex items-start justify-between border-b border-gray-100 pb-4 mb-6">
            <div class="flex items-center gap-2.5">
              <div class="p-2 bg-amber-50 rounded-xl text-amber-600">
                <Calendar class="w-5 h-5" />
              </div>
              <div>
                <h2 class="text-sm font-black text-gray-800 uppercase tracking-wide">1. Kalender / Pemilih Tanggal</h2>
                <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-0.5">DatePicker.svelte</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-8 items-start">
            <div class="flex flex-col gap-4">
              <DatePicker
                label="Tanggal Pembiasaan"
                bind:value={selectedDate}
                placeholder="Pilih tanggal hari ini"
                {highlightedDates}
              />
              
              <div class="bg-gray-50 p-4 rounded-xl border border-gray-100 text-xs font-semibold text-gray-600">
                <div class="flex items-center gap-1.5 text-teal-600 mb-2">
                  <Info class="w-3.5 h-3.5 shrink-0" />
                  <span>Petunjuk Visual</span>
                </div>
                Dot teal <span class="inline-block w-2 h-2 rounded-full bg-teal-500 mx-1"></span> menandakan tanggal yang sudah terisi jurnal pembiasaan/kehadiran.
              </div>
            </div>

            <!-- Value display panel -->
            <div class="flex flex-col gap-3 justify-center bg-gray-50/50 p-6 rounded-2xl border border-gray-100">
              <span class="text-[9px] font-black text-gray-400 uppercase tracking-widest">Reaktif Bounding</span>
              <div class="text-xs font-bold text-gray-700">
                Format DB (YYYY-MM-DD): 
                <span class="ml-1 px-2.5 py-1 bg-white border border-gray-200 rounded-lg text-amber-700 font-mono font-black">{selectedDate || 'Belum dipilih'}</span>
              </div>
              <div class="text-[10px] text-gray-400 leading-relaxed mt-2 font-medium">
                Komponen ini otomatis membaca format tanggal secara presisi, mendukung validasi min/max date, dan terintegrasi dengan click-outside handler.
              </div>
            </div>
          </div>
        </section>

        <!-- CARD 2: Dropdown Choice -->
        <section class="card bg-white p-6 rounded-2xl border border-gray-200/60 shadow-sm">
          <div class="flex items-start justify-between border-b border-gray-100 pb-4 mb-6">
            <div class="flex items-center gap-2.5">
              <div class="p-2 bg-teal-50 rounded-xl text-teal-600">
                <Check class="w-5 h-5" />
              </div>
              <div>
                <h2 class="text-sm font-black text-gray-800 uppercase tracking-wide">2. Dropdown Pilihan Kustom</h2>
                <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-0.5">DropdownChoice.svelte</p>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            <!-- Single Select -->
            <div class="flex flex-col gap-4">
              <DropdownChoice
                label="Pilih Guru Penilai (Single Select)"
                options={assessorOptions}
                bind:value={selectedAssessor}
                placeholder="Pilih guru..."
              />
              
              <div class="bg-gray-50 p-4 rounded-xl border border-gray-100 text-xs font-semibold">
                <span class="text-[9px] font-black text-gray-400 uppercase tracking-widest block mb-1">State Aktif</span>
                <span class="font-mono text-teal-700 font-black">"{selectedAssessor}"</span>
              </div>
            </div>

            <!-- Multi Select with search -->
            <div class="flex flex-col gap-4">
              <DropdownChoice
                label="Pilih Siswa Terlibat (Multi Select + Cari)"
                options={studentOptions}
                bind:value={selectedStudents}
                multiple={true}
                searchable={true}
                placeholder="Cari & pilih siswa..."
              />
              
              <div class="bg-gray-50 p-4 rounded-xl border border-gray-100 text-xs font-semibold">
                <span class="text-[9px] font-black text-gray-400 uppercase tracking-widest block mb-1">Array ID Terpilih</span>
                <span class="font-mono text-teal-700 font-black">{JSON.stringify(selectedStudents)}</span>
              </div>
            </div>
          </div>
        </section>

      </div>

      <!-- RIGHT COLUMN: Clock & Timer (lg:col-span-4) -->
      <div class="lg:col-span-4 flex flex-col gap-8">
        
        <!-- CARD 3: Clock & Timer Showcase -->
        <section class="card bg-white p-6 rounded-2xl border border-gray-200/60 shadow-sm flex flex-col">
          <div class="flex items-start justify-between border-b border-gray-100 pb-4 mb-6 w-full">
            <div class="flex items-center gap-2.5">
              <div class="p-2 bg-blue-50 rounded-xl text-blue-600">
                <Clock class="w-5 h-5" />
              </div>
              <div>
                <h2 class="text-sm font-black text-gray-800 uppercase tracking-wide">3. Jam & Penghitung Waktu</h2>
                <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-0.5">ClockTimer.svelte</p>
              </div>
            </div>
          </div>

          <div class="flex justify-center w-full">
            <ClockTimer
              mode="clock"
              allowSwitch={true}
              initialDuration={15}
              onComplete={handleTimerComplete}
              onLap={handleStopwatchLap}
            />
          </div>

          <!-- Timepicker Showcase -->
          <div class="mt-6 border-t border-gray-100 pt-6 w-full">
            <ClockTimer
              mode="timepicker"
              label="Uji Pemilih Waktu (Time Picker)"
              bind:value={selectedTime}
              placeholder="Pilih waktu..."
            />
            <div class="mt-3 bg-gray-50 p-3 rounded-xl border border-gray-100 text-xs font-semibold text-gray-600 flex justify-between">
              <span>Nilai Binded:</span>
              <span class="font-mono text-amber-700 font-black">"{selectedTime}"</span>
            </div>
          </div>

          <!-- Logs for actions -->
          {#if timerLogs.length > 0 || stopwatchLaps.length > 0}
            <div class="mt-6 border-t border-gray-100 pt-4 w-full">
              <span class="text-[9px] font-black text-gray-400 uppercase tracking-widest block mb-2">Aktivitas Terkini</span>
              
              {#if timerLogs.length > 0}
                <div class="flex flex-col gap-1.5 mb-3">
                  {#each timerLogs.slice(0, 3) as log}
                    <div class="text-[10px] font-semibold text-rose-600 bg-rose-50 px-2.5 py-1.5 rounded-lg border border-rose-100/50">
                      🔔 {log}
                    </div>
                  {/each}
                </div>
              {/if}
              
              {#if stopwatchLaps.length > 0}
                <div class="text-[10px] font-semibold text-gray-500">
                  Lap terakhir dicatat: <span class="font-mono text-gray-700 font-bold">{stopwatchLaps[0]}</span> (Total: {stopwatchLaps.length})
                </div>
              {/if}
            </div>
          {/if}
        </section>

      </div>

    </div>
  </main>
</div>
