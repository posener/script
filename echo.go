package script

import "strings"

// Echo writes to stdout.
//
// Shell command: `echo <s>`
func Echo(s string) Stream {
	return Stream{Command: command{
		name:   "echo",
		Reader: strings.NewReader(s + "\n"),
	}}
}
