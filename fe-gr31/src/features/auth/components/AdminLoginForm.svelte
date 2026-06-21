<script lang="ts">
	import { onMount } from 'svelte';
	import { loginAdmin } from '../logic/authLogic';
	import { roleFromToken } from '../../../stores/authStore';
	import SubmitButton from '../../shared/components/SubmitButton.svelte';

	let email = $state('');
	let password = $state('');
	let rememberMe = $state(false);
	let errorMsg = $state('');
	let showPassword = $state(false);
	let submitBtnEl = $state<HTMLButtonElement>();

	onMount(() => {
		const remember = localStorage.getItem('admin_remember') === 'true';
		if (remember) {
			email = localStorage.getItem('admin_email') || '';
			password = localStorage.getItem('admin_password') || '';
			rememberMe = true;
		}
	});

	async function handleSubmit(handlers: { resolve: () => void; reject: () => void }) {
		errorMsg = '';
		const success = await loginAdmin(email, password);
		if (success) {
			if (rememberMe) {
				localStorage.setItem('admin_remember', 'true');
				localStorage.setItem('admin_email', email);
				localStorage.setItem('admin_password', password);
			} else {
				localStorage.removeItem('admin_remember');
				localStorage.removeItem('admin_email');
				localStorage.removeItem('admin_password');
			}
			handlers.resolve();
			// Baca role langsung dari token yang baru disimpan (hindari race condition store)
			const freshToken =
				localStorage.getItem('adminToken') ?? localStorage.getItem('studentToken');
			const role = roleFromToken(freshToken);
			if (role === 'guru_bk' || role === 'admin_bk' || role === 'bk') {
				window.location.href = '/bk';
			} else if (role === 'piket' || role === 'admin_piket') {
				window.location.href = '/piket';
			} else {
				window.location.href = '/admin';
			}
		} else {
			errorMsg = 'Kombinasi email dan kata sandi salah. Silakan coba lagi.';
			handlers.reject();
		}
	}
</script>

<form
	class="w-full space-y-6"
	onsubmit={(e) => e.preventDefault()}
	onkeydown={(e) => {
		if (e.key === 'Enter' && (e.target as HTMLElement).tagName !== 'BUTTON') {
			e.preventDefault();
			submitBtnEl?.click();
		}
	}}
>
	{#if errorMsg}
		<div
			class="rounded-2xl border border-rose-100 bg-rose-50 p-4 text-xs leading-relaxed font-bold text-rose-600"
		>
			{errorMsg}
		</div>
	{/if}

	<!-- Email Field -->
	<div class="space-y-1.5 text-left">
		<label
			for="email"
			class="block text-[11px] font-extrabold tracking-wider text-gray-500 uppercase"
		>
			Email Guru / Admin
		</label>
		<div class="relative">
			<!-- Mail Icon -->
			<svg
				class="absolute top-1/2 left-4 h-5 w-5 -translate-y-1/2 text-gray-400"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
				<polyline points="22,6 12,13 2,6" />
			</svg>
			<input
				type="email"
				id="email"
				placeholder="Masukkan email Anda"
				bind:value={email}
				class="w-full rounded-2xl border border-slate-200/80 bg-[#f8fafc] py-3.5 pr-4 pl-12 text-sm text-gray-800 placeholder-gray-400 transition-all duration-200 outline-none focus:border-[#0070f3] focus:bg-white focus:ring-4 focus:ring-blue-100/40"
				required
			/>
		</div>
	</div>

	<!-- Password Field -->
	<div class="space-y-1.5 text-left">
		<label
			for="password"
			class="block text-[11px] font-extrabold tracking-wider text-gray-500 uppercase"
		>
			Kata Sandi
		</label>
		<div class="relative">
			<!-- Padlock Icon -->
			<svg
				class="absolute top-1/2 left-4 h-5 w-5 -translate-y-1/2 text-gray-400"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
				<path d="M7 11V7a5 5 0 0 1 10 0v4" />
			</svg>
			<input
				type={showPassword ? 'text' : 'password'}
				id="password"
				placeholder="Masukkan kata sandi"
				bind:value={password}
				class="w-full rounded-2xl border border-slate-200/80 bg-[#f8fafc] py-3.5 pr-12 pl-12 text-sm text-gray-800 placeholder-gray-400 transition-all duration-200 outline-none focus:border-[#0070f3] focus:bg-white focus:ring-4 focus:ring-blue-100/40"
				required
			/>
			<!-- Eye Toggle Icon -->
			<button
				type="button"
				onclick={() => (showPassword = !showPassword)}
				class="absolute top-1/2 right-4 -translate-y-1/2 text-gray-400 transition-colors hover:text-gray-600 focus:outline-none"
				aria-label={showPassword ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi'}
			>
				{#if showPassword}
					<svg
						class="h-5 w-5"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path
							d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"
						/>
						<line x1="1" y1="1" x2="23" y2="23" />
					</svg>
				{:else}
					<svg
						class="h-5 w-5"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
						<circle cx="12" cy="12" r="3" />
					</svg>
				{/if}
			</button>
		</div>
	</div>

	<!-- Remember Me Checkbox -->
	<div class="flex items-center text-left">
		<label class="group flex cursor-pointer items-center gap-2.5 select-none">
			<input
				type="checkbox"
				bind:checked={rememberMe}
				class="h-4.5 w-4.5 cursor-pointer rounded border-slate-300 text-[#0070f3] transition-colors focus:ring-[#0070f3]"
			/>
			<span
				class="text-xs font-semibold text-slate-500 transition-colors group-hover:text-slate-700"
				>Ingat saya</span
			>
		</label>
	</div>

	<!-- Submit Button -->
	<div class="pt-2">
		<SubmitButton
			bind:el={submitBtnEl}
			label="Masuk"
			loadingLabel="Memproses Masuk..."
			className="w-full py-3.5 bg-linear-to-r! from-[#0070f3]! to-[#29b6f6]! hover:from-[#0060cb]! hover:to-[#02a5f4]! rounded-2xl! text-white font-bold text-sm shadow-[0_4px_15px_rgba(0,112,243,0.25)] hover:shadow-[0_6px_20px_rgba(0,112,243,0.35)] transition-all active:scale-[0.99] cursor-pointer border-none"
			onclick={handleSubmit}
		/>
	</div>
</form>
