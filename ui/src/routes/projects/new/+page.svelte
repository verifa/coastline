<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { createHttpStore } from '$lib/http/store';
	import type { components } from '$lib/oapi/spec';
	import Button from '@smui/button';
	import Textfield from '@smui/textfield';
	import { writable } from 'svelte/store';

	type NewProject = components['schemas']['NewProject'];
	type Project = components['schemas']['Project'];

	const createProjectStore = createHttpStore<Project>();

	const store = writable<NewProject>({
		name: ''
	});

	createProjectStore.subscribe((value) => {
		if (value.ok) {
			goto(`${base}`);
		} else {
			// TODO: something went wrong
		}
	});

	function handleOnSubmit() {
		createProjectStore.post('/projects', {}, $store);
	}
</script>

<h1>Create project</h1>

<form id="project" on:submit|preventDefault={handleOnSubmit}>
	<label for="name">Name:</label><br />
	<Textfield type="text" id="name" name="name" bind:value={$store.name} /><br />

	<Button type="submit" form="project">Submit</Button>
</form>
