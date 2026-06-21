<script lang="ts">
	import { X, Download, ExternalLink } from 'lucide-svelte';

	// Svelte 5 Props
	let {
		show = false,
		fotoUrl = '',
		namaSiswa = '',
		kelas = '',
		status = '',
		tanggal = '',
		alasan = '',
		onclose
	} = $props<{
		show: boolean;
		fotoUrl: string;
		namaSiswa: string;
		kelas: string;
		status: string;
		tanggal: string;
		alasan: string;
		onclose: () => void;
	}>();

	// Escape key listener to close modal
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && show) {
			onclose();
		}
	}

	$effect(() => {
		if (show) {
			window.addEventListener('keydown', handleKeydown);
		}
		return () => {
			window.removeEventListener('keydown', handleKeydown);
		};
	});
</script>

{#if show}
	<!-- Overlay -->
	<div
		class="animate-fade-in fixed inset-0 z-50 flex items-center justify-center bg-slate-900/65 p-4 text-slate-700 backdrop-blur-xs"
		onclick={onclose}
		onkeydown={(e) => {
			if (e.key === 'Enter' || e.key === ' ') onclose();
		}}
		role="button"
		tabindex="0"
	>
		<!-- Modal Box -->
		<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
		<div
			class="flex max-h-[85vh] w-full max-w-lg scale-100 transform flex-col overflow-hidden rounded-2xl bg-white shadow-2xl transition-all duration-300"
			onclick={(e) => e.stopPropagation()}
			onkeydown={() => {}}
			role="document"
			tabindex="-1"
		>
			<!-- Header -->
			<div
				class="flex items-center justify-between border-b border-slate-100 bg-slate-50/50 px-5 py-4"
			>
				<div class="text-left">
					<h3 class="text-xs font-bold tracking-wide text-slate-800 uppercase">
						Lampiran Surat {status === 'sakit' ? 'Sakit' : 'Izin'}
					</h3>
					<p class="mt-0.5 text-[10px] font-semibold text-slate-400">
						{namaSiswa} ({kelas}) • {tanggal}
					</p>
				</div>
				<button
					onclick={onclose}
					class="hover:text-slate-655 cursor-pointer rounded-lg border-none bg-transparent p-1.5 text-slate-400 transition-colors hover:bg-slate-100"
					title="Tutup"
				>
					<X class="h-4.5 w-4.5" />
				</button>
			</div>

			<!-- Content -->
			<div class="flex flex-1 flex-col gap-4 overflow-y-auto bg-slate-50/10 p-5">
				{#if alasan}
					<div class="rounded-xl border border-amber-100/60 bg-amber-50/60 p-3.5 text-left text-xs">
						<span class="mb-1 block font-bold text-amber-800">Alasan Ketidakhadiran:</span>
						<p class="text-slate-650 leading-relaxed font-medium italic">"{alasan}"</p>
					</div>
				{/if}

				{#if fotoUrl}
					<div
						class="relative flex min-h-[300px] w-full items-center justify-center overflow-hidden rounded-xl border border-slate-200/50 bg-slate-50 p-2 shadow-[inset_0_1px_2px_rgba(0,0,0,0.02)]"
					>
						<img
							src={fotoUrl}
							alt="Lampiran Surat {namaSiswa}"
							class="max-h-[50vh] max-w-full rounded-lg object-contain shadow-sm"
						/>
					</div>
				{:else}
					<div class="p-12 text-center font-medium text-slate-400">
						Tidak ada dokumen surat yang dilampirkan.
					</div>
				{/if}
			</div>

			<!-- Footer Actions -->
			<div
				class="flex items-center justify-end gap-2.5 border-t border-slate-100 bg-slate-50/50 px-5 py-3.5"
			>
				{#if fotoUrl}
					<a
						href={fotoUrl}
						download="Surat_Izin_{namaSiswa.replace(/\s+/g, '_')}"
						target="_blank"
						class="border-slate-250 hover:border-slate-350 text-slate-650 shadow-xxs flex items-center gap-1.5 rounded-xl border bg-white px-3.5 py-2 text-xs font-bold no-underline transition-colors hover:bg-slate-50"
					>
						<Download class="h-3.5 w-3.5" />
						Unduh
					</a>
					<a
						href={fotoUrl}
						target="_blank"
						class="flex items-center gap-1.5 rounded-xl border-none bg-slate-800 px-3.5 py-2 text-xs font-bold text-white no-underline shadow-xs transition-all hover:bg-slate-900 active:scale-98"
					>
						<ExternalLink class="h-3.5 w-3.5 text-white" />
						Buka Tab Baru
					</a>
				{/if}
				<button
					onclick={onclose}
					class="border-slate-250 cursor-pointer rounded-xl border bg-white px-4 py-2 text-xs font-bold text-slate-600 transition-colors hover:bg-slate-50"
				>
					Tutup
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}
	.animate-fade-in {
		animation: fadeIn 0.2s ease-out forwards;
	}
</style>
