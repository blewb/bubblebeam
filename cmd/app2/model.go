package main

import (
	"time"

	"github.com/blewb/bubblebeam/span"
	"github.com/blewb/bubblebeam/stream"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type appState int

const (
	StateSelectDate appState = iota
	StateMain
)

type panelFocus int

const (
	FocusEntries panelFocus = iota
	FocusJobs
	FocusItems
)

const (
	MIN_WIDTH  = 80
	MIN_HEIGHT = 20
)

type Assignment struct {
	Job     stream.Job
	JobItem stream.JobItem
}

type jobItemsMsg struct {
	jobID int64
	items []stream.JobItem
	err   error
}

type model struct {
	data  span.Span
	api   *stream.API
	state appState
	focus panelFocus

	dayDates   []time.Time
	dateCursor int

	day         int
	entryCursor int

	searchInput textinput.Model
	searchJobs  []stream.Job
	jobCursor   int

	itemList    []stream.JobItem
	itemCursor  int
	itemLoading bool
	selectedJob stream.Job

	assignments  map[[2]int]Assignment
	jobItemCache map[int64][]stream.JobItem

	width  int
	height int
}

func initialModel(sp span.Span, api *stream.API) model {

	ti := textinput.New()
	ti.Placeholder = "Type to search jobs..."
	ti.CharLimit = 64
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(colorBlue2)

	m := model{
		data:         sp,
		api:          api,
		state:        StateSelectDate,
		focus:        FocusEntries,
		searchInput:  ti,
		searchJobs:   make([]stream.Job, 0),
		assignments:  make(map[[2]int]Assignment),
		jobItemCache: make(map[int64][]stream.JobItem),
	}

	m.calculateDates()

	return m

}

func (m *model) calculateDates() {

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	m.dayDates = make([]time.Time, len(m.data.Days))

	if len(m.data.Days) == 0 {
		return
	}

	firstWeekday := m.data.Days[0].Weekday
	diff := int(today.Weekday()) - int(firstWeekday)
	if diff < 0 {
		diff += 7
	}
	m.dayDates[0] = today.AddDate(0, 0, -diff)

	for i := 1; i < len(m.data.Days); i++ {
		prevDate := m.dayDates[i-1]
		targetWeekday := m.data.Days[i].Weekday
		dayDiff := int(targetWeekday) - int(prevDate.Weekday())
		if dayDiff <= 0 {
			dayDiff += 7
		}
		m.dayDates[i] = prevDate.AddDate(0, 0, dayDiff)
	}

}

func (m model) Init() tea.Cmd {
	return nil
}

func fetchJobItems(api *stream.API, jobID int64) tea.Cmd {
	return func() tea.Msg {
		items, err := api.GetJobItems(jobID)
		return jobItemsMsg{jobID: jobID, items: items, err: err}
	}
}
