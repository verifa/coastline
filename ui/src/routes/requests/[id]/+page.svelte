<script lang="ts">
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

	requestStore.get(`/requests/${$page.params.id}`);
	requestStore.subscribe((value) => {
		if (!value.requested || value.fetching) {
			return;
		} else if (!value.ok) {
			console.log('something wrong');
		}
	});

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
	}
	function handleReject() {
		const review: NewReview = {
			status: 'reject',
			type: 'user'
		};
		reviewStore.post(`/requests/${$requestStore.data?.id}/review`, {}, review);
	}
</script>

<Breadcrumb {links} page="Request" />

<h1>Request</h1>

{#if $requestStore.fetching}
	<Loading text="Loading" />
{:else if $requestStore.ok && $requestStore.data}
	<div class="flex items-start">
		<div class="flex-1">
			<div class="dropdown">
				<label tabindex="0" class="btn btn-outline">Review</label>
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
		<div class="flex-1">
			<!-- <h2 class="m-0">Reviews</h2> -->
			<div>
				{#if $requestStore.data.reviews}
					<ul class="steps steps-vertical">
						{#each $requestStore.data.reviews as review}
							<li class="step step-primary">{review.type} | {review.status}</li>
						{/each}
					</ul>
				{/if}
			</div>
		</div>
	</div>

	<h2>Spec</h2>
	<div class="mockup-code">
		<pre><code>// TODO: show spec here.</code></pre>
		<pre><code
				>// Do we want to render based on the OpenAPI spec for which the request was created?</code
			></pre>
		<pre><code
				>// Or render based on the current (if it exists) OpenAPI spec for the request template?</code
			></pre>
		<pre><code>// The request template, if it exists, might have changed...</code></pre>
	</div>
	<!-- <RequestObjectForm schemaObj={$store.data.spec} /> -->
{/if}
