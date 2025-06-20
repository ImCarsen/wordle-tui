package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"wordle-tui/game"

	tea "github.com/charmbracelet/bubbletea"
)

func run(state *game.GameState) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	p := tea.NewProgram(
		initialModel(state),
		tea.WithContext(ctx),
		tea.WithAltScreen(), // fullscreen terminal
	)

	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
