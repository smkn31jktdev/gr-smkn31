<script lang="ts">
	import { onMount } from 'svelte';
	import { Download, RefreshCw, SlidersHorizontal } from 'lucide-svelte';
	import Table from './table/Table.svelte';
	import Pagination from './pagination/Pagination.svelte';
	import { addToast } from '../../../../stores/uiStore';
	import type { Kehadiran } from '../../../student/types/student.types';
	import AttendanceDetailModal from '../../../admin/components/dashboard/AttendanceDetailModal.svelte';
	import { getUploadUrl } from '../../../../api/client';
	import DatePicker from '../../../shared/components/DatePicker.svelte';
	import SearchBar from '../../../shared/components/SearchBar.svelte';
	import Badges from '../../../shared/components/Badges.svelte';
	import FilterPill from '../../../shared/components/FilterPill.svelte';
	import {
		logs,
		total,
		page,
		limit,
		loading,
		hasMore,
		kelasList,
		selectedKelas,
		selectedTanggal,
		searchQuery,
		selectedStatus,
		urutkanWaktu,
		urutkanNama,
		loadData,
		loadKelasJurusan,
		handleFilter,
		handleDelete
	} from '../../logic/piketMonitoringLogic';

	onMount(() => {
		loadKelasJurusan();
		loadData();
	});

	// Client-side filtering & sorting fallback
	let filteredLogs = $derived.by(() => {
		let list = [...$logs];

		// Sort by waktuAbsen
		if ($urutkanWaktu === 'terbaru') {
			list.sort((a, b) => {
				const hasA = !!a.waktuAbsen;
				const hasB = !!b.waktuAbsen;
				if (hasA && !hasB) return -1;
				if (!hasA && hasB) return 1;
				if (!hasA && !hasB) return 0;
				return b.waktuAbsen.localeCompare(a.waktuAbsen); // latest check-in first
			});
		} else if ($urutkanWaktu === 'terlama') {
			list.sort((a, b) => {
				const hasA = !!a.waktuAbsen;
				const hasB = !!b.waktuAbsen;
				if (hasA && !hasB) return -1;
				if (!hasA && hasB) return 1;
				if (!hasA && !hasB) return 0;
				return a.waktuAbsen.localeCompare(b.waktuAbsen); // earliest check-in first
			});
		}

		// Sort by namaSiswa
		if ($urutkanNama === 'nama_asc') {
			list.sort((a, b) => {
				const nA = a.namaSiswa || '';
				const nB = b.namaSiswa || '';
				return nA.localeCompare(nB);
			});
		} else if ($urutkanNama === 'nama_desc') {
			list.sort((a, b) => {
				const nA = a.namaSiswa || '';
				const nB = b.namaSiswa || '';
				return nB.localeCompare(nA);
			});
		}

		return list;
	});

	// Mutually exclusive sort handlers
	function handleSortWaktuChange() {
		if ($urutkanWaktu !== 'default') {
			urutkanNama.set('default');
		}
	}

	function handleSortNamaChange() {
		if ($urutkanNama !== 'default') {
			urutkanWaktu.set('default');
		}
	}

	// Calculate counts based on today's logs
	let totalCount = $derived($logs.length);
	let hadirCount = $derived($logs.filter((l) => l.status === 'hadir').length);
	let izinSakitCount = $derived(
		$logs.filter((l) => l.status === 'izin' || l.status === 'sakit').length
	);
	let alpaCount = $derived($logs.filter((l) => l.status === 'tidak_hadir').length);

	// Indonesian date formatter
	function formatIndonesianDateStr(dateStr: string): string {
		if (!dateStr) return '';
		const parts = dateStr.split('-');
		if (parts.length !== 3) return dateStr;
		const d = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]));
		const days = ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'];
		const months = [
			'Januari',
			'Februari',
			'Maret',
			'April',
			'Mei',
			'Juni',
			'Juli',
			'Agustus',
			'September',
			'Oktober',
			'November',
			'Desember'
		];
		return `${days[d.getDay()]}, ${d.getDate()} ${months[d.getMonth()]} ${d.getFullYear()}`;
	}

	// PDF download simulation
	function handleDownloadPDF() {
		const currentKelasVal = $selectedKelas || 'Semua Kelas';
		addToast(`Mempersiapkan unduhan PDF Rekap Absensi untuk ${currentKelasVal}...`, 'info');

		setTimeout(() => {
			try {
				const dummyContent =
					`REKAPITULASI HARIAN ABSENSI SISWA - SMKN 31 JAKARTA\n` +
					`Tanggal: ${formatIndonesianDateStr($selectedTanggal)}\n` +
					`Kelas: ${currentKelasVal}\n` +
					`Total Kehadiran: ${totalCount} siswa\n` +
					`-----------------------------------------\n` +
					`Siswa Hadir: ${hadirCount}\n` +
					`Siswa Izin/Sakit: ${izinSakitCount}\n` +
					`Siswa Alpa: ${alpaCount}\n` +
					`-----------------------------------------\n` +
					`Tanggal Unduh: ${new Date().toLocaleDateString('id-ID')}\n` +
					`Status Unduhan: SUKSES`;

				const blob = new Blob([dummyContent], { type: 'application/pdf' });
				const url = URL.createObjectURL(blob);
				const link = document.createElement('a');
				link.href = url;
				link.download = `Rekap_Absen_${$selectedTanggal.replace(/-/g, '_')}_${currentKelasVal.replace(/\s+/g, '_')}.pdf`;
				link.click();
				URL.revokeObjectURL(url);
				addToast(`Berhasil mengunduh rekap PDF untuk ${currentKelasVal}`, 'success');
			} catch (err) {
				addToast('Gagal mengunduh file PDF', 'error');
			}
		}, 1000);
	}

	// Attendance detail modal state
	let attendanceModalOpen = $state(false);
	let selectedAttendanceLog = $state<any>(null);

	function openAttendanceDetail(log: Kehadiran) {
		selectedAttendanceLog = log;
		attendanceModalOpen = true;
	}
