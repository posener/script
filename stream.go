// Package script provides helper functions to write scripts.
//
// Inspired by https://github.com/bitfield/script, with a small modifications:
//
// * Output between streamed commands is a stream and not loaded to memory.
//
// * Better representation of errors and stderr.
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
	// Parent points to the command before the current command.
	Parent *Stream
}

// PipeTo pipes the current stream to a given command and return the new stream.
func (s Stream) PipeTo(c Command) Stream {
	return Stream{Command: c, Parent: &s}
}

// Close closes all the commands in the current stream and return the errors that occured in all
// of the commands by invoking the Error() function on each one of them.
func (s Stream) Close() error {
	var errors *multierror.Error
	for cur := &s; cur != nil; cur = cur.Parent {
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
