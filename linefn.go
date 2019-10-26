package script

import (
	"bufio"
	"io"
)

// LineModifier modifies a line.
type Modifer interface {
	// Modify a line. The input of this function will always be a single line from the input of the
	// stream, without the trailing '\n'. It should return the output of the stream and should
	// append a trailing '\n' if it want it to be a line in the output.
	//
	// When EOF of input stream is met, the function will be called once more with a nil line value
	// to enable output any buffered data.
	//
	// When the return modified value is nil, the line will be dicarded.
	//
	// When the returned eof value is true, the Read will return that error.
	Modify(line []byte) (modifed []byte, err error)
}

// LineModifierFn is a function that modifies a line.
type ModifierFn func(line []byte) (modifed []byte, err error)

func (m ModifierFn) Modify(line []byte) (modifed []byte, err error) {
	return m(line)
}

// LineFn applies modifier on every line of the input.
func (s Stream) LineFn(name string, modifier Modifer) Stream {
	return s.PipeTo(command{
		name:   name,
		Reader: &lineFn{r: bufio.NewReader(s.Command), modifier: modifier},
	})
}

type lineFn struct {
	r        *bufio.Reader
	modifier Modifer
	// partialOut stores leftover of a line that was not fully read by output.
	partialOut []byte
	err        error
}

func (u *lineFn) Read(out []byte) (n int, err error) {
	if len(u.partialOut) > 0 {
		u.partialOut, n = copyBytes(out, u.partialOut)
		return n, nil
	}
	if u.err != nil {
		return 0, u.err
	}

	// partialIn stores a line that was not fully read from input.
	var partialIn []byte

	for {
		line, isPrefix, err := u.r.ReadLine()
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			// Remember that we have EOF for next read call.
			u.err = io.EOF
		}
		if len(partialIn) > 0 {
			line = append(partialIn, line...)
			partialIn = nil
		}
		if isPrefix {
			partialIn = line
			continue
		}

		line, err = u.modifier.Modify(line)
		if err != nil {
			u.err = err
		}

		u.partialOut, n = copyBytes(out, line)
		return n, nil
	}
}

func copyBytes(dst, src []byte) (leftover []byte, n int) {
	n = len(src)
	if n > len(dst) {
		n = len(dst)
	}
	copy(dst[:n], src[:n])
	return src[n:], n
}
