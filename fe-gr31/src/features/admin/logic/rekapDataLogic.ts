import { listG7Rekap } from './adminLogic';
import { exportToCSV } from './exportLogic';
import type { G7Rekap } from '../types/admin.types';

export interface RekapSummaryStats {
  terpilihCount: number;
  totalArsipCount: number;
}

// Fetch selected period count and total count across all periods
export async function loadRekapSummary(bulan: string, kelas = ''): Promise<RekapSummaryStats> {
  try {
    const selectedRes = await listG7Rekap({ bulan, kelas }, 1, 1);
    const totalRes = await listG7Rekap({}, 1, 1);

    return {
      terpilihCount: selectedRes.total || 0,
      totalArsipCount: totalRes.total || 0
    };
  } catch (error) {
    console.error('Error loading rekap summary stats:', error);
    return { terpilihCount: 0, totalArsipCount: 0 };
  }
}

// Fetch top 5 preview records to minimize query times
export async function loadRekapPreview(bulan: string, kelas = ''): Promise<G7Rekap[]> {
  try {
    const res = await listG7Rekap({ bulan, kelas }, 1, 5);
    return res.items || [];
  } catch (error) {
    console.error('Error loading rekap preview data:', error);
    return [];
  }
}

/**
 * Fetch all records for CSV export by paging through limits of 100
 */
export async function fetchAllRekapForExport(bulan: string, kelas = ''): Promise<G7Rekap[]> {
  let allItems: G7Rekap[] = [];
  let page = 1;
  const limit = 100;
  let hasMore = true;

  try {
    while (hasMore) {
      const res = await listG7Rekap({ bulan, kelas }, page, limit);
      if (res.items && res.items.length > 0) {
        allItems = [...allItems, ...res.items];
      }
      hasMore = res.hasMore;
      page++;
      // Guard to prevent infinite looping
      if (page > 50) break;
    }
  } catch (error) {
    console.error('Error fetching all records for export:', error);
  }

  return allItems;
}

/**
 * Process array of G7Rekap items and export to CSV mapping the 18 sub-indicators (omitting Ramadan-only indicators outside Ramadan)
 */
export function handleExportG7RekapCSV(items: G7Rekap[], selectedBulan: string, selectedKelas = '', isRamadan = false) {
  const headers = [
    'NIS', 'Nama Siswa', 'Kelas', 'Agama', 'Periode', 'Hari Tercatat',
    '1. Bangun Pagi', 
    '2a. Berdoa Sendiri & Ortu', 
    '2b. Sholat Fajar (Islam)', 
    '2c. Sholat 5 Waktu Berjamaah (Islam)', 
    '2d. Zikir & Doa (Islam)', 
    '2e. Sholat Dhuha (Islam)'
  ];

  if (isRamadan) {
    headers.push(
      '2f. Sholat Sunah Rowatib (Islam)', 
      '2g. Sholat Tarawih & Witir (Islam)', 
      '2h. Puasa Ramadhan (Islam)'
    );
  }

  headers.push(
    '2i. Zakat / Infaq / Sodaqoh', 
    '3. Berolahraga', 
    '4. Makan Sehat & Bergizi', 
    '5a. Membaca Kitab Suci', 
    '5b. Membaca Buku Umum/Novel/Hobi', 
    '5c. Membaca Buku Mapel', 
    '5d. Mengerjakan Tugas/PR', 
    '6. Bermasyarakat', 
    '7. Tidur Cepat',
    'Nilai Perolehan', 'Nilai Maks', 'Nilai Akhir ((Perolehan/Nilai Maks)x100)'
  );

  const rows = items.map(item => {
    const row = [
      item.nis || '',
      item.namaSiswa || '',
      item.kelas || '',
      item.agama || 'islam',
      item.bulanTahun || '',
      item.hariTercatat ?? 0,
      item.skor?.bangunPagi ?? 0,
      item.skor?.ibadahDoa ?? 0,
      item.skor?.ibadahSholatFajar ?? 0,
      item.skor?.ibadahSholat5Waktu ?? 0,
      item.skor?.ibadahZikir ?? 0,
      item.skor?.ibadahDhuha ?? 0
    ];

    if (isRamadan) {
      row.push(
        item.skor?.ibadahRowatib ?? 0,
        item.skor?.ibadahTarawih ?? 0,
        item.skor?.ibadahPuasa ?? 0
      );
    }

    row.push(
      item.skor?.ibadahZakat ?? 0,
      item.skor?.olahraga ?? 0,
      item.skor?.makanSehat ?? 0,
      item.skor?.belajarKitabSuci ?? 0,
      item.skor?.belajarBukuUmum ?? 0,
      item.skor?.belajarBukuMapel ?? 0,
      item.skor?.belajarTugas ?? 0,
      item.skor?.bermasyarakat ?? 0,
      item.skor?.tidurCepat ?? 0,
      item.nilaiPerolehan ?? 0,
      item.nilaiMaks ?? 90,
      Math.round(item.nilaiAkhir ?? 0)
    );

    return row;
  });

  const classSuffix = selectedKelas ? `_${selectedKelas.replace(/\s+/g, '')}` : '';
  const filename = `Rekap_G7_${selectedBulan}${classSuffix}`;
  exportToCSV(headers, rows, filename);
}
