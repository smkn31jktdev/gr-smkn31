<script lang="ts">
	import { User, CalendarDays, AlertCircle, Sparkles, ChartColumn, Info, Briefcase } from 'lucide-svelte';

	// Svelte 5 Props destructuring
	let { reportType, classSummary } = $props<{
		reportType: 'bulanan' | 'mingguan';
		classSummary: {
			totalHadir: number;
			totalIzinSakit: number;
			totalIzin: number;
			totalSakit: number;
			totalAlpa: number;
			totalMagang: number;
			activeStudents: number;
			hariEfektif: number;
			rate: number;
		};
	}>();

	let isAttention = $derived(classSummary.rate < 75);
</script>

<!-- Summary Cards for Class Report -->
<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
	<!-- Card 1: Total Siswa Aktif -->
	<div
		class="flex items-center gap-3 rounded-2xl border border-slate-100/80 bg-white p-4 shadow-xs"
	>
		<div
			class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-slate-100 bg-slate-50 text-slate-500"
		>
			<User class="h-4.5 w-4.5" />
		</div>
		<div class="text-left">
			<span class="text-slate-450 block text-[9px] font-bold tracking-wider uppercase"
				>Siswa Aktif</span
			>
			<span class="mt-0.5 block text-sm font-bold text-slate-800">
				{classSummary.activeStudents}
				<span class="text-[10px] font-normal text-slate-400">Siswa</span>
			</span>
		</div>
	</div>

	<!-- Card 2: Hari Efektif (bulanan dari backend) / Izin+Sakit (mingguan) -->
	<div
		class="flex items-center gap-3 rounded-2xl border border-slate-100/80 bg-white p-4 shadow-xs"
	>
		<div
			class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border border-slate-100 bg-slate-50 text-slate-500"
		>
			<CalendarDays class="h-4.5 w-4.5" />
		</div>
		<div class="text-left">
			{#if reportType === 'bulanan'}
				<span class="text-slate-450 block text-[9px] font-bold tracking-wider uppercase"
					>Hari Kerja Efektif</span
				>
				<span class="mt-0.5 block text-sm font-bold text-slate-800">
					{classSummary.hariEfektif}
					<span class="text-[10px] font-normal text-slate-400">Hari</span>
				</span>
			{:else}
				<span class="text-slate-450 block text-[9px] font-bold tracking-wider uppercase"
					>Izin & Sakit (Mingguan)</span
				>
				<span class="mt-0.5 block text-sm font-bold text-slate-800">
					{classSummary.totalIzinSakit}
					<span class="text-[10px] font-normal text-slate-400">Kali</span>
				</span>
			{/if}
		</div>
	</div>

	<!-- Card 4: Persentase Kehadiran (dari data real) -->
	<div
		class="flex items-center gap-3 rounded-2xl border p-4 text-left shadow-xs transition-colors {isAttention
			? 'border-rose-100 bg-rose-50/35'
			: 'border-slate-100 bg-white'}"
	>
		<div
			class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border transition-colors {isAttention
				? 'border-rose-200/20 bg-rose-100/50 text-rose-600'
				: 'text-slate-550 border-slate-100 bg-slate-50'}"
		>
			<ChartColumn class="h-4.5 w-4.5" />
		</div>
		<div>
			<span
				class="block text-[9px] font-bold tracking-wider uppercase transition-colors {isAttention
					? 'text-rose-500'
					: 'text-slate-450'}">Kehadiran Rata-rata</span
			>
			<span
				class="mt-0.5 block text-sm font-extrabold {isAttention
					? 'text-rose-700'
					: 'text-slate-800'}"
			>
				{classSummary.rate}%
				<span
					class="block text-[9px] font-normal {isAttention ? 'text-rose-500' : 'text-slate-400'}"
				>
					{isAttention ? 'Perlu Perhatian' : 'Sangat Baik'}
				</span>
			</span>
		</div>
	</div>
</div>

<!-- Info Alert Banner -->
<div
	class="shadow-xxs flex gap-3 rounded-2xl border border-[#00a294]/15 bg-[#00a294]/5 p-4 text-left"
>
	<Info class="mt-0.5 h-4 w-4 shrink-0 text-[#00a294]" />
	<div class="space-y-1">
		<h4 class="text-xs font-bold text-slate-800">Metodologi Perhitungan Laporan</h4>
		<p class="text-[11px] leading-relaxed font-medium text-slate-500">
			{#if reportType === 'bulanan'}
				Persentase Kehadiran Kelas dihitung: <code
					>(Total Hadir / (Siswa × {classSummary.hariEfektif} Hari Kerja)) × 100</code
				>. Hari Kerja Efektif diambil dari data rekap bulanan backend. Siswa berstatus Magang (PKL)
				tidak dihitung dalam pembagi.
			{:else}
				Persentase Kehadiran Mingguan dihitung berdasarkan 5 hari sekolah efektif (Senin s.d.
				Jumat): <code>(Total Hadir / (Siswa × 5 Hari Kerja)) × 100</code>.
			{/if}
		</p>
	</div>
</div>
