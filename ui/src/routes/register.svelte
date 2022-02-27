<script lang="ts">
	import '../../static/main.css';
	import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import { serverURL } from '../utils';
	import TitleBar from '../components/TitleBar.svelte';

	let error = '';
	let formID = 'register-form';
	let message = 'Enter a valid Duke NetID.<br>'
		.concat('<small>Only students can register.</small>');

	const processRegister = () => {
		const form = document.forms[formID];
		const data = new FormData(form);
		const id = data.get('username');

		if (id != '') {
			window.location.replace(`/register-confirm?net_id=${id}`);
		}
	};
</script>

<main>
	<TitleBar />

	<Form id={formID} title="Register" submitFn={processRegister} {error} {message}>
		<FormInput id="username" title="NetID" />
	</Form>
</main>
