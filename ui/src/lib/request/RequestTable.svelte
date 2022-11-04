<script lang="ts">
	import { base } from '$app/paths';
	import type { components } from '$lib/oapi/gen/types';
	import StatusBadge from './StatusBadge.svelte';

	type Request = components['schemas']['Request'];
	export let requests: Request[] = [];
</script>

<div class="overflow-x-auto w-full">
	<table class="table table-auto w-full">
		<thead>
			<tr>
				<th>Status</th>
				<th>Request</th>
				<th>Project</th>
				<th>Service</th>
				<th />
			</tr>
		</thead>
		<tbody>
			{#each requests as request}
				<tr>
					<td>
						<StatusBadge {request} />
					</td>
					<td>
						{request.kind}
					</td>
					<td>
						{request.project.name}
					</td>
					<td>{request.service.name}</td>
					<th>
						<a href="{base}/requests/{request.id}" role="button" class="btn btn-sm">View</a>
					</th>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
