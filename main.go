package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"wordle-tui/game"
	"wordle-tui/words"
)

func main() {
	reset := flag.Bool("reset", false, "Reset saved game state")
	word := flag.String("word", "", "Custom word to use if mode is custom")

	flag.Parse()

	// Words for fallback if the nytimes is not reachable
	if err := words.Load(); err != nil {
		fmt.Println("Failed to load words:", err)
		os.Exit(1)
	}

	if *reset {
		err := game.ClearSavedState()
		if err != nil {
			fmt.Println("Failed to reset saved state:", err)
			os.Exit(1)
		}
	}

	if *word == "" {
		state, err := game.LoadSavedState()
		if err != nil || state == nil || !state.IsToday() {
			wotd := getWordOfTheDay()
			if wotd == "" {
				wotd = getFallbackWord()
			}
			run(game.New(wotd, true))
		} else {
			run(state)
		}
	} else {
		var word = *word
		letters := strings.Split(word, "")
		if len(letters) != 5 {
			fmt.Println("Word must be 5 letters")
			os.Exit(1)
		}

		run(game.New(strings.ToLower(word), false))
	}
}
