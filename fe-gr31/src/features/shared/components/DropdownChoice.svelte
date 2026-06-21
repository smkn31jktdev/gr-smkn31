<script lang="ts">
  import { ChevronDown, Search, X, Check } from 'lucide-svelte';

  // Interface for option objects
  export interface DropdownOption {
    value: string | number;
    label: string;
    description?: string;
    badge?: string;
    badgeColor?: string; // e.g. bg-blue-100 text-blue-800
    avatar?: string; // Image URL
    icon?: any; // Lucide Component
  }

  // Props utilizing Svelte 5 Runes
  let {
    options = [],
    value = $bindable(), // Can be single value or array for multiple
    multiple = false,
    searchable = false,
    placeholder = 'Pilih opsi...',
    label = '',
    disabled = false,
    align = 'left',
    onchange
  } = $props<{
    options: (DropdownOption | string)[];
    value: any;
    multiple?: boolean;
    searchable?: boolean;
    placeholder?: string;
    label?: string;
    disabled?: boolean;
    align?: 'left' | 'right';
    onchange?: (value: any) => void;
  }>();

  // Normalize options to DropdownOption objects
  const normalizedOptions = $derived.by(() => {
    return options.map((opt: DropdownOption | string) => {
      if (typeof opt === 'string') {
        return { value: opt, label: opt } as DropdownOption;
      }
      return opt;
    });
  });

  // State using Runes
  let isOpen = $state(false);
  let searchQuery = $state('');
  let focusedIndex = $state(-1);
  let containerEl = $state<HTMLDivElement | null>(null);

  // Filtered options based on search query
  const filteredOptions = $derived.by(() => {
    if (!searchQuery) return normalizedOptions;
    const query = searchQuery.toLowerCase().trim();
    return normalizedOptions.filter((opt: DropdownOption) => 
      opt.label.toLowerCase().includes(query) || 
      (opt.description && opt.description.toLowerCase().includes(query))
    );
  });

  // Handle single vs multiple initial values safely
  $effect(() => {
    if (multiple && !Array.isArray(value)) {
      value = value ? [value] : [];
    } else if (!multiple && Array.isArray(value)) {
      value = value[0] || '';
    }
  });

  // Keyboard navigation & helpers
  function toggleDropdown() {
    if (disabled) return;
    isOpen = !isOpen;
    if (isOpen) {
      focusedIndex = -1;
      searchQuery = '';
    }
  }

  function handleSelect(optionVal: string | number) {
    if (multiple) {
      const arr = Array.isArray(value) ? [...value] : [];
      const idx = arr.indexOf(optionVal);
      if (idx > -1) {
        arr.splice(idx, 1);
      } else {
        arr.push(optionVal);
      }
      value = arr;
    } else {
      value = optionVal;
      isOpen = false;
    }
    
    if (onchange) {
      onchange(value);
    }
  }

  function removeMultiValue(val: string | number, e: MouseEvent) {
    e.stopPropagation();
    if (disabled) return;
    if (Array.isArray(value)) {
      value = value.filter(v => v !== val);
      if (onchange) onchange(value);
    }
  }

  // Keyboard events
  function handleKeyDown(e: KeyboardEvent) {
    if (disabled) return;
    
    if (!isOpen) {
      if (e.key === 'ArrowDown' || e.key === 'ArrowUp' || e.key === 'Enter') {
        e.preventDefault();
        isOpen = true;
        focusedIndex = 0;
      }
      return;
    }

    switch (e.key) {
      case 'Escape':
        e.preventDefault();
        isOpen = false;
        break;
      case 'ArrowDown':
        e.preventDefault();
        focusedIndex = (focusedIndex + 1) % filteredOptions.length;
        break;
      case 'ArrowUp':
        e.preventDefault();
        focusedIndex = (focusedIndex - 1 + filteredOptions.length) % filteredOptions.length;
        break;
      case 'Enter':
        e.preventDefault();
        if (focusedIndex >= 0 && focusedIndex < filteredOptions.length) {
          handleSelect(filteredOptions[focusedIndex].value);
        }
        break;
      case 'Tab':
        isOpen = false;
        break;
    }
  }

  // Helper to check selection state
  function isSelected(optionVal: string | number): boolean {
    if (multiple && Array.isArray(value)) {
      return value.includes(optionVal);
    }
    return value === optionVal;
  }

  // Derived selected labels
  const selectedLabel = $derived.by(() => {
    if (multiple) return '';
    const opt = normalizedOptions.find((o: DropdownOption) => o.value === value);
    return opt ? opt.label : '';
  });

  const selectedOptions = $derived.by(() => {
    if (!multiple || !Array.isArray(value)) return [];
    return normalizedOptions.filter((o: DropdownOption) => value.includes(o.value));
  });

  // Click outside listener
  function handleDocumentClick(event: MouseEvent) {
    if (isOpen && containerEl && !containerEl.contains(event.target as Node)) {
      isOpen = false;
    }
  }

    $effect(() => {
    document.addEventListener('click', handleDocumentClick, true);
    return () => {
      document.removeEventListener('click', handleDocumentClick, true);
    };
  });

  const uniqueId = 'dropdown-' + Math.random().toString(36).substring(2, 9);
</script>

