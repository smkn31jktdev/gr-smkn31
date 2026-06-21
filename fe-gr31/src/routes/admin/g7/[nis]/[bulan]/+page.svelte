<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { G7DetailState } from '../../../../../features/admin/logic/g7DetailLogic.svelte';
	import G7DetailCard from '../../../../../features/admin/components/monitoring/card/G7DetailCard.svelte';
	import G7AssessorCard from '../../../../../features/admin/components/monitoring/card/G7AssessorCard.svelte';
	import G7HabitsInput from '../../../../../features/admin/components/monitoring/g7/G7HabitsInput.svelte';

	// Instantiate the reactive state
	const state = new G7DetailState(
		$page.params.nis ?? '',
		$page.params.bulan ?? ''
	);

	onMount(() => {
		state.loadDetail();
	});
</script>

<div class="space-y-6">
	<!-- Header Card -->
	<div class="flex items-center justify-between">
		<a
			href="/admin/g7"
			class="flex items-center gap-1 text-xs font-bold text-muted hover:text-foreground hover:underline"
		>
			← Kembali ke Daftar Rekap
		</a>
		{#if state.isReadOnly}
			<span
				class="text-xxs rounded-full bg-emerald-100 px-3 py-1 font-bold text-emerald-800 uppercase"
			>
				Status: Final (Terkunci)
			</span>
		{/if}
	</div>

	{#if state.loading}
		<div class="flex flex-col items-center justify-center p-12 text-muted">
			<span class="spinner border-top-primary mb-3 h-8 w-8 border-2 border-primary/20"></span>
			<p class="text-sm font-semibold">Memuat penilaian rekap...</p>
		</div>
	{:else if state.rekap}
		<!-- Profile summary banner -->
		<G7DetailCard rekap={state.rekap} />

		<!-- Forms Layout -->
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
			<!-- Left side: 18 indicators inputs -->
			<div class="space-y-4 lg:col-span-2">
				<G7HabitsInput
					bind:skor={state.skor}
					suggest={state.suggest}
					isReadOnly={state.isReadOnly}
				/>
			</div>

			<!-- Right side: Assessors metadata and status -->
			<div class="space-y-6">
				<G7AssessorCard
					bind:waliKelas={state.waliKelas}
					bind:guruBK={state.guruBK}
					bind:status={state.status}
					isReadOnly={state.isReadOnly}
					handleSave={(handlers) => state.handleSave(handlers)}
				/>
			</div>
		</div>
	{/if}
</div>
