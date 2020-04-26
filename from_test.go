package script

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcho(t *testing.T) {
	t.Parallel()

	s, err := Echo("hello world").ToString()
	require.NoError(t, err)
	assert.Equal(t, "hello world\n", s)
}

func TestStdin(t *testing.T) {
	// Create a temporary file to fake stdin.
	fakeStdin, err := ioutil.TempFile("", "")
	require.NoError(t, err)
	defer os.Remove(fakeStdin.Name())
	defer fakeStdin.Close()
	_, err = fakeStdin.WriteString("hello world\n")
	require.NoError(t, err)
	_, err = fakeStdin.Seek(0, 0)
	require.NoError(t, err)

	// Temporarely replace stdin with the temporary file.
	stdin := os.Stdin
	os.Stdin = fakeStdin
	defer func() { os.Stdin = stdin }()

	// Test
	s, err := Stdin().ToString()
	require.NoError(t, err)
	assert.Equal(t, "hello world\n", s)
}

func TestWriter(t *testing.T) {
	t.Parallel()
	got, err := Writer("json", func(w io.Writer) error { return json.NewEncoder(w).Encode("foo") }).ToString()
	require.NoError(t, err)
	assert.Equal(t, "\"foo\"\n", got)
}

func TestWriter_failure(t *testing.T) {
	t.Parallel()
	_, err := Writer("fail", func(w io.Writer) error { return errors.New("failed") }).ToString()
	assert.Error(t, err)
}
