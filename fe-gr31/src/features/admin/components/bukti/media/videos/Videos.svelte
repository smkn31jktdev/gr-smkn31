<script lang="ts">
  import { X, ExternalLink } from 'lucide-svelte';

  interface Props {
    open: boolean;
    videos: string[];
    studentName: string;
  }

  let {
    open = $bindable(),
    videos,
    studentName
  }: Props = $props();

  // YouTube Helper
  function getYouTubeEmbedUrl(url: string): string | null {
    if (!url) return null;
    const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*/;
    const match = url.match(regExp);
    if (match && match[2].length === 11) {
      return `https://www.youtube.com/embed/${match[2]}`;
    }
    return null;
  }
</script>

{#if open}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/80 backdrop-blur-sm p-4 animate-fade-in">
    <div class="relative w-full max-w-3xl bg-white rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
      <!-- Modal Header -->
      <div class="px-6 py-4 border-b border-slate-100 flex items-center justify-between">
        <div>
          <h3 class="text-sm font-black text-slate-700 uppercase">Dokumentasi Video</h3>
          <p class="text-[10px] text-slate-400 font-bold mt-0.5">{studentName}</p>
        </div>
        <button 
          onclick={() => open = false}
          class="p-1.5 rounded-lg hover:bg-slate-100 text-slate-400 hover:text-slate-655 transition-colors border-none cursor-pointer"
        >
          <X class="w-4 h-4" />
        </button>
      </div>

      <!-- Video Content List -->
      <div class="p-6 space-y-5 overflow-y-auto max-h-[70vh] custom-scrollbar">
        {#each videos as videoLink, idx}
          {@const embedUrl = getYouTubeEmbedUrl(videoLink)}
          <div class="space-y-2.5">
            <span class="block text-[10px] font-bold text-slate-400 uppercase tracking-widest">VIDEO #{idx + 1}</span>
            
            {#if embedUrl}
              <!-- Embedded Iframe -->
              <div class="aspect-video w-full rounded-2xl overflow-hidden border border-slate-100 bg-slate-50 shadow-xxs">
                <iframe 
                  src={embedUrl} 
                  title="YouTube video player" 
                  frameborder="0" 
                  allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" 
                  allowfullscreen 
                  class="w-full h-full"
                ></iframe>
              </div>
            {/if}

            <!-- External Link -->
            <a 
              href={videoLink} 
              target="_blank" 
              class="inline-flex items-center gap-1.5 text-xs text-[#00a294] font-bold hover:underline"
            >
              <ExternalLink class="w-3.5 h-3.5" />
              Buka di Tab Baru
            </a>
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  /* Custom scrollbar styling for a clean sleek feel */
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 99px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
  }
</style>
