package script

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

// To writes the output of the stream to an io.Writer.
func (s Stream) To(w io.Writer) error {
	var errors *multierror.Error
	if _, err := io.Copy(w, s); err != nil {
		errors = multierror.Append(errors, err)
	}
	if err := s.Close(); err != nil {
		errors = multierror.Append(errors, err)
	}
	return errors.ErrorOrNil()
}

// ToStdout pipes the stdout of the stream to screen.
func (s Stream) ToStdout() error {
	return s.To(os.Stdout)
}

// ToString reads stdout of the stream and returns it as a string.
func (s Stream) ToString() (string, error) {
	var out bytes.Buffer
	err := s.To(&out)
	return out.String(), err
}

// ToFile dumps the output of the stream to a file.
func (s Stream) ToFile(path string) error {
	err := makeDir(path)
	if err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return s.To(f)
}

// AppendFile appends the output of the stream to a file.
func (s Stream) AppendFile(path string) error {
	err := makeDir(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		return s.ToFile(path)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	return s.To(f)
}

// ToTempFile dumps the output of the stream to a temporary file and returns the temporary files'
// path.
func (s Stream) ToTempFile() (path string, err error) {
	f, err := ioutil.TempFile("", "script-")
	if err != nil {
		return "", err
	}
	defer f.Close()

	return f.Name(), s.To(f)
}

// Discard executes the stream pipeline but discards the output.
func (s Stream) Discard() error {
	return s.To(ioutil.Discard)
}

func makeDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0775)
}
