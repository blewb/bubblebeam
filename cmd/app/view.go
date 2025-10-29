package main

import "fmt"

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
	return "Select Date"
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
