<script lang="ts">
	import '../../static/main.css';

	import { isUserLoggedIn, logout, serverURL } from '$lib/utils';
	import { onMount } from 'svelte';

	import FooterBar from '../components/FooterBar.svelte';
	import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import TitleBar from '../components/TitleBar.svelte';

	let error = '';
	let formID = 'contact-form';
	let title = 'Send a message';

	let loggedIn = false;
	let username = '';

	const checkEmail = (email: String): boolean => {
		return email.includes('@duke.edu');
	};

	const submitFn = async () => {
		let form = document.forms[formID];
		const data = new FormData(form);

		let body = {};
		for (let entry of data.entries()) {
			body[entry[0]] = entry[1];
		}

		if (!checkEmail(body['email'])) {
			error = 'Please enter a valid Duke email';
		}

		const res = await fetch(`${serverURL}/contact`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(body)
		});

		if (!res.ok) {
			error = 'Error sending message';
		}

		console.log(res);
	};

	const handleLogout = () => {
		if (logout()) {
			loggedIn = false;
			username = '';
		} else {
			console.error('Failed to logout');
		}
	};

	onMount(() => {
		[loggedIn, username] = isUserLoggedIn();
	});
</script>

<main>
	<TitleBar {username} {loggedIn} on:logout={() => handleLogout()} />

	<Form id={formID} {title} {error} {submitFn}>
		<FormInput id="email" title="Email" />
		<FormInput id="subject" title="Subject" />
		<FormInput id="message" title="Message" isTextarea={true} />
	</Form>

	<FooterBar />
</main>
