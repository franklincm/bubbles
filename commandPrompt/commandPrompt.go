package commandprompt

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textinput textinput.Model
	style     lipgloss.Style
	editing   bool
	input     string

	InputAbort  key.Binding
	InputAccept key.Binding
	InputShow   key.Binding

	Prompt string
}

type PromptInput string
type PromptEditing bool

func (m Model) PromptInput() tea.Msg {
	return PromptInput(m.input)
}

func (m Model) PromptEditing() tea.Msg {
	return PromptEditing(m.editing)
}

func New() Model {
	ti := textinput.New()
	ti.Prompt = "» "

	return Model{
		textinput: ti,
		style:     lipgloss.NewStyle(),

		InputAbort:  key.NewBinding(key.WithKeys("esc")),
		InputAccept: key.NewBinding(key.WithKeys("enter")),
		InputShow:   key.NewBinding(key.WithKeys(":")),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.PromptInput)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.InputShow) && !m.editing {
			m.editing = true
			m.textinput.Focus()
			m.textinput, cmd = m.textinput.Update(msg)
			m.textinput.SetValue("")
			cmd = m.PromptEditing
			return m, cmd

		} else if key.Matches(msg, m.InputAbort) {
			m.textinput.Reset()
			m.textinput.Blur()
			m.editing = false
			cmd = m.PromptEditing
			return m, cmd
		} else if key.Matches(msg, m.InputAccept) {
			m.input = m.textinput.Value()
			m.textinput.Reset()
			m.textinput.Blur()
			m.editing = false
			cmds = append(cmds, m.PromptInput)
			cmds = append(cmds, m.PromptEditing)
		} else {
			m.textinput, cmd = m.textinput.Update(msg)
			return m, cmd
		}

	}
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.editing {
		return m.style.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.textinput.View(),
			),
		)
	}

	return ""
}
