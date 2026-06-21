<script lang="ts">
	import { onMount } from 'svelte';
	import JurnalForm from '../../../features/student/components/kegiatan/KegiatanForm.svelte';
	import { listJurnalSiswa } from '../../../features/student/logic/kegiatanLogic';
	import { getRekapBulananSiswa } from '../../../features/student/logic/kehadiranLogic';
	import type { RekapBulanan } from '../../../features/student/logic/kehadiranLogic';
	import type { G7Jurnal } from '../../../features/student/types/student.types';
	import { RefreshCw, Check, X, ChevronLeft, ChevronRight } from 'lucide-svelte';

	// Helpers bulan

	function toBulanStr(d: Date): string {
		return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
	}

	function labelBulan(s: string): string {
		const [y, m] = s.split('-').map(Number);
		return new Date(y, m - 1, 1).toLocaleDateString('id-ID', { month: 'long', year: 'numeric' });
	}

	function geserBulan(s: string, delta: number): string {
		const [y, m] = s.split('-').map(Number);
		const d = new Date(y, m - 1 + delta, 1);
		return toBulanStr(d);
	}

	// State

	const bulanSekarang = toBulanStr(new Date());

	// Kehadiran — sumber: rekap_bulanan, navigasi antar bulan
	let bulanKehadiran = $state(bulanSekarang);
	let rekap = $state<RekapBulanan | null>(null);
	let loadingKehadiran = $state(false);

	// Jurnal G7 — sumber: kebiasaan_hebat
	let listJurnals = $state<G7Jurnal[]>([]);
	let loadingJurnal = $state(false);

	// Tombol next hanya aktif jika bulan yang ditampilkan bukan bulan sekarang
	let bisaMaju = $derived(bulanKehadiran < bulanSekarang);

	// Load

	// Fetch rekap kehadiran
	async function loadRekap() {
		loadingKehadiran = true;
		rekap = await getRekapBulananSiswa(bulanKehadiran);
		loadingKehadiran = false;
	}

	// Fetch riwayat jurnal
	async function loadJurnal() {
		loadingJurnal = true;
		const res = await listJurnalSiswa(bulanSekarang, bulanSekarang, 1, 31);
		listJurnals = [...res.items].sort((a, b) => b.tanggal.localeCompare(a.tanggal));
		loadingJurnal = false;
	}

	onMount(() => {
		loadRekap();
		loadJurnal();
	});

	// Navigasi bulan — masing-masing trigger fetch ulang
	function prevBulan() {
		bulanKehadiran = geserBulan(bulanKehadiran, -1);
		loadRekap();
	}
	function nextBulan() {
		if (bisaMaju) {
			bulanKehadiran = geserBulan(bulanKehadiran, 1);
			loadRekap();
		}
	}
</script>

