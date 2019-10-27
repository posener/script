# script

[![Build Status](https://travis-ci.org/posener/script.svg?branch=master)](https://travis-ci.org/posener/script)
[![codecov](https://codecov.io/gh/posener/script/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/script)
[![GoDoc](https://godoc.org/github.com/posener/script?status.svg)](http://godoc.org/github.com/posener/script)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/script.svg)](https://goreadme.herokuapp.com)

Package script provides helper functions to write scripts.

Inspired by [https://github.com/bitfield/script](https://github.com/bitfield/script), with some small modifications:

* Output between streamed commands is a stream and not loaded to memory.

* Better representation of errors and stderr.

The script chain is represented by the
[`Stream`](https://godoc.org/github.com/posener/script#Stream) type. While each command in the
stream is abstracted by the [`Command`](https://godoc.org/github.com/posener/script#Command)
struct, which enable extending this library freely.

#### Examples

##### HelloWorld

A simple "hello world" example that creats a stream and pipe it to the stdout.

```golang
// Create an "hello world" stream and use the ToStdout method to write it to stdout.
Echo("hello world").ToStdout()
```

##### Iterate

An example that shows how to iterate scanned lines.

```golang
// Stream can be any stream, in this case we have echoed 3 lines.
stream := Echo("first\nsecond\nthird")

// To iterate over the stream lines, it is better not to read it into memory and split over the
// lines, but use the `bufio.Scanner`:
defer stream.Close()
scanner := bufio.NewScanner(stream)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}
```

##### PipeTo

An example that shows how to create custom commands using the `PipeTo` method with a `PipeFn`
function.

```golang
Echo("1\n2\n3").PipeTo(func(r io.Reader) Command {
    // Create a command that sums up all numbers in input.

    // Use buffered reader to read lines from input.
    buf := bufio.NewReader(r)

    // Store the sum of all numbers.
    sum := 0

    // Create a reade function (for the sake of example function, such that all code will fit
    // into the example function body - a proper more readable way to do it was to create a new
    // type with a state and a read function).
    read := func(b []byte) (int, error) {
        // Read next line from input.
        line, _, err := buf.ReadLine()

        // if EOF write sum to output.
        if err == io.EOF {
            return copy(b, append([]byte(strconv.Itoa(sum)), '\n')), nil
        }

        // Convert the line to a number and add it to the sum.
        if i, err := strconv.Atoi(string(line)); err == nil {
            sum += i
        }

        // We don't write anything to output, so we return 0 bytes with no error.
        return 0, nil
    }

    return Command{Name: "sum", Reader: readerFn(read)}
})
```


---

Created by [goreadme](https://github.com/apps/goreadme)
