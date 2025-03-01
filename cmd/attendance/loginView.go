package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func loginView(m model) string {
	help := [][2]string{{"tab/arrowkeys", "move cursor"}, {"enter", "save"}, {"esc", "quit"}}
	helpMsg := renderHelpMsg(help)

	return fmt.Sprintf("%s %s\n%s %s\n\n%s",
		labelStyle.Render("Username"),
		m.usernameInput.View(),
		labelStyle.Render("Password"),
		m.passwordInput.View(),
		helpMsg,
	)
}

func loginUpdate(m model, msg tea.Msg) (model, tea.Cmd) {
	cmds := make([]tea.Cmd, 2)
	m.usernameInput, cmds[0] = m.usernameInput.Update(msg)
	m.passwordInput, cmds[1] = m.passwordInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:

		// as far as a i understand from their examples, return textinput.Blink shouldn't be needed here..
		// strange
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			if m.usernameInput.Focused() {
				m.usernameInput.Blur()
				m.passwordInput.Focus()
				return m, textinput.Blink
			} else {
				m.passwordInput.Blur()
				m.usernameInput.Focus()
				return m, textinput.Blink
			}

		case "enter":
			m.client.login(m.usernameInput.Value(), m.passwordInput.Value())
			m.view = 0
			return m, textinput.Blink
		}
	}

	return m, tea.Batch(cmds...)
}
