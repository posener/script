package script

import (
	"fmt"
	"regexp"
)

// Regexp filters only line that match the given regexp.
//
// Shell command: `grep <re>`.
func (s Stream) Regexp(re *regexp.Regexp) Stream {
	return s.Modify(Regexp{Re: re})
}

// Regexp is a modifier that filters only line that match `Re`. If Invert was set only line that did
// not match the regex will be returned.
//
// Shell command: `grep [-v <Invert>] <Re>`.
type Regexp struct {
	Re     *regexp.Regexp
	Invert bool
}

func (g Regexp) Modify(line []byte) (modifed []byte, err error) {
	if line == nil {
		return nil, nil
	}
	if g.Re.Match(line) != g.Invert {
		return append(line, '\n'), nil
	}
	return nil, nil
}

func (g Regexp) Name() string {
	return fmt.Sprintf("grep(%v, invert=%v)", g.Re, g.Invert)
}
