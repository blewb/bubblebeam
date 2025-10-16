package span

import (
	"fmt"
	"strings"
)

type Timestamp struct {
	Hour, Minute int
}

type Entry struct {
	Warnings    []string  `json:"warnings"`
	Description string    `json:"description"`
	Tag         string    `json:"tag"`
	Start       Timestamp `json:"start"`
	End         Timestamp `json:"end"`
	Duration    int       `json:"duration"`
}

const (
	REF_LENGTH = 32
)

func newEntry(desc, tag string, start, end Timestamp, warnings []string) Entry {

	entry := Entry{
		Warnings:    warnings,
		Description: desc,
		Tag:         tag,
		Start:       start,
		End:         end,
	}

	entry.CalculateDuration()

	return entry

}

func (e *Entry) CalculateDuration() {

	dur := (e.End.Hour - e.Start.Hour) * 60
	if dur < 0 {
		e.Warnings = append(e.Warnings, "negative duration")
		return
	}

	dur += e.End.Minute - e.Start.Minute
	if dur < 0 {
		e.Warnings = append(e.Warnings, "negative duration")
		return
	}

	e.Duration = dur

}

func (e Entry) RefName() string {

	var b strings.Builder

	if len(e.Description) > REF_LENGTH {
		b.WriteString(e.Description[:REF_LENGTH] + "...")
	} else {
		b.WriteString(e.Description)
	}

	b.WriteString(" #")

	if len(e.Tag) > REF_LENGTH {
		b.WriteString(e.Tag[:REF_LENGTH] + "...")
	} else {
		b.WriteString(e.Tag)
	}

	return b.String()

}

func (t Timestamp) GetSlot() int {
	return t.Hour*4 + t.Minute/15
}

func (e Entry) Render() {

	fmt.Println("-------------")
	fmt.Printf("%s #%s\n", e.Description, e.Tag)
	fmt.Println("-------------")

	fmt.Printf("%02d:%02d - %02d:%02d\n", e.Start.Hour, e.Start.Minute, e.End.Hour, e.End.Minute)

	if len(e.Warnings) > 0 {
		for _, warn := range e.Warnings {
			fmt.Println("*", warn)
		}
	}

	fmt.Println("")

}
