<script lang="ts">
  import { onMount } from 'svelte';
  import { Play, Pause, RotateCcw } from 'lucide-svelte';

  let {
    initialDuration = 60,
    onComplete
  } = $props<{
    initialDuration?: number;
    onComplete?: () => void;
  }>();

  let timerMax = $state(initialDuration);
  let timerRemaining = $state(initialDuration);
  let isTimerRunning = $state(false);
  let timerInterval: any = null;

  // Custom Input States
  let inputMinutes = $state(Math.floor(initialDuration / 60));
  let inputSeconds = $state(initialDuration % 60);

  function updateTimerMaxFromInputs() {
    const totalSecs = (inputMinutes * 60) + inputSeconds;
    timerMax = totalSecs > 0 ? totalSecs : 60;
    if (!isTimerRunning) {
      timerRemaining = timerMax;
    }
  }

  function startTimer() {
    if (isTimerRunning) return;
    if (timerRemaining <= 0) {
      timerRemaining = timerMax;
    }
    isTimerRunning = true;
    
    timerInterval = setInterval(() => {
      if (timerRemaining > 0) {
        timerRemaining -= 1;
      } else {
        stopTimer();
        if (onComplete) onComplete();
        playNotificationTone();
      }
    }, 1000);
  }

  function pauseTimer() {
    isTimerRunning = false;
    if (timerInterval) {
      clearInterval(timerInterval);
      timerInterval = null;
    }
  }

  function stopTimer() {
    pauseTimer();
  }

  function resetTimer() {
    pauseTimer();
    timerRemaining = timerMax;
  }

  // Derived values for timer progress ring
  const timerDashArray = 2 * Math.PI * 45; // r=45
  const timerDashOffset = $derived(
    timerMax > 0 ? timerDashArray * (1 - timerRemaining / timerMax) : 0
  );

  const timerDisplay = $derived.by(() => {
    const m = String(Math.floor(timerRemaining / 60)).padStart(2, '0');
    const s = String(timerRemaining % 60).padStart(2, '0');
    return `${m}:${s}`;
  });

  function playNotificationTone() {
    try {
      const audioCtx = new (window.AudioContext || (window as any).webkitAudioContext)();
      const oscillator = audioCtx.createOscillator();
      const gainNode = audioCtx.createGain();
      oscillator.connect(gainNode);
      gainNode.connect(audioCtx.destination);
      oscillator.type = 'sine';
      oscillator.frequency.value = 523.25; // C5 note
      gainNode.gain.setValueAtTime(0, audioCtx.currentTime);
      gainNode.gain.linearRampToValueAtTime(0.5, audioCtx.currentTime + 0.1);
      gainNode.gain.exponentialRampToValueAtTime(0.01, audioCtx.currentTime + 0.8);
      oscillator.start(audioCtx.currentTime);
      oscillator.stop(audioCtx.currentTime + 0.8);
    } catch (e) {
      console.warn('AudioContext not supported', e);
    }
  }

  onMount(() => {
    return () => {
      if (timerInterval) clearInterval(timerInterval);
    };
  });
</script>

<div class="flex flex-col items-center py-2 w-full">
  <!-- Circular Progress Countdown SVG -->
  <div class="relative w-44 h-44 mb-6 flex items-center justify-center">
    <svg class="w-full h-full transform -rotate-90" viewBox="0 0 100 100">
      <circle cx="50" cy="50" r="45" fill="none" stroke="#f3f4f6" stroke-width="5" />
      <circle
        cx="50" cy="50" r="45"
        fill="none"
        stroke="var(--color-primary)"
        stroke-width="5"
        stroke-linecap="round"
        stroke-dasharray={timerDashArray}
        stroke-dashoffset={timerDashOffset}
        class="transition-all duration-1000 ease-linear"
        class:stroke-red-500={timerRemaining < 10 && timerRemaining > 0}
      />
    </svg>
    
    <!-- Timer digital counter inside circle -->
    <div class="absolute flex flex-col items-center">
      <span class="text-3xl font-black text-gray-800 font-display tracking-tight leading-none" class:animate-pulse={isTimerRunning && timerRemaining < 10}>
        {timerDisplay}
      </span>
      {#if timerRemaining === 0}
        <span class="text-[9px] font-black text-red-500 uppercase tracking-widest mt-1 animate-bounce">Selesai!</span>
      {:else}
        <span class="text-[9px] font-bold text-gray-400 uppercase tracking-widest mt-1">Sisa Waktu</span>
      {/if}
    </div>
  </div>

  <!-- Settings Inputs (when not running) -->
  {#if !isTimerRunning}
    <div class="flex items-center gap-2 mb-6">
      <div class="flex flex-col items-center">
        <span class="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Menit</span>
        <input
          type="number"
          min="0"
          max="59"
          bind:value={inputMinutes}
          oninput={updateTimerMaxFromInputs}
          class="w-14 text-center border border-gray-200 rounded-lg p-1 text-sm font-bold bg-white text-gray-800 focus:outline-none focus:border-amber-400"
        />
      </div>
      <span class="font-bold text-gray-400 self-end mb-1.5">:</span>
      <div class="flex flex-col items-center">
        <span class="text-[9px] font-bold text-gray-400 uppercase tracking-wider mb-1">Detik</span>
        <input
          type="number"
          min="0"
          max="59"
          bind:value={inputSeconds}
          oninput={updateTimerMaxFromInputs}
          class="w-14 text-center border border-gray-200 rounded-lg p-1 text-sm font-bold bg-white text-gray-800 focus:outline-none focus:border-amber-400"
        />
      </div>
    </div>
  {:else}
    <!-- Running display text -->
    <div class="h-14 flex items-center justify-center mb-6">
      <span class="text-xs font-semibold text-gray-500 italic">Timer sedang berjalan...</span>
    </div>
  {/if}

  <!-- Controls Buttons -->
  <div class="flex items-center gap-3">
    <button
      type="button"
      onclick={resetTimer}
      class="p-2.5 rounded-xl border border-gray-200 text-gray-500 hover:text-gray-700 hover:bg-gray-50 transition-all cursor-pointer shadow-sm"
      title="Reset"
    >
      <RotateCcw class="w-4 h-4" />
    </button>

    {#if isTimerRunning}
      <button
        type="button"
        onclick={pauseTimer}
        class="flex items-center gap-1.5 px-6 py-2.5 rounded-xl bg-amber-450 text-white font-bold text-sm hover:bg-amber-500 active:scale-95 transition-all shadow-md cursor-pointer border-none"
      >
        <Pause class="w-4 h-4" />
        Pause
      </button>
    {:else}
      <button
        type="button"
        onclick={startTimer}
        disabled={timerRemaining <= 0}
        class="flex items-center gap-1.5 px-6 py-2.5 rounded-xl bg-primary text-white font-bold text-sm hover:bg-primary-hover active:scale-95 transition-all shadow-md cursor-pointer border-none disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <Play class="w-4 h-4" />
        Mulai
      </button>
    {/if}
  </div>
</div>
