<script lang="ts">
	import { onMount } from 'svelte';
	import { currentUser, isStudentRole } from '../../stores/authStore';
	import { fade, scale } from 'svelte/transition';
	import { getKehadiranBulanan } from '../../features/student/logic/kehadiranLogic';
	import type { KehadiranHariItem } from '../../features/student/logic/kehadiranLogic';
	import { getDashboardG7 } from '../../features/student/logic/kegiatanLogic';
	import type { G7DashboardSiswa } from '../../features/student/logic/kegiatanLogic';
	import StudentProfileCard from '../../features/student/components/StudentProfileCard.svelte';
	import DailyProgressBar from '../../features/student/components/DailyProgressBar.svelte';
	import QuickAttendance from '../../features/student/components/absensi/QuickAttendance.svelte';
	import DailyHabitsList from '../../features/student/components/kegiatan/DailyHabitsList.svelte';
	import MonthlyAttendanceCalendar from '../../features/student/components/absensi/MonthlyAttendanceCalendar.svelte';
	import { HelpCircle } from 'lucide-svelte';

	let kehadiranBulanan = $state<KehadiranHariItem[]>([]);
	let dashboardG7 = $state<G7DashboardSiswa | null>(null);
	let loadingHistory = $state(false);
	let loadingJurnal = $state(false);

	const now = new Date();
	const currentYear = now.getFullYear();
	const currentMonth = now.getMonth();
	const todayDateStr = now.toLocaleDateString('sv-SE');

	const todayStrIndonesian = now.toLocaleDateString('id-ID', {
		weekday: 'long',
		day: 'numeric',
		month: 'long',
		year: 'numeric'
	});

	const monthNameIndonesian = now.toLocaleDateString('id-ID', {
		month: 'long',
		year: 'numeric'
	});

	async function loadDashboardData() {
		loadingHistory = true;
		loadingJurnal = true;

		const bulanStr = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}`;

		const [kehadiranRes, g7Res] = await Promise.all([
			getKehadiranBulanan(bulanStr),
			getDashboardG7()
		]);

		kehadiranBulanan = kehadiranRes?.kehadiran ?? [];
		dashboardG7 = g7Res;

		loadingHistory = false;
		loadingJurnal = false;
	}

	// Promotion congratulations popup states
	let showPromotionPopup = $state(false);
	let studentName = $derived($currentUser?.nama || '');
	let studentClass = $derived($currentUser?.kelas || '');
	let displayGrade = $derived.by(() => {
		if (!studentClass) return '';
		return studentClass.trim().split(/\s+/)[0];
	});

	function closePromotionPopup() {
		showPromotionPopup = false;
		if ($currentUser && studentClass) {
			const nisn = $currentUser.nis || $currentUser.nisn || $currentUser.id || 'user';
			const storageKey = `g7_promotion_popup_shown_${nisn}_${studentClass}`;
			localStorage.setItem(storageKey, 'true');
		}
	}

	onMount(() => {
		loadDashboardData();

		// Check promotion popup trigger
		if ($currentUser && isStudentRole($currentUser.role)) {
			const now = new Date();
			// July 1st or later (month >= 6: July is 6)
			const isJulyOrLater = now.getMonth() >= 6;

			if (isJulyOrLater && studentClass) {
				const nisn = $currentUser.nis || $currentUser.nisn || $currentUser.id || 'user';
				const storageKey = `g7_promotion_popup_shown_${nisn}_${studentClass}`;
				const alreadyShown = localStorage.getItem(storageKey);
				if (!alreadyShown) {
					showPromotionPopup = true;
				}
			}
		}
	});

	let todayKehadiran = $derived(kehadiranBulanan.find((k) => k.tanggal === todayDateStr));
	let completedHabitsCount = $derived(dashboardG7?.progresHariIni ?? 0);
	let todayJurnal = $derived(dashboardG7?.jurnalHariIni ?? null);

	// Buat grid kalender bulanan dari data kehadiran
	let calendarDays = $derived.by(() => {
		const days = [];
		const firstDayIndex = new Date(currentYear, currentMonth, 1).getDay();
		const daysInMonth = new Date(currentYear, currentMonth + 1, 0).getDate();

		for (let i = 0; i < firstDayIndex; i++) {
			days.push({ day: null, dateStr: '', status: 'empty', waktu: '' });
		}

		const startOfToday = new Date();
		startOfToday.setHours(0, 0, 0, 0);

		for (let d = 1; d <= daysInMonth; d++) {
			const dayDate = new Date(currentYear, currentMonth, d);
			const dateStr = dayDate.toLocaleDateString('sv-SE');
			const attendance = kehadiranBulanan.find((k) => k.tanggal === dateStr);

			let status = 'belum_absen';
			let waktu = '';

			if (dayDate > startOfToday) {
				status = 'future';
			} else if (attendance) {
				status = attendance.status;
				waktu = attendance.waktuAbsen;
			}

			days.push({ day: d, dateStr, status, waktu });
		}

		return days;
	});
</script>

<div class="relative space-y-6 pb-10 font-sans select-none">
	<!-- Dashboard Page Header Section -->
	<div
		class="flex flex-col justify-between gap-4 border-b border-slate-100 pb-5 md:flex-row md:items-center"
	>
		<div>
			<h1 class="mb-2 font-sans text-3xl leading-none font-extrabold tracking-tight text-slate-800">
				Dashboard
			</h1>
			<div class="flex items-center gap-1.5 text-xs font-semibold text-slate-400">
				<!-- Calendar icon -->
				<svg
					class="h-4 w-4 text-slate-400"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
					stroke-width="2"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
					/>
				</svg>
				<span>{todayStrIndonesian}</span>
			</div>
		</div>

		<!-- Active activity greeting pill -->
		<div
			class="shadow-xxs inline-flex items-center gap-2 rounded-2xl border border-amber-100 bg-amber-50/50 px-4.5 py-2 text-xs font-bold text-amber-700 md:self-center"
		>
			<!-- Clock icon -->
			<svg
				class="h-4.5 w-4.5 animate-spin text-amber-500"
				style="animation-duration: 20s;"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
				stroke-width="2"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364-6.364l-.707.707M6.343 17.657l-.707.707m12.728 0l-.707-.707M6.343 6.343l-.707-.707M12 8v4l3 3"
				/>
			</svg>
			<span>Selamat Beraktivitas!</span>
		</div>
	</div>

	<!-- Cards Grid Layout -->
	<div class="grid grid-cols-1 gap-6 md:grid-cols-12">
		<!-- Profile Info Card -->
		<StudentProfileCard />

		<!-- Daily Progress Bar Card -->
		<DailyProgressBar {completedHabitsCount} />
	</div>

	<!-- Bottom Details and Calendars Layout -->
	<div class="grid grid-cols-1 gap-6 lg:grid-cols-12">
		<!-- Habit Activities List (Col span 7) -->
		<div
			class="space-y-6 rounded-3xl border border-slate-100/90 bg-white p-6 shadow-[0_10px_35px_rgba(0,0,0,0.01)] lg:col-span-7"
		>
			<!-- Section Header -->
			<div class="flex items-center justify-between border-b border-slate-50 pb-4">
				<div class="text-left">
					<h3 class="text-sm font-black tracking-tight text-slate-800">Aktivitas Hari Ini</h3>
					<p class="mt-0.5 text-[10px] font-bold text-slate-400">
						Kelola absensi harian dan pantau progress 7 kebiasaan baik Anda.
					</p>
				</div>
			</div>
			<div class="mb-6">
				<QuickAttendance {todayKehadiran} onload={loadDashboardData} />
			</div>

			<DailyHabitsList {loadingJurnal} {todayJurnal} />
		</div>

		<!-- Monthly Attendance Calendar -->
		<MonthlyAttendanceCalendar
			{calendarDays}
			{loadingHistory}
			{monthNameIndonesian}
			onrefresh={loadDashboardData}
		/>
	</div>

	<!-- Floating Counseling BK Help button -->
	<a
		href="/siswa/chat"
		class="fixed right-6 bottom-6 z-20 flex h-12 w-12 cursor-pointer items-center justify-center rounded-full border border-[#4db6ac]/20 bg-[#4db6ac] text-white shadow-lg transition-all duration-200 hover:scale-105 hover:bg-[#3ca59b] hover:shadow-xl active:scale-95"
		aria-label="Konsultasi BK"
	>
		<HelpCircle class="h-6 w-6 text-white" />
	</a>

	<!-- Celebratory Promotion Congratulations Popup -->
	{#if showPromotionPopup}
		<div class="fixed inset-0 z-50 flex items-center justify-center p-4">
			<!-- Backdrop -->
			<div
				transition:fade={{ duration: 250 }}
				onclick={closePromotionPopup}
				class="fixed inset-0 cursor-pointer bg-slate-900/60 backdrop-blur-xs"
				role="presentation"
			></div>

			<!-- Celebratory Card -->
			<div
				transition:scale={{ duration: 300, start: 0.95 }}
				class="animate-in fade-in zoom-in-95 relative z-10 flex w-full max-w-md flex-col items-center overflow-hidden rounded-3xl border border-slate-100 bg-white p-8 text-center shadow-2xl duration-300 select-none"
			>
				<!-- Subtle decorative glow layers -->
				<div
					class="absolute -top-10 -right-10 h-32 w-32 rounded-full bg-[#4db6ac]/10 blur-xl"
				></div>
				<div
					class="absolute -bottom-10 -left-10 h-32 w-32 rounded-full bg-amber-500/10 blur-xl"
				></div>

				<!-- Congratulations Icon Badge -->
				<div
					class="relative mb-6 flex h-20 w-20 items-center justify-center rounded-3xl border border-amber-100 bg-amber-50 text-amber-500 shadow-[0_8px_20px_rgba(245,158,11,0.15)] transition-transform duration-300 hover:scale-105"
				>
					<svg
						class="h-10 w-10"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
						stroke-width="2"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"
						/>
					</svg>
					<!-- Sparkle dots -->
					<span class="absolute -top-1 -right-1 flex h-4 w-4">
						<span
							class="absolute inline-flex h-full w-full animate-ping rounded-full bg-amber-400 opacity-75"
						></span>
						<span class="relative inline-flex h-4 w-4 rounded-full bg-amber-500"></span>
					</span>
				</div>

				<!-- Greetings & Title -->
				<h3 class="mb-1 text-xl font-black tracking-tight text-slate-800">
					Halo {studentName}!
				</h3>
				<span class="mb-4 block text-xs font-black tracking-wider text-[#4db6ac] uppercase">
					Selamat Kenaikan Kelas!
				</span>

				<!-- Motivational Message Content -->
				<div
					class="mb-8 space-y-3.5 text-left text-xs leading-relaxed font-semibold text-slate-500"
				>
					<p>
						Selamat atas kerja kerasmu! Sekarang kamu telah resmi naik ke kelas <span
							class="rounded-lg border border-[#4db6ac]/20 bg-[#e0f2f1]/40 px-2 py-0.5 font-bold text-[#00a294]"
							>{displayGrade}</span
						>. Ini adalah awal dari babak baru yang penuh tantangan dan peluang emas.
					</p>
					<p>
						Mari kita sambut jenjang kelas <span class="font-bold text-slate-800"
							>{displayGrade}</span
						>
						ini dengan energi positif! Tetap konsisten menjalankan dan mengisi jurnal
						<strong class="text-slate-700">7 Pembiasaan Anak Indonesia Hebat (G7KAIH)</strong> setiap
						hari. Karakter unggul dan kesuksesan masa depanmu dibangun dari komitmen serta kebiasaan baik
						kecil yang kamu lakukan hari ini.
					</p>
				</div>

				<!-- Action button -->
				<button
					onclick={closePromotionPopup}
					class="w-full transform cursor-pointer rounded-2xl border-none bg-[#4db6ac] px-6 py-3.5 text-xs font-bold text-white shadow-md transition-all duration-200 hover:bg-[#3ca59b] hover:shadow-lg active:scale-98"
				>
					Siap, Mulai Hari Ini!
				</button>
			</div>
		</div>
	{/if}
</div>
