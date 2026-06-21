<script lang="ts">
	import { submitKehadiran } from '../../logic/kehadiranLogic';
	import type { KehadiranHariItem } from '../../logic/kehadiranLogic';
	import { addToast } from '../../../../stores/uiStore';
	import { CheckCircle2, QrCode, Briefcase, CalendarRange } from 'lucide-svelte';

	// Sub-components
	import AbsenJarak from './AbsenJarak.svelte';
	import AbsenQR from './AbsenQR.svelte';
	import AbsenMagang from './AbsenMagang.svelte';
	import AbsenIzin from './AbsenIzin.svelte';

	let { todayKehadiran, onload } = $props<{
		todayKehadiran: KehadiranHariItem | undefined;
		onload: () => void;
	}>();

	// State absensi & verifikasi jarak (bindable dengan AbsenJarak)
	let distanceVerified = $state(false);
	let currentCoords = $state<{ lat: number; lng: number } | undefined>(undefined);
	let currentAccuracy = $state<number | undefined>(undefined);
	let distanceFromSchool = $state<number | undefined>(undefined);
	let isWithinRange = $state(false);

	// Modals/Absen triggers
	let showAttendanceModal = $state(false);
	let showQRModal = $state(false);
	let showMagangAbsen = $state(false);

	// Weekend check
	const today = new Date();
	const isWeekendDay = today.getDay() === 0 || today.getDay() === 6;

	async function handleDirectAbsen() {
		if (!distanceVerified || !currentCoords) {
			addToast('Harap verifikasi jarak terlebih dahulu!', 'error');
			return;
		}
		if (!isWithinRange) {
			addToast('Tidak dapat melakukan absensi Hadir karena Anda berada di luar radius sekolah.', 'error');
			return;
		}

		const success = await submitKehadiran({
			status: 'hadir',
			koordinat: currentCoords,
			akurasi: currentAccuracy
		});

		if (success) {
			onload();
		}
	}
</script>

