<script lang="ts">
  import { Search, X } from 'lucide-svelte';

  let {
    value = $bindable(''),
    placeholder = 'Cari...',
    class: wrapperClass = '',
    inputClass = '',
    rounded = 'xl',
    size = 'md',
    showClear = true,
    oninput,
    onclear
  }: {
    value?: string;
    placeholder?: string;
    class?: string;
    inputClass?: string;
    rounded?: 'full' | 'xl';
    size?: 'sm' | 'md';
    showClear?: boolean;
    oninput?: (e: Event) => void;
    onclear?: () => void;
  } = $props();

  function handleInput(e: Event) {
    const target = e.target as HTMLInputElement;
    value = target.value;
    if (oninput) {
      oninput(e);
    }
  }

  function handleClear() {
    value = '';
    if (onclear) {
      onclear();
    } else if (oninput) {
      // Trigger input event after clear if oninput is provided
      const dummyEvent = new Event('input', { bubbles: true });
      oninput(dummyEvent);
    }
  }
</script>

<div class="relative {wrapperClass}">
  <Search 
    class="text-slate-400 absolute left-3.5 top-1/2 -translate-y-1/2 
           {size === 'sm' ? 'w-3.5 h-3.5' : 'w-4 h-4'}" 
  />
  
  <input
    type="text"
    {placeholder}
    {value}
    oninput={handleInput}
    class="w-full bg-slate-50/60 hover:bg-slate-50 border border-slate-200/50 hover:border-slate-300/80 
           focus:border-[#00a294] focus:ring-3 focus:ring-teal-50 focus:bg-white outline-none transition-all
           {rounded === 'full' ? 'rounded-full' : 'rounded-xl'}
           {size === 'sm' ? 'pl-9 pr-8 py-1.5 text-xs font-semibold' : 'pl-11 pr-10 py-2.5 text-xs'}
           {inputClass}"
  />

  {#if showClear && value}
    <button
      type="button"
      onclick={handleClear}
      class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 bg-transparent border-none cursor-pointer p-0.5 rounded transition-colors"
      aria-label="Clear search"
    >
      <X class="{size === 'sm' ? 'w-3 h-3' : 'w-3.5 h-3.5'}" />
    </button>
  {/if}
</div>
