package main

import (
	"github.com/blewb/bubblebeam/span"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	data   span.Span
	day    int
	cursor int
}

func initialModel(sp span.Span) model {
	return model{
		data:   sp,
		day:    0,
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
