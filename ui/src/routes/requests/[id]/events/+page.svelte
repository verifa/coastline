<script lang="ts">
	import { getContext } from 'svelte';
	import type { getRequestFunc } from '../request';
	const getRequest: getRequestFunc = getContext('request');
	const request = getRequest();
</script>

<ul class="steps steps-vertical">
	{#each request.reviews as review}
		<li class="step step-primary">{review.type} | {review.status}</li>
	{/each}
</ul>
