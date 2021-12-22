package filelist

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListOpts__Inside(t *testing.T) {
	req := &http.Request{}
	opts, err := ReadListOptions(req)
	require.NoError(t, err)

	require.True(t, opts.Inside(time.Now()))

	future := time.Now().Add(7 * 24 * time.Hour)
	require.False(t, opts.Inside(future))

	past := time.Now().Add(-30 * 7 * 24 * time.Hour)
	require.False(t, opts.Inside(past))
}

func TestListOpts__Dates(t *testing.T) {
	eastern, _ := time.LoadLocation("America/New_York")

	opts := ListOpts{
		StartDate: time.Date(2021, time.December, 21, 10, 30, 0, 0, eastern),
		EndDate:   time.Date(2021, time.December, 27, 10, 30, 0, 0, eastern),
	}

	dates := opts.Dates()
	require.Len(t, dates, 6)
	require.Equal(t, "2021-12-21", dates[0])
	require.Equal(t, "2021-12-22", dates[1])
	require.Equal(t, "2021-12-23", dates[2])
	require.Equal(t, "2021-12-24", dates[3])
	require.Equal(t, "2021-12-25", dates[4])
	require.Equal(t, "2021-12-26", dates[5])
}
