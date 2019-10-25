package script

import (
	"io"
	"os"
)

// Cat outputs the contents of files.
//
// Shell command: cat <path>.
func Cat(paths ...string) Stream {
	s := Stream{stage: "cat"}
	var readers []io.Reader

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			s.appendError(err, "open path: %s", path)
		} else {
			readers = append(readers, f)
			s.closers = append(s.closers, f)
		}
	}

	s.Reader = io.MultiReader(readers...)

	return s
}
