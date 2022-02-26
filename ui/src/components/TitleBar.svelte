<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let netID: string;
	export let loggedIn: boolean;

	const dispatch = createEventDispatcher();
</script>

<div class="TitleBar">
	<div>
		<p style="margin-left: 10px;">GroupDuke</p>
	</div>
	<div style="display: inline-flex;">
		{#if loggedIn}
			<p>{netID}<span on:click={() => dispatch('logout')} class="link">Logout</span></p>
		{:else}
			<div class="flex-row">
				<p>
					<span on:click={() => dispatch('register')} class="link">Register</span>
					<span on:click={() => dispatch('login')} class="link">Login</span>
				</p>
			</div>
		{/if}
	</div>
</div>

<style>
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

	div {
		vertical-align: text-bottom;
	}

	.link {
		margin-left: 20px;
		color: white;
	}

	.link:hover {
		cursor: pointer;
		text-decoration: underline;
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
</style>
