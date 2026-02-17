package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	colorBlue1 = lipgloss.Color("#00ccff")
	colorBlue2 = lipgloss.Color("#00bbff")
	colorBlue3 = lipgloss.Color("#00aaff")
	colorGrey  = lipgloss.Color("240")
	colorWhite = lipgloss.Color("255")
	colorDim   = lipgloss.Color("245")
)

var (
	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#00aaff")).
			Foreground(lipgloss.Color("#000000")).
			Bold(true)

	headerRowStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Bold(true)

	cursorStyle = lipgloss.NewStyle().
			Foreground(colorBlue1).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)
)

var headerTitle = bubble1.Render("B") + bubble2.Render("U") + bubble3.Render("B") + bubble1.Render("B") + bubble2.Render("L") + bubble3.Render("E") + bubble1.Render("B") + bubble2.Render("E") + bubble3.Render("A") + bubble1.Render("M")

var (
	bubble1 = lipgloss.NewStyle().Bold(true).Foreground(colorBlue1)
	bubble2 = lipgloss.NewStyle().Bold(true).Foreground(colorBlue2)
	bubble3 = lipgloss.NewStyle().Bold(true).Foreground(colorBlue3)
)

func renderTitle() string {
	return headerTitle
}

func renderPanelHeader(title, icon string, width int) string {

	titleRendered := lipgloss.NewStyle().
		Foreground(colorBlue2).
		Bold(true).
		Render(title)

	iconRendered := lipgloss.NewStyle().
		Foreground(colorDim).
		Render(icon)

	gap := width - lipgloss.Width(titleRendered) - lipgloss.Width(iconRendered)
	if gap < 0 {
		gap = 0
	}

	return titleRendered + strings.Repeat(" ", gap) + iconRendered

}

func renderPanel(title, icon, content string, totalW, totalH int, focused bool) string {

	innerW := max(1, totalW-2)
	innerH := max(1, totalH-2)

	header := renderPanelHeader(title, icon, innerW)

	bodyH := max(0, innerH-1)

	body := lipgloss.NewStyle().
		Width(innerW).
		Height(bodyH).
		Render(content)

	full := header + "\n" + body

	borderColor := colorGrey
	if focused {
		borderColor = colorBlue2
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(innerW).
		Render(full)

}
