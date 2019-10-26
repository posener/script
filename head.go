package script

import "io"

// Head reads only the n first lines of the given reader.
//
// Shell command: `head -n <n>`.
func (s Stream) Head(n int) Stream {
	s.stage = "head"
	h := head(n)
	return s.LineFn(&h)
}

type head int

func (n *head) Modify(line []byte) ([]byte, error) {
	if line == nil || *n <= 0 {
		return nil, io.EOF
	}
	*n--
	return append(line, '\n'), nil
}
