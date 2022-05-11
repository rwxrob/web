// Copyright 2022 web Robert Muhlestein
// SPDX-License-Identifier: Apache-2.0

// Package web provides the Bonzai command branch of the same name.
package web

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

// main branch
var Cmd = &Z.Cmd{

	Name:      `web`,
	Summary:   `common web requests`,
	Version:   `v0.3.0`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Source:    `git@github.com:rwxrob/web.git`,
	Issues:    `github.com/rwxrob/web/issues`,

	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, vars.Cmd, // common
		get, // post, put, del|delete, patch, dl|download
	},

	Description: `
		The {{cmd .Name}} command contains common web requests and
		strives for more command line usability than {{pkg "net/http"}} or
		even {{exe "curl"}} or {{exe "w3m"}}. In particular, the interface
		design is purposefully simple and stateful. The high-level {{pre
		"pkg"}} library can be used independently from the {{cmd .Name}}
		composable command.`,
}

var get = &Z.Cmd{

	Name:    `get`,
	Summary: `submit http get request`,
	MinArgs: 1,
	MaxArgs: 2,

	Call: func(_ *Z.Cmd, args ...string) error {
		req := Req{U: args[0], D: ""}
		if err := req.Submit(); err != nil {
			return err
		}
		fmt.Println(req.D)
		return nil
	},
}
