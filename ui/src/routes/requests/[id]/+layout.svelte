<script lang="ts">
	import { setContext } from 'svelte';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import Breadcrumb from '$lib/Breadcrumb.svelte';
	import { createHttpStore } from '$lib/http/store';
	import Loading from '$lib/Loading.svelte';
	import type { components } from '$lib/oapi/gen/types';
	import StatusBadge from '$lib/request/StatusBadge.svelte';

	type Request = components['schemas']['Request'];
	type NewReview = components['schemas']['NewReview'];

	const requestStore = createHttpStore<Request>();
	const reviewStore = createHttpStore<Request>();

	setContext('request', () => $requestStore.data);

	let currentId = $page.params.id;
	function getRequest() {
		requestStore.get(`/requests/${$page.params.id}`);
	}
	// Make initial call to get requests
	getRequest();

	// Subscribe to page in case request id changes
	page.subscribe((value) => {
		if (currentId != value.params.id) {
			currentId = value.params.id;
			getRequest();
		}
	});

	requestStore.subscribe((value) => {
		if (!value.requested || value.fetching) {
			return;
		} else if (!value.ok) {
			console.log('something wrong');
		}
	});

	interface tab {
		name: string;
		url: string;
	}

	const tabs = (): tab[] => {
		const baseUrl = `${base}/requests/${$requestStore.data?.id}`;
		return [
			{
				name: 'Spec',
				url: baseUrl
			},
			{
				name: 'Events',
				url: `${baseUrl}/events`
			},
			{
				name: 'Workflows',
				url: `${baseUrl}/workflows`
			},
			{
				name: 'Triggers',
				url: `${baseUrl}/triggers`
			}
		];
	};

	const links = [
		{
			path: '/requests',
			text: 'Requests'
		}
	];

	function handleApprove() {
		const review: NewReview = {
			status: 'approve',
			type: 'user'
		};
		reviewStore.post(`/requests/${$requestStore.data?.id}/review`, {}, review);
		getRequest();
	}
	function handleReject() {
		const review: NewReview = {
			status: 'reject',
			type: 'user'
		};
		reviewStore.post(`/requests/${$requestStore.data?.id}/review`, {}, review);
		getRequest();
	}
</script>

<Breadcrumb {links} page="Request" />

<h1>Request</h1>

{#if $requestStore.fetching}
	<Loading text="Loading" />
{:else if $requestStore.ok && $requestStore.data}
	<div class="flex flex-col">
		<div class="dropdown">
			<label for="" tabindex="0" class="btn btn-outline">Review</label>
			<ul tabindex="0" class="dropdown-content menu p-2 shadow bg-base-100 rounded-box w-52">
				<li><button on:click={handleApprove} class="btn btn-success">Approve</button></li>
				<li><button on:click={handleReject} class="btn btn-error">Reject</button></li>
			</ul>
		</div>
		<div class="overflow-x-auto">
			<table class="table table-auto max-w-md">
				<tbody>
					<tr>
						<th>Project</th>
						<td>{$requestStore.data.project.name}</td>
					</tr>
					<tr>
						<th>Service</th>
						<td>{$requestStore.data.service.name}</td>
					</tr>
					<tr>
						<th>Status</th>
						<td><StatusBadge request={$requestStore.data} /></td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>

	<div class="flex flex-col space-y-8">
		<div class="tabs">
			{#each tabs() as tab}
				<a
					href={tab.url}
					class="tab tab-bordered no-underline {$page.url.pathname == tab.url ? 'tab-active' : ''}"
					>{tab.name}</a
				>
			{/each}
		</div>
		<slot request={$requestStore.data} />
	</div>
{/if}
