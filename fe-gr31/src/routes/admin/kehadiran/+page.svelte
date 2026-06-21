<script lang="ts">
  import ProtectedRoute from '../../../features/shared/components/ProtectedRoute.svelte';
  import AdminKehadiranSiswa from '../../../features/admin/components/kehadiran-siswa/AdminKehadiranSiswa.svelte';
</script>

<ProtectedRoute allowedRoles={['piket', 'admin_piket']}>
  <AdminKehadiranSiswa />
</ProtectedRoute>
