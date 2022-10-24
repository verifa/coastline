<script lang="ts">
	import { base } from '$app/paths';
	import Breadcrumb from '$lib/Breadcrumb.svelte';
	import { createHttpStore } from '$lib/http/store';

	import type { components } from '$lib/oapi/gen/types';
	import RequestTable from '$lib/request/RequestTable.svelte';

	type RequestResp = components['schemas']['RequestsResp'];

	const requestsStore = createHttpStore<RequestResp>();

	requestsStore.get('/requests');
</script>

<Breadcrumb page="Requests" />

<h1>Requests</h1>

{#if $requestsStore.fetching}
	<h2>Loading requests</h2>
{:else if $requestsStore.ok && $requestsStore.data}
	<RequestTable requests={$requestsStore.data.requests} />
{/if}

<a href="{base}/requests/new" class="btn btn-primary">New Request</a>
