export const BASE = import.meta.env.VITE_API_BASE_URL ?? 'https://api.gr31.tech';

export function getUploadUrl(path: string | null | undefined): string {
	if (!path) return '';
	if (path.startsWith('http://') || path.startsWith('https://')) {
		return path;
	}
	const cleanPath = path.startsWith('/') ? path : '/' + path;
	return `${BASE}${cleanPath}`;
}

export interface ApiResponse<T = any> {
	data: T | null;
	error: string | null;
	status: number;
}

// Handle Auth
function handleAuthError(status: number) {
	if (typeof window === 'undefined') return;
	if (status === 401) {
		localStorage.removeItem('adminToken');
		localStorage.removeItem('studentToken');
		localStorage.removeItem('currentUser');
		if (!window.location.pathname.includes('/login')) {
			window.location.href = '/login';
		}
	}

}

export async function apiRequest<T = any>(
	path: string,
	options: RequestInit = {}
): Promise<ApiResponse<T>> {
	const isBrowser = typeof window !== 'undefined';
	const token = isBrowser
		? (localStorage.getItem('adminToken') ?? localStorage.getItem('studentToken'))
		: null;

	const method = (options.method || 'GET').toUpperCase();
	const isWrite = method === 'POST' || method === 'PUT' || method === 'DELETE' || method === 'PATCH';

	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(token ? { Authorization: `Bearer ${token}` } : {}),
		...(isBrowser && isWrite ? { 'X-Idempotency-Key': crypto.randomUUID() } : {}),
		...(options.headers as Record<string, string>)
	};

	try {
		const response = await fetch(`${BASE}${path}`, {
			...options,
			headers
		});

		if (response.status === 204) {
			return { data: null, error: null, status: 204 };
		}

		const json = await response.json().catch(() => ({}));
		if (!response.ok) {
			handleAuthError(response.status);
			const errMsg = json.message ?? json.error ?? response.statusText ?? 'Terjadi kesalahan';
			return { data: null, error: errMsg, status: response.status };
		}

		const resultData = json.data !== undefined ? json.data : json;
		return { data: resultData as T, error: null, status: response.status };
	} catch (error: any) {
		return { data: null, error: error.message ?? 'Tidak dapat terhubung ke server', status: 0 };
	}
}

// Special upload helper for multipart files (permits / photos)
export async function apiUpload<T = any>(
	path: string,
	formData: FormData
): Promise<ApiResponse<T>> {
	const isBrowser = typeof window !== 'undefined';
	const token = isBrowser
		? (localStorage.getItem('adminToken') ?? localStorage.getItem('studentToken'))
		: null;

	const headers: Record<string, string> = {
		...(token ? { Authorization: `Bearer ${token}` } : {})
	};

	try {
		const response = await fetch(`${BASE}${path}`, {
			method: 'POST',
			body: formData,
			headers
		});

		const json = await response.json().catch(() => ({}));
		if (!response.ok) {
			const errMsg = json.message ?? json.error ?? response.statusText ?? 'Gagal upload';
			return { data: null, error: errMsg, status: response.status };
		}

		const resultData = json.data !== undefined ? json.data : json;
		return { data: resultData as T, error: null, status: response.status };
	} catch (error: any) {
		return { data: null, error: error.message ?? 'Gagal upload file ke server', status: 0 };
	}
}
