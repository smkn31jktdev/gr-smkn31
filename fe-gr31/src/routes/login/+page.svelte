<script lang="ts">
	import { onMount } from 'svelte';
	import {
		currentUser,
		isAuthenticated,
		clearAuth,
		isAdminRole,
		isStudentRole
	} from '../../stores/authStore';

	// Redirect to respective dashboards if already authenticated
	onMount(() => {
		if ($isAuthenticated && $currentUser) {
			const role = $currentUser.role;
			if (isStudentRole(role)) {
				window.location.href = '/siswa';
			} else if (role === 'guru_bk' || role === 'admin_bk' || role === 'bk') {
				window.location.href = '/bk';
			} else if (role === 'piket' || role === 'admin_piket') {
				window.location.href = '/piket';
			} else if (isAdminRole(role)) {
				window.location.href = '/admin';
			} else {
				// Role tidak dikenal → bersihkan sesi rusak agar bisa login ulang.
				clearAuth();
			}
		}
	});
</script>

<div
	class="relative flex min-h-screen flex-col items-center justify-between overflow-hidden bg-[#fcfdfe] py-10 select-none"
>
	<div
		class="pointer-events-none absolute inset-0 bg-[radial-gradient(#e5e7eb_1.5px,transparent_1.5px)] bg-size-[16px_16px] opacity-40"
	></div>

	<div
		class="pointer-events-none absolute top-0 left-0 h-[600px] w-[600px] -translate-x-1/3 -translate-y-1/3 rounded-full bg-[#29b6f6]/8 blur-[100px]"
	></div>
	<div
		class="pointer-events-none absolute right-0 bottom-0 h-[600px] w-[600px] translate-x-1/3 translate-y-1/3 rounded-full bg-[#0070f3]/4 blur-[120px]"
	></div>

	<header
		class="z-10 flex w-full max-w-5xl justify-center px-6"
	>
		<div class="flex items-center gap-1.5 opacity-80">
			<img
				src="/assets/img/navbar.png"
				alt="SMK Negeri 31 Jakarta Logo"
				class="h-8 w-auto object-contain"
			/>
		</div>
	</header>

	<main class="z-10 my-6 flex w-full max-w-4xl flex-col items-center px-6 text-center">
		<h2
			class="mb-2.5 font-display text-2xl leading-none font-black tracking-tight text-gray-800 uppercase sm:text-[28px]"
		>
			PILIH PERAN LOGIN
		</h2>
		<p
			class="mb-10 max-w-md text-xs leading-relaxed font-medium text-gray-400 sm:text-sm"
		>
			Selamat datang di portal pembiasaan karakter 7 KAIH. Silakan pilih peranan akses masuk Anda di
			bawah ini.
		</p>

		<div class="grid w-full max-w-2xl grid-cols-1 gap-8 px-2 sm:grid-cols-2">
			<a
				href="/login/admin"
				class="group flex min-h-[350px] flex-col overflow-hidden rounded-xl border border-gray-100/80 bg-white shadow-[0_8px_30px_rgba(0,0,0,0.03)] transition-all duration-500 hover:-translate-y-2 hover:shadow-[0_20px_50px_rgba(0,112,243,0.15)]"
			>
				<div
					class="relative flex h-32 items-center justify-center overflow-hidden bg-linear-to-br from-[#0070f3] to-[#29b6f6]"
				>
					<div
						class="absolute inset-0 bg-[linear-gradient(to_right,#ffffff08_1px,transparent_1px),linear-gradient(to_bottom,#ffffff08_1px,transparent_1px)] bg-size-[10px_10px]"
					></div>

					<div
						class="flex h-16 w-16 items-center justify-center rounded-2xl border border-white/20 bg-white/12 text-white shadow-inner backdrop-blur-md transition-transform duration-300 group-hover:scale-110"
					>
						<!-- User + Gear SVG -->
						<svg
							class="h-8 w-8 text-white"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
							<circle cx="9" cy="7" r="4" />
							<circle cx="18" cy="18" r="3" />
							<path d="M18 14v1M18 21v1M14 18h1M21 18h1" />
						</svg>
					</div>
				</div>

				<div class="flex flex-1 flex-col justify-between p-6 text-left">
					<div class="space-y-2">
						<h3 class="font-sans text-sm font-black tracking-tight text-gray-800 uppercase">
							Akses Guru & Staf
						</h3>
						<p class="text-[11px] leading-relaxed font-medium text-gray-400">
							Verifikasi jurnal harian siswa, kelola data kehadiran GPS, and lakukan rekapitulasi
							penilaian G7 bulanan.
						</p>
					</div>
					<div class="mt-4 flex items-center justify-between border-t border-gray-50 pt-4">
						<span class="text-[10px] font-bold tracking-wider text-[#0070f3] uppercase"
							>Masuk Sekarang</span
						>
						<span
							class="flex h-7 w-7 items-center justify-center rounded-full bg-blue-50 text-xs text-[#0070f3] transition-transform duration-200 group-hover:translate-x-1.5"
							>→</span
						>
					</div>
				</div>
			</a>

			<a
				href="/login/siswa"
				class="group flex min-h-[350px] flex-col overflow-hidden rounded-xl border border-gray-100/80 bg-white shadow-[0_8px_30px_rgba(0,0,0,0.03)] transition-all duration-500 hover:-translate-y-2 hover:shadow-[0_20px_50px_rgba(77,182,172,0.15)]"
			>
				<div
					class="relative flex h-32 items-center justify-center overflow-hidden bg-linear-to-br from-[#4db6ac] to-[#80cbc4]"
				>
					<div
						class="absolute inset-0 bg-[linear-gradient(to_right,#ffffff08_1px,transparent_1px),linear-gradient(to_bottom,#ffffff08_1px,transparent_1px)] bg-size-[10px_10px]"
					></div>

					<div
						class="flex h-16 w-16 items-center justify-center rounded-2xl border border-white/20 bg-white/12 text-white shadow-inner backdrop-blur-md transition-transform duration-300 group-hover:scale-110"
					>
						<svg
							class="h-8 w-8 text-white"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							stroke-width="2"
							stroke-linecap="round"
							stroke-linejoin="round"
						>
							<path d="M22 10v6M2 10l10-5 10 5-10 5z" />
							<path d="M6 12v5c0 2 2 3 6 3s6-1 6-3v-5" />
						</svg>
					</div>
				</div>

				<div class="flex flex-1 flex-col justify-between p-6 text-left">
					<div class="space-y-2">
						<h3 class="font-sans text-sm font-black tracking-tight text-gray-800 uppercase">
							Akses Siswa
						</h3>
						<p class="text-[11px] leading-relaxed font-medium text-gray-400">
							Laporkan jurnal 7 pembiasaan harian KAIH, lakukan absensi masuk GPS mandiri, and
							unggah lampiran bukti foto/video.
						</p>
					</div>
					<div class="mt-4 flex items-center justify-between border-t border-gray-50 pt-4">
						<span class="text-[10px] font-bold tracking-wider text-[#4db6ac] uppercase"
							>Masuk Sekarang</span
						>
						<span
							class="flex h-7 w-7 items-center justify-center rounded-full bg-teal-50 text-xs text-[#4db6ac] transition-transform duration-200 group-hover:translate-x-1.5"
							>→</span
						>
					</div>
				</div>
			</a>
		</div>

		<!-- Clean return link -->
		<a
			href="/"
			class="text-xxs mt-12 flex cursor-pointer items-center gap-1.5 rounded-xl border border-gray-200 px-5 py-2.5 font-bold text-gray-500 transition-all duration-200 hover:border-gray-300 hover:bg-gray-50"
		>
			← Kembali ke Beranda
		</a>
	</main>

	<!-- Footer Branding -->
	<footer
		class="z-10 text-center"
	>
		<span class="text-[9px] font-bold tracking-wider text-gray-400 uppercase">
			Gerakan Ramah Anak SMK Negeri 31 Jakarta
		</span>
	</footer>
</div>
