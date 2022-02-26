<script lang="ts">
    import { createEventDispatcher, getContext } from 'svelte';
    const { close } = getContext('simple-modal');

	export let id: string;
	export let title: string;
	export let message: string = '';
	export let error: string = '';
	export let submitFn: Function;
	// export let validateFn: Function;

	let dispatch = createEventDispatcher();
</script>

<div id="form-modal">
	<h3>{title}</h3>
	{#if message != ''}
		<p>{message}</p>
	{/if}
	<div>
		<form {id}>
			<slot />
		</form>

		<div class="flex-row button-row">
			<button type="button" on:click|preventDefault={() => close()}> Close </button>
			<button type="button" on:click|preventDefault={() => submitFn()}>Submit</button>
		</div>
		{#if error}
			<strong>{error}</strong>
		{/if}
	</div>
</div>

<style scoped>
	label {
		margin-right: 25px;
	}

	input {
		vertical-align: middle;
		float: right;
	}

	#form-modal {
		margin: auto;
		margin-top: 20px;
		width: 75%;
	}

	div {
		margin-bottom: 10px;
	}

	.flex-row {
		display: flex;
		justify-content: flex-end;
	}

    .button-row {
        float: right;
        width: 100%;
    }

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

	#error-msg {
		text-align: center;
		color: rgb(197, 33, 33);
	}
</style>
