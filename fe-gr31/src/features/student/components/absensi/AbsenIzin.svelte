<script lang="ts">
	import Modal from '../../../shared/components/Modal.svelte';
	import KehadiranForm from './KehadiranForm.svelte';

	let {
		show = $bindable(false),
		onload
	} = $props<{
		show: boolean;
		onload: () => void;
	}>();
</script>

<Modal
	{show}
	title=""
	onclose={() => {
		show = false;
		onload();
	}}
>
	<div class="px-2">
		<KehadiranForm defaultStatus="izin" allowedStatuses={['izin', 'sakit']} />
	</div>
</Modal>
