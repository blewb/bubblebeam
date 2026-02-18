package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/blewb/bubblebeam/span"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func (m model) View() string {

	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	if m.width < MIN_WIDTH || m.height < MIN_HEIGHT {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			dimStyle.Render("Terminal too small. Please resize to at least 80×20."),
		)
	}

	switch m.state {
	case StateSelectDate:
		return m.viewSelectDate()
	case StateMain:
		return m.viewMain()
	}

	return "Loading..."

}

// ─── Date Selection ─────────────────────────────────────────

func (m model) viewSelectDate() string {

	var b strings.Builder

	title := renderTitle()
	b.WriteString(title)
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Bold(true).Foreground(colorBlue1).Render("  Select Dates"))
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("  Assign a date to each day. ↑↓ select · ←→ adjust ±7 days · Enter confirm"))
	b.WriteString("\n\n")

	for i, day := range m.data.Days {

		cursor := "  "
		if i == m.dateCursor {
			cursor = cursorStyle.Render("▸ ")
		}

		weekday := fmt.Sprintf("%-10s", day.Weekday.String())
		date := m.dayDates[i].Format("Mon Jan 2, 2006")

		if i == m.dateCursor {
			b.WriteString(cursor + selectedStyle.Render(fmt.Sprintf(" %-10s  %s ", weekday, date)))
		} else {
			b.WriteString(cursor + fmt.Sprintf(" %-10s  %s", weekday, dimStyle.Render(date)))
		}
		b.WriteString("\n")

	}

	if len(m.data.Days) == 0 {
		b.WriteString(dimStyle.Render("  No days found in the input file."))
		b.WriteString("\n")
	}

	content := b.String()

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)

}

// ─── Main Three-Panel Layout ────────────────────────────────

func (m model) viewMain() string {

	title := renderTitle()
	dayTabs := m.dayTabs()
	headerLine := renderHeaderLine(title, dayTabs, m.width)

	helpLine := m.helpText()

	availH := max(6, m.height-3)

	panel1H := availH / 2
	panel23H := availH - panel1H

	panel2W := (m.width * 2) / 3
	panel3W := m.width - panel2W

	p1BodyH := max(1, panel1H-3)
	p1Content := m.viewEntries(m.width-2, p1BodyH)
	p1Icon := m.entriesIcon()
	p1 := renderPanel("Entries", p1Icon, p1Content, m.width, panel1H, m.focus == FocusEntries)

	p23BodyH := max(1, panel23H-3)
	p2Content := m.viewJobSearch(panel2W-2, p23BodyH)
	p2Icon := m.jobSearchIcon()
	p2 := renderPanel("Search", p2Icon, p2Content, panel2W, panel23H, m.focus == FocusJobs)

	p3Content := m.viewJobItems(panel3W-2, p23BodyH)
	p3Icon := m.jobItemsIcon()
	p3 := renderPanel("Job Items", p3Icon, p3Content, panel3W, panel23H, m.focus == FocusItems)

	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, p2, p3)

	return headerLine + "\n" + p1 + "\n" + bottomRow + "\n" + helpLine

}

func (m model) helpText() string {

	switch m.state {
	case StateMain:
		switch m.focus {
		case FocusEntries:
			return dimStyle.Render(" ↑↓ navigate · ←→ change day · Enter select · Tab next panel · Esc quit")
		case FocusJobs:
			return dimStyle.Render(" Type to search · ↑↓ navigate results · Enter select · Esc back")
		case FocusItems:
			return dimStyle.Render(" ↑↓ navigate · Enter confirm · Esc back")
		}
	}

	return ""

}

// ─── Panel 1: Entries ───────────────────────────────────────

func (m model) entriesIcon() string {

	if len(m.data.Days) == 0 {
		return ""
	}

	d := m.data.Days[m.day]
	total := span.DurationAsString(d.Duration)

	if len(m.data.Days) > 1 {
		if m.day < len(m.dayDates) {
			return fmt.Sprintf("%s · %s", total, formatDayDate(m.dayDates[m.day]))
		}
		return fmt.Sprintf("%s · Day %d/%d", total, m.day+1, len(m.data.Days))
	}

	if len(m.dayDates) > 0 {
		return fmt.Sprintf("%s · %s", total, formatDayDate(m.dayDates[0]))
	}

	return total

}

