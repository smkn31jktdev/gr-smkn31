<script lang="ts">
  import { onMount } from 'svelte';
  import { Play, Pause, RotateCcw, Plus } from 'lucide-svelte';

  let {
    onLap
  } = $props<{
    onLap?: (lapTime: string, laps: string[]) => void;
  }>();

  let stopwatchMs = $state(0);
  let isStopwatchRunning = $state(false);
  let stopwatchLastTime = 0;
  let stopwatchFrameId: number | null = null;
  let laps = $state<string[]>([]);

  function updateStopwatch(time: number) {
    if (!isStopwatchRunning) return;
    const delta = time - stopwatchLastTime;
    stopwatchMs += delta;
    stopwatchLastTime = time;
    stopwatchFrameId = requestAnimationFrame(updateStopwatch);
  }

  function startStopwatch() {
    if (isStopwatchRunning) return;
    isStopwatchRunning = true;
    stopwatchLastTime = performance.now();
    stopwatchFrameId = requestAnimationFrame(updateStopwatch);
  }

  function pauseStopwatch() {
    isStopwatchRunning = false;
    if (stopwatchFrameId !== null) {
      cancelAnimationFrame(stopwatchFrameId);
      stopwatchFrameId = null;
    }
  }

  function resetStopwatch() {
    pauseStopwatch();
    stopwatchMs = 0;
    laps = [];
  }

  function recordLap() {
    const formatted = stopwatchDisplay;
    laps = [formatted, ...laps];
    if (onLap) onLap(formatted, laps);
  }

  const stopwatchDisplay = $derived.by(() => {
    const totalSecs = Math.floor(stopwatchMs / 1000);
    const ms = Math.floor((stopwatchMs % 1000) / 10);
    const s = totalSecs % 60;
    const m = Math.floor(totalSecs / 60);
    
    const mm = String(m).padStart(2, '0');
    const ss = String(s).padStart(2, '0');
    const cc = String(ms).padStart(2, '0');
    return `${mm}:${ss}.${cc}`;
  });

  onMount(() => {
    return () => {
      if (stopwatchFrameId !== null) cancelAnimationFrame(stopwatchFrameId);
    };
  });
</script>

<div class="flex flex-col items-center py-2 w-full">
  <!-- Stopwatch face -->
  <div class="text-4xl font-black font-display text-gray-800 tracking-tight mb-8 py-6 bg-gray-50 border border-gray-100 rounded-2xl w-full text-center shadow-inner">
    {stopwatchDisplay}
  </div>

  <!-- Controls Buttons -->
  <div class="flex items-center gap-3 mb-6">
    <button
      type="button"
      onclick={resetStopwatch}
      class="p-2.5 rounded-xl border border-gray-200 text-gray-500 hover:text-gray-700 hover:bg-gray-50 transition-all cursor-pointer shadow-sm"
      title="Reset"
    >
      <RotateCcw class="w-4 h-4" />
    </button>

    {#if isStopwatchRunning}
      <button
        type="button"
        onclick={pauseStopwatch}
        class="flex items-center gap-1.5 px-5 py-2.5 rounded-xl bg-amber-450 text-white font-bold text-sm hover:bg-amber-500 active:scale-95 transition-all shadow-md cursor-pointer border-none"
      >
        <Pause class="w-4 h-4" />
        Pause
      </button>
      
      <button
        type="button"
        onclick={recordLap}
        class="p-2.5 rounded-xl border border-gray-200 text-gray-600 hover:text-gray-800 hover:bg-gray-50 transition-all cursor-pointer shadow-sm"
        title="Catat Putaran"
      >
        <Plus class="w-4 h-4" />
      </button>
    {:else}
      <button
        type="button"
        onclick={startStopwatch}
        class="flex items-center gap-1.5 px-6 py-2.5 rounded-xl bg-primary text-white font-bold text-sm hover:bg-primary-hover active:scale-95 transition-all shadow-md cursor-pointer border-none"
      >
        <Play class="w-4 h-4" />
        Mulai
      </button>
    {/if}
  </div>

  <!-- Laps Record List -->
  {#if laps.length > 0}
    <div class="w-full border-t border-gray-100 pt-4 max-h-40 overflow-y-auto px-1">
      <span class="text-[9px] font-black text-gray-400 uppercase tracking-widest block mb-2">Catatan Putaran</span>
      <div class="flex flex-col gap-1.5">
        {#each laps as lap, idx}
          <div class="flex items-center justify-between text-xs font-semibold py-1.5 px-2.5 bg-gray-50 rounded-lg border border-gray-100/50">
            <span class="text-gray-400">Putaran {laps.length - idx}</span>
            <span class="text-gray-700 font-mono">{lap}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
