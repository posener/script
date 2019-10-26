package script

import (
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Command abstract a single command
type Command struct {
	// Reader reads the output of the command.
	io.Reader
	// Closer closes the command.
	io.Closer
	// Name is the name of the command.
	Name string
	// errors contains all errors that occured in the command.
	errors *multierror.Error
}

func (c Command) Close() error {
	if c.Closer == nil {
		return nil
	}
	return c.Closer.Close()
}

func (c Command) error() error {
	return c.errors.ErrorOrNil()
}

// Append an error to the command.
func (c *Command) AppendError(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	err = errors.Wrap(err, c.Name+": "+fmt.Sprintf(format, args...))
	c.errors = multierror.Append(c.errors, err)
}
