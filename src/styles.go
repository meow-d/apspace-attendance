package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle     = lipgloss.NewStyle().Padding(1, 2)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Italic(true)  // red
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true) // green
	quitMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))             // grey
)
