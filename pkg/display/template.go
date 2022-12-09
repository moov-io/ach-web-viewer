package display

import (
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/markbates/pkger"
)

var templateFuncs template.FuncMap = map[string]interface{}{
	"dateTime": func(when string, pattern string) string {
		tt, _ := time.Parse("2006-01-02", when)
		return tt.Format(pattern)
	},
	"startDateParam": func(end time.Time) string {
		start := end.Add(-7 * 24 * time.Hour)
		return fmt.Sprintf("?startDate=%s&endDate=%s", start.Format("2006-01-02"), end.Format("2006-01-02"))
	},
	"endDateParam": func(start time.Time) string {
		end := start.Add(7 * 24 * time.Hour)
		return fmt.Sprintf("?startDate=%s&endDate=%s", start.Format("2006-01-02"), end.Format("2006-01-02"))
	},
}

func InitTemplate(name, path string) *template.Template {
	fd, err := pkger.Open(path)
	if err != nil {
		panic(fmt.Sprintf("error opening %s: %v", path, err))
	}
	defer fd.Close()

	bs, err := io.ReadAll(fd)
	if err != nil {
		panic(fmt.Sprintf("error reading %s: %v", fd.Name(), err))
	}

	return template.Must(template.New(name).Funcs(templateFuncs).Parse(string(bs)))
}
