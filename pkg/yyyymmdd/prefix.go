package yyyymmdd

import (
	"slices"
	"time"
)

// Given two time.Time values generate a start and end prefix (in yyyy-mm-dd format)
// which serves as prefixs usable to filter.
//
// Examples:
//
//	2023-12-23 to 2023-12-31 produces 2023-12-2, 2023-12-3
//	2023-12-23 to 2024-01-10 produces 2023-12-2, 2023-12-3, 2024-01-0, 2024-01-10
func Prefixes(start, end time.Time) []string {
	var out []string

	// For now just iterate over each day and chop off the trailing day digit
	for {
		if start.After(end) {
			break
		}

		// Add the current day to our list
		ts := start.Format("2006-01-02")

		// Only when the end day is 10, 20, 30 we can extend the timestamp
		if (start.Month() == end.Month()) && (start.Day() == end.Day()) && end.Day()%10 == 0 {
			// do nothing
		} else {
			ts = ts[:len(ts)-1] // chop off the last digit
		}

		out = append(out, ts)

		start = start.Add(24 * time.Hour)
	}

	slices.Sort(out)
	return slices.Compact(out)
}
