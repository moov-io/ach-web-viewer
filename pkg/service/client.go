// generated-from:9e5795883635a80199c38497589e39b956ddfa90e17e6597be4077b792f076f1 DO NOT REMOVE, DO UPDATE

package service

import (
	"net/http"
	"time"

	"github.com/moov-io/base/log"
)

type ClientConfig struct {
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
}

func NewInternalClient(logger log.Logger, config *ClientConfig, name string) *http.Client {
	if config == nil {
		config = &ClientConfig{
			Timeout:             60 * time.Second,
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 20,
			MaxConnsPerHost:     20,
		}
	}

	// Default settings we approve of
	internalClient := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			ForceAttemptHTTP2:   true,
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			MaxConnsPerHost:     config.MaxConnsPerHost,
		},
	}

	return internalClient
}
