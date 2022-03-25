//go:build ignore

package main

import (
	"github.com/opslevel/cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	doc.GenMarkdownTree(cmd.GetRootCmd(), "../docs/")
}
