package script

import (
	"bytes"
	"fmt"
	"io"
)

// Head reads only the n first lines of the given reader. If n is a negative number, the last (-n)
// lines will be returned.
//
// Shell command: `head -n <n>` / `tail -n <-n>`.
func (s Stream) Head(n int) Stream {
	var mod Modifer
	if n >= 0 {
		h := head(n)
		mod = &h
	} else {
		t := make(tail, 0, -n)
		mod = &t
	}
	return s.LineFn(fmt.Sprintf("head(%d)", n), mod)
}

type head int

func (n *head) Modify(line []byte) ([]byte, error) {
	if line == nil || *n <= 0 {
		return nil, io.EOF
	}
	*n--
	return append(line, '\n'), nil
}

type tail [][]byte

func (t *tail) Modify(line []byte) ([]byte, error) {
	if line == nil {
		return append(bytes.Join(*t, []byte{'\n'}), '\n'), io.EOF
	}

	// Shift all lines and append the new line.
	if len(*t) < cap(*t) {
		*t = append(*t, line)
	} else {
		for i := 0; i < len(*t)-1; i++ {
			(*t)[i] = (*t)[i+1]
		}
		(*t)[len(*t)-1] = line
	}

	return nil, nil
}
