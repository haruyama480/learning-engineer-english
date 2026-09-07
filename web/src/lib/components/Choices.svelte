<script lang="ts">
	import { answersMatch } from '$lib/quiz';

	let {
		choices,
		answers,
		selected = '',
		status = 'idle',
		onselect
	}: {
		choices: string[];
		answers: string[];
		selected?: string;
		status?: 'idle' | 'ok' | 'ng';
		onselect: (choice: string) => void;
	} = $props();

	const labels = ['A', 'B', 'C', 'D'];
</script>

<div class="choices" role="listbox" aria-label="英語の選択肢">
	{#each choices as choice, i (choice)}
		{@const isAnswer = answersMatch(choice, answers)}
		{@const isPicked = selected === choice}
		<button
			type="button"
			class="choice"
			role="option"
			aria-selected={isPicked}
			class:picked={isPicked && status === 'idle'}
			class:ok={status !== 'idle' && isAnswer}
			class:ng={status === 'ng' && isPicked && !isAnswer}
			disabled={status !== 'idle'}
			onclick={() => onselect(choice)}
		>
			<span class="choice-key">{labels[i] ?? i + 1}</span>
			<span class="choice-text">{choice}</span>
		</button>
	{/each}
</div>
