import { apiRequest } from '../../../api/client';
import type { G7Jurnal } from '../../student/types/student.types';

export interface HabitMetadata {
  title: string;
  subtitle: string;
  dbKey: 'bangun' | 'ibadah' | 'makan' | 'olahraga' | 'belajar' | 'bermasyarakat' | 'tidur';
}

/**
 * Fetch daily G7 student journals with query and pagination filters
 */
export async function listG7DailyJournals(filters: {
  tanggal?: string;
  kelas?: string;
  q?: string;
  page?: number;
  limit?: number;
  nis?: string;
  nisn?: string;
}) {
  const query = new URLSearchParams();
  if (filters.tanggal) query.append('tanggal', filters.tanggal);
  if (filters.kelas) query.append('kelas', filters.kelas);
  if (filters.q) query.append('q', filters.q);
  if (filters.page) query.append('page', String(filters.page));
  if (filters.limit) query.append('limit', String(filters.limit));
  const nisVal = filters.nis || filters.nisn;
  if (nisVal) query.append('nisn', nisVal);

  const { data, error } = await apiRequest<{ items: G7Jurnal[]; total: number; hasMore: boolean }>(
    `/v1/admin/g7?${query.toString()}`
  );
  if (error) return { items: [] as G7Jurnal[], total: 0, hasMore: false };
  return {
    items: data?.items || [],
    total: data?.total || 0,
    hasMore: data?.hasMore || false
  };
}

// Maps route parameter key to localized title, subtitle and DB key
export function getHabitMetadata(habitKey: string): HabitMetadata {
  switch (habitKey) {
    case 'bangun-pagi':
      return {
        title: 'Bangun Pagi',
        subtitle: 'Rekapitulasi aktivitas bangun pagi siswa.',
        dbKey: 'bangun'
      };
    case 'beribadah':
      return {
        title: 'Beribadah',
        subtitle: 'Rekapitulasi pembiasaan ibadah harian siswa.',
        dbKey: 'ibadah'
      };
    case 'makan-sehat':
      return {
        title: 'Makan Sehat',
        subtitle: 'Rekapitulasi pola makan sehat dan bergizi siswa.',
        dbKey: 'makan'
      };
    case 'olahraga':
      return {
        title: 'Berolahraga',
        subtitle: 'Rekapitulasi pelaksanaan aktivitas olahraga siswa.',
        dbKey: 'olahraga'
      };
    case 'belajar':
      return {
        title: 'Gemar Belajar',
        subtitle: 'Rekapitulasi kebiasaan membaca dan belajar mandiri siswa.',
        dbKey: 'belajar'
      };
    case 'bermasyarakat':
      return {
        title: 'Bermasyarakat',
        subtitle: 'Rekapitulasi kegiatan sosial kemasyarakatan siswa.',
        dbKey: 'bermasyarakat'
      };
    case 'tidur-cukup':
      return {
        title: 'Tidur Cepat & Cukup',
        subtitle: 'Rekapitulasi pola tidur awal sebelum pukul 22:00.',
        dbKey: 'tidur'
      };
    default:
      return {
        title: 'Bangun Pagi',
        subtitle: 'Rekapitulasi aktivitas bangun pagi siswa.',
        dbKey: 'bangun'
      };
  }
}

// Parses raw waktu/keterangan string from AstraDB into dynamic structured properties

