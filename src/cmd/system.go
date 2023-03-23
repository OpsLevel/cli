package cmd

import (
	"fmt"
	"github.com/creasty/defaults"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createSystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Create a system",
	Long: `Create a system
cat << EOF | opslevel create system -f -
name: "My System"
description: "Hello World System"
ownerId: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NjY"
parent:
	alias: "Name of parent domain"
note: "Additional system details"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readSystemCreateInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateSystem(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

func init() {
	createCmd.AddCommand(createSystemCmd)
}

func readSystemCreateInput() (*opslevel.SystemCreateInput, error) {
	readCreateConfigFile()
	evt := &opslevel.SystemCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
