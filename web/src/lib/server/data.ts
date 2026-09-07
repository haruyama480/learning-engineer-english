import { readFileSync } from 'node:fs';
import path from 'node:path';
import { parse } from 'yaml';
import type { QuizItem, QuizPayload, Source, VocabEntry } from '$lib/types';

const dataDir = path.resolve(process.cwd(), '../vocabulary/mercari');

type VocabularyFile = {
	source: Source;
	entries: VocabEntry[];
};

type QuizFile = {
	source: Source;
	quizzes: QuizItem[];
};

export function loadQuizPayload(): QuizPayload {
	const vocabulary = parse(
		readFileSync(path.join(dataDir, 'vocabulary.yaml'), 'utf8')
	) as VocabularyFile;
	const quizzes = parse(readFileSync(path.join(dataDir, 'quizzes.yaml'), 'utf8')) as QuizFile;

	return {
		source: vocabulary.source,
		entries: vocabulary.entries.map((entry) => ({
			id: entry.id,
			ja: entry.ja,
			ja_reading: entry.ja_reading,
			en: entry.en,
			pos: entry.pos,
			notes: entry.notes,
			lists: entry.lists
		})),
		quizzes: quizzes.quizzes.map((quiz) => ({
			id: quiz.id,
			vocab_id: quiz.vocab_id,
			prompt_ja: quiz.prompt_ja,
			blanked_en: quiz.blanked_en,
			answers: quiz.answers,
			hint: quiz.hint
		}))
	};
}
