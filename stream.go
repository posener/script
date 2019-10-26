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
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Stream chains commands in a way that enable piping the output of one command to the other.
type Stream struct {
	// Reader holds the output of the current stream stage.
	io.Reader
	// stage holds the name of the current stage.
	stage string
	// closers holds all close functions for the chain of commands until this stage.
	closers []io.Closer
	// errors hold all the errors of all the chain of commands in the stream.
	errors *multierror.Error
}

// ToScreen pipes the stdout and stderr to screen.
func (s Stream) ToScreen() error {
	s.writeTo(os.Stdout)
	s.close()
	return s.error()
}

// ToString reads stdout and returns it as a string.
func (s Stream) ToString() (string, error) {
	var out bytes.Buffer
	s.writeTo(&out)
	s.close()
	return out.String(), s.error()
}

// writeTo writes the output of the stream to an io.Writer.
func (s *Stream) writeTo(w io.Writer) {
	_, err := io.Copy(w, s)
	s.appendError(err, "copy to writer")
}

// close the stream.
//
// Should be called after the Reader was readen and should be called only once.
func (s *Stream) close() {
	for _, c := range s.closers {
		err := c.Close()
		s.appendError(err, "closing")
	}
}

func (p *Stream) appendCloser(c io.Closer) {
	if c == nil {
		return
	}
	p.closers = append(p.closers, c)
}

func (s *Stream) appendError(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	err = errors.Wrap(err, s.stage+": "+fmt.Sprintf(format, args...))
	s.errors = multierror.Append(s.errors, err)
}

func (s Stream) error() error {
	return s.errors.ErrorOrNil()
}
