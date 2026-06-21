<script lang="ts">
  import { getUploadUrl } from '../../../../../api/client';

  let { log } = $props<{ log: any }>();
</script>

<div class="space-y-3">
  <h5 class="text-[10px] font-black uppercase tracking-wider text-slate-400 text-left">Detail Izin / Sakit</h5>
  
  <!-- Reason Block -->
  {#if log.alasan}
    <div class="bg-amber-50/50 border border-amber-100/60 rounded-xl p-3.5 text-left">
      <span class="text-[9px] font-black text-amber-800 uppercase block tracking-wider">Keterangan / Alasan:</span>
      <p class="text-slate-650 leading-relaxed font-semibold italic text-xs mt-1">
        "{log.alasan}"
      </p>
    </div>
  {/if}

  <!-- Document Image Preview -->
  {#if log.fotoIzin}
    <div class="space-y-2">
      <span class="text-[9px] font-black text-slate-400 uppercase block tracking-wider text-left">Surat Bukti Lampiran:</span>
      <div class="relative flex min-h-[220px] max-h-[350px] w-full items-center justify-center overflow-hidden rounded-xl border border-slate-200/50 bg-slate-50 p-2 shadow-[inset_0_1px_2px_rgba(0,0,0,0.02)]">
        <img
          src={getUploadUrl(log.fotoIzin)}
          alt="Surat Keterangan {log.namaSiswa}"
          class="max-h-[330px] max-w-full rounded-lg object-contain shadow-sm"
        />
      </div>
    </div>
  {:else}
    <div class="bg-white border border-slate-100 rounded-xl p-8 text-center text-xs font-bold text-slate-400">
      Tidak ada berkas surat bukti yang dilampirkan.
    </div>
  {/if}
</div>
