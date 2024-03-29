package yyyymmdd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPrefixes(t *testing.T) {
	start := time.Date(2023, time.December, 23, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	expected := []string{"2023-12-2", "2023-12-3"}
	require.ElementsMatch(t, expected, Prefixes(start, end))

	end = time.Date(2024, time.January, 4, 0, 0, 0, 0, time.UTC)
	expected = append(expected, "2024-01-0")
	require.ElementsMatch(t, expected, Prefixes(start, end))

	end = time.Date(2024, time.January, 10, 0, 0, 0, 0, time.UTC)
	expected = append(expected, "2024-01-10")
	require.ElementsMatch(t, expected, Prefixes(start, end))

	end = time.Date(2024, time.January, 11, 0, 0, 0, 0, time.UTC)
	expected = []string{"2023-12-2", "2023-12-3", "2024-01-0", "2024-01-1"}
	require.ElementsMatch(t, expected, Prefixes(start, end))

	end = time.Date(2024, time.January, 20, 0, 0, 0, 0, time.UTC)
	expected = []string{"2023-12-2", "2023-12-3", "2024-01-0", "2024-01-1", "2024-01-20"}
	require.ElementsMatch(t, expected, Prefixes(start, end))

	end = time.Date(2024, time.January, 25, 0, 0, 0, 0, time.UTC)
	expected = []string{"2023-12-2", "2023-12-3", "2024-01-0", "2024-01-1", "2024-01-2"}
	require.ElementsMatch(t, expected, Prefixes(start, end))

	end = time.Date(2024, time.January, 30, 0, 0, 0, 0, time.UTC)
	expected = []string{"2023-12-2", "2023-12-3", "2024-01-0", "2024-01-1", "2024-01-2", "2024-01-30"}
	require.ElementsMatch(t, expected, Prefixes(start, end))
}
