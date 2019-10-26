package script

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExec(t *testing.T) {
	t.Parallel()

	t.Run("without stdin", func(t *testing.T) {
		stdout, err := Exec("echo", "hello world").ToString()

		require.NoError(t, err)
		assert.Equal(t, "hello world\n", stdout)
	})

	t.Run("with stdin", func(t *testing.T) {
		stdout, err := Echo("hello world").Exec("cat").ToString()

		require.NoError(t, err)
		assert.Equal(t, "hello world\n", stdout)
	})

	t.Run("exit code", func(t *testing.T) {
		stdout, err := Exec("false").ToString()

		assert.Error(t, err)
		assert.Equal(t, "", stdout)
	})

	t.Run("stderr", func(t *testing.T) {
		var stderr bytes.Buffer
		stdout, err := ExecHandleStderr(&stderr, "cat", "no-such-file", "testdata/a.txt").ToString()

		assert.Error(t, err)
		assert.Equal(t, "a\n", stdout) // Content of testdata/a.txt
		assert.Equal(t, "cat: no-such-file: No such file or directory\n", stderr.String())
	})
}
