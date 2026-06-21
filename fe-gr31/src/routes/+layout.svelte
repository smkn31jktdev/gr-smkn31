<script lang="ts">
	import './layout.css';
	import Toast from '../features/shared/components/Toast.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<title>GR31 SMKN 31 Jakarta</title>
</svelte:head>

<Toast />

{@render children()}
