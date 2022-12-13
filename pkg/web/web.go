package web

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/moov-io/ach-web-viewer/pkg/display"
	webdisplay "github.com/moov-io/ach-web-viewer/pkg/display/web"
	"github.com/moov-io/ach-web-viewer/pkg/filelist"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/base/log"
	"github.com/moov-io/base/strx"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
)

func AppendRoutes(env *service.Environment, listers filelist.Listers, basePath string) {
	env.PublicRouter.Methods("GET").Path("/").HandlerFunc(listFiles(env.Logger, listers, basePath))
	env.PublicRouter.Methods("GET").PathPrefix("/sources/{sourceID}/").HandlerFunc(getFile(env.Logger, env.Config.Display, listers, basePath))

	dir, _ := pkger.Open("/webui/")
	if dir != nil {
		env.PublicRouter.Methods("GET").Path("/style.css").Handler(http.StripPrefix(basePath, http.FileServer(dir)))
	}
}

type listFile struct {
	Path      string
	CleanPath string
	Filename  string
	CreatedAt time.Time
}

type listFileGroup struct {
	DateTime string
	Files    []listFile
}

type listFilesSource struct {
	ID     string
	Type   string
	Groups []listFileGroup
}

type listFilesOptions struct {
	TimeRangeMin, TimeRangeMax time.Time
}

type listFilesTemplate struct {
	BaseURL string
	Sources []listFilesSource
	Options listFilesOptions
}

var listFilesTmpl = display.InitTemplate("list-files", "/webui/index.html.tpl")

func listFiles(logger log.Logger, listers filelist.Listers, basePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := filelist.ReadListOptions(r)
		if err != nil {
			logger.Set("service", log.String("web")).Error().LogErrorf("problem reading list params: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		resp, err := listers.GetFiles(opts)
		if err != nil {
			logger.Set("service", log.String("web")).Error().LogErrorf("problem listing files: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html")

		response := listFilesTemplate{
			BaseURL: baseURL(basePath),
			Options: listFilesOptions{
				TimeRangeMin: opts.StartDate,
				TimeRangeMax: opts.EndDate,
			},
		}
		for _, files := range resp {
			var listings []listFile
			for i := range files.Files {
				fullName := fmt.Sprintf("%s%s", files.Files[i].StoragePath, files.Files[i].Name)

				fullPath := filepath.Join(basePath, "sources", files.SourceID, fullName)

				listings = append(listings, listFile{
					Path:      fullPath,
					CleanPath: cleanPath(basePath, fullPath),
					Filename:  files.Files[i].Name,
					CreatedAt: files.Files[i].CreatedAt,
				})
			}
			response.Sources = append(response.Sources, listFilesSource{
				ID:     files.SourceID,
				Type:   files.SourceType,
				Groups: groupFileListings(listings),
			})
		}
		err = listFilesTmpl.Execute(w, response)
		if err != nil {
			fmt.Printf("ERROR: rendering template: %v\n", err)
		}
	}
}

// cleanPath will strip fluff from paths so the web output is more clear where files come from.
func cleanPath(baseURL, path string) string {
	return filepath.Dir(strings.TrimPrefix(path, filepath.Join(baseURL, "sources"))) + "/"
}

func groupFileListings(listings []listFile) (out []listFileGroup) {
	format := func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	for i := range listings {
		found := false
		for k := range out {
			if out[k].DateTime == format(listings[i].CreatedAt) {
				found = true
				out[k].Files = append(out[k].Files, listings[i])
			}
		}
		if !found {
			out = append(out, listFileGroup{
				DateTime: format(listings[i].CreatedAt),
				Files:    []listFile{listings[i]},
			})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].DateTime > out[j].DateTime
	})
	return
}

type getFileTemplate struct {
	Filename     string
	BaseURL      string
	Contents     template.HTML
	Valid        error
	Metadata     map[string]string
	HelpfulLinks service.HelpfulLinks
}

var getFileTmpl = display.InitTemplate("get-file", "/webui/file.html.tpl")

func getFile(logger log.Logger, cfg service.DisplayConfig, listers filelist.Listers, basePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceID := mux.Vars(r)["sourceID"]
		fullPath := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("%s/sources/%s/", basePath, sourceID))

		file, err := listers.GetFile(sourceID, fullPath)
		if err != nil {
			logger.Warn().Logf("ERROR: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html")

		var contents bytes.Buffer
		if file != nil && file.Contents != nil {
			// Change the output formatting based on query params
			qq := r.URL.Query()
			params := webdisplay.DisplayParams{
				Format: strx.Or(qq.Get("format"), cfg.Format),
				SortBy: qq.Get("sort"),
				Order:  qq.Get("order"),
			}
			webdisplay.File(&contents, params, cfg.Masking, file.Contents)
		} else {
			contents.WriteString("file not found...")
		}

		validationError := errors.New("missing / partial file")
		if file != nil && file.Contents != nil {
			validationError = file.Contents.Validate()
		}

		err = getFileTmpl.Execute(w, getFileTemplate{
			Filename:     filepath.Base(fullPath),
			BaseURL:      baseURL(basePath),
			Contents:     template.HTML(contents.String()),
			Valid:        validationError,
			Metadata:     setMetadata(file),
			HelpfulLinks: cfg.HelpfulLinks,
		})
		if err != nil {
			fmt.Printf("ERROR: rendering template: %v\n", err)
		}
	}
}

func setMetadata(file *filelist.File) map[string]string {
	out := make(map[string]string)

	if file != nil {
		out["Created At"] = file.CreatedAt.Format(time.RFC1123)
	}

	return out
}

func baseURL(basePath string) string {
	cleaned := path.Clean(basePath)
	if cleaned == "." {
		return "/"
	}
	return cleaned + "/"
}
