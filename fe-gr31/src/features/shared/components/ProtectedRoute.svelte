<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import {
		authToken,
		currentUser,
		clearAuth,
		isAdminRole,
		isStudentRole
	} from '../../../stores/authStore';

	let { allowedRoles = [], children } = $props<{
		allowedRoles: string[];
		children?: import('svelte').Snippet;
	}>();

	let isAuthorized = $state(false);

	onMount(() => {
		const token = $authToken;
		const user = $currentUser;

		if (!token || !user) {
			goto('/login', { replaceState: true });
			return;
		}
		const currentRole = user.role;

		let hasAccess = false;
		for (const allowed of allowedRoles) {
			if (allowed === 'siswa' && isStudentRole(currentRole)) {
				hasAccess = true;
			} else if (allowed === 'admin' && isAdminRole(currentRole)) {
				hasAccess = true;
			} else if (
				allowed === 'piket' &&
				(currentRole === 'admin' || currentRole === 'super_admin' || currentRole === 'piket' || currentRole === 'admin_piket')
			) {
				hasAccess = true;
			} else if (
				allowed === 'guru_bk' &&
				(currentRole === 'guru_bk' || currentRole === 'admin_bk' || currentRole === 'bk')
			) {
				hasAccess = true;
			} else if (allowed === currentRole) {
				hasAccess = true;
			} else if (currentRole === 'super_admin' && allowed !== 'siswa') {
				hasAccess = true;
			}
		}

		if (hasAccess) {
			isAuthorized = true;
			return;
		}

		if (isStudentRole(currentRole)) {
			goto('/siswa', { replaceState: true });
		} else if (currentRole === 'guru_bk' || currentRole === 'admin_bk' || currentRole === 'bk') {
			goto('/bk', { replaceState: true });
		} else if (currentRole === 'piket' || currentRole === 'admin_piket') {
			goto('/piket', { replaceState: true });
		} else if (isAdminRole(currentRole)) {
			goto('/admin', { replaceState: true });
		} else {
			clearAuth();
			goto('/login', { replaceState: true });
		}
	});
</script>

{#if isAuthorized && children}
	{@render children()}
{/if}