<div class="relative w-full" bind:this={containerEl} onkeydown={handleKeyDown} role="presentation">
  {#if label}
    <label for={uniqueId} class="block text-xs font-semibold text-gray-700 mb-1.5 uppercase tracking-wider">
      {label}
    </label>
  {/if}

  <!-- Dropdown Trigger Input -->
  <div class="relative">
    <div
      id={uniqueId}
      tabindex={disabled ? -1 : 0}
      role="combobox"
      aria-expanded={isOpen}
      aria-haspopup="listbox"
      aria-controls="dropdown-options"
      onclick={toggleDropdown}
      onkeydown={(e) => {
        if (e.key === ' ' || e.key === 'Enter') {
          e.preventDefault();
          toggleDropdown();
        }
      }}
      class="w-full flex items-center justify-between text-left border border-gray-200 rounded-xl px-4 py-3 bg-white text-sm text-gray-800 transition-all duration-200 hover:border-gray-300 focus:outline-none focus:border-[#4db6ac] focus:ring-3 focus:ring-teal-50 select-none min-h-[46px] cursor-pointer"
      class:opacity-60={disabled}
      class:bg-gray-50={disabled}
      class:cursor-not-allowed={disabled}
      class:border-[#4db6ac]={isOpen}
      class:ring-3={isOpen}
      class:ring-teal-50={isOpen}
    >
      <div class="flex flex-wrap gap-1.5 items-center overflow-hidden flex-1 mr-2 bg-transparent">
        {#if multiple}
          {#if selectedOptions.length > 0}
            {#each selectedOptions as selectedOpt}
              <span class="inline-flex items-center gap-1 bg-teal-50 border border-teal-200/60 text-teal-900 text-xs font-bold px-2.5 py-1 rounded-lg">
                {#if selectedOpt.avatar}
                  <img src={selectedOpt.avatar} alt={selectedOpt.label} class="w-3.5 h-3.5 rounded-full object-cover" />
                {/if}
                {selectedOpt.label}
                <button
                  type="button"
                  class="text-[#4db6ac] hover:text-[#3ca59b] transition-colors p-0.5 cursor-pointer border-none bg-transparent"
                  onclick={(e) => removeMultiValue(selectedOpt.value, e)}
                >
                  <X class="w-3 h-3" />
                </button>
              </span>
            {/each}
          {:else}
            <span class="text-gray-400">{placeholder}</span>
          {/if}
        {:else}
          {#if value !== undefined && value !== null && value !== ''}
            {@const activeOpt = normalizedOptions.find((o: DropdownOption) => o.value === value)}
            <div class="flex items-center gap-2">
              {#if activeOpt?.avatar}
                <img src={activeOpt.avatar} alt={activeOpt.label} class="w-5 h-5 rounded-full object-cover shrink-0" />
              {/if}
              {#if activeOpt?.icon}
                {@const IconComponent = activeOpt.icon}
                <IconComponent class="w-4 h-4 text-gray-500 shrink-0" />
              {/if}
              <span class="font-medium text-gray-800">{selectedLabel}</span>
            </div>
          {:else}
            <span class="text-gray-400">{placeholder}</span>
          {/if}
        {/if}
      </div>

      <ChevronDown class="w-4 h-4 text-gray-400 transition-transform duration-200 shrink-0 {isOpen ? 'rotate-180' : ''}" />
    </div>
  </div>

  <!-- Dropdown Options Panel -->
  {#if isOpen}
    <div 
      id="dropdown-options" 
      class="absolute z-50 mt-2 p-2 bg-white rounded-xl border border-gray-200 shadow-xl select-none max-h-72 overflow-hidden flex flex-col transition-all duration-300 min-w-full w-max max-w-[280px] sm:max-w-[360px]"
      class:left-0={align === 'left'}
      class:right-0={align === 'right'}
    >
      
      <!-- Search Input -->
      {#if searchable}
        <div class="relative mb-2 shrink-0">
          <Search class="absolute left-3 top-2.5 w-4 h-4 text-gray-400" />
          <input
            type="text"
            placeholder="Cari..."
            bind:value={searchQuery}
            class="w-full text-xs border border-gray-200 rounded-lg pl-9 pr-4 py-2 bg-gray-50 text-gray-800 placeholder-gray-400 focus:outline-none focus:border-[#4db6ac] focus:bg-white"
          />
        </div>
      {/if}

      <!-- Options List -->
      <div class="overflow-y-auto flex-1 flex flex-col gap-0.5 max-h-52">
        {#each filteredOptions as option, idx}
          {@const selected = isSelected(option.value)}
          {@const focused = focusedIndex === idx}
          
          <button
            type="button"
            onclick={() => handleSelect(option.value)}
            class="w-full flex items-center justify-between text-left px-3 py-2.5 rounded-lg text-xs font-semibold transition-colors cursor-pointer"
            class:bg-teal-50={selected}
            class:text-teal-955={selected}
            class:bg-gray-50={focused && !selected}
            class:text-gray-800={!selected}
            class:hover:bg-gray-50={!selected}
          >
            <div class="flex items-center gap-2.5 flex-1 min-w-0">
              {#if option.avatar}
                <img src={option.avatar} alt={option.label} class="w-6 h-6 rounded-full object-cover shrink-0" />
              {/if}
              {#if option.icon}
                {@const IconComponent = option.icon}
                <IconComponent class="w-4 h-4 text-gray-500 shrink-0" />
              {/if}

              <div class="flex flex-col min-w-0">
                <span class="truncate">{option.label}</span>
                {#if option.description}
                  <span class="text-[10px] font-medium text-gray-400 truncate mt-0.5">{option.description}</span>
                {/if}
              </div>
            </div>

            <div class="flex items-center gap-2 ml-2 shrink-0">
              {#if option.badge}
                <span class="inline-flex items-center text-[9px] font-black uppercase px-2 py-0.5 rounded-md {option.badgeColor || 'bg-gray-100 text-gray-700'}">
                  {option.badge}
                </span>
              {/if}
              {#if selected}
                <Check class="w-4 h-4 text-[#4db6ac]" />
              {/if}
            </div>
          </button>
        {:else}
          <div class="py-6 text-center text-xs font-bold text-gray-400">
            Tidak ada opsi ditemukan
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
