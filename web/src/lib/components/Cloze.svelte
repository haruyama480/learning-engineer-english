<script lang="ts">
	import { longestAnswerLength, splitCloze } from '$lib/quiz';

	let {
		blanked,
		filled = '',
		widthFrom = [],
		status = 'idle'
	}: {
		blanked: string;
		filled?: string;
		widthFrom?: string[];
		status?: 'idle' | 'ok' | 'ng';
	} = $props();

	const parts = $derived(splitCloze(blanked));
	const width = $derived(longestAnswerLength(widthFrom.length > 0 ? widthFrom : [filled || '________']));
</script>

<p class="cloze">
	{parts.before}<span
		class="blank"
		class:ok={status === 'ok'}
		class:ng={status === 'ng'}
		class:filled={filled.length > 0}
		style:min-width="{width}ch">{filled || '\u00a0'}</span
	>{parts.after}
</p>
