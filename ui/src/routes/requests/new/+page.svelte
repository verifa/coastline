<script lang="ts">
	import { createHttpStore } from '$lib/http/store';
	import { getRequestSpecs } from '$lib/oapi/parse';
	import type { RequestSpec } from '$lib/oapi/parse';

	import Autocomplete from '@smui-extra/autocomplete';

	import type { OpenAPI3 } from 'openapi-typescript';
	import type { components } from '$lib/oapi/spec';
	import SpecForm from './specForm.svelte';
	import { writable } from 'svelte/store';
	import Button, { Label } from '@smui/button';

	type ProjectsResp = components['schemas']['ProjectsResp'];
	type ServicesResp = components['schemas']['ServicesResp'];
	type NewRequest = components['schemas']['NewRequest'];

	const projectStore = createHttpStore<ProjectsResp>();
	const serviceStore = createHttpStore<ServicesResp>();
	const requestsSpecStore = createHttpStore<OpenAPI3>();
	const specStore = writable<{ [key: string]: any }>({});

	specStore.subscribe((value) => {
		console.log(value);
	});

	projectStore.get('/projects');
	serviceStore.get('/services');
	requestsSpecStore.get('/requestsspec');

	let specs: RequestSpec[];
	let selectedRequest: number;

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
		console.log($specStore);
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
				required={true}
				options={specItems}
				getOptionLabel={(item) => (item ? item.label : '')}
				bind:value={selectedSpec}
				label="Type"
			/>

			{#if selectedSpec && selectedSpec.index >= 0}
				<h2>Spec</h2>
				<SpecForm store={specStore} spec={specs[selectedSpec.index].spec} />
				<Button variant={'raised'} class="w-40">
					<Label>Submit</Label>
				</Button>
			{/if}
		{/if}
	</div>
</form>
