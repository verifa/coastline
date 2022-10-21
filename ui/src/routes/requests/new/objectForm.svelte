<script lang="ts">
	import type { SchemaObject } from 'openapi-typescript';
	import { getInitialPropValue, propFromSchema, type Property } from './spec';
	import FieldForm from './fieldForm.svelte';

	export let store: { [key: string]: any };
	export let schemaObj: SchemaObject;

	function getProperties(obj: SchemaObject): Property[] {
		let properties: Property[] = [];
		// Reset the store
		store = {};
		for (const key in obj.properties) {
			const propObj: SchemaObject = obj.properties[key];
			store[key] = getInitialPropValue(propObj);
			properties.push(propFromSchema(key, obj, propObj));
		}
		return properties;
	}

	// Reactive declaration of properties in case schemaObj changes,
	// e.g. if user changes the request type
	$: properties = getProperties(schemaObj);
</script>

<div class="ml-4 flex flex-col space-y-4">
	{#each properties as prop}
		<FieldForm bind:store={store[prop.name]} {prop} />
	{/each}
</div>
