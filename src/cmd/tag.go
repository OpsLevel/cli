package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2023"
	"golang.org/x/exp/slices"

	"github.com/spf13/cobra"
)

var resourceType string

var createTagCmd = &cobra.Command{
	Use:   "tag --type=RESOURCE_TYPE [--assign] RESOURCE_ID KEY VALUE",
	Short: "Create/assign a tag",
	Long:  "Create/assign a tag",
	Example: `
opslevel create tag --type=Team ID|ALIAS KEY VALUE
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
				input.Id = opslevel.ID(resource)
			} else {
				input.Alias = resource
				input.Type = opslevel.TaggableResource(resourceType)
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
				input.Id = opslevel.ID(resource)
			} else {
				input.Alias = resource
				input.Type = opslevel.TaggableResource(resourceType)
			}

			result, err := getClientGQL().CreateTag(input)
			cobra.CheckErr(err)
			fmt.Printf("created tag for %s: %s\n", resource, result.Id)
		}
	},
}

var getObjectTagCmd = &cobra.Command{
	Use:   "tag --type=RESOURCE_TYPE RESOURCE_ID KEY",
	Short: "Get tags on an object matching key",
	Long:  "Get tags on an object matching key",
	Example: `
opslevel get tag --type=Service ID|ALIAS KEY | jq
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

		output := []opslevel.Tag{}
		for _, tag := range tags.Nodes {
			if tagKey == tag.Key {
				output = append(output, tag)
			}
		}

		common.PrettyPrint(output)
	},
}

var listObjectTagCmd = &cobra.Command{
	Use:     "tag --type=RESOURCE_TYPE RESOURCE_ID",
	Aliases: []string{"tags"},
	Short:   "Get all tags on an object",
	Long:    "Get all tags on an object",
	Example: `
opslevel list tag --type=Service ID|ALIAS
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
opslevel update tag TAG_ID KEY VALUE
`,
	Args:       cobra.ExactArgs(3),
	ArgAliases: []string{"TAG_ID", "KEY", "VALUE"},
	Run: func(cmd *cobra.Command, args []string) {
		tag := args[0]
		key := args[1]
		value := args[2]

		input := opslevel.TagUpdateInput{
			Id:    opslevel.ID(tag),
			Key:   key,
			Value: value,
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
opslevel delete tag TAG_ID
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
	createCmd.AddCommand(createTagCmd)
	createTagCmd.Flags().StringVar(&resourceType, "type", "", "resource type")
	createTagCmd.Flags().Bool("assign", false, "assign a tag instead of creating it")

	getCmd.AddCommand(getObjectTagCmd)
	getObjectTagCmd.Flags().StringVar(&resourceType, "type", "", "resource type")

	listCmd.AddCommand(listObjectTagCmd)
	listObjectTagCmd.Flags().StringVar(&resourceType, "type", "", "resource type")

	updateCmd.AddCommand(updateTagCmd)

	deleteCmd.AddCommand(deleteTagCmd)
}

func validateResourceTypeArg() error {
	if resourceType == "" {
		return errors.New("must specify a taggable resource type using --type=RESOURCE_TYPE")
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
		return errors.New("not a taggable resource type: " + resourceType)
	}
	resourceType = lowercaseMap[lowercaseInput]

	return nil
}
