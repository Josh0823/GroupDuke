<script lang="ts">
	import '../../../../static/main.css';

	import { serverURL } from '$lib/utils';

	import FooterBar from '../../../components/FooterBar.svelte';
	import Form from '../../../components/Form.svelte';
	import FormInput from '../../../components/FormInput.svelte';
	import TitleBar from '../../../components/TitleBar.svelte';

	export let username: string;
	export let pin: string;

	const formID = 'reset-password-form';
	let error = '';
	let message = '';

	const submitFn = async () => {
		const form = document.forms[formID];
		const data = new FormData(form);

		const password = data.get('password').toString();

		const res = await fetch(`${serverURL}/confirm-reset-password`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ username: username, password: password, pin: pin })
		});

		if (res.ok) {
			window.location.replace('/login');
		} else {
			error = 'Error resetting password';
		}
	};
</script>

<main>
	<TitleBar disabled={true} />

	<Form id={formID} title={`Reset Password for ${username}`} {submitFn} {error} {message}>
		<FormInput id="password" title="New Password"  type="password"/>
	</Form>

	<FooterBar />
</main>
