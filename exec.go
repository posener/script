package script

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

var stdin = Stream{Command: command{Reader: os.Stdin, name: "stdin"}}

// Exec executes a command and returns a stream of the stdout of the command.
func Exec(command string, args ...string) Stream {
	return stdin.exec(nil, command, args...)
}

// Exec executes a command, returns a stream of the stdout of the command and enable collecting the
// stderr of the command.
//
// If the errWriter is nil, it will be ignored.
//
// For example, collecting the stderr to memory can be done by providing a `&bytes.Buffer` as
// `errWriter`. Writing it to stderr can be done by providing `os.Stderr` as `errWriter`. Logging it
// to a file can be done by providing an `os.File` as the `errWriter`.
func ExecHandleStderr(errWriter io.Writer, cmd string, args ...string) Stream {
	return stdin.exec(errWriter, cmd, args...)
}

// Exec executes a command and returns a stream of the stdout of the command.
func (s Stream) Exec(cmd string, args ...string) Stream {
	return s.exec(nil, cmd, args...)
}

// Exec executes a command, returns a stream of the stdout of the command and enable collecting the
// stderr of the command.
//
// If the errWriter is nil, it will be ignored.
func (s Stream) ExecHandleStderr(errWriter io.Writer, cmd string, args ...string) Stream {
	return s.exec(errWriter, cmd, args...)
}

func (s Stream) exec(errWriter io.Writer, name string, args ...string) Stream {
	c := command{name: fmt.Sprintf("exec(%v, %+v)", name, args)}

	cmd := exec.Command(name, args...)

	// Pipe previous command output to stdin if available.
	if s.Command != nil {
		cmd.Stdin = s.Command
	}

	// Pipe stdout to the current command output.
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
