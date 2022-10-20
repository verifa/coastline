<script lang="ts">
	import type { SchemaObject } from 'openapi-typescript';
	import { getInitialPropValue, propFromSchema, type Property } from './spec';
	import FieldForm from './fieldForm.svelte';

	export let spec: SchemaObject;
	export let store: { [key: string]: any };

	let properties: Property[] = [];
	for (const key in spec.properties) {
		const schemaObj: SchemaObject = spec.properties[key];
		store[key] = getInitialPropValue(schemaObj);
		properties.push(propFromSchema(key, spec, schemaObj));
	}
</script>

<div class="ml-4 flex flex-col space-y-4">
	{#each properties as prop}
		<FieldForm bind:store={store[prop.name]} {prop} />
	{/each}
</div>
