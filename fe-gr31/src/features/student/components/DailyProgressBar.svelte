<script lang="ts">
	let { completedHabitsCount } = $props<{ completedHabitsCount: number }>();
</script>

<div
	class="flex min-h-[170px] flex-col justify-between rounded-3xl border border-slate-100/90 bg-white p-6 shadow-[0_10px_35px_rgba(0,0,0,0.01)] md:col-span-7"
>
	<div class="flex items-start justify-between gap-4">
		<!-- Progress title -->
		<div class="flex items-center gap-3">
			<div
				class="shadow-xxs flex h-11 w-11 items-center justify-center rounded-2xl border border-amber-100 bg-amber-50 text-amber-500"
			>
				<!-- Wave graph icon -->
				<svg
					class="h-5.5 w-5.5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
				</svg>
			</div>
			<div>
				<h3 class="text-sm font-black tracking-tight text-slate-800">Progres Harian</h3>
				<p class="mt-0.5 text-[11px] font-bold text-slate-400">Selesaikan kebiasaan baikmu</p>
			</div>
		</div>
		<!-- Ratio value -->
		<span class="text-2xl font-black text-slate-800"
			>{completedHabitsCount}<span class="text-xs font-bold text-slate-400">/7</span></span
		>
	</div>

	<!-- Linear progress bar -->
	<div class="my-4">
		<div class="h-2 w-full overflow-hidden rounded-full bg-slate-100">
			<div
				class="h-full rounded-full bg-[#4db6ac] transition-all duration-500 ease-out"
				style="width: {(completedHabitsCount / 7) * 100}%"
			></div>
		</div>
	</div>

	<!-- Advice subtext -->
	<div class="flex justify-end">
		<span class="text-[10px] font-bold tracking-wide text-[#4db6ac]">
			{#if completedHabitsCount >= 7}
				Kamu Hebat!
			{:else}
				Ayo selesaikan kegiatanmu!
			{/if}
		</span>
	</div>
</div>
