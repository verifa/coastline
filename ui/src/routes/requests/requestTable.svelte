<script lang="ts">
	import type { components } from '$lib/oapi/spec';

	type Request = components['schemas']['Request'];
	export let requests: Request[] = [];
</script>

<div class="overflow-x-auto w-full">
	<table class="table w-full">
		<!-- head -->
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
					<th>
						<span class="badge badge-lg">{request.status}</span>
					</th>
					<td>
						{request.type}
					</td>
					<td>
						{request.project.name}
					</td>
					<td>{request.service.name}</td>
					<th>
						<button class="btn btn-sm">Action</button>
					</th>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
