<script lang="ts">
	import { createHttpStore } from '$lib/http/store';

	const loginStore = createHttpStore()

	// function handleLogin() {
	// 	loginStore.get("/login")
	// }
</script>

<h1>Login</h1>

<a href="/api/v1/login">Login</a>
