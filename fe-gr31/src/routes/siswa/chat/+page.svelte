<script lang="ts">
  import StudentChatRoom from '../../../features/student/components/chat/StudentChatRoom.svelte';
</script>

<div class="space-y-6">
  <div>
    <h2 class="text-xl font-extrabold tracking-tight text-foreground font-display uppercase">Konsultasi BK & Walas</h2>
    <p class="text-xs text-muted mt-0.5">Ruang konsultasi pribadi untuk bimbingan konseling dan aduan secara rahasia</p>
  </div>

  <!-- Counseling room view -->
  <StudentChatRoom />
</div>
