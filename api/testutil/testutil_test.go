package testutil_test

import (
	"testing"
	"time"

	"codeberg.org/jessienyan/booruview/testutil"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {
	require.Equal(t, "2026-04-01T01:23:45Z", testutil.Time().Format(time.RFC3339))
}
