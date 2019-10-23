package script

import (
	"bufio"
	"strings"
)

// Count represents the output of `wc` shell command.
type Count struct {
	Lines, Words, Chars int
}

// Wc counts the number of lines, words and characters.
//
// Shell command: `wc`.
func (p Pipe) Wc() Count {
	defer p.Close()
	var c Count
	scanner := bufio.NewScanner(p.Out)
	for scanner.Scan() {
		c.Lines++
		c.Chars += len(scanner.Text())
		c.Words += countWords(scanner.Text())
	}
	c.Chars += c.Lines
	p.appendError(scanner.Err())
	return c
}

func countWords(s string) int {
	// TODO: improve performance.
	return len(strings.Fields(s))
}
