import { get, writable } from 'svelte/store';
import { listStudents, listG7Rekap, getG7Statistik, getG7RekapSemester, getG7RekapKelas } from './adminLogic';
import { currentUser } from '../../../stores/authStore';
import { addToast } from '../../../stores/uiStore';
import { BASE, apiRequest } from '../../../api/client';

// Reactive state for the admin dashboard
export const dashboardLoading = writable<boolean>(false);
export const activeAdminName = writable<string>('Super Admin 31');
export const activeAdminEmail = writable<string>('smkn31jktdev@gmail.com');
export const totalSiswaCount = writable<number>(607);
export const pendingWalasAduan = writable<any[]>([]);

// Left Panel: Laporan Bulanan
export const selectedMonth = writable<string>('2026-06');
export const laporanSearchQuery = writable<string>('');
export const G7ReportsList = writable<any[]>([]);
export const G7ReportsCount = writable<number>(0);

// Left Panel: Laporan Semester
export const selectedSemester = writable<string>('2025/2026-genap');
export const semesterSearchQuery = writable<string>('');
export const G7SemesterReportsList = writable<any[]>([]);
export const G7SemesterReportsCount = writable<number>(0);

// Right Panel: List Siswa
export const studentSearchQuery = writable<string>('');
export const studentsList = writable<any[]>([]);
export const totalStudentsFiltered = writable<number>(0);

// Mapping of semester values to month lists
const semesterMonthsMap: Record<string, string[]> = {
  '2025/2026-genap': ['2026-01', '2026-02', '2026-03', '2026-04', '2026-05', '2026-06'],
  '2025/2026-ganjil': ['2025-07', '2025-08', '2025-09', '2025-10', '2025-11', '2025-12']
};

// Predicate helper
function getPredikat(akhir: number): string {
  if (akhir >= 90) return 'Istimewa';
  if (akhir >= 80) return 'Sangat Baik';
  if (akhir >= 70) return 'Baik';
  if (akhir >= 60) return 'Cukup';
  return 'Kurang';
}

// Helper to format month to Indonesian label like "Juni 2026"
export function formatMonthIndonesian(monthStr: string): string {
  const [year, month] = monthStr.split('-');
  const months = [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ];
  const idx = parseInt(month, 10) - 1;
  return `${months[idx] || month} ${year}`;
}

// Fetch pending complaints for the logged-in homeroom teacher (Walas)
export async function loadPendingAduanForWalas() {
  const user = get(currentUser);
  if (!user) return;

  const isWalasUser = user.role === 'walas' || 
                      user.role === 'guru_wali' || 
                      user.is_walas === true || 
                      user.isWalas === true;

  if (isWalasUser && user.nama) {
    try {
      const { data: aduanData } = await apiRequest<any>('/v1/admin/aduan');
      const items = (aduanData?.items || []) as any[];

      const nowMs = new Date().getTime();
      const sevenDaysMs = 7 * 24 * 60 * 60 * 1000;

      const filteredAduan = items.filter(a => {
        const isNotClosed = a.status !== 'closed';
        const matchesWalas = (a.walas && a.walas.toLowerCase() === user.nama.toLowerCase()) ||
                             (a.kelas && user.kelas && a.kelas.toLowerCase() === user.kelas.toLowerCase());

        if (!isNotClosed || !matchesWalas) return false;

        const createdAtMs = new Date(a.createdAt).getTime();
        const ageMs = nowMs - createdAtMs;
        const isWithin7Days = ageMs >= 0 && ageMs <= sevenDaysMs;

        if (!isWithin7Days) return false;

        // Check if the latest message was sent by student (meaning it needs response from admin)
        const hasMessages = a.messages && a.messages.length > 0;
        const lastMessageFromStudent = hasMessages && a.messages[a.messages.length - 1].role === 'student';

        return lastMessageFromStudent;
      });

      pendingWalasAduan.set(filteredAduan);
    } catch (e) {
      console.error('Error loading pending walas aduan:', e);
    }
  } else {
    pendingWalasAduan.set([]);
  }
}

