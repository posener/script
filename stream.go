// Package script provides helper functions to write scripts.
//
// Inspired by https://github.com/bitfield/script, with some small modifications:
//
// * Output between streamed commands is a stream and not loaded to memory.
//
// * Better representation of errors and stderr.
//
// The script chain is represented by the
// [`Stream`](https://godoc.org/github.com/posener/script#Stream) type. While each command in the
// stream is abstracted by the [`Command`](https://godoc.org/github.com/posener/script#Command)
// interface, which enable extending this library freely.
package script

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/go-multierror"
)

// Stream is a chain of commands: the stdout of one command feeds the following one.
type Stream struct {
	// Command is the current command of the stream.
	Command
	// parent points to the command before the current command.
	parent *Stream
}

// Stdin startst a stream from stdin.
func Stdin() Stream {
	return Stream{Command: command{Reader: os.Stdin, name: "stdin"}}
}

// PipeFn is a function that returns a command given input reader for that command.
type PipeFn func(io.Reader) Command

// PipeTo pipes the current stream to a given command and return the new stream.
func (s Stream) PipeTo(pipeFn PipeFn) Stream {
	return Stream{Command: pipeFn(s.Command), parent: &s}
}

// Close closes all the commands in the current stream and return the errors that occured in all
// of the commands by invoking the Error() function on each one of them.
func (s Stream) Close() error {
	var errors *multierror.Error
	for cur := &s; cur != nil; cur = cur.parent {
		if err := cur.Command.Error(); err != nil {
			errors = multierror.Append(errors, err)
		}
		if err := cur.Command.Close(); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}

// ToScreen pipes the stdout of the stream to screen.
func (c Stream) ToScreen() error {
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
