package script

import (
	"os/exec"
)

// Exec executes a command, and pipes its stdout.
func Exec(command string, args ...string) Pipe {
	var p Pipe
	return p.Exec(command, args...)
}

// Exec executes a command, and pipes its stdout.
//
// If the pipe already contains a reader, it will pipe it into the command line.
func (p Pipe) Exec(command string, args ...string) Pipe {
	cmd := exec.Command(command, args...)

	// pipe previous pipe to stdin if available.
	if p.Out != nil {
		cmd.Stdin = p.Out
	}

	// pipe stdout and stderr to the new pipe.
	cmdOut, err := cmd.StdoutPipe()
	p.appendError(err)
	p.Out = cmdOut

	cmdErr, err := cmd.StderrPipe()
	p.appendError(err)
	p.appendErrorReader(cmdErr)

	// start the process
	err = cmd.Start()
	p.appendError(err)

	p.appendCloser(closerFn(func() error { return cmd.Wait() }))

	return p
}

type closerFn func() error

func (f closerFn) Close() error {
	return f()
}