// Load all initial data for the dashboard
export async function loadDashboardData() {
  dashboardLoading.set(true);
  try {
    const user = get(currentUser);
    if (user) {
      activeAdminName.set(user.nama || 'Super Admin 31');
      activeAdminEmail.set(user.email || 'smkn31jktdev@gmail.com');
    }

    const monthVal = get(selectedMonth);

    // 1. Get G7 Stats (for total students and count of graded)
    const stats = await getG7Statistik(monthVal);
    if (stats) {
      totalSiswaCount.set(stats.totalSiswa || 607);
      G7ReportsCount.set(stats.sudahDinilai || 421);
    }

    // 2. Fetch G7 Reports (Laporan Bulanan) using Rekap Kelas Lengkap
    const rekapKelas = await getG7RekapKelas(monthVal, "");
    if (rekapKelas && rekapKelas.siswa) {
      const query = get(laporanSearchQuery).toLowerCase().trim();
      const mappedSiswa = rekapKelas.siswa.map(item => ({
        ...item,
        bulanTahun: monthVal
      }));
      const filtered = query
        ? mappedSiswa.filter(item =>
            item.namaSiswa.toLowerCase().includes(query) ||
            item.nis.toLowerCase().includes(query) ||
            item.kelas.toLowerCase().includes(query)
          )
        : mappedSiswa;
      G7ReportsList.set(filtered);
      if (rekapKelas.sudahDinilai !== undefined) {
        G7ReportsCount.set(rekapKelas.sudahDinilai);
      }
      if (rekapKelas.totalSiswa !== undefined) {
        totalSiswaCount.set(rekapKelas.totalSiswa);
      }
    } else {
      G7ReportsList.set([]);
    }

    // 3. Fetch Students List (for right panel)
    const studentRes = await listStudents(1, 100, get(studentSearchQuery));
    studentsList.set(studentRes.items || []);
    totalStudentsFiltered.set(studentRes.total || 0);

    // 4. Load pending aduan for walas
    await loadPendingAduanForWalas();

  } catch (error) {
    console.error('Error loading dashboard data:', error);
    addToast('Gagal memuat data dashboard', 'error');
  } finally {
    dashboardLoading.set(false);
  }
}

// Load semester reports (from single backend aggregated query)
export async function loadSemesterData() {
  dashboardLoading.set(true);
  try {
    const user = get(currentUser);
    if (user) {
      activeAdminName.set(user.nama || 'Super Admin 31');
      activeAdminEmail.set(user.email || 'smkn31jktdev@gmail.com');
    }

    const semVal = get(selectedSemester);
    
    // Fetch aggregated semester data from single endpoint
    const list = await getG7RekapSemester(semVal) || [];

    // Filter by search query
    const query = get(semesterSearchQuery).toLowerCase().trim();
    const filteredList = query
      ? list.filter(item =>
          item.namaSiswa.toLowerCase().includes(query) ||
          item.nis.toLowerCase().includes(query) ||
          item.kelas.toLowerCase().includes(query)
        )
      : list;

    G7SemesterReportsList.set(filteredList);
    G7SemesterReportsCount.set(filteredList.length);

  } catch (error) {
    console.error('Error loading semester data:', error);
    addToast('Gagal memuat data semester', 'error');
  } finally {
    dashboardLoading.set(false);
  }
}

