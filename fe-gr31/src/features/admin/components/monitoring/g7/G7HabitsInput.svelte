<script lang="ts">
	import { onMount } from 'svelte';
	import { Lightbulb } from 'lucide-svelte';
	import type { SkorG7, G7SuggestResponse } from '../../../types/admin.types';
	import { SUB_INDICATORS } from '../../../../../const/g7';
	import { apiRequest } from '../../../../../api/client';

	let {
		skor = $bindable(),
		suggest,
		isReadOnly = false
	}: {
		skor: SkorG7;
		suggest: G7SuggestResponse | null;
		isReadOnly: boolean;
	} = $props();

	let ramadanPeriods = $state<{ start_date: string; end_date: string }[]>([]);

	onMount(async () => {
		try {
			const res = await apiRequest<{ ramadan: { start_date: string; end_date: string }[] }>('/v1/puasa/calendar');
			if (res.data && res.data.ramadan) {
				ramadanPeriods = res.data.ramadan;
			}
		} catch (e) {
			console.error(e);
		}
	});

	let isRamadan = $derived.by(() => {
		if (!suggest || !suggest.bulanTahun) return false;
		const targetBulan = suggest.bulanTahun; // YYYY-MM
		return ramadanPeriods.some(p => {
			const startMonth = p.start_date.substring(0, 7);
			const endMonth = p.end_date.substring(0, 7);
			return targetBulan === startMonth || targetBulan === endMonth;
		});
	});

	let filteredIndicators = $derived.by(() => {
		return SUB_INDICATORS.filter(ind => {
			if (ind.key === 'ibadahRowatib' || ind.key === 'ibadahTarawih' || ind.key === 'ibadahPuasa') {
				return isRamadan;
			}
			return true;
		});
	});
</script>

<div class="space-y-5 card p-6">
	<div class="border-b border-border pb-3">
		<h3 class="text-sm font-bold text-foreground">Skor Indikator Penilaian</h3>
		<p class="text-xxs mt-0.5 text-muted">
			Isi skor 0 jika tidak relevan, 1–5 jika relevan (Kurang s.d. Istimewa)
		</p>
	</div>

	<div class="space-y-4">
		{#each filteredIndicators as ind}
			<div
				class="flex flex-col justify-between gap-3 rounded-xl border bg-gray-50/50 p-3 sm:flex-row sm:items-center"
			>
				<div class="flex-1">
					<h4 class="text-xs font-bold text-foreground">{ind.label}</h4>
					<p class="text-xxs mt-0.5 text-muted">{ind.desc}</p>
					{#if suggest}
						<span
							class="mt-1.5 inline-flex items-center gap-1 rounded-md border border-amber-100 bg-amber-50 px-1.5 py-0.5 text-[9px] font-bold text-amber-800"
						>
							<Lightbulb class="w-3 h-3 text-amber-500 shrink-0" />
							Saran Auto-Suggest: {suggest.skor[ind.key as keyof SkorG7] ?? 0}
						</span>
					{/if}
				</div>

				<div class="flex shrink-0 items-center gap-2">
					<input
						type="number"
						min="0"
						max="5"
						bind:value={skor[ind.key as keyof SkorG7]}
						disabled={isReadOnly}
						class="input w-16 px-2 py-1 text-center text-xs font-bold"
					/>
					<span class="text-xxs font-bold text-muted">/ 5</span>
				</div>
			</div>
		{/each}
	</div>
</div>
