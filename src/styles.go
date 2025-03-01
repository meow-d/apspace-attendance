package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle     = lipgloss.NewStyle().Padding(1, 2)
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF06B7")).Bold(true) // pink
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Italic(true)     // red
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Italic(true)    // green
	helpMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))                // grey
)

func renderHelpMsg(help [][2]string) string {
	var s string
	for _, h := range help {
		s += helpMsgStyle.Bold(true).Render(h[0]) + helpMsgStyle.Render(" "+h[1]+"   ")
	}
	return s
}
