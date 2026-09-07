<script lang="ts">
	import Choices from '$lib/components/Choices.svelte';
	import Cloze from '$lib/components/Cloze.svelte';
	import { answersMatch, attachChoices, availableLists, buildSession } from '$lib/quiz';
	import type { PlayQuestion, QuizView } from '$lib/types';

	let { data } = $props();

	type Phase = 'setup' | 'play' | 'results';
	type Status = 'idle' | 'ok' | 'ng';

	const lists = $derived(availableLists(data.entries));
	let selectedLists = $state<number[]>([1, 2, 3, 4, 5]);
	let shuffleOn = $state(true);
	let onePerVocab = $state(false);
	let hideVocab = $state(false);

	let phase = $state<Phase>('setup');
	let session = $state<PlayQuestion[]>([]);
	let index = $state(0);
	let selected = $state('');
	let status = $state<Status>('idle');
	let hintOn = $state(false);
	let correctCount = $state(0);
	let misses = $state<PlayQuestion[]>([]);

	const previewCount = $derived(
		buildSession(data.quizzes, data.entries, {
			lists: selectedLists,
			shuffle: false,
			onePerVocab
		}).length
	);
	const current = $derived(session[index]);
	const total = $derived(session.length);
	const progress = $derived(total === 0 ? 0 : ((index + (status === 'idle' ? 0 : 1)) / total) * 100);
	const filled = $derived(status === 'idle' ? '' : (current?.answers[0] ?? ''));

	function toggleList(list: number) {
		if (selectedLists.includes(list)) {
			selectedLists = selectedLists.filter((item) => item !== list);
			return;
		}
		selectedLists = [...selectedLists, list].sort((a, b) => a - b);
	}

	function toPlay(items: QuizView[]): PlayQuestion[] {
		return attachChoices(items, data.quizzes, data.entries);
	}

	function start(from: QuizView[] | null = null) {
		session = toPlay(
			from ??
				buildSession(data.quizzes, data.entries, {
					lists: selectedLists,
					shuffle: shuffleOn,
					onePerVocab
				})
		);
		if (session.length === 0) return;
		index = 0;
		correctCount = 0;
		misses = [];
		resetQuestion();
		phase = 'play';
	}

	function resetQuestion() {
		selected = '';
		status = 'idle';
		hintOn = false;
	}

	function pick(choice: string) {
		if (!current || status !== 'idle') return;
		selected = choice;
		if (answersMatch(choice, current.answers)) {
			status = 'ok';
			correctCount += 1;
			return;
		}
		status = 'ng';
		misses = [...misses, current];
	}

	function reveal() {
		if (!current || status !== 'idle') return;
		status = 'ng';
		misses = [...misses, current];
		selected = current.answers[0];
	}

	function next() {
		if (index + 1 >= session.length) {
			phase = 'results';
			return;
		}
		index += 1;
		resetQuestion();
	}

	function onKey(event: KeyboardEvent) {
		if (phase !== 'play' || !current) return;
		if (event.key === 'Enter' && status !== 'idle') {
			event.preventDefault();
			next();
			return;
		}
		if (status !== 'idle') return;
		const map: Record<string, number> = { '1': 0, '2': 1, '3': 2, '4': 3, a: 0, b: 1, c: 2, d: 3 };
		const choiceIndex = map[event.key.toLowerCase()];
		if (choiceIndex === undefined) return;
		const choice = current.choices[choiceIndex];
		if (!choice) return;
		event.preventDefault();
		pick(choice);
	}
</script>

<svelte:window onkeydown={onKey} />

<header class="top">
	<div>
		<p class="kicker">Engineer English</p>
		<h1>選択クイズ</h1>
	</div>
	<p class="lede">日本語の例文を見て、英語の空欄に入る語句を選んでください。</p>
</header>

