<script lang="ts">
	import '../../static/main.css';
	import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import { onMount } from 'svelte';
    import { serverURL } from '../utils';
	import TitleBar from '../components/TitleBar.svelte';

	let error = '';
	let formID = 'register-confirm-form';
	let message = 'Check your Duke email for a pin to register.<br>'
		.concat('<small>Please use a <strong>unique password</strong> for this site ')
		.concat("(I'm not very good at security)</small>");
	let netID: string;
	let title = 'Register';

	const confirmRegistration = () => {
		return null;
	};

	const processRegister = async () => {
		let form = document.forms[formID];
		const data = new FormData(form);

		let body = {};
		for (let entry of data.entries()) {
			body[entry[0]] = entry[1];
		}

		// TODO hash password

		const res = await fetch(`${serverURL}/register`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		});

		form.reset();

		if (res.ok) {
			window.location.replace('/login');
		} else if (res.status == 401) {
			error = 'Registration failed - NetID already registered';
		} else {
			error = 'Registration failed - Retry in a bit';
		}
	};

	onMount(async () => {
		const urlParams = new URLSearchParams(window.location.search);
		netID = urlParams.get('net_id');
		console.log(netID);

		if (netID) {
			title += ` - ${netID}`;
		} else {
			window.location.replace('/register');
		}
	});
</script>

<main>
	<TitleBar />

	<Form id={formID} {title} submitFn={confirmRegistration} {error} {message}>
		<FormInput id="pin" title="Pin" />
		<FormInput id="password" title="Password" type="password" />
	</Form>
</main>
