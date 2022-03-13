<script lang="ts">
	export let disabled = false;
</script>

<main>
	{#if !disabled}
		<p><a href="/help" class="link">Help</a></p>
		<p><a href="/contact" class="link">Contact</a></p>
	{/if}
</main>

<style>
	a {
		color: white;
		text-decoration: none;
	}

	a:hover {
		cursor: pointer;
		text-decoration: underline;
	}

	main {
		width: 100%;
		min-height: 41.5px;
		margin: 0;
		padding-top: 2px;
		padding-bottom: 2px;
		padding-left: 10px;
		background-color: rgb(0, 48, 135);
		color: white;

		display: inline-flex;
		flex-direction: row;
		justify-content: flex-start;

		position: absolute;
		bottom: 0;
	}

	p {
		margin: 10px;
	}
</style>
