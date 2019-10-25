package script

import "io"

// Head reads only the n first lines of the given reader.
//
// Shell command: `head -n <n>`.
func (s Stream) Head(n int) Stream {
	s.stage = "head"
	s.Reader = &head{r: s, n: n}
	return s
}

type head struct {
	r io.Reader
	n int
}

func (h *head) Read(p []byte) (n int, err error) {
	if h.n <= 0 {
		return 0, io.EOF
	}

	n, err = h.r.Read(p)
	for i := range p[:n] {
		if p[i] == '\n' {
			h.n--
		}
		if h.n == 0 {
			return i + 1, io.EOF
		}
	}
	return n, err
}
