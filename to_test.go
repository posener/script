package script

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToFile(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "script")
	require.NoError(t, err)

	path := filepath.Join(dir, "file")

	err = Echo("hello world").ToFile(path)
	require.NoError(t, err)
	defer os.Remove(path)

	got, err := Cat(path).ToString()
	require.NoError(t, err)

	assert.Equal(t, "hello world\n", got)
}

func TestAppendFile(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir("", "script")
	require.NoError(t, err)

	path := filepath.Join(dir, "file")

	err = Echo("hello world").AppendFile(path)
	require.NoError(t, err)
	defer os.Remove(path)

	err = Echo("hello world").AppendFile(path)
	require.NoError(t, err)

	got, err := Cat(path).ToString()
	require.NoError(t, err)

	assert.Equal(t, "hello world\nhello world\n", got)
}

func TestToTempFile(t *testing.T) {
	t.Parallel()

	tmp, err := Echo("hello world").ToTempFile()
	require.NoError(t, err)
	defer os.Remove(tmp)

	got, err := Cat(tmp).ToString()
	require.NoError(t, err)

	assert.Equal(t, "hello world\n", got)
}

func TestIterate(t *testing.T) {
	t.Parallel()
	out := []byte{}
	Echo("a\nb\nc").Iterate(func(l []byte) error {
		out = append(out, l...)
		return nil
	})
	assert.Equal(t, out, []byte("abc"))
}
