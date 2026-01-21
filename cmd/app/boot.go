package main

import (
	"github.com/blewb/bubblebeam/span"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func BuildTable(rows []table.Row) table.Model {

	t := table.New(
		table.WithColumns(getColumns(30, 10)), // Arbitrary starting size, as it gets dynamically resized
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

	return t

}

// TODO: Make this a list instead
func BuildSearchTable() table.Model {

	t := table.New(
		table.WithColumns(getSearchColumns()),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#00bbff")).
		Bold(true)
	t.SetStyles(s)

	return t

}

func BuildPaginator(sp span.Span) paginator.Model {

	p := paginator.New()
	p.Type = paginator.Dots
	p.PerPage = 1
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	p.SetTotalPages(len(sp.Days))

	return p

}

func BuildTextinput(plc string, size int) textinput.Model {

	ti := textinput.New()
	ti.Placeholder = plc
	ti.CharLimit = size
	ti.Width = size

	return ti

}
