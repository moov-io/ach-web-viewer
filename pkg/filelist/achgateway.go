package filelist

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/moov-io/ach"
	"github.com/moov-io/ach-web-viewer/pkg/service"
)

type achgatewayLister struct {
	endpoint string
	client   *http.Client
	sourceID string
	shards   []string
}

func newACHGatewayLister(sourceID string, cfg service.ACHGatewayConfig) (*achgatewayLister, error) {
	if cfg.Endpoint == "" {
		return nil, errors.New("missing endpoint")
	}

	timeout := 10 * time.Second
	if cfg.Timeout > 0*time.Second {
		timeout = cfg.Timeout
	}

	return &achgatewayLister{
		endpoint: cfg.Endpoint,
		client: &http.Client{
			Timeout: timeout,
		},
		sourceID: sourceID,
		shards:   cfg.Shards,
	}, nil
}

func (a *achgatewayLister) SourceID() string {
	return a.sourceID
}

func (a *achgatewayLister) GetFiles(opts ListOpts) (Files, error) {
	out := Files{
		SourceID:   a.sourceID,
		SourceType: "ACHGateway",
	}
	for i := range a.shards {
		files, err := a.getFiles(a.shards[i])
		if err != nil {
			return out, fmt.Errorf("error getting %s files: %w", a.shards[i], err)
		}
		out.Files = append(out.Files, files...)
	}
	return out, nil
}

func (a *achgatewayLister) getFiles(shard string) ([]File, error) {
	req, err := http.NewRequest("GET", a.endpoint+"/shards/"+shard+"/files", nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Files []struct {
			Filename string    `json:"filename"`
			Path     string    `json:"path"`
			ModTime  time.Time `json:"modTime"`
		} `json:"files"`
	}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, err
	}

	var out []File
	for i := range wrapper.Files {
		_, fname := filepath.Split(wrapper.Files[i].Filename)

		out = append(out, File{
			Name:        fname,
			StoragePath: fmt.Sprintf("shards/%s/files/", shard),
			CreatedAt:   wrapper.Files[i].ModTime,
		})
	}
	return out, nil
}

func (a *achgatewayLister) GetFile(path string) (*File, error) {
	req, err := http.NewRequest("GET", a.endpoint+"/"+path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var wrapper struct {
		Filename       string    `json:"filename"`
		ContentsBase64 string    `json:"contentsBase64"`
		Valid          error     `json:"valid"`
		ModTime        time.Time `json:"modTime"`
	}
	err = json.NewDecoder(resp.Body).Decode(&wrapper)
	if err != nil {
		return nil, err
	}
	dir, _ := filepath.Split(path)

	contents, err := base64.StdEncoding.DecodeString(wrapper.ContentsBase64)
	if err != nil {
		return nil, err
	}

	file, err := ach.NewReader(bytes.NewReader(contents)).Read()
	if err != nil {
		return nil, err
	}

	return &File{
		Name:        wrapper.Filename,
		StoragePath: dir,
		Contents:    &file,
		CreatedAt:   wrapper.ModTime,
	}, nil
}
