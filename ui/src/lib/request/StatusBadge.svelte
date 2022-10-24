<script lang="ts">
	import type { components } from '$lib/oapi/gen/types';

	type Request = components['schemas']['Request'];

	export let request: Request;

	const text = request.status.split('_').join(' ').toUpperCase();

	const color = (): string => {
		switch (request.status) {
			case 'rejected':
				return 'badge-error';
			case 'pending_approval':
				return 'badge-warning';
			case 'approved':
				return 'badge-success';
		}
		return '';
	};
</script>

<span class="badge badge-lg {color()}">{text}</span>
