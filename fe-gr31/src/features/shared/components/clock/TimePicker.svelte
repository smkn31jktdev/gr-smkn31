<script lang="ts">
  import { Clock, X } from 'lucide-svelte';

  let {
    value = $bindable(''),
    placeholder = 'Pilih waktu...',
    label = '',
    disabled = false,
    onchange
  } = $props<{
    value?: string;
    placeholder?: string;
    label?: string;
    disabled?: boolean;
    onchange?: (time: string) => void;
  }>();

  let isPopoverOpen = $state(false);
  let containerEl = $state<HTMLDivElement | null>(null);

  // Split selected hour and minute
  let selectedHour = $state('12');
  let selectedMinute = $state('00');
  let initialTimeOnOpen = '';

  // Sync state with incoming value
  $effect(() => {
    if (value) {
      const parts = value.split(':');
      if (parts.length >= 2) {
        selectedHour = parts[0];
        selectedMinute = parts[1];
      }
    }
  });

  function togglePopover() {
    if (disabled) return;
    isPopoverOpen = !isPopoverOpen;
    if (isPopoverOpen) {
      initialTimeOnOpen = value || '12:00';
      const parts = initialTimeOnOpen.split(':');
      if (parts.length >= 2) {
        selectedHour = parts[0];
        selectedMinute = parts[1];
      }
    }
  }

  function selectHour(hStr: string) {
    selectedHour = hStr;
  }

  function selectMinute(mStr: string) {
    selectedMinute = mStr;
  }

  function selectQuickTime(timeStr: string) {
    const parts = timeStr.split(':');
    if (parts.length >= 2) {
      selectedHour = parts[0];
      selectedMinute = parts[1];
    }
  }

  function confirmSelection() {
    value = `${selectedHour}:${selectedMinute}`;
    if (onchange) onchange(value);
    isPopoverOpen = false;
  }

  function cancelSelection() {
    const parts = initialTimeOnOpen.split(':');
    if (parts.length >= 2) {
      selectedHour = parts[0];
      selectedMinute = parts[1];
    }
    isPopoverOpen = false;
  }

  function handleClear() {
    value = '';
    selectedHour = '12';
    selectedMinute = '00';
    if (onchange) onchange('');
  }

  // Click outside to close popover
  function handleDocumentClick(event: MouseEvent) {
    if (isPopoverOpen && containerEl && !containerEl.contains(event.target as Node)) {
      cancelSelection();
    }
  }

  $effect(() => {
    document.addEventListener('click', handleDocumentClick, true);
    return () => {
      document.removeEventListener('click', handleDocumentClick, true);
    };
  });

  const uniqueId = 'timepicker-' + Math.random().toString(36).substring(2, 9);
</script>

