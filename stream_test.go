package script

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A simple "hello world" example that creats a stream and pipe it to the stdout.
func Example_helloWorld() {
	// Create an "hello world" stream and use the ToStdout method to write it to stdout.
	Echo("hello world").ToStdout()

	// Output: hello world
}

// An example that shows how to iterate scanned lines.
func Example_iterate() {
	// Stream can be any stream, in this case we have echoed 3 lines.
	stream := Echo("first\nsecond\nthird")

	// To iterate over the stream lines, it is better not to read it into memory and split over the
	// lines, but use the `bufio.Scanner`:
	defer stream.Close()
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Output: first
	// second
	// third
}

// An example that shows how to create custom commands using the `PipeTo` method with a `PipeFn`
// function.
func Example_pipeTo() {
	Echo("1\n2\n3").PipeTo(func(r io.Reader) Command {
		// Create a command that sums up all numbers in input.

		// Use buffered reader to read lines from input.
		buf := bufio.NewReader(r)

		// Store the sum of all numbers.
		sum := 0

		// Create a reade function (for the sake of example function, such that all code will fit
		// into the example function body - a proper more readable way to do it was to create a new
		// type with a state and a read function).
		read := func(b []byte) (int, error) {
			// Read next line from input.
			line, _, err := buf.ReadLine()

			// if EOF write sum to output.
			if err == io.EOF {
				return copy(b, append([]byte(strconv.Itoa(sum)), '\n')), nil
			}

			// Convert the line to a number and add it to the sum.
			if i, err := strconv.Atoi(string(line)); err == nil {
				sum += i
			}

			// We don't write anything to output, so we return 0 bytes with no error.
			return 0, nil
		}

		return Command{Name: "sum", Reader: readerFn(read)}
	})

	// Output: 6
}

type readerFn func(b []byte) (int, error)

func (f readerFn) Read(b []byte) (int, error) {
	return f(b)
}

type sum struct {
	r *bufio.Reader
	n int
}

func (s sum) Read(b []byte) (int, error) {
	line, _, err := s.r.ReadLine()
	// if EOF write sum to output.
	if err == io.EOF {
		return copy(b, append([]byte(strconv.Itoa(s.n)), '\n')), nil
	}
	if i, err := strconv.Atoi(string(line)); err == nil {
		s.n += i
	}
	return 0, nil
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
