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

<div class="navbar bg-primary text-primary-content not-prose">
	<div class="navbar-start space-x-16">
		<a href={base} class="btn btn-ghost normal-case text-xl">Coastline</a>
	</div>
	<div class="navbar-center">
		<a href="{base}/projects" class="btn btn-ghost normal-case">Projects</a>
		<a href="{base}/services" class="btn btn-ghost normal-case">Services</a>
		<a href="{base}/requests" class="btn btn-ghost normal-case">Requests</a>
		<a href="{base}/templates" class="btn btn-ghost normal-case">Templates</a>
	</div>

	<div class="navbar-end">
		{#if $session.initialized && $session.authenticated}
			<div class="dropdown dropdown-end">
				<label
					for=""
					tabindex="0"
					class="btn bg-primary-content btn-outline btn-circle avatar placeholder"
				>
					<div class="w-10 h-10 rounded-full">
						{#if $session.user?.picture}
							<img class="object-fill" src={$session.user.picture} alt="user" />
						{:else}
							<span class="text-xl">JO</span>
						{/if}
					</div>
				</label>
				<ul
					tabindex="0"
					class="menu menu-compact dropdown-content mt-3 p-2 shadow bg-base-100 rounded-box w-52"
				>
					<li>
						<a class="text-primary" href="{base}/profile">Profile</a>
					</li>
					<li><button class="text-primary" on:click={handleLogout}>Logout</button></li>
				</ul>
			</div>
			<!-- <button on:click={handleLogout} class="btn">Logout</button> -->
		{/if}
	</div>
</div>