export function parseHabitDetails(item: G7Jurnal, dbKey: string) {
  const act = item[dbKey as keyof G7Jurnal] as any;
  if (!act) {
    return { done: false, display: '-' };
  }

  const done = act.done || false;
  const waktu = act.waktu || '-';
  const ket = act.keterangan || '';

  switch (dbKey) {
    case 'bangun':
    case 'tidur': {
      const doa = ket.includes('Membaca Doa: Ya') ? 'Ya' : 'Tidak';
      return {
        done,
        waktu,
        doa,
        display: done ? `${waktu} (Berdoa: ${doa})` : 'Belum mengisi'
      };
    }
    case 'ibadah': {
      const parts = [];
      if (ket.includes('Berdoa')) parts.push('Berdoa Ortu');
      if (ket.includes('Sholat Fajar')) parts.push('Sholat Fajar');
      if (ket.includes('Sholat 5 Waktu')) parts.push('Sholat 5 Waktu');
      if (ket.includes('Zikir')) parts.push('Zikir');
      if (ket.includes('Sholat Dhuha')) parts.push('Dhuha');
      if (ket.includes('Sholat Sunah Rowatib')) parts.push('Rowatib');
      const infaqMatch = ket.match(/Infaq: Rp ([\d.]+)/);
      if (infaqMatch) parts.push(`Infaq: Rp ${infaqMatch[1]}`);

      return {
        done,
        doaOrtu: ket.includes('Berdoa'),
        sholatFajar: ket.includes('Sholat Fajar'),
        sholat5Waktu: ket.includes('Sholat 5 Waktu'),
        zikir: ket.includes('Zikir'),
        dhuha: ket.includes('Sholat Dhuha'),
        rowatib: ket.includes('Sholat Sunah Rowatib'),
        infaq: infaqMatch ? `Rp ${infaqMatch[1]}` : '-',
        display: parts.length > 0 ? parts.join(', ') : (done ? 'Melakukan ibadah' : 'Belum mengisi')
      };
    }
    case 'makan': {
      const utMatch = ket.match(/Makanan Utama: ([^,]+)/);
      const lMatch = ket.match(/Lauk: ([^,]+)/);
      const sayur = ket.includes('Sayur/Buah: Ya') ? 'Ya' : 'Tidak';
      const susu = ket.includes('Susu/Suplemen: Ya') ? 'Ya' : 'Tidak';

      return {
        done,
        utama: utMatch ? utMatch[1] : '-',
        lauk: lMatch ? lMatch[1] : '-',
        sayurBuah: sayur,
        susuSuplemen: susu,
        display: done ? `Utama: ${utMatch ? utMatch[1] : '-'}, Lauk: ${lMatch ? lMatch[1] : '-'}` : 'Belum mengisi'
      };
    }
    case 'olahraga': {
      const aktMatch = ket.match(/Aktivitas: ([^,]+)/);
      const durMatch = ket.match(/Durasi: (\d+) menit/);
      const aktivitas = aktMatch ? aktMatch[1] : '-';
      const durasi = durMatch ? `${durMatch[1]} menit` : '-';

      return {
        done,
        aktivitas,
        durasi,
        display: done ? `${aktivitas} (${durasi})` : 'Belum mengisi'
      };
    }
    case 'belajar': {
      const parts = [];
      if (ket.includes('Kitab Suci')) parts.push('Baca Kitab Suci');
      if (ket.includes('Buku Umum')) parts.push('Baca Buku Umum');
      if (ket.includes('Buku Pelajaran')) parts.push('Buku Pelajaran');
      if (ket.includes('Tugas/PR')) parts.push('Tugas/PR');

      return {
        done,
        kitabSuci: ket.includes('Kitab Suci'),
        bukuUmum: ket.includes('Buku Umum'),
        bukuMapel: ket.includes('Buku Pelajaran'),
        tugasPR: ket.includes('Tugas/PR'),
        display: parts.length > 0 ? parts.join(', ') : (done ? 'Belajar Mandiri' : 'Belum mengisi')
      };
    }
    case 'bermasyarakat': {
      const kegMatch = ket.match(/Kegiatan: ([^,]+)/);
      const lokMatch = ket.match(/Lokasi: ([^,]+)/);
      const ot = ket.includes('Diketahui OT/RT: Ya') ? 'Ya' : 'Tidak';
      const kegiatan = kegMatch ? kegMatch[1] : '-';
      const lokasi = lokMatch ? lokMatch[1] : '-';

      return {
        done,
        kegiatan,
        lokasi,
        diketahuiOT: ot,
        display: done ? `${kegiatan} di ${lokasi}` : 'Belum mengisi'
      };
    }
    default:
      return { done, display: done ? 'Selesai' : 'Belum mengisi' };
  }
}
