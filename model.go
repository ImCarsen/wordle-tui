package main

import (
	"fmt"
	"strings"
	"unicode"
	"wordle-tui/game"
	"wordle-tui/words"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	state   *game.GameState
	input   string
	message string

	width, height int
}

var (
	blockStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true).
			MarginRight(1).
			Align(lipgloss.Center)

	greenBlock   = blockStyle.Background(lipgloss.Color("#538d4e")).Foreground(lipgloss.Color("#ffffff"))
	yellowBlock  = blockStyle.Background(lipgloss.Color("#b59f3b")).Foreground(lipgloss.Color("#ffffff"))
	grayBlock    = blockStyle.Background(lipgloss.Color("#3a3a3c")).Foreground(lipgloss.Color("#ffffff"))
	neutralBlock = blockStyle.Background(lipgloss.Color("#565758")).Foreground(lipgloss.Color("#ffffff"))

	messageStyle = lipgloss.NewStyle().
			Bold(true).
			MarginTop(1).
			Align(lipgloss.Center).
			Foreground(lipgloss.Color("#ffcc00"))

	lossStyle = messageStyle.Foreground(lipgloss.Color("#ff5555"))

	keybindStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Align(lipgloss.Center).
			MarginTop(1)
)

func initialModel(state *game.GameState) model {
	if state == nil {
		return model{
			state: game.New(getWordOfTheDay(), true),
		}
	}
	return model{
		state: state,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.state.SaveState()
			return m, tea.Quit
		case tea.KeyCtrlC:
			if m.state.Daily && m.state.Over {
				m.message = "Copied to clipboard!"
				clipboard.WriteAll(m.shareGrid())
				return m, nil
			}
			return m, nil
		case tea.KeyEnter:
			if len(m.input) != len(m.state.Answer) {
				m.message = "Not enough letters"
				return m, nil
			}

			if m.state.Daily && !words.Guesses[m.input] {
				m.message = "Not in word list"
				return m, nil
			}

			if len(m.input) == 5 && !m.state.Over {
				m.state.Guesses = append(m.state.Guesses, m.input)
				fb := checkGuess(m.state.Answer, m.input)
				m.state.Feedback = append(m.state.Feedback, fb)

				for i, l := range m.input {
					r := unicode.ToLower(l)
					res := fb[i]
					prev := m.state.LetterStatus[r]
					if prev == "green" {
						continue
					}
					if prev == "yellow" && res == "gray" {
						continue
					}
					m.state.LetterStatus[r] = res
				}

				if m.input == m.state.Answer {
					m.state.Won = true
					m.state.Over = true
					m.message = "You win!"
				} else if len(m.state.Guesses) >= 6 {
					m.state.Over = true
					m.message = "Out of guesses!"
				}
				m.input = ""
			}

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.message = ""
			}

		default:
			if len(msg.String()) == 1 && !m.state.Over && len(m.input) < 5 {
				ch := rune(msg.String()[0])
				if unicode.IsLetter(ch) {
					m.input += strings.ToLower(string(ch))
					m.message = ""
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Render guesses
	for i, guess := range m.state.Guesses {
		for j, letter := range guess {
			style := grayBlock
			switch m.state.Feedback[i][j] {
			case "green":
				style = greenBlock
			case "yellow":
				style = yellowBlock
			}
			b.WriteString(style.Render(strings.ToUpper(string(letter))))
		}
		b.WriteString("\n")
	}

	for i := len(m.state.Guesses); i < 6; i++ {
		for j := 0; j < 5; j++ {
			b.WriteString(grayBlock.Render(" "))
		}
		b.WriteString("\n")
	}

	// Input
	if !m.state.Over {
		b.WriteString("\n")
		for i := 0; i < 5; i++ {
			if i < len(m.input) {
				b.WriteString(neutralBlock.Render(strings.ToUpper(string(m.input[i]))))
			} else {
				b.WriteString(neutralBlock.Render(" "))
			}
		}
	}

	// Keys
	b.WriteString("\n\n")
	keyboardRows := []string{
		"QWERTYUIOP",
		"ASDFGHJKL",
		"ZXCVBNM",
	}

	for _, row := range keyboardRows {
		var indent string
		switch len(row) {
		case 9:
			indent = " "
		case 7:
			indent = "  "
		}

		b.WriteString(indent)
		for _, key := range row {
			l := unicode.ToLower(key)
			var style lipgloss.Style
			switch m.state.LetterStatus[l] {
			case "green":
				style = greenBlock
			case "yellow":
				style = yellowBlock
			case "gray":
				style = grayBlock
			default:
				style = blockStyle
			}
			b.WriteString(style.Render(string(key)))
		}
		b.WriteString("\n")
	}

	// Result message
	if m.state.Over {
		if m.state.Won {
			b.WriteString(messageStyle.Render("\n🎉 You win!"))
		} else {
			b.WriteString(lossStyle.Render("\n❌ You lose! Word was: " + strings.ToUpper(m.state.Answer)))
		}

		// Add share grid
		b.WriteString("\n\n")
		b.WriteString(m.shareGrid())
	}

	if m.message != "" && !m.state.Over {
		b.WriteString("\n")
		b.WriteString(messageStyle.Render(m.message))
	}

	b.WriteString("\n\n")

	if m.state.Daily && m.state.Over {
		b.WriteString(keybindStyle.Render("ENTER: submit • BACKSPACE: delete • CTRL+C: copy to clipboard • ESC: quit"))
	} else {
		b.WriteString(keybindStyle.Render("ENTER: submit • BACKSPACE: delete • ESC: quit"))
	}

	content := b.String()

	return lipgloss.PlaceVertical(m.height, lipgloss.Center,
		lipgloss.PlaceHorizontal(m.width, lipgloss.Center, content),
	)
}

func (m model) shareGrid() string {
	var sb strings.Builder

	num := fmt.Sprint(len(m.state.Guesses))
	sb.WriteString("Wordle " + fmt.Sprint(getWordleNumber()) + " " + num + "/6\n\n")
	// Build a string for each guess row
	for _, fb := range m.state.Feedback {
		for _, res := range fb {
			switch res {
			case "green":
				sb.WriteString("🟩")
			case "yellow":
				sb.WriteString("🟨")
			default:
				sb.WriteString("⬛")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
