<script lang="ts">
	import Icon from '$lib/icons/Icon.svelte';
	import { getContext } from 'svelte';
	import type { getRequestFunc } from '../request';
	const getRequest: getRequestFunc = getContext('request');
	const request = getRequest();
</script>

<div class="overflow-x-auto w-full">
	<table class="table w-full">
		<thead>
			<tr>
				<th>Status</th>
				<th>ID</th>
				<th />
			</tr>
		</thead>
		<tbody>
			{#each request.triggers as trigger}
				<tr>
					<th>
						<Icon name="check-circle-solid" class="w-10 h-10 text-success" />
					</th>
					<td>{trigger.id}</td>
					<th>
						<button class="btn btn-ghost btn-xs">details</button>
					</th>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
