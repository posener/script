package script

import (
	"fmt"
	"os/exec"
)

// Exec executes a command, and pipes its stdout.
// TODO: make command return two streams? one for out and one for err?
func Exec(command string, args ...string) ErrStream {
	var s Stream
	return s.Exec(command, args...)
}

// Exec executes a command, and pipes its stdout.
//
// If the pipe already contains a reader, it will pipe it into the command line.
// The output is an ErrStream which contains the standard output and standard error of the executed
// command. Both of them are represented as stream and can be used as streams.
func (s Stream) Exec(command string, args ...string) ErrStream {
	s.stage = fmt.Sprintf("exec %v %+v", command, args)
	cmd := exec.Command(command, args...)

	// pipe previous pipe to stdin if available.
	if s.Reader != nil {
		cmd.Stdin = s.Reader
	}

	// pipe stdout and stderr to the new pipe.
	cmdOut, err := cmd.StdoutPipe()
	s.appendError(err, "pipe stdout")

	cmdErr, err := cmd.StderrPipe()
	s.appendError(err, "pipe stderr")

	// start the process
	err = cmd.Start()
	s.appendError(err, "start process")

	s.appendCloser(closerFn(func() error { return cmd.Wait() }))

	return newErrStream(s, cmdOut, cmdErr)
}

type closerFn func() error

func (f closerFn) Close() error { return f() }
