<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { createHttpStore } from '$lib/http/store';
	import type { components } from '$lib/oapi/gen/types';
	import { writable } from 'svelte/store';

	type NewProject = components['schemas']['NewProject'];
	type Project = components['schemas']['Project'];

	const createProjectStore = createHttpStore<Project>();

	const store = writable<NewProject>({
		name: ''
	});

	createProjectStore.subscribe((value) => {
		if (value.ok) {
			goto(`${base}/projects`);
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
	<div class="flex flex-col space-y-4 w-full max-w-md">
		<div class="form-control flex-1">
			<label for="name" class="label">
				<span class="label-text">Name</span>
			</label>
			<input
				id="name"
				type="text"
				placeholder="Project name"
				class="input input-bordered"
				bind:value={$store.name}
				required
			/>
		</div>
		<div>
			<button type="submit" class="btn btn-primary {$createProjectStore.fetching ? 'loading' : ''} "
				>Submit</button
			>
		</div>
	</div>
</form>
