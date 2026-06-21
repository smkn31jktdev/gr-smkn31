<script lang="ts">
	import type { G7Jurnal } from '../../types/student.types';

	let { loadingJurnal, todayJurnal } = $props<{
		loadingJurnal: boolean;
		todayJurnal: G7Jurnal | null;
	}>();

	const habitsMetadata = [
		{ key: 'bangun' as const, name: 'Bangun Pagi', defaultTime: '-' },
		{ key: 'ibadah' as const, name: 'Beribadah', defaultTime: '-' },
		{ key: 'makan' as const, name: 'Makan Sehat', defaultTime: '-' },
		{ key: 'olahraga' as const, name: 'Olahraga', defaultTime: '-' },
		{ key: 'belajar' as const, name: 'Belajar', defaultTime: '-' },
		{ key: 'bermasyarakat' as const, name: 'Bermasyarakat', defaultTime: '-' },
		{ key: 'tidur' as const, name: 'Tidur', defaultTime: '22:00' }
	];
</script>

<!-- 2-Column Habits List -->
{#if loadingJurnal}
	<p class="py-10 text-center text-xs font-medium text-slate-400">Memuat data jurnal harian...</p>
{:else}
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
		{#each habitsMetadata as item}
			{@const entry = todayJurnal?.[item.key]}
			{@const done = entry?.done ?? false}

			<div
				class="flex items-center justify-between rounded-2xl border border-slate-100/70 bg-slate-50/25 p-3.5 font-sans transition-all duration-150"
			>
				<div class="flex items-center gap-3">
					<!-- Circular Icon Badge -->
					<div
						class="flex h-8 w-8 items-center justify-center rounded-full border text-sm transition-all"
						class:bg-emerald-50={done}
						class:border-emerald-100={done}
						class:text-emerald-500={done}
						class:bg-slate-50={!done}
						class:border-slate-100={!done}
						class:text-slate-300={!done}
					>
						{done ? '✓' : '×'}
					</div>
					<!-- Labels -->
					<div class="text-left font-sans">
						<p class="text-xs font-bold text-slate-700">{item.name}</p>
						<p class="mt-0.5 text-[10px] font-semibold text-slate-400">
							{#if done && entry?.waktu}
								Jam {entry.waktu}
							{:else}
								{item.defaultTime}
							{/if}
						</p>
					</div>
				</div>

				<!-- Status tag -->
				<span
					class="rounded-md px-2 py-0.5 text-[9px] font-extrabold tracking-wide"
					class:bg-emerald-100={done}
					class:text-emerald-600={done}
					class:bg-slate-100={!done}
					class:text-slate-400={!done}
				>
					{done ? 'SELESAI' : 'BELUM'}
				</span>
			</div>
		{/each}
	</div>

	{#if !todayJurnal}
		<p class="pt-2 text-center text-[10px] font-semibold text-slate-400">
			Belum ada jurnal untuk hari ini. Isi di menu
			<a href="/siswa/kegiatan" class="font-bold text-[#4db6ac] underline">Kegiatan</a>.
		</p>
	{/if}
{/if}
