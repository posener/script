package script

import "strings"

// Echo writes to stdout.
//
// Shell command: `echo <s>`
func Echo(s string) Stream {
	return Stream{command: Command{
		Name:   "echo",
		Reader: strings.NewReader(s + "\n"),
	}}
}
