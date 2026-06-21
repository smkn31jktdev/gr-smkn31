<script lang="ts">
  import { Clock as ClockIcon, Hourglass, Timer as TimerIcon } from 'lucide-svelte';
  import Clock from './clock/Clock.svelte';
  import Timer from './clock/Timer.svelte';
  import Stopwatch from './clock/Stopwatch.svelte';
  import TimePicker from './clock/TimePicker.svelte';

  let {
    mode = 'clock',
    allowSwitch = true,
    initialDuration = 60,
    value = $bindable(''),
    placeholder = 'Pilih waktu...',
    label = '',
    disabled = false,
    onComplete,
    onLap,
    onchange
  } = $props<{
    mode?: 'clock' | 'timer' | 'stopwatch' | 'timepicker';
    allowSwitch?: boolean;
    initialDuration?: number;
    value?: string;
    placeholder?: string;
    label?: string;
    disabled?: boolean;
    onComplete?: () => void;
    onLap?: (lapTime: string, laps: string[]) => void;
    onchange?: (time: string) => void;
  }>();

  let activeMode = $state(mode);

  $effect(() => {
    activeMode = mode;
  });
</script>

{#if activeMode === 'timepicker'}
  <TimePicker
    bind:value={value}
    {placeholder}
    {label}
    {disabled}
    {onchange}
  />
{:else}
  <div class="card flex flex-col items-center max-w-sm w-full mx-auto select-none bg-radial from-white to-gray-50 border border-gray-100 shadow-xl overflow-hidden transition-all duration-300">
    
    <!-- Tab Switches -->
    {#if allowSwitch}
      <div class="flex bg-gray-100/80 p-1 rounded-xl w-full mb-6 border border-gray-200/50">
        <button
          type="button"
          onclick={() => activeMode = 'clock'}
          class="flex-1 flex items-center justify-center gap-1.5 py-2 text-xs font-bold rounded-lg transition-all cursor-pointer border-none bg-transparent"
          class:bg-white={activeMode === 'clock'}
          class:text-gray-800={activeMode === 'clock'}
          class:shadow-sm={activeMode === 'clock'}
          class:text-gray-500={activeMode !== 'clock'}
          class:hover:text-gray-700={activeMode !== 'clock'}
        >
          <ClockIcon class="w-3.5 h-3.5" />
          Jam
        </button>
        <button
          type="button"
          onclick={() => activeMode = 'timer'}
          class="flex-1 flex items-center justify-center gap-1.5 py-2 text-xs font-bold rounded-lg transition-all cursor-pointer border-none bg-transparent"
          class:bg-white={activeMode === 'timer'}
          class:text-gray-800={activeMode === 'timer'}
          class:shadow-sm={activeMode === 'timer'}
          class:text-gray-500={activeMode !== 'timer'}
          class:hover:text-gray-700={activeMode !== 'timer'}
        >
          <Hourglass class="w-3.5 h-3.5" />
          Timer
        </button>
        <button
          type="button"
          onclick={() => activeMode = 'stopwatch'}
          class="flex-1 flex items-center justify-center gap-1.5 py-2 text-xs font-bold rounded-lg transition-all cursor-pointer border-none bg-transparent"
          class:bg-white={activeMode === 'stopwatch'}
          class:text-gray-800={activeMode === 'stopwatch'}
          class:shadow-sm={activeMode === 'stopwatch'}
          class:text-gray-500={activeMode !== 'stopwatch'}
          class:hover:text-gray-700={activeMode !== 'stopwatch'}
        >
          <TimerIcon class="w-3.5 h-3.5" />
          Stopwatch
        </button>
      </div>
    {/if}

    {#if activeMode === 'clock'}
      <Clock />
    {:else if activeMode === 'timer'}
      <Timer {initialDuration} {onComplete} />
    {:else if activeMode === 'stopwatch'}
      <Stopwatch {onLap} />
    {/if}
  </div>
{/if}
