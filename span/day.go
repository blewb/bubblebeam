package span

import (
	"cmp"
	"fmt"
	"slices"
	"time"
)

type Day struct {
	Entries  []Entry      `json:"entries"`
	Weekday  time.Weekday `json:"weekday"`
	Duration int          `json:"duration"`
}

func newDay(w int) Day {
	return Day{
		Entries: make([]Entry, 0),
		Weekday: time.Weekday(w),
	}
}

const (
	HOURS_IN_DAY  = 24
	SLOTS_IN_HOUR = 4
	SLOTS_IN_DAY  = HOURS_IN_DAY * SLOTS_IN_HOUR
)

func (d *Day) Validate() {

	slices.SortFunc(d.Entries, func(a, b Entry) int {
		return cmp.Compare(a.End.GetSlot(), b.End.GetSlot())
	})

	starts := make([][]int, SLOTS_IN_DAY)
	ends := make([][]int, SLOTS_IN_DAY)

	var (
		entry                     Entry
		startSlot, endSlot, index int
	)

	for index, entry = range d.Entries {
		startSlot, endSlot = entry.Start.GetSlot(), entry.End.GetSlot()
		starts[startSlot] = append(starts[startSlot], index)
		ends[endSlot] = append(ends[endSlot], index)
		d.Duration += entry.Duration
	}

	current := make(map[int]struct{})
	conflicts := make(map[[2]int]struct{})

	for slot := 0; slot < SLOTS_IN_DAY; slot++ {

		if len(ends[slot]) > 0 {
			for _, index = range ends[slot] {
				delete(current, index)
			}
		}

		if len(starts[slot]) > 0 {
			for _, index = range starts[slot] {
				current[index] = struct{}{}
			}
		}

		registerConflicts(current, &conflicts)

	}

	if len(conflicts) > 0 {
		d.writeConflicts(conflicts)
	}

}

func (d *Day) writeConflicts(conflicts map[[2]int]struct{}) {

	var a, b int

	for conflict := range conflicts {

		a, b = conflict[0], conflict[1]

		d.Entries[a].Warnings = append(d.Entries[a].Warnings, fmt.Sprintf("conflicts with %s", d.Entries[b].RefName()))
		d.Entries[b].Warnings = append(d.Entries[b].Warnings, fmt.Sprintf("conflicts with %s", d.Entries[a].RefName()))

	}

}

func registerConflicts(current map[int]struct{}, conflicts *map[[2]int]struct{}) {

	n := len(current)
	if n < 2 {
		return
	}

	indexes := make([]int, 0, n)

	var (
		cf   [2]int
		i, j int
	)

	for i = range current {
		indexes = append(indexes, i)
	}

	for i = 0; i < n-1; i++ {
		for j = i + 1; j < n; j++ {
			cf[0], cf[1] = indexes[i], indexes[j]
			(*conflicts)[cf] = struct{}{}
		}
	}

}
