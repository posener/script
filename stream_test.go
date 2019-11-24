package script

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
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

// An example that shows how to create custom commands using the `Through` method with a `PipeFn`
// function.
func Example_through() {
	Echo("1\n2\n3").Through(PipeFn(func(r io.Reader) (io.Reader, error) {
		// Create a command that sums up all numbers in input.
		//
		// In this example we create a reader function such that the whole code will fit into the
		// example function body. A more proper and readable way to do it was to create a new
		// type with a state that implements the `io.Reader` interface.

		// Use buffered reader to read lines from input.
		buf := bufio.NewReader(r)

		// Store the sum of all numbers.
		sum := 0

		// Read function reads the next line and adds it to the sum. If it gets and EOF error, it
		// writes the sum to the output and returns an EOF.
		read := func(b []byte) (int, error) {
			// Read next line from input.
			line, _, err := buf.ReadLine()

			// if EOF write sum to output.
			if err == io.EOF {
				return copy(b, append([]byte(strconv.Itoa(sum)), '\n')), io.EOF
			}
			if err != nil {
				return 0, err
			}

			// Convert the line to a number and add it to the sum.
			if i, err := strconv.Atoi(string(line)); err == nil {
				sum += i
			}

			// We don't write anything to output, so we return 0 bytes with no error.
			return 0, nil
		}

		return readerFn(read), nil
	})).ToStdout()

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
