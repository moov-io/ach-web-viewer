package api

import (
	"encoding/json"
	"net/http"

	"github.com/moov-io/ach-web-viewer/pkg/filelist"
	"github.com/moov-io/base/log"

	"github.com/gorilla/mux"
)

func AppendRoutes(logger log.Logger, router *mux.Router, listers filelist.Listers) {
	router.Methods("GET").Path("/files").HandlerFunc(listFiles(logger, listers))
}

type wrapper struct {
	Sources map[string]filelist.Files
}

func listFiles(logger log.Logger, listers filelist.Listers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opts, err := filelist.ReadListOptions(r)
		if err != nil {
			logger.Set("service", log.String("api")).Error().LogErrorf("problem reading list params: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		resp, err := listers.GetFiles(opts)
		if err != nil {
			logger.Set("service", log.String("api")).Error().LogErrorf("problem listing files: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(wrapper{
			Sources: resp,
		})
	}
}
