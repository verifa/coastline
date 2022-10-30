<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { createHttpStore } from '$lib/http/store';
	import Icon from '$lib/icons/Icon.svelte';
	import type { components } from '$lib/oapi/gen/types';
	import { writable } from 'svelte/store';

	type NewService = components['schemas']['NewService'];
	type Service = components['schemas']['Service'];

	const createServiceStore = createHttpStore<Service>();

	const store = writable<NewService>({
		name: '',
		labels: {}
	});

	interface Label {
		key: string;
		value: string;
	}
	let labels: Label[] = [];

	createServiceStore.subscribe((value) => {
		if (value.ok) {
			goto(`${base}/services`);
		} else {
			// TODO: something went wrong
		}
	});

	store.subscribe((value) => {
		console.log(value);
	});

	function handleOnSubmit() {
		$store.labels = labels.reduce((obj, label) => {
			return { ...obj, [label.key]: label.value };
		}, {});

		// console.log('labels: ', $store.labels);
		// console.log('test: ', { a: 'a', b: 'b' });
		createServiceStore.post('/services', {}, $store);
	}

	function handleAddLabel() {
		labels = [
			...labels,
			{
				key: '',
				value: ''
			}
		];
		console.log('labels: ', labels);
	}
</script>

<h1>Create service</h1>

<form id="service" on:submit|preventDefault={handleOnSubmit}>
	<div class="flex flex-col space-y-4 w-full max-w-md">
		<div class="form-control">
			<label for="name" class="label">
				<span class="label-text">Name</span>
			</label>
			<input
				id="name"
				type="text"
				placeholder="Service name"
				class="input input-bordered"
				bind:value={$store.name}
				required
			/>
		</div>
		<div>
			<div class="flex flex-col">
				<!-- TODO: this does not bind at the moment!! -->
				{#each labels as label}
					<div class="flex  items-center space-x-4">
						<button type="button" on:click={() => console.log('TODO!')}
							><Icon name="x-mark-mini" class="w-5" /></button
						>
						<div class="flex">
							<input placeholder="key" bind:value={label.key} required />
							<input placeholder="value" bind:value={label.value} required />
						</div>
					</div>
				{/each}
			</div>
			<button type="button" class="btn btn-sm btn-outline" on:click={handleAddLabel}
				>Add label</button
			>
		</div>
		<div>
			<button type="submit" class="btn btn-primary {$createServiceStore.fetching ? 'loading' : ''} "
				>Submit</button
			>
		</div>
	</div>
</form>
