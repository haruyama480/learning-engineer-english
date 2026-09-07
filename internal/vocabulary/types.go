// Package vocabulary defines the YAML schemas for engineer vocabulary
// entries and Japanese-to-English cloze quizzes.
package vocabulary

const (
	// BlankPlaceholder is the token that replaces the target span in a cloze quiz.
	BlankPlaceholder = "____"

	MinExamples = 2
	MaxExamples = 3
)

// VocabularyFile is the root document of vocabulary.yaml.
type VocabularyFile struct {
	Source  Source  `yaml:"source"`
	Entries []Entry `yaml:"entries"`
}

// Source records where the word list came from.
type Source struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	License  string `yaml:"license"`
	Adapted  bool   `yaml:"adapted"`
	Notes    string `yaml:"notes,omitempty"`
}

// Entry is one vocabulary item with example sentences.
type Entry struct {
	ID        string    `yaml:"id"`
	JA        string    `yaml:"ja"`
	JAReading string    `yaml:"ja_reading,omitempty"`
	EN        string    `yaml:"en"`
	POS       string    `yaml:"pos,omitempty"`
	Notes     string    `yaml:"notes,omitempty"`
	Lists     []int     `yaml:"lists"`
	Examples  []Example `yaml:"examples"`
}

// Example is a Japanese/English sentence pair.
// Blank is the exact English span to hide in the cloze quiz; it must occur in EN.
type Example struct {
	JA    string `yaml:"ja"`
	EN    string `yaml:"en"`
	Blank string `yaml:"blank"`
}

// QuizFile is the root document of quizzes.yaml.
type QuizFile struct {
	Source  Source `yaml:"source"`
	Quizzes []Quiz `yaml:"quizzes"`
}

// Quiz is a Japanese-prompt, English fill-in-the-blank item.
type Quiz struct {
	ID        string   `yaml:"id"`
	VocabID   string   `yaml:"vocab_id"`
	PromptJA  string   `yaml:"prompt_ja"`
	BlankedEN string   `yaml:"blanked_en"`
	Answers   []string `yaml:"answers"`
	Hint      string   `yaml:"hint,omitempty"`
}
