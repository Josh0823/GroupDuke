<script lang="ts">
	import '../../static/main.css';
	import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import { serverURL } from '$lib/utils';
	import TitleBar from '../components/TitleBar.svelte';

	let error = '';
	let formID = 'register-form';
	let message = 'Enter a valid Duke NetID.<br>'.concat(
		'<small>Only students can register.</small>'
	);
	let emailSent = false;
	let netID: string;

	const processRegister = async () => {
		const form = document.forms[formID];
		const data = new FormData(form);
		const username = data.get('username');
		const password = data.get('password');

		if (username === '' || password === '') {
			return;
		}

		netID = username.toString();

		const res = await fetch(`${serverURL}/register`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ username: username, password: password })
		});

		if (res.ok) {
			emailSent = true;
		} else {
			error = 'Error validating student netID';
		}
	};
</script>

<main>
	<TitleBar />

	{#if !emailSent}
		<Form id={formID} title="Register" submitFn={processRegister} {error} {message}>
			<FormInput id="username" title="NetID" />
			<FormInput id="password" title="Password" type="password" />
		</Form>
	{:else}
		<div>
			<p>Check {netID}@duke.edu to confirm your registration</p>
			<small
				>No email? Click here to <a href="/register" on:click={() => (emailSent = false)}>resend</a
				></small
			>
		</div>
	{/if}
</main>

<style>
	div {
		text-align: center;
	}
</style>
