package main

import (
	"fmt"
	"time"

	"github.com/blewb/bubblebeam/span"
	"github.com/blewb/bubblebeam/stream"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

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

func initialModel(sp span.Span, api *stream.API, launchState modelState) model {

	columns := getColumns(30, 10)

	n := len(sp.Days[0].Entries)
	rows := make([]table.Row, 0, n)

	for e, entry := range sp.Days[0].Entries {
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

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("#00aaff")).
		Bold(false)
	t.SetStyles(s)

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 1
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(sp.Days))

	now := time.Now()
	dates, sd := span.GetDatestamps(now, SELECTION_DAYS)
	today := now.Format(time.DateOnly)

	return model{
		data:         sp,
		api:          api,
		day:          0,
		state:        launchState,
		table:        t,
		paginator:    p,
		dates:        dates,
		today:        today,
		selectedDate: sd,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
