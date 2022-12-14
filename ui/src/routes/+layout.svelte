<script lang="ts">
	import '../app.postcss';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	import { session } from '$lib/session/store';
	import { createHttpStore } from '$lib/http/store';
	import type { User } from '$lib/session/session';
	import { page } from '$app/stores';

	import { base } from '$app/paths';
	import NavBar from './NavBar.svelte';

	const authStore = createHttpStore<User>();

	$: isLoginPage = (): boolean => {
		return $page.url.pathname === `${base}/login`;
	};

	authStore.subscribe((resp) => {
		// Do nothing if request is in progress
		if (resp.fetching) {
			return;
		}
		if (resp.ok && resp.data) {
			session.login(resp.data);
		} else if (resp.status === 401) {
			session.logout();
		}
	});
	authStore.get('/userinfo');

	onMount(() => {
		// Subscribe to the session store to handle events if user is not authenticated
		session.subscribe((session) => {
			if (!session.initialized) {
				return;
			} else if (!session.authenticated) {
				if (!isLoginPage()) {
					goto(`${base}/login`);
				}
			} else if (session.initialized && session.authenticated) {
				if (isLoginPage()) {
					goto(`${base}`);
				}
			}
		});
	});
</script>

<!-- Configure TailwindCSS typographic defaults for whole site -->
<div class="prose mx-auto max-w-full">
	<NavBar />
	<div class="relative mx-auto max-w-full md:px-6 mt-6">
		{#if $authStore.fetching}
			<h2>Authenticating...</h2>
		{:else if $authStore.error}
			<h2>Error: {$authStore.error.message}</h2>
		{:else if $authStore.ok}
			<slot />
		{:else if isLoginPage()}
			<slot />
		{/if}
	</div>
</div>
