package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniq(t *testing.T) {
	t.Parallel()

	out, err := Echo("a\na\nb\nbb\na").Uniq(false).ToString()
	require.NoError(t, err)
	assert.Equal(t, "a\nb\nbb\na\n", out)
}
