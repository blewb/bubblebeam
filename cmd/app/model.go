package main

import (
	"fmt"
	"time"

	"github.com/blewb/bubblebeam/span"
	"github.com/blewb/bubblebeam/stream"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type modelState int

const (
	StateLoading modelState = iota
	StateSelectDate
	StateListEntries
	StateSelectJob
	StateConfirm
)

type model struct {
	data         span.Span
	api          *stream.API
	day          int
	state        modelState
	table        table.Model
	paginator    paginator.Model
	dates        []span.Datestamp
	today        string
	selectedDate int
	width        int
	searchInput  textinput.Model
	searchJobs   []stream.Job
	searchTable  table.Model
}

func initialModel(sp span.Span, api *stream.API, launchState modelState) model {

	now := time.Now()
	dates, sd := span.GetDatestamps(now, SELECTION_DAYS)
	today := now.Format(time.DateOnly)

	m := model{
		data:         sp,
		api:          api,
		day:          0,
		state:        launchState,
		paginator:    BuildPaginator(sp),
		dates:        dates,
		today:        today,
		selectedDate: sd,
		searchInput:  BuildTextinput("Job Number/Name/Company", 32),
		searchJobs:   make([]stream.Job, 0),
	}

	rows := m.GetEntryRows()
	m.table = BuildTable(rows)

	m.searchTable = BuildSearchTable()

	return m

}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) GetEntryRows() []table.Row {

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
			"--",
		})
	}

	return rows

}

// Sum of the fixed width columns below, plus 2 per border between/around
const tableFixedSpace = 43

func getColumns(w1, w2 int) []table.Column {
	return []table.Column{
		{Title: "#", Width: 3},
		{Title: "Start", Width: 6},
		{Title: "End", Width: 6},
		{Title: "Time", Width: 6},
		{Title: "Description", Width: w1},
		{Title: "Tag", Width: w2},
		{Title: "Status", Width: 6},
	}
}

func getSearchColumns() []table.Column {
	return []table.Column{
		{Title: "#", Width: 3},
		{Title: "Number", Width: 6},
		{Title: "Name", Width: 24},
		{Title: "Client", Width: 24},
	}
}
