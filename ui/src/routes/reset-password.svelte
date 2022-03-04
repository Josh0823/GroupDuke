<script lang="ts">
	import '../../static/main.css';

	import { isUserLoggedIn, serverURL } from '$lib/utils';
	import { onMount } from 'svelte';

	import FooterBar from '../components/FooterBar.svelte';
	import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import TitleBar from '../components/TitleBar.svelte';

	let formID = 'reset-password';
	let title = 'Reset Password';
	let message =
		'Enter your NetID and a link to reset your password will be sent to the your Duke email';
	let error = '';

	let emailSent = false;
	let username: string;

	const submitFn = async () => {
		const form = document.forms[formID];
		const data = new FormData(form);
		const user = data.get('username');

		if (user === '') {
			return;
		}

		username = user.toString();

		const res = await fetch(`${serverURL}/reset-password`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ username: username })
		});

		if (res.ok) {
			emailSent = true;
		} else {
			error = 'Error requesting password reset';
		}
	};
</script>

<main>
	<TitleBar />

	{#if !emailSent}
		<Form id={formID} {title} {message} {error} {submitFn}>
			<FormInput id="username" title="NetID" />
		</Form>
	{:else}
		<div>
			<p>Check {username}@duke.edu to confirm your registration</p>
			<small
				>No email? Click here to <a href="/register" on:click={() => (emailSent = false)}>resend</a
				></small
			>
		</div>
	{/if}

	<FooterBar />
</main>

<style>
	div {
		text-align: center;
	}
</style>
