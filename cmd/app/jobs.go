package main

import (
	"fmt"
	"strings"

	"github.com/blewb/bubblebeam/stream"
	"github.com/charmbracelet/bubbles/table"
)

func (m *model) SearchJobs() {

	found := make([]stream.Job, 0)
	term := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	m.searchTable.SetCursor(0)

	if len(term) < 2 {
		m.searchJobs = found
		m.searchTable.SetRows([]table.Row{})
		return
	}

	for _, job := range m.api.GetJobs() {

		if strings.Contains(job.Search, term) {
			found = append(found, job)
		}

	}

	rows := make([]table.Row, 0, len(found))
	for i, fj := range found {
		rows = append(rows, []string{
			fmt.Sprintf("%d", i+1),
			fj.Number,
			fj.Name,
			fj.Company,
		})
	}

	m.searchJobs = found
	m.searchTable.SetRows(rows)

}
