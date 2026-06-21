<script lang="ts">
  import { Download, ExternalLink } from 'lucide-svelte';
  import { getUploadUrl } from '../../../../../api/client';

  let { log, onclose } = $props<{
    log: any;
    onclose: () => void;
  }>();
</script>

<div class="flex items-center justify-end gap-2 border-t border-slate-100 bg-slate-50/50 px-5 py-3.5">
  {#if (log.status === 'izin' || log.status === 'sakit') && log.fotoIzin}
    <a
      href={getUploadUrl(log.fotoIzin)}
      download="Surat_Izin_{log.namaSiswa.replace(/\s+/g, '_')}"
      target="_blank"
      class="border-slate-200 hover:border-slate-300 text-slate-600 shadow-xxs flex items-center gap-1.5 rounded-xl border bg-white px-3.5 py-2 text-xs font-bold no-underline transition-colors hover:bg-slate-50 cursor-pointer"
    >
      <Download class="h-3.5 w-3.5" />
      Unduh
    </a>
    <a
      href={getUploadUrl(log.fotoIzin)}
      target="_blank"
      class="flex items-center gap-1.5 rounded-xl border-none bg-slate-800 px-3.5 py-2 text-xs font-bold text-white no-underline shadow-xs transition-all hover:bg-slate-900 active:scale-98 cursor-pointer"
    >
      <ExternalLink class="h-3.5 w-3.5 text-white" />
      Buka Tab Baru
    </a>
  {/if}
  <button
    onclick={onclose}
    class="border-slate-200 cursor-pointer rounded-xl border bg-white px-4 py-2 text-xs font-bold text-slate-600 transition-colors hover:bg-slate-50"
  >
    Tutup
  </button>
</div>
