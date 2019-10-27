package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCut(t *testing.T) {
	t.Parallel()

	t.Run("multiple lines", func(t *testing.T) {
		got, err := Echo("a\tbb\tccc\nddd\tee\tf").Cut(1, 3).ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\tccc\nddd\tf\n", got)
	})

	t.Run("double separator", func(t *testing.T) {
		got, err := Echo("a\t\tb").Cut(1, 2).ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\t\n", got)
	})

	t.Run("field out of range", func(t *testing.T) {
		got, err := Echo("a\tb").Cut(1, 3).ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\n", got)
	})

	t.Run("custom delimiter", func(t *testing.T) {
		got, err := Echo("a b c").Modify(Cut{Delim: []byte{' '}, Fields: []int{1, 3}}).ToString()
		require.NoError(t, err)
		assert.Equal(t, "a c\n", got)
	})
}
