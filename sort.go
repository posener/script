package script

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
)

// Sort returns a stream with lines ordered alphabetically.
//
// Shell command: `wc`.
func (s Stream) Sort(reverse bool) Stream {
	s.stage = fmt.Sprintf("sort %v", reverse)
	var lines []string
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	s.appendError(scanner.Err(), "scanning stream")

	sort.Slice(lines, func(i, j int) bool { return (lines[i] < lines[j]) != reverse })

	var out strings.Builder
	for _, line := range lines {
		_, err := out.WriteString(line + "\n")
		s.appendError(err, "writing line %q", line)
	}
	s.Reader = strings.NewReader(out.String())
	return s
}
