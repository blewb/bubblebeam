package main

import "fmt"

func (m model) View() string {

	s := fmt.Sprintln("Bubblebeam")
	s += "==========\n\n"

	for _, day := range m.data.Days {

		s += fmt.Sprintln(day.Weekday.String())

		for _, ent := range day.Entries {
			s += fmt.Sprintln("\t", ent)
		}

	}

	return s
}
