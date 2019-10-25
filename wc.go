package script

import (
	"bufio"
	"fmt"
	"strings"
)

// Count represents the output of `wc` shell command.
type Count struct {
	Stream
	Lines, Words, Chars int
}

// Wc counts the number of lines, words and characters.
//
// Shell command: `wc`.
func (s Stream) Wc() Count {
	s.stage = "wc"

	defer s.close()
	c := Count{Stream: s}
	scanner := bufio.NewScanner(c)
	for scanner.Scan() {
		c.Lines++
		c.Chars += len(scanner.Text()) + 1
		c.Words += countWords(scanner.Text())
	}
	c.appendError(scanner.Err(), "scanning stream")
	c.Reader = strings.NewReader(c.String())
	return c
}

func (c Count) String() string {
	return fmt.Sprintf("%d\t%d\t%d\n", c.Lines, c.Words, c.Chars)
}

func countWords(s string) int {
	// TODO: improve performance.
	return len(strings.Fields(s))
}