// Trigger specific student report PDF download
export async function downloadStudentReportPDF(nis: string, nama: string, bulan: string) {
  addToast(`Mengunduh Laporan PDF ${nama}...`, 'info');
  
  try {
    const isBrowser = typeof window !== 'undefined';
    const token = isBrowser
      ? (localStorage.getItem('adminToken') ?? localStorage.getItem('studentToken'))
      : null;

    const response = await fetch(`${BASE}/v1/admin/g7/rekap/${nis}/${bulan}/pdf`, {
      method: 'GET',
      headers: {
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const htmlText = await response.text();

    const html2pdf = await getHtml2Pdf();
    if (!html2pdf) {
      throw new Error('html2pdf failed to load');
    }

    // Parse the HTML text into a DOM node to ensure html2pdf handles it properly and triggers direct download
    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = htmlText;

    const opt = {
      margin:       [10, 10, 10, 10], // in mm
      filename:     `Laporan_G7_${nama.replace(/\s+/g, '_')}_${bulan}.pdf`,
      image:        { type: 'jpeg', quality: 0.98 },
      html2canvas:  { scale: 2, useCORS: true },
      jsPDF:        { unit: 'mm', format: 'a4', orientation: 'portrait' }
    };

    await html2pdf().set(opt).from(tempDiv).save();
    addToast(`Laporan PDF ${nama} berhasil diunduh`, 'success');
  } catch (err) {
    console.error(err);
    addToast('Gagal mengunduh laporan PDF', 'error');
  }
}

// Trigger specific student semester report PDF download
export async function downloadStudentSemesterReportPDF(nis: string, nama: string, semester: string, avgScore: number, predikat: string) {
  const semesterLabel = semester === '2025/2026-genap' ? 'Genap 2025/2026' : 'Ganjil 2025/2026';
  addToast(`Mengunduh PDF Laporan Semester ${nama} - ${semesterLabel}...`, 'info');
  
  try {
    const html2pdf = await getHtml2Pdf();
    if (!html2pdf) {
      throw new Error('html2pdf failed to load');
    }

    const htmlContent = `
      <div style="font-family: Arial, sans-serif; padding: 30px; color: #000; line-height: 1.6; max-width: 800px; margin: auto;">
        <div style="text-align: center; margin-bottom: 40px; border-bottom: 2px solid #003399; padding-bottom: 20px;">
          <h2 style="margin: 0; color: #003399; font-size: 16pt; font-weight: bold;">LAPORAN REKAPITULASI SEMESTER G7</h2>
          <h3 style="margin: 5px 0 0 0; color: #003399; font-size: 13pt;">SMK NEGERI 31 JAKARTA</h3>
          <p style="margin: 5px 0 0 0; font-size: 10pt; color: #555; font-weight: bold;">TAHUN AJARAN 2025/2026</p>
        </div>
        
        <table style="width: 100%; margin-bottom: 40px; border-collapse: collapse; font-size: 11pt;">
          <tr>
            <td style="width: 200px; font-weight: bold; padding: 8px 0; color: #334155;">Nama Siswa</td>
            <td style="padding: 8px 0; font-weight: bold;">: ${nama}</td>
          </tr>
          <tr>
            <td style="font-weight: bold; padding: 8px 0; color: #334155;">NIS / NISN</td>
            <td style="padding: 8px 0; font-weight: bold;">: ${nis}</td>
          </tr>
          <tr>
            <td style="font-weight: bold; padding: 8px 0; color: #334155;">Semester</td>
            <td style="padding: 8px 0;">: ${semesterLabel}</td>
          </tr>
          <tr>
            <td style="font-weight: bold; padding: 8px 0; color: #334155;">Rata-rata Nilai Akhir G7</td>
            <td style="padding: 8px 0; font-weight: bold; color: #003399; font-size: 13pt;">: ${avgScore}</td>
          </tr>
          <tr>
            <td style="font-weight: bold; padding: 8px 0; color: #334155;">Predikat</td>
            <td style="padding: 8px 0;">: <span style="background-color: #f1f5f9; padding: 4px 10px; border-radius: 6px; font-weight: bold; color: #334155; border: 1px solid #cbd5e1;">${predikat}</span></td>
          </tr>
          <tr>
            <td style="font-weight: bold; padding: 8px 0; color: #334155;">Status Verifikasi</td>
            <td style="padding: 8px 0; font-weight: bold; color: #10b981;">: FINAL SEMESTER</td>
          </tr>
        </table>
        
        <div style="margin-top: 60px; border-top: 1px dashed #cbd5e1; padding-top: 20px; font-size: 10pt; color: #64748b; text-align: center;">
          Dokumen ini sah dikeluarkan secara otomatis oleh Sistem Monitoring Karakter (G7) SMKN 31 Jakarta.
        </div>
      </div>
    `;

    const opt = {
      margin:       [15, 15, 15, 15],
      filename:     `Laporan_Semester_G7_${nama.replace(/\s+/g, '_')}_${semester}.pdf`,
      image:        { type: 'jpeg', quality: 0.98 },
      html2canvas:  { scale: 2, useCORS: true },
      jsPDF:        { unit: 'mm', format: 'a4', orientation: 'portrait' }
    };

    const tempDiv = document.createElement('div');
    tempDiv.innerHTML = htmlContent;

    await html2pdf().set(opt).from(tempDiv).save();
    addToast(`Laporan Semester ${nama} berhasil diunduh`, 'success');
  } catch (err) {
    console.error(err);
    addToast('Gagal mengunduh laporan semester PDF', 'error');
  }
}

// Helper to dynamically load html2pdf from CDN
async function getHtml2Pdf(): Promise<any> {
  if (typeof window === 'undefined') return null;
  if ((window as any).html2pdf) {
    return (window as any).html2pdf;
  }
  
  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = 'https://cdnjs.cloudflare.com/ajax/libs/html2pdf.js/0.10.1/html2pdf.bundle.min.js';
    script.onload = () => resolve((window as any).html2pdf);
    script.onerror = (err) => reject(err);
    document.head.appendChild(script);
  });
}

// Helper to dynamically load JSZip from CDN
async function getJSZip(): Promise<any> {
  if (typeof window === 'undefined') return null;
  if ((window as any).JSZip) {
    return (window as any).JSZip;
  }
  
  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = 'https://cdnjs.cloudflare.com/ajax/libs/jszip/3.10.1/jszip.min.js';
    script.onload = () => resolve((window as any).JSZip);
    script.onerror = (err) => reject(err);
    document.head.appendChild(script);
  });
}

