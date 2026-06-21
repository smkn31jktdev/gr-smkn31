<script lang="ts">
	interface CalendarDay {
		day: number | null;
		dateStr: string;
		status: string;
		waktu: string;
	}

	let { calendarDays, loadingHistory, monthNameIndonesian, onrefresh } = $props<{
		calendarDays: CalendarDay[];
		loadingHistory: boolean;
		monthNameIndonesian: string;
		onrefresh: () => void;
	}>();
</script>

<div
	class="flex flex-col justify-between rounded-3xl border border-slate-100/90 bg-white p-6 shadow-[0_10px_35px_rgba(0,0,0,0.01)] lg:col-span-5"
>
	<div>
		<div class="mb-4 flex items-center justify-between border-b border-slate-50 pb-4">
			<div class="text-left">
				<h3 class="text-sm font-black tracking-tight text-slate-800">Kehadiran Bulanan</h3>
				<p class="mt-0.5 text-[10px] font-bold text-slate-400">{monthNameIndonesian}</p>
			</div>
			<button
				onclick={onrefresh}
				disabled={loadingHistory}
				class="cursor-pointer border-none bg-transparent text-[10px] font-extrabold text-[#4db6ac] hover:underline disabled:opacity-50"
			>
				{loadingHistory ? 'Loading...' : 'Segarkan'}
			</button>
		</div>

		<!-- Weekly Headers -->
		<div
			class="mb-2.5 grid grid-cols-7 gap-2 text-center text-[9px] font-extrabold text-slate-400 uppercase"
		>
			<div>Min</div>
			<div>Sen</div>
			<div>Sel</div>
			<div>Rab</div>
			<div>Kam</div>
			<div>Jum</div>
			<div>Sab</div>
		</div>

		<!-- Day Cells Grid -->
		{#if loadingHistory}
			<p class="py-16 text-center text-xs font-medium text-slate-400">Memuat kalender...</p>
		{:else}
			<div class="grid grid-cols-7 gap-1.5">
				{#each calendarDays as cell}
					{#if cell.status === 'empty'}
						<div
							class="aspect-square rounded-xl border border-transparent bg-slate-50/35"
						></div>
					{:else}
						<div
							class="group relative flex aspect-square cursor-pointer flex-col items-center justify-between rounded-xl border px-0.5 py-1 transition-all
								{cell.status === 'hadir'
									? 'border-emerald-100 bg-emerald-50 text-emerald-600'
									: ''}
								{cell.status === 'izin' ? 'border-sky-100 bg-sky-50 text-sky-600' : ''}
								{cell.status === 'sakit' ? 'border-amber-100 bg-amber-50 text-amber-600' : ''}
								{cell.status === 'magang' ? 'border-purple-100 bg-purple-50 text-purple-600' : ''}
								{cell.status === 'tidak_hadir'
									? 'border-rose-100 bg-rose-50 text-rose-600'
									: ''}
								{cell.status === 'belum_absen'
									? 'border-slate-100/70 bg-slate-50 text-slate-400'
									: ''}
								{cell.status === 'future'
									? 'border-transparent bg-slate-50/20 text-slate-200'
									: ''}"
						>
							<span class="text-[10px] font-extrabold">{cell.day}</span>

							<!-- Cell indicator dots -->
							{#if cell.status === 'hadir'}
								<span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span>
							{:else if cell.status === 'izin'}
								<span class="h-1.5 w-1.5 rounded-full bg-sky-400"></span>
							{:else if cell.status === 'sakit'}
								<span class="h-1.5 w-1.5 rounded-full bg-amber-400"></span>
							{:else if cell.status === 'magang'}
								<span class="h-1.5 w-1.5 rounded-full bg-purple-500"></span>
							{:else if cell.status === 'tidak_hadir'}
								<span class="h-1.5 w-1.5 rounded-full bg-rose-400"></span>
							{:else}
								<span class="h-1.5 w-1.5 rounded-full bg-transparent"></span>
							{/if}

							<!-- Cell Details Tooltip -->
							{#if cell.status !== 'future' && cell.status !== 'belum_absen'}
								<div
									class="pointer-events-none absolute bottom-full z-20 mb-1 rounded bg-slate-900/90 px-2 py-0.5 text-[9px] font-bold whitespace-nowrap text-white opacity-0 shadow-sm transition-opacity group-hover:opacity-100"
								>
									{cell.status.toUpperCase()}
									{cell.waktu ? '(' + cell.waktu.substring(0, 5) + ')' : ''}
								</div>
							{/if}
						</div>
					{/if}
				{/each}
			</div>
		{/if}
	</div>

	<!-- Legend Indicator explanation footer -->
	<div
		class="mt-4 flex flex-wrap justify-center gap-3 border-t border-slate-50 pt-3 text-[9px] font-bold text-slate-400"
	>
		<span class="flex items-center gap-1.5"
			><span class="h-1.5 w-1.5 rounded-full bg-emerald-500"></span> Hadir</span
		>
		<span class="flex items-center gap-1.5"
			><span class="h-1.5 w-1.5 rounded-full bg-sky-400"></span> Izin</span
		>
		<span class="flex items-center gap-1.5"
			><span class="h-1.5 w-1.5 rounded-full bg-amber-400"></span> Sakit</span
		>
		<span class="flex items-center gap-1.5"
			><span class="h-1.5 w-1.5 rounded-full bg-purple-500"></span> Magang</span
		>
		<span class="flex items-center gap-1.5"
			><span class="h-1.5 w-1.5 rounded-full bg-rose-400"></span> Alpa</span
		>
	</div>
</div>
