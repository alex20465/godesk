package main

import (
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

var CLI struct {
	Connect ConnectCmd `cmd: "" help:"Establish the first initial connection / pairing to the desk device over bluetooth"`
	Up      UpCmd      `cmd help:"move the desk one step UP"`
	Down    DownCmd    `cmd help:"move the desk one step DOWN"`
	Goto    GotoCmd    `cmd help:"move the desk to a specific position"`
	Status  StatusCmd  `cmd help:"current desk status"`
}

func main() {
	log.SetLevel(log.ErrorLevel)

	ctx := kong.Parse(&CLI)
	err := ctx.Run(&Context{})
	ctx.FatalIfErrorf(err)
}
