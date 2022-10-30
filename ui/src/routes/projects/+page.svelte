<script lang="ts">
	import { base } from '$app/paths';
	import Breadcrumb from '$lib/Breadcrumb.svelte';
	import { createHttpStore } from '$lib/http/store';

	import type { components } from '$lib/oapi/gen/types';

	type ProjectResp = components['schemas']['ProjectsResp'];

	const projectsStore = createHttpStore<ProjectResp>();

	projectsStore.get('/projects');
</script>

<Breadcrumb page="Projects" />

<h1>Projects</h1>

{#if $projectsStore.fetching}
	<h2>Loading projects</h2>
{:else if $projectsStore.ok && $projectsStore.data}
	<div class="overflow-x-auto w-full">
		<table class="table table-auto w-full">
			<thead>
				<tr>
					<th>Name</th>
					<th />
				</tr>
			</thead>
			<tbody>
				{#each $projectsStore.data.projects as project}
					<tr>
						<td>
							{project.name}
						</td>
						<th>
							<a href="{base}/requests/new?project_id={project.id}" role="button" class="btn btn-sm"
								>New Request</a
							>
						</th>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}

<a href="{base}/projects/new" class="btn btn-primary">New Project</a>
