// generated-from:f0d07dce83f211366844ebe556c6f4909087f76e808087e32719a9e510feec0d DO NOT REMOVE, DO UPDATE

package test

import (
	"testing"

	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/base/log"
	"github.com/moov-io/base/stime"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

type TestEnvironment struct {
	Assert     *require.Assertions
	StaticTime stime.StaticTimeService

	service.Environment
}

func NewEnvironment(t *testing.T, router *mux.Router) *TestEnvironment {
	testEnv := &TestEnvironment{}

	testEnv.PublicRouter = router
	testEnv.Assert = require.New(t)
	testEnv.Logger = log.NewDefaultLogger()
	testEnv.StaticTime = stime.NewStaticTimeService()
	testEnv.TimeService = testEnv.StaticTime

	cfg, err := service.LoadConfig(testEnv.Logger)
	if err != nil {
		t.Fatal(err)
	}
	testEnv.Config = cfg

	_, err = service.NewEnvironment(&testEnv.Environment)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(testEnv.Shutdown)

	return testEnv
}
