export interface Siswa {
  id: string;
  nis: string;
  nama: string;
  kelas: string;
  walas: string;      // Guru Wali (Permanent tutor, mapped to s.wali_id)
  waliKelas?: string;  // Walas (Dynamic class advisor, matching teacher of the class)
  agama?: string;
  email?: string;
  isOnline: boolean;
  fotoProfil?: string;
}

export interface Admin {
  id: string;
  nama: string;
  email: string;
  role: 'super_admin' | 'guru_bk' | 'admin_bk' | 'guru_wali' | 'walas' | 'admin' | 'piket' | 'admin_piket';
  isWalas?: boolean;
  kelas?: string;
  fotoProfil?: string;
}

export interface SiswaLoginResponse {
  accessToken: string;
  refreshToken: string;
  siswa: Siswa;
}

export interface AdminLoginResponse {
  accessToken: string;
  refreshToken: string;
  admin: Admin;
}
