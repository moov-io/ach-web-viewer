package service

import (
	achwebviewer "github.com/moov-io/ach-web-viewer"
	"github.com/moov-io/base/config"
	"github.com/moov-io/base/log"
)

func LoadConfig(logger log.Logger) (*Config, error) {
	configService := config.NewService(logger)

	global := &GlobalConfig{}
	if err := configService.LoadFromFS(global, achwebviewer.ConfigDefaults); err != nil {
		return nil, err
	}

	cfg := &global.ACHWebViewer

	return cfg, nil
}
