<script lang="ts">
  import { page } from '$app/stores';
  import HabitsMonitoring from '../../../../features/admin/components/monitoring/HabitsMonitoring.svelte';

  let kegiatan = $derived($page.params.kegiatan ?? 'bangun-pagi');
</script>

<HabitsMonitoring {kegiatan} />
