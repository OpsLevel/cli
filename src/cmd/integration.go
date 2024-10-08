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
	IntegrationTypeAWS   IntegrationType = "aws"
	IntegrationTypeAzure IntegrationType = "azure"
	IntegrationTypeGCP   IntegrationType = "googleCloud"
)

var AllIntegrationType = []IntegrationType{IntegrationTypeAWS, IntegrationTypeAzure, IntegrationTypeGCP}

var IntegrationConfigCurrentVersion = "1"

type IntegrationInputType struct {
	Version string `yaml:"version"`
	Kind    IntegrationType
	Spec    map[string]interface{}
}

type IntegrationInput interface {
	opslevel.AWSIntegrationInput | opslevel.AzureResourcesIntegrationInput | opslevel.GoogleCloudIntegrationInput
}

func validateIntegrationInput() (*IntegrationInputType, error) {
	input, err := readResourceInput[IntegrationInputType]()
	if err != nil {
		return nil, err
	}
	if input.Version != CheckConfigCurrentVersion {
		return nil, fmt.Errorf("supported config version is '%s' but found '%s'",
			IntegrationConfigCurrentVersion, input.Version)
	}
	switch input.Kind {
	case IntegrationTypeAWS, IntegrationTypeAzure, IntegrationTypeGCP:
		return input, nil
	default:
		return nil, fmt.Errorf("unsupported integration kind: '%s' (must be one of: %+v)",
			input.Kind, AllIntegrationType)
	}
}

func readIntegrationInput[T IntegrationInput](input *IntegrationInputType) (T, error) {
	var output T
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
  awsTagsOverrideOwnership: true
  ownershipTagKeys: ["owner","service","app"]
EOF

cat << EOF | opslevel create integration -f -
version: 1
kind: azure
spec:
  name: "Azure New"
  tenantId: "12345678-1234-1234-1234-123456789abc"
  subscriptionId: "12345678-1234-1234-1234-123456789def"
  clientId: "XXX_CLIENT_ID_XXX"
  clientSecret: "XXX_CLIENT_SECRET_XXX"
EOF

cat << EOF | opslevel create integration -f -
version: 1
kind: googleCloud
spec:
  name: "GCP New"
  ownershipTagKeys:
	- owner
    - team
  privateKey: "XXX_PRIVATE_KEY_XXX"
  clientEmail: "service-account-123@appspot.gserviceaccount.com"
  tagsOverrideOwnership: false
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, validateErr := validateIntegrationInput()
		cobra.CheckErr(validateErr)

		var result *opslevel.Integration
		switch input.Kind {
		case IntegrationTypeAWS:
			awsInput, err := readIntegrationInput[opslevel.AWSIntegrationInput](input)
			cobra.CheckErr(err)
			result, err = getClientGQL().CreateIntegrationAWS(awsInput)
			cobra.CheckErr(err)
		case IntegrationTypeAzure:
			azureInput, err := readIntegrationInput[opslevel.AzureResourcesIntegrationInput](input)
			cobra.CheckErr(err)
			result, err = getClientGQL().CreateIntegrationAzureResources(azureInput)
			cobra.CheckErr(err)
		case IntegrationTypeGCP:
			gcpInput, err := readIntegrationInput[opslevel.GoogleCloudIntegrationInput](input)
			cobra.CheckErr(err)
			result, err = getClientGQL().CreateIntegrationGCP(gcpInput)
			cobra.CheckErr(err)
		default:
			cobra.CheckErr(fmt.Errorf("cannot use unexpected input kind: '%s'", input.Kind))
		}

		fmt.Printf("Created %s integration '%s' with id '%s'\n", input.Kind, result.Name, result.Id)
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
	Example: `cat << EOF | opslevel update integration Z2lkOi8vb123456789 -f -
version: 1
kind: aws
spec:
  awsTagsOverrideOwnership: true
  ownershipTagKeys: ["owner","service","app"]
EOF

cat << EOF | opslevel update integration Z2lkOi8vb123456789 -f -
version: 1
kind: azure
spec:
  name: "Azure Dev"
  clientId: "XXX_CLIENT_ID_XXX"
  clientSecret: "XXX_CLIENT_SECRET_XXX"
EOF

cat << EOF | opslevel update integration Z2lkOi8vb123456789 -f -
version: 1
kind: googleCloud
spec:
  name: "GCP Dev"
  ownershipTagKeys:
	- opslevel_team
	- team
  privateKey: "XXX_NEW_PRIVATE_KEY_XXX"
  tagsOverrideOwnership: true
EOF
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		input, validateErr := validateIntegrationInput()
		cobra.CheckErr(validateErr)

		var result *opslevel.Integration
		switch input.Kind {
		case IntegrationTypeAWS:
			awsInput, err := readIntegrationInput[opslevel.AWSIntegrationInput](input)
			cobra.CheckErr(err)
			result, err = getClientGQL().UpdateIntegrationAWS(args[0], awsInput)
			cobra.CheckErr(err)
		case IntegrationTypeAzure:
			azureInput, err := readIntegrationInput[opslevel.AzureResourcesIntegrationInput](input)
			cobra.CheckErr(err)
			result, err = getClientGQL().UpdateIntegrationAzureResources(args[0], azureInput)
			cobra.CheckErr(err)
		case IntegrationTypeGCP:
			gcpInput, err := readIntegrationInput[opslevel.GoogleCloudIntegrationInput](input)
			cobra.CheckErr(err)
			result, err = getClientGQL().UpdateIntegrationGCP(args[0], gcpInput)
			cobra.CheckErr(err)
		default:
			cobra.CheckErr(fmt.Errorf("cannot use unexpected input kind: '%s'", input.Kind))
		}

		fmt.Printf("Updated %s integration '%s' with id '%s'\n", input.Kind, result.Name, result.Id)
	},
}

var reactivateIntegrationCmd = &cobra.Command{
	Use:        "reactivate ID",
	Short:      "Reactivate an integration",
	Long:       `Reactivate an integration that was invalidated or deactivated`,
	Example:    `opslevel update integration reactivate Z2lkOi8vb123456789`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		integration, err := getClientGQL().IntegrationReactivate(args[0])
		cobra.CheckErr(err)

		fmt.Printf("reactivated integration with id '%s'\n", integration.Id)
	},
}

func init() {
	createCmd.AddCommand(createIntegrationCmd)
	getCmd.AddCommand(getIntegrationCmd)
	listCmd.AddCommand(listIntegrationCmd)
	updateCmd.AddCommand(updateIntegrationCmd)

	updateIntegrationCmd.AddCommand(reactivateIntegrationCmd)
}
