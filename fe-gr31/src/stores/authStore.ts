import { writable, derived } from 'svelte/store';

// Helper to access localStorage safely in SvelteKit
const isBrowser = typeof window !== 'undefined';

// Role helpers
export const ADMIN_ROLES = ['admin', 'super_admin', 'guru_bk', 'guru_wali', 'piket', 'admin_bk', 'walas'];

export function isAdminRole(role?: string | null): boolean {
	return !!role && ADMIN_ROLES.includes(role);
}

export function isStudentRole(role?: string | null): boolean {
	// Role Student
	return role === 'student' || role === 'siswa';
}

// Decode payload JWT
function decodeJwt(token: string | null): any {
	if (!token) return null;
	try {
		const part = token.split('.')[1];
		if (!part) return null;
		let base64 = part.replace(/-/g, '+').replace(/_/g, '/');
		const pad = base64.length % 4;
		if (pad) base64 += '='.repeat(4 - pad);
		const json = decodeURIComponent(
			atob(base64)
				.split('')
				.map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
				.join('')
		);
		return JSON.parse(json);
	} catch {
		return null;
	}
}

export function roleFromToken(token: string | null): string | null {
	const payload = decodeJwt(token);
	return payload?.role ?? null;
}

// State
const initialToken = isBrowser
	? localStorage.getItem('adminToken') || localStorage.getItem('studentToken') || null
	: null;

const storedUser = isBrowser ? JSON.parse(localStorage.getItem('currentUser') || 'null') : null;
const initialUser = (() => {
	console.log(
		'[authStore Debug] initializing initialUser. storedUser:',
		storedUser,
		'initialToken:',
		initialToken ? initialToken.substring(0, 15) + '...' : null
	);
	if (!storedUser) return null;
	const tokenRole = roleFromToken(initialToken);

	console.log('[authStore Debug] tokenRole:', tokenRole);
	const res = tokenRole ? { ...storedUser, role: tokenRole } : { ...storedUser, role: storedUser.role || 'admin' };
	console.log('[authStore Debug] final initialUser:', res);
	return res;
})();

export const currentUser = writable<any>(initialUser);
export const authToken = writable<string | null>(initialToken);
export const isLoading = writable<boolean>(false);

export const isAuthenticated = derived(authToken, ($authToken) => $authToken !== null);

export function setAuth(user: any, token: string) {
	console.log(
		'[authStore Debug] setAuth called with user:',
		user,
		'token:',
		token ? token.substring(0, 15) + '...' : null
	);
	const tokenRole = roleFromToken(token);

	const normalizedUser = { ...user, role: tokenRole ?? user?.role ?? 'admin' };
	console.log('[authStore Debug] normalizedUser:', normalizedUser);

	currentUser.set(normalizedUser);
	authToken.set(token);

	if (isBrowser) {
		localStorage.setItem('currentUser', JSON.stringify(normalizedUser));
		if (isStudentRole(normalizedUser.role)) {
			localStorage.setItem('studentToken', token);
			localStorage.removeItem('adminToken');
		} else {
			localStorage.setItem('adminToken', token);
			localStorage.removeItem('studentToken');
		}
	}
}

export function clearAuth() {
	currentUser.set(null);
	authToken.set(null);

	if (isBrowser) {
		localStorage.removeItem('currentUser');
		localStorage.removeItem('studentToken');
		localStorage.removeItem('adminToken');
	}
}
