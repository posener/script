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
	"io/ioutil"
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
	s.writeAndClose(os.Stdout)
	return s.error()
}

// ToString reads stdout and returns it as a string.
func (s Stream) ToString() (string, error) {
	var out bytes.Buffer
	s.writeAndClose(&out)
	return out.String(), s.error()
}

// ToFile dumps the output of the stream to a file.
func (s Stream) ToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	s.writeAndClose(f)
	return s.error()
}

// Discard executes the stream pipeline but discards the output.
func (s Stream) Discard() error {
	s.writeAndClose(ioutil.Discard)
	return s.error()
}

// writeAndClose writes the output of the stream to an io.Writer.
func (s *Stream) writeAndClose(w io.Writer) {
	_, err := io.Copy(w, s)
	s.appendError(err, "copy to writer")
	s.close()
}

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
