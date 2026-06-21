<script lang="ts">
	import { currentUser } from '../../../stores/authStore';

	let name = $derived($currentUser?.nama || 'Siswa');
	let kelas = $derived($currentUser?.kelas || 'Kelas');
	let nis = $derived($currentUser?.nis || '-');
	let initial = $derived(name.charAt(0).toUpperCase());
</script>

<div
	class="flex flex-col justify-between overflow-hidden rounded-3xl border border-slate-100/90 bg-white shadow-[0_10px_35px_rgba(0,0,0,0.01)] md:col-span-5"
>
	<!-- Top header area with light teal tint -->
	<div
		class="relative flex items-center justify-center border-b border-slate-50 bg-[#e0f2f1]/30 py-8"
	>
		<div
			class="relative z-10 flex h-16 w-16 items-center justify-center rounded-full border-2 border-white bg-[#4db6ac] text-xl font-black text-white shadow-xs"
		>
			{initial}
		</div>
	</div>
	<!-- Bottom details -->
	<div class="flex flex-1 flex-col justify-center p-6 text-center">
		<h3 class="mb-1 text-lg font-black tracking-tight text-slate-800">{name}</h3>
		<p
			class="inline-flex items-center justify-center gap-1.5 text-xs leading-none font-bold text-slate-400"
		>
			<!-- User Profile Icon -->
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
					d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<span>{kelas} • {nis}</span>
		</p>
	</div>
</div>
