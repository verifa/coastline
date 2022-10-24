<script lang="ts">
	import type { components } from '$lib/oapi/gen/types';
	import { createHttpStore } from '$lib/http/store';
	import { session } from '$lib/session/store';
	import { base } from '$app/paths';
	import Loading from '$lib/Loading.svelte';

	type ProjectsResp = components['schemas']['ProjectsResp'];
	type ServicesResp = components['schemas']['ServicesResp'];
	type RequestsResp = components['schemas']['RequestsResp'];

	const projectStore = createHttpStore<ProjectsResp>();
	const serviceStore = createHttpStore<ServicesResp>();
	const requestStore = createHttpStore<RequestsResp>();

	projectStore.get('/projects');
	serviceStore.get('/services');
	requestStore.get('/requests');
</script>

<h1>Welcome {$session.user?.email}!</h1>

{#if $projectStore.ok && $projectStore.data}
	{#if $projectStore.data.projects.length === 0}
		<h2>No projects...</h2>
	{:else}
		<h2>Project list</h2>
		<ul>
			{#each $projectStore.data.projects as project}
				<li>{project.name}</li>
			{/each}
		</ul>
	{/if}
	<a href="{base}/projects/new" class="btn btn-primary">New project</a>
{:else if $projectStore.error}
	<h2>Error: {$projectStore.error.message}</h2>
{:else if $projectStore.fetching}
	<h2>Loading projects...</h2>
{/if}

{#if $serviceStore.ok && $serviceStore.data}
	{#if $serviceStore.data.services.length === 0}
		<h2>No services...</h2>
	{:else}
		<h2>Services list</h2>
		<ul>
			{#each $serviceStore.data.services as service}
				<li>{service.name}</li>
			{/each}
		</ul>
	{/if}
{:else if $serviceStore.error}
	<h2>Error: {$serviceStore.error.message}</h2>
{:else if $serviceStore.fetching}
	<h2>Loading services...</h2>
{/if}

{#if $requestStore.ok && $requestStore.data}
	{#if $requestStore.data.requests.length === 0}
		<h2>No requests...</h2>
	{:else}
		<h2>Requests list</h2>
		<ul>
			{#each $requestStore.data.requests as request}
				<li>{request.type} in {request.project?.name}</li>
			{/each}
		</ul>
	{/if}
{:else if $requestStore.fetching}
	<Loading text="Loading" />
{/if}

<a href="{base}/requests" class="btn btn-primary">Requests</a>
<a href="{base}/requests/new" class="btn btn-primary">New Request</a>
