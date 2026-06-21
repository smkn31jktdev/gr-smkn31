<script lang="ts">
	import { getGeolocation } from '../../logic/kehadiranLogic';
	import { addToast } from '../../../../stores/uiStore';
	import { MapPin, Navigation, RefreshCw, Loader2 } from 'lucide-svelte';

	let {
		distanceVerified = $bindable(false),
		currentCoords = $bindable(undefined),
		currentAccuracy = $bindable(undefined),
		distanceFromSchool = $bindable(undefined),
		isWithinRange = $bindable(false)
	} = $props<{
		distanceVerified?: boolean;
		currentCoords?: { lat: number; lng: number } | undefined;
		currentAccuracy?: number | undefined;
		distanceFromSchool?: number | undefined;
		isWithinRange?: boolean;
	}>();

	let locationLoading = $state(false);

	// Formula Haversine di client-side
	function calculateHaversineDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
		const R = 6371e3; // radius bumi dalam meter
		const dLat = ((lat2 - lat1) * Math.PI) / 180;
		const dLon = ((lon2 - lon1) * Math.PI) / 180;
		const a =
			Math.sin(dLat / 2) * Math.sin(dLat / 2) +
			Math.cos((lat1 * Math.PI) / 180) *
				Math.cos((lat2 * Math.PI) / 180) *
				Math.sin(dLon / 2) *
				Math.sin(dLon / 2);
		const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
		return R * c;
	}

	async function handleVerifyDistance() {
		locationLoading = true;
		const loc = await getGeolocation();
		if (loc) {
			currentCoords = loc.koordinat;
			currentAccuracy = loc.akurasi;

			// Ambil koordinat sekolah dari localStorage atau default
			const sekolahLat = Number(localStorage.getItem('config_sekolah_lat') || '-6.1819399');
			const sekolahLng = Number(localStorage.getItem('config_sekolah_lng') || '106.8518572');
			const maxRadius = Number(localStorage.getItem('config_radius_meter') || '80');

			const dist = calculateHaversineDistance(
				currentCoords.lat,
				currentCoords.lng,
				sekolahLat,
				sekolahLng
			);
			distanceFromSchool = Math.round(dist);
			isWithinRange = dist <= maxRadius;
			distanceVerified = true;

			if (isWithinRange) {
				addToast(`Lokasi terverifikasi! Jarak ke sekolah: ${distanceFromSchool}m (Dalam radius ${maxRadius}m)`, 'success');
			} else {
				addToast(`Anda berada di luar radius sekolah (${distanceFromSchool}m). Hubungi guru piket/kelas jika ingin mengajukan izin/magang`, 'warning');
			}
		}
		locationLoading = false;
	}
</script>

<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
	<!-- Bagian Kiri: Verifikasi Jarak -->
	<div class="flex items-center gap-3">
		<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-teal-50 text-teal-500">
			<MapPin class="h-5.5 w-5.5" />
		</div>
		<div class="text-left">
			<h4 class="text-xs font-black text-slate-800">Verifikasi Jarak Lokasi</h4>
			{#if distanceVerified}
				<p class="mt-0.5 text-[10px] font-bold" class:text-emerald-600={isWithinRange} class:text-amber-600={!isWithinRange}>
					{isWithinRange ? 'Terverifikasi' : 'Di luar radius'} ({distanceFromSchool}m dari sekolah)
				</p>
			{:else}
				<p class="mt-0.5 text-[10px] font-bold text-slate-400">
					Wajib verifikasi jarak untuk memulai absensi harian
				</p>
			{/if}
		</div>
	</div>

	<!-- Tombol Verifikasi Jarak -->
	<div>
		{#if distanceVerified}
			<button
				type="button"
				onclick={handleVerifyDistance}
				disabled={locationLoading}
				class="shadow-xxs flex items-center gap-1.5 rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-xs font-bold text-slate-600 transition-all hover:bg-slate-50 active:scale-95 disabled:opacity-50 cursor-pointer"
			>
				{#if locationLoading}
					<Loader2 class="h-3.5 w-3.5 animate-spin text-slate-400" />
					Memverifikasi...
				{:else}
					<RefreshCw class="h-3.5 w-3.5 text-slate-500" />
					Verifikasi Ulang
				{/if}
			</button>
		{:else}
			<button
				type="button"
				onclick={handleVerifyDistance}
				disabled={locationLoading}
				class="shadow-xxs flex items-center gap-1.5 rounded-xl bg-teal-500 px-4.5 py-2 text-xs font-bold text-white transition-all hover:bg-teal-600 active:scale-95 disabled:opacity-50 cursor-pointer"
			>
				{#if locationLoading}
					<Loader2 class="h-3.5 w-3.5 animate-spin text-white" />
					Membaca GPS...
				{:else}
					<Navigation class="h-3.5 w-3.5 text-white" />
					Verifikasi Jarak
				{/if}
			</button>
		{/if}
	</div>
</div>
