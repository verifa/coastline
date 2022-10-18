<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	import { session } from '$lib/session/store';
	import { createHttpStore } from '$lib/http/store';
	import type { UserResponse } from '$lib/session/session';
	import { page } from '$app/stores';

	import TopAppBar, { Row, Section, Title, AutoAdjust } from '@smui/top-app-bar';
	import IconButton from '@smui/icon-button';
	import Button, { Label } from '@smui/button';

	const authStore = createHttpStore<UserResponse>();
	const logoutStore = createHttpStore();

	$: isLoginPage = (): boolean => {
		return $page.url.pathname === '/login';
	};
	authStore.subscribe((resp) => {
		// Do nothing if request is in progress
		if (resp.fetching) {
			return;
		}
		if (resp.ok && resp.data) {
			session.login(resp.data.user);
		} else if (resp.status === 401) {
			session.logout();
		}
	});
	authStore.get('/authenticate');

	onMount(() => {
		if ('theme' in localStorage) {
			const theme = localStorage['theme'];
			if (theme === 'dark') {
				darkTheme = true;
			} else {
				darkTheme = false;
			}
		} else {
			// Get preference for dark theme from browser
			darkTheme = window.matchMedia('(prefers-color-scheme: dark)').matches;
		}

		// Subscribe to the session store to handle events if user is not authenticated
		session.subscribe((session) => {
			if (!session.initialized) {
				return;
			} else if (!session.authenticated) {
				if (!isLoginPage()) {
					goto('/login');
				}
			} else if (session.initialized && session.authenticated) {
				if (isLoginPage()) {
					goto('/');
				}
			}
		});
	});

	let topAppBar: TopAppBar;
	let darkTheme: boolean = false;

	$: modeLabel = `switch to ${darkTheme ? 'light' : 'dark'} mode`;
	$: modeIcon = darkTheme ? 'light_mode' : 'dark_mode';

	function toggleMode() {
		darkTheme = !darkTheme;
		// Remember choice: store the theme setting in local storage
		localStorage.setItem('theme', darkTheme ? 'dark' : 'light');
	}

	function handleLogout() {
		logoutStore.get('/logout');
		logoutStore.subscribe((value) => {
			if (value.ok) {
				session.logout();
			}
		});
	}
</script>

<TopAppBar bind:this={topAppBar} variant="fixed">
	<Row>
		<Section>
			<IconButton class="material-icons">menu</IconButton>
			<Title>Fixed</Title>
		</Section>
		<Section align={'end'}>
			<!-- Light/Dark mode -->
			<IconButton
				aria-label={modeLabel}
				class="material-icons"
				on:click={toggleMode}
				title={modeLabel}
			>
				{modeIcon}
			</IconButton>
			<Button color="secondary" on:click={handleLogout}>
				<Label>Logout</Label>
			</Button>
		</Section>
	</Row>
</TopAppBar>
<AutoAdjust {topAppBar}>
	{#if $authStore.fetching}
		<h2>Authenticating...</h2>
	{:else if $authStore.error}
		<h2>Error: {$authStore.error.message}</h2>
	{:else if $authStore.ok}
		<slot />
	{:else if isLoginPage()}
		<slot />
	{/if}
</AutoAdjust>

<svelte:head>
	{#if darkTheme}
		<link rel="stylesheet" href="/smui-dark.css" />
	{:else}
		<link rel="stylesheet" href="/smui.css" />
	{/if}
</svelte:head>
