<script lang="ts">
	import { createHttpStore } from '$lib/http/store';

	import type { SchemaObject, OpenAPI3 } from 'openapi-typescript';
	import type { components } from '$lib/oapi/gen/types';
	import { session } from '$lib/session/store';
	import { writable } from 'svelte/store';
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { page } from '$app/stores';

	import RequestObjectForm from '$lib/request/RequestObjectForm.svelte';
	import Loading from '$lib/Loading.svelte';

	type ProjectsResp = components['schemas']['ProjectsResp'];
	type ServicesResp = components['schemas']['ServicesResp'];
	type RequestTemplatesResp = components['schemas']['RequestTemplatesResp'];
	type NewRequest = components['schemas']['NewRequest'];
	type Request = components['schemas']['Request'];

	const projectStore = createHttpStore<ProjectsResp>();
	const serviceStore = createHttpStore<ServicesResp>();
	const requestTemplatesStore = createHttpStore<RequestTemplatesResp>();
	const templateSpecStore = createHttpStore<OpenAPI3>();
	const requestsSubmitStore = createHttpStore<Request>();

	const requestStore = writable<NewRequest>({
		project_id: '',
		service_id: '',
		type: '',
		status: 'pending_approval',
		spec: {}
	});

	projectStore.get('/projects');

	// TODO: remove this, but it's useful for debugging at the moment
	requestStore.subscribe((value) => {
		console.log(value);
	});

	requestsSubmitStore.subscribe((value) => {
		if (value.ok) {
			goto(`${base}/requests`);
		} else {
			console.log('store: ', value);
			console.log('data: ', value.data);
		}
	});

	function getSchemaObjFromOpenAPI(spec: OpenAPI3): SchemaObject {
		if (spec.components) {
			for (const key in spec.components.schemas) {
				return spec.components.schemas[key];
			}
		}
		return undefined;
	}

	function handleProjectChange() {
		serviceStore.get('/services');
	}

	function handleServiceChange() {
		requestTemplatesStore.get(`/services/${$requestStore.service_id}/templates`);
	}

	function handleRequestTemplateChange() {
		templateSpecStore.get(`/templates/${$requestStore.type}/openapi`);
	}

	function handleSubmit() {
		requestsSubmitStore.post('/requests', {}, $requestStore);
	}

	// Check search params
	const projectId = $page.url.searchParams.get('project_id');
	if (projectId) {
		$requestStore.project_id = projectId;
		handleProjectChange();
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
				<select
					id="project"
					class="select select-bordered"
					bind:value={$requestStore.project_id}
					on:change={handleProjectChange}
					required
				>
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
				<select
					id="service"
					class="select select-bordered"
					bind:value={$requestStore.service_id}
					on:change={handleServiceChange}
					required
				>
					<option disabled selected value={''}>Select service</option>
					{#each $serviceStore.data.services as service}
						<option value={service.id}>{service.name}</option>
					{/each}
				</select>
			</div>
		{/if}

		{#if $requestTemplatesStore.fetching}
			<h2>Loading</h2>
		{:else if $requestTemplatesStore.ok && $requestTemplatesStore.data}
			<div class="form-control w-full max-w-xs">
				<label for="request" class="label">
					<span class="label-text">Request</span>
				</label>
				<select
					id="request"
					class="select select-bordered"
					bind:value={$requestStore.type}
					on:change={handleRequestTemplateChange}
				>
					<option disabled selected value={undefined}>Select request</option>
					{#each $requestTemplatesStore.data.templates as template}
						<option value={template.type}>{template.type}</option>
					{/each}
				</select>
			</div>
			<!-- TODO: when this changes we need to reset the store... -->

			{#if $templateSpecStore.fetching}
				<Loading />
			{:else if $templateSpecStore.ok && $templateSpecStore.data}
				<!-- Uncomment this to inspect the OpenAPI JSON -->
				<!-- {JSON.stringify(selectedTemplate.spec)} -->
				<h2>Spec</h2>

				<div class="mockup-code not-prose relative">
					<RequestObjectForm
						bind:store={$requestStore.spec}
						depth={0}
						schemaObj={getSchemaObjFromOpenAPI($templateSpecStore.data)}
					/>
				</div>
				<!-- <ObjectForm bind:store={$requestStore.spec} schemaObj={selectedTemplate.spec} /> -->
				<div>
					<button class="btn btn-primary">Submit</button>
				</div>
			{/if}
			<!--
				Hardcode some error message for now if validation fails.
				TODO: We should build a component/store for handling these kinds of things.
			 -->
			{#if $requestsSubmitStore.requested && !$requestsSubmitStore.fetching && !$requestsSubmitStore.ok}
				<div class="alert alert-error shadow-lg">
					<div>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="stroke-current flex-shrink-0 h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
							/></svg
						>
						<span>Status: {$requestsSubmitStore.status}: {$requestsSubmitStore.text}</span>
					</div>
				</div>
			{/if}
		{/if}
	</div>
</form>
