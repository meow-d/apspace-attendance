package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func inputView(m model) string {
	errorMsg := ""
	// TODO urm..... urm..... stop doing this....?
	if m.status == "success" {
		errorMsg = successStyle.Render("Success! Submit another one?") + "\n\n"
	} else if m.status != "" {
		errorMsg = errorStyle.Render(m.status) + "\n\n"
	}

	help := [][2]string{{"enter", "submit"}, {"l", "logout"}, {"esc", "quit"}}
	helpMsg := renderHelpMsg(help)

	return fmt.Sprintf(
		"%s%s\n%s \n \n%s",
		errorMsg,
		labelStyle.Render("Attendance code"),
		m.codeInput.View(),
		helpMsg,
	)
}

func inputUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.code = m.codeInput.Value()
			m.view = 1
			return m, tea.Batch(m.spinner.Tick, attendance(m))
		case "l":
			m.client.logout()
			m.view = 3
			return m, textinput.Blink
		}
	}

	m.codeInput, cmd = m.codeInput.Update(msg)
	m.codeInput.SetValue(filterNumbers(m.codeInput.Value()))
	return m, cmd
}

func attendance(m model) tea.Cmd {
	return func() tea.Msg {
		err := m.client.submitAttendance(m.code)
		if err != nil {
			return statusMsg(err.Error())
		}
		return statusMsg("success")
	}
}
