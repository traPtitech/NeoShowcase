package main

import (
	"github.com/traPtitech/neoshowcase/pkg/cli"
	"log"
)

var (
	version  = "UNKNOWN"
	revision = "UNKNOWN"
)

func main() {
	example2()
	example()
}

func Main() {
	cli.Version = version
	cli.Revision = revision
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
