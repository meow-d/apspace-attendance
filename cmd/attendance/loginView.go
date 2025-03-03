package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func loginView(m model) string {
	statusMessage := ""
	if m.status != "" {
		statusMessage = errorStyle.Render(m.status) + "\n\n"
	}

	help := [][2]string{{"tab/arrowkeys", "move cursor"}, {"enter", "save"}, {"esc", "quit"}}
	helpMsg := renderHelpMsg(help)

	return fmt.Sprintf("%s%s %s\n%s %s\n\n%s",
		statusMessage,
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
			return validateAndLogin(m)
		}
	}

	return m, tea.Batch(cmds...)
}

func validateAndLogin(m model) (model, tea.Cmd) {
	if err := validateUsername(m.usernameInput.Value()); err != nil {
		return m, func() tea.Msg {
			return statusMsg(err.Error())
		}
	}
	if err := validateExists(m.passwordInput.Value()); err != nil {
		return m, func() tea.Msg {
			return statusMsg(err.Error())
		}
	}

	m.client.login(m.usernameInput.Value(), m.passwordInput.Value())
	m.view = 0
	m.status = ""
	return m, textinput.Blink
}
