package script

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Example that shows how to write piped content to screen.
func Example() {
	Echo("hello world").ToScreen()
	// Output: hello world
}

func TestToFile(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "script")
	require.NoError(t, err)

	path := filepath.Join(dir, "file")

	err = Echo("hello world").ToFile(path)
	require.NoError(t, err)

	got, err := Cat(path).ToString()
	require.NoError(t, err)

	assert.Equal(t, "hello world\n", got)
}
