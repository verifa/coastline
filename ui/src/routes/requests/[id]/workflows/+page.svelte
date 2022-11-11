<script lang="ts">
	import Icon from '$lib/icons/Icon.svelte';
	import type { components } from '$lib/oapi/gen/types';
	import { getContext } from 'svelte';
	import type { getRequestFunc } from '../request';

	type Trigger = components['schemas']['Trigger'];
	type Workflow = components['schemas']['Workflow'];

	const getRequest: getRequestFunc = getContext('request');
	const request = getRequest();

	const workflows = request.triggers.reduce((prev: Workflow[], cur: Trigger) => {
		prev.push(...cur.workflows);
		return prev;
	}, []);
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
			{#each workflows as workflow}
				<tr>
					<th>
						{#if workflow.error === ''}
							<Icon name="check-circle-solid" class="w-10 h-10 text-success" />
						{:else}
							<Icon name="x-circle-solid" class="w-10 h-10 text-error" />
						{/if}
					</th>
					<td>{workflow.id}</td>
					<td>
						{#if workflow.error === ''}
							{JSON.stringify(workflow.output, undefined, 2)}
						{:else}
							Error: {workflow.error}
						{/if}
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
