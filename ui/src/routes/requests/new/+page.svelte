<script lang="ts">
	import { createHttpStore } from '$lib/http/store';
	import { getRequestSpecs } from '$lib/oapi/parse';
	import type { RequestSpec } from '$lib/oapi/parse';

	import Autocomplete from '@smui-extra/autocomplete';

	import type { OpenAPI3 } from 'openapi-typescript';
	import Button from '@smui/button/src/Button.svelte';

	const store = createHttpStore<OpenAPI3>();
	store.get('/requestsspec');

	let specs: RequestSpec[];
	let selectedRequest: number;

	store.subscribe((value) => {
		if (value.ok && value.data) {
			specs = getRequestSpecs(value.data);
		}
	});
</script>

<h1>New Request</h1>

{#if $store.fetching}
	<h2>Loading</h2>
{:else if $store.ok}
	<h2>Form</h2>

	<label for="type">Request type:</label>
	<Button>Hello</Button>
	<Autocomplete
		getOptionLabel={(item) => {
			console.log('item: ', item);
			return item ? item.type : '';
		}}
		options={specs}
		label="Standard"
	/>
	<!-- bind:value={valueStandard} -->
	<select
		name="type"
		id="type"
		form="request"
		bind:value={selectedRequest}
		on:change={() => console.log(selectedRequest)}
	>
		<option value="-1" disabled selected>Select request</option>
		{#each specs as spec, index}
			<option value={index}>{spec.type}</option>
		{/each}
	</select>
	{#if selectedRequest >= 0}
		<h3>Render custom form here...</h3>
		<pre>{JSON.stringify(specs[selectedRequest].spec)}</pre>
	{/if}
{/if}