// Trigger download all student reports
export async function downloadAllReportsPDF(bulan: string, count: number) {
  addToast(`Mempersiapkan unduhan masal ${count} laporan untuk bulan ${formatMonthIndonesian(bulan)}...`, 'info');
  
  try {
    const JSZip = await getJSZip();
    if (!JSZip) {
      throw new Error('JSZip failed to load');
    }
    
    const zip = new JSZip();
    const list = get(G7ReportsList);
    const token = typeof window !== 'undefined'
      ? (localStorage.getItem('adminToken') ?? localStorage.getItem('studentToken'))
      : null;
      
    let successCount = 0;
    
    for (const report of list) {
      try {
        const response = await fetch(`${BASE}/v1/admin/g7/rekap/${report.nis}/${bulan}/pdf`, {
          method: 'GET',
          headers: {
            ...(token ? { Authorization: `Bearer ${token}` } : {})
          }
        });
        if (response.ok) {
          const htmlText = await response.text();
          const fileName = `Laporan_G7_${report.namaSiswa.replace(/\s+/g, '_')}_${bulan}.html`;
          zip.file(fileName, htmlText);
          successCount++;
        }
      } catch (err) {
        console.error(`Gagal mengunduh laporan untuk ${report.namaSiswa}:`, err);
      }
    }
    
    if (successCount === 0) {
      addToast('Tidak ada laporan yang berhasil diunduh', 'error');
      return;
    }
    
    const content = await zip.generateAsync({ type: 'blob' });
    const url = URL.createObjectURL(content);
    const link = document.createElement('a');
    link.href = url;
    link.download = `Bulk_Laporan_G7_${bulan}.zip`;
    link.click();
    URL.revokeObjectURL(url);
    addToast(`Berhasil mengunduh ${successCount} laporan dalam format ZIP`, 'success');
  } catch (err) {
    console.error(err);
    addToast('Gagal memproses unduhan masal ZIP', 'error');
  }
}

// Trigger download all student semester reports
export async function downloadAllSemesterReportsPDF(semester: string, count: number) {
  const semesterLabel = semester === '2025/2026-genap' ? 'Genap 2025/2026' : 'Ganjil 2025/2026';
  addToast(`Mempersiapkan unduhan masal ${count} laporan semester untuk ${semesterLabel}...`, 'info');
  
  try {
    const JSZip = await getJSZip();
    if (!JSZip) {
      throw new Error('JSZip failed to load');
    }
    
    const zip = new JSZip();
    const list = get(G7SemesterReportsList);
    
    for (const report of list) {
      const txtContent = `LAPORAN REKAPITULASI SEMESTER G7 - SMKN 31 JAKARTA\n` +
        `Nama Siswa: ${report.namaSiswa}\n` +
        `NIS: ${report.nis}\n` +
        `Kelas: ${report.kelas}\n` +
        `Semester: ${semesterLabel}\n` +
        `Rata-rata Nilai Akhir G7: ${report.nilaiAkhir}\n` +
        `Predikat: ${report.predikat}\n` +
        `Total Bulan Tercatat: ${report.monthsCount}\n` +
        `Status Verifikasi: FINAL SEMESTER\n` +
        `-----------------------------------------\n` +
        `Dokumen ini sah dikeluarkan oleh Sistem Monitoring Karakter SMKN 31.`;
      
      const fileName = `Laporan_Semester_G7_${report.namaSiswa.replace(/\s+/g, '_')}_${semester}.txt`;
      zip.file(fileName, txtContent);
    }
    
    const content = await zip.generateAsync({ type: 'blob' });
    const url = URL.createObjectURL(content);
    const link = document.createElement('a');
    link.href = url;
    link.download = `Bulk_Laporan_Semester_G7_${semester}.zip`;
    link.click();
    URL.revokeObjectURL(url);
    addToast(`Berhasil mengunduh ${list.length} laporan semester dalam format ZIP`, 'success');
  } catch (err) {
    console.error(err);
    addToast('Gagal memproses unduhan masal semester ZIP', 'error');
  }
}

// Trigger Manual Book PDF download (as TXT)
export function downloadManualBookPDF() {
  addToast('Mengunduh Manual Book Dashboard Guru Wali...', 'info');
  
  setTimeout(() => {
    try {
      const dummyContent = `MANUAL BOOK PANDUAN PENGGUNAAN DASHBOARD GURU WALI\n` +
        `SMK NEGERI 31 JAKARTA\n` +
        `-----------------------------------------\n` +
        `1. Cara login akun admin\n` +
        `2. Verifikasi jurnal harian siswa\n` +
        `3. Pengisian rekap bulanan G7\n` +
        `4. Cetak dan export data rekap`;
      
      const blob = new Blob([dummyContent], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `Manual_Book_Dashboard_Guru_Wali.txt`;
      link.click();
      URL.revokeObjectURL(url);
      addToast('Manual book berhasil diunduh', 'success');
    } catch (err) {
      addToast('Gagal mengunduh manual book', 'error');
    }
  }, 1000);
}
