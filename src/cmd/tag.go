package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

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
		err := validateResourceTypeArg(resourceType)
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
			fmt.Printf("updated new tag on %s: %s\n", resource, result[0].Id)
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
			fmt.Printf("added new tag on %s: %s\n", resource, result.Id)
		}
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

	updateCmd.AddCommand(updateTagCmd)

	deleteCmd.AddCommand(deleteTagCmd)
}

func validateResourceTypeArg(resourceType string) error {
	if resourceType == "" {
		return errors.New("must specify a taggable resource type using --type=RESOURCE_TYPE")
	}
	if !slices.Contains(opslevel.AllTaggableResource, resourceType) {
		return errors.New("not a taggable resource type: " + resourceType)
	}

	return nil
}
