<script lang="ts">
	import { onMount } from 'svelte';
	import { Loader2, ChevronLeft, ChevronRight, Search, Download, RefreshCw } from 'lucide-svelte';
	import { addToast } from '../../../../stores/uiStore';
	import { currentUser } from '../../../../stores/authStore';
	import {
		activeTab,
		reportType,
		loading,
		selectedKelas,
		selectedBulan,
		selectedWeekMonday,
		studentSearchQuery,
		foundStudents,
		selectedStudent,
		items,
		weeklyClassItems,
		summaryByClass,
		summaryRange,
		rawWeeklyLogs,
		studentLogs,
		studentSummary,
		studentMonthlyTrend,
		semesterTrend,
		classTotalHariEfektif,
		availableMonths,
		availableKelas,
		changeWeek,
		searchForStudents,
		loadData,
		loadKelas,
		formatDateStr,
		studentRekapList,
		studentRekapPage,
		studentRekapLimit,
		studentRekapTotal,
		studentRekapHasMore,
		studentRekapLoading,
		changeStudentRekapPage,
		setStudentRekapLimit,
		resetStudentRekapPage
	} from '../../logic/piketLaporanLogic';

	// Sub-components
	import LaporanCards from './LaporanCards.svelte';
	import LaporanCharts from './LaporanCharts.svelte';
	import LaporanTableKelas from './LaporanTableKelas.svelte';
	import LaporanDetailSiswa from './LaporanDetailSiswa.svelte';
	import LaporanLombaKebersihan from './LaporanLombaKebersihan.svelte';
	import DropdownChoice from '../../../shared/components/DropdownChoice.svelte';

	onMount(() => {
		loadKelas();
		loadData();
	});

	// Week range string formatter
	function formatDateRangeStr(monStr: string): string {
		if (!monStr) return '';
		const mon = new Date(monStr);
		const fri = new Date(mon);
		fri.setDate(fri.getDate() + 4);
		const formatShort = (d: Date) => {
			const months = [
				'Jan',
				'Feb',
				'Mar',
				'Apr',
				'Mei',
				'Jun',
				'Jul',
				'Agu',
				'Sep',
				'Okt',
				'Nov',
				'Des'
			];
			return `${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`;
		};
		return `${formatShort(mon)} - ${formatShort(fri)}`;
	}

	function handleFilter() {
		resetStudentRekapPage();
		loadData();
	}

	function selectStudent(student: any) {
		selectedStudent.set(student);
		studentSearchQuery.set('');
		foundStudents.set([]);
		loadData();
	}

	function clearSelectedStudent() {
		selectedStudent.set(null);
		studentLogs.set([]);
		studentSummary.set({ hadir: 0, izin: 0, sakit: 0, alpa: 0, magang: 0 });
	}

	// PDF download simulation
	function handleDownloadPDF() {
		const currentKelasVal = $selectedKelas || 'Semua Kelas';
		const currentType =
			$activeTab === 'kelas' ? `Kelas ${currentKelasVal}` : `Siswa ${$selectedStudent?.nama || ''}`;
		const currentPeriod =
			$reportType === 'bulanan' ? $selectedBulan : formatDateRangeStr($selectedWeekMonday);
		addToast(
			`Mempersiapkan unduhan PDF Laporan untuk ${currentType} (${currentPeriod})...`,
			'info'
		);
		setTimeout(() => {
			try {
				const dummyContent =
					`LAPORAN KEHADIRAN ABSENSI - SMKN 31 JAKARTA\n` +
					`Tipe Laporan: ${$activeTab === 'kelas' ? 'Rekap Kelas' : 'Laporan Siswa'}\n` +
					`Periode: ${currentPeriod}\n` +
					`Subjek: ${currentType}\n` +
					`Tanggal Unduh: ${new Date().toLocaleDateString('id-ID')}\n` +
					`Status Unduhan: SUKSES`;
				const blob = new Blob([dummyContent], { type: 'application/pdf' });
				const url = URL.createObjectURL(blob);
				const link = document.createElement('a');
				link.href = url;
				link.download = `Laporan_${$activeTab}_${$reportType}_${currentPeriod.replace(/\s+/g, '_')}.pdf`;
				link.click();
				URL.revokeObjectURL(url);
				addToast('Laporan PDF berhasil diunduh', 'success');
			} catch {
				addToast('Gagal mengunduh file PDF', 'error');
			}
		}, 1000);
	}

	// Reactive class summary
	let classSummary = $derived.by(() => {
		if ($reportType === 'bulanan' && $summaryRange) {
			const sr = $summaryRange;
			const classes = $summaryByClass;
			const activeClasses = classes.filter((c) => c.tingkatKehadiran > 0);
			let calculatedRate = sr.tingkatKehadiran;
			if (activeClasses.length > 0) {
				const sum = activeClasses.reduce((acc, c) => acc + c.tingkatKehadiran, 0);
				calculatedRate = parseFloat((sum / activeClasses.length).toFixed(2));
			}
			return {
				totalHadir: sr.totalHadir,
				totalIzinSakit: sr.totalIzin + sr.totalSakit,
				totalIzin: sr.totalIzin,
				totalSakit: sr.totalSakit,
				totalAlpa: sr.totalAlpa,
				totalMagang: sr.totalMagang,
				activeStudents: sr.totalSiswa,
				hariEfektif: sr.hariEfektif,
				rate: calculatedRate
			};
		}
		// Mingguan: hitung dari weeklyClassItems
		const list = $weeklyClassItems;
		let totalHadir = 0, totalIzin = 0, totalSakit = 0, totalAlpa = 0, totalMagang = 0;
		list.forEach((item) => {
			totalHadir  += item.totalHadir  || 0;
			totalIzin   += item.totalIzin   || 0;
			totalSakit  += item.totalSakit  || 0;
			totalAlpa   += item.totalAlpa   || 0;
			totalMagang += item.totalMagang || 0;
		});
		const activeStudents = list.length;
		const hariEfektif = $classTotalHariEfektif > 0 ? $classTotalHariEfektif : 5;

		const classes = $summaryByClass;
		const activeClasses = classes.filter((c) => c.tingkatKehadiran > 0);
		let calculatedRate = 0;
		if (activeClasses.length > 0) {
			const sum = activeClasses.reduce((acc, c) => acc + c.tingkatKehadiran, 0);
			calculatedRate = parseFloat((sum / activeClasses.length).toFixed(2));
		} else {
			const possibleAttendances = activeStudents * hariEfektif;
			calculatedRate = possibleAttendances > 0 ? parseFloat(((totalHadir / possibleAttendances) * 100).toFixed(1)) : 0;
		}

		return {
			totalHadir,
			totalIzinSakit: totalIzin + totalSakit,
			totalIzin,
			totalSakit,
			totalAlpa,
			totalMagang,
			activeStudents,
			hariEfektif,
			rate: calculatedRate
		};
	});

	// ── Weekly chart: jumlah hadir per hari Senin–Jumat ──────────────────────────
	let weeklyChartData = $derived.by(() => {
		const days = ['Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat'];
		const counts = [0, 0, 0, 0, 0];
		const mon = new Date($selectedWeekMonday);
		const dates = [0, 1, 2, 3, 4].map((idx) => {
			const d = new Date(mon);
			d.setDate(d.getDate() + idx);
			return formatDateStr(d);
		});
		$rawWeeklyLogs.forEach((log) => {
			const dateIdx = dates.indexOf(log.tanggal);
			if (dateIdx !== -1 && (log.status === 'hadir' || log.status === 'magang')) counts[dateIdx]++;
		});
		return days.map((dayName, idx) => ({ label: dayName, val: counts[idx] }));
	});

	function shortenKelasName(name: string): string {
		if (!name) return '';
		const parts = name.split(/\s+/);
		if (parts.length === 0) return '';
		
		const first = parts[0];
		const isRoman = /^(X|XI|XII|IX|VIII|VII|I|II|III|IV|V|VI)$/i.test(first);
		
		if (isRoman) {
			if (parts.length === 2) {
				const second = parts[1];
				if (second.length > 4) {
					const lowerSec = second.toLowerCase();
					if (lowerSec.startsWith('akuntan')) return `${first} AK`;
					if (lowerSec.startsWith('perbankan') || lowerSec.startsWith('layanan')) return `${first} LP`;
					if (lowerSec.startsWith('desain') || lowerSec.startsWith('dkv')) return `${first} DKV`;
					return `${first} ${second.slice(0, 3).toUpperCase()}`;
				}
				return name;
			} else if (parts.length > 2) {
				const initials = parts.slice(1)
					.map(p => p[0])
					.join('')
					.toUpperCase();
				return `${first} ${initials}`;
			}
		}
		return name;
	}

	// ── Monthly/semester chart: tren persentase kehadiran per kelas ──────────────
	let monthlyChartData = $derived.by(() => {
		const classes = $summaryByClass;
		if (classes.length > 0) {
			return classes.map((item) => ({
				label: shortenKelasName(item.kelas),
				val: parseFloat(item.tingkatKehadiran.toFixed(1))
			}));
		}
		return [];
	});
