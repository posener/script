package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {
	t.Parallel()

	t.Run("sort", func(t *testing.T) {
		out, err := Echo("ab\na\nb").Sort(false).ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\nab\nb\n", out)
	})

	t.Run("sort reversed", func(t *testing.T) {
		out, err := Echo("ab\na\nb").Sort(true).ToString()
		require.NoError(t, err)
		assert.Equal(t, "b\nab\na\n", out)
	})
}
