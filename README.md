# script

[![GoDoc](https://godoc.org/github.com/posener/script?status.svg)](http://godoc.org/github.com/posener/script)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/script.svg)](https://goreadme.herokuapp.com)

Package script provides helper functions to write scripts.

Inspired by [https://github.com/bitfield/script](https://github.com/bitfield/script), with a small modifications:

* Output between streamed commands is a stream and not loaded to memory.

* Better representation of errors and stderr.

#### Examples

Example that shows how to write piped content to screen.

```golang
Echo("hello world").ToScreen()
```


---

Created by [goreadme](https://github.com/apps/goreadme)
