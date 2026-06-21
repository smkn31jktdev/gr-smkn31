<script lang="ts">
	import SubmitButton from '../../../../shared/components/SubmitButton.svelte';
	import DropdownChoice from '../../../../shared/components/DropdownChoice.svelte';

	let {
		waliKelas = $bindable(''),
		guruBK = $bindable(''),
		status = $bindable('draft'),
		isReadOnly = false,
		handleSave
	}: {
		waliKelas: string;
		guruBK: string;
		status: 'draft' | 'reviewed' | 'final';
		isReadOnly: boolean;
		handleSave: (handlers: { resolve: () => void; reject: () => void }) => void;
	} = $props();
</script>

<div class="space-y-4 card p-6">
	<div class="border-b border-border pb-3">
		<h3 class="text-sm font-bold text-foreground">Metadata Penilai & Status</h3>
		<p class="text-xxs mt-0.5 text-muted">Lengkapi data tim verifikator sekolah</p>
	</div>

	<!-- Assessor names -->
	<div class="space-y-3.5">
		<div>
			<!-- svelte-ignore a11y_label_has_associated_control -->
			<label class="text-xxs mb-1.5 block font-bold tracking-wider text-muted uppercase"
				>Nama Guru Wali</label
			>
			<input
				type="text"
				placeholder="Masukkan nama guru wali"
				bind:value={waliKelas}
				disabled={isReadOnly}
				class="input text-xs"
			/>
		</div>

		<div>
			<!-- svelte-ignore a11y_label_has_associated_control -->
			<label class="text-xxs mb-1.5 block font-bold tracking-wider text-muted uppercase"
				>Nama Guru BK</label
			>
			<input
				type="text"
				placeholder="Masukkan nama guru BK"
				bind:value={guruBK}
				disabled={isReadOnly}
				class="input text-xs"
			/>
		</div>



		<div>
			<!-- svelte-ignore a11y_label_has_associated_control -->
			<label class="text-xxs mb-1.5 block font-bold tracking-wider text-muted uppercase"
				>Status Penilaian</label
			>
			<div class="text-left">
				<DropdownChoice
					options={[
						{ value: 'draft', label: 'Draft (Belum Selesai)' },
						{ value: 'reviewed', label: 'Reviewed (Menunggu Final)' },
						{ value: 'final', label: 'Final (Kunci Nilai)' }
					]}
					bind:value={status}
					disabled={isReadOnly}
					placeholder="Status Penilaian"
				/>
			</div>
		</div>
	</div>

	<!-- Submit area -->
	{#if !isReadOnly}
		<div class="border-t border-border pt-4">
			<SubmitButton
				label="Simpan Hasil Nilai"
				loadingLabel="Menyimpan..."
				className="w-full py-3"
				onclick={handleSave}
			/>
		</div>
	{/if}
</div>
