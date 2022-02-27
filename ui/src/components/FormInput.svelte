<script lang="ts">
	export let id: string;
	export let placeholder: string = '';
	export let title: string;
	export let type: string = 'text';
</script>

<div class="form-group">
	<label for={id}>{title}</label>
	<input name={id} {id} {placeholder} {type} required />
</div>

<style>
	input {
		float: right;
		margin-bottom: 15px;
	}

	.form-group {
		display: inline-flex;
		flex-direction: row;
		justify-content: space-between;
		width: 100%;
	}
</style>