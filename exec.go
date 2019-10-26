package script

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

var stdin = Stream{Command: command{Reader: os.Stdin, name: "stdin"}}

// Exec executes a command, and pipes its stdout.
func Exec(command string, args ...string) Stream {
	return stdin.exec(nil, command, args...)
}

// ExecHandleStderr executes a command, and pipes its stdout and enable collecting the stderr of the
// command.
func ExecHandleStderr(errWriter io.Writer, cmd string, args ...string) Stream {
	return stdin.exec(errWriter, cmd, args...)
}

// Exec executes a command, and pipes its stdout.
//
// If the pipe already contains a reader, it will pipe it into the command line.
func (s Stream) Exec(cmd string, args ...string) Stream {
	return s.exec(nil, cmd, args...)
}

// ExecHandleStderr executes a command, and pipes its stdout and enable collecting the stderr of the
// command.
func (s Stream) ExecHandleStderr(errWriter io.Writer, cmd string, args ...string) Stream {
	return s.exec(errWriter, cmd, args...)
}

func (s Stream) exec(errWriter io.Writer, name string, args ...string) Stream {
	c := command{name: fmt.Sprintf("exec(%v, %+v)", name, args)}

	cmd := exec.Command(name, args...)

	// pipe previous pipe to stdin if available.
	if s.Command != nil {
		cmd.Stdin = s.Command
	}

	// pipe stdout and stderr to the new pipe.
	cmdOut, err := cmd.StdoutPipe()
	c.appendError(err, "pipe stdout")
	c.Reader = cmdOut

	if errWriter == nil {
		errWriter = ioutil.Discard
	}
	cmd.Stderr = errWriter

	// start the process
	err = cmd.Start()
	c.appendError(err, "start process")

	c.Closer = closerFn(func() error { return cmd.Wait() })

	return s.PipeTo(c)
}

type closerFn func() error

func (f closerFn) Close() error { return f() }
