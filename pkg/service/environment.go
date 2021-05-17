// generated-from:b43099dcdb0795138b8654aca3676c36612c35678fe726cb5383616e99a62226 DO NOT REMOVE, DO UPDATE

package service

import (
	"net/http"

	_ "github.com/moov-io/ach-web-viewer"
	"github.com/moov-io/base/config"
	"github.com/moov-io/base/log"
	"github.com/moov-io/base/stime"

	"github.com/gorilla/mux"
)

// Environment - Contains everything thats been instantiated for this service.
type Environment struct {
	Logger         log.Logger
	Config         *Config
	TimeService    stime.TimeService
	InternalClient *http.Client

	PublicRouter *mux.Router
	Shutdown     func()
}

// NewEnvironment - Generates a new default environment. Overrides can be specified via configs.
func NewEnvironment(env *Environment) (*Environment, error) {
	if env == nil {
		env = &Environment{}
	}

	env.Shutdown = func() {}

	if env.Logger == nil {
		env.Logger = log.NewDefaultLogger()
	}

	if env.Config == nil {
		cfg, err := LoadConfig(env.Logger)
		if err != nil {
			return nil, err
		}

		env.Config = cfg
	}

	if env.InternalClient == nil {
		env.InternalClient = NewInternalClient(env.Logger, env.Config.Clients, "internal-client")
	}

	if env.TimeService == nil {
		env.TimeService = stime.NewSystemTimeService()
	}

	// router
	if env.PublicRouter == nil {
		env.PublicRouter = mux.NewRouter()
	}

	return env, nil
}

func LoadConfig(logger log.Logger) (*Config, error) {
	configService := config.NewService(logger)

	global := &GlobalConfig{}
	if err := configService.Load(global); err != nil {
		return nil, err
	}

	cfg := &global.ACHWebViewer

	return cfg, nil
}
