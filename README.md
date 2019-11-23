# script

[![Build Status](https://travis-ci.org/posener/script.svg?branch=master)](https://travis-ci.org/posener/script)
[![codecov](https://codecov.io/gh/posener/script/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/script)
[![GoDoc](https://godoc.org/github.com/posener/script?status.svg)](http://godoc.org/github.com/posener/script)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/script.svg)](https://goreadme.herokuapp.com)

Package script provides helper functions to write scripts.

Inspired by [https://github.com/bitfield/script](https://github.com/bitfield/script), with some improvements:

* Output between streamed commands is a stream and not loaded to memory.

* Better representation and handling of errors.

* Proper incocation, usage and handling of stderr of custom commands.

The script chain is represented by a
[`Stream`](https://godoc.org/github.com/posener/script#Stream) object. While each command in the
stream is abstracted by the [`Command`](https://godoc.org/github.com/posener/script#Command)
struct. This library provides basic functionality, but can be extended freely.

## Functions

### func [AppendFile](https://github.com/posener/script/blob/master/to.go#L81)

`func AppendFile(path string) (io.WriteCloser, error)`

### func [File](https://github.com/posener/script/blob/master/to.go#L73)

`func File(path string) (io.WriteCloser, error)`

#### Examples

##### HelloWorld

A simple "hello world" example that creats a stream and pipe it to the stdout.

```golang
// Create an "hello world" stream and use the ToStdout method to write it to stdout.
Echo("hello world").ToStdout()
```

 Output:

```
hello world

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

 Output:

```
first
second
third

```

##### PipeTo

An example that shows how to create custom commands using the `PipeTo` method with a `PipeFn`
function.

```golang
Echo("1\n2\n3").PipeTo(func(r io.Reader) Command {
    // Create a command that sums up all numbers in input.
    //
    // In this example we create a reader function such that the whole code will fit into the
    // example function body. A more proper and readable way to do it was to create a new
    // type with a state that implements the `io.Reader` interface.

    // Use buffered reader to read lines from input.
    buf := bufio.NewReader(r)

    // Store the sum of all numbers.
    sum := 0

    // Read function reads the next line and adds it to the sum. If it gets and EOF error, it
    // writes the sum to the output and returns an EOF.
    read := func(b []byte) (int, error) {
        // Read next line from input.
        line, _, err := buf.ReadLine()

        // if EOF write sum to output.
        if err == io.EOF {
            return copy(b, append([]byte(strconv.Itoa(sum)), '\n')), io.EOF
        }
        if err != nil {
            return 0, err
        }

        // Convert the line to a number and add it to the sum.
        if i, err := strconv.Atoi(string(line)); err == nil {
            sum += i
        }

        // We don't write anything to output, so we return 0 bytes with no error.
        return 0, nil
    }

    return Command{Name: "sum", Reader: readerFn(read)}
}).ToStdout()
```

 Output:

```
6

```


---

Created by [goreadme](https://github.com/apps/goreadme)
