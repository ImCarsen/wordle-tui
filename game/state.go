package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type GameState struct {
	Answer       string
	Guesses      []string
	Feedback     [][]string
	LetterStatus map[rune]string
	Won, Over    bool
	Day          string
	Daily        bool
}

func New(answer string, daily bool) *GameState {
	return &GameState{
		Answer:       answer,
		Guesses:      []string{},
		Feedback:     [][]string{},
		LetterStatus: map[rune]string{},
		Won:          false,
		Over:         false,
		Day:          time.Now().Format("2006-01-02"),
		Daily:        daily,
	}
}

func getSavePath() string {
	dir, _ := os.UserConfigDir()
	return filepath.Join(dir, "wordle-tui", "state.json")
}

func LoadSavedState() (*GameState, error) {
	path := getSavePath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s GameState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *GameState) SaveState() error {
	// Don't save custom
	if !s.Daily {
		return nil
	}
	path := getSavePath()
	data, _ := json.MarshalIndent(s, "", "  ")
	os.MkdirAll(filepath.Dir(path), 0700)
	return os.WriteFile(path, data, 0600)
}

func ClearSavedState() error {
	return os.Remove(getSavePath())
}

func (s *GameState) IsToday() bool {
	today := time.Now().Format("2006-01-02")
	return s.Day == today
}
