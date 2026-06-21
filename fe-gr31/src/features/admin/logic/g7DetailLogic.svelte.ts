import { getG7RekapDetail, getG7Suggest, saveG7Rekap } from './adminLogic';
import { addToast } from '../../../stores/uiStore';
import { goto } from '$app/navigation';
import type { G7Rekap, SkorG7, G7SuggestResponse } from '../types/admin.types';

export class G7DetailState {
	nis: string;
	bulan: string;

	loading = $state(true);
	rekap = $state<G7Rekap | null>(null);
	suggest = $state<G7SuggestResponse | null>(null);
	isReadOnly = $derived(this.rekap?.status === 'final');

	skor = $state<SkorG7>({
		bangunPagi: 3,
		ibadahDoa: 3,
		ibadahSholatFajar: 0,
		ibadahSholat5Waktu: 0,
		ibadahZikir: 0,
		ibadahDhuha: 0,
		ibadahRowatib: 0,
		ibadahTarawih: 0,
		ibadahPuasa: 0,
		ibadahZakat: 3,
		olahraga: 3,
		makanSehat: 3,
		belajarKitabSuci: 3,
		belajarBukuUmum: 3,
		belajarBukuMapel: 3,
		belajarTugas: 3,
		bermasyarakat: 3,
		tidurCepat: 3
	});

	waliKelas = $state('');
	guruBK = $state('');
	status = $state<'draft' | 'reviewed' | 'final'>('draft');

	constructor(nis: string, bulan: string) {
		this.nis = nis;
		this.bulan = bulan;
	}

	async loadDetail() {
		this.loading = true;
		const data = await getG7RekapDetail(this.nis, this.bulan);
		const advisory = await getG7Suggest(this.nis, this.bulan);

		if (data) {
			this.rekap = data;
			this.skor = { ...data.skor };
			this.waliKelas = data.waliKelas || '';
			this.guruBK = data.guruBK || '';
			this.status = data.status || 'draft';
		}
		if (advisory) {
			this.suggest = advisory;
		}
		this.loading = false;
	}

	async handleSave(handlers: { resolve: () => void; reject: () => void }) {
		if (!this.rekap) {
			handlers.reject();
			return;
		}

		// Validation checks for finalization
		if (this.status === 'final') {
			if (!this.waliKelas.trim() || !this.guruBK.trim()) {
				addToast(
					'Finalisasi memerlukan nama 2 penilai lengkap (Guru Wali dan BK)',
					'warning'
				);
				handlers.reject();
				return;
			}
		}

		const payload: Partial<G7Rekap> = {
			nis: this.nis,
			bulanTahun: this.bulan,
			skor: this.skor,
			waliKelas: this.waliKelas,
			guruBK: this.guruBK,
			status: this.status
		};

		const success = await saveG7Rekap(payload);
		if (success) {
			handlers.resolve();
			goto('/admin/g7');
		} else {
			handlers.reject();
		}
	}
}
