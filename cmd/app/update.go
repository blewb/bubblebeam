package main

import (
	"fmt"

	"github.com/blewb/bubblebeam/span"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		rem := msg.Width - tableFixedSpace
		larger := int(float64(rem) * 0.75)
		smaller := rem - larger

		cols := getColumns(larger, smaller)
		m.table.SetColumns(cols)
		m.table.SetHeight(max(7, (msg.Height/2)-3))
		m.width = msg.Width

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		}
	}

	switch m.state {
	case StateSelectDate:
		return m.UpdateSelectDate(msg)
	case StateListEntries:
		return m.UpdateListEntries(msg)
	case StateSelectJob:
		return m.UpdateSelectJob(msg)
	case StateSelectItem:
		return m.UpdateSelectItem(msg)
	}

	return m, cmd
}

func (m *model) UpdateSelectDate(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "left":
			if m.selectedDate-7 >= 0 {
				m.selectedDate -= 7
			}

		case "right":
			if m.selectedDate+7 < len(m.dates) {
				m.selectedDate += 7
			}

		case "up":
			if m.selectedDate > 0 {
				m.selectedDate--
			}

		case "down":
			if m.selectedDate < len(m.dates)-1 {
				m.selectedDate++
			}

		}
	}

	return m, cmd

}

func (m *model) UpdateListEntries(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "left":
			if len(m.data.Days) > 1 {
				if m.day > 0 {
					m.day--
				} else {
					m.day = len(m.data.Days) - 1
				}
				m.UpdateTable()
			}

		case "right":
			if len(m.data.Days) > 1 {
				if m.day < len(m.data.Days)-1 {
					m.day++
				} else {
					m.day = 0
				}
				m.UpdateTable()
			}

		case "space", "enter":
			m.state = StateSelectJob
			m.searchInput.Focus()
		}
	}

	m.paginator.Page = m.day
	m.table, cmd = m.table.Update(msg)

	return m, cmd

}

func (m *model) UpdateSelectJob(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+left", "ctrl+up":
			m.state = StateListEntries
			m.searchInput.Blur()

		case "space", "enter":
			m.state = StateSelectItem
			m.searchInput.Blur()

			job := m.searchJobs[m.searchTable.Cursor()]
			var err error
			m.itemList, err = m.api.GetJobItems(job.ID)
			if err != nil {
				// oh crap
			}
			rows := make([]table.Row, 0, len(m.itemList))
			for i, itm := range m.itemList {
				rows = append(rows, []string{
					fmt.Sprintf("%d", i+1),
					itm.Name,
					itm.Description,
					fmt.Sprintf("%s / %s",
						span.DurationAsString(itm.LoggedMinutes),
						span.DurationAsString(itm.PlannedMinutes),
					),
				})
			}
			m.itemTable.SetRows(rows)
		}
	}

	before := m.searchInput.Value()
	m.searchInput, cmd = m.searchInput.Update(msg)

	if before != m.searchInput.Value() {
		m.SearchJobs()
	}

	m.searchTable, cmd = m.searchTable.Update(msg)
	return m, cmd

}

func (m *model) UpdateSelectItem(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+left", "ctrl+up":
			m.state = StateSelectJob
			m.searchInput.Focus()

		}
	}

	m.itemTable, cmd = m.itemTable.Update(msg)
	return m, cmd

}

func (m *model) UpdateTable() {

	m.table.SetRows(m.GetEntryRows())
	m.table.SetCursor(0)
	m.table.MoveUp(100) // Hacky

}
