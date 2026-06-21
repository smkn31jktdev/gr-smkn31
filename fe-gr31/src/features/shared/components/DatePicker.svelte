<script lang="ts">
  import { Calendar, ChevronLeft, ChevronRight, X } from 'lucide-svelte';
  
  let {
    value = $bindable(''),
    placeholder = 'Pilih tanggal',
    label = '',
    disabled = false,
    minDate = '',
    maxDate = '',
    highlightedDates = [],
    onchange
  } = $props<{
    value?: string;
    placeholder?: string;
    label?: string;
    disabled?: boolean;
    minDate?: string;
    maxDate?: string;
    highlightedDates?: string[];
    onchange?: (date: string) => void;
  }>();

  // Helper arrays for Indonesian localization
  const monthNames = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];
  const dayNames = ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'];

  // State using Svelte 5 Runes
  let isOpen = $state(false);
  let containerEl = $state<HTMLDivElement | null>(null);

  // Initialize view month/year based on value or today
  let today = new Date();
  let initialDate = value ? new Date(value) : today;
  let viewYear = $state(isNaN(initialDate.getTime()) ? today.getFullYear() : initialDate.getFullYear());
  let viewMonth = $state(isNaN(initialDate.getTime()) ? today.getMonth() : initialDate.getMonth());

  // Watch for external value changes to sync view
  $effect(() => {
    if (value) {
      const d = new Date(value);
      if (!isNaN(d.getTime())) {
        viewYear = d.getFullYear();
        viewMonth = d.getMonth();
      }
    }
  });

  // Derived: Formatted value to display in input
  const displayValue = $derived.by(() => {
    if (!value) return '';
    const d = new Date(value);
    if (isNaN(d.getTime())) return '';
    
    const day = d.getDate();
    const month = monthNames[d.getMonth()];
    const year = d.getFullYear();
    const daysOfWeek = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
    const dayName = daysOfWeek[d.getDay()];
    
    return `${dayName}, ${day} ${month} ${year}`;
  });

  // Derived: Days grid for the calendar
  const calendarGrid = $derived.by(() => {
    const grid = [];
    // First day of current view month
    const firstDay = new Date(viewYear, viewMonth, 1);
    // Number of days in current view month
    const daysInCurMonth = new Date(viewYear, viewMonth + 1, 0).getDate();
    // Day of the week for the 1st day (0 = Sunday, 1 = Monday, etc.)
    const startDayOfWeek = firstDay.getDay();

    // Previous month filler days
    const prevMonthDays = new Date(viewYear, viewMonth, 0).getDate();
    for (let i = startDayOfWeek - 1; i >= 0; i--) {
      const dayNum = prevMonthDays - i;
      const prevMonth = viewMonth === 0 ? 11 : viewMonth - 1;
      const prevYear = viewMonth === 0 ? viewYear - 1 : viewYear;
      grid.push({
        day: dayNum,
        month: prevMonth,
        year: prevYear,
        isCurrentMonth: false,
        dateString: formatDateString(prevYear, prevMonth, dayNum)
      });
    }

    // Current month days
    for (let i = 1; i <= daysInCurMonth; i++) {
      grid.push({
        day: i,
        month: viewMonth,
        year: viewYear,
        isCurrentMonth: true,
        dateString: formatDateString(viewYear, viewMonth, i)
      });
    }

    // Next month filler days (to complete 42 days grid)
    const remainingDays = 42 - grid.length;
    for (let i = 1; i <= remainingDays; i++) {
      const nextMonth = viewMonth === 11 ? 0 : viewMonth + 1;
      const nextYear = viewMonth === 11 ? viewYear + 1 : viewYear;
      grid.push({
        day: i,
        month: nextMonth,
        year: nextYear,
        isCurrentMonth: false,
        dateString: formatDateString(nextYear, nextMonth, i)
      });
    }

    return grid;
  });

  // Helper: Format YYYY-MM-DD from parts
  function formatDateString(year: number, month: number, day: number): string {
    const mm = String(month + 1).padStart(2, '0');
    const dd = String(day).padStart(2, '0');
    return `${year}-${mm}-${dd}`;
  }

  // Action handlers
  function toggleCalendar() {
    if (disabled) return;
    isOpen = !isOpen;
  }

  function handleSelectDate(dateStr: string) {
    if (isDateDisabled(dateStr)) return;
    value = dateStr;
    isOpen = false;
    if (onchange) {
      onchange(dateStr);
    }
  }

  function handlePrevMonth() {
    if (viewMonth === 0) {
      viewMonth = 11;
      viewYear -= 1;
    } else {
      viewMonth -= 1;
    }
  }

  function handleNextMonth() {
    if (viewMonth === 11) {
      viewMonth = 0;
      viewYear += 1;
    } else {
      viewMonth += 1;
    }
  }

  function handleClear() {
    value = '';
    if (onchange) {
      onchange('');
    }
  }

  // Date limit checks
  function isDateDisabled(dateStr: string): boolean {
    if (minDate && dateStr < minDate) return true;
    if (maxDate && dateStr > maxDate) return true;
    return false;
  }

  // Click outside to close
  function handleDocumentClick(event: MouseEvent) {
    if (isOpen && containerEl && !containerEl.contains(event.target as Node)) {
      isOpen = false;
    }
  }

  // Listen to window click to handle click outside
  $effect(() => {
    document.addEventListener('click', handleDocumentClick, true);
    return () => {
      document.removeEventListener('click', handleDocumentClick, true);
    };
  });

  const uniqueId = 'datepicker-' + Math.random().toString(36).substring(2, 9);
