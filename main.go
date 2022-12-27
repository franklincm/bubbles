package main

import (
	"fmt"
	"os"
	"strings"

	commandprompt "github.com/franklincm/bubbles/commandPrompt"

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
	err        error
	input      []string
	prompt     tea.Model
	quitting   bool
	showprompt bool
}

func initialModel() model {

	prompt := commandprompt.New()
	prompt.InputShow = key.NewBinding(key.WithKeys(":"))

	return model{
		prompt: prompt,
	}
}

func (m model) Init() tea.Cmd {
	m.input = make([]string, 3)
	return tea.Batch(
		m.prompt.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case commandprompt.PromptInput:
		if msg == "quit" || msg == "q" {
			m.quitting = true
			return m, tea.Quit
		}

		m.input = append(m.input, string(msg))
		return m, nil

	case commandprompt.PromptEditing:
		m.showprompt = bool(msg)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, quitKeys) && !m.showprompt {
			m.quitting = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.prompt, cmd = m.prompt.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.showprompt {
		return lipgloss.
			NewStyle().
			Width(25).
			BorderStyle(lipgloss.NormalBorder()).
			Render(m.prompt.View()) + strings.Join(m.input, "\n") + "\n"
	}
	return lipgloss.
		NewStyle().
		Width(25).
		Render("") + strings.Join(m.input, "\n") + "\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
