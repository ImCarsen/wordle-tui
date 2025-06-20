package words

import (
	_ "embed"
	"strings"
)

//go:embed answers.txt
var wordsFile string

//go:embed guesses.txt
var guessesFile string

var (
	Words   []string
	Guesses map[string]bool
)

func Load() error {
	for _, line := range strings.Split(wordsFile, "\n") {
		word := strings.TrimSpace(line)
		if word != "" {
			Words = append(Words, word)
		}
	}

	Guesses = make(map[string]bool)
	for _, line := range strings.Split(guessesFile, "\n") {
		word := strings.TrimSpace(line)
		if word != "" {
			Guesses[word] = true
		}
	}

	return nil
}
