<script lang="ts">
	import type { SchemaObject } from 'openapi-typescript';
	import Textfield from '@smui/textfield';
	import { writable } from 'svelte/store';

	export let spec: SchemaObject;

	type Property = {
		name: string;
		schema: SchemaObject;
		is_required: boolean;
	};

	const store = writable<Record<string, any>>({});

	store.subscribe((value) => {
		console.log(value);
	});

	let properties: Property[] = [];
	for (const key in spec.properties) {
		$store[key] = '';
		properties.push({
			name: key,
			schema: spec.properties[key],
			is_required: key in (spec.required || [])
		});
	}
</script>

<!-- {"type":"object","required":["name"],"properties":{"name":{"type":"string"}}} -->

<p>{JSON.stringify(spec)}</p>

{#each properties as prop}
	{#if prop.schema.type == 'string'}
		<Textfield bind:value={$store[prop.name]} label={prop.name} required={prop.is_required} />
	{/if}
{/each}