</script>

<div class="space-y-5 pb-10 font-sans text-slate-700 select-none">
	<!-- Header Section -->
	<div class="flex flex-wrap items-start justify-between gap-4 border-b border-slate-100 pb-5">
		<div class="space-y-1 text-left">
			<h1 class="text-xl font-bold tracking-tight text-slate-800">Monitoring Absensi Siswa</h1>
			<p class="text-xs font-medium text-slate-400">
				Pantau kehadiran, izin, dan siswa tanpa keterangan berdasarkan kelas.
			</p>
		</div>

		<div class="flex items-center gap-3">
			<!-- Custom Date Picker -->
			<div class="w-48 text-left">
				<DatePicker
					bind:value={$selectedTanggal}
					onchange={handleFilter}
					placeholder="Pilih tanggal"
				/>
			</div>

			<!-- Segarkan / Loading State -->
			<button
				onclick={loadData}
				disabled={$loading}
				class="shadow-xxs flex cursor-pointer items-center justify-center rounded-xl border border-slate-200 bg-white p-2.5 text-slate-500 transition-all hover:bg-slate-50"
				title="Segarkan data"
			>
				<RefreshCw class="h-4 w-4 {$loading ? 'animate-spin' : ''}" />
			</button>
		</div>
	</div>

	<!-- Search and Counters Box (Pill Style, Borderless/Soft) -->
	<div class="flex flex-col items-center justify-between gap-4 md:flex-row">
		<!-- Search Bar Pill -->
		<SearchBar
			bind:value={$searchQuery}
			placeholder="Cari nama, kelas, atau nis..."
			oninput={handleFilter}
			rounded="full"
			class="w-full md:flex-1"
		/>

		<!-- Count Badges Pill on the right of Search -->
		<Badges total={totalCount} hadir={hadirCount} izinSakit={izinSakitCount} alpa={alpaCount} />
	</div>

	<!-- Filters Bar (Redesigned as clean horizontal pills toolbar) -->
	<div
		class="flex flex-wrap items-center gap-4 rounded-2xl border border-slate-100/80 bg-white p-4 text-left shadow-xs"
	>
		<div
			class="text-slate-450 flex items-center gap-2 pr-1 text-xs font-bold tracking-wider uppercase"
		>
			<SlidersHorizontal class="h-4 w-4 text-slate-400" />
			<span>Filter</span>
		</div>

		<!-- Kelas select pill -->
		<FilterPill
			options={[
				{ value: '', label: 'Semua Kelas' },
				...$kelasList.map((k) => ({ value: k, label: `Kelas ${k}` }))
			]}
			bind:value={$selectedKelas}
			onchange={handleFilter}
			placeholder="Semua Kelas"
			minWidth="130px"
		/>

		<!-- Status select pill -->
		<FilterPill
			options={[
				{ value: '', label: 'Semua Status' },
				{ value: 'hadir', label: 'Masuk' },
				{ value: 'izin_sakit', label: 'Ijin/Sakit' },
				{ value: 'magang', label: 'Magang' },
				{ value: 'tidak_hadir', label: 'Tanpa Keterangan' }
			]}
			bind:value={$selectedStatus}
			onchange={handleFilter}
			placeholder="Semua Status"
			minWidth="130px"
		/>

		<!-- Sort Nama select pill (pushed to right on wider screens) -->
		<FilterPill
			options={[
				{ value: 'default', label: 'Urutan Nama: Default' },
				{ value: 'nama_asc', label: 'Nama: A - Z' },
				{ value: 'nama_desc', label: 'Nama: Z - A' }
			]}
			bind:value={$urutkanNama}
			onchange={handleSortNamaChange}
			placeholder="Urutan Nama"
			minWidth="145px"
			class="sm:ml-auto"
		/>

		<!-- Sort Waktu select pill -->
		<FilterPill
			options={[
				{ value: 'default', label: 'Urutan Waktu: Default' },
				{ value: 'terbaru', label: 'Urutan Waktu: Terbaru' },
				{ value: 'terlama', label: 'Urutan Waktu: Terlama' }
			]}
			bind:value={$urutkanWaktu}
			onchange={handleSortWaktuChange}
			placeholder="Urutan Waktu"
			minWidth="145px"
			align="right"
		/>
	</div>

	<!-- Rekap PDF Info Block -->
	<div class="flex flex-col items-center justify-between gap-3.5 px-1 sm:flex-row">
		<p class="text-left text-[11px] font-medium text-slate-400">
			Rekap PDF harian menggunakan tanggal terpilih dengan cakupan <span
				class="text-slate-650 font-semibold"
				>{$selectedKelas || 'Semua Kelas'} ({totalCount} siswa)</span
			>
		</p>

		<button
			onclick={handleDownloadPDF}
			class="flex cursor-pointer items-center gap-1.5 rounded-xl border-none bg-slate-800 px-4.5 py-2.5 text-xs font-bold text-white shadow-xs transition-all hover:bg-slate-900 active:scale-98"
		>
			<Download class="h-3.5 w-3.5 text-white" />
			Download Rekap PDF
		</button>
	</div>

	<!-- Logs Table (Siswa, Kelas Lengkap, Status, Waktu, Keterangan) -->
	<Table loading={$loading} {filteredLogs} onDelete={handleDelete} onOpenPermit={openAttendanceDetail} />

	<!-- Pagination Controls -->
	<Pagination
		bind:limit={$limit}
		bind:page={$page}
		total={$total}
		hasMore={$hasMore}
		loading={$loading}
		onFilter={handleFilter}
		onLoadData={loadData}
	/>
</div>

<AttendanceDetailModal
	bind:show={attendanceModalOpen}
	log={selectedAttendanceLog}
	onclose={() => (attendanceModalOpen = false)}
/>
