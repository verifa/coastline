<script lang="ts">
	import { createHttpStore } from '$lib/http/store';
	import { getRequestSpecs } from '$lib/oapi/parse';
	import type { RequestSpec } from '$lib/oapi/parse';

	import type { OpenAPI3 } from 'openapi-typescript';
	import type { components } from '$lib/oapi/spec';
	import ObjectForm from './objectForm.svelte';
	import { session } from '$lib/session/store';
	import { writable } from 'svelte/store';

	type ProjectsResp = components['schemas']['ProjectsResp'];
	type ServicesResp = components['schemas']['ServicesResp'];
	type NewRequest = components['schemas']['NewRequest'];
	type Request = components['schemas']['Request'];

	const projectStore = createHttpStore<ProjectsResp>();
	const serviceStore = createHttpStore<ServicesResp>();
	const requestsSpecStore = createHttpStore<OpenAPI3>();
	const requestsSubmitStore = createHttpStore<Request>();

	const requestStore = writable<NewRequest>({
		project_id: '',
		service_id: '',
		type: '',
		requested_by: '',
		spec: {}
	});

	requestStore.subscribe((value) => {
		console.log(value);
	});

	projectStore.get('/projects');
	serviceStore.get('/services');
	requestsSpecStore.get('/requestsspec');

	let specs: RequestSpec[];
	let selectedSpec: RequestSpec;

	requestsSpecStore.subscribe((value) => {
		if (value.ok && value.data) {
			specs = getRequestSpecs(value.data);
		}
	});

	requestsSubmitStore.subscribe((value) => {
		console.log(value);
	});

	function handleSubmit() {
		$requestStore.type = selectedSpec.type;
		$requestStore.requested_by = $session.user?.id ? $session.user?.id : 'anonymous';
		requestsSubmitStore.post('/requests', {}, $requestStore);
	}
</script>

<h1>New Request</h1>

<form on:submit|preventDefault={handleSubmit}>
	<div class="flex flex-col space-y-4">
		{#if $projectStore.fetching}
			<h2>Loading projects</h2>
		{:else if $projectStore.ok && $projectStore.data}
			<div class="form-control w-full max-w-xs">
				<label for="project" class="label">
					<span class="label-text">Project</span>
				</label>
				<select id="project" class="select select-bordered" bind:value={$requestStore.project_id}>
					<option disabled selected value={''}>Select project</option>
					{#each $projectStore.data.projects as project}
						<option value={project.id}>{project.name}</option>
					{/each}
				</select>
			</div>
		{/if}

		{#if $serviceStore.fetching}
			<h2>Loading services</h2>
		{:else if $serviceStore.ok && $serviceStore.data}
			<div class="form-control w-full max-w-xs">
				<label for="service" class="label">
					<span class="label-text">Service</span>
				</label>
				<select id="service" class="select select-bordered" bind:value={$requestStore.service_id}>
					<option disabled selected value={''}>Select service</option>
					{#each $serviceStore.data.services as service}
						<option value={service.id}>{service.name}</option>
					{/each}
				</select>
			</div>
		{/if}

		{#if $requestsSpecStore.fetching}
			<h2>Loading</h2>
		{:else if $requestsSpecStore.ok}
			<div class="form-control w-full max-w-xs">
				<label for="request" class="label">
					<span class="label-text">Request</span>
				</label>
				<select id="request" class="select select-bordered" bind:value={selectedSpec}>
					<option disabled selected value={undefined}>Select request</option>
					{#each specs as spec}
						<option value={spec}>{spec.type}</option>
					{/each}
				</select>
			</div>
			<!-- TODO: when this changes we need to reset the store... -->

			{#if selectedSpec}
				<h2>Spec</h2>
				<ObjectForm bind:store={$requestStore.spec} spec={selectedSpec.spec} />
				<div>
					<button class="btn btn-primary">Submit</button>
				</div>
			{/if}
		{/if}
	</div>
</form>
