<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>エンジニア英語 選択クイズ</title>
	<meta
		name="description"
		content="日本語の例文を見て、英語の空欄に入る語句を選ぶエンジニア向け語彙クイズ"
	/>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="wrap">
	{@render children()}
</div>