func (m model) viewEntries(w, h int) string {

	if len(m.data.Days) == 0 {
		return dimStyle.Render("No entries loaded")
	}

	var b strings.Builder

	day := m.data.Days[m.day]
	entries := day.Entries

	tableW := w - 2
	if tableW < 1 {
		tableW = 1
	}

	numW := 3
	startW := 6
	endW := 6
	durW := 6
	tagW := 10
	statusW := 12
	separators := 6
	fixed := numW + startW + endW + durW + tagW + statusW + separators
	descW := tableW - fixed
	if descW < 10 {
		descW = 10
	}

	hdr := fmt.Sprintf(" %-*s %-*s %-*s %-*s %-*s %-*s %-*s ",
		numW, "#",
		startW, "Start",
		endW, "End",
		durW, "Time",
		descW, "Desc",
		tagW, "Tag",
		statusW, "Status",
	)
	b.WriteString(headerRowStyle.Render(hdr) + "\n")
	b.WriteString(dimStyle.Render(strings.Repeat("─", w)) + "\n")
	h -= 2

	start, end := visibleRange(m.entryCursor, len(entries), h)

	for i := start; i < end; i++ {

		e := entries[i]
		status := m.getEntryStatus(m.day, i)

		desc := e.Description
		if len(desc) > descW {
			desc = desc[:descW-2] + ".."
		}

		tag := e.Tag
		if len(tag) > tagW {
			tag = tag[:tagW-2] + ".."
		}

		row := fmt.Sprintf(" %-*d %-*s %-*s %-*s %-*s %-*s %-*s ",
			numW, i+1,
			startW, e.Start.Render(),
			endW, e.End.Render(),
			durW, e.DurationString(),
			descW, desc,
			tagW, tag,
			statusW, status,
		)

		if i == m.entryCursor {
			b.WriteString(selectedStyle.Render(row))
		} else {
			b.WriteString(row)
		}

		if i < end-1 {
			b.WriteString("\n")
		}

	}

	return b.String()

}

func (m model) getEntryStatus(dayIdx, entryIdx int) string {

	key := [2]int{dayIdx, entryIdx}
	if a, ok := m.assignments[key]; ok {
		return "✓ " + a.Job.Number
	}
	return "–"

}

// ─── Panel 2: Job Search ────────────────────────────────────

func (m model) jobSearchIcon() string {

	if len(m.searchJobs) > 0 {
		return fmt.Sprintf("%d found", len(m.searchJobs))
	}
	return ""

}

func (m model) viewJobSearch(w, h int) string {

	var b strings.Builder

	b.WriteString(padLine(w, m.searchInput.View()))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.Repeat("─", w)))
	b.WriteString("\n")
	h -= 2

	if len(m.searchJobs) == 0 {
		if len(m.searchInput.Value()) >= 2 {
			b.WriteString(dimStyle.Render(padLine(w, "No matches")))
		} else {
			b.WriteString(dimStyle.Render(padLine(w, "Type at least 2 characters")))
		}
		return b.String()
	}

	tableW := w - 2
	if tableW < 1 {
		tableW = 1
	}

	numW := 3
	numberW := 8
	separators := 3
	fixed := numW + numberW + separators
	remaining := tableW - fixed
	if remaining < 4 {
		remaining = 4
	}
	nameW := (remaining * 2) / 3
	compW := remaining - nameW

	start, end := visibleRange(m.jobCursor, len(m.searchJobs), h)

	for i := start; i < end; i++ {

		job := m.searchJobs[i]

		name := job.Name
		if len(name) > nameW {
			name = name[:nameW-2] + ".."
		}

		comp := job.Company
		if len(comp) > compW {
			comp = comp[:compW-2] + ".."
		}

		row := fmt.Sprintf(" %-*d %-*s %-*s %-*s ",
			numW, i+1,
			numberW, job.Number,
			nameW, name,
			compW, comp,
		)

		if i == m.jobCursor {
			b.WriteString(selectedStyle.Render(row))
		} else {
			b.WriteString(row)
		}

		if i < end-1 {
			b.WriteString("\n")
		}

	}

	return b.String()

}

