<script lang="ts">
  import { toasts, type ToastItem } from '../../../stores/uiStore';
  import { fly } from 'svelte/transition';
</script>

<div class="fixed top-4 right-4 z-50 flex flex-col gap-3 max-w-sm w-full pointer-events-none">
  {#each $toasts as toast (toast.id)}
    <div
      transition:fly={{ y: -20, duration: 300 }}
      class="pointer-events-auto p-4 rounded-xl shadow-lg border backdrop-blur-md flex items-start gap-3 transition-all duration-300"
      class:bg-blue-50={toast.type === 'info'}
      class:border-blue-200={toast.type === 'info'}
      class:text-blue-800={toast.type === 'info'}
      class:bg-emerald-50={toast.type === 'success'}
      class:border-emerald-200={toast.type === 'success'}
      class:text-emerald-800={toast.type === 'success'}
      class:bg-amber-50={toast.type === 'warning'}
      class:border-amber-200={toast.type === 'warning'}
      class:text-amber-800={toast.type === 'warning'}
      class:bg-rose-50={toast.type === 'error'}
      class:border-rose-200={toast.type === 'error'}
      class:text-rose-800={toast.type === 'error'}
    >
      <!-- Icon indicator -->
      <div class="mt-0.5 shrink-0">
        {#if toast.type === 'info'}
          <svg class="w-5 h-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        {:else if toast.type === 'success'}
          <svg class="w-5 h-5 text-emerald-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        {:else if toast.type === 'warning'}
          <svg class="w-5 h-5 text-amber-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        {:else if toast.type === 'error'}
          <svg class="w-5 h-5 text-rose-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        {/if}
      </div>

      <!-- Message -->
      <div class="flex-1 text-sm font-medium leading-relaxed">
        {toast.message}
      </div>

      <!-- Close button -->
      <button
        onclick={() => toasts.update($toasts => $toasts.filter(t => t.id !== toast.id))}
        class="shrink-0 p-0.5 rounded-lg hover:bg-black/5 text-gray-400 hover:text-gray-600 transition-colors"
        aria-label="Tutup"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  {/each}
</div>
