package script

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Command abstract a single command
type Command interface {
	// Reader reads the output of the command.
	io.Reader
	// Closer closes the command.
	io.Closer
	// Name returns the name of the command.
	Name() string
	// Error returns error occured while executing the command.
	Error() error
}

// command is a simple implementation of the command interface.
type command struct {
	io.Reader
	io.Closer
	name   string
	errors *multierror.Error
}

func (c command) Close() error {
	if c.Closer == nil {
		return nil
	}
	return c.Closer.Close()
}

func (c command) Name() string {
	return c.name
}

func (c command) Error() error {
	return c.errors.ErrorOrNil()
}

func (c *command) appendError(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	err = errors.Wrap(err, c.Name()+": "+fmt.Sprintf(format, args...))
	c.errors = multierror.Append(c.errors, err)
}
