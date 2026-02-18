package main

import (
	"github.com/blewb/bubblebeam/stream"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.searchInput.Width = max(10, (m.width*2)/3-6)
		return m, nil

	case jobItemsMsg:
		if msg.err == nil {
			m.jobItemCache[msg.jobID] = msg.items
		}
		if m.selectedJob.ID == msg.jobID {
			m.itemLoading = false
			if msg.err == nil {
				m.itemList = msg.items
				m.itemError = ""
			} else {
				m.itemList = nil
				m.itemError = msg.err.Error()
			}
			m.itemCursor = 0
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	}

	switch m.state {
	case StateSelectDate:
		return m.updateSelectDate(msg)
	case StateMain:
		return m.updateMain(msg)
	}

	return m, nil

}

// ─── Date Selection ─────────────────────────────────────────

func (m *model) updateSelectDate(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.dateCursor > 0 {
				m.dateCursor--
			}

		case "down", "j":
			if m.dateCursor < len(m.data.Days)-1 {
				m.dateCursor++
			}

		case "left", "h":
			if len(m.dayDates) > 0 {
				m.dayDates[m.dateCursor] = m.dayDates[m.dateCursor].AddDate(0, 0, -7)
			}

		case "right", "l":
			if len(m.dayDates) > 0 {
				m.dayDates[m.dateCursor] = m.dayDates[m.dateCursor].AddDate(0, 0, 7)
			}

		case "enter":
			m.state = StateMain
			return m, nil

		}
	}

	return m, nil

}

// ─── Main Panel Layout ─────────────────────────────────────

func (m *model) updateMain(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			return m.cycleFocus(true)
		case "shift+tab":
			return m.cycleFocus(false)
		case "esc":
			return m.handleEsc()
		}
	}

	switch m.focus {
	case FocusEntries:
		return m.updateEntries(msg)
	case FocusJobs:
		return m.updateJobs(msg)
	case FocusItems:
		return m.updateItems(msg)
	}

	return m, nil

}

func (m *model) cycleFocus(forward bool) (tea.Model, tea.Cmd) {

	var cmd tea.Cmd

	if forward {
		switch m.focus {
		case FocusEntries:
			m.focus = FocusJobs
			cmd = m.searchInput.Focus()
		case FocusJobs:
			m.focus = FocusItems
			m.searchInput.Blur()
		case FocusItems:
			m.focus = FocusEntries
		}
	} else {
		switch m.focus {
		case FocusEntries:
			m.focus = FocusItems
		case FocusJobs:
			m.focus = FocusEntries
			m.searchInput.Blur()
		case FocusItems:
			m.focus = FocusJobs
			cmd = m.searchInput.Focus()
		}
	}

	return m, cmd

}

func (m *model) handleEsc() (tea.Model, tea.Cmd) {

	switch m.focus {
	case FocusEntries:
		return m, tea.Quit
	case FocusJobs:
		m.focus = FocusEntries
		m.searchInput.Blur()
	case FocusItems:
		m.focus = FocusJobs
		cmd := m.searchInput.Focus()
		return m, cmd
	}

	return m, nil

}

// ─── Panel 1: Entries ───────────────────────────────────────

func (m *model) updateEntries(msg tea.Msg) (tea.Model, tea.Cmd) {

	if len(m.data.Days) == 0 {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.entryCursor > 0 {
				m.entryCursor--
			}

		case "down", "j":
			if len(m.data.Days) > 0 && m.entryCursor < len(m.data.Days[m.day].Entries)-1 {
				m.entryCursor++
			}

		case "left", "h":
			if len(m.data.Days) > 1 {
				if m.day > 0 {
					m.day--
				} else {
					m.day = len(m.data.Days) - 1
				}
				m.entryCursor = 0
			}

		case "right", "l":
			if len(m.data.Days) > 1 {
				if m.day < len(m.data.Days)-1 {
					m.day++
				} else {
					m.day = 0
				}
				m.entryCursor = 0
			}

		case "enter":
			m.focus = FocusJobs
			m.searchInput.SetValue("")
			m.searchJobs = nil
			m.jobCursor = 0
			cmd := m.searchInput.Focus()
			return m, cmd

		}
	}

	return m, nil

}

// ─── Panel 2: Job Search ────────────────────────────────────

func (m *model) updateJobs(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.jobCursor > 0 {
				m.jobCursor--
			}
			return m, nil

		case "down", "j":
			if m.jobCursor < len(m.searchJobs)-1 {
				m.jobCursor++
			}
			return m, nil

		case "enter":
			if len(m.searchJobs) > 0 {
				return m.selectJob()
			}
			return m, nil

		}
	}

	before := m.searchInput.Value()
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	if m.searchInput.Value() != before {
		m.searchJobs = searchJobs(m.api, m.searchInput.Value())
		m.jobCursor = 0
	}

	return m, cmd

}

func (m *model) selectJob() (tea.Model, tea.Cmd) {

	job := m.searchJobs[m.jobCursor]
	m.selectedJob = job
	m.searchInput.Blur()
	m.focus = FocusItems
	m.itemCursor = 0
	m.itemError = ""

	if cached, ok := m.jobItemCache[job.ID]; ok {
		m.itemList = cached
		return m, nil
	}

	m.itemLoading = true
	m.itemList = nil
	return m, fetchJobItems(m.api, job.ID)

}

// ─── Panel 3: Job Items ────────────────────────────────────

func (m *model) updateItems(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			if m.itemCursor > 0 {
				m.itemCursor--
			}

		case "down", "j":
			if m.itemCursor < len(m.itemList)-1 {
				m.itemCursor++
			}

		case "enter":
			if len(m.itemList) > 0 {
				return m.confirmItem()
			}

		}
	}

	return m, nil

}

func (m *model) confirmItem() (tea.Model, tea.Cmd) {

	if len(m.data.Days) == 0 {
		return m, nil
	}

	day := m.data.Days[m.day]
	if m.entryCursor >= len(day.Entries) {
		return m, nil
	}

	item := m.itemList[m.itemCursor]
	key := [2]int{m.day, m.entryCursor}

	m.assignments[key] = Assignment{
		Job:     m.selectedJob,
		JobItem: item,
	}

	m.focus = FocusEntries
	m.searchInput.Blur()
	m.searchInput.SetValue("")
	m.searchJobs = nil
	m.itemList = nil
	m.selectedJob = stream.Job{}

	if m.entryCursor < len(day.Entries)-1 {
		m.entryCursor++
	}

	return m, nil

}
