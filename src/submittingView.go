package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func submittingView(m model) string {
	return fmt.Sprintf("%s Submitting code %s...", m.spinner.View(), m.code)
}

func submittingUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		m.status = string(msg)
		m.input.SetValue("")
		m.view = 0
		return m, textinput.Blink
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}
