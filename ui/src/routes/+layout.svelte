<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount, setContext } from 'svelte';
	

	import {user} from '$lib/auth/store';
	import { createHttpStore } from '$lib/http/store';
	import type { UserResponse } from '$lib/auth/user';
	import { page } from '$app/stores';

    const http = createHttpStore<UserResponse>()

    const isLoginPage = (): boolean => {
        return $page.url.pathname === "/login"
    }
    
    onMount(()=> {
        http.subscribe((resp) => {
            if (resp.ok && resp.data) {
                user.set(resp.data.user)
            } else if (!resp.error) {
                if (resp.status === 401) {
                    if (isLoginPage()) {
                        goto("/login")
                    }
                }
            }
        })
        http.get("/authenticate")
    })

    console.log()


</script>

{#if $http.fetching}
    <h2>Authenticating...</h2>
{:else if $http.error}
    <h2>Error: {$http.error.message}</h2>
{:else if $http.ok}
	<slot/>
{:else if isLoginPage()}
	<slot/>
{/if}