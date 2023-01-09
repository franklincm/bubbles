package main

import (
	"fmt"
	"math"
	"os"

	"github.com/franklincm/bubbles/tabs"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

var quitKeys = key.NewBinding(
	key.WithKeys("q", "ctrl+c"),
	key.WithHelp("", "press q to quit"),
)

type model struct {
	err      error
	tabs     tea.Model
	cursor   int
	quitting bool
}

func initialModel() model {

	headers := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
	}

	tabs := tabs.New(headers)

	return model{
		tabs: tabs,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) {
			m.quitting = true
			return m, tea.Quit
		} else if key.Matches(msg, key.NewBinding(key.WithKeys("l"))) {
			numHeadings := len(m.tabs.(tabs.Model).GetHeadings())
			m.cursor = int(math.Min(float64(m.cursor+1), float64(numHeadings-1)))
			m.tabs = m.tabs.(tabs.Model).SetFocused(m.cursor)
		} else if key.Matches(msg, key.NewBinding(key.WithKeys("h"))) {
			m.cursor = int(math.Max(float64(m.cursor-1), 0))
			m.tabs = m.tabs.(tabs.Model).SetFocused(m.cursor)
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return lipgloss.
		NewStyle().
		Render(m.tabs.View() + "\n")
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
