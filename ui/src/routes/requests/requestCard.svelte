<!-- THIS WAS A HACK WIP... Can be deleted, probably -->
<script lang="ts">
	import type { components } from '$lib/oapi/gen/types';

	type Request = components['schemas']['Request'];
	export let requests: Request[] = [];
</script>

<div class="grid grid-cols-3">
	{#each requests as request}
		<div class="card w-96 bg-base-100 shadow-xl">
			<div class="card-body">
				<h3 class="card-title">{request.type}</h3>
				<div class="badge badge-warning">Needs approval</div>
				<ol class="list-none">
					<li>Project: {request.project?.name}</li>
					<li>Service: {request.service?.name}</li>
					<li>Requested by: {request.requested_by}</li>
					<!-- <li>{JSON.stringify(request.spec)}</li> -->
				</ol>
			</div>
		</div>
	{/each}
</div>
