<script lang="ts">
  import { Navigation, ShieldAlert, MapPin, ExternalLink } from 'lucide-svelte';
  import { getGoogleMapsLink } from './absensi';

  let { log } = $props<{ log: any }>();
</script>

<div class="space-y-3">
  <h5 class="text-[10px] font-black uppercase tracking-wider text-slate-400 text-left">Verifikasi Lokasi GPS</h5>
  
  <div class="grid grid-cols-2 gap-3">
    <div class="bg-white border border-slate-100 rounded-xl p-3 flex items-center gap-2.5 text-left shadow-xxs">
      <Navigation class="w-4 h-4 text-teal-500 shrink-0" />
      <div class="min-w-0">
        <span class="text-[9px] font-bold text-slate-400 block uppercase">Jarak ke Sekolah</span>
        <span class="text-[11px] font-black text-slate-800 block mt-0.5">
          {log.jarak !== undefined ? `${Math.round(log.jarak)} meter` : 'N/A'}
        </span>
      </div>
    </div>
    <div class="bg-white border border-slate-100 rounded-xl p-3 flex items-center gap-2.5 text-left shadow-xxs">
      <ShieldAlert class="w-4 h-4 text-teal-500 shrink-0" />
      <div class="min-w-0">
        <span class="text-[9px] font-bold text-slate-400 block uppercase">Akurasi GPS</span>
        <span class="text-[11px] font-black text-slate-800 block mt-0.5">
          {log.akurasi !== undefined ? `±${Math.round(log.akurasi)} meter` : 'N/A'}
        </span>
      </div>
    </div>
  </div>

  <!-- Coordinates & Google Maps Link -->
  {#if log.koordinat}
    <div class="bg-white border border-slate-100 rounded-xl p-3.5 text-left shadow-xxs space-y-2.5">
      <div class="flex items-center gap-2">
        <MapPin class="w-4 h-4 text-slate-455 shrink-0" />
        <div>
          <span class="text-[9px] font-bold text-slate-400 block uppercase">Koordinat Geofence</span>
          <span class="text-[10px] font-mono font-bold text-slate-600 block mt-0.5">
            {log.koordinat.lat.toFixed(7)}, {log.koordinat.lng.toFixed(7)}
          </span>
        </div>
      </div>
      
      <a
        href={getGoogleMapsLink(log.koordinat.lat, log.koordinat.lng)}
        target="_blank"
        class="flex items-center justify-center gap-1.5 w-full py-2 bg-slate-50 hover:bg-slate-100 border border-slate-200/60 rounded-lg text-[10px] font-extrabold text-slate-700 hover:text-slate-900 transition-colors no-underline cursor-pointer"
      >
        <ExternalLink class="w-3.5 h-3.5 text-slate-500" />
        Buka Lokasi di Google Maps
      </a>
    </div>
  {/if}
</div>
