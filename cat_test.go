package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCat(t *testing.T) {
	t.Parallel()

	t.Run("One file", func(t *testing.T) {
		got, err := Cat("testdata/a.txt").ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\n", got)
	})

	t.Run("Multiple files", func(t *testing.T) {
		got, err := Cat("testdata/a.txt", "testdata/b.txt").ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\nbb\n", got)
	})

	t.Run("No such file", func(t *testing.T) {
		got, err := Cat("testdata/c.txt").ToString()
		assert.Error(t, err)
		assert.Equal(t, "", got)
	})
}
