package script

import (
	"bufio"
	"io"
	"reflect"
)

// Modifer modifies input lines to output. On each line of the input the Modify method is called,
// and the modifier can change it, omit it, or break the iteration.
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
	// Name returns the name of the command that will represent this modifier.
	Name() string
}

// ModifierFn is a function for modifying input lines.
type ModifierFn func(line []byte) (modifed []byte, err error)

func (m ModifierFn) Modify(line []byte) (modifed []byte, err error) {
	return m(line)
}

func (m ModifierFn) Name() string {
	return reflect.TypeOf(m).Name()
}

// Modify applies modifier on every line of the input.
func (s Stream) Modify(modifier Modifer) Stream {
	return s.PipeTo(pipeModifier(modifier))
}

func pipeModifier(m Modifer) PipeFn {
	return func(stdin io.Reader) Command {
		return Command{
			Name:   m.Name(),
			Reader: &modifier{r: bufio.NewReader(stdin), modifier: m},
		}
	}
}

type modifier struct {
	r        *bufio.Reader
	modifier Modifer
	// partialOut stores leftover of a line that was not fully read by output.
	partialOut []byte
	err        error
}

func (m *modifier) Read(out []byte) (n int, err error) {
	if len(m.partialOut) > 0 {
		m.partialOut, n = copyBytes(out, m.partialOut)
		return n, nil
	}
	if m.err != nil {
		return 0, m.err
	}

	// partialIn stores a line that was not fully read from input.
	var partialIn []byte

	for {
		line, isPrefix, err := m.r.ReadLine()
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			// Remember that we have EOF for next read call.
			m.err = io.EOF
		}
		if len(partialIn) > 0 {
			line = append(partialIn, line...)
			partialIn = nil
		}
		if isPrefix {
			partialIn = line
			continue
		}

		line, err = m.modifier.Modify(line)
		if err != nil {
			m.err = err
		}

		m.partialOut, n = copyBytes(out, line)
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
