package main

//go:generate go run gen.go

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
