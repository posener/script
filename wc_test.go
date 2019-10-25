package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWc(t *testing.T) {
	t.Parallel()

	t.Run("Multiple lines words and chars", func(t *testing.T) {
		wc := Echo("a b c\nd e \ng ").Wc()

		assert.Equal(t, 3, wc.Lines)
		assert.Equal(t, 6, wc.Words)
		assert.Equal(t, 14, wc.Chars)

		out, err := wc.ToString()
		require.NoError(t, err)
		assert.Equal(t, "3\t6\t14\n", out)
	})

	t.Run("Empty text", func(t *testing.T) {
		wc := Echo("").Wc()

		assert.Equal(t, 1, wc.Lines)
		assert.Equal(t, 0, wc.Words)
		assert.Equal(t, 1, wc.Chars)
	})
}
