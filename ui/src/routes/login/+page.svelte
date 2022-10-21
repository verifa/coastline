<script lang="ts">
	import { goto } from '$app/navigation';

	function handleLogin() {
		goto('/api/v1/login');
	}
</script>

<h1>Login</h1>

<button class="btn btn-primary" on:click={handleLogin}>Login</button>
