<script lang="ts">
	import '../../static/main.css';
	import Form from '../components/Form.svelte';
    import FormInput from '../components/FormInput.svelte';
	import { serverURL } from '../utils';
	import TitleBar from '../components/TitleBar.svelte';

	let error = '';
	let formID = 'login-form';

	const processLogin = async () => {
		let form = document.forms[formID];
		const data = new FormData(form);

		let user = {};
		for (let entry of data.entries()) {
			user[entry[0]] = entry[1];
		}

		const res = await fetch(`${serverURL}/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(user)
		});

		form.reset();

		if (res.ok) {
			window.location.replace('/');
		} else if (res.status == 401) {
			error = 'Login failed - Invalid NetID or Password';
		} else {
			error = 'Login failed - Retry';
		}
	};
</script>

<main>
	<TitleBar on:register={() => window.location.replace('/register')} />

	<Form id={formID} title="Login" submitFn={processLogin} {error}>
		<FormInput id="username" title="NetID" />
		<FormInput id="password" title="Password" type="password" />
	</Form>
</main>
