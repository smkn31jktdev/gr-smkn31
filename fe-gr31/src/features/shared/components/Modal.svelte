<script lang="ts">
  import { fade, scale } from 'svelte/transition';
  
  let { 
    show = false, 
    title = "", 
    onclose, 
    children,
    footer 
  } = $props<{
    show: boolean;
    title?: string;
    onclose: () => void;
    children?: import('svelte').Snippet;
    footer?: import('svelte').Snippet;
  }>();
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <!-- Backdrop -->
    <div 
      onclick={onclose}
      transition:fade={{ duration: 200 }} 
      class="fixed inset-0 bg-black/60 backdrop-blur-xs cursor-pointer"
      role="presentation"
    ></div>

    <!-- Modal Content -->
    <div
      transition:scale={{ duration: 200, start: 0.95 }}
      class="bg-surface border border-border w-full max-w-lg rounded-2xl shadow-xl z-10 flex flex-col max-h-[90vh] overflow-hidden"
    >
      <!-- Header -->
      {#if title}
        <div class="p-5 border-b border-border flex items-center justify-between">
          <h3 class="text-lg font-bold text-foreground">{title}</h3>
          <button
            onclick={onclose}
            class="p-1 rounded-lg hover:bg-gray-100 text-muted hover:text-foreground transition-colors"
            aria-label="Tutup"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      {:else}
        <button
          onclick={onclose}
          class="absolute top-4 right-4 p-1.5 rounded-lg hover:bg-gray-100 text-muted hover:text-foreground transition-colors z-20"
          aria-label="Tutup"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      {/if}

      <!-- Body -->
      <div class="p-6 overflow-y-auto flex-1 text-sm text-foreground/95 leading-relaxed">
        {@render children?.()}
      </div>

      <!-- Footer -->
      {#if footer}
        <div class="p-5 border-t border-border bg-gray-50 flex items-center justify-end gap-3">
          {@render footer()}
        </div>
      {/if}
    </div>
  </div>
{/if}
