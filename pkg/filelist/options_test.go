package filelist

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestListOpts__Inside(t *testing.T) {
	req := &http.Request{}
	opts, err := ReadListOptions(req)
	require.NoError(t, err)

	require.True(t, opts.Inside(time.Now()))

	future := time.Now().Add(7 * 24 * time.Hour)
	require.False(t, opts.Inside(future))

	past := time.Now().Add(-30 * 7 * 24 * time.Hour)
	require.False(t, opts.Inside(past))
}
