package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"wordle-tui/words"
)

type WordleResponse struct {
	Solution string `json:"solution"`
}

var (
	epoch = time.Date(2021, 6, 19, 0, 0, 0, 0, time.UTC)
)

func getFallbackWord() string {
	today := time.Now().UTC()
	days := int(today.Sub(epoch).Hours() / 24)
	index := days % len(words.Words)
	return words.Words[index]
}

func getWordOfTheDay() string {
	now := time.Now()
	url := fmt.Sprintf(
		"https://www.nytimes.com/svc/wordle/v2/%d-%02d-%02d.json",
		now.Year(), now.Month(), now.Day(),
	)

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var wordleResp WordleResponse
	err = json.Unmarshal(body, &wordleResp)
	if err != nil {
		return ""
	}

	return wordleResp.Solution
}

func checkGuess(answer, guess string) []string {
	result := make([]string, len(guess))
	answerRunes := []rune(answer)
	guessRunes := []rune(guess)

	// Mark exact matches
	for i := range guessRunes {
		if guessRunes[i] == answerRunes[i] {
			result[i] = "green"
			answerRunes[i] = '-' // mark as used
		}
	}

	for i := range guessRunes {
		if result[i] == "" {
			for j := range answerRunes {
				if guessRunes[i] == answerRunes[j] {
					result[i] = "yellow"
					answerRunes[j] = '-' // mark as used
					break
				}
			}
		}
	}

	// Mark absent
	for i := range result {
		if result[i] == "" {
			result[i] = "gray"
		}
	}

	return result
}

func getWordleNumber() int {
	today := time.Now().UTC()
	days := int(today.Sub(epoch).Hours() / 24)
	return days
}
