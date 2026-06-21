<script lang="ts">
  let { 
    label = "Simpan", 
    loadingLabel = "Menyimpan...", 
    disabled = false, 
    className = "", 
    el = $bindable(),
    onclick 
  } = $props<{
    label?: string;
    loadingLabel?: string;
    disabled?: boolean;
    className?: string;
    el?: HTMLButtonElement;
    onclick: (handlers: { resolve: () => void; reject: () => void }) => Promise<void> | void;
  }>();

  let isSubmitting = $state(false);

  async function handleClick() {
    if (isSubmitting || disabled || !onclick) return;
    isSubmitting = true;
    
    try {
      await onclick({
        resolve: () => {
          setTimeout(() => { isSubmitting = false; }, 1500); // Reset after 1.5s delay
        },
        reject: () => {
          isSubmitting = false;
        }
      });
    } catch (e) {
      isSubmitting = false;
    }
  }
</script>

<button
  bind:this={el}
  type="button"
  onclick={handleClick}
  disabled={isSubmitting || disabled}
  class="btn-primary {className}"
  aria-busy={isSubmitting}
>
  {#if isSubmitting}
    <span class="flex items-center justify-center gap-2">
      <span class="spinner" aria-hidden="true"></span>
      {loadingLabel}
    </span>
  {:else}
    {label}
  {/if}
</button>
