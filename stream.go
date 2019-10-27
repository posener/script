// Package script provides helper functions to write scripts.
//
// Inspired by https://github.com/bitfield/script, with some improvements:
//
// * Output between streamed commands is a stream and not loaded to memory.
//
// * Better representation and handling of errors.
//
// * Proper incocation, usage and handling of stderr of custom commands.
//
// The script chain is represented by a
// (`Stream`) https://godoc.org/github.com/posener/script#Stream object. While each command in the
// stream is abstracted by the (`Command`) https://godoc.org/github.com/posener/script#Command
// struct. This library provides basic functionality, but can be extended freely.
package script

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/go-multierror"
)

// Stream is a chain of commands: the stdout of each command in the stream feeds the following one.
// The stream object have different method that allow manipulating it, most of them resemble well
// known linux commands.
//
// If a command which does not exist in this library is required, the `PipeTo` function should be
// used, which allows constructing a command object from a given input `io.Reader`.
//
// The Stream object output can be used to be written to the stdout, to a file or to a string using
// the `.To*` methods. It also exposes an `io.ReadCloser` interface which allows the user using the
// stream output for any other usecase.
type Stream struct {
	// Command is the current command of the stream.
	command Command
	// parent points to the command before the current command.
	parent *Stream
}

// Stdin startst a stream from stdin.
func Stdin() Stream {
	return Stream{command: Command{Reader: os.Stdin, Name: "stdin"}}
}

// PipeFn is a function that returns a command given input reader for that command.
type PipeFn func(io.Reader) Command

// PipeTo pipes the current stream to a new command and return the new stream. This function should
// be used to add custom commands that are not available in this library.
func (s Stream) PipeTo(pipeFn PipeFn) Stream {
	c := pipeFn(s.command)
	if c.Reader == nil {
		panic("a command must contain a reader")
	}
	if c.Name == "" {
		panic("A command must contain a name")
	}
	return Stream{command: c, parent: &s}
}

// Read implements the io.Reader interface.
func (s Stream) Read(b []byte) (int, error) {
	return s.command.Read(b)
}

// Close closes all the commands in the current stream and return the errors that occured in all
// of the commands by invoking the Error() function on each one of them.
func (s Stream) Close() error {
	var errors *multierror.Error
	for cur := &s; cur != nil; cur = cur.parent {
		if err := cur.command.error(); err != nil {
			errors = multierror.Append(errors, err)
		}
		if err := cur.command.Close(); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

// ToStdout pipes the stdout of the stream to screen.
func (c Stream) ToStdout() error {
	return writeAndClose(c, os.Stdout)
}

// ToString reads stdout of the stream and returns it as a string.
func (s Stream) ToString() (string, error) {
	var out bytes.Buffer
	err := writeAndClose(s, &out)
	return out.String(), err
}

// ToFile dumps the output of the stream to a file.
func (s Stream) ToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return writeAndClose(s, f)
}

// Discard executes the stream pipeline but discards the output.
func (s Stream) Discard() error {
	return writeAndClose(s, ioutil.Discard)
}

// writeAndClose writes the output of the stream to an io.Writer.
func writeAndClose(s Stream, w io.Writer) error {
	var errors *multierror.Error
	if _, err := io.Copy(w, s); err != nil {
		errors = multierror.Append(errors, err)
	}
	if err := s.Close(); err != nil {
		errors = multierror.Append(errors, err)
	}
	return errors.ErrorOrNil()
}
