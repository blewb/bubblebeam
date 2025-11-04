package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	arrowLeft  = '‹'
	arrowRight = '›'
)

func (m model) View() string {

	s := appTitle + "\n\n"

	switch m.state {
	case StateSelectDate:
		return s + m.ViewSelectDate()
	case StateListEntries:
		return s + m.ViewListEntries()
	case StateSelectJob:
		return s + m.ViewSelectJob()
	case StateConfirm:
		return s + m.ViewConfirm()
	}

	return s + m.ViewLoading()

}

func (m model) ViewLoading() string {
	return "Booting"
}

func (m model) ViewSelectDate() string {

	s := "\n"

	week := make([]string, 0, 7)

	for d, date := range m.dates {

		if d == m.selectedDate {
			week = append(week, selectedCardStyle.Render(date.Friendly)+" ")
		} else {
			week = append(week, cardStyle.Render(date.Friendly)+" ")
		}

		if len(week) == 7 {
			s += lipgloss.JoinHorizontal(lipgloss.Top, week...) + "\n"
			week = week[:0]
		}

	}

	return s

}

func (m model) ViewListEntries() string {

	n := len(m.data.Days)
	singleDay := n == 1

	s := fmt.Sprintf("%c %s %c\n", arrowLeft, m.data.Days[m.day].Weekday.String(), arrowRight)

	if !singleDay {
		s += m.paginator.View()
		s += "\n"
	}

	s += "\n"

	s += baseStyle.Render(m.table.View())

	return s

}

func (m model) ViewSelectJob() string {
	return "Select Job"
}

func (m model) ViewConfirm() string {
	return "Confirm"
}
