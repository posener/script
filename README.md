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

### func [AppendFile](https://github.com/posener/script/blob/master/to.go#L82)

`func AppendFile(path string) (io.WriteCloser, error)`

### func [File](https://github.com/posener/script/blob/master/to.go#L74)

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


---

Created by [goreadme](https://github.com/apps/goreadme)
