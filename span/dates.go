package span

import (
	"time"
)

type Datestamp struct {
	Formal   string // Date as YYYY-MM-DD for technical/API use
	Friendly string // Date as Month Day, for reading
	Weekday  time.Weekday
}

// Stamp Range and Stamp Date (Index)
func GetDatestamps(day time.Time, days int) (sr []Datestamp, sd int) {

	sr = make([]Datestamp, days)
	today := day.Format(time.DateOnly)

	wday := day.Weekday()
	if wday != time.Monday {
		offset := time.Duration(wday) - 1
		if offset == 0 {
			offset = 6
		}
		day = day.Add(time.Hour * -24 * offset)
	}
	day = day.Add(time.Hour * 24 * -7)

	for i := range sr {
		sr[i] = Datestamp{
			Formal:   day.Format(time.DateOnly),
			Friendly: day.Format("Jan 2"),
			Weekday:  day.Weekday(),
		}
		if sr[i].Formal == today {
			sd = i
		}
		day = day.Add(time.Hour * 24)
	}

	return

}
