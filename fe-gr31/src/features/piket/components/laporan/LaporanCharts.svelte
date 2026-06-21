<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import Chart from 'chart.js/auto';

	// Svelte 5 Props destructuring
	let {
		reportType,
		selectedBulan,
		selectedKelas,
		formattedWeekRange,
		weeklyChartData,
		monthlyChartData
	} = $props<{
		reportType: 'bulanan' | 'mingguan';
		selectedBulan: string;
		selectedKelas: string;
		formattedWeekRange: string;
		// Monthly chart: tren persentase kehadiran per kelas (dari backend)
		weeklyChartData: Array<{ label: string; val: number }>;
		monthlyChartData: Array<{ label: string; val: number }>;
	}>();

	let activeChartData = $derived(reportType === 'bulanan' ? monthlyChartData : weeklyChartData);

	let maxVal = $derived.by(() => {
		if (reportType === 'bulanan') return 100;
		const vals = activeChartData.map((d: { label: string; val: number }) => d.val);
		return vals.length > 0 ? Math.max(10, ...vals) : 10;
	});

	let chartWidth = $derived(Math.max(500, activeChartData.length * 45));

	let canvasEl = $state<HTMLCanvasElement | null>(null);
	let chartInstance: Chart | null = null;

	// Update or create chart instance
	$effect(() => {
		if (!canvasEl || activeChartData.length === 0) {
			if (chartInstance) {
				chartInstance.destroy();
				chartInstance = null;
			}
			return;
		}

		const ctx = canvasEl.getContext('2d');
		if (!ctx) return;

		// Create beautiful gradients
		// Green Gradient (>= 80)
		const gradientGreen = ctx.createLinearGradient(0, 0, 0, 220);
		gradientGreen.addColorStop(0, '#34d399'); // emerald-400
		gradientGreen.addColorStop(1, '#059669'); // emerald-600

		// Yellow Gradient (70 <= val < 80)
		const gradientYellow = ctx.createLinearGradient(0, 0, 0, 220);
		gradientYellow.addColorStop(0, '#fbbf24'); // amber-400
		gradientYellow.addColorStop(1, '#d97706'); // amber-600

		// Red/Rose Gradient (< 70)
		const gradientRed = ctx.createLinearGradient(0, 0, 0, 220);
		gradientRed.addColorStop(0, '#fb7185'); // rose-400
		gradientRed.addColorStop(1, '#e11d48'); // rose-600

		const labels = activeChartData.map((d: { label: string; val: number }) => d.label);
		const dataValues = activeChartData.map((d: { label: string; val: number }) => d.val);

		const backgroundColors = dataValues.map((val: number) => {
			if (reportType === 'bulanan') {
				if (val >= 80) return gradientGreen;
				if (val >= 70) return gradientYellow;
				return gradientRed;
			}
			return gradientGreen;
		});

		const borderColors = dataValues.map((val: number) => {
			if (reportType === 'bulanan') {
				if (val >= 80) return '#059669';
				if (val >= 70) return '#d97706';
				return '#e11d48';
			}
			return '#059669';
		});

		const datasetLabel = reportType === 'bulanan' ? 'Persentase Kehadiran' : 'Siswa Hadir';

		if (chartInstance) {
			// Update data and scales
			chartInstance.data.labels = labels;
			chartInstance.data.datasets[0].label = datasetLabel;
			chartInstance.data.datasets[0].data = dataValues;
			chartInstance.data.datasets[0].backgroundColor = backgroundColors;
			chartInstance.data.datasets[0].borderColor = borderColors;
			
			// Adjust y-axis max
			if (chartInstance.options.scales?.y) {
				chartInstance.options.scales.y.max = reportType === 'bulanan' ? 100 : undefined;
			}
			chartInstance.update();
		} else {
			chartInstance = new Chart(canvasEl, {
				type: 'bar',
				data: {
					labels: labels,
					datasets: [{
						label: datasetLabel,
						data: dataValues,
						backgroundColor: backgroundColors,
						borderColor: borderColors,
						borderWidth: 1,
						borderRadius: 5,
						borderSkipped: false,
						barPercentage: 0.55,
						categoryPercentage: 0.85
					}]
				},
				options: {
					responsive: true,
					maintainAspectRatio: false,
					plugins: {
						legend: {
							display: false
						},
						tooltip: {
							backgroundColor: 'rgba(15, 23, 42, 0.95)', // slate-900 with opacity
							titleColor: '#ffffff',
							bodyColor: '#ffffff',
							titleFont: {
								family: 'Plus Jakarta Sans, sans-serif',
								size: 11,
								weight: 'bold'
							},
							bodyFont: {
								family: 'Plus Jakarta Sans, sans-serif',
								size: 11
							},
							padding: 8,
							cornerRadius: 8,
							displayColors: false,
							callbacks: {
								label: function(context) {
									let val = context.parsed.y;
									return `${context.dataset.label}: ${val}${reportType === 'bulanan' ? '%' : ' siswa'}`;
								}
							}
						}
					},
					scales: {
						y: {
							min: 0,
							max: reportType === 'bulanan' ? 100 : undefined,
							grid: {
								color: '#f1f5f9'
							},
							border: {
								dash: [3, 3],
								color: '#e2e8f0'
							},
							ticks: {
								callback: function(value) {
									return value + (reportType === 'bulanan' ? '%' : '');
								},
								color: '#94a3b8',
								font: {
									family: 'Plus Jakarta Sans, sans-serif',
									size: 8,
									weight: 'bold'
								}
							}
						},
						x: {
							grid: {
								display: false
							},
							border: {
								color: '#e2e8f0'
							},
							ticks: {
								color: '#64748b',
								font: {
									family: 'Plus Jakarta Sans, sans-serif',
									size: 9,
									weight: 'bold'
								},
								autoSkip: false
							}
						}
					}
				}
			});
		}
	});

	onDestroy(() => {
		if (chartInstance) {
			chartInstance.destroy();
			chartInstance = null;
		}
	});

	function formatMonthYear(val: string) {
		if (!val) return '';
		const [year, month] = val.split('-');
		const monthNames = [
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
		const mIndex = parseInt(month, 10) - 1;
		return `${monthNames[mIndex] || month} ${year}`;
	}
</script>

<!-- Dynamic SVG Bar Chart Block -->
<div class="rounded-2xl border border-slate-100/80 bg-white p-5 text-left shadow-xs">
	<div class="mb-5 flex items-center justify-between border-b border-slate-100 pb-4">
		<div>
			<h3 class="text-xs font-bold tracking-wider text-slate-800 uppercase">
				{#if reportType === 'bulanan'}
					{#if selectedKelas}
						Persentase Kehadiran Kelas {selectedKelas}
					{:else}
						Persentase Kehadiran Per Kelas
					{/if}
				{:else}
					Grafik Harian Minggu Ini
				{/if}
			</h3>
			<p class="text-[10px] font-medium text-slate-400">
				{#if reportType === 'bulanan'}
					{#if selectedKelas}
						Persentase kehadiran kelas {selectedKelas} pada bulan {formatMonthYear(selectedBulan)} (sumber: rekap bulanan)
					{:else}
						Persentase kehadiran masing-masing kelas pada bulan {formatMonthYear(selectedBulan)} (sumber: rekap bulanan)
					{/if}
				{:else}
					Jumlah siswa hadir kelas {selectedKelas || 'Semua Kelas'} dari {formattedWeekRange}
				{/if}
			</p>
		</div>
	</div>

	<!-- Responsive Chart.js Bar Chart -->
	<div class="custom-scrollbar w-full overflow-x-auto py-2">
		{#if activeChartData.length === 0}
			<p class="py-8 text-center text-xs font-medium text-slate-400">Belum ada data untuk ditampilkan.</p>
		{:else}
			<div style="min-width: {chartWidth}px; height: 260px; position: relative;">
				<canvas bind:this={canvasEl}></canvas>
			</div>
		{/if}
	</div>
</div>

<style>
	/* Custom scrollbar styling for a clean sleek feel */
	.custom-scrollbar::-webkit-scrollbar {
		height: 5px;
	}
	.custom-scrollbar::-webkit-scrollbar-track {
		background: transparent;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: #e2e8f0;
		border-radius: 99px;
	}
	.custom-scrollbar::-webkit-scrollbar-thumb:hover {
		background: #cbd5e1;
	}
</style>
