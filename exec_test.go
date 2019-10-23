package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()

	t.Run("Without stdin", func(t *testing.T) {
		got, err := Exec("echo", "hello world").ToString()
		require.NoError(t, err)
		assert.Equal(t, "hello world\n", got)
	})

	t.Run("With stdin", func(t *testing.T) {
		got, err := Echo("hello world").Exec("cat").ToString()
		require.NoError(t, err)
		assert.Equal(t, "hello world\n", got)
	})

	t.Run("exit code", func(t *testing.T) {
		got, err := Exec("false").ToString()
		assert.Error(t, err)
		assert.Equal(t, "", got)
	})
}