<div class="space-y-6">
	<div class="text-left">
		<h2 class="font-display text-xl font-black tracking-tight text-slate-800">Kegiatan Harian</h2>
		<p class="mt-0.5 text-xs font-bold text-slate-500">
			Catat aktivitas 7 kebiasaan baikmu setiap hari.
		</p>
	</div>

	<JurnalForm />

	<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
		<div class="card p-6">
			<div class="mb-4 flex items-center justify-between">
				<div>
					<h3 class="text-sm font-bold text-foreground">Riwayat Jurnal G7</h3>
					<p class="text-xxs mt-0.5 text-muted">{labelBulan(bulanSekarang)}</p>
				</div>
				<button
					onclick={loadJurnal}
					disabled={loadingJurnal}
					class="inline-flex cursor-pointer items-center gap-1.5 border-none bg-transparent text-xs font-bold text-primary hover:underline"
				>
					<RefreshCw class="h-3.5 w-3.5 {loadingJurnal ? 'animate-spin' : ''}" />
					Segarkan
				</button>
			</div>

			{#if loadingJurnal}
				<p class="py-8 text-center text-xs text-muted">Memuat jurnal...</p>
			{:else if listJurnals.length === 0}
				<p class="py-8 text-center text-xs text-muted">Belum ada jurnal bulan ini.</p>
			{:else}
				<div class="overflow-x-auto">
					<table class="w-full border-collapse text-left text-xs">
						<thead>
							<tr
								class="border-b border-border bg-gray-50 text-[10px] font-bold text-muted uppercase"
							>
								<th class="p-2.5">Tanggal</th>
								<th class="p-2.5 text-center">Bangun</th>
								<th class="p-2.5 text-center">Ibadah</th>
								<th class="p-2.5 text-center">Makan</th>
								<th class="p-2.5 text-center">Gerak</th>
								<th class="p-2.5 text-center">Belajar</th>
								<th class="p-2.5 text-center">Sosial</th>
								<th class="p-2.5 text-center">Tidur</th>
								<th class="p-2.5 text-center">Skor</th>
							</tr>
						</thead>
						<tbody>
							{#each listJurnals as item}
								<tr class="border-b border-border hover:bg-gray-50/50">
									<td class="p-2.5 font-semibold whitespace-nowrap text-foreground"
										>{item.tanggal}</td
									>
									{#each ['bangun', 'ibadah', 'makan', 'olahraga', 'belajar', 'bermasyarakat', 'tidur'] as col}
										<td class="p-2.5 text-center">
											{#if (item as any)[col]?.done}
												<Check class="mx-auto h-3.5 w-3.5 text-emerald-500" />
											{:else}
												<X class="mx-auto h-3.5 w-3.5 text-rose-300" />
											{/if}
										</td>
									{/each}
									<td class="p-2.5 text-center font-bold text-primary">{item.totalDone}/7</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>

		<!-- Kehadiran Bulanan (langsung dari rekap_bulanan) -->
		<div class="flex flex-col gap-4 card p-6">
			<div class="flex items-center justify-between">
				<div>
					<h3 class="text-sm font-bold text-foreground">Kehadiran Bulanan</h3>
				</div>
				<button
					onclick={loadRekap}
					disabled={loadingKehadiran}
					class="inline-flex cursor-pointer items-center gap-1 border-none bg-transparent text-xs font-bold text-primary hover:underline"
				>
					<RefreshCw class="h-3 w-3 {loadingKehadiran ? 'animate-spin' : ''}" />
					Segarkan
				</button>
			</div>

			<!-- Navigasi bulan: ‹ Juni 2026 › -->
			<div
				class="flex items-center justify-between rounded-xl border border-slate-100 bg-slate-50/50 px-3 py-2"
			>
				<button
					onclick={prevBulan}
					class="flex h-7 w-7 cursor-pointer items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-100 active:scale-95"
					aria-label="Bulan sebelumnya"
				>
					<ChevronLeft class="h-4 w-4" />
				</button>
				<span class="text-sm font-extrabold text-slate-700">{labelBulan(bulanKehadiran)}</span>
				<button
					onclick={nextBulan}
					disabled={!bisaMaju}
					class="flex h-7 w-7 cursor-pointer items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-500 transition hover:bg-slate-100 active:scale-95 disabled:cursor-default disabled:opacity-30"
					aria-label="Bulan berikutnya"
				>
					<ChevronRight class="h-4 w-4" />
				</button>
			</div>

			{#if loadingKehadiran}
				<p class="py-8 text-center text-xs text-muted">Memuat rekap dari database...</p>
			{:else if !rekap || rekap.totalHariEfektif === 0}
				<p class="py-8 text-center text-xs text-muted">Belum ada data kehadiran bulan ini.</p>
			{:else}
				<!-- 4 tile: Hadir / Izin / Sakit / Alpa -->
				<div class="grid grid-cols-4 gap-2">
					<div class="rounded-xl border border-emerald-100 bg-emerald-50 p-3 text-center">
						<p class="text-[10px] font-bold text-emerald-600 uppercase">Hadir</p>
						<p class="text-2xl font-black text-emerald-700">{rekap.totalHadir}</p>
					</div>
					<div class="rounded-xl border border-sky-100 bg-sky-50 p-3 text-center">
						<p class="text-[10px] font-bold text-sky-600 uppercase">Izin</p>
						<p class="text-2xl font-black text-sky-700">{rekap.totalIzin}</p>
					</div>
					<div class="rounded-xl border border-amber-100 bg-amber-50 p-3 text-center">
						<p class="text-[10px] font-bold text-amber-600 uppercase">Sakit</p>
						<p class="text-2xl font-black text-amber-700">{rekap.totalSakit}</p>
					</div>
					<div class="rounded-xl border border-rose-100 bg-rose-50 p-3 text-center">
						<p class="text-[10px] font-bold text-rose-600 uppercase">Alpa</p>
						<p class="text-2xl font-black text-rose-700">{rekap.totalTidakHadir}</p>
					</div>
				</div>

				<!-- Magang/PKL -->
				{#if rekap.totalMagang > 0}
					<div
						class="flex items-center justify-between rounded-xl border border-violet-100 bg-violet-50 px-4 py-2"
					>
						<span class="text-xs font-bold text-violet-700">Magang / PKL</span>
						<span class="text-lg font-black text-violet-700">{rekap.totalMagang} hari</span>
					</div>
				{/if}

				<!-- Bar persentase kehadiran -->
				<div class="space-y-1.5">
					<div class="flex items-center justify-between text-xs font-bold">
						<span class="text-slate-600">Persentase Kehadiran</span>
						<span class="text-emerald-600">{rekap.persentaseHadir.toFixed(1)}%</span>
					</div>
					<div class="h-2.5 w-full overflow-hidden rounded-full bg-slate-100">
						<div
							class="h-full rounded-full bg-emerald-500 transition-all duration-500"
							style="width: {Math.min(rekap.persentaseHadir, 100)}%"
						></div>
					</div>
					<p class="text-right text-[10px] font-medium text-slate-400">
						{rekap.totalHariEfektif} hari efektif tercatat
					</p>
				</div>
			{/if}
		</div>
	</div>
</div>
