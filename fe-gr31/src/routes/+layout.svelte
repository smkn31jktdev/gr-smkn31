<script lang="ts">
	import './layout.css';
	import Toast from '../features/shared/components/Toast.svelte';
	import { onMount } from 'svelte';
	import { version } from '$app/environment';
	import { updated } from '$app/state';
	import { addToast } from '../stores/uiStore';

	let { children } = $props();

	async function clearCachesAndReload() {
		// 1. Clear Cache Storage
		if ('caches' in window) {
			try {
				const keys = await caches.keys();
				await Promise.all(keys.map((key) => caches.delete(key)));
				console.log('Cache storage cleared successfully');
			} catch (e) {
				console.error('Failed to clear Cache Storage:', e);
			}
		}

		// 2. Unregister all service workers
		if ('serviceWorker' in navigator) {
			try {
				const registrations = await navigator.serviceWorker.getRegistrations();
				await Promise.all(registrations.map((r) => r.unregister()));
				console.log('Service workers unregistered successfully');
			} catch (e) {
				console.error('Failed to unregister Service Workers:', e);
			}
		}

		// 3. Clear session storage
		try {
			sessionStorage.clear();
			console.log('Session storage cleared successfully');
		} catch (e) {
			console.error('Failed to clear session storage:', e);
		}

		// 4. Force a hard reload from the server to bypass browser cache
		console.log('Forcing page reload...');
		window.location.reload();
	}

	onMount(async () => {
		// Only run cache checking logic in production, or when version is actually defined
		if (!version || version === 'development') {
			return;
		}

		// 1. Check for immediate update on startup
		try {
			const updateAvailable = await updated.check();
			if (updateAvailable) {
				console.log('New update detected on startup!');
				addToast('Memuat pembaruan sistem terbaru...', 'info');
				// Give the toast a brief moment to show before reloading
				setTimeout(clearCachesAndReload, 1000);
				return;
			}
		} catch (e) {
			console.warn('Failed to check for updates on startup:', e);
		}

		// 2. Fallback check using localStorage comparison
		const lastVersion = localStorage.getItem('last_app_version');
		if (lastVersion !== version) {
			console.log(`App version changed from ${lastVersion} to ${version}`);
			localStorage.setItem('last_app_version', version);
			addToast('Pembaruan sistem berhasil diterapkan!', 'success');
			// Clear caches to ensure we start fresh on the new version
			setTimeout(clearCachesAndReload, 1000);
		}
	});

	// Dynamic check for updates while the app is running
	$effect(() => {
		if (updated.current) {
			console.log('New update detected during app runtime!');
			addToast('Versi baru aplikasi tersedia. Memperbarui...', 'info');
			setTimeout(clearCachesAndReload, 2000);
		}
	});
</script>

<svelte:head>
	<title>GR31 SMKN 31 Jakarta</title>
</svelte:head>

<Toast />

{@render children()}
