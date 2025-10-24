package main

import "fmt"

const (
	arrowLeft  = '‹'
	arrowRight = '›'
)

func (m model) View() string {

	switch m.state {
	case StateStart:
		return m.ViewStart()
	case StateList:
		return m.ViewList()
	case StateConfirm:
		return m.ViewConfirm()
	}

	return baseStyle.Render(m.table.View())

}

func (m model) ViewStart() string {
	return "TODO"
}

func (m model) ViewList() string {

	s := appTitle + "\n\n"

	n := len(m.data.Days)
	singleDay := n == 1

	s += fmt.Sprintf("%c %s %c\n", arrowLeft, m.data.Days[m.day].Weekday.String(), arrowRight)

	if !singleDay {
		s += m.paginator.View()
		s += "\n"
	}

	s += "\n"

	s += baseStyle.Render(m.table.View())

	return s

}

func (m model) ViewConfirm() string {
	return "TODO"
}
