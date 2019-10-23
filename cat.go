package script

import (
	"io"
	"os"
)

// Cat outputs the contents of files.
//
// Shell command: cat <path>.
func Cat(paths ...string) Pipe {
	var (
		p       Pipe
		readers []io.Reader
	)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			p.appendError(err)
		} else {
			readers = append(readers, f)
			p.closers = append(p.closers, f)
		}
	}

	p.Out = io.MultiReader(readers...)

	return p
}
