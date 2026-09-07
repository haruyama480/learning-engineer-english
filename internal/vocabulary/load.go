package vocabulary

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadVocabulary reads a VocabularyFile from path.
func LoadVocabulary(path string) (*VocabularyFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read vocabulary: %w", err)
	}
	var file VocabularyFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse vocabulary: %w", err)
	}
	return &file, nil
}

// LoadQuizzes reads a QuizFile from path.
func LoadQuizzes(path string) (*QuizFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read quizzes: %w", err)
	}
	var file QuizFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parse quizzes: %w", err)
	}
	return &file, nil
}

// ValidateVocabulary checks IDs, example counts, and that each blank occurs in EN.
func ValidateVocabulary(file *VocabularyFile) error {
	if file == nil {
		return fmt.Errorf("vocabulary file is nil")
	}
	seen := make(map[string]struct{}, len(file.Entries))
	for i, entry := range file.Entries {
		if entry.ID == "" {
			return fmt.Errorf("entries[%d]: empty id", i)
		}
		if _, ok := seen[entry.ID]; ok {
			return fmt.Errorf("duplicate id %q", entry.ID)
		}
		seen[entry.ID] = struct{}{}
		if entry.JA == "" || entry.EN == "" {
			return fmt.Errorf("entry %q: ja and en are required", entry.ID)
		}
		n := len(entry.Examples)
		if n < MinExamples || n > MaxExamples {
			return fmt.Errorf("entry %q: want %d-%d examples, got %d", entry.ID, MinExamples, MaxExamples, n)
		}
		for j, ex := range entry.Examples {
			if ex.JA == "" || ex.EN == "" {
				return fmt.Errorf("entry %q examples[%d]: ja and en are required", entry.ID, j)
			}
			if ex.Blank == "" {
				return fmt.Errorf("entry %q examples[%d]: blank is required", entry.ID, j)
			}
			if !strings.Contains(ex.EN, ex.Blank) {
				return fmt.Errorf("entry %q examples[%d]: blank %q not found in en %q", entry.ID, j, ex.Blank, ex.EN)
			}
		}
	}
	return nil
}

// QuizzesFromVocabulary converts each example into a cloze quiz.
func QuizzesFromVocabulary(file *VocabularyFile) (*QuizFile, error) {
	if err := ValidateVocabulary(file); err != nil {
		return nil, err
	}
	out := &QuizFile{
		Source:  file.Source,
		Quizzes: make([]Quiz, 0, len(file.Entries)*MinExamples),
	}
	for _, entry := range file.Entries {
		for i, ex := range entry.Examples {
			blanked := strings.Replace(ex.EN, ex.Blank, BlankPlaceholder, 1)
			out.Quizzes = append(out.Quizzes, Quiz{
				ID:        fmt.Sprintf("%s-%d", entry.ID, i+1),
				VocabID:   entry.ID,
				PromptJA:  ex.JA,
				BlankedEN: blanked,
				Answers:   []string{ex.Blank},
				Hint:      entry.EN,
			})
		}
	}
	return out, nil
}
