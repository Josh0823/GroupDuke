<script lang="ts">
	import '../../static/main.css';

	import { isUserLoggedIn, logout } from '$lib/utils';
	import { onMount } from 'svelte';

	import FooterBar from '../components/FooterBar.svelte';
	import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import TitleBar from '../components/TitleBar.svelte';

	let formID = 'contact-form';
	let title = 'Send a message';
	let submitFn = null;

	let loggedIn = false;
	let username = '';

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

	<Form id={formID} {title} {submitFn}>
		<FormInput id="email" title="Email" />
		<FormInput id="subject" title="Subject" />
		<FormInput id="message" title="Message" isTextarea={true} />
	</Form>

	<FooterBar />
</main>
