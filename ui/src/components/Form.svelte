<script lang="ts">
	export let id: string;
	export let title: string;
	export let message: string = '';
	export let error: string = '';
	export let submitFn: Function;
	export let cancelFn: Function = null;
	// export let validateFn: Function;
</script>

<main>
	<h3>{title}</h3>
	{#if message != ''}
		<p>{@html message}</p>
	{/if}
	<div>
		<form {id}>
			<slot />
		</form>

		<div class="flex-row button-row">
			{#if cancelFn}
				<button type="button" on:click|preventDefault={() => cancelFn()}>Cancel</button>
			{/if}
			<button type="button" on:click|preventDefault={() => submitFn()}>Submit</button>
		</div>
		{#if error}
			<strong class="error-msg">{error}</strong>
		{/if}
	</div>
</main>

<style scoped>
	button {
		margin-top: 10px;
		margin-left: 10px;
		vertical-align: middle;

		background-color: rgb(0, 48, 135);
		color: white;
		border: 1px solid rgb(0, 48, 135);
		border-radius: 10px;
		padding: 5px 10px 5px 10px;
	}

	button:hover {
		cursor: pointer;
	}

	div {
		margin-bottom: 10px;
	}

	main {
		margin: auto;
		margin-top: 20px;
		width: 75%;
		max-width: 400px;
		padding: 20px;
		border-radius: 5%;
		background-color: rgb(247, 247, 247);
		border: 1px solid rgb(200, 200, 200);
	}

	p {
		margin-top: 0px;
	}

	.button-row {
		width: 100%;
	}

	.error-msg {
		text-align: center;
		color: rgb(197, 33, 33);
	}

	.flex-row {
		display: flex;
		justify-content: flex-end;
	}
</style>
