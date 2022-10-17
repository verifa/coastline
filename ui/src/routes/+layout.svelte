<script lang="ts">
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	

	import {user} from '$lib/auth/store';
	import { createHttpStore } from '$lib/http/store';
	import type { UserResponse } from '$lib/auth/user';
	import { page } from '$app/stores';

    import TopAppBar, {
      Row,
      Section,
      Title,
      AutoAdjust,
    } from '@smui/top-app-bar';
    import IconButton from '@smui/icon-button';
   
    
    const http = createHttpStore<UserResponse>()
    
    $: isLoginPage = (): boolean => {
        return $page.url.pathname === "/login"
    }
        
    onMount(()=> {
        // Get preference for dark theme from browser
        darkTheme = window.matchMedia('(prefers-color-scheme: dark)').matches;

        http.subscribe((resp) => {
            if (resp.ok && resp.data) {
                user.set(resp.data.user)
            } else if (!resp.error) {
                if (resp.status === 401) {
                    if (!isLoginPage()) {
                        goto("/login")
                    }
                }
            }
        })
        http.get("/authenticate")
    })
    
    let topAppBar: TopAppBar;
    let darkTheme: boolean = false;

    $: modeLabel = `switch to ${darkTheme ? 'light' : 'dark'} mode`;
    $: modeIcon = darkTheme ? 'light_mode' : 'dark_mode';

    function toggleMode() {
        darkTheme = !darkTheme
    }

</script>

<TopAppBar bind:this={topAppBar} variant="fixed">
    <Row>
      <Section>
        <IconButton class="material-icons">menu</IconButton>
        <Title>Fixed</Title>
      </Section>
      <Section align={"end"}>
        <IconButton
        aria-label="{modeLabel}"
        class="material-icons"
        on:click="{toggleMode}"
        title="{modeLabel}"
      >{modeIcon}
    </IconButton>
    </Section>
    </Row>
</TopAppBar>
<AutoAdjust {topAppBar}>

{#if $http.fetching}
    <h2>Authenticating...</h2>
{:else if $http.error}
    <h2>Error: {$http.error.message}</h2>
{:else if $http.ok}
	<slot/>
{:else if isLoginPage()}
	<slot/>
{/if}
</AutoAdjust>

<svelte:head>
  {#if darkTheme}
    <!-- TODO: toggle dark theme at a later date... -->
    <!-- <link rel="stylesheet" href="/smui-dark.css" media="screen" /> -->
    <link rel="stylesheet" href="/smui.css" />
  {:else}
    <link rel="stylesheet" href="/smui.css" />
    {/if}
</svelte:head>