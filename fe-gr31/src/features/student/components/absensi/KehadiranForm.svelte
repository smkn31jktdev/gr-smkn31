<script lang="ts">
  import { onMount } from 'svelte';
  import { getGeolocation, submitKehadiran, uploadIzinFile } from '../../logic/kehadiranLogic';
  import { addToast } from '../../../../stores/uiStore';
  import SubmitButton from '../../../shared/components/SubmitButton.svelte';
  import type { LatLng } from '../../types/student.types';
  import { School, Mail, Frown, RefreshCw } from 'lucide-svelte';
  import { getUploadUrl } from '../../../../api/client';

  import ClockTimer from '../../../shared/components/ClockTimer.svelte';

  let {
    defaultStatus = 'hadir',
    allowedStatuses = ['hadir', 'izin', 'sakit']
  } = $props<{
    defaultStatus?: 'hadir' | 'izin' | 'sakit';
    allowedStatuses?: Array<'hadir' | 'izin' | 'sakit'>;
  }>();

  let status = $state<'hadir' | 'izin' | 'sakit'>(defaultStatus);
  let alasan = $state('');
  let locationLoading = $state(false);
  
  let koordinat = $state<LatLng | undefined>(undefined);
  let akurasi = $state<number | undefined>(undefined);
  let fotoIzin = $state('');
  let selectedFile = $state<File | null>(null);

  // Auto get location if status is 'hadir'
  async function fetchLocation() {
    if (status !== 'hadir') return;
    locationLoading = true;
    const loc = await getGeolocation();
    if (loc) {
      koordinat = loc.koordinat;
      akurasi = loc.akurasi;
      addToast('Lokasi GPS berhasil didapatkan', 'success');
    }
    locationLoading = false;
  }

  // React to status changes
  $effect(() => {
    if (status === 'hadir') {
      fetchLocation();
    } else {
      koordinat = undefined;
      akurasi = undefined;
    }
  });

  onMount(() => {
    if (status === 'hadir') {
      fetchLocation();
    }
  });

  async function handleFileChange(e: Event) {
    const input = e.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      selectedFile = input.files[0];
      // Upload file right away to get the URL
      addToast('Mengunggah berkas...', 'info');
      const url = await uploadIzinFile(selectedFile);
      if (url) {
        fotoIzin = url;
      } else {
        selectedFile = null;
      }
    }
  }

  async function handleSubmit(handlers: { resolve: () => void; reject: () => void }) {
    if (status === 'hadir' && !koordinat) {
      addToast('GPS koordinat diperlukan untuk absensi Hadir', 'error');
      handlers.reject();
      return;
    }

    if ((status === 'izin' || status === 'sakit') && !fotoIzin) {
      addToast('Harap unggah surat/bukti izin/sakit terlebih dahulu', 'error');
      handlers.reject();
      return;
    }

    const success = await submitKehadiran({
      status,
      alasan,
      koordinat,
      akurasi,
      fotoIzin: fotoIzin || undefined
    });

    if (success) {
      handlers.resolve();
      // Reset form
      alasan = '';
      fotoIzin = '';
      selectedFile = null;
      if (status === 'hadir') {
        fetchLocation();
      }
    } else {
      handlers.reject();
    }
  }
</script>

