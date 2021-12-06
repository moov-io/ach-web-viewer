package filelist

import (
	"fmt"
	"net/http"
	"time"
)

type ListOpts struct {
	StartDate, EndDate time.Time
}

func ReadListOptions(r *http.Request) (ListOpts, error) {
	opts := ListOpts{
		StartDate: time.Now().Add(-7 * 24 * time.Hour),
		EndDate:   time.Now().Add(1 * time.Hour),
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
