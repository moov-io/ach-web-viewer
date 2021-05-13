// generated-from:eeef6aa1788990e489bf5aef075e8cfa2b5710b4a45e80c359926a0fca7ab500 DO NOT REMOVE, DO UPDATE

package service_test

import (
	"testing"

	"github.com/moov-io/ach-web-viewer/pkg/service"
	"github.com/moov-io/base/config"
	"github.com/moov-io/base/log"

	"github.com/stretchr/testify/require"
)

func Test_ConfigLoading(t *testing.T) {
	logger := log.NewNopLogger()

	ConfigService := config.NewService(logger)

	gc := &service.GlobalConfig{}
	err := ConfigService.Load(gc)
	require.Nil(t, err)
}