<div class="card p-6 space-y-6">
  <div>
    <h3 class="text-base font-bold text-foreground">Isi Kehadiran Harian</h3>
    <p class="text-xs text-muted mt-0.5">Absen masuk harian sesuai jadwal sekolah</p>
  </div>

  <!-- Real-time Clock for Student Reference -->
  {#if status === 'hadir'}
    <div class="border-y border-gray-100 py-4 w-full flex justify-center bg-gray-50/30 rounded-2xl">
      <ClockTimer mode="clock" allowSwitch={false} />
    </div>
  {/if}

  <div class="space-y-4">
    <!-- Status Select -->
    <div>
      <label class="block text-xs font-bold uppercase tracking-wider text-muted mb-2">Status Kehadiran</label>
      <div
        class="grid gap-3"
        class:grid-cols-3={allowedStatuses.length === 3}
        class:grid-cols-2={allowedStatuses.length === 2}
        class:grid-cols-1={allowedStatuses.length === 1}
      >
        {#if allowedStatuses.includes('hadir')}
          <button
            type="button"
            onclick={() => status = 'hadir'}
            class="flex flex-col items-center justify-center p-3.5 rounded-xl border-2 font-semibold text-sm transition-all cursor-pointer"
            class:border-primary={status === 'hadir'}
            class:bg-primary-light={status === 'hadir'}
            class:text-primary={status === 'hadir'}
            class:border-border={status !== 'hadir'}
            class:text-muted={status !== 'hadir'}
          >
            <School class="w-5 h-5 mb-1.5" />
            Hadir
          </button>
        {/if}

        {#if allowedStatuses.includes('izin')}
          <button
            type="button"
            onclick={() => status = 'izin'}
            class="flex flex-col items-center justify-center p-3.5 rounded-xl border-2 font-semibold text-sm transition-all cursor-pointer"
            class:border-primary={status === 'izin'}
            class:bg-primary-light={status === 'izin'}
            class:text-primary={status === 'izin'}
            class:border-border={status !== 'izin'}
            class:text-muted={status !== 'izin'}
          >
            <Mail class="w-5 h-5 mb-1.5" />
            Izin
          </button>
        {/if}

        {#if allowedStatuses.includes('sakit')}
          <button
            type="button"
            onclick={() => status = 'sakit'}
            class="flex flex-col items-center justify-center p-3.5 rounded-xl border-2 font-semibold text-sm transition-all cursor-pointer"
            class:border-primary={status === 'sakit'}
            class:bg-primary-light={status === 'sakit'}
            class:text-primary={status === 'sakit'}
            class:border-border={status !== 'sakit'}
            class:text-muted={status !== 'sakit'}
          >
            <Frown class="w-5 h-5 mb-1.5" />
            Sakit
          </button>
        {/if}
      </div>
    </div>

    <!-- Hadir Location Status -->
    {#if status === 'hadir'}
      <div class="p-4 rounded-xl border bg-gray-50 flex items-center justify-between gap-3 text-sm">
        <div class="flex items-center gap-2.5">
          <span class="relative flex h-3 w-3">
            {#if locationLoading}
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-amber-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-3 w-3 bg-amber-500"></span>
            {:else if koordinat}
              <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
            {:else}
              <span class="relative inline-flex rounded-full h-3 w-3 bg-rose-500"></span>
            {/if}
          </span>
          <div>
            <p class="font-bold text-foreground">Lokasi GPS Anda</p>
            {#if locationLoading}
              <p class="text-xs text-muted mt-0.5">Membaca satelit GPS...</p>
            {:else if koordinat}
              <p class="text-xs text-muted mt-0.5">Latitude: {koordinat.lat.toFixed(6)}, Longitude: {koordinat.lng.toFixed(6)} (Akurasi: {akurasi?.toFixed(1)}m)</p>
            {:else}
              <p class="text-xs text-rose-600 mt-0.5">Gagal membaca koordinat GPS</p>
            {/if}
          </div>
        </div>

        <button
          type="button"
          onclick={fetchLocation}
          disabled={locationLoading}
          class="text-xs font-bold text-primary hover:text-primary-hover flex items-center gap-1 hover:underline cursor-pointer"
        >
          <RefreshCw class="w-3.5 h-3.5 {locationLoading ? 'animate-spin' : ''}" />
          Segarkan
        </button>
      </div>
    {/if}

    <!-- Permit File Input -->
    {#if status === 'izin' || status === 'sakit'}
      <div class="space-y-3">
        <div>
          <label class="block text-xs font-bold uppercase tracking-wider text-muted mb-2">Unggah Surat Bukti (Maks 5MB)</label>
          <input
            type="file"
            accept="image/*,.pdf"
            onchange={handleFileChange}
            class="input file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-xs file:font-semibold file:bg-primary-light file:text-primary hover:file:bg-primary-light/80"
          />
        </div>

        {#if fotoIzin}
          <div class="p-3 bg-emerald-50 border border-emerald-200 rounded-xl text-emerald-800 text-xs font-semibold flex items-center justify-between">
            <span>✓ Berkas berhasil diunggah dan ditautkan</span>
            <a href={getUploadUrl(fotoIzin)} target="_blank" class="underline text-primary">Lihat Berkas</a>
          </div>
        {/if}

        <div>
          <label for="alasan" class="block text-xs font-bold uppercase tracking-wider text-muted mb-2">Alasan / Keterangan</label>
          <textarea
            id="alasan"
            placeholder="Tuliskan keterangan detail alasan Anda..."
            bind:value={alasan}
            class="input min-h-[80px] py-2"
          ></textarea>
        </div>
      </div>
    {/if}

    <!-- Submit Action -->
    <div class="pt-2">
      <SubmitButton
        label="Kirim Absen"
        loadingLabel="Mengirim..."
        className="w-full py-3"
        onclick={handleSubmit}
        disabled={status === 'hadir' && !koordinat && !locationLoading}
      />
    </div>
  </div>
</div>
