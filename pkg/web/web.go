package web

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	webdisplay "github.com/moov-io/ach-web-viewer/pkg/display/web"
	"github.com/moov-io/ach-web-viewer/pkg/filelist"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/base/log"

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
	Path     string
	Filename string
}

type listFilesSource struct {
	ID    string
	Type  string
	Files []listFile
}

type listFilesTemplate struct {
	BaseURL string
	Sources []listFilesSource
}

var listFilesTmpl = initTemplate("list-files", "/webui/index.html.tpl")

func listFiles(logger log.Logger, listers filelist.Listers, basePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := listers.GetFiles()
		if err != nil {
			logger.Set("service", log.String("web")).Error().LogErrorf("problem listing files: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html")

		response := listFilesTemplate{
			BaseURL: baseURL(basePath),
		}
		for _, files := range resp {
			var listings []listFile
			for i := range files.Files {
				fullName := fmt.Sprintf("%s%s", files.Files[i].StoragePath, files.Files[i].Name)

				listings = append(listings, listFile{
					Path:     path.Join(basePath, "sources", files.SourceID, fullName),
					Filename: files.Files[i].Name,
				})
			}
			response.Sources = append(response.Sources, listFilesSource{
				ID:    files.SourceID,
				Type:  files.SourceType,
				Files: listings,
			})
		}
		err = listFilesTmpl.Execute(w, response)
		if err != nil {
			fmt.Printf("ERROR: rendering template: %v\n", err)
		}
	}
}

type getFileTemplate struct {
	Filename string
	BaseURL  string
	Contents string
	Valid    error
}

var getFileTmpl = initTemplate("get-file", "/webui/file.html.tpl")

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
		webdisplay.File(&contents, &cfg, file)

		err = getFileTmpl.Execute(w, getFileTemplate{
			Filename: filepath.Base(fullPath),
			BaseURL:  baseURL(basePath),
			Contents: contents.String(),
			Valid:    file.Validate(),
		})
		if err != nil {
			fmt.Printf("ERROR: rendering template: %v\n", err)
		}

		// TODO(adam): include err := file.Validate() response
	}
}

func baseURL(basePath string) string {
	cleaned := path.Clean(basePath)
	if cleaned == "." {
		return "/"
	}
	return cleaned + "/"
}

func initTemplate(name, path string) *template.Template {
	fd, err := pkger.Open(path)
	if err != nil {
		panic(fmt.Sprintf("error opening %s: %v", path, err))
	}
	defer fd.Close()

	bs, err := ioutil.ReadAll(fd)
	if err != nil {
		panic(fmt.Sprintf("error reading %s: %v", fd.Name(), err))
	}

	return template.Must(template.New(name).Parse(string(bs)))
}
