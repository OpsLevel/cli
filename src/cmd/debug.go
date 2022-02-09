package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use: "debug",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := getClientGQL().ListCategories()

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ .Name | cyan }} ({{ .Index | red }})", promptui.IconSelect),
			Inactive: "    {{ .Name | cyan }} ({{ .Index | red }})",
			Selected: fmt.Sprintf("%s {{ .Name | faint }}", promptui.IconGood),
		}

		prompt := promptui.Select{
			Label:     "Select Day",
			Items:     list,
			Templates: templates,
			Size:      len(list),
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		fmt.Printf("You choose %q\n", result)
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
