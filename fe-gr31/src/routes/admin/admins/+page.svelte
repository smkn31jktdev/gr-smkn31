<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import AdminTable from '../../../features/admin/components/add/admin/AdminTable.svelte';
  import { currentUser } from '../../../stores/authStore';

  onMount(() => {
    // Extra safety double check
    if ($currentUser?.role !== 'super_admin') {
      goto('/admin', { replaceState: true });
    }
  });
</script>

<div class="space-y-6">
  <div>
    <h2 class="text-xl font-extrabold tracking-tight text-foreground font-display uppercase">Pusat Kelola Guru & Staf</h2>
    <p class="text-xs text-muted mt-0.5">Tambah dan hapus dewan guru BK, piket, wali kelas, atau super admin akses</p>
  </div>

  <AdminTable />
</div>
