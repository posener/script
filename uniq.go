package script

import (
	"bytes"
	"fmt"
)

// Uniq report or omit repeated lines.
//
// Shell command: `uniq`.
func (s Stream) Uniq(writeCountPrefix bool) Stream {
	s.stage = fmt.Sprintf("uniq %v", writeCountPrefix)
	s = s.LineFn(&uniq{writeCountPrefix: writeCountPrefix})
	return s
}

type uniq struct {
	writeCountPrefix bool
	// last stores the last written line.
	last []byte
	// count is the number of times the `last` was in input.
	count int
}

func (u *uniq) Modify(line []byte) ([]byte, error) {
	if line != nil /* not EOF */ && bytes.Equal(line, u.last) /* Repeated line */ {
		u.count++
		return nil, nil
	}

	// Output the last seen line.
	var out []byte
	if u.count > 0 {
		if u.writeCountPrefix {
			out = []byte(fmt.Sprintf("%d\t", u.count))
		}
		out = append(out, u.last...)
		out = append(out, '\n')
	}

	// Remember the line without the '\n' suffix or count prefix.
	u.last = line
	u.count = 1

	return out, nil
}
