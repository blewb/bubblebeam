package span

import (
	"reflect"
	"testing"
)

func TestParseTimestamp(t *testing.T) {

	tests := map[string]struct {
		expected   Timestamp
		input, err string
	}{
		`empty`: {
			expected: Timestamp{},
			input:    "",
			err:      "timestamp must be in the format hh:mm",
		},
		`simple`: {
			expected: Timestamp{10, 0},
			input:    "10:00",
			err:      "",
		},
		`spaced`: {
			expected: Timestamp{9, 15},
			input:    " 9:15",
			err:      "",
		},
		`halfa`: {
			expected: Timestamp{14, 30},
			input:    "14:30",
			err:      "",
		},
		`hightide`: {
			expected: Timestamp{},
			input:    "24:30",
			err:      "hours must be between 0 and 23",
		},
		`negative`: {
			expected: Timestamp{},
			input:    "-9:15",
			err:      "hours must be between 0 and 23",
		},
		`negativeMinutes`: {
			expected: Timestamp{},
			input:    "12:-5",
			err:      "minutes must be 0, 15, 30 or 45",
		},
		`wrong`: {
			expected: Timestamp{},
			input:    "not a timestamp",
			err:      "timestamp must be in the format hh:mm",
		},
		`malf`: {
			expected: Timestamp{},
			input:    "3X:12",
			err:      "unable to parse hours",
		},
		`malfMinutes`: {
			expected: Timestamp{},
			input:    "17:I2",
			err:      "unable to parse minutes",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			warnings := make([]string, 0)
			got := parseTimestamp(tt.input, &warnings)
			if len(warnings) > 0 {
				if warnings[len(warnings)-1] != tt.err {
					t.Errorf("got: %s, expected: %s", warnings[len(warnings)-1], tt.err)
				}
			} else {
				if got != tt.expected {
					t.Errorf("got: %v, expected: %v", got, tt.expected)
				}
			}
		})
	}

}

func TestParseRange(t *testing.T) {

	tests := map[string]struct {
		expected   [2]Timestamp
		input, err string
	}{
		`empty`: {
			expected: [2]Timestamp{},
			input:    "",
			err:      "range must be in the format hh:mm - hh:mm",
		},
		`simple`: {
			expected: [2]Timestamp{{10, 0}, {11, 0}},
			input:    "10:00 - 11:00",
			err:      "",
		},
		`spaced`: {
			expected: [2]Timestamp{{9, 15}, {11, 30}},
			input:    " 9:15 - 11:30",
			err:      "",
		},
		`mixed`: {
			expected: [2]Timestamp{{14, 45}, {0, 0}},
			input:    "14:45 - xx:xx",
			err:      "unable to parse hours",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			warnings := make([]string, 0)
			got := parseRange(tt.input, &warnings)
			if len(warnings) > 0 {
				if warnings[len(warnings)-1] != tt.err {
					t.Errorf("got: %s, expected: %s", warnings[len(warnings)-1], tt.err)
				}
			} else {
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("got: %+v, expected: %+v", got, tt.expected)
				}
			}
		})
	}

}

func TestParseDescriptionAndTag(t *testing.T) {

	tests := map[string]struct {
		input, desc, tag string
	}{
		`empty`: {
			input: "",
			desc:  "",
			tag:   "",
		},
		`notag`: {
			input: "some line with no tag ",
			desc:  "some line with no tag",
			tag:   "",
		},
		`tagged`: {
			input: "did a bit of work #somejob",
			desc:  "did a bit of work",
			tag:   "somejob",
		},
		`multitag`: {
			input: " did a bit of #work #some job ",
			desc:  "did a bit of #work",
			tag:   "some job",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			desc, tag := parseDescriptionAndTag(tt.input)
			if desc != tt.desc {
				t.Errorf("got: %s, expected: %s", desc, tt.desc)
			} else if tag != tt.tag {
				t.Errorf("got: %s, expected: %s", tag, tt.tag)
			}
		})
	}

}

func TestParseEntry(t *testing.T) {

	tests := map[string]struct {
		entry      Entry
		input, err string
	}{
		`empty`: {
			entry: Entry{},
			input: "",
			err:   "unable to parse entry, input too short",
		},
		`simple`: {
			entry: Entry{
				Warnings:    []string{},
				Description: "did some work",
				Tag:         "job",
				Start:       Timestamp{9, 0},
				End:         Timestamp{10, 0},
				Duration:    60,
			},
			input: " 9:00 - 10:00 did some work #job",
			err:   "",
		},
		`notag`: {
			entry: Entry{
				Warnings:    []string{},
				Description: "did lots of work",
				Tag:         "",
				Start:       Timestamp{10, 15},
				End:         Timestamp{12, 0},
				Duration:    105,
			},
			input: "10:15 - 12:00 did lots of work ",
			err:   "",
		},
		`timely`: {
			entry: Entry{
				Warnings: []string{
					"unable to parse hours",
				},
				Description: "end time with no start #job",
				Tag:         "job",
				Start:       Timestamp{0, 0},
				End:         Timestamp{2, 30},
				Duration:    150,
			},
			input: "xx:xx -  2:30 end time with no start #job #job",
			err:   "",
		},
		`untimely`: {
			entry: Entry{
				Warnings: []string{
					"negative duration",
				},
				Description: "desc",
				Tag:         "",
				Start:       Timestamp{15, 0},
				End:         Timestamp{14, 0},
				Duration:    0,
			},
			input: "15:00 - 14:00 desc",
			err:   "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			entry, err := parseEntry(tt.input)
			if err != nil && err.Error() != tt.err {
				t.Errorf("got: %s, expected: %s", err.Error(), tt.err)
			} else if !reflect.DeepEqual(entry, tt.entry) {
				t.Errorf("got: %+v, expected: %+v", entry, tt.entry)
			}
		})
	}

}
