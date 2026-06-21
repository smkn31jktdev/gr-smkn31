<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import {
		currentUser,
		isAuthenticated,
		clearAuth,
		isAdminRole,
		isStudentRole
	} from '../stores/authStore';

	onMount(() => {
		if ($isAuthenticated && $currentUser) {
			const role = $currentUser.role;
			if (isStudentRole(role)) {
				goto('/siswa', { replaceState: true });
			} else if (role === 'guru_bk' || role === 'admin_bk' || role === 'bk') {
				goto('/bk', { replaceState: true });
			} else if (role === 'piket' || role === 'admin_piket') {
				goto('/piket', { replaceState: true });
			} else if (isAdminRole(role)) {
				goto('/admin', { replaceState: true });
			} else {
				// Token ada tapi role tidak dikenal (sesi lama / data rusak):
				clearAuth();
			}
		}
	});

	const habits = [
		{ name: 'Bangun Pagi', img: '/assets/img/bangun.png' },
		{ name: 'Beribadah', img: '/assets/img/beribadah.png' },
		{ name: 'Berolahraga', img: '/assets/img/olahraga.png' },
		{ name: 'Makan Sehat', img: '/assets/img/makan.png' },
		{ name: 'Gemar Belajar', img: '/assets/img/belajar.png' },
		{ name: 'Bermasyarakat', img: '/assets/img/organisasi.png' },
		{ name: 'Tidur Cepat', img: '/assets/img/tidur.png' }
	];
</script>

<div
	class="relative flex min-h-screen flex-col justify-between overflow-hidden bg-white select-none"
>
	<!-- Dynamic soft lighting glow overlays matching the reference look -->
	<div
		class="pointer-events-none absolute top-0 left-0 h-[500px] w-[500px] -translate-x-1/3 -translate-y-1/3 rounded-full bg-[#29b6f6]/10 blur-[100px]"
	></div>
	<div
		class="pointer-events-none absolute right-0 bottom-0 h-[500px] w-[500px] translate-x-1/3 translate-y-1/3 rounded-full bg-[#0070f3]/5 blur-[100px]"
	></div>

	<!-- Header -->
	<header class="z-10 mx-auto flex w-full max-w-7xl items-center justify-between px-6 py-5">
		<!-- Navbar Logo -->
		<div class="flex items-center">
			<img
				src="/assets/img/navbar.png"
				alt="Gerakan Ramah Anak Logo"
				class="h-10 w-auto object-contain"
			/>
		</div>

		<!-- Login Button Pill -->
		<a
			href="/login"
			class="flex cursor-pointer items-center gap-1.5 rounded-full bg-[#0070f3] px-7 py-2.5 text-xs font-bold text-white shadow-md transition-all duration-200 hover:bg-blue-700 hover:shadow-lg sm:text-sm"
		>
			Login <span class="text-xs font-bold">→</span>
		</a>
	</header>

	<!-- Main Hero Layout -->
	<main class="z-10 flex flex-1 items-center">
		<div
			class="mx-auto grid w-full max-w-7xl grid-cols-1 items-center gap-12 px-6 py-8 lg:grid-cols-12 lg:gap-8"
		>
			<!-- Left Side: Circular Wheel (G7 Habits) -->
			<div class="order-2 flex justify-center lg:order-1 lg:col-span-6">
				<div
					class="relative flex h-80 w-80 items-center justify-center [--r:110px] sm:h-[380px] sm:w-[380px] sm:[--r:140px] lg:h-[440px] lg:w-[440px] lg:[--r:165px]"
				>
					<!-- Thin circular grey connector -->
					<div
						class="pointer-events-none absolute h-[calc(var(--r)*2)] w-[calc(var(--r)*2)] rounded-full border border-gray-100"
					></div>

					<!-- Central Circular Brand Logo -->
					<div
						class="absolute z-10 flex h-32 w-32 items-center justify-center overflow-hidden rounded-full border border-gray-100/80 bg-white p-2.5 shadow-[0_10px_40px_rgba(0,0,0,0.06)] sm:h-40 sm:w-40 lg:h-48 lg:w-48"
					>
						<img
							src="/assets/img/7kaih.png"
							alt="7 KAIH Logo"
							class="h-full w-full object-contain"
						/>
					</div>

					<!-- Circular Orbiting Habit Images -->
					{#each habits as habit, idx}
						{@const angle = (idx * 360) / habits.length - 90}
						{@const rad = (angle * Math.PI) / 180}
						{@const cos = Math.cos(rad)}
						{@const sin = Math.sin(rad)}

						<div
							class="group absolute flex h-12 w-12 cursor-pointer items-center justify-center rounded-full border border-gray-100 bg-white p-2.5 shadow-[0_8px_30px_rgba(0,0,0,0.05)] transition-all duration-300 hover:scale-110 hover:border-[#29b6f6]/40 hover:shadow-[0_12px_35px_rgba(0,0,0,0.1)] sm:h-16 sm:w-16 sm:p-3 lg:h-20 lg:w-20 lg:p-4"
							style="left: calc(50% + var(--r) * {cos}); top: calc(50% + var(--r) * {sin}); transform: translate(-50%, -50%);"
						>
							<img src={habit.img} alt={habit.name} class="h-full w-full object-contain" />
							<!-- Tooltip label -->
							<span
								class="pointer-events-none absolute -bottom-8 z-20 rounded bg-gray-900/90 px-2 py-0.5 text-[9px] font-bold whitespace-nowrap text-white opacity-0 shadow-sm backdrop-blur-xs transition-opacity duration-200 group-hover:opacity-100"
							>
								{habit.name}
							</span>
						</div>
					{/each}
				</div>
			</div>

			<!-- Right Side: Hero Titles and Badges -->
			<div
				class="order-1 flex flex-col items-center space-y-6 text-center lg:order-2 lg:col-span-6 lg:items-start lg:text-left"
			>
				<!-- Subtle Pill Badge -->
				<span
					class="inline-flex items-center rounded-full bg-[#e3f2fd] px-4 py-1.5 text-[10px] font-black tracking-widest text-[#0070f3] uppercase"
				>
					Gerakan Ramah Anak
				</span>

				<!-- High-Impact Bold Header in Unbounded Font -->
				<h1
					class="font-display text-4xl leading-none font-black tracking-tight text-[#29b6f6] uppercase sm:text-[46px] md:text-[52px]"
				>
					SMK Negeri 31 <br />
					<span class="text-[#4db6ac]">Jakarta</span>
				</h1>

				<!-- Subheading Description -->
				<p class="max-w-md text-sm leading-relaxed font-medium text-gray-500 sm:text-base">
					Berfokus pada kegiatan siswa hebat melalui pembiasaan positif setiap hari.
				</p>
			</div>
		</div>
	</main>

	<!-- Footer -->
	<footer class="z-10 w-full border-t border-gray-50 bg-white/50 py-6 text-center backdrop-blur-xs">
		<p class="text-[9px] font-black tracking-widest text-gray-400 uppercase">
			Hak Cipta © 2026 SMKN 31 Jakarta. Semua Hak Dilindungi.
		</p>
	</footer>
</div>
