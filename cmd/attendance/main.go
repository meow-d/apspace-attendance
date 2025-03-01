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
	// debug
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// check args
	args := os.Args
	if !(len(args) <= 2) {
		fmt.Println("Too many arguments")
		return
	}

	// main app
	model := initialModel(args)
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

// model

type statusMsg string

// TODO: use enum for view?
type model struct {
	client        *Client
	view          int // 0: input, 1: submitting, 3: login
	usernameInput textinput.Model
	passwordInput textinput.Model
	codeInput     textinput.Model
	spinner       spinner.Model
	code          string
	status        string
}

func initialModel(args []string) model {
	c := NewClient()

	ui := textinput.New()
	ui.Placeholder = "TP000000"
	ui.CharLimit = 8
	ui.Width = 20
	ui.Prompt = ""
	ui.SetValue(c.Auth.Username)
	ui.Focus()

	pi := textinput.New()
	pi.Placeholder = "•••••••••••••"
	pi.EchoMode = textinput.EchoPassword
	pi.EchoCharacter = '•'
	pi.Prompt = ""
	pi.CharLimit = 40
	pi.Width = 20

	ci := textinput.New()
	ci.Placeholder = "000"
	ci.CharLimit = 3
	ci.Width = 3
	ci.Prompt = ""
	ci.Validate = validateCode
	ci.Focus()

	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	v := 0
	if len(args) == 2 {
		ci.SetValue(args[1])
		v = 1
	}
	if c.Auth.Username == "" || c.Auth.Password == "" {
		v = 3
	}

	return model{
		client:        c,
		view:          v,
		usernameInput: ui,
		passwordInput: pi,
		codeInput:     ci,
		spinner:       s,
		code:          ci.Value(),
	}
}

func (m model) Init() tea.Cmd {
	if m.view == 1 {
		return tea.Batch(attendance(m), m.spinner.Tick)
	}
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	switch m.view {
	case 3:
		return loginUpdate(m, msg)
	case 0:
		return inputUpdate(m, msg)
	default:
		return submittingUpdate(m, msg)
	}
}

func (m model) View() string {
	var s string
	switch m.view {
	case 3:
		s = loginView(m)
	case 0:
		s = inputView(m)
	default:
		s = submittingView(m)
	}
	return appStyle.Render(s)
}
