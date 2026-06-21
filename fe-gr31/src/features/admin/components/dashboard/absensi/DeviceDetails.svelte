<script lang="ts">
  import { Smartphone } from 'lucide-svelte';

  let { devInfo } = $props<{ devInfo: any }>();
</script>

<div class="space-y-3">
  <h5 class="text-[10px] font-black uppercase tracking-wider text-slate-400 text-left">Informasi Perangkat</h5>
  <div class="bg-white border border-slate-100 rounded-xl p-3.5 text-left shadow-xxs flex items-start gap-3">
    <Smartphone class="w-4.5 h-4.5 text-slate-400 mt-0.5 shrink-0" />
    <div class="space-y-1">
      <span class="text-[9px] font-bold text-slate-400 block uppercase">Device Identifikasi</span>
      <span class="text-[11px] font-bold text-slate-700 block">
        {devInfo.model || 'Unknown Device'}
      </span>
      <span class="text-[9px] font-medium text-slate-400 block">
        Platform: <span class="capitalize text-slate-600 font-bold">{devInfo.platform || 'Web'}</span> • OS: <span class="text-slate-600 font-bold">{devInfo.osVersion || 'Unknown'}</span>
      </span>
    </div>
  </div>
</div>
