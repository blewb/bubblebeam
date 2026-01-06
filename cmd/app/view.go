package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	arrowLeft  = '‹'
	arrowRight = '›'

	weekdaySize = 10 // Enough space for longest day (Wed) plus one space
)

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

func Header(width int, centre, right string) string {

	hWidth := width / 3

	headerLeft = headerLeft.Width(hWidth)
	headerCentre = headerCentre.Width(hWidth)
	headerRight = headerRight.Width(width - (2 * hWidth))

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		headerLeft.Render(headerTitle),
		headerCentre.Render(centre),
		headerRight.Render(right+" "), // Trailing space because it looks nicer
	)

}

func (m model) View() string {

	switch m.state {
	case StateSelectDate:
		return m.ViewSelectDate()
	case StateListEntries:
		return m.ViewListEntries()
	case StateSelectJob:
		return m.ViewSelectJob()
	case StateConfirm:
		return m.ViewConfirm()
	}

	return m.ViewLoading()

}

func (m model) ViewLoading() string {
	return "Booting"
}

func (m model) ViewSelectDate() string {

	s := "\n"

	row := m.selectedDate - 3

	for ; row < m.selectedDate; row++ {
		if row >= 0 {
			s += fmt.Sprintf("%s   %-10s %s\n", weekdayPadding, m.dates[row].Weekday.String(), m.dates[row].Friendly)
		} else {
			s += fmt.Sprintf("%s   ---\n", weekdayPadding)
		}
	}

	for d := range m.data.Days {
		s += fmt.Sprintf("%-10s | %-10s %s\n", m.data.Days[d].Weekday.String(), m.dates[row].Weekday.String(), m.dates[row].Friendly)
		row++
	}

	for t := 0; t < 3; t++ {
		if row < len(m.dates) {
			s += fmt.Sprintf("%s   %-10s %s\n", weekdayPadding, m.dates[row].Weekday.String(), m.dates[row].Friendly)
		} else {
			s += fmt.Sprintf("%s   ---\n", weekdayPadding)
		}
		row++
	}

	return s

}

func (m model) ViewListEntries() string {

	var centre string
	if len(m.data.Days) > 1 {
		centre = fmt.Sprintf("%c %s %c", arrowLeft, m.paginator.View(), arrowRight)
	}

	right := fmt.Sprintf("%s", m.data.Days[m.day].Weekday.String())

	s := Header(m.width, centre, right) + "\n"
	s += baseStyle.Render(m.table.View())

	return s

}

func (m model) ViewSelectJob() string {
	return "Select Job"
}

func (m model) ViewConfirm() string {
	return "Confirm"
}
