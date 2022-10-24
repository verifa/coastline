<script lang="ts">
	import { base } from '$app/paths';

	interface Link {
		path: string;
		text: string;
	}

	export let links: Link[] = [];
	export let page: string;
</script>

<div class="text-sm breadcrumbs">
	<ul>
		<li><a href={base}>Home</a></li>
		{#each links as link, index}
			<li><a href={base + link.path}>{link.text}</a></li>
		{/each}
		<li>{page}</li>
	</ul>
</div>
