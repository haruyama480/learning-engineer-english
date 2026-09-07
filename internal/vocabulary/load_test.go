package vocabulary

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}

func TestLoadAndValidateVocabulary(t *testing.T) {
	path := filepath.Join(repoRoot(t), "vocabulary/mercari/vocabulary.yaml")
	file, err := LoadVocabulary(path)
	if err != nil {
		t.Fatalf("LoadVocabulary: %v", err)
	}
	if err := ValidateVocabulary(file); err != nil {
		t.Fatalf("ValidateVocabulary: %v", err)
	}
	if len(file.Entries) < 100 {
		t.Fatalf("expected at least 100 entries, got %d", len(file.Entries))
	}
}

func TestQuizzesMatchCommittedFile(t *testing.T) {
	root := repoRoot(t)
	vocab, err := LoadVocabulary(filepath.Join(root, "vocabulary/mercari/vocabulary.yaml"))
	if err != nil {
		t.Fatalf("LoadVocabulary: %v", err)
	}
	got, err := QuizzesFromVocabulary(vocab)
	if err != nil {
		t.Fatalf("QuizzesFromVocabulary: %v", err)
	}
	want, err := LoadQuizzes(filepath.Join(root, "vocabulary/mercari/quizzes.yaml"))
	if err != nil {
		t.Fatalf("LoadQuizzes: %v", err)
	}
	if len(got.Quizzes) != len(want.Quizzes) {
		t.Fatalf("quiz count: generated %d, committed %d; run go run ./cmd/genquiz", len(got.Quizzes), len(want.Quizzes))
	}
	if !reflect.DeepEqual(got.Quizzes, want.Quizzes) {
		for i := range got.Quizzes {
			if !reflect.DeepEqual(got.Quizzes[i], want.Quizzes[i]) {
				t.Fatalf("quiz[%d] mismatch\ngot:  %+v\nwant: %+v", i, got.Quizzes[i], want.Quizzes[i])
			}
		}
		t.Fatal("quizzes differ")
	}
	for i := range got.Quizzes {
		if !strings.Contains(got.Quizzes[i].BlankedEN, BlankPlaceholder) {
			t.Fatalf("quiz %s: blanked_en has no placeholder: %q", got.Quizzes[i].ID, got.Quizzes[i].BlankedEN)
		}
	}
}
