<script lang="ts">
	import { getContext } from 'svelte';
	import type { getRequestFunc } from './request';
	const getRequest: getRequestFunc = getContext('request');
	const request = getRequest();
</script>

<div class="mockup-code">
	<pre><code>// TODO: show spec here.</code></pre>
	<pre><code
			>// Do we want to render based on the OpenAPI spec for which the request was created?</code
		></pre>
	<pre><code
			>// Or render based on the current (if it exists) OpenAPI spec for the request template?</code
		></pre>
	<pre><code>// The request template, if it exists, might have changed...</code></pre>
	<pre><code>{JSON.stringify(request.spec, null, 2)}</code></pre>
</div>
