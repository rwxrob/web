# 🌳 Go Bonzai™ Common Web Requests

[![GoDoc](https://godoc.org/github.com/rwxrob/web?status.svg)](https://godoc.org/github.com/rwxrob/web)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

This `web` Bonzai branch contains common web requests and strives for
more command line usability than the raw `net/http` package or even
`curl` or `w3m`. In particular, the interface design is purposefully
simple and stateful. The high-level `pkg` can be used separate
from the `web` composable command.

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/rwxrob/web/web@latest
```

Composed

```go
package z

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/web"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, web.Cmd},
}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C web web
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.
