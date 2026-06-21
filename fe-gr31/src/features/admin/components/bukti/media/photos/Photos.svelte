<script lang="ts">
  import { X, ChevronLeft, ChevronRight } from 'lucide-svelte';

  interface Props {
    open: boolean;
    photos: string[];
    studentName: string;
  }

  let {
    open = $bindable(),
    photos,
    studentName
  }: Props = $props();

  let activePhotoIndex = $state(0);

  // Reset indicator index when photos list changes
  $effect(() => {
    if (photos) {
      activePhotoIndex = 0;
    }
  });
</script>

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4 animate-fade-in">
    <div class="relative w-full max-w-4xl bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h3 class="text-sm font-black text-slate-700 uppercase">Dokumentasi Foto</h3>
          <p class="text-[10px] text-slate-400 font-bold mt-0.5">{studentName}</p>
        </div>
        <button 
          onclick={() => open = false}
          class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-655 transition-colors border-none cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Main Photo Area -->
      <div class="flex-1 bg-slate-50 flex items-center justify-center relative p-6 min-h-[300px] overflow-hidden">
        <!-- Main Image -->
        {#if photos && photos.length > 0}
          <img 
            src={photos[activePhotoIndex]} 
            alt="Bukti Foto" 
            class="max-w-full max-h-[50vh] sm:max-h-[60vh] object-contain rounded-xl shadow-xs" 
          />
        {/if}

        <!-- Prev button -->
        {#if photos && photos.length > 1}
          <button 
            onclick={() => activePhotoIndex = (activePhotoIndex - 1 + photos.length) % photos.length}
            class="absolute left-4 p-2.5 rounded-full bg-white/95 border border-slate-100 shadow-md text-slate-600 hover:text-slate-800 transition-all cursor-pointer border-none"
          >
            <ChevronLeft class="w-4 h-4" />
          </button>
          
          <!-- Next button -->
          <button 
            onclick={() => activePhotoIndex = (activePhotoIndex + 1) % photos.length}
            class="absolute right-4 p-2.5 rounded-full bg-white/95 border border-slate-100 shadow-md text-slate-600 hover:text-slate-800 transition-all cursor-pointer border-none"
          >
            <ChevronRight class="w-4 h-4" />
          </button>
        {/if}
      </div>

      <!-- Indicator Footer -->
      {#if photos && photos.length > 1}
        <div class="px-6 py-3 border-t border-slate-100 flex items-center justify-center gap-1.5 bg-white shrink-0">
          {#each photos as _, idx}
            <button 
              onclick={() => activePhotoIndex = idx}
              class="w-2 h-2 rounded-full transition-all border-none cursor-pointer {activePhotoIndex === idx ? 'bg-slate-800 w-4' : 'bg-slate-200 hover:bg-slate-300'}"
            ></button>
          {/each}
        </div>
      {/if}
    </div>
  </div>
{/if}
