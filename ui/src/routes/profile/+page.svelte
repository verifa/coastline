<script lang="ts">
	import { session } from '$lib/session/store';
</script>

<h1>Profile</h1>

{#if $session.user}
	<div class="max-w-md">
		<table class="table w-full">
			<tbody>
				<tr>
					<td>Issuer</td>
					<td>{$session.user.iss}</td>
				</tr>
				<tr>
					<td>Subject</td>
					<td>{$session.user.sub}</td>
				</tr>
				<tr>
					<td>Name</td>
					<td>{$session.user.name}</td>
				</tr>
				<tr>
					<td>Email</td>
					<td>{$session.user.email}</td>
				</tr>
				<tr>
					<td>Groups</td>
					<td>{$session.user.groups.join(', ')}</td>
				</tr>
			</tbody>
		</table>
	</div>
{/if}
