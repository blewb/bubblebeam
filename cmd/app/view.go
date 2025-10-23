package main

import "fmt"

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

	s := "# BUBBLEBEAM\n############\n\n"

	n := len(m.data.Days)
	singleDay := n == 1

	s += fmt.Sprintf("%c %s %c\n", '<', m.data.Days[m.day].Weekday.String(), '>')

	if !singleDay {
		for d := 0; d < n; d++ {
			if d == m.day {
				s += "X "
			} else {
				s += "o "
			}
		}
		s += "\n"
	}

	s += "\n"

	s += baseStyle.Render(m.table.View())

	return s

}

func (m model) ViewConfirm() string {
	return "TODO"
}
