<script lang="ts">
	import { createHttpStore } from '$lib/http/store';
	import { getRequestSpecs } from '$lib/oapi/parse';
	import type { RequestSpec } from '$lib/oapi/parse';

	import Autocomplete from '@smui-extra/autocomplete';

	import type { OpenAPI3 } from 'openapi-typescript';
	import type { components } from '$lib/oapi/spec';

	type ProjectsResp = components['schemas']['ProjectsResp'];
	type ServicesResp = components['schemas']['ServicesResp'];

	const projectStore = createHttpStore<ProjectsResp>();
	const serviceStore = createHttpStore<ServicesResp>();
	const specStore = createHttpStore<OpenAPI3>();

	projectStore.get('/projects');
	serviceStore.get('/services');
	specStore.get('/requestsspec');

	let specs: RequestSpec[];
	let selectedRequest: number;
	
	let projectItems: Item[];
	let selectedProject: Item;

	let serviceItems: Item[];

	specStore.subscribe((value) => {
		if (value.ok && value.data) {
			specs = getRequestSpecs(value.data);
		}
	});

	projectStore.subscribe((value) => {
		if (value.ok && value.data) {
			for(const project in value.data.projects) {
				console.log("Foop")
				console.log(project.name)
				projectItems.push({id: project.id, label: project.name})
				console.log(projectItems)
			}
		}
	});

	serviceStore.subscribe((value) => {
		if (value.ok && value.data) {
			value.data.services.forEach(service => {
				serviceItems.push({id: service.id, label: service.name})
			});
		}
	});

	type Item = {
		id: string;
		label: string;
	}
</script>

<h1>New Request</h1>

{#if $projectStore.fetching}
	<h2>Loading projects</h2>
{:else if $projectStore.ok}
	<Autocomplete
		options={projectItems}
		getOptionLabel={(item) => item ? item.label : 'N/A'}
		bind:value={selectedProject}
		label="Project"
	/>
{/if}

{#if $specStore.fetching}
	<h2>Loading</h2>
{:else if $specStore.ok}
	<h2>Form</h2>
	<label for="type">Project:</label>
	<Autocomplete
		getOptionLabel={(item) => {
			console.log('item: ', item);
			return item ? item.type : '';
		}}
		options={specs}
		label="Standard"
	/>
	<!-- bind:value={valueStandard} -->
	<select
		name="type"
		id="type"
		form="request"
		bind:value={selectedRequest}
		on:change={() => console.log(selectedRequest)}
	>
		<option value="-1" disabled selected>Select request</option>
		{#each specs as spec, index}
			<option value={index}>{spec.type}</option>
		{/each}
	</select>
	{#if selectedRequest >= 0}
		<h3>Render custom form here...</h3>
		<pre>{JSON.stringify(specs[selectedRequest].spec)}</pre>
	{/if}
{/if}
