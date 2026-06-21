<script lang="ts">
	import { onMount } from 'svelte';
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
			window.location.href = '/login';
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
			window.location.href = '/siswa';
		} else if (currentRole === 'guru_bk' || currentRole === 'admin_bk' || currentRole === 'bk') {
			window.location.href = '/bk';
		} else if (currentRole === 'piket' || currentRole === 'admin_piket') {
			window.location.href = '/piket';
		} else if (isAdminRole(currentRole)) {
			window.location.href = '/admin';
		} else {
			clearAuth();
			window.location.href = '/login';
		}
	});
</script>

{#if isAuthorized && children}
	{@render children()}
{/if}
