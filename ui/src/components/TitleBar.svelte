<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	export let disabled: boolean = false;
	export let loggedIn: boolean = false;
	export let netID: string = '';
</script>

<div class="TitleBar">
	<div>
		<p style="margin-left: 10px;">
			<span on:click={() => window.location.replace('/')} class="link">GroupDuke</span>
		</p>
	</div>
	<div style="display: inline-flex;">
		{#if loggedIn && !disabled}
			<p>{netID}<span on:click={() => dispatch('logout')} class="link">Logout</span></p>
		{:else if !disabled}
			<div class="flex-row">
				<p>
					<span on:click={() => window.location.replace('register')} class="link">Register</span>
					<span on:click={() => window.location.replace('login')} class="link">Login</span>
				</p>
			</div>
		{/if}
	</div>
</div>

<style>
	div {
		vertical-align: text-bottom;
	}

	p {
		margin: 10px 20px 10px 20px;
	}

	.flex-row {
		display: inline-flex;
		justify-content: flex-end;
	}

	.flex-row p {
		margin-right: 20px;
	}

	.link {
		margin-left: 20px;
		color: white;
	}

	.link:hover {
		cursor: pointer;
		text-decoration: underline;
	}

	.TitleBar {
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
</style>
