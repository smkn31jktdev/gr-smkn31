<script lang="ts">
  let {
    status,
    waktuAbsen
  }: {
    status: string;
    waktuAbsen?: string;
  } = $props();

  const isTerlambat = $derived(status === 'hadir' && waktuAbsen && waktuAbsen.substring(0, 5) > '06:35');

  const badgeClass = $derived.by(() => {
    if (isTerlambat) {
      return 'bg-rose-50 text-rose-700 border-rose-100/40';
    }
    switch (status) {
      case 'hadir':
        return 'bg-teal-50 text-[#00a294] border-teal-100/30';
      case 'izin':
        return 'bg-slate-50 text-slate-650 border-slate-200/50';
      case 'sakit':
        return 'bg-amber-50/80 text-amber-700 border-amber-100/40';
      case 'magang':
        return 'bg-indigo-50/80 text-indigo-750 border-indigo-100/40';
      case 'belum':
        return 'bg-slate-100 text-slate-400 border-slate-200/50';
      default:
        return 'bg-rose-50 text-rose-700 border-rose-100/40';
    }
  });

  const displayLabel = $derived.by(() => {
    if (isTerlambat) return 'terlambat';
    if (status === 'tidak_hadir' || status === 'alpa') return 'tanpa keterangan';
    if (status === 'hadir') return 'masuk';
    return status;
  });
</script>

<span class="inline-block px-2.5 py-0.5 rounded-md border text-[9px] font-bold uppercase tracking-wider {badgeClass}">
  {displayLabel}
</span>
