package filelist

import (
	"fmt"
	"net/http"
	"time"
)

type ListOpts struct {
	StartDate, EndDate time.Time
	Pattern            string
}

func (opts ListOpts) Inside(when time.Time) bool {
	afterStart := opts.StartDate.Before(when)
	beforeEnd := when.Before(opts.EndDate)
	return afterStart && beforeEnd
}

func DefaultListOptions(when time.Time) ListOpts {
	return ListOpts{
		StartDate: when.Add(-7 * 24 * time.Hour),
		EndDate:   when.Add(1 * time.Hour),
	}
}

func ReadListOptions(r *http.Request) (ListOpts, error) {
	opts := DefaultListOptions(time.Now())

	if r == nil || r.URL == nil {
		return opts, nil
	}

	qry := r.URL.Query()

	if v := qry.Get("startDate"); v != "" {
		if tt, err := time.Parse("2006-01-02", v); err != nil {
			return opts, fmt.Errorf("unable to read %s: %w", v, err)
		} else {
			opts.StartDate = tt
		}
	}
	if v := qry.Get("endDate"); v != "" {
		if tt, err := time.Parse("2006-01-02", v); err != nil {
			return opts, fmt.Errorf("unable to read %s: %w", v, err)
		} else {
			opts.EndDate = tt
		}
	}

	opts.Pattern = qry.Get("pattern")

	return opts, nil
}
