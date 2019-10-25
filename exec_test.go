package script

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()

	t.Run("Without stdin", func(t *testing.T) {
		stdout, stderr, err := Exec("echo", "hello world").ToString()

		require.NoError(t, err)
		assert.Equal(t, "hello world\n", stdout)
		assert.Equal(t, "", stderr)
	})

	t.Run("With stdin", func(t *testing.T) {
		stdout, stderr, err := Echo("hello world").Exec("cat").ToString()

		require.NoError(t, err)
		assert.Equal(t, "hello world\n", stdout)
		assert.Equal(t, "", stderr)
	})

	t.Run("stderr", func(t *testing.T) {
		stdout, stderr, err := Exec("cat", "no-such-file").ToString()

		assert.Error(t, err)
		assert.Equal(t, "", stdout)
		assert.Equal(t, "cat: no-such-file: No such file or directory\n", stderr)
	})

	t.Run("stderr only", func(t *testing.T) {
		stderr, err := Exec("cat", "no-such-file").Err.ToString()

		assert.Error(t, err)
		assert.Equal(t, "cat: no-such-file: No such file or directory\n", stderr)
	})

	t.Run("exit code", func(t *testing.T) {
		stdout, stderr, err := Exec("false").ToString()

		assert.Error(t, err)
		assert.Equal(t, "", stdout)
		assert.Equal(t, "", stderr)
	})
}
