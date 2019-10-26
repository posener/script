package script

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

// Exec executes a command, and pipes its stdout.
func Exec(command string, args ...string) Stream {
	var s Stream
	return s.exec(nil, command, args...)
}

// ExecHandleStderr executes a command, and pipes its stdout and enable collecting the stderr of the
// command.
func ExecHandleStderr(errWriter io.Writer, command string, args ...string) Stream {
	var s Stream
	return s.exec(errWriter, command, args...)
}

// Exec executes a command, and pipes its stdout.
//
// If the pipe already contains a reader, it will pipe it into the command line.
func (s Stream) Exec(command string, args ...string) Stream {
	return s.exec(nil, command, args...)
}

// ExecHandleStderr executes a command, and pipes its stdout and enable collecting the stderr of the
// command.
func (s Stream) ExecHandleStderr(errWriter io.Writer, command string, args ...string) Stream {
	return s.exec(errWriter, command, args...)
}

func (s Stream) exec(errWriter io.Writer, command string, args ...string) Stream {
	s.stage = fmt.Sprintf("exec %v %+v", command, args)
	cmd := exec.Command(command, args...)

	// pipe previous pipe to stdin if available.
	if s.Reader != nil {
		cmd.Stdin = s.Reader
	}

	// pipe stdout and stderr to the new pipe.
	cmdOut, err := cmd.StdoutPipe()
	s.appendError(err, "pipe stdout")
	s.Reader = cmdOut

	if errWriter == nil {
		errWriter = ioutil.Discard
	}
	cmd.Stderr = errWriter

	// start the process
	err = cmd.Start()
	s.appendError(err, "start process")

	s.appendCloser(closerFn(func() error { return cmd.Wait() }))

	return s
}

type closerFn func() error

func (f closerFn) Close() error { return f() }
