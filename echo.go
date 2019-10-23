package script

import "strings"

// Echo writes to stdout.
//
// Shell command: `echo <s>`
func Echo(s string) Pipe {
	return Pipe{Out: strings.NewReader(s + "\n")}
}
