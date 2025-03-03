package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func inputView(m model) string {
	// TODO proper way to handle success messages
	statusMessage := ""
	if m.status == "success" {
		statusMessage = successStyle.Render("Success! Submit another one?") + "\n\n"
	} else if m.status != "" {
		statusMessage = errorStyle.Render(m.status) + "\n\n"
	}

	help := [][2]string{{"enter", "submit"}, {"l", "logout"}, {"esc", "quit"}}
	helpMsg := renderHelpMsg(help)

	return fmt.Sprintf(
		"%s%s\n%s \n \n%s",
		statusMessage,
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
			return attendance(m)
		case "l":
			m.client.logout()
			m.status = ""
			m.view = 3
			return m, textinput.Blink
		}
	}

	m.codeInput, cmd = m.codeInput.Update(msg)
	m.codeInput.SetValue(filterNumbers(m.codeInput.Value()))
	return m, cmd
}

func attendance(m model) (model, tea.Cmd) {
	if err := validateCode(m.codeInput.Value()); err != nil {
		return m, func() tea.Msg {
			return statusMsg(err.Error())
		}
	}

	m.code = m.codeInput.Value()
	m.view = 1
	attendanceFunc := func() tea.Msg {
		err := m.client.submitAttendance(m.code)
		if err != nil {
			return statusMsg(err.Error())
		}
		return statusMsg("success")
	}
	return m, tea.Batch(attendanceFunc, m.spinner.Tick)
}
