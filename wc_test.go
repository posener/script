package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWc(t *testing.T) {
	t.Parallel()

	t.Run("Multiple lines words and chars", func(t *testing.T) {
		wc := Echo("a b c\nd e \ng ").Wc()
		assert.Equal(t, Count{Lines: 3, Words: 6, Chars: 14}, wc)
	})

	t.Run("Empty text", func(t *testing.T) {
		wc := Echo("").Wc()
		assert.Equal(t, Count{Lines: 1, Chars: 1}, wc)
	})
}
