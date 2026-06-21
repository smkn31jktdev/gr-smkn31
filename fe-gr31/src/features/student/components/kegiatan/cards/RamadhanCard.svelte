<script lang="ts">
	import { Moon, Sparkles, Loader } from 'lucide-svelte';

	let {
		puasa = $bindable(false),
		tarawihRokaat = $bindable(0),
		onsave,
		loading = false
	} = $props();

	const setRokaat = (val: number) => {
		tarawihRokaat = val;
	};
</script>

<div
	class="space-y-6 rounded-3xl border border-amber-100 bg-amber-50/20 p-6 shadow-[0_8px_30px_rgb(245,158,11,0.03)] relative overflow-hidden"
>
	<!-- Decorative background glows -->
	<div class="absolute -right-24 -top-24 h-48 w-48 rounded-full bg-amber-200/20 blur-3xl pointer-events-none"></div>

	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-3">
			<div
				class="flex h-9 w-9 items-center justify-center rounded-full bg-amber-100 text-amber-600 shadow-xs"
			>
				<Moon class="h-5 w-5 fill-amber-500/20" />
			</div>
			<div>
				<h3 class="text-sm font-black tracking-tight text-slate-800">Ibadah Tambahan Ramadhan</h3>
				<p class="text-[10px] text-amber-600 font-bold mt-0.5">Khusus di Bulan Suci Ramadhan</p>
			</div>
		</div>
	</div>

	<!-- Grid -->
	<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
		<!-- Puasa Ramadhan -->
		<div class="space-y-3">
			<label class="block text-[10px] font-black tracking-wider text-slate-400 uppercase">Puasa Hari Ini</label>
			<label
				class="flex cursor-pointer items-center gap-3 rounded-2xl border transition-all p-5 select-none bg-white shadow-xs {puasa ? 'border-amber-400 bg-amber-50/10' : 'border-slate-100 hover:bg-slate-50/30'}"
			>
				<input
					type="checkbox"
					bind:checked={puasa}
					class="h-5 w-5 rounded border-slate-200 text-amber-500 focus:ring-amber-400 cursor-pointer"
				/>
				<div class="text-left">
					<span class="text-xs leading-tight font-black text-slate-700 block">Saya Berpuasa</span>
					<span class="text-[9px] text-slate-400 font-bold block mt-0.5">Mendapatkan pahala wajib puasa Ramadhan</span>
				</div>
			</label>
		</div>

		<!-- Sholat Tarawih -->
		<div class="space-y-3">
			<label class="block text-[10px] font-black tracking-wider text-slate-400 uppercase">Tarawih & Witir</label>
			<div class="rounded-2xl border border-slate-100 bg-white p-4 shadow-xs space-y-3">
				<div class="flex items-center justify-between">
					<span class="text-xs font-black text-slate-700">Jumlah Rakaat</span>
					<span class="text-xs font-mono font-black text-amber-600 bg-amber-50 px-2.5 py-1 rounded-lg border border-amber-100">{tarawihRokaat} Rakaat</span>
				</div>

				<div class="grid grid-cols-4 gap-1.5">
					<button
						type="button"
						onclick={() => setRokaat(0)}
						class="py-2 rounded-xl text-[10px] font-black tracking-wide border cursor-pointer transition-all {tarawihRokaat === 0 ? 'bg-slate-800 text-white border-slate-800 shadow-xs' : 'bg-slate-50 hover:bg-slate-100 text-slate-600 border-slate-200/60'}"
					>
						Tidak
					</button>
					<button
						type="button"
						onclick={() => setRokaat(11)}
						class="py-2 rounded-xl text-[10px] font-black tracking-wide border cursor-pointer transition-all {tarawihRokaat === 11 ? 'bg-amber-500 text-white border-amber-500 shadow-xs' : 'bg-slate-50 hover:bg-slate-100 text-slate-600 border-slate-200/60'}"
					>
						11 Rk
					</button>
					<button
						type="button"
						onclick={() => setRokaat(23)}
						class="py-2 rounded-xl text-[10px] font-black tracking-wide border cursor-pointer transition-all {tarawihRokaat === 23 ? 'bg-amber-500 text-white border-amber-500 shadow-xs' : 'bg-slate-50 hover:bg-slate-100 text-slate-600 border-slate-200/60'}"
					>
						23 Rk
					</button>
					<div class="relative">
						<input
							type="number"
							placeholder="Lain"
							bind:value={tarawihRokaat}
							class="w-full text-center py-2 rounded-xl text-[10px] font-black text-slate-600 border border-slate-200/60 focus:border-amber-400 focus:outline-none bg-slate-50/50 [appearance:textfield] [&::-webkit-outer-spin-button]:appearance-none [&::-webkit-inner-spin-button]:appearance-none"
						/>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Save Action -->
	<div class="flex justify-end border-t border-amber-100/30 pt-4 mt-2">
		<button
			type="button"
			onclick={onsave}
			disabled={loading}
			class="shadow-sm inline-flex shrink-0 cursor-pointer items-center justify-center gap-1.5 rounded-xl bg-amber-500 px-8 py-2.5 text-xs font-black text-white transition-all hover:bg-amber-600 active:scale-[0.98] disabled:cursor-not-allowed disabled:bg-slate-200 border-none"
		>
			{#if loading}
				<Loader class="h-3.5 w-3.5 animate-spin" />
				Menyimpan...
			{:else}
				Simpan Data Ramadhan
			{/if}
		</button>
	</div>
</div>
