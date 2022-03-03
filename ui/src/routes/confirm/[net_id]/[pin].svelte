<script lang="ts">
	import '../../../../static/main.css';

	import { onMount } from 'svelte';
	import { serverURL } from '$lib/utils';

	import TitleBar from '../../../components/TitleBar.svelte';

	export let username: string;
	export let pin: string;

	let running = true;
	let success = false;

	onMount(async () => {
		const res = await fetch(`${serverURL}/confirm-registration`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ username: username, pin: pin })
		});

		running = false;
		success = res.ok;
	});
</script>

<main>
	<TitleBar disabled={true} />

	{#if running}
		<p>Checking registration...</p>
	{:else if success}
		<p>Registration succeeded</p>
		<p><a href="/login">Login here</a></p>
	{:else}
		<p>Registration failed.</p>
		<p><a href="/register">Try again</a></p>
	{/if}
</main>

<style>
	p {
		text-align: center;
	}
</style>
