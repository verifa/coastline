<script lang="ts">
	import { getInitialPropValue, type Property } from './spec';

	export let store: any;
	export let prop: Property;
	export let depth: number;

	const indent: string = '  '.repeat(depth);

	let showDescription: boolean = false;

	function toggleDescription() {
		showDescription = !showDescription;
	}
</script>

{#if showDescription && prop.description}
	<!-- TODO: move this out -->
	<pre><code><span class="text-gray-500">{indent + '// ' + prop.description}</span></code></pre>
{/if}

{#if prop.schema.type == 'array'}
	<span>
		<span class="hover:cursor-pointer" on:click={toggleDescription}>{indent + prop.name}</span
		>{': []'}
		<span
			on:click={() => {
				// Append without .push() as we need to trigger a reactive update
				store = [...store, getInitialPropValue(prop.schema.items)];
				console.log('wooow: ', store);
			}}
			class="text-success hover:cursor-pointer">+</span
		>
	</span>
{:else}
	<span class="hover:cursor-pointer" on:click={toggleDescription}>{indent + prop.name}</span>: <slot
	/>
{/if}
