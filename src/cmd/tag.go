package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2025"
	"golang.org/x/exp/slices"

	"github.com/spf13/cobra"
)

var resourceType string

var exampleTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Example tag to assign to a resource",
	Long:  `Example tag to assign to a resource`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample(opslevel.TagInput{
			Key:   "example_key",
			Value: "example_value",
		}))
	},
}

var createTagCmd = &cobra.Command{
	Use:   "tag --type=RESOURCE_TYPE [--assign] ID|ALIAS KEY VALUE",
	Short: "Create/assign a tag",
	Long:  "Create/assign a tag",
	Example: `
opslevel create tag --type=Service my-service-alias foo bar
opslevel create tag --type=Team my-team-alias foo bar
opslevel create tag --type=Infra my-infra-alias foo bar
`,
	Args:       cobra.ExactArgs(3),
	ArgAliases: []string{"RESOURCE_ID", "KEY", "VALUE"},
	Run: func(cmd *cobra.Command, args []string) {
		err := validateResourceTypeArg()
		cobra.CheckErr(err)

		resource := args[0]
		key := args[1]
		value := args[2]

		if cmd.Flag("assign").Changed {
			tagInput := opslevel.TagInput{
				Key:   key,
				Value: value,
			}
			input := opslevel.TagAssignInput{Tags: []opslevel.TagInput{tagInput}}

			if opslevel.IsID(resource) {
				input.Id = opslevel.RefOf(opslevel.ID(resource))
			} else {
				input.Alias = opslevel.RefOf(resource)
				resourceType := opslevel.TaggableResource(resourceType)
				input.Type = &resourceType
			}

			result, err := getClientGQL().AssignTag(input)
			cobra.CheckErr(err)
			fmt.Printf("assigned tag for %s: %s\n", resource, result[0].Id)
		} else {
			input := opslevel.TagCreateInput{
				Key:   key,
				Value: value,
			}
			if opslevel.IsID(resource) {
				input.Id = opslevel.NewID(resource)
			} else {
				input.Alias = &resource
				resourceType := opslevel.TaggableResource(resourceType)
				input.Type = &resourceType
			}

			result, err := getClientGQL().CreateTag(input)
			cobra.CheckErr(err)
			fmt.Printf("created tag for %s: %s\n", resource, result.Id)
		}
	},
}

var getObjectTagCmd = &cobra.Command{
	Use:   "tag --type=RESOURCE_TYPE ID|ALIAS KEY",
	Short: "Get tags on an object matching key",
	Long:  "Get tags on an object matching key",
	Example: `
opslevel get tag --type=Service my-service-alias foo
`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"RESOURCE_ID", "KEY"},
	Run: func(cmd *cobra.Command, args []string) {
		err := validateResourceTypeArg()
		cobra.CheckErr(err)

		resourceIdentifier := args[0]
		tagKey := args[1]

		client := getClientGQL()
		result, err := client.GetTaggableResource(opslevel.TaggableResource(resourceType), resourceIdentifier)
		cobra.CheckErr(err)

		tags, err := result.GetTags(client, nil)
		cobra.CheckErr(err)

		output := make([]opslevel.Tag, 0)
		for _, tag := range tags.Nodes {
			if tagKey == tag.Key {
				output = append(output, tag)
			}
		}

		common.PrettyPrint(output)
	},
}

var listObjectTagCmd = &cobra.Command{
	Use:     "tag --type=RESOURCE_TYPE ID|ALIAS",
	Aliases: []string{"tags"},
	Short:   "Get all tags on an object",
	Long:    "Get all tags on an object",
	Example: `
opslevel list tag --type=Service my-service-alias -o json | jq
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"RESOURCE_ID"},
	Run: func(cmd *cobra.Command, args []string) {
		err := validateResourceTypeArg()
		cobra.CheckErr(err)

		resourceIdentifier := args[0]

		client := getClientGQL()
		result, err := client.GetTaggableResource(opslevel.TaggableResource(resourceType), resourceIdentifier)
		cobra.CheckErr(err)

		tags, err := result.GetTags(client, nil)
		cobra.CheckErr(err)

		common.PrettyPrint(tags.Nodes)
	},
}

var updateTagCmd = &cobra.Command{
	Use:   "tag TAG_ID KEY VALUE",
	Short: "Update a tag",
	Long:  "Update a tag",
	Example: `
opslevel update tag XXX_TAG_ID_XXX foo baz
`,
	Args:       cobra.ExactArgs(3),
	ArgAliases: []string{"TAG_ID", "KEY", "VALUE"},
	Run: func(cmd *cobra.Command, args []string) {
		tag := args[0]
		key := args[1]
		value := args[2]

		input := opslevel.TagUpdateInput{
			Id:    opslevel.ID(tag),
			Key:   &key,
			Value: &value,
		}

		result, err := getClientGQL().UpdateTag(input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(result, "", "    "))
	},
}

var deleteTagCmd = &cobra.Command{
	Use:   "tag TAG_ID",
	Short: "Delete a tag",
	Long:  "Delete a tag",
	Example: `
opslevel delete tag XXX_TAG_ID_XXX
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"TAG_ID"},
	Run: func(cmd *cobra.Command, args []string) {
		tag := opslevel.ID(args[0])

		err := getClientGQL().DeleteTag(tag)
		cobra.CheckErr(err)
		fmt.Printf("deleted a tag: %s\n", tag)
	},
}

func init() {
	exampleCmd.AddCommand(exampleTagCmd)
	createCmd.AddCommand(createTagCmd)
	createTagCmd.Flags().StringVar(&resourceType, "type", "", "resource type")
	createTagCmd.Flags().Bool("assign", false, "if a tag with the same key already exists it will be updated, otherwise a new tag will be created.")

	getCmd.AddCommand(getObjectTagCmd)
	getObjectTagCmd.Flags().StringVar(&resourceType, "type", "", "resource type")

	listCmd.AddCommand(listObjectTagCmd)
	listObjectTagCmd.Flags().StringVar(&resourceType, "type", "", "resource type")

	updateCmd.AddCommand(updateTagCmd)

	deleteCmd.AddCommand(deleteTagCmd)
}

func validateResourceTypeArg() error {
	if resourceType == "" {
		return fmt.Errorf("must specify a taggable resource type using --type=RESOURCE_TYPE")
	}

	// if ProperCase, continue
	if slices.Contains(opslevel.AllTaggableResource, resourceType) {
		return nil
	}

	// if lowercase, check if it exists. if not, error out
	lowercaseInput := strings.ToLower(resourceType)
	if lowercaseInput == "infra" {
		resourceType = string(opslevel.TaggableResourceInfrastructureresource)
		return nil
	}
	lowercaseMap := make(map[string]string)
	for _, s := range opslevel.AllTaggableResource {
		lowercaseMap[strings.ToLower(s)] = s
	}
	if lowercaseMap[lowercaseInput] == "" {
		return fmt.Errorf("not a taggable resource type: '%s'", resourceType)
	}
	resourceType = lowercaseMap[lowercaseInput]

	return nil
}
