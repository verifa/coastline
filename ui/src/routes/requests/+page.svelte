<script lang="ts">
	import { base } from '$app/paths';
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
	<div class="grid grid-cols-3">
		{#each $requestsStore.data.requests as request}
			<div class="card w-96 bg-base-100 shadow-xl">
				<div class="card-body">
					<h3 class="card-title">{request.type}</h3>
					<div class="badge badge-warning">Needs approval</div>
					<ol class="list-none">
						<li>Project: {request.project?.name}</li>
						<li>Service: {request.service?.name}</li>
						<li>Requested by: {request.requested_by}</li>
						<!-- <li>{JSON.stringify(request.spec)}</li> -->
					</ol>
				</div>
			</div>
		{/each}
	</div>
{/if}

<a href="{base}/requests/new" class="btn btn-primary">New Request</a>
