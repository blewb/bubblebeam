package span

import (
	"errors"
	"strconv"
	"strings"
)

const (
	TIMESTAMP_SIZE = 5
	RANGE_SIZE     = 13
)

func parseTimestamp(input string, warnings *[]string) Timestamp {

	if len(input) != TIMESTAMP_SIZE {
		if warnings != nil {
			*warnings = append(*warnings, "timestamp must be in the format hh:mm")
		}
		return Timestamp{}
	}

	parts := strings.Split(input, ":")

	if len(parts) != 2 || len(parts[0]) != 2 || len(parts[1]) != 2 {
		if warnings != nil {
			*warnings = append(*warnings, "timestamp must be in the format hh:mm")
		}
		return Timestamp{}
	}

	ts := Timestamp{}
	var err error

	// Parse

	ts.Hour, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		if warnings != nil {
			*warnings = append(*warnings, "unable to parse hours")
		}
		return ts
	}

	ts.Minute, err = strconv.Atoi(parts[1])
	if err != nil {
		if warnings != nil {
			*warnings = append(*warnings, "unable to parse minutes")
		}
		return ts
	}

	// Validate

	if ts.Hour < 0 || ts.Hour > 23 {
		if warnings != nil {
			*warnings = append(*warnings, "hours must be between 0 and 23")
		}
		return ts
	}

	if ts.Minute < 0 || ts.Minute > 45 || ts.Minute%15 != 0 {
		if warnings != nil {
			*warnings = append(*warnings, "minutes must be 0, 15, 30 or 45")
		}
		return ts
	}

	return ts

}

func parseRange(input string, warnings *[]string) [2]Timestamp {

	rng := [2]Timestamp{}

	if len(input) != RANGE_SIZE {
		if warnings != nil {
			*warnings = append(*warnings, "range must be in the format hh:mm - hh:mm")
		}
		return rng
	}

	rng[0] = parseTimestamp(input[0:TIMESTAMP_SIZE], warnings)
	rng[1] = parseTimestamp(input[RANGE_SIZE-TIMESTAMP_SIZE:], warnings)

	return rng

}

func parseDescriptionAndTag(input string) (string, string) {

	if len(input) == 0 {
		return "", ""
	}

	h := strings.LastIndex(input, "#")

	if h < 0 {
		return strings.TrimSpace(input), ""
	}

	desc := strings.TrimSpace(input[0:h])
	tag := strings.TrimSpace(input[h+1:])

	return desc, tag

}

func parseEntry(input string) (Entry, error) {

	if len(input) < RANGE_SIZE {
		return Entry{}, errors.New("unable to parse entry, input too short")
	}

	warnings := make([]string, 0)

	ranges := parseRange(input[0:RANGE_SIZE], &warnings)
	desc, tag := parseDescriptionAndTag(input[RANGE_SIZE:])

	return newEntry(desc, tag, ranges[0], ranges[1], warnings), nil

}
