// generated-from:978eb7e8497019d58e3ef1a92840745f9415cc1bceb815251e7a716fdeb0d674 DO NOT REMOVE, DO UPDATE

package main

import (
	"fmt"
	"os"

	achwebviewer "github.com/moov-io/ach-web-viewer"
	"github.com/moov-io/ach-web-viewer/pkg/api"
	"github.com/moov-io/ach-web-viewer/pkg/filelist"
	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/ach-web-viewer/pkg/web"
	"github.com/moov-io/base/log"
)

func main() {
	env := &service.Environment{
		Logger: log.NewDefaultLogger().Set("app", log.String("ach-web-viewer")).Set("version", log.String(achwebviewer.Version)),
	}

	env, err := service.NewEnvironment(env)
	if err != nil {
		env.Logger.Fatal().LogErrorf("Error loading up environment: %v", err)
		os.Exit(1)
	}
	defer env.Shutdown()

	termListener := service.NewTerminationListener()

	// Register API routes with listers
	listers := createFileListers(env.Config.Sources)
	api.AppendRoutes(env.Logger, env.PublicRouter, listers)
	web.AppendRoutes(env, listers)

	stopServers := env.RunServers(termListener)
	defer stopServers()

	service.AwaitTermination(env.Logger, termListener)
}

func createFileListers(cfg service.Sources) filelist.Listers {
	listers, err := filelist.NewListers(cfg)
	if err != nil {
		panic(fmt.Sprintf("ERROR initializing listers: %v", err))
	}
	return listers
}
