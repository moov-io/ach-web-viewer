// generated-from:4581abea4a84e388c10c1c743717ada3cba817eafdb73bcdf332dd49ec28ecb1 DO NOT REMOVE, DO UPDATE

package service_test

import (
	"testing"

	"github.com/moov-io/base/log"
	"github.com/stretchr/testify/assert"

	"github.com/moov-io/ach-web-viewer/pkg/service"
)

func Test_Environment_Startup(t *testing.T) {
	a := assert.New(t)

	env := &service.Environment{
		Logger: log.NewDefaultLogger(),
	}

	env, err := service.NewEnvironment(env)
	a.Nil(err)

	t.Cleanup(env.Shutdown)
}
