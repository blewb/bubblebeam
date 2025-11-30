package main

import (
	"fmt"
	"strings"
)

const (
	arrowLeft  = '‹'
	arrowRight = '›'

	weekdaySize = 10 // Enough space for longest day (Wed) plus one space
)

var weekdayPadding string = strings.Repeat(" ", weekdaySize)

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
