package script

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrep(t *testing.T) {
	t.Parallel()

	t.Run("grep", func(t *testing.T) {
		got, err := Echo("a\nb\na\nc").Grep(regexp.MustCompile(`^a`)).ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\na\n", got)
	})

	t.Run("invert", func(t *testing.T) {
		got, err := Echo("a\nb\na\nc").Modify(Grep{Re: regexp.MustCompile(`^a`), Inverse: true}).ToString()
		require.NoError(t, err)
		assert.Equal(t, "b\nc\n", got)
	})
}
