<script lang="ts">
	import { listKehadiranSiswa, submitKehadiran, uploadIzinFile } from '../../logic/kehadiranLogic';
	import { addToast } from '../../../../stores/uiStore';
	import Modal from '../../../shared/components/Modal.svelte';
	import { Loader2 } from 'lucide-svelte';

	let {
		show = $bindable(false),
		onload
	} = $props<{
		show: boolean;
		onload: () => void;
	}>();

	let magangLoading = $state(false);
	let magangAlasan = $state('');
	let magangFoto = $state('');
	let magangSubmitLoading = $state(false);
	let showMagangModal = $state(false);

	async function checkMagangStatus() {
		magangLoading = true;
		try {
			const res = await listKehadiranSiswa(undefined, undefined, 1, 50);
			const previousMagang = res.items.find(
				(item: any) => item.status === 'magang' && item.fotoIzin
			);

			if (previousMagang) {
				addToast('Riwayat verifikasi magang terdeteksi. Mengirim absensi...', 'info');
				const success = await submitKehadiran({
					status: 'magang',
					alasan: previousMagang.alasan || 'Absensi Magang Mandiri',
					fotoIzin: previousMagang.fotoIzin
				});
				if (success) {
					onload();
				}
				show = false;
			} else {
				magangAlasan = '';
				magangFoto = '';
				showMagangModal = true;
			}
		} catch (err) {
			console.error(err);
			addToast('Gagal memverifikasi status magang', 'error');
			show = false;
		} finally {
			magangLoading = false;
		}
	}

	async function handleMagangFileChange(e: Event) {
		const input = e.target as HTMLInputElement;
		if (input.files && input.files[0]) {
			const file = input.files[0];
			addToast('Mengunggah berkas bukti magang...', 'info');
			const url = await uploadIzinFile(file);
			if (url) {
				magangFoto = url;
			}
		}
	}

	async function handleMagangSubmit() {
		if (!magangFoto) {
			addToast('Harap unggah bukti magang terlebih dahulu!', 'error');
			return;
		}
		if (!magangAlasan.trim()) {
			addToast('Harap isi keterangan/lokasi magang!', 'error');
			return;
		}

		magangSubmitLoading = true;
		const success = await submitKehadiran({
			status: 'magang',
			alasan: magangAlasan.trim(),
			fotoIzin: magangFoto
		});

		magangSubmitLoading = false;
		if (success) {
			showMagangModal = false;
			show = false;
			onload();
		}
	}

	$effect(() => {
		if (show) {
			checkMagangStatus();
		}
	});
</script>

<!-- Modal Registrasi Magang Pertama Kali -->
<Modal
	show={showMagangModal}
	title=""
	onclose={() => {
		showMagangModal = false;
		show = false;
	}}
>
	<div class="p-6 space-y-6">
		<div>
			<h3 class="text-base font-bold text-slate-800">Registrasi Absensi Magang</h3>
			<p class="text-xs text-slate-400 mt-1">
				Karena ini adalah absensi magang pertama Anda, silakan unggah bukti magang (Surat Tugas/MOU/Pernyataan) dan masukkan keterangan lokasi.
			</p>
		</div>

		<div class="space-y-4">
			<!-- Upload Proof -->
			<div>
				<label class="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-2">Unggah Bukti Magang (Maks 5MB)</label>
				<input
					type="file"
					accept="image/*,.pdf"
					onchange={handleMagangFileChange}
					class="input w-full file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-xs file:font-semibold file:bg-purple-50 file:text-purple-700 hover:file:bg-purple-100"
				/>
			</div>

			{#if magangFoto}
				<div class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-800 text-xs font-semibold flex items-center justify-between">
					<span>✓ Berkas berhasil diunggah</span>
					<a href={magangFoto} target="_blank" class="underline text-purple-700 font-bold">Lihat Berkas</a>
				</div>
			{/if}

			<!-- Description -->
			<div>
				<label for="magang_alasan" class="block text-xs font-bold uppercase tracking-wider text-slate-500 mb-2">Lokasi / Keterangan Magang</label>
				<textarea
					id="magang_alasan"
					placeholder="Contoh: Magang di PT. Maju Jaya (Divisi IT Support)..."
					bind:value={magangAlasan}
					class="w-full border border-slate-200 rounded-xl p-3 text-sm focus:outline-none focus:ring-2 focus:ring-purple-500/20 focus:border-purple-500 min-h-[80px]"
				></textarea>
			</div>

			<!-- Submit Button -->
			<div class="pt-2">
				<button
					type="button"
					onclick={handleMagangSubmit}
					disabled={magangSubmitLoading || !magangFoto || !magangAlasan.trim()}
					class="w-full py-3 bg-purple-600 hover:bg-purple-700 text-white rounded-xl font-bold text-sm transition-all duration-150 active:scale-98 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer"
				>
					{#if magangSubmitLoading}
						<Loader2 class="h-4 w-4 animate-spin text-white" />
						Memproses Registrasi...
					{:else}
						Registrasi & Kirim Absen
					{/if}
				</button>
			</div>
		</div>
	</div>
</Modal>

{#if magangLoading}
	<!-- Back-drop/Loading spinner when auto-verifying first -->
	<div class="fixed inset-0 z-50 bg-slate-900/20 backdrop-blur-xs flex items-center justify-center">
		<div class="bg-white p-5 rounded-2xl border border-slate-100 flex flex-col items-center gap-3 shadow-lg">
			<Loader2 class="h-8 w-8 animate-spin text-purple-600" />
			<span class="text-xs font-bold text-slate-650">Memverifikasi status magang...</span>
		</div>
	</div>
{/if}
