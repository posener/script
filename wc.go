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
	defer s.Close()

	var count Count
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		count.Lines++
		count.Chars += len(scanner.Text()) + 1
		count.Words += countWords(scanner.Text())
	}
	c := command{
		name:   "wc",
		Reader: strings.NewReader(count.String()),
	}
	c.appendError(scanner.Err(), "scanning stream")
	count.Stream = s.PipeTo(c)
	return count
}

func (c Count) String() string {
	return fmt.Sprintf("%d\t%d\t%d\n", c.Lines, c.Words, c.Chars)
}

func countWords(s string) int {
	// TODO: improve performance.
	return len(strings.Fields(s))
}
