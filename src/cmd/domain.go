package cmd

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Create a domain",
	Long: `Create a domain

cat << EOF | opslevel create domain -f -
name: "My Domain"
description: "Hello World Domain"
ownerId: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NjY"
note: "Additional details"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readDomainCreateInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateDomain(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

func init() {
	createCmd.AddCommand(createDomainCmd)
}

func readDomainCreateInput() (*opslevel.DomainCreateInput, error) {
	readCreateConfigFile()
	evt := &opslevel.DomainCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
