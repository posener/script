package script

import (
	"io"
	"os"

	"github.com/hashicorp/go-multierror"
)

// Cat outputs the contents of files.
//
// Shell command: cat <path>.
func Cat(paths ...string) Stream {
	c := command{name: "cat"}
	var (
		readers []io.Reader
		closers multicloser
	)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			c.appendError(err, "open path: %s", path)
		} else {
			readers = append(readers, f)
			closers = append(closers, f)
		}
	}

	c.Reader = io.MultiReader(readers...)
	c.Closer = closers

	return Stream{Command: c}
}

type multicloser []io.Closer

func (mc multicloser) Close() error {
	var errors *multierror.Error
	for _, c := range mc {
		if err := c.Close(); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}
