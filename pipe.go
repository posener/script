// Package script provides helper functions to write scripts.
//
// Inspired from: https://github.com/bitfield/script, with a small modifications:
//
// * Output between piped commands is a stream and not loaded to memory.
//
// * Better representation of errors.
package script

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// Pipe is a data stracture that chains scriptable commands.
type Pipe struct {
	// Out holds the stdout of the current pipe stage.
	Out io.Reader
	// Err holds the stderr of the current pipe stage.
	Err io.Reader
	// closers holds all close functions for the chain of pipes until this stage.
	closers []io.Closer
}

// ToString reads stdout and returns it as a string, and read stderr and returns it as error if
// exists.
func (p Pipe) ToString() (string, error) {
	out, err := ioutil.ReadAll(p.Out)
	p.appendError(err)

	err = readError(p.Err)

	if closeError := p.Close(); closeError != nil {
		err = multierror.Append(err, closeError).ErrorOrNil()
	}

	return string(out), err
}

// ToScreen pipes the stdout and stderr to screen.
func (p Pipe) ToScreen() {
	defer func() {
		if err := p.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}()
	io.Copy(os.Stdout, p.Out)
	if p.Err != nil {
		io.Copy(os.Stderr, p.Err)
	}
}

// Close the pipe, should be called after the Out and Err were readen, and should be called only
// once.
func (p Pipe) Close() error {
	var err *multierror.Error
	for _, c := range p.closers {
		err = multierror.Append(err, c.Close())
	}
	return err.ErrorOrNil()
}

func (p *Pipe) appendCloser(c io.Closer) {
	if c == nil {
		return
	}
	p.closers = append(p.closers, c)
}

func (p *Pipe) appendError(err error) {
	if err == nil {
		return
	}
	p.appendErrorReader(strings.NewReader(err.Error() + "\n"))
}

func (p *Pipe) appendErrorReader(r io.Reader) {
	if r == nil {
		return
	}
	if p.Err == nil {
		p.Err = r
	} else {
		p.Err = io.MultiReader(p.Err, r)
	}
}

func readError(r io.Reader) error {
	if r == nil {
		return nil
	}
	msg, err := ioutil.ReadAll(r)
	if err != nil || len(msg) == 0 {
		return err
	}
	return fmt.Errorf(string(msg))
}
