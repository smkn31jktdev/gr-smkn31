<script lang="ts">
	import type { G7Rekap } from '../../../types/admin.types';

	let { rekap }: { rekap: G7Rekap | null } = $props();
</script>

{#if rekap}
	<!-- Profile summary banner -->
	<div class="card border-none bg-linear-to-r from-primary to-primary-hover p-6 text-white">
		<h2 class="font-display text-xl font-bold uppercase">{rekap.namaSiswa}</h2>
		<p class="mt-1 text-xs font-semibold text-blue-100">
			NIS: {rekap.nis} | Kelas: {rekap.kelas}
		</p>
		<div class="mt-6 grid grid-cols-3 gap-4 border-t border-white/20 pt-4 text-center text-xs">
			<div>
				<p class="font-semibold text-blue-100/70">Hari Tercatat</p>
				<p class="mt-1 text-lg font-black">{rekap.hariTercatat} Hari</p>
			</div>
			<div>
				<p class="font-semibold text-blue-100/70">Nilai Akhir Saat Ini</p>
				<p class="mt-1 text-lg font-black">{rekap.nilaiAkhir.toFixed(2)}</p>
			</div>
			<div>
				<p class="font-semibold text-blue-100/70">Predikat</p>
				<p class="mt-1 text-lg font-black">{rekap.predikat || 'Kurang'}</p>
			</div>
		</div>
	</div>
{/if}
