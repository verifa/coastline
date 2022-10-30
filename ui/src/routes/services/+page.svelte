<script lang="ts">
	import { base } from '$app/paths';
	import Breadcrumb from '$lib/Breadcrumb.svelte';
	import { createHttpStore } from '$lib/http/store';

	import type { components } from '$lib/oapi/gen/types';

	type ServicesResp = components['schemas']['ServicesResp'];

	const servicesStore = createHttpStore<ServicesResp>();

	servicesStore.get('/services');
</script>

<Breadcrumb page="Services" />

<h1>Services</h1>

{#if $servicesStore.fetching}
	<h2>Loading services</h2>
{:else if $servicesStore.ok && $servicesStore.data}
	<div class="overflow-x-auto w-full">
		<table class="table table-auto w-full">
			<thead>
				<tr>
					<th>Name</th>
					<th />
				</tr>
			</thead>
			<tbody>
				{#each $servicesStore.data.services as service}
					<tr>
						<td>
							{service.name}
						</td>
						<th>
							<a href="{base}/requests/new" role="button" class="btn btn-sm">New Request</a>
						</th>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

<a href="{base}/services/new" class="btn btn-primary">New Service</a>
