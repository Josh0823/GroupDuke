<script lang="ts">
	export let id: string;
	export let placeholder: string = '';
	export let title: string;
	export let type: string = 'text';

	export let isTextarea = false;
</script>

<div class="form-group">
	<label for={id}>{@html title}</label>
	{#if !isTextarea}
		<input name={id} {id} {placeholder} {type} required />
	{:else}
		<textarea name={id} {id} value={placeholder} {type} required />
	{/if}
</div>

<style>
	input {
		float: right;
		margin-bottom: 15px;
	}

	textarea {
		height: 50px;
		width: 183px;
	}

	.form-group {
		display: inline-flex;
		flex-direction: row;
		justify-content: space-between;
		width: 100%;
	}
</style>
