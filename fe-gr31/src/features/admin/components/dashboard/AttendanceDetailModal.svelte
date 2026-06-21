<script lang="ts">
  import { X } from 'lucide-svelte';
  import StudentHeader from './absensi/StudentHeader.svelte';
  import DateTimeCard from './absensi/DateTimeCard.svelte';
  import GpsVerification from './absensi/GpsVerification.svelte';
  import DeviceDetails from './absensi/DeviceDetails.svelte';
  import PermitDetails from './absensi/PermitDetails.svelte';
  import EmptyState from './absensi/EmptyState.svelte';
  import FooterActions from './absensi/FooterActions.svelte';

  // Svelte 5 Props
  let {
    show = $bindable(false),
    log = null,
    onclose
  } = $props<{
    show: boolean;
    log: any;
    onclose: () => void;
  }>();

  let devInfo = $derived(log ? (log.deviceInfo || log.device) : null);

  // Escape key listener to close modal
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape' && show) {
      onclose();
    }
  }

  $effect(() => {
    if (show) {
      window.addEventListener('keydown', handleKeydown);
    }
    return () => {
      window.removeEventListener('keydown', handleKeydown);
    };
  });
</script>

{#if show && log}
  <!-- Overlay -->
  <div
    class="animate-fade-in fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 p-4 text-slate-700 backdrop-blur-xs"
    onclick={onclose}
    onkeydown={(e) => {
      if (e.key === 'Enter' || e.key === ' ') onclose();
    }}
    role="button"
    tabindex="0"
  >
    <!-- Modal Box -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div
      class="flex max-h-[90vh] w-full max-w-lg scale-100 transform flex-col overflow-hidden rounded-2xl bg-white shadow-2xl transition-all duration-300 border border-slate-100"
      onclick={(e) => e.stopPropagation()}
      onkeydown={() => {}}
      role="document"
      tabindex="-1"
    >
      <!-- Header -->
      <div class="flex items-center justify-between border-b border-slate-100 bg-slate-50/50 px-5 py-4">
        <div class="text-left">
          <h3 class="text-xs font-black tracking-wider text-slate-800 uppercase">
            Detail Informasi Absensi
          </h3>
          <p class="mt-0.5 text-[10px] font-bold text-slate-400">
            Log Kehadiran Siswa
          </p>
        </div>
        <button
          onclick={onclose}
          class="hover:text-slate-600 cursor-pointer rounded-lg border-none bg-transparent p-1.5 text-slate-400 transition-colors hover:bg-slate-100"
          title="Tutup"
        >
          <X class="h-4.5 w-4.5" />
        </button>
      </div>

      <!-- Content -->
      <div class="flex flex-1 flex-col gap-4 overflow-y-auto bg-slate-50/10 p-5 custom-scrollbar">
        
        <!-- Student Profil Header -->
        <StudentHeader {log} />

        <!-- Date and Time Card -->
        <DateTimeCard {log} />

        <!-- Conditional view based on status -->
        {#if log.status === 'hadir' || log.status === 'magang'}
          <!-- GPS Verification Details Section -->
          <GpsVerification {log} />

          <!-- Device Information Section -->
          {#if devInfo}
            <DeviceDetails {devInfo} />
          {/if}
        
        {:else if log.status === 'izin' || log.status === 'sakit'}
          <!-- Permit/Attachment Section -->
          <PermitDetails {log} />

        {:else}
          <!-- Empty/Placeholder State for Belum Absen -->
          <EmptyState />
        {/if}

      </div>

      <!-- Footer Actions -->
      <FooterActions {log} {onclose} />
    </div>
  </div>
{/if}

<style>
  @keyframes fadeIn {
    from { opacity: 0; }
    to { opacity: 1; }
  }
  .animate-fade-in {
    animation: fadeIn 0.2s ease-out forwards;
  }

  /* Custom scrollbar */
  .custom-scrollbar::-webkit-scrollbar {
    width: 4px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #cbd5e1;
    border-radius: 99px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #94a3b8;
  }
</style>
