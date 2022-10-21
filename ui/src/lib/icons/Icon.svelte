<script lang="ts">
	import { icons } from '$lib/icons/icons';

	export let name: string;
	export { className as class };
	let className: string;

	let displayIcon = icons.find((e) => e.name === name);
</script>

{#if displayIcon}
	<svg
		class={className}
		xmlns="http://www.w3.org/2000/svg"
		fill={'currentColor'}
		viewBox="0 0 {displayIcon.box} {displayIcon.box}">{@html displayIcon.svg}</svg
	>
{:else}
	<h3>ERROR: couldn't find icon</h3>
{/if}
