package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// For the main table
var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var (
	bubble1 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ccff"))
	bubble2 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00bbff"))
	bubble3 = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00aaff"))

	headerLeft   = lipgloss.NewStyle().AlignHorizontal(lipgloss.Left)
	headerCentre = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
	headerRight  = lipgloss.NewStyle().AlignHorizontal(lipgloss.Right)
)

var (
	weekdayPadding string = strings.Repeat(" ", weekdaySize)
	headerTitle    string = bubble1.Render("B") + bubble2.Render("U") + bubble3.Render("B") + bubble1.Render("B") + bubble2.Render("L") + bubble3.Render("E") + bubble1.Render("B") + bubble2.Render("E") + bubble3.Render("A") + bubble1.Render("M")
)