</script>

<div class="space-y-6 pb-10 font-sans text-slate-700 select-none">
	<!-- Header -->
	<div class="flex items-start justify-between border-b border-slate-100 pb-5">
		<div class="space-y-1 text-left">
			<h1 class="text-xl font-bold tracking-tight text-slate-800">Laporan Kehadiran</h1>
			<p class="text-xs font-medium text-slate-400">
				Analisis kehadiran siswa per bulan, per minggu, dan per individu berdasarkan data absensi.
			</p>
		</div>

		<!-- Period picker (top-right) -->
		<div class="flex items-center gap-3">
			{#if $activeTab !== 'lomba'}
				{#if $reportType === 'bulanan'}
					<div class="relative min-w-[150px] text-left">
						<DropdownChoice
							options={availableMonths}
							bind:value={$selectedBulan}
							onchange={handleFilter}
							placeholder="Pilih Bulan"
						/>
					</div>
				{:else}
					<div
						class="shadow-xxs flex items-center overflow-hidden rounded-xl border border-slate-200 bg-white"
					>
						<button
							onclick={() => changeWeek(-7)}
							class="text-slate-550 cursor-pointer border-none bg-transparent p-2.5 transition-colors hover:bg-slate-50"
						>
							<ChevronLeft class="h-4 w-4" />
						</button>
						<span class="px-4 text-xs font-bold text-slate-700 select-none">
							{formatDateRangeStr($selectedWeekMonday)}
						</span>

						<button
							onclick={() => changeWeek(7)}
							class="text-slate-550 cursor-pointer border-none bg-transparent p-2.5 transition-colors hover:bg-slate-50"
						>
							<ChevronRight class="h-4 w-4" />
						</button>
					</div>
				{/if}
			{/if}
		</div>
	</div>

	<!-- Primary Tabs: Laporan Kelas vs Laporan Per Siswa vs Lomba Kebersihan -->
	<div class="flex gap-1 border-b border-slate-100 pb-px">
		<button
			onclick={() => {
				activeTab.set('kelas');
				loadData();
			}}
			class="cursor-pointer border-b-2 bg-transparent px-5 py-3 text-xs font-bold transition-all"
			class:border-[#00a294]={$activeTab === 'kelas'}
			class:text-[#00a294]={$activeTab === 'kelas'}
			class:border-transparent={$activeTab !== 'kelas'}
			class:text-slate-400={$activeTab !== 'kelas'}
		>
			Laporan Kelas
		</button>
		<button
			onclick={() => {
				activeTab.set('siswa');
				resetStudentRekapPage();
				loadData();
			}}
			class="cursor-pointer border-b-2 bg-transparent px-5 py-3 text-xs font-bold transition-all"
			class:border-[#00a294]={$activeTab === 'siswa'}
			class:text-[#00a294]={$activeTab === 'siswa'}
			class:border-transparent={$activeTab !== 'siswa'}
			class:text-slate-400={$activeTab !== 'siswa'}
		>
			Laporan Per Siswa
		</button>
		{#if $currentUser?.role === 'super_admin'}
		<button
			onclick={() => {
				activeTab.set('lomba');
			}}
			class="cursor-pointer border-b-2 bg-transparent px-5 py-3 text-xs font-bold transition-all"
			class:border-[#00a294]={$activeTab === 'lomba'}
			class:text-[#00a294]={$activeTab === 'lomba'}
			class:border-transparent={$activeTab !== 'lomba'}
			class:text-slate-400={$activeTab !== 'lomba'}
		>
			Lomba Kebersihan
		</button>
		{/if}
	</div>

	<!-- Filter Row -->
	{#if $activeTab !== 'lomba'}
		<div
			class="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-slate-100/80 bg-white p-4 shadow-xs"
		>
		<div class="flex w-full flex-wrap items-center gap-4 sm:w-auto">
			<!-- Bulanan / Mingguan toggle -->
			<div class="flex shrink-0 rounded-full border border-slate-200 bg-slate-50/50 p-1">
				<button
					onclick={() => {
						reportType.set('bulanan');
						resetStudentRekapPage();
						loadData();
					}}
					class="cursor-pointer rounded-full border-none px-4.5 py-1.5 text-xs font-bold transition-all"
					class:bg-[#00a294]={$reportType === 'bulanan'}
					class:text-white={$reportType === 'bulanan'}
					class:shadow-xs={$reportType === 'bulanan'}
					class:bg-transparent={$reportType !== 'bulanan'}
					class:text-slate-450={$reportType !== 'bulanan'}
					class:hover:text-slate-700={$reportType !== 'bulanan'}
				>
					Laporan Bulanan
				</button>
				<button
					onclick={() => {
						reportType.set('mingguan');
						resetStudentRekapPage();
						loadData();
					}}
					class="cursor-pointer rounded-full border-none px-4.5 py-1.5 text-xs font-bold transition-all"
					class:bg-[#00a294]={$reportType === 'mingguan'}
					class:text-white={$reportType === 'mingguan'}
					class:shadow-xs={$reportType === 'mingguan'}
					class:bg-transparent={$reportType !== 'mingguan'}
					class:text-slate-450={$reportType !== 'mingguan'}
					class:hover:text-slate-700={$reportType !== 'mingguan'}
				>
					Laporan Mingguan
				</button>
			</div>



			<!-- Student search (visible under Laporan Per Siswa) -->
			{#if $activeTab === 'siswa'}
				<div class="flex items-center gap-3">
					{#if !$selectedStudent}
						<div class="relative w-72">
							<span class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3.5">
								<Search class="h-3.5 w-3.5 text-slate-400" />
							</span>
							<input
								type="text"
								placeholder="Cari nama atau NIS siswa..."
								bind:value={$studentSearchQuery}
								oninput={searchForStudents}
								class="w-full rounded-xl border border-slate-200 bg-slate-50/50 py-2.5 pr-4 pl-9 text-xs text-slate-700 placeholder-slate-400 transition-all outline-none focus:border-[#00a294] focus:bg-white"
							/>
							{#if $foundStudents.length > 0}
								<div
									class="border-slate-150 absolute right-0 left-0 z-10 mt-1 max-h-56 overflow-y-auto rounded-xl border bg-white text-left shadow-md"
								>
									{#each $foundStudents as student}
										<button
											onclick={() => selectStudent(student)}
											class="block w-full cursor-pointer border-b border-none border-slate-100 bg-transparent px-4 py-2.5 text-left text-xs font-medium text-slate-700 last:border-0 hover:bg-slate-50"
										>
											<span class="block font-bold uppercase">{student.nama}</span>
											<span class="mt-0.5 font-mono text-[10px] text-slate-400"
												>{student.nis} - {student.kelas}</span
											>
										</button>
									{/each}
								</div>
							{/if}
						</div>
					{:else}
						<div
							class="flex items-center gap-2 rounded-xl border border-[#00a294]/20 bg-[#00a294]/5 px-3.5 py-1.5"
						>
							<div class="text-left leading-tight">
								<span class="block text-xs font-bold text-slate-800 uppercase"
									>{$selectedStudent?.nama}</span
								>
								<span class="text-slate-450 font-mono text-[9px]"
									>{$selectedStudent?.nis} - {$selectedStudent?.kelas}</span
								>
							</div>
							<button
								onclick={clearSelectedStudent}
								class="ml-2 cursor-pointer border-none bg-transparent text-xs font-bold text-rose-600 hover:text-rose-700"
							>
								Ganti
							</button>
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<div class="flex items-center gap-2.5 w-full sm:w-auto justify-end">
			<!-- Refresh Button -->
			<button
				onclick={() => loadData()}
				disabled={$activeTab === 'siswa' && !$selectedStudent}
				class="flex cursor-pointer items-center gap-1.5 rounded-xl border border-slate-200 bg-white hover:bg-slate-50 px-4.5 py-2.5 text-xs font-bold text-slate-700 shadow-xxs transition-all active:scale-98 disabled:cursor-not-allowed disabled:opacity-40"
			>
				<RefreshCw class="h-3.5 w-3.5 text-slate-650" />
				Refresh Data
			</button>

			<!-- Download PDF -->
			<button
				onclick={handleDownloadPDF}
				disabled={$activeTab === 'siswa' && !$selectedStudent}
				class="flex cursor-pointer items-center gap-1.5 rounded-xl border-none bg-slate-800 px-4.5 py-2.5 text-xs font-bold text-white shadow-xs transition-all hover:bg-slate-900 active:scale-98 disabled:cursor-not-allowed disabled:opacity-40"
			>
				<Download class="h-3.5 w-3.5 text-white" />
				Download Rekap PDF
			</button>
		</div>
	</div>
{/if}

	<!-- Content Render Area -->
	{#if $loading}
		<div
			class="flex flex-col items-center justify-center gap-3 rounded-2xl border border-slate-100/80 bg-white p-16 shadow-xs"
		>
			<Loader2 class="h-8 w-8 animate-spin text-slate-400" />
			<p class="text-xs font-semibold text-slate-400">Memproses rekap laporan kehadiran...</p>
		</div>
	{:else}
		<!-- TAB 1: LAPORAN KELAS -->
		{#if $activeTab === 'kelas'}
			{@const activeList = $reportType === 'bulanan' ? $items : $weeklyClassItems}
			<div class="space-y-6">
				<LaporanCards reportType={$reportType} {classSummary} />
				<LaporanCharts
					reportType={$reportType}
					selectedBulan={$selectedBulan}
					selectedKelas={$selectedKelas}
					formattedWeekRange={formatDateRangeStr($selectedWeekMonday)}
					{weeklyChartData}
					{monthlyChartData}
				/>
				<LaporanTableKelas 
					{activeList} 
					summaryByClass={$summaryByClass} 
					selectedKelas={$selectedKelas} 
				/>
			</div>
		{:else if $activeTab === 'siswa'}
			<!-- TAB 2: LAPORAN PER SISWA -->
			{#if $selectedStudent}
				<LaporanDetailSiswa
					studentSummary={$studentSummary}
					studentLogs={$studentLogs}
					monthlyTrend={$studentMonthlyTrend}
					reportType={$reportType}
				/>
			{:else}
				<!-- Daftar rekap semua siswa (terpaginasi) -->
				{@const rows = $studentRekapList}
				<div class="space-y-4">
					<!-- Toolbar: page size + info -->
					<div
						class="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-slate-100/80 bg-white p-4 shadow-xs"
					>
						<div class="flex items-center gap-2">
							<span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase"
								>Tampilkan</span
							>
							{#each [{ v: 50, l: '50' }, { v: 100, l: '100' }, { v: 0, l: 'Semua' }] as opt}
								<button
									onclick={() => setStudentRekapLimit(opt.v)}
									class="cursor-pointer rounded-lg border px-3 py-1.5 text-xs font-bold transition-all"
									class:border-[#00a294]={$studentRekapLimit === opt.v}
									class:bg-[#00a294]={$studentRekapLimit === opt.v}
									class:text-white={$studentRekapLimit === opt.v}
									class:border-slate-200={$studentRekapLimit !== opt.v}
									class:bg-white={$studentRekapLimit !== opt.v}
									class:text-slate-600={$studentRekapLimit !== opt.v}
								>
									{opt.l}
								</button>
							{/each}
							<span class="ml-1 text-xs font-medium text-slate-400">nama siswa</span>
						</div>
						<div class="text-xs font-semibold text-slate-500">
							Total <span class="font-bold text-slate-800">{$studentRekapTotal}</span> siswa
						</div>
					</div>

					<!-- Tabel rekap per siswa -->
					<div
						class="overflow-hidden rounded-2xl border border-slate-100/80 bg-white shadow-xs"
					>
						<div class="overflow-x-auto">
							<table class="w-full border-collapse text-left text-xs">
								<thead>
									<tr
										class="border-b border-slate-100 bg-slate-50/80 text-[10px] font-bold tracking-wider text-slate-400 uppercase"
									>
										<th class="w-12 p-4 pl-6 text-center">No</th>
										<th class="p-4">NIS</th>
										<th class="p-4">Nama Siswa</th>
										<th class="p-4">Kelas</th>
										<th class="w-16 p-4 text-center">Hadir</th>
										<th class="w-16 p-4 text-center">Izin</th>
										<th class="w-16 p-4 text-center">Sakit</th>
										<th class="w-16 p-4 text-center">Alpa</th>
										<th class="w-16 p-4 text-center">Magang</th>
										<th class="w-24 p-4 pr-6 text-center">Kehadiran</th>
									</tr>
								</thead>
								<tbody>
									{#if $studentRekapLoading}
										<tr>
											<td colspan="10" class="p-10 text-center">
												<Loader2 class="mx-auto h-6 w-6 animate-spin text-slate-300" />
											</td>
										</tr>
									{:else if rows.length === 0}
										<tr>
											<td colspan="10" class="p-10 text-center font-medium text-slate-400">
												Tidak ada data rekap siswa untuk periode ini.
											</td>
										</tr>
									{:else}
										{#each rows as item, idx}
											{@const startNo =
												$studentRekapLimit > 0
													? ($studentRekapPage - 1) * $studentRekapLimit
													: 0}
											<tr
												class="border-b border-slate-50 transition-colors hover:bg-slate-50/20 {!item.adaData
													? 'opacity-60'
													: ''}"
											>
												<td class="p-4 pl-6 text-center font-bold text-slate-400"
													>{startNo + idx + 1}</td
												>
												<td class="text-slate-550 p-4 font-mono font-medium">{item.nis}</td>
												<td class="p-4">
													<button
														onclick={() =>
															selectStudent({
																nis: item.nis,
																nama: item.namaSiswa,
																kelas: item.kelas
															})}
														class="cursor-pointer border-none bg-transparent p-0 text-left font-bold tracking-wide text-slate-800 uppercase hover:text-[#00a294] hover:underline"
													>
														{item.namaSiswa}
													</button>
													{#if !item.adaData}
														<span
															class="ml-1.5 rounded-md bg-slate-100 px-1.5 py-0.5 text-[9px] font-bold text-slate-400"
															>Belum ada data</span
														>
													{/if}
												</td>
												<td class="text-slate-550 p-4 font-semibold">{item.kelas}</td>
												<td class="p-4 text-center font-bold text-emerald-600">{item.totalHadir}</td>
												<td class="text-slate-550 p-4 text-center font-semibold">{item.totalIzin}</td>
												<td class="text-slate-550 p-4 text-center font-semibold">{item.totalSakit}</td>
												<td class="p-4 text-center font-bold text-rose-600">{item.totalAlpa}</td>
												<td class="text-slate-550 p-4 text-center font-semibold">{item.totalMagang}</td>
												<td
													class="p-4 pr-6 text-center font-bold {item.tingkatKehadiran >= 80
														? 'text-emerald-600'
														: item.tingkatKehadiran >= 70
															? 'text-amber-600'
															: 'text-rose-600'}"
												>
													{item.tingkatKehadiran.toFixed(1)}%
												</td>
											</tr>
										{/each}
									{/if}
								</tbody>
							</table>
						</div>

						<!-- Footer pagination -->
						{#if $studentRekapLimit > 0 && $studentRekapTotal > 0}
							<div
								class="flex items-center justify-between border-t border-slate-100 px-5 py-3"
							>
								<span class="text-xs font-medium text-slate-400">
									Halaman {$studentRekapPage} dari {Math.max(
										1,
										Math.ceil($studentRekapTotal / $studentRekapLimit)
									)}
								</span>
								<div class="flex items-center gap-2">
									<button
										onclick={() => changeStudentRekapPage(-1)}
										disabled={$studentRekapPage <= 1}
										class="shadow-xxs flex cursor-pointer items-center gap-1 rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-xs font-bold text-slate-700 transition-all hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
									>
										<ChevronLeft class="h-3.5 w-3.5" /> Sebelumnya
									</button>
									<button
										onclick={() => changeStudentRekapPage(1)}
										disabled={!$studentRekapHasMore}
										class="shadow-xxs flex cursor-pointer items-center gap-1 rounded-lg border border-slate-200 bg-white px-3 py-1.5 text-xs font-bold text-slate-700 transition-all hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-40"
									>
										Berikutnya <ChevronRight class="h-3.5 w-3.5" />
									</button>
								</div>
							</div>
						{/if}
					</div>
				</div>
			{/if}
		{:else if $activeTab === 'lomba'}
			<LaporanLombaKebersihan />
		{/if}
	{/if}
</div>
