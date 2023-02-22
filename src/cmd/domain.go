package cmd

import (
    "github.com/creasty/defaults"
    "github.com/opslevel/opslevel-go/v2023"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var createDomainCmd = &cobra.Command{
    Use:   "domain",
    Short: "Create a domain",
    Long:  `Create a domain`,
    Example: `
opslevel create domain -f domain.yaml
cat << EOF | opslevel create domain -f -
name: "My Domain"
description: "Hello World Domain"
owner: "my-team-or-group"
note: |
  This is my domain
EOF
`,
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var getDomainCmd = &cobra.Command{
    Use:   "domain {ID|EMAIL}",
    Short: "Get details about a domain",
    Long:  `Get details about a domain`,
    Example: `
opslevel get domain my-domain | jq
`,
    Args:       cobra.ExactArgs(1),
    ArgAliases: []string{"ID"},
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var listDomainCmd = &cobra.Command{
    Use:     "domain",
    Aliases: []string{"domains"},
    Short:   "Lists domains",
    Long:    `Lists domains`,
    Example: `
opslevel list domain
opslevel list domain -o json | jq 'map({"key": .Name, "value": .Owner}) | from_entries
`,
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var updateDomainCmd = &cobra.Command{
    Use:   "domain {ID|ALIAS}",
    Short: "Update a domain",
    Long:  `Update a domain`,
    Example: `
cat << EOF | opslevel update domain my-domain -f -
owner: "my-other-group"
`,
    Args:       cobra.ExactArgs(1),
    ArgAliases: []string{"ID"},
    Run: func(cmd *cobra.Command, args []string) {
    },
}

var deleteDomainCmd = &cobra.Command{
    Use:   "domain {ID|ALIAS}",
    Short: "Delete a domain",
    Long:  `Delete a domaian`,
    Example: `
opslevel delete domain my-domain
`,
    Args:       cobra.ExactArgs(1),
    ArgAliases: []string{"ID"},
    Run: func(cmd *cobra.Command, args []string) {
    },
}

func init() {
    createCmd.AddCommand(createDomainCmd)
    getCmd.AddCommand(getDomainCmd)
    listCmd.AddCommand(listDomainCmd)
    updateCmd.AddCommand(updateDomainCmd)
    deleteCmd.AddCommand(deleteDomainCmd)
}

func readDomainInput() (*opslevel.UserInput, error) {
    readCreateConfigFile()
    evt := &opslevel.UserInput{}
    viper.Unmarshal(&evt)
    if err := defaults.Set(evt); err != nil {
        return nil, err
    }
    return evt, nil
}
