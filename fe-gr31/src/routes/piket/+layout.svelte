<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import ProtectedRoute from '../../features/shared/components/ProtectedRoute.svelte';
	import PiketSidebar from '../../features/shared/components/sidebar/PiketSidebar.svelte';
	import { sidebarCollapsed } from '../../stores/uiStore';
	import { currentUser } from '../../stores/authStore';
	import { Menu } from 'lucide-svelte';

	let { children } = $props();

	let name = $derived($currentUser?.nama || 'Guru Piket');
	let roleLabel = $derived('Guru Piket');
	let initial = $derived(name.charAt(0).toUpperCase());

	onMount(() => {
		// Izinkan akses jika role "piket" (dari JWT) ATAU email mengandung "piket@"
		const role = $currentUser?.role?.toLowerCase() || '';
		const email = $currentUser?.email?.toLowerCase() || '';
		const isPiket = role === 'piket' || email.includes('piket@');
		if (!isPiket) {
			goto('/admin', { replaceState: true });
		}
	});
</script>

<ProtectedRoute allowedRoles={['admin']}>
	<div class="flex min-h-screen bg-slate-50/50 font-sans">
		<!-- Sidebar -->
		<PiketSidebar />

		<!-- Main Content Area -->
		<div class="flex h-screen flex-1 flex-col overflow-hidden">
			<!-- Top header bar (navbar) -->
			<header
				class="z-10 flex h-16 shrink-0 items-center justify-between border-b border-slate-100 bg-white px-4 sm:px-8 select-none"
			>
				<div class="flex items-center gap-4">
					<button
						aria-label="Menu"
						onclick={() => sidebarCollapsed.update((v) => !v)}
						class="cursor-pointer rounded-xl border border-slate-100 bg-white p-2 text-slate-400 transition-all hover:border-slate-200 hover:bg-slate-50 hover:text-slate-600"
					>
						<Menu class="h-5 w-5" />
					</button>
				</div>

				<!-- Right side: Profile info -->
				<div class="flex items-center gap-3">
					<div class="hidden text-right sm:block">
						<h4 class="text-xs leading-none font-black text-slate-700">{name}</h4>
						<span
							class="mt-1 block text-[9px] font-extrabold tracking-wider text-slate-400 uppercase"
							>{roleLabel}</span
						>
					</div>

					<!-- Avatar icon -->
					<div
						class="shadow-xxs flex h-9 w-9 items-center justify-center rounded-2xl border border-blue-100 bg-blue-50 text-xs font-black text-[#0070f3]"
					>
						{initial}
					</div>
				</div>
			</header>

			<!-- Page Content Viewport -->
			<main class="flex-1 overflow-y-auto bg-slate-50/20 p-4 sm:p-8">
				{@render children()}
			</main>
		</div>
	</div>
</ProtectedRoute>
