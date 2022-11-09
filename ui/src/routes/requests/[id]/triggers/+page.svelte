<script lang="ts">
	import Icon from '$lib/icons/Icon.svelte';
	import type { components } from '$lib/oapi/gen/types';
	import { getContext } from 'svelte';
	import type { getRequestFunc } from '../request';

	type Trigger = components['schemas']['Trigger'];

	const getRequest: getRequestFunc = getContext('request');
	const request = getRequest();

	function isSuccessful(trigger: Trigger): boolean {
		return trigger.tasks.filter((task) => task.error !== '').length === 0;
	}
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
						{#if isSuccessful(trigger)}
							<Icon name="check-circle-solid" class="w-10 h-10 text-success" />
						{:else}
							<Icon name="x-circle-solid" class="w-10 h-10 text-error" />
						{/if}
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
