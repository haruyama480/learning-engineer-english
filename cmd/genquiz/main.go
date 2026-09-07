// Command genquiz builds quizzes.yaml from vocabulary.yaml.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/haruyama480/learning-engineer-english/internal/vocabulary"
	"gopkg.in/yaml.v3"
)

func main() {
	in := flag.String("in", "vocabulary/mercari/vocabulary.yaml", "path to vocabulary.yaml")
	out := flag.String("out", "vocabulary/mercari/quizzes.yaml", "path to write quizzes.yaml")
	flag.Parse()

	if err := run(*in, *out); err != nil {
		fmt.Fprintf(os.Stderr, "genquiz: %v\n", err)
		os.Exit(1)
	}
}

func run(in, out string) error {
	vocab, err := vocabulary.LoadVocabulary(in)
	if err != nil {
		return err
	}
	quizzes, err := vocabulary.QuizzesFromVocabulary(vocab)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(quizzes); err != nil {
		return fmt.Errorf("encode quizzes: %w", err)
	}
	if err := enc.Close(); err != nil {
		return fmt.Errorf("close encoder: %w", err)
	}
	if err := os.WriteFile(out, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write quizzes: %w", err)
	}
	fmt.Printf("wrote %d quizzes to %s\n", len(quizzes.Quizzes), out)
	return nil
}
