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

Example that shows how to write piped content to screen.

```golang
Echo("hello world").ToScreen()
```


---

Created by [goreadme](https://github.com/apps/goreadme)
