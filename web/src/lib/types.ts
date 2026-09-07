export type Source = {
	name: string;
	url: string;
	license: string;
	adapted: boolean;
	notes?: string;
};

export type VocabEntry = {
	id: string;
	ja: string;
	ja_reading?: string;
	en: string;
	pos?: string;
	notes?: string;
	lists: number[];
};

export type QuizItem = {
	id: string;
	vocab_id: string;
	prompt_ja: string;
	blanked_en: string;
	answers: string[];
	hint: string;
};

export type QuizView = QuizItem & {
	vocab: VocabEntry;
};

export type PlayQuestion = QuizView & {
	choices: string[];
};

export type QuizPayload = {
	source: Source;
	entries: VocabEntry[];
	quizzes: QuizItem[];
};
