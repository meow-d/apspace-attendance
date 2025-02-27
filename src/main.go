package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	args := os.Args

	if len(args) <= 2 {
		auth, err := getAuth()
		if err != nil {
			fmt.Println(err)
			return
		}

		model := initialModel(args, auth)
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			log.Fatal(err)
		}

	} else {
		fmt.Println("Too many arguments")
		return
	}
}

// types
type statusMsg string

// model
type model struct {
	client  Client
	view    int
	input   textinput.Model
	spinner spinner.Model
	code    string
	status  string
}

func initialModel(args []string, auth Auth) model {
	c := NewClient(auth)

	i := textinput.New()
	i.Placeholder = "000"
	i.Focus()
	i.CharLimit = 3
	i.Width = 3
	i.Prompt = ""
	i.Validate = validateCode

	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	v := 0
	if len(args) == 2 {
		i.SetValue(args[1])
		v = 1
	}

	return model{
		client:  c,
		view:    v,
		spinner: s,
		input:   i,
	}
}

func (m model) Init() tea.Cmd {
	if m.view == 0 {
		return textinput.Blink
	}
	return tea.Batch(attendance(m), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "ctrl+c" {
			return m, tea.Quit
		}
	}

	switch m.view {
	case 0:
		return inputUpdate(m, msg)
	default:
		return submittingUpdate(m, msg)
	}
}

func (m model) View() string {
	var s string
	switch m.view {
	case 0:
		s = inputView(m)
	case 1:
		s = submittingView(m)
	default:
		s = "Invalid view"
	}
	return appStyle.Render(s)
}
