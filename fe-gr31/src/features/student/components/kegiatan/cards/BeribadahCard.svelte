<script lang="ts">
	import { Sparkles, Info, Loader } from 'lucide-svelte';

	let {
		doa_ortu = $bindable(false),
		sholat_fajar = $bindable(false),
		sholat_5waktu = $bindable(false),
		zikir = $bindable(false),
		dhuha = $bindable(false),
		rowatib = $bindable(false),
		infaq = $bindable(0),
		isIslam = true,
		isRamadan = false,
		ramadanDay = 0,
		onsave,
		loading = false
	} = $props();

	// Helper to handle numeric binding for Infaq input
	let infaqInputVal = $state(infaq || '');

	$effect(() => {
		infaq = Number(infaqInputVal) || 0;
	});
</script>

<div
	class="space-y-6 rounded-3xl border border-slate-100/90 bg-white p-6 shadow-[0_8px_30px_rgb(0,0,0,0.015)]"
>
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-3">
			<div
				class="flex h-9 w-9 items-center justify-center rounded-full bg-emerald-50 text-emerald-500"
			>
				<Sparkles class="h-5 w-5" />
			</div>
			<h3 class="text-sm font-black tracking-tight text-slate-800">Beribadah</h3>
		</div>
		{#if isRamadan && ramadanDay > 0}
			<span class="rounded-full bg-emerald-100 px-3.5 py-1 text-[10px] font-black text-emerald-800 uppercase tracking-wider">
				Puasa Hari ke-{ramadanDay}
			</span>
		{/if}
	</div>

	<!-- Information Banner -->
	<div class="flex gap-3 rounded-2xl border border-blue-100/80 bg-blue-50/50 p-4 text-left">
		<Info class="mt-0.5 h-4 w-4 shrink-0 text-blue-500" />
		<p class="text-[10px] leading-normal font-bold text-blue-700/90">
			Tanda * wajib diisi oleh siswa muslim. Wanita muslim haid tetap dihitung melaksanakan (akan otomatis mendapatkan nilai penuh untuk periode haid saat rekap bulanan).
		</p>
	</div>

	<!-- Checkboxes Grid -->
	<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
		<!-- Doa -->
		<label
			class="flex cursor-pointer items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50/50 p-4 transition-colors select-none hover:bg-slate-50"
		>
			<input
				type="checkbox"
				bind:checked={doa_ortu}
				class="h-4.5 w-4.5 rounded border-slate-200 text-[#4db6ac] focus:ring-[#4db6ac]"
			/>
			<span class="text-xs leading-tight font-bold text-slate-600"
				>Berdoa untuk diri sendiri & orang tua</span
			>
		</label>

		{#if isIslam}
			<!-- Sholat Fajar -->
			<label
				class="flex cursor-pointer items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50/50 p-4 transition-colors select-none hover:bg-slate-50"
			>
				<input
					type="checkbox"
					bind:checked={sholat_fajar}
					class="h-4.5 w-4.5 rounded border-slate-200 text-[#4db6ac] focus:ring-[#4db6ac]"
				/>
				<span class="text-xs leading-tight font-bold text-slate-600">Sholat Fajar / Qoblal Subuh *</span>
			</label>

			<!-- Sholat 5 Waktu -->
			<label
				class="flex cursor-pointer items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50/50 p-4 transition-colors select-none hover:bg-slate-50"
			>
				<input
					type="checkbox"
					bind:checked={sholat_5waktu}
					class="h-4.5 w-4.5 rounded border-slate-200 text-[#4db6ac] focus:ring-[#4db6ac]"
				/>
				<span class="text-xs leading-tight font-bold text-slate-600">Sholat 5 Waktu Berjamaah *</span>
			</label>

			<!-- Zikir -->
			<label
				class="flex cursor-pointer items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50/50 p-4 transition-colors select-none hover:bg-slate-50"
			>
				<input
					type="checkbox"
					bind:checked={zikir}
					class="h-4.5 w-4.5 rounded border-slate-200 text-[#4db6ac] focus:ring-[#4db6ac]"
				/>
				<span class="text-xs leading-tight font-bold text-slate-600">Zikir & Doa sehabis Sholat *</span>
			</label>

			<!-- Dhuha -->
			<label
				class="flex cursor-pointer items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50/50 p-4 transition-colors select-none hover:bg-slate-50"
			>
				<input
					type="checkbox"
					bind:checked={dhuha}
					class="h-4.5 w-4.5 rounded border-slate-200 text-[#4db6ac] focus:ring-[#4db6ac]"
				/>
				<span class="text-xs leading-tight font-bold text-slate-600">Sholat Dhuha *</span>
			</label>

			{#if isRamadan}
				<!-- Rowatib -->
				<label
					class="flex cursor-pointer items-center gap-3 rounded-2xl border border-slate-100 bg-slate-50/50 p-4 transition-colors select-none hover:bg-slate-50"
				>
					<input
						type="checkbox"
						bind:checked={rowatib}
						class="h-4.5 w-4.5 rounded border-slate-200 text-[#4db6ac] focus:ring-[#4db6ac]"
					/>
					<span class="text-xs leading-tight font-bold text-slate-600">Sholat Sunah Rowatib *</span>
				</label>
			{/if}
		{/if}
	</div>

	<!-- Infaq / Sedekah -->
	<div
		class="flex flex-col justify-between gap-6 border-t border-slate-50 pt-4 md:flex-row md:items-end"
	>
		<div class="flex-1">
			<label class="mb-1.5 block text-[10px] font-bold tracking-wider text-slate-400 uppercase"
				>Infaq / Sedekah</label
			>
			<p class="mb-2 text-[10px] font-bold text-slate-400">
				Masukkan nominal rupiah jika Anda bersedekah hari ini
			</p>
			<div class="relative max-w-sm">
				<span class="absolute top-1/2 left-3.5 -translate-y-1/2 text-xs font-bold text-slate-400"
					>Rp</span
				>
				<input
					type="number"
					placeholder="0"
					bind:value={infaqInputVal}
					class="w-full rounded-xl border border-slate-100 bg-slate-50/50 py-2 pr-4 pl-10 text-xs font-bold text-slate-600 transition-all focus:border-[#4db6ac] focus:outline-none"
				/>
			</div>
		</div>

		<!-- Save Button -->
		<button
			type="button"
			onclick={onsave}
			disabled={loading}
			class="shadow-xxs inline-flex shrink-0 cursor-pointer items-center justify-center gap-1.5 rounded-xl bg-[#4db6ac] px-8 py-2.5 text-xs font-black text-white transition-all hover:bg-[#3ca59b] active:scale-[0.98] disabled:cursor-not-allowed disabled:bg-slate-200"
		>
			{#if loading}
				<Loader class="h-3.5 w-3.5 animate-spin" />
				Menyimpan...
			{:else}
				Simpan Data Ibadah
			{/if}
		</button>
	</div>
</div>
