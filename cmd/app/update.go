package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "left":
			if m.day > 0 {
				m.day--
			} else {
				m.day = len(m.data.Days) - 1
			}
			m.UpdateTable()

		case "right":
			if m.day < len(m.data.Days)-1 {
				m.day++
			} else {
				m.day = 0
			}
			m.UpdateTable()

		}
	}

	m.paginator.Page = m.day
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

// Doesn't work just yet
// Of course it doesn't if you pass a copy, dummy...
func (m *model) UpdateTable() {

	n := len(m.data.Days[m.day].Entries)
	rows := make([]table.Row, 0, n)

	for e, entry := range m.data.Days[m.day].Entries {
		rows = append(rows, []string{
			fmt.Sprintf("%d", e+1),
			entry.Start.Render(),
			entry.End.Render(),
			entry.DurationString(),
			entry.Description,
			entry.Tag,
		})
	}

	m.table.SetRows(rows)
	m.table.SetHeight(n + 2)
	m.table.SetCursor(0)

}
