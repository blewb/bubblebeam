package span

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Span struct {
	Days []Day `json:"days"`
}

var weekdays = map[string]int{
	`sunday`:    0,
	`monday`:    1,
	`tuesday`:   2,
	`wednesday`: 3,
	`thursday`:  4,
	`friday`:    5,
	`saturday`:  6,
}

func (sp *Span) Read(useSample int) error {

	filepath := getFilePath(useSample)

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var (
		day     Day
		entry   Entry
		line    string
		nextday int
	)

	currday := -1

	for scanner.Scan() {

		line = scanner.Text()
		if len(line) < 3 || line[0:3] == "---" {
			continue
		}

		nextday = getDay(line)

		if nextday >= 0 && nextday != currday {
			if len(day.Entries) > 0 {
				sp.Days = append(sp.Days, day)
			}
			currday = nextday
			day = newDay(currday)
			continue
		}

		if currday >= 0 {
			entry, err = parseEntry(line)
			if err != nil {
				fmt.Println(err)
				continue
			}
			day.Entries = append(day.Entries, entry)
		}

	}

	if len(day.Entries) > 0 {
		day.Validate()
		sp.Days = append(sp.Days, day)
	}

	return nil

}

func (sp Span) Render() {

	if len(sp.Days) == 0 {
		fmt.Println("")
		fmt.Println("<empty span>")
		fmt.Println("")
		return
	}

	var wd string

	for _, day := range sp.Days {

		wd = day.Weekday.String()

		fmt.Println("")

		fmt.Println(strings.Repeat("#", len(wd)+8))
		fmt.Printf("##  %s  ##\n", wd)
		fmt.Println(strings.Repeat("#", len(wd)+8))

		fmt.Println("")

		for _, entry := range day.Entries {
			entry.Render()
		}
		fmt.Println("")
	}

}

func getDay(line string) int {

	// Shortest/longest entry in weekdays map
	if len(line) < 6 || len(line) > 9 {
		return -1
	}

	line = strings.ToLower(line)

	if _, ok := weekdays[line]; !ok {
		return -1
	}

	return weekdays[line]

}

func getFilePath(useSample int) string {

	if useSample > 0 {
		return fmt.Sprintf("sample/%d.txt", useSample)
	}

	// @TODO: read proper file from somewhere? .env?

	return ""

}
