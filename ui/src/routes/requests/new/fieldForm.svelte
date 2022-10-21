<script lang="ts">
	import { getInitialPropValue, propFromSchema, type Property } from './spec';
	import ObjectForm from './objectForm.svelte';
	import Icon from '$lib/icons/Icon.svelte';

	export let prop: Property;
	export let store: any;
</script>

{#if prop.schema.enum}
	<!-- 
    key field sets a function that makes sure each value is a string
    which is a requirement of this component.
 -->
	<div class="form-control w-full max-w-xs">
		<label class="label">
			<span class="label-text">{prop.name}</span>
			<select class="select select-bordered" bind:value={store}>
				{#each prop.schema.enum as item}
					<option>{item}</option>
				{/each}
			</select>
		</label>
	</div>
{:else if prop.schema.type == 'string'}
	<div class="form-control w-full max-w-xs">
		<label class="label">
			<span class="label-text">{prop.name}</span>
			<input
				type="text"
				placeholder={prop.name}
				class="input input-bordered w-full max-w-xs"
				bind:value={store}
				required={prop.is_required}
			/>
		</label>
	</div>
{:else if prop.schema.type == 'number'}
	<div class="form-control w-full max-w-xs">
		<label class="label">
			<span class="label-text">{prop.name}</span>
			<input
				type="number"
				placeholder={prop.name}
				class="input input-bordered w-full max-w-xs"
				bind:value={store}
				required={prop.is_required}
			/>
		</label>
	</div>
{:else if prop.schema.type == 'integer'}
	<div class="form-control w-full max-w-xs">
		<label class="label">
			<span class="label-text">{prop.name}</span>
			<input
				type="number"
				placeholder={prop.name}
				class="input input-bordered w-full max-w-xs"
				bind:value={store}
				required={prop.is_required}
			/>
		</label>
	</div>
{:else if prop.schema.type == 'boolean'}
	<div class="form-control">
		<label class="label">
			<span class="label-text w-20">{prop.name}</span>
			<input type="checkbox" class="toggle" bind:checked={store} />
		</label>
	</div>
{:else if prop.schema.type == 'array'}
	<div class="">
		<label class="flex items-center space-x-4">
			<span class="label-text">{prop.name}</span>
			<button
				type="button"
				class="btn btn-circle btn-xs"
				on:click={() => {
					// Append without .push() as we need to trigger a reactive update
					store = [...store, getInitialPropValue(prop.schema.items)];
				}}
			>
				<Icon name="add-mini" class="w-5 h-5" />
			</button>
		</label>
	</div>
	<ol class="list-none">
		{#each store as item, index}
			<li>
				<div class="flex items-center space-x-4">
					<button
						type="button"
						class="btn btn-circle btn-xs"
						on:click={() => {
							store.splice(index, 1);
							// Need to reassign variable to trigger a reactive update
							store = store;
						}}
					>
						<Icon name="x-mark-mini" class="w-5 h-5" />
					</button>
					<svelte:self
						bind:store={item}
						prop={propFromSchema(prop.name, prop.schema, prop.schema.items)}
					/>
				</div>
			</li>
		{/each}
	</ol>
{:else if prop.schema.type == 'object'}
	<span class="label-text">{prop.name}:</span>
	<ObjectForm bind:store schemaObj={prop.schema} />
{:else}
	<h3>Error: unsupported type {prop.schema.type}</h3>
{/if}
