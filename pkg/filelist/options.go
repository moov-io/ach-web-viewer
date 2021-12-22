package filelist

import (
	"fmt"
	"net/http"
	"time"
)

type ListOpts struct {
	StartDate, EndDate time.Time
}

func (opts ListOpts) Inside(when time.Time) bool {
	afterStart := opts.StartDate.Before(when)
	beforeEnd := when.Before(opts.EndDate)
	return afterStart && beforeEnd
}

func (opts ListOpts) Dates() []string {
	if opts.StartDate.IsZero() || opts.EndDate.IsZero() {
		return nil
	}

	forward := 1
	out := []string{
		opts.StartDate.Format("2006-01-02"),
	}
	for {
		when := opts.StartDate.Add(time.Duration(forward) * 24 * time.Hour)
		if !opts.Inside(when) {
			break
		}
		forward++
		out = append(out, when.Format("2006-01-02"))
	}
	return out
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

	return opts, nil
}
