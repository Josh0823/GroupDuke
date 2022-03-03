<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	export let disabled: boolean = false;
	export let loggedIn: boolean = false;
	export let username: string = '';
</script>

<main>
	<div>
		<p style="margin-left: 10px;">
			<a href="/">GroupDuke</a>
		</p>
	</div>
	<div style="display: inline-flex;">
		{#if loggedIn && !disabled}
			<p>{username}<span on:click={() => dispatch('logout')}>Logout</span></p>
		{:else if !disabled}
			<div class="flex-row">
				<p>
					<a href="/register">Register</a>
					<a href="/login">Login</a>
				</p>
			</div>
		{/if}
	</div>
</main>

<style>
	a {
		margin-left: 20px;
		color: white;
		text-decoration: none;
	}

	a:hover {
		cursor: pointer;
		text-decoration: underline;
	}

	div {
		vertical-align: text-bottom;
	}

	main {
		width: 100%;
		margin: 0;
		padding-top: 10px;
		padding-bottom: 10px;
		background-color: rgb(0, 48, 135);
		color: white;

		display: inline-flex;
		flex-direction: row;
		justify-content: space-between;
	}

	p {
		margin: 10px 20px 10px 20px;
	}

	span {
		margin-left: 20px;
		color: white;
		text-decoration: none;
	}

	span:hover {
		cursor: pointer;
		text-decoration: underline;
	}

	.flex-row {
		display: inline-flex;
		justify-content: flex-end;
	}

	.flex-row p {
		margin-right: 20px;
	}
</style>
