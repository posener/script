package script

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Sort returns a stream with lines ordered alphabetically.
//
// Shell command: `wc`.
func (s Stream) Sort(reverse bool) Stream {
	return s.PipeTo(pipeSort(reverse))
}

func pipeSort(reverse bool) PipeFn {
	return func(stdin io.Reader) Command {
		c := Command{Name: fmt.Sprintf("sort(%v)", reverse)}
		var lines []string
		scanner := bufio.NewScanner(stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		c.AppendError(scanner.Err(), "scanning stream")

		sort.Slice(lines, func(i, j int) bool { return (lines[i] < lines[j]) != reverse })

		var out strings.Builder
		for _, line := range lines {
			_, err := out.WriteString(line + "\n")
			c.AppendError(err, "writing line %q", line)
		}
		c.Reader = strings.NewReader(out.String())
		return c
	}
}
