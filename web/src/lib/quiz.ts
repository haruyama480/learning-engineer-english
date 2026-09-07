import type { PlayQuestion, QuizItem, QuizView, VocabEntry } from './types';

export const BLANK = '____';
export const CHOICE_COUNT = 4;

export type SessionOptions = {
	lists: number[];
	shuffle: boolean;
	onePerVocab: boolean;
};

export function normalizeAnswer(input: string): string {
	return input
		.normalize('NFKC')
		.trim()
		.replace(/\s+/g, ' ')
		.replace(/[’‘]/g, "'")
		.replace(/[“”]/g, '"')
		.toLowerCase();
}

export function answersMatch(input: string, answers: string[]): boolean {
	const got = normalizeAnswer(input);
	if (!got) return false;
	return answers.some((answer) => normalizeAnswer(answer) === got);
}

export function splitCloze(blanked: string): { before: string; after: string } {
	const index = blanked.indexOf(BLANK);
	if (index < 0) {
		return { before: blanked, after: '' };
	}
	return {
		before: blanked.slice(0, index),
		after: blanked.slice(index + BLANK.length)
	};
}

export function longestAnswerLength(answers: string[]): number {
	return answers.reduce((max, answer) => Math.max(max, answer.length), 8);
}

function distractorScore(quiz: QuizItem, candidate: QuizItem, vocabById: Map<string, VocabEntry>): number {
	const self = vocabById.get(quiz.vocab_id);
	const other = vocabById.get(candidate.vocab_id);
	let score = 0;
	if (candidate.vocab_id === quiz.vocab_id) score += 5;
	if (self?.pos && other?.pos === self.pos) score += 2;
	if (self && other?.lists.some((list) => self.lists.includes(list))) score += 1;
	const lengthGap = Math.abs((candidate.answers[0]?.length ?? 0) - (quiz.answers[0]?.length ?? 0));
	score += Math.max(0, 3 - lengthGap / 5);
	return score;
}

export function buildChoices(
	quiz: QuizItem,
	pool: QuizItem[],
	entries: VocabEntry[],
	random: () => number = Math.random
): string[] {
	const correct = quiz.answers[0];
	if (!correct) return [];

	const vocabById = new Map(entries.map((entry) => [entry.id, entry]));
	const used = new Set<string>(quiz.answers.map(normalizeAnswer));
	const ranked = pool
		.filter((item) => item.id !== quiz.id)
		.map((item) => ({ item, score: distractorScore(quiz, item, vocabById) }))
		.sort((a, b) => b.score - a.score);

	const distractors: string[] = [];
	for (const { item } of ranked) {
		const answer = item.answers[0];
		if (!answer) continue;
		const key = normalizeAnswer(answer);
		if (used.has(key)) continue;
		used.add(key);
		distractors.push(answer);
		if (distractors.length >= CHOICE_COUNT - 1) break;
	}

	return shuffle([correct, ...distractors], random);
}

export function attachChoices(
	items: QuizView[],
	pool: QuizItem[],
	entries: VocabEntry[]
): PlayQuestion[] {
	return items.map((item) => ({
		...item,
		choices: buildChoices(item, pool, entries)
	}));
}

export function shuffle<T>(items: readonly T[], random: () => number = Math.random): T[] {
	const out = items.slice();
	for (let i = out.length - 1; i > 0; i--) {
		const j = Math.floor(random() * (i + 1));
		[out[i], out[j]] = [out[j], out[i]];
	}
	return out;
}

export function buildSession(
	quizzes: QuizItem[],
	entries: VocabEntry[],
	options: SessionOptions
): QuizView[] {
	const byId = new Map(entries.map((entry) => [entry.id, entry]));
	const selected = new Set(options.lists);
	let picked: QuizView[] = [];

	for (const quiz of quizzes) {
		const vocab = byId.get(quiz.vocab_id);
		if (!vocab) continue;
		if (selected.size > 0 && !vocab.lists.some((list) => selected.has(list))) continue;
		picked.push({ ...quiz, vocab });
	}

	if (options.onePerVocab) {
		const groups = new Map<string, QuizView[]>();
		for (const quiz of picked) {
			const group = groups.get(quiz.vocab_id) ?? [];
			group.push(quiz);
			groups.set(quiz.vocab_id, group);
		}
		picked = [...groups.values()].map((group) => group[Math.floor(Math.random() * group.length)]);
	}

	if (options.shuffle) {
		picked = shuffle(picked);
	}

	return picked;
}

export function availableLists(entries: VocabEntry[]): number[] {
	return [...new Set(entries.flatMap((entry) => entry.lists))].sort((a, b) => a - b);
}