{#if phase === 'setup'}
	<section class="card">
		<div class="field">
			<h2>リスト</h2>
			<div class="row">
				{#each lists as list (list)}
					<button
						type="button"
						class="chip"
						aria-pressed={selectedLists.includes(list)}
						onclick={() => toggleList(list)}
					>
						List {list}
					</button>
				{/each}
			</div>
		</div>

		<div class="field">
			<h2>出題</h2>
			<div class="row">
				<label class="opt">
					<input type="checkbox" bind:checked={shuffleOn} />
					順番をシャッフル
				</label>
				<label class="opt">
					<input type="checkbox" bind:checked={onePerVocab} />
					語彙ごとに1問
				</label>
				<label class="opt">
					<input type="checkbox" bind:checked={hideVocab} />
					語彙を隠す
				</label>
			</div>
		</div>

		<div class="actions">
			<button class="btn" type="button" disabled={previewCount === 0} onclick={() => start()}>
				開始する
			</button>
			<span class="count">{previewCount}問</span>
		</div>
	</section>
{:else if phase === 'play' && current}
	<section class="card">
		<div class="progress">
			<div class="progress-top">
				<span>{index + 1} / {total}</span>
				<span>正解 {correctCount} · ミス {misses.length}</span>
			</div>
			<div class="bar" aria-hidden="true"><span style:width="{progress}%"></span></div>
		</div>

		<div class="meta">
			{#each current.vocab.lists as list (list)}
				<span class="tag">List {list}</span>
			{/each}
			{#if current.vocab.pos}
				<span class="tag">{current.vocab.pos}</span>
			{/if}
		</div>

		{#if !hideVocab || hintOn}
			<p class="vocab">
				{current.vocab.ja}
				{#if current.vocab.ja_reading && current.vocab.ja_reading !== current.vocab.ja}
					<span class="reading">（{current.vocab.ja_reading}）</span>
				{/if}
			</p>
		{:else}
			<p class="vocab"><span class="reading">語彙はヒントで表示</span></p>
		{/if}

		<p class="prompt">{current.prompt_ja}</p>

		{#key current.id}
			<Cloze
				blanked={current.blanked_en}
				{filled}
				widthFrom={current.choices}
				{status}
			/>
			<Choices
				choices={current.choices}
				answers={current.answers}
				{selected}
				{status}
				onselect={pick}
			/>
		{/key}

		{#if hintOn}
			<p class="hint">ヒント: {current.hint}</p>
		{/if}

		{#if status === 'ok'}
			<p class="feedback ok" role="status">正解です。</p>
		{:else if status === 'ng'}
			<p class="feedback ng" role="status">
				正解は <strong>{current.answers.join(' / ')}</strong>
			</p>
		{/if}

		<div class="actions">
			{#if status === 'idle'}
				<button class="btn ghost" type="button" onclick={() => (hintOn = true)}>ヒント</button>
				<button class="btn ghost" type="button" onclick={reveal}>答えを見る</button>
			{:else}
				<button class="btn" type="button" onclick={next}>
					{index + 1 >= total ? '結果を見る' : '次へ'}
				</button>
			{/if}
		</div>
	</section>
{:else}
	<section class="card">
		<p class="kicker">結果</p>
		<p class="score">{correctCount} / {total}</p>
		<p class="lede">
			{total === 0 ? '問題がありません。' : `正答率 ${Math.round((correctCount / total) * 100)}%`}
		</p>

		{#if misses.length > 0}
			<h2 class="section-title">間違えた問題</h2>
			<ul class="misses">
				{#each misses as miss (miss.id)}
					<li>
						<p class="ja">{miss.vocab.ja} — {miss.prompt_ja}</p>
						<p class="en">{miss.answers[0]}</p>
					</li>
				{/each}
			</ul>
		{/if}

		<div class="actions">
			<button class="btn" type="button" onclick={() => start()}>もう一度</button>
			{#if misses.length > 0}
				<button class="btn ghost" type="button" onclick={() => start(misses)}>ミスだけ復習</button>
			{/if}
			<button class="btn ghost" type="button" onclick={() => (phase = 'setup')}>設定に戻る</button>
		</div>
	</section>
{/if}

<p class="foot">
	語彙は
	<a href={data.source.url}>{data.source.name}</a>
	（{data.source.license}）を改変して利用しています。
</p>
