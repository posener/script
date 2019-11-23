package script

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexp(t *testing.T) {
	t.Parallel()

	t.Run("grep", func(t *testing.T) {
		got, err := Echo("a\nb\na\nc").Regexp(regexp.MustCompile(`^a`)).ToString()
		require.NoError(t, err)
		assert.Equal(t, "a\na\n", got)
	})

	t.Run("invert", func(t *testing.T) {
		got, err := Echo("a\nb\na\nc").Modify(Regexp{Re: regexp.MustCompile(`^a`), Invert: true}).ToString()
		require.NoError(t, err)
		assert.Equal(t, "b\nc\n", got)
	})
}