// ─── Panel 3: Job Items ────────────────────────────────────

func (m model) jobItemsIcon() string {

	if m.itemLoading {
		return "⟳"
	}
	if m.itemError != "" {
		return "error"
	}
	if len(m.itemList) > 0 {
		return fmt.Sprintf("%d items", len(m.itemList))
	}
	return ""

}

func (m model) viewJobItems(w, h int) string {

	if m.itemLoading {
		return dimStyle.Render(padLine(w, "Loading job items..."))
	}

	if m.itemError != "" {
		msg := "Error: " + m.itemError
		return dimStyle.Render(padLine(w, msg))
	}

	if len(m.itemList) == 0 {
		if m.selectedJob.ID > 0 {
			return dimStyle.Render(padLine(w, "No items available"))
		}
		return dimStyle.Render(padLine(w, "Select a job first"))
	}

	var b strings.Builder

	jobInfo := fmt.Sprintf("[%s] %s", m.selectedJob.Number, m.selectedJob.Name)
	b.WriteString(lipgloss.NewStyle().Foreground(colorBlue2).Render(padLine(w, jobInfo)))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.Repeat("─", w)))
	b.WriteString("\n")
	h -= 2

	tableW := w - 2
	if tableW < 1 {
		tableW = 1
	}

	numW := 3
	timeW := 14
	separators := 2
	nameW := tableW - numW - timeW - separators
	if nameW < 5 {
		nameW = 5
	}

	start, end := visibleRange(m.itemCursor, len(m.itemList), h)

	for i := start; i < end; i++ {

		item := m.itemList[i]

		name := item.Name
		if len(name) > nameW {
			name = name[:nameW-2] + ".."
		}

		timeStr := fmt.Sprintf("%s/%s",
			span.DurationAsString(item.LoggedMinutes),
			span.DurationAsString(item.PlannedMinutes),
		)

		row := fmt.Sprintf(" %-*d %-*s %-*s ",
			numW, i+1,
			nameW, name,
			timeW, timeStr,
		)

		if i == m.itemCursor {
			b.WriteString(selectedStyle.Render(row))
		} else {
			b.WriteString(row)
		}

		if i < end-1 {
			b.WriteString("\n")
		}

	}

	return b.String()

}

// ─── Helpers ────────────────────────────────────────────────

func (m model) dayTabs() string {

	if len(m.data.Days) <= 1 {
		return ""
	}

	var tabs strings.Builder
	for i, day := range m.data.Days {
		name := day.Weekday.String()[:3]
		if i == m.day {
			tabs.WriteString(lipgloss.NewStyle().Foreground(colorBlue1).Bold(true).Render(name))
		} else {
			tabs.WriteString(dimStyle.Render(name))
		}
		if i < len(m.data.Days)-1 {
			tabs.WriteString(dimStyle.Render(" · "))
		}
	}

	left := dimStyle.Render("◀")
	right := dimStyle.Render("▶")

	return left + " " + tabs.String() + " " + right

}

func renderHeaderLine(title, right string, width int) string {

	if right == "" {
		return title
	}

	gap := width - lipgloss.Width(title) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}

	return title + strings.Repeat(" ", gap) + right

}

func padLine(width int, text string) string {

	if width <= 2 {
		return text
	}

	max := width - 2
	trimmed := ansi.Truncate(text, max, "")
	return " " + trimmed + strings.Repeat(" ", max-lipgloss.Width(trimmed)) + " "

}

func formatDayDate(date time.Time) string {

	day := date.Day()
	month := date.Format("Jan")
	suffix := ordinalSuffix(day)

	return fmt.Sprintf("%d%s %s", day, suffix, month)

}

func ordinalSuffix(day int) string {

	if day%100 >= 11 && day%100 <= 13 {
		return "th"
	}

	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}

}

func visibleRange(cursor, total, viewHeight int) (int, int) {

	if total <= viewHeight || viewHeight <= 0 {
		return 0, total
	}

	start := cursor - viewHeight/2
	if start < 0 {
		start = 0
	}

	end := start + viewHeight
	if end > total {
		end = total
		start = end - viewHeight
		if start < 0 {
			start = 0
		}
	}

	return start, end

}
