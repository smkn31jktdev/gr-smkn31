<script lang="ts">
	import { submitKehadiran } from '../../logic/kehadiranLogic';
	import { addToast } from '../../../../stores/uiStore';
	import Modal from '../../../shared/components/Modal.svelte';
	import { CheckCircle2, Loader2 } from 'lucide-svelte';

	let {
		show = $bindable(false),
		onload
	} = $props<{
		show: boolean;
		onload: () => void;
	}>();

	let qrScanning = $state(false);
	let qrSuccess = $state(false);
	let cameraError = $state<string | null>(null);

	let videoElement = $state<HTMLVideoElement | null>(null);
	let stream = $state<MediaStream | null>(null);
	let jsQR: any = null;
	let scanningActive = false;
	let scanCanvas: HTMLCanvasElement | null = null;

	function getTodayDateString(): string {
		const d = new Date();
		const year = d.getFullYear();
		const month = String(d.getMonth() + 1).padStart(2, '0');
		const day = String(d.getDate()).padStart(2, '0');
		return `${year}-${month}-${day}`;
	}

	async function loadJsQR(): Promise<void> {
		if (jsQR) return;
		if ((window as any).jsQR) {
			jsQR = (window as any).jsQR;
			return;
		}
		return new Promise<void>((resolve, reject) => {
			const script = document.createElement('script');
			script.src = 'https://cdn.jsdelivr.net/npm/jsqr@1.4.0/dist/jsQR.min.js';
			script.onload = () => {
				jsQR = (window as any).jsQR;
				resolve();
			};
			script.onerror = () => {
				reject(new Error('Gagal memuat pustaka scanner QR'));
			};
			document.head.appendChild(script);
		});
	}

	async function startCamera() {
		cameraError = null;
		scanningActive = true;
		qrScanning = true;
		qrSuccess = false;
		try {
			if (!jsQR) {
				await loadJsQR();
			}

			const mediaStream = await navigator.mediaDevices.getUserMedia({
				video: { facingMode: 'environment' }
			});
			stream = mediaStream;
			if (videoElement) {
				videoElement.srcObject = mediaStream;
				requestAnimationFrame(scanLoop);
			}
		} catch (err: any) {
			console.error('Gagal mengakses kamera:', err);
			cameraError = 'Tidak dapat mengakses kamera. Pastikan izin kamera telah diberikan.';
			addToast('Gagal mengakses kamera: ' + (err.message || err), 'error');
		}
	}

	function stopCamera() {
		scanningActive = false;
		if (stream) {
			stream.getTracks().forEach((track) => track.stop());
			stream = null;
		}
		if (videoElement) {
			videoElement.srcObject = null;
		}
	}

	function scanLoop() {
		if (!scanningActive || !videoElement || videoElement.paused || videoElement.ended) {
			return;
		}

		if (videoElement.readyState === videoElement.HAVE_ENOUGH_DATA) {
			if (!scanCanvas) {
				scanCanvas = document.createElement('canvas');
			}
			const canvas = scanCanvas;
			canvas.width = videoElement.videoWidth;
			canvas.height = videoElement.videoHeight;
			const ctx = canvas.getContext('2d');
			if (ctx) {
				ctx.drawImage(videoElement, 0, 0, canvas.width, canvas.height);
				const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
				if (jsQR) {
					const code = jsQR(imageData.data, imageData.width, imageData.height, {
						inversionAttempts: 'dontInvert'
					});
					if (code) {
						handleScannedQR(code.data);
						return;
					}
				}
			}
		}

		requestAnimationFrame(scanLoop);
	}

	async function handleScannedQR(scannedText: string) {
		scanningActive = false;
		stopCamera();

		const today = getTodayDateString();
		const expectedToken = `SMKN31-ATTENDANCE-KEY-${today}`;

		if (scannedText === expectedToken) {
			qrSuccess = true;
			qrScanning = false;
			addToast('QR Code terverifikasi!', 'success');

			const success = await submitKehadiran({
				status: 'hadir',
				tipe: `qr:${scannedText}`
			});

			if (success) {
				onload();
				setTimeout(() => {
					show = false;
				}, 1500);
			} else {
				show = false;
			}
		} else {
			addToast('QR Code tidak valid atau untuk hari yang berbeda!', 'error');
			setTimeout(() => {
				if (show) {
					startCamera();
				}
			}, 2000);
		}
	}

	$effect(() => {
		if (show) {
			startCamera();
		} else {
			scanningActive = false;
			stopCamera();
		}
		return () => {
			scanningActive = false;
			stopCamera();
		};
	});
</script>

<Modal
	{show}
	title=""
	onclose={() => {
		show = false;
	}}
>
	<div class="p-6 text-center space-y-6">
		<div>
			<h3 class="text-base font-bold text-slate-800">Scan QR Code Sekolah</h3>
			<p class="text-xs text-slate-400 mt-1">Arahkan kamera ke QR Code absensi di gerbang sekolah</p>
		</div>

		<div class="relative mx-auto aspect-square w-64 overflow-hidden rounded-2xl border-4 border-teal-500/30 bg-slate-900 flex flex-col items-center justify-center text-white">
			<!-- Video Stream -->
			{#if show && !cameraError && !qrSuccess}
				<video
					bind:this={videoElement}
					autoplay
					playsinline
					class="absolute inset-0 w-full h-full object-cover"
				>
					<track kind="captions" />
				</video>
			{/if}

			<!-- Simulated scanning line -->
			{#if qrScanning}
				<div class="absolute inset-x-0 top-0 h-1 bg-teal-400 shadow-[0_0_15px_#2dd4bf] animate-scan-line z-10"></div>
				<!-- Subtle frame border -->
				<div class="absolute inset-8 border-2 border-dashed border-teal-400/60 rounded-xl pointer-events-none z-10"></div>
				<p class="absolute bottom-4 left-0 right-0 text-center text-[10px] font-bold text-teal-300 bg-slate-900/80 py-1 px-2 mx-auto w-max rounded-md z-10">Mendeteksi QR Code...</p>
			{:else if qrSuccess}
				<div class="absolute inset-0 bg-slate-900/90 flex flex-col items-center justify-center z-20 animate-fade-in">
					<CheckCircle2 class="h-20 w-20 text-emerald-400 animate-bounce" />
					<p class="mt-4 text-xs font-black text-emerald-400">Scan Berhasil!</p>
				</div>
			{:else if cameraError}
				<div class="p-4 text-center space-y-2 z-10">
					<p class="text-xs text-rose-400 font-bold">{cameraError}</p>
					<p class="text-[10px] text-slate-500">Mencoba menggunakan simulasi...</p>
				</div>
			{/if}
		</div>

		<div class="text-xs text-slate-500 bg-slate-50 p-3.5 rounded-xl border border-slate-100 font-semibold leading-relaxed">
			{#if qrScanning}
				Memproses data enkripsi QR Code kehadiran harian...
			{:else}
				Absensi berhasil terekam di sistem SMKN 31 Jakarta.
			{/if}
		</div>
	</div>
</Modal>

<style>
	@keyframes scan-anim {
		0%, 100% { top: 0%; }
		50% { top: 100%; }
	}
	:global(.animate-scan-line) {
		position: absolute;
		animation: scan-anim 2s ease-in-out infinite;
	}
</style>
