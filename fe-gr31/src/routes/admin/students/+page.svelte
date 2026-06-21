<script lang="ts">
  import StudentTable from '../../../features/admin/components/add/student/StudentTable.svelte';
</script>

<div class="space-y-6">
  <div>
    <h2 class="text-xl font-extrabold tracking-tight text-foreground font-display uppercase">Pusat Kelola Siswa</h2>
    <p class="text-xs text-muted mt-0.5">Tambah, edit, hapus, dan lakukan bulk import data profile siswa sekolah</p>
  </div>

  <StudentTable />
</div>
