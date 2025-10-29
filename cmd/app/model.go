package main

import (
	"fmt"

	"github.com/blewb/bubblebeam/span"
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
	data      span.Span
	day       int
	state     modelState
	table     table.Model
	paginator paginator.Model
}

func initialModel(sp span.Span, launchState modelState) model {

	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Start", Width: 6},
		{Title: "End", Width: 6},
		{Title: "Time", Width: 6},
		{Title: "Description", Width: 64},
		{Title: "Tag", Width: 16},
	}

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
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(n+1),
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

	return model{
		data:      sp,
		day:       0,
		state:     launchState,
		table:     t,
		paginator: p,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
