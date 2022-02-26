<script lang="ts">
	import { createEventDispatcher, getContext } from 'svelte';
    const { close } = getContext('simple-modal');

	import FormInput from '../components/FormInput.svelte';
	import ModalForm from '../components/ModalForm.svelte';

	export let serverURL: string;

	let error: string = '';
	let dispatch = createEventDispatcher();

	const processLogin = () => {
		let form = document.forms['login-form'];
		const data = new FormData(form);

		let user = {};
		for (let entry of data.entries()) {
			user[entry[0]] = entry[1];
		}

		fetch(`${serverURL}/login`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(user)
		}).then((res) => {
			if (res.ok) {
                // FIND A WAY TO CATCH THIS EVENT
                dispatch('login', user['username']);
                close();
			} else {
				error = 'Login failed - Retry';
				form.reset();
			}
		});
	};
</script>

<ModalForm
	id="login-form"
	title="Login"
	submitFn={processLogin}
	{error}
	on:close={() => dispatch('close')}
>
	<FormInput id="username" title="NetID" />
	<FormInput id="password" title="Password" type="password" />
</ModalForm>
