package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	arrowLeft  = '‹'
	arrowRight = '›'

	weekdaySize = 10 // Enough space for longest day (Wed) plus one space
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
	case StateSelectItem:
		return m.ViewSelectItem()
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
	s += baseStyle.Render(m.table.View()) + "\n"

	return s

}

func (m model) ViewSelectJob() string {

	var centre string
	if len(m.data.Days) > 1 {
		centre = fmt.Sprintf("%c %s %c", arrowLeft, m.paginator.View(), arrowRight)
	}

	right := fmt.Sprintf("%s", m.data.Days[m.day].Weekday.String())

	s := Header(m.width, centre, right) + "\n"
	s += baseStyle.Render(m.table.View())

	s += "\n"
	s += m.searchInput.View()
	s += "\n"

	if len(m.searchJobs) > 0 {
		s += baseStyle.Render(m.searchTable.View()) + "\n"
	}

	return s

}

func (m model) ViewSelectItem() string {

	var centre string
	if len(m.data.Days) > 1 {
		centre = fmt.Sprintf("%c %s %c", arrowLeft, m.paginator.View(), arrowRight)
	}

	right := fmt.Sprintf("%s", m.data.Days[m.day].Weekday.String())

	s := Header(m.width, centre, right) + "\n"
	s += baseStyle.Render(m.table.View())

	job := m.searchJobs[m.searchTable.Cursor()]

	s += "\n"
	s += fmt.Sprintf("[%s] %s", job.Number, job.Name)
	s += "\n"

	if len(m.itemList) > 0 {
		s += baseStyle.Render(m.itemTable.View()) + "\n"
	}

	return s

}

func (m model) ViewConfirm() string {
	return "Confirm"
}
