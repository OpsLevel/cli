package main

import (
	"fmt"

	"github.com/opslevel/cli/cmd"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.Execute(fmt.Sprintf("%s-%0.12s", version, commit))
}
