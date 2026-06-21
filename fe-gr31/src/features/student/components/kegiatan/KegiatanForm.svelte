<script lang="ts">
	import { onMount } from 'svelte';
	import { getJurnalForDate, submitJurnal } from '../../logic/kegiatanLogic';
	import { listKehadiranSiswa } from '../../logic/kehadiranLogic';
	import { currentUser } from '../../../../stores/authStore';
	import { addToast } from '../../../../stores/uiStore';
	import type { Kehadiran, G7Jurnal } from '../../types/student.types';
	import { Loader } from 'lucide-svelte';

	import StudentProfileCard from './cards/StudentProfileCard.svelte';
	import DateKehadiranCard from './cards/DateKehadiranCard.svelte';
	import BangunPagiCard from './cards/BangunPagiCard.svelte';
	import TidurMalamCard from './cards/TidurMalamCard.svelte';
	import BeribadahCard from './cards/BeribadahCard.svelte';
	import MakanSehatCard from './cards/MakanSehatCard.svelte';
	import OlahragaCard from './cards/OlahragaCard.svelte';
	import BelajarCard from './cards/BelajarCard.svelte';
	import BermasyarakatCard from './cards/BermasyarakatCard.svelte';

	import Modal from '../../../shared/components/Modal.svelte';
	import KehadiranForm from '../absensi/KehadiranForm.svelte';
	import RamadhanCard from './cards/RamadhanCard.svelte';
	import { apiRequest } from '../../../../api/client';

	const getTodayStr = () => new Date().toLocaleDateString('sv-SE');

	let tanggal = $state(getTodayStr());
	let loadingJurnal = $state(false);
	let loadingKehadiran = $state(false);
	let savingCard = $state(false);
	let showAttendanceModal = $state(false);
	let fullJurnalData = $state<G7Jurnal | null>(null);
	let todayKehadiran = $state<Kehadiran | null>(null);
	let isRamadan = $state(false);
	let ramadanDay = $state(0);

	let name = $derived($currentUser?.nama || 'Siswa');
	let kelas = $derived($currentUser?.kelas || 'Kelas');
	let nis = $derived($currentUser?.nis || '-');

	// G7 Form State
	let formState = $state({
		bangun: { waktu: '', doa: false },
		tidur: { waktu: '', doa: false },
		ibadah: {
			doa_ortu: false,
			sholat_fajar: false,
			sholat_5waktu: false,
			zikir: false,
			dhuha: false,
			rowatib: false,
			infaq: 0,
			puasa: false,
			tarawihRokaat: 0
		},
		makan: { utama: '', lauk: '', sayurBuah: false, susuSuplemen: false },
		olahraga: { aktivitas: '', durasi: 30 },
		belajar: {
			belajarMandiri: false,
			kitabSuci: false,
			bukuUmum: false,
			bukuMapel: false,
			tugasPR: false
		},
		bermasyarakat: { kegiatan: '', lokasi: '', waktu: '', diketahuiOT: false }
	});

	let checkingRamadan = false;
	async function checkRamadan() {
		if (checkingRamadan) return;
		checkingRamadan = true;
		try {
			const res = await apiRequest<{ is_ramadan: boolean; ramadan_day: number }>(`/v1/puasa/check?tanggal=${tanggal}`);
			if (res.data) {
				isRamadan = res.data.is_ramadan;
				ramadanDay = res.data.ramadan_day;
			} else {
				isRamadan = false;
				ramadanDay = 0;
			}
		} catch (e) {
			console.error(e);
			isRamadan = false;
			ramadanDay = 0;
		} finally {
			checkingRamadan = false;
		}
	}

	async function loadData() {
		loadingJurnal = true;
		loadingKehadiran = true;
		try {
			await checkRamadan();
			const jData = await getJurnalForDate(tanggal);
			fullJurnalData = jData;
			if (jData) parseJurnal(jData);
			else resetState();

			const attRes = await listKehadiranSiswa(tanggal, tanggal);
			todayKehadiran = attRes.items.length > 0 ? attRes.items[0] : null;
		} catch (e) {
			addToast('Gagal memuat data harian', 'error');
		} finally {
			loadingJurnal = false;
			loadingKehadiran = false;
		}
	}

	function resetState() {
		formState.bangun = { waktu: '', doa: false };
		formState.tidur = { waktu: '', doa: false };
		formState.ibadah = {
			doa_ortu: false,
			sholat_fajar: false,
			sholat_5waktu: false,
			zikir: false,
			dhuha: false,
			rowatib: false,
			infaq: 0,
			puasa: false,
			tarawihRokaat: 0
		};
		formState.makan = { utama: '', lauk: '', sayurBuah: false, susuSuplemen: false };
		formState.olahraga = { aktivitas: '', durasi: 30 };
		formState.belajar = {
			belajarMandiri: false,
			kitabSuci: false,
			bukuUmum: false,
			bukuMapel: false,
			tugasPR: false
		};
		formState.bermasyarakat = { kegiatan: '', lokasi: '', waktu: '', diketahuiOT: false };
	}

	function parseJurnal(data: G7Jurnal) {
		if (data.bangun) {
			formState.bangun.waktu = data.bangun.waktu || '';
			formState.bangun.doa = data.bangun.keterangan?.includes('Membaca Doa: Ya') ?? false;
		} else formState.bangun = { waktu: '', doa: false };

		if (data.tidur) {
			formState.tidur.waktu = data.tidur.waktu || '';
			formState.tidur.doa = data.tidur.keterangan?.includes('Membaca Doa: Ya') ?? false;
		} else formState.tidur = { waktu: '', doa: false };

		if (data.ibadah) {
			const ket = data.ibadah.keterangan || '';
			formState.ibadah.doa_ortu = ket.includes('Berdoa');
			formState.ibadah.sholat_fajar = ket.includes('Sholat Fajar');
			formState.ibadah.sholat_5waktu = ket.includes('Sholat 5 Waktu');
			formState.ibadah.zikir = ket.includes('Zikir');
			formState.ibadah.dhuha = ket.includes('Sholat Dhuha');
			formState.ibadah.rowatib = ket.includes('Sholat Sunah Rowatib');
			const infaqMatch = ket.match(/Infaq: Rp ([\d.]+)/);
			formState.ibadah.infaq = infaqMatch ? Number(infaqMatch[1].replace(/\./g, '')) || 0 : 0;
			
			// Parse Ramadan fields
			formState.ibadah.puasa = ket.includes('Puasa: Ya');
			const tarawihMatch = ket.match(/Tarawih: (\d+) Rokaat/);
			formState.ibadah.tarawihRokaat = tarawihMatch ? Number(tarawihMatch[1]) : 0;
		} else {
			formState.ibadah = {
				doa_ortu: false,
				sholat_fajar: false,
				sholat_5waktu: false,
				zikir: false,
				dhuha: false,
				rowatib: false,
				infaq: 0,
				puasa: false,
				tarawihRokaat: 0
			};
		}

		if (data.makan) {
			const ket = data.makan.keterangan || '';
			const utMatch = ket.match(/Makanan Utama: ([^,]+)/);
			const lMatch = ket.match(/Lauk: ([^,]+)/);
			formState.makan.utama = utMatch ? utMatch[1] : '';
			formState.makan.lauk = lMatch ? lMatch[1] : '';
			formState.makan.sayurBuah = ket.includes('Sayur/Buah: Ya');
			formState.makan.susuSuplemen = ket.includes('Susu/Suplemen: Ya');
		} else formState.makan = { utama: '', lauk: '', sayurBuah: false, susuSuplemen: false };

		if (data.olahraga) {
			const ket = data.olahraga.keterangan || '';
			const aktMatch = ket.match(/Aktivitas: ([^,]+)/);
			const durMatch = ket.match(/Durasi: (\d+) menit/);
			formState.olahraga.aktivitas = aktMatch ? aktMatch[1] : '';
			formState.olahraga.durasi = durMatch ? Number(durMatch[1]) : 30;
		} else formState.olahraga = { aktivitas: '', durasi: 30 };

		if (data.belajar) {
			formState.belajar.belajarMandiri = data.belajar.done;
			const ket = data.belajar.keterangan || '';
			formState.belajar.kitabSuci = ket.includes('Kitab Suci');
			formState.belajar.bukuUmum = ket.includes('Buku Umum');
			formState.belajar.bukuMapel = ket.includes('Buku Pelajaran');
			formState.belajar.tugasPR = ket.includes('Tugas/PR');
		} else
			formState.belajar = {
				belajarMandiri: false,
				kitabSuci: false,
				bukuUmum: false,
				bukuMapel: false,
				tugasPR: false
			};

		if (data.bermasyarakat) {
			formState.bermasyarakat.waktu = data.bermasyarakat.waktu || '';
			const ket = data.bermasyarakat.keterangan || '';
			const kegMatch = ket.match(/Kegiatan: ([^,]+)/);
			const lokMatch = ket.match(/Lokasi: ([^,]+)/);
			formState.bermasyarakat.kegiatan = kegMatch ? kegMatch[1] : '';
			formState.bermasyarakat.lokasi = lokMatch ? lokMatch[1] : '';
			formState.bermasyarakat.diketahuiOT = ket.includes('Diketahui OT/RT: Ya');
		} else formState.bermasyarakat = { kegiatan: '', lokasi: '', waktu: '', diketahuiOT: false };
	}

	async function submitCard(payload: any) {
		savingCard = true;
		try {
			const merged = {
				bangun: fullJurnalData?.bangun || { done: false, waktu: '', keterangan: '' },
				ibadah: fullJurnalData?.ibadah || { done: false, waktu: '', keterangan: '' },
				makan: fullJurnalData?.makan || { done: false, waktu: '', keterangan: '' },
				olahraga: fullJurnalData?.olahraga || { done: false, waktu: '', keterangan: '' },
				belajar: fullJurnalData?.belajar || { done: false, waktu: '', keterangan: '' },
				bermasyarakat: fullJurnalData?.bermasyarakat || { done: false, waktu: '', keterangan: '' },
				tidur: fullJurnalData?.tidur || { done: false, waktu: '', keterangan: '' },
				...payload
			};
			const ok = await submitJurnal(tanggal, merged);
			if (ok) await loadData();
		} finally {
			savingCard = false;
		}
	}

	const saveBangun = async () =>
		submitCard({
			bangun: {
				done: !!formState.bangun.waktu,
				waktu: formState.bangun.waktu,
				keterangan: `Membaca Doa: ${formState.bangun.doa ? 'Ya' : 'Tidak'}`
			}
		});

	const saveTidur = async () =>
		submitCard({
			tidur: {
				done: !!formState.tidur.waktu,
				waktu: formState.tidur.waktu,
				keterangan: `Membaca Doa: ${formState.tidur.doa ? 'Ya' : 'Tidak'}`
			}
		});

	async function saveIbadah() {
		const list: string[] = [];
		if (formState.ibadah.doa_ortu) list.push('Berdoa');
		if (formState.ibadah.sholat_fajar) list.push('Sholat Fajar');
		if (formState.ibadah.sholat_5waktu) list.push('Sholat 5 Waktu');
		if (formState.ibadah.zikir) list.push('Zikir');
		if (formState.ibadah.dhuha) list.push('Sholat Dhuha');
		if (formState.ibadah.rowatib) list.push('Sholat Sunah Rowatib');

		const done = list.length > 0 || formState.ibadah.infaq > 0 || formState.ibadah.puasa || formState.ibadah.tarawihRokaat > 0;
		const infaqStr = formState.ibadah.infaq.toLocaleString('id-ID');

		let extra = '';
		if (isRamadan) {
			extra = `. Puasa: ${formState.ibadah.puasa ? 'Ya' : 'Tidak'}. Tarawih: ${formState.ibadah.tarawihRokaat} Rokaat`;
		}

		await submitCard({
			ibadah: {
				done,
				waktu: '',
				keterangan: `Ibadah: ${list.join(', ') || 'Nihil'}. Infaq: Rp ${infaqStr}${extra}`
			}
		});
	}

	const saveMakan = async () =>
		submitCard({
			makan: {
				done: !!formState.makan.utama,
				waktu: '',
				keterangan: `Makanan Utama: ${formState.makan.utama || 'Nihil'}, Lauk: ${formState.makan.lauk || 'Nihil'}, Sayur/Buah: ${formState.makan.sayurBuah ? 'Ya' : 'Tidak'}, Susu/Suplemen: ${formState.makan.susuSuplemen ? 'Ya' : 'Tidak'}`
			}
		});

	const saveOlahraga = async () =>
		submitCard({
			olahraga: {
				done: !!formState.olahraga.aktivitas,
				waktu: '',
				keterangan: `Aktivitas: ${formState.olahraga.aktivitas || 'Nihil'}, Durasi: ${formState.olahraga.durasi} menit`
			}
		});

	async function saveBelajar() {
		const list: string[] = [];
		if (formState.belajar.kitabSuci) list.push('Kitab Suci');
		if (formState.belajar.bukuUmum) list.push('Buku Umum');
		if (formState.belajar.bukuMapel) list.push('Buku Pelajaran');
		if (formState.belajar.tugasPR) list.push('Tugas/PR');
		await submitCard({
			belajar: {
				done: formState.belajar.belajarMandiri,
				waktu: '',
				keterangan: `Belajar Mandiri: ${formState.belajar.belajarMandiri ? 'Ya' : 'Tidak'}, Topik: ${list.join(', ') || 'Nihil'}`
			}
		});
	}

	const saveBermasyarakat = async () =>
		submitCard({
			bermasyarakat: {
				done: !!formState.bermasyarakat.kegiatan,
				waktu: formState.bermasyarakat.waktu,
				keterangan: `Kegiatan: ${formState.bermasyarakat.kegiatan || 'Nihil'}, Lokasi: ${formState.bermasyarakat.lokasi || 'Nihil'}, Diketahui OT/RT: ${formState.bermasyarakat.diketahuiOT ? 'Ya' : 'Tidak'}`
			}
		});

	// Realtime clock for unsaved card times (updates every second)
	let currentTime = $state('');

	$effect(() => {
		const updateTime = () => {
			const now = new Date();
			const hh = String(now.getHours()).padStart(2, '0');
			const mm = String(now.getMinutes()).padStart(2, '0');
			currentTime = `${hh}:${mm}`;

			// Only update if viewing today's date and the values haven't been saved yet
			if (tanggal === getTodayStr()) {
				if (!fullJurnalData?.bangun?.waktu) {
					formState.bangun.waktu = currentTime;
				}
				if (!fullJurnalData?.tidur?.waktu) {
					formState.tidur.waktu = currentTime;
				}
			}
		};

		updateTime();
		const interval = setInterval(updateTime, 1000);
		return () => clearInterval(interval);
	});

	onMount(() => {
		loadData();
	});
