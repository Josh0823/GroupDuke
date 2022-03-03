<script lang="ts">
	import '../../static/main.css';

	import { isUserLoggedIn } from '$lib/utils';
	import { onMount } from 'svelte';

    import FooterBar from '../components/FooterBar.svelte';
    import Form from '../components/Form.svelte';
	import FormInput from '../components/FormInput.svelte';
	import TitleBar from '../components/TitleBar.svelte';

	const processSubmit = () => {
		return true;
	};

	let formID = 'reset-password';
	let title = 'Reset Password';
	let submitFn = processSubmit;

	let loggedIn = false;
	let username = '';

	onMount(() => {
		[loggedIn, username] = isUserLoggedIn();
	});
</script>

<main>
	<TitleBar {loggedIn} {username} />

	<Form id={formID} {title} {submitFn}>
		<FormInput id="username" title="NetID" />
	</Form>

    <FooterBar />
</main>
