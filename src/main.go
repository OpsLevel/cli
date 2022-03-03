package main

import (
	"github.com/opslevel/cli/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.Execute(version, commit)
}
