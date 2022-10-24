<script lang="ts">
	import { base } from '$app/paths';
	import { createHttpStore } from '$lib/http/store';
	import { session } from '$lib/session/store';

	const logoutStore = createHttpStore();

	function handleLogout() {
		logoutStore.get('/logout');
		logoutStore.subscribe((value) => {
			if (value.ok) {
				session.logout();
			}
		});
	}
</script>

<div class="navbar bg-primary text-primary-content">
	<div class="navbar-start space-x-16">
		<a href={base} class="btn btn-ghost normal-case text-xl">Coastline</a>
		<a href="{base}/requests" class="btn btn-ghost normal-case">Requests</a>
	</div>

	<div class="navbar-center" />

	<div class="navbar-end">
		{#if $session.initialized && $session.authenticated}
			<button on:click={handleLogout} class="btn">Logout</button>
		{/if}
	</div>
</div>
