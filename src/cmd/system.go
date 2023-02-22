package cmd

import (
    "github.com/creasty/defaults"
    "github.com/opslevel/opslevel-go/v2023"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var createSystemCmd = &cobra.Command{
    Use:   "system",
    Short: "Create a system",
    Long:  `Create a system`,
    Example: `
opslevel create system -f system.yaml
cat << EOF | opslevel create system -f -
name: "My System"
description: "Hello World System"
owner: "my-team-or-group"
parent: "my-domain"
note: |
  This is my system
EOF
`,
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var getSystemCmd = &cobra.Command{
    Use:   "system {ID|EMAIL}",
    Short: "Get details about a system",
    Long:  `Get details about a system`,
    Example: `
opslevel get system my-system | jq
`,
    Args:       cobra.ExactArgs(1),
    ArgAliases: []string{"ID"},
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var listSystemCmd = &cobra.Command{
    Use:     "system",
    Aliases: []string{"systems"},
    Short:   "Lists systems",
    Long:    `Lists systems`,
    Example: `
opslevel list system
opslevel list system -o json | jq 'map({"key": .Name, "value": .Owner}) | from_entries
`,
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var updateSystemCmd = &cobra.Command{
    Use:   "system {ID|ALIAS}",
    Short: "Update a system",
    Long:  `Update a system`,
    Example: `
cat << EOF | opslevel update system my-system -f -
owner: "my-other-group"
parent: "my-other-domain"
`,
    Args:       cobra.ExactArgs(1),
    ArgAliases: []string{"ID"},
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var deleteSystemCmd = &cobra.Command{
    Use:   "system {ID|ALIAS}",
    Short: "Delete a system",
    Long:  `Delete a domaian`,
    Example: `
opslevel delete system my-system
`,
    Args:       cobra.ExactArgs(1),
    ArgAliases: []string{"ID"},
    Run: func(cmd *cobra.Command, args []string) {
    },
}

func init() {
    createCmd.AddCommand(createSystemCmd)
    getCmd.AddCommand(getSystemCmd)
    listCmd.AddCommand(listSystemCmd)
    updateCmd.AddCommand(updateSystemCmd)
    deleteCmd.AddCommand(deleteSystemCmd)
}

func readSystemInput() (*opslevel.UserInput, error) {
    readCreateConfigFile()
    evt := &opslevel.UserInput{}
    viper.Unmarshal(&evt)
    if err := defaults.Set(evt); err != nil {
        return nil, err
    }
    return evt, nil
}
