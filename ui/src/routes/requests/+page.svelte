<script lang="ts">
	import { createHttpStore } from '$lib/http/store';

	import type { components } from '$lib/oapi/spec';

	type RequestResp = components['schemas']['RequestsResp'];

	const requestsStore = createHttpStore<RequestResp>();

	requestsStore.get('/requests');
</script>

<h1>Requests</h1>

{#if $requestsStore.fetching}
	<h2>Loading requests</h2>
{:else if $requestsStore.ok && $requestsStore.data}
	{#each $requestsStore.data.requests as request}
		<h3>{request.id}</h3>
		<ul>
			<li>{request.project?.name}</li>
			<li>{request.service?.name}</li>
			<li>{request.requested_by}</li>
			<li>{request.type}</li>
			<li>{JSON.stringify(request.spec)}</li>
		</ul>
	{/each}
{/if}
