<script lang="ts">
	import type { Property } from './spec';

	export let store: any;
	export let prop: Property;

	let placeholder: string = prop.is_required ? '<required>' : '<optional>';
</script>

{#if prop.schema.enum}
	<select
		class="inline px-1 bg-inherit focus:outline-none text-white text-right invalid:text-gray-500"
		bind:value={store}
		required
	>
		{#if !prop.schema.default}
			<option value="" disabled selected>select</option>
		{/if}
		{#each prop.schema.enum as item}
			<option value={item}>{item}</option>
		{/each}
	</select>
{:else if prop.schema.type == 'string'}
	<input
		class="inline py-0 px-2 bg-inherit placeholder:text-gray-500 border-2 border-transparent focus:border-inherit focus:outline-none"
		{placeholder}
		bind:value={store}
		type="text"
		required={prop.is_required}
	/>
{:else if prop.schema.type == 'number'}
	<input
		class="inline py-0 px-2 bg-inherit border-2 border-transparent focus:border-inherit focus:outline-none"
		{placeholder}
		bind:value={store}
		type="number"
		required={prop.is_required}
	/>
{:else if prop.schema.type == 'integer'}
	<input
		class="inline py-0 px-2 bg-inherit border-2 border-transparent focus:border-inherit focus:outline-none"
		{placeholder}
		bind:value={store}
		step="1"
		type="number"
		required={prop.is_required}
	/>
{:else if prop.schema.type == 'boolean'}
	<input
		type="checkbox"
		class="inline-block align-middle focus:outline-none toggle toggle-sm border-primary checked:bg-success"
		bind:checked={store}
	/>
{:else}
	<span class="px-2 border-2 border-error">Error: unsupported type {prop.schema.type}</span>
{/if}
