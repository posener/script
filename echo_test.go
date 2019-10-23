package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcho(t *testing.T) {
	t.Parallel()

	s, err := Echo("hello world").ToString()
	require.NoError(t, err)
	assert.Equal(t, "hello world\n", s)
}