</script>

<div class="space-y-6">
	<StudentProfileCard {name} {nis} {kelas} />

	<DateKehadiranCard
		bind:tanggal
		maxDate={getTodayStr()}
		kehadiran={todayKehadiran}
		{loadingKehadiran}
		onchange={loadData}
		onabsen={() => (showAttendanceModal = true)}
	/>

	{#if loadingJurnal}
		<div
			class="flex flex-col items-center justify-center rounded-3xl border border-slate-100/90 bg-white p-16 text-slate-400"
		>
			<Loader class="mb-3 h-8 w-8 animate-spin text-[#4db6ac]" />
			<p class="text-xs font-bold">Memuat formulir kegiatan...</p>
		</div>
	{:else}
		<!-- Bangun Pagi & Tidur Malam -->
		<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
			<BangunPagiCard
				bind:waktu={formState.bangun.waktu}
				bind:doa={formState.bangun.doa}
				onsave={saveBangun}
				loading={savingCard}
			/>
			<TidurMalamCard
				bind:waktu={formState.tidur.waktu}
				bind:doa={formState.tidur.doa}
				onsave={saveTidur}
				loading={savingCard}
			/>
		</div>

		<!-- Beribadah -->
		<BeribadahCard
			bind:doa_ortu={formState.ibadah.doa_ortu}
			bind:sholat_fajar={formState.ibadah.sholat_fajar}
			bind:sholat_5waktu={formState.ibadah.sholat_5waktu}
			bind:zikir={formState.ibadah.zikir}
			bind:dhuha={formState.ibadah.dhuha}
			bind:rowatib={formState.ibadah.rowatib}
			bind:infaq={formState.ibadah.infaq}
			isRamadan={isRamadan}
			ramadanDay={ramadanDay}
			onsave={saveIbadah}
			loading={savingCard}
		/>

		{#if isRamadan}
			<RamadhanCard
				bind:puasa={formState.ibadah.puasa}
				bind:tarawihRokaat={formState.ibadah.tarawihRokaat}
				onsave={saveIbadah}
				loading={savingCard}
			/>
		{/if}

		<!-- Makan, Olahraga, Belajar -->
		<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
			<MakanSehatCard
				bind:utama={formState.makan.utama}
				bind:lauk={formState.makan.lauk}
				bind:sayurBuah={formState.makan.sayurBuah}
				bind:susuSuplemen={formState.makan.susuSuplemen}
				onsave={saveMakan}
				loading={savingCard}
			/>
			<OlahragaCard
				bind:aktivitas={formState.olahraga.aktivitas}
				bind:durasi={formState.olahraga.durasi}
				onsave={saveOlahraga}
				loading={savingCard}
			/>
			<BelajarCard
				bind:belajarMandiri={formState.belajar.belajarMandiri}
				bind:kitabSuci={formState.belajar.kitabSuci}
				bind:bukuUmum={formState.belajar.bukuUmum}
				bind:bukuMapel={formState.belajar.bukuMapel}
				bind:tugasPR={formState.belajar.tugasPR}
				onsave={saveBelajar}
				loading={savingCard}
			/>
		</div>

		<!-- Bermasyarakat -->
		<BermasyarakatCard
			bind:kegiatan={formState.bermasyarakat.kegiatan}
			bind:lokasi={formState.bermasyarakat.lokasi}
			bind:waktu={formState.bermasyarakat.waktu}
			bind:diketahuiOT={formState.bermasyarakat.diketahuiOT}
			onsave={saveBermasyarakat}
			loading={savingCard}
		/>
	{/if}
</div>

<Modal
	show={showAttendanceModal}
	title=""
	onclose={() => {
		showAttendanceModal = false;
		loadData();
	}}
>
	<div class="px-2">
		<KehadiranForm />
	</div>
</Modal>