</script>

<div class="relative w-full" bind:this={containerEl}>
  {#if label}
    <label for={uniqueId} class="block text-xs font-semibold text-gray-700 mb-1.5 uppercase tracking-wider">
      {label}
    </label>
  {/if}

  <div class="relative">
    <button
      id={uniqueId}
      type="button"
      {disabled}
      onclick={toggleCalendar}
      class="w-full flex items-center justify-between text-left border border-gray-200 rounded-xl pl-4 pr-10 py-3 bg-white text-sm text-gray-800 transition-all duration-200 hover:border-gray-300 focus:outline-none focus:border-amber-400 focus:ring-3 focus:ring-amber-100 disabled:opacity-60 disabled:bg-gray-50 disabled:cursor-not-allowed select-none"
      class:border-amber-400={isOpen}
      class:ring-3={isOpen}
      class:ring-amber-100={isOpen}
    >
      <div class="flex items-center gap-2.5 overflow-hidden">
        <Calendar class="w-4 h-4 text-gray-400 shrink-0" />
        {#if value}
          <span class="truncate text-gray-800 font-medium">{displayValue}</span>
        {:else}
          <span class="truncate text-gray-400">{placeholder}</span>
        {/if}
      </div>
    </button>

    {#if value && !disabled}
      <button
        type="button"
        class="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded-full text-gray-400 hover:text-gray-600 hover:bg-gray-100 shrink-0 transition-colors z-10 cursor-pointer"
        onclick={(e) => {
          e.stopPropagation();
          handleClear();
        }}
        aria-label="Bersihkan"
      >
        <X class="w-3.5 h-3.5" />
      </button>
    {/if}
  </div>

  {#if isOpen}
    <div
      class="absolute z-50 left-0 right-0 mt-2 p-4 bg-white rounded-xl border border-gray-200 shadow-xl select-none min-w-[320px] transition-all duration-300 transform origin-top-left"
    >
      <!-- Month & Year Navigation -->
      <div class="flex items-center justify-between mb-4">
        <button
          type="button"
          class="p-1.5 rounded-lg border border-gray-100 hover:bg-gray-50 text-gray-600 transition-colors cursor-pointer"
          onclick={handlePrevMonth}
          aria-label="Bulan sebelumnya"
        >
          <ChevronLeft class="w-4 h-4" />
        </button>
        
        <span class="font-bold text-gray-800 tracking-wide">
          {monthNames[viewMonth]} {viewYear}
        </span>
        
        <button
          type="button"
          class="p-1.5 rounded-lg border border-gray-100 hover:bg-gray-50 text-gray-600 transition-colors cursor-pointer"
          onclick={handleNextMonth}
          aria-label="Bulan berikutnya"
        >
          <ChevronRight class="w-4 h-4" />
        </button>
      </div>

      <!-- Day Names Header -->
      <div class="grid grid-cols-7 gap-1 text-center mb-2">
        {#each dayNames as day}
          <span class="text-xs font-bold text-gray-400 uppercase py-1">
            {day}
          </span>
        {/each}
      </div>

      <!-- Day Cells Grid -->
      <div class="grid grid-cols-7 gap-1">
        {#each calendarGrid as cell}
          {@const isSelected = value === cell.dateString}
          {@const isToday = formatDateString(today.getFullYear(), today.getMonth(), today.getDate()) === cell.dateString}
          {@const isInactive = !cell.isCurrentMonth}
          {@const isDisabled = isDateDisabled(cell.dateString)}
          {@const hasHighlight = highlightedDates.includes(cell.dateString)}

          <button
            type="button"
            disabled={isDisabled}
            onclick={() => handleSelectDate(cell.dateString)}
            class="relative aspect-square flex flex-col items-center justify-center rounded-lg text-xs font-semibold transition-all duration-150 cursor-pointer disabled:opacity-30 disabled:cursor-not-allowed"
            class:bg-amber-400={isSelected}
            class:text-white={isSelected}
            class:hover:bg-amber-500={isSelected}
            class:text-gray-800={!isSelected && !isInactive && !isDisabled}
            class:text-gray-300={isInactive && !isSelected}
            class:hover:bg-gray-100={!isSelected && !isDisabled}
            class:border={isToday && !isSelected}
            class:border-amber-400={isToday && !isSelected}
          >
            <span>{cell.day}</span>
            
            <!-- Highlight Indicator Dot -->
            {#if hasHighlight}
              <span
                class="absolute bottom-1 w-1.5 h-1.5 rounded-full transition-colors duration-200"
                class:bg-white={isSelected}
                class:bg-teal-500={!isSelected}
              ></span>
            {/if}
          </button>
        {/each}
      </div>

      <!-- Footer / Quick Action -->
      <div class="mt-4 pt-3 border-t border-gray-100 flex items-center justify-between">
        <button
          type="button"
          class="text-xs font-bold text-teal-600 hover:text-teal-700 transition-colors cursor-pointer"
          onclick={() => handleSelectDate(formatDateString(today.getFullYear(), today.getMonth(), today.getDate()))}
        >
          Hari Ini
        </button>
        {#if value}
          <span class="text-[10px] font-medium text-gray-400">
            Terpilih: {value}
          </span>
        {/if}
      </div>
    </div>
  {/if}
</div>
