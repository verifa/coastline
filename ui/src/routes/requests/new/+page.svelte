<script lang="ts">
	import { createHttpStore } from '$lib/http/store';
	import { getRequestSpecs } from '$lib/oapi/parse';
	import type { RequestSpec } from '$lib/oapi/parse';

	import Autocomplete from '@smui-extra/autocomplete';

	import type { OpenAPI3 } from 'openapi-typescript';
	import type { components } from '$lib/oapi/spec';
	import ObjectForm from './objectForm.svelte';
	import { writable } from 'svelte/store';
	import Button, { Label } from '@smui/button';

	type ProjectsResp = components['schemas']['ProjectsResp'];
	type ServicesResp = components['schemas']['ServicesResp'];
	type NewRequest = components['schemas']['NewRequest'];

	const projectStore = createHttpStore<ProjectsResp>();
	const serviceStore = createHttpStore<ServicesResp>();
	const requestsSpecStore = createHttpStore<OpenAPI3>();
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

	let specItems: Item[] = [];
	let selectedSpec: Item;

	let projectItems: Item[] = [];
	let selectedProject: Item;

	let serviceItems: Item[] = [];
	let selectedService: Item;

	requestsSpecStore.subscribe((value) => {
		if (value.ok && value.data) {
			specs = getRequestSpecs(value.data);
			specs.forEach((spec, index) => {
				specItems.push({ index, id: spec.type, label: spec.type });
			});
		}
	});

	projectStore.subscribe((value) => {
		if (value.ok && value.data) {
			value.data.projects.forEach((project, index) => {
				projectItems.push({ index, id: project.id, label: project.name });
			});
		}
	});

	serviceStore.subscribe((value) => {
		if (value.ok && value.data) {
			value.data.services.forEach((service, index) => {
				serviceItems.push({ index, id: service.id, label: service.name });
			});
		}
	});

	type Item = {
		index: number;
		id: string;
		label: string;
	};
</script>

<h1>New Request</h1>

<form
	on:submit|preventDefault={() => {
		$requestStore.project_id = selectedProject.id;
		$requestStore.service_id = selectedService.id;
		$requestStore.type = specs[selectedSpec.index].type;
		console.log($requestStore);
		// TODO: actually submit it!!
	}}
>
	<div class="flex flex-col space-y-4">
		{#if $projectStore.fetching}
			<h2>Loading projects</h2>
		{:else if $projectStore.ok}
			<Autocomplete
				options={projectItems}
				getOptionLabel={(item) => (item ? item.label : '')}
				bind:value={selectedProject}
				label="Project"
			/>
		{/if}

		{#if $serviceStore.fetching}
			<h2>Loading services</h2>
		{:else if $serviceStore.ok}
			<Autocomplete
				options={serviceItems}
				getOptionLabel={(item) => (item ? item.label : '')}
				bind:value={selectedService}
				label="Service"
			/>
		{/if}

		{#if $requestsSpecStore.fetching}
			<h2>Loading</h2>
		{:else if $requestsSpecStore.ok}
			<Autocomplete
				options={specItems}
				getOptionLabel={(item) => (item ? item.label : '')}
				bind:value={selectedSpec}
				label="Type"
			/>
			<!-- TODO: when this changes we need to reset the store... -->

			{#if selectedSpec && selectedSpec.index >= 0}
				<h2>Spec</h2>
				<ObjectForm bind:store={$requestStore.spec} spec={specs[selectedSpec.index].spec} />
				<Button variant={'raised'} class="w-40">
					<Label>Submit</Label>
				</Button>
			{/if}
		{/if}
	</div>
</form>
