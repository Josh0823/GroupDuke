<script lang="ts">
	import '../../static/main.css';

	import { isUserLoggedIn, logout } from '$lib/utils';
	import { onMount } from 'svelte';

	import FooterBar from '../components/FooterBar.svelte';
	import TitleBar from '../components/TitleBar.svelte';

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

	<div class="content">
		<div class="info-box">
			<h3>Where do I find the link to share a GroupMe group?</h3>
			<p>
				Go to the group you want to share and click the group's avatar. Then click on Settings and
				look for the Share Group option. You can copy the link from there.
			</p>
		</div>

		<div class="info-box">
			<h3>How do I reset my password?</h3>
			<p>
				You can click <a href="/reset-password">here</a>
				to reset your password.
			</p>
		</div>

		<div class="info-box">
			<h3>Can I remove a group from the database?</h3>
			<p>
				If you uploaded a group to the database you can remove it at any time by clicking on the
				group in the table and then clicking the trash can icon.
			</p>
			<br />
			<p>
				If you did not upload a group but would still like it removed, send an email from your duke
				email account to
				<a href="mainto:groupduke2023@gmail.com"> GroupDuke2023@gmail.com </a>
				with the details.
			</p>
		</div>
	</div>

	<FooterBar />
</main>

<style>
	.content {
		display: flex;
		flex-direction: column;
		justify-content: flex-start;

		width: 50%;
		margin: auto;
		margin-top: 20px;
	}

	.content h3 {
		text-align: center;
		margin: 0;
		margin-bottom: 5px;
	}

	.content p {
		margin: 0;
	}

	.info-box {
		border-radius: 25px;
		padding: 25px;
		margin-bottom: 40px;
		max-width: 600px;
		background-color: rgb(247, 247, 247);
		border: 1px solid rgb(200, 200, 200);
	}
</style>
