package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUniq(t *testing.T) {
	t.Parallel()

	out, err := Echo("a\na\nb\nbb\na").Uniq().ToString()
	require.NoError(t, err)
	assert.Equal(t, "a\nb\nbb\na\n", out)
}

func TestUniq_count(t *testing.T) {
	t.Parallel()

	out, err := Echo("a\na\nb\nbb\na").Modify(&Uniq{WriteCount: true}).ToString()
	require.NoError(t, err)
	assert.Equal(t, "2\ta\n1\tb\n1\tbb\n1\ta\n", out)
}
