<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { SchemaObject } from 'openapi-typescript';
	import { getInitialPropValue, propFromSchema, type Property } from './spec';
	import InputField from './fieldInput.svelte';
	import SchemaObjectForm from './schemaObjectForm.svelte';
	import FieldProperty from './fieldProperty.svelte';

	export let store: any;
	export let parent: SchemaObject;
	export let prop: Property;
	export let depth: number = 0;

	const dispatch = createEventDispatcher();

	function handleDelete() {
		dispatch('delete');
	}

	const indent: string = '  '.repeat(depth);

	const isArrayElement: boolean = parent.type === 'array';
</script>

{#if prop.schema.type == 'object'}
	<pre><code><FieldProperty bind:store {prop} {depth} /></code></pre>
	<SchemaObjectForm bind:store schemaObj={prop.schema} depth={depth + 1} />
{:else if prop.schema.type == 'array'}
	<pre><code
			><FieldProperty bind:store {prop} {depth} /> {#if isArrayElement}<span
					class="text-error hover:cursor-pointer"
					on:click={handleDelete}>x</span
				>{/if}</code
		></pre>
	{#each store as item, index}
		<svelte:self
			bind:store={item}
			parent={prop.schema}
			prop={propFromSchema(prop.name, prop.schema, prop.schema.items)}
			depth={depth + 1}
			on:delete={() => {
				store.splice(index, 1);
				// Need to reassign variable to trigger a reactive update
				store = store;
			}}
		/>
	{/each}
{:else if isArrayElement}
	<pre><code
			><span>{indent + ' - '}</span><InputField bind:store {prop} /><span
				class="text-error hover:cursor-pointer"
				on:click={handleDelete}>x</span
			></code
		></pre>
{:else}
	<pre><code
			><FieldProperty bind:store {prop} {depth}><InputField bind:store {prop} /></FieldProperty
			></code
		></pre>
{/if}
