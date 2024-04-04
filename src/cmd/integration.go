package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/spf13/cobra"
)

type IntegrationType string

const (
	IntegrationTypeAWS IntegrationType = "aws"
)

var IntegrationConfigCurrentVersion = "1"

type IntegrationInputType struct {
	Version string `yaml:"version"`
	Kind    IntegrationType
	Spec    map[string]interface{}
}

type IntegrationInput interface {
	opslevel.AWSIntegrationInput
}

func readIntegrationInput[T IntegrationInput]() (T, error) {
	var output T
	input, err := readResourceInput[IntegrationInputType]()
	if err != nil {
		return output, err
	}
	if input.Version != CheckConfigCurrentVersion {
		return output, fmt.Errorf("supported config version is '%s' but found '%s'",
			IntegrationConfigCurrentVersion, input.Version)
	}
	// TODO: need to use input.Kind and a switch statement - but currently we only support AWS
	if err := mapstructure.Decode(input.Spec, &output); err != nil {
		return output, err
	}
	return output, nil
}

var createIntegrationCmd = &cobra.Command{
	Use:     "integration",
	Aliases: []string{"integrations", "int", "ints"},
	Short:   "Create an integration",
	Long:    `Create an integration`,
	Example: `cat << EOF | opslevel create integration -f -
version: 1
kind: aws
spec:
  name: "Prod"
  iamRole: "arn:aws:iam::XXXXX:role/opslevel-integration"
  externalId: "XXXXXX"
  ownershipTagOverrides: true
  ownershipTagKeys: ["owner","service","app"]
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readIntegrationInput[opslevel.AWSIntegrationInput]()
		cobra.CheckErr(err)
		resp, err := getClientGQL().CreateIntegrationAWS(input)
		cobra.CheckErr(err)
		fmt.Printf("Created Integration '%s' with id '%s'\n", resp.Name, resp.Id)
	},
}

var getIntegrationCmd = &cobra.Command{
	Use:        "integration ID",
	Aliases:    []string{"int"},
	Short:      "Get details about a integration",
	Long:       `Get details about a integration`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		integration, err := getClientGQL().GetIntegration(opslevel.ID(key))
		cobra.CheckErr(err)
		common.PrettyPrint(integration)
	},
}

var listIntegrationCmd = &cobra.Command{
	Use:     "integration",
	Aliases: []string{"integrations", "int", "ints"},
	Short:   "Lists integrations",
	Long:    `Lists integrations`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListIntegrations(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "TYPE", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", item.Name, item.Type, item.Alias(), item.Id)
			}
			w.Flush()
		}
	},
}

var updateIntegrationCmd = &cobra.Command{
	Use:     "integration ID",
	Aliases: []string{"int"},
	Short:   "Update an integration",
	Long:    `Update an integration`,
	Example: `cat << EOF | opslevel update integration XXXXXXXX -f -
version: 1
kind: aws
spec:
  ownershipTagOverrides: true
  ownershipTagKeys: ["owner","service","app"]
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readIntegrationInput[opslevel.AWSIntegrationInput]()
		cobra.CheckErr(err)
		resp, err := getClientGQL().UpdateIntegrationAWS(args[0], input)
		cobra.CheckErr(err)
		fmt.Printf("Updated Integration '%s' with id '%s'\n", resp.Name, resp.Id)
	},
}

func init() {
	createCmd.AddCommand(createIntegrationCmd)
	getCmd.AddCommand(getIntegrationCmd)
	listCmd.AddCommand(listIntegrationCmd)
	updateCmd.AddCommand(updateIntegrationCmd)
}