<!-- Widget Absensi Mandiri (Kotak Kecil) -->
{#if todayKehadiran}
	<div class="rounded-2xl border {todayKehadiran.status === 'hadir' && todayKehadiran.waktuAbsen && todayKehadiran.waktuAbsen.substring(0, 5) > '06:35' ? 'border-rose-100 bg-rose-50/40' : 'border-emerald-100 bg-emerald-50/40'} p-4 flex items-center justify-between transition-all duration-300">
		<div class="flex items-center gap-3">
			<div class="flex h-10 w-10 items-center justify-center rounded-xl {todayKehadiran.status === 'hadir' && todayKehadiran.waktuAbsen && todayKehadiran.waktuAbsen.substring(0, 5) > '06:35' ? 'bg-rose-100 text-rose-600' : 'bg-emerald-100 text-emerald-600'}">
				<CheckCircle2 class="h-5 w-5" />
			</div>
			<div class="text-left">
				<h4 class="text-xs font-black text-slate-800">Kehadiran Hari Ini Tercatat</h4>
				<p class="mt-0.5 text-[10px] font-bold text-slate-500">
					Status: <span class="capitalize font-extrabold {todayKehadiran.status === 'hadir' && todayKehadiran.waktuAbsen && todayKehadiran.waktuAbsen.substring(0, 5) > '06:35' ? 'text-rose-700' : 'text-emerald-700'}">{todayKehadiran.status === 'hadir' && todayKehadiran.waktuAbsen && todayKehadiran.waktuAbsen.substring(0, 5) > '06:35' ? 'terlambat' : todayKehadiran.status}</span> pada pukul {todayKehadiran.waktuAbsen.substring(0, 5)} WIB
				</p>
			</div>
		</div>
		<div class="text-[9px] font-extrabold tracking-wider {todayKehadiran.status === 'hadir' && todayKehadiran.waktuAbsen && todayKehadiran.waktuAbsen.substring(0, 5) > '06:35' ? 'text-rose-700 bg-rose-100/60' : 'text-emerald-700 bg-emerald-100/60'} px-3 py-1.5 rounded-lg uppercase">
			TERKIRIM
		</div>
	</div>
{:else if isWeekendDay}
	<div class="rounded-2xl border border-amber-100 bg-amber-50/40 p-4 flex items-center justify-between transition-all duration-300">
		<div class="flex items-center gap-3">
			<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-amber-100 text-amber-600">
				<CalendarRange class="h-5 w-5" />
			</div>
			<div class="text-left">
				<h4 class="text-xs font-black text-slate-800">Absensi Hari Ini Ditutup</h4>
				<p class="mt-0.5 text-[10px] font-bold text-slate-550 leading-relaxed max-w-[220px] sm:max-w-none">
					Absensi hanya dibuka pada hari efektif sekolah (Senin s.d. Jumat).
				</p>
			</div>
		</div>
		<div class="text-[9px] font-extrabold tracking-wider text-amber-700 bg-amber-100/60 px-3 py-1.5 rounded-lg uppercase">
			TUTUP
		</div>
	</div>
{:else}
	<div class="rounded-2xl border border-slate-100 bg-slate-50/40 p-4 transition-all duration-300">
		
		<!-- Sub-component 1: Jarak / GPS Check -->
		<AbsenJarak 
			bind:distanceVerified 
			bind:currentCoords 
			bind:currentAccuracy 
			bind:distanceFromSchool 
			bind:isWithinRange 
		/>

		<!-- Bagian Bawah: Action Buttons -->
		<div class="mt-4 border-t border-slate-100/70 pt-4">
			<div class="grid grid-cols-2 gap-3 sm:grid-cols-4">
				<!-- Button Absen Hadir -->
				<button
					type="button"
					onclick={handleDirectAbsen}
					disabled={!distanceVerified || !isWithinRange}
					class="flex flex-col items-center justify-center rounded-xl border p-3 text-center transition-all duration-150
						{!distanceVerified
							? 'border-slate-100 bg-slate-50/50 text-slate-300 cursor-not-allowed'
							: !isWithinRange
							? 'border-slate-100 bg-slate-50/50 text-slate-300 cursor-not-allowed'
							: 'border-teal-100 bg-white text-teal-600 hover:bg-teal-50 active:scale-97 cursor-pointer'}"
				>
					<CheckCircle2 class="mb-1.5 h-5 w-5" />
					<span class="text-[10px] font-extrabold uppercase tracking-wide">Absen Hadir</span>
				</button>

				<!-- Button Izin / Sakit -->
				<button
					type="button"
					onclick={() => showAttendanceModal = true}
					class="flex flex-col items-center justify-center rounded-xl border p-3 text-center transition-all duration-150
						border-sky-100 bg-white text-sky-600 hover:bg-sky-50 active:scale-97 cursor-pointer"
				>
					<CalendarRange class="mb-1.5 h-5 w-5" />
					<span class="text-[10px] font-extrabold uppercase tracking-wide">Izin / Sakit</span>
				</button>

				<!-- Button Magang -->
				<button
					type="button"
					onclick={() => showMagangAbsen = true}
					class="flex flex-col items-center justify-center rounded-xl border p-3 text-center transition-all duration-150
						border-purple-100 bg-white text-purple-600 hover:bg-purple-50 active:scale-97 cursor-pointer"
				>
					<Briefcase class="mb-1.5 h-5 w-5" />
					<span class="text-[10px] font-extrabold uppercase tracking-wide">Magang</span>
				</button>

				<!-- Button Scan QR -->
				<button
					type="button"
					onclick={() => showQRModal = true}
					class="flex flex-col items-center justify-center rounded-xl border p-3 text-center transition-all duration-150
						border-amber-100 bg-white text-amber-600 hover:bg-amber-50 active:scale-97 cursor-pointer"
				>
					<QrCode class="mb-1.5 h-5 w-5" />
					<span class="text-[10px] font-extrabold uppercase tracking-wide">Scan QR</span>
				</button>
			</div>
		</div>
	</div>
{/if}

<!-- Sub-component 2: Absen QR Modal -->
<AbsenQR bind:show={showQRModal} {onload} />

<!-- Sub-component 3: Absen Magang Modal -->
<AbsenMagang bind:show={showMagangAbsen} {onload} />

<!-- Sub-component 4: Absen Izin Modal -->
<AbsenIzin bind:show={showAttendanceModal} {onload} />
