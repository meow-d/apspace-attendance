package main

import (
	"fmt"
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

	quitMsg := quitMsgStyle.Bold(true).Render("q") + quitMsgStyle.Render(" to quit")

	return fmt.Sprintf(
		"%sEnter your attendance code\n%s \n \n%s",
		errorMsg,
		m.input.View(),
		quitMsg,
	)
}

func inputUpdate(m model, msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "q":
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.code = m.input.Value()
			m.view = 1
			return m, tea.Batch(m.spinner.Tick, attendance(m))
		}
	}

	m.input, cmd = m.input.Update(msg)
	m.input.SetValue(filterNumbers(m.input.Value()))
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
