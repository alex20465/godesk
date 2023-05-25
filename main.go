package main

import (
	"github.com/alecthomas/kong"
)

var CLI struct {
	Connect ConnectCmd `cmd: "" help:"Establish the first initial connection / pairing to the desk device over bluetooth"`
	Up      UpCmd      `cmd help:"move the desk one step UP"`
	Down    DownCmd    `cmd help:"move the desk one step DOWN"`
}

func main() {
	ctx := kong.Parse(&CLI)
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
