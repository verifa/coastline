<script lang="ts">
	import { createHttpStore } from '$lib/http/store';
	import { user } from '$lib/auth/store'
	import { dataset_dev } from 'svelte/internal';

	interface Project {
		name: string
	}

	interface ProjectResponse {
		projects: Project[]
	}

	const projectStore = createHttpStore<ProjectResponse>()
	const logoutStore = createHttpStore()

	projectStore.get("/projects")

	function handleLogout() {
		logoutStore.get("/logout")
	}
</script>

<h1>Welcome {$user.email}!</h1>
{#if $projectStore.ok && $projectStore.data}
	{#if $projectStore.data.projects.length === 0}
	<h2>No projects...</h2>
	{:else}
	<ul>
		{#each $projectStore.data.projects as project}
			<li>{project.name}</li>
		{/each}
	</ul>
	{/if}
{:else if $projectStore.error}
<h2>Error: {$projectStore.error.message}</h2>
{:else if $projectStore.fetching}
<h2>Loading...</h2>
{/if}


<button on:click={handleLogout}>Logout</button>