<div class="relative w-full" bind:this={containerEl}>
  {#if label}
    <label for={uniqueId} class="block text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-1.5">
      {label}
    </label>
  {/if}

  <div class="relative">
    <button
      id={uniqueId}
      type="button"
      {disabled}
      onclick={togglePopover}
      class="w-full flex items-center justify-between text-left border border-slate-100 rounded-xl pl-4 pr-10 py-2 bg-slate-50/50 text-xs font-bold text-slate-650 transition-all duration-205 hover:border-slate-200 focus:outline-none focus:border-[#4db6ac] focus:bg-white disabled:opacity-60 disabled:cursor-not-allowed select-none min-h-[38px]"
      class:border-[#4db6ac]={isPopoverOpen}
      class:bg-white={isPopoverOpen}
      class:ring-3={isPopoverOpen}
      class:ring-teal-50={isPopoverOpen}
    >
      <div class="flex items-center gap-2 overflow-hidden">
        <Clock class="w-3.5 h-3.5 text-slate-400 shrink-0" />
        {#if value}
          <span class="truncate text-slate-800">{value}</span>
        {:else}
          <span class="truncate text-slate-400">{placeholder}</span>
        {/if}
      </div>
    </button>

    {#if value && !disabled}
      <button
        type="button"
        class="absolute right-2.5 top-1/2 -translate-y-1/2 p-0.5 rounded-full text-gray-400 hover:text-gray-600 hover:bg-gray-100 shrink-0 transition-colors z-10 cursor-pointer border-none bg-transparent"
        onclick={(e) => {
          e.stopPropagation();
          handleClear();
        }}
        aria-label="Bersihkan"
      >
        <X class="w-3 h-3" />
      </button>
    {/if}
  </div>

  {#if isPopoverOpen}
    <div class="absolute z-55 left-0 mt-2 p-5 bg-white rounded-3xl border border-gray-150 shadow-[0_12px_40px_rgba(0,0,0,0.08)] select-none w-[275px] flex flex-col items-center">
      <!-- Title "Pilih Waktu" -->
      <span class="text-xs font-bold text-slate-800 tracking-wide mb-1.5">Pilih Waktu</span>
      
      <!-- Large display "14:16" in teal/green -->
      <div class="text-3xl font-black text-[#4db6ac] font-display tracking-tight mb-4 select-all">
        {selectedHour}:{selectedMinute}
      </div>

      <!-- Scroll Selection Columns -->
      <div class="flex justify-center gap-5 w-full mb-5">
        <!-- Hours Column -->
        <div class="flex flex-col items-center">
          <span class="text-[9px] font-black text-gray-400 uppercase tracking-widest mb-1.5">Jam</span>
          <div class="flex flex-col gap-0.5 h-36 overflow-y-auto w-20 border border-slate-100 rounded-2xl p-1.5 custom-time-scroll scroll-smooth text-center bg-slate-50/50">
            {#each Array(24) as _, h}
              {@const hStr = String(h).padStart(2, '0')}
              {@const isSelected = selectedHour === hStr}
              <button
                type="button"
                onclick={() => selectHour(hStr)}
                class="py-1 text-xs font-extrabold rounded-xl transition-all cursor-pointer border-none bg-transparent"
                class:bg-[#4db6ac]={isSelected}
                class:text-white={isSelected}
                class:hover:bg-teal-50={!isSelected}
                class:text-slate-600={!isSelected}
              >
                {hStr}
              </button>
            {/each}
          </div>
        </div>

        <!-- Divider -->
        <span class="font-black text-gray-300 self-center mt-4 text-xl">:</span>

        <!-- Minutes Column -->
        <div class="flex flex-col items-center">
          <span class="text-[9px] font-black text-gray-400 uppercase tracking-widest mb-1.5">Menit</span>
          <div class="flex flex-col gap-0.5 h-36 overflow-y-auto w-20 border border-slate-100 rounded-2xl p-1.5 custom-time-scroll scroll-smooth text-center bg-slate-50/50">
            {#each Array(60) as _, m}
              {@const mStr = String(m).padStart(2, '0')}
              {@const isSelected = selectedMinute === mStr}
              <button
                type="button"
                onclick={() => selectMinute(mStr)}
                class="py-1 text-xs font-extrabold rounded-xl transition-all cursor-pointer border-none bg-transparent"
                class:bg-[#4db6ac]={isSelected}
                class:text-white={isSelected}
                class:hover:bg-teal-50={!isSelected}
                class:text-slate-600={!isSelected}
              >
                {mStr}
              </button>
            {/each}
          </div>
        </div>
      </div>

      <!-- Quick selections "Waktu Cepat" -->
      <div class="w-full flex flex-col items-center mb-5 border-t border-slate-100 pt-3">
        <span class="text-[9px] font-black text-gray-450 uppercase tracking-widest mb-2.5">Waktu Cepat</span>
        <div class="flex flex-wrap gap-1.5 justify-center">
          {#each ['04:00', '04:30', '12:00', '22:00'] as qTime}
            <button
              type="button"
              onclick={() => selectQuickTime(qTime)}
              class="px-3 py-1.5 bg-slate-50 hover:bg-slate-100 active:scale-95 text-slate-700 text-[10px] font-extrabold rounded-xl transition-all border-none cursor-pointer"
            >
              {qTime}
            </button>
          {/each}
        </div>
      </div>

      <!-- Action buttons Batal & Pilih -->
      <div class="flex gap-3 w-full border-t border-slate-100 pt-3.5">
        <button
          type="button"
          onclick={cancelSelection}
          class="flex-1 py-2.5 bg-slate-50 hover:bg-slate-100 text-slate-700 text-xs font-black rounded-xl transition-colors cursor-pointer border-none"
        >
          Batal
        </button>
        <button
          type="button"
          onclick={confirmSelection}
          class="flex-1 py-2.5 bg-[#4db6ac] hover:bg-[#3ca59b] text-white text-xs font-black rounded-xl transition-colors cursor-pointer border-none"
        >
          Pilih
        </button>
      </div>
    </div>
  {/if}
</div>

<style>
  :global(.custom-time-scroll::-webkit-scrollbar) {
    width: 5px;
  }
  :global(.custom-time-scroll::-webkit-scrollbar-track) {
    background: transparent;
  }
  :global(.custom-time-scroll::-webkit-scrollbar-thumb) {
    background: #10b981;
    border-radius: 99px;
  }
  :global(.custom-time-scroll::-webkit-scrollbar-thumb:hover) {
    background: #059669;
  }
</style>
