<script lang="ts">
	import type { SchemaObject } from 'openapi-typescript';
	import Textfield from '@smui/textfield';
	import Switch from '@smui/switch';
	import FormField from '@smui/form-field';
	import List, { Item, Graphic, Separator, Text } from '@smui/list';
	import type { Writable } from 'svelte/store';
	import IconButton from '@smui/icon-button';

	export let spec: SchemaObject;
	export let store: Writable<{ [key: string]: any }>;

	type Property = {
		name: string;
		schema: SchemaObject;
		is_required: boolean;
	};

	let properties: Property[] = [];
	for (const key in spec.properties) {
		const prop: SchemaObject = spec.properties[key];
		prop.required;
		// Basic types:
		// Ref: https://swagger.io/docs/specification/data-models/data-types/
		//   string (this includes dates and files)
		//   number
		//   integer
		//   boolean
		//   array
		//   object
		switch (prop.type) {
			case 'string': {
				$store[key] = prop.default || '';
				break;
			}
			case 'number': {
				$store[key] = prop.default || 0;
				break;
			}
			case 'integer': {
				$store[key] = prop.default || 0;
				break;
			}
			case 'boolean': {
				$store[key] = prop.default || false;
				break;
			}
			case 'array': {
				$store[key] = prop.default || [];
				break;
			}
			case 'object': {
				$store[key] = prop.default || {};
				break;
			}
			default: {
				console.log('error: unsupported spec type: ', prop.type);
				break;
			}
		}
		properties.push({
			name: key,
			schema: prop,
			is_required: (spec.required || []).includes(key)
		});
	}
</script>

<div class="ml-4 flex flex-col space-y-4">
	{#each properties as prop}
		{#if prop.schema.type == 'string'}
			<Textfield
				label={prop.name}
				type={prop.schema.type}
				required={prop.is_required}
				bind:value={$store[prop.name]}
			/>
		{:else if prop.schema.type == 'number'}
			<Textfield
				label={prop.name}
				type={prop.schema.type}
				required={prop.is_required}
				bind:value={$store[prop.name]}
			/>
		{:else if prop.schema.type == 'integer'}
			<Textfield
				label={prop.name}
				type={'number'}
				required={prop.is_required}
				bind:value={$store[prop.name]}
			/>
		{:else if prop.schema.type == 'boolean'}
			<div class="text-left">
				<FormField align={'end'}>
					<Switch bind:checked={$store[prop.name]} />
					<span slot="label">Fields of grain.</span>
				</FormField>
			</div>
		{:else if prop.schema.type == 'array'}
			<div>
				<div class="flex items-center">
					<p>{prop.name}</p>
					<IconButton
						type="button"
						class="material-icons"
						on:click={() => {
							store.update((value) => {
								// TODO: need to know what kind of value
								value[prop.name].push('');
								return value;
							});
						}}>add</IconButton
					>
				</div>
				<List>
					{#each $store[prop.name] as item, index}
						<Item
							on:SMUI:action={() => {
								store.update((value) => {
									value[prop.name].splice(index, 1);
									return value;
								});
							}}
						>
							<Graphic class="material-icons">clear</Graphic>
							<Text>{item}</Text>
						</Item>
					{/each}
				</List>
			</div>
		{:else if prop.schema.type == 'object'}
			<p>{prop.name}:</p>
			<svelte:self {store} spec={prop.schema} />
		{:else}
			<h3>Error: unsupported type {prop.schema.type}</h3>
		{/if}
	{/each}
</div>
