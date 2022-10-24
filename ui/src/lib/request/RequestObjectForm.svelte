<script lang="ts">
	import type { SchemaObject } from 'openapi-typescript';
	import { getInitialPropValue, propFromSchema, type Property } from './spec';
	import RequestPropertyForm from './RequestPropertyForm.svelte';

	export let store: { [key: string]: any } = {};
	export let schemaObj: SchemaObject;
	export let depth: number = 0;

	const indent: string = '  '.repeat(depth);

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

{#each properties as prop}
	<RequestPropertyForm bind:store={store[prop.name]} parent={schemaObj} {prop} {depth} />
{/each}
