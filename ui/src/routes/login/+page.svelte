<script lang="ts">
	import Button, { Label } from '@smui/button';
</script>

<h1>Login</h1>

<!--
	Use link to login rather than API request so that redirects are handled properly
 -->
<Button href="/api/v1/login">
	<Label>Login</Label>
</Button>
