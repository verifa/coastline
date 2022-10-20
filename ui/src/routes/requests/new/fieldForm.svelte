<script lang="ts">
	import IconButton from '@smui/icon-button';
	import Textfield from '@smui/textfield';
	import { getInitialPropValue, propFromSchema, type Property } from './spec';
	import FormField from '@smui/form-field';
	import Switch from '@smui/switch';
	import ObjectForm from './objectForm.svelte';
	import Select, { Option } from '@smui/select';

	export let prop: Property;
	export let store: any;
</script>

{#if prop.schema.enum}
	<!-- 
    key field sets a function that makes sure each value is a string
    which is a requirement of this component.
 -->
	<Select
		bind:value={store}
		label={prop.name}
		required={prop.is_required}
		key={(val) => String(val)}
	>
		{#each prop.schema.enum as item}
			<Option value={item}>{item}</Option>
		{/each}
	</Select>
{:else if prop.schema.type == 'string'}
	<Textfield
		label={prop.name}
		type={prop.schema.type}
		required={prop.is_required}
		bind:value={store}
	/>
{:else if prop.schema.type == 'number'}
	<Textfield
		label={prop.name}
		type={prop.schema.type}
		required={prop.is_required}
		bind:value={store}
	/>
{:else if prop.schema.type == 'integer'}
	<Textfield label={prop.name} type={'number'} required={prop.is_required} bind:value={store} />
{:else if prop.schema.type == 'boolean'}
	<div class="text-left">
		<FormField align={'end'}>
			<Switch bind:checked={store} />
			<span slot="label">Fields of grain.</span>
		</FormField>
	</div>
{:else if prop.schema.type == 'array'}
	<div class="flex items-center">
		<p>{prop.name}</p>
		<IconButton
			type="button"
			class="material-icons"
			on:click={() => {
				// Append without .push() as we need to trigger a reactive update
				store = [...store, getInitialPropValue(prop.schema.items)];
			}}>add</IconButton
		>
	</div>
	<ol>
		{#each store as item, index}
			<li>
				<div class="flex items-center">
					<IconButton
						type="button"
						class="material-icons"
						on:click={() => {
							store.splice(index, 1);
							// Need to reassign variable to trigger a reactive update
							store = store;
						}}>clear</IconButton
					>
					<svelte:self
						bind:store={item}
						prop={propFromSchema(prop.name, prop.schema, prop.schema.items)}
					/>
				</div>
			</li>
		{/each}
	</ol>
{:else if prop.schema.type == 'object'}
	<p>{prop.name}:</p>
	<ObjectForm bind:store spec={prop.schema} />
{:else}
	<h3>Error: unsupported type {prop.schema.type}</h3>
{/if}
