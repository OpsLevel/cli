package cmd

import (
	"errors"
	"fmt"

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
		if resourceType == "" {
			err := errors.New("must specify a taggable resource type using --type=RESOURCE_TYPE")
			cobra.CheckErr(err)
		}
		if !slices.Contains(opslevel.AllTaggableResource, resourceType) {
			err := errors.New("not a taggable resource type: " + resourceType)
			cobra.CheckErr(err)
		}

		resource := args[0]
		key := args[1]
		value := args[2]

		if cmd.Flag("assign").Changed {
			var input opslevel.TagAssignInput
			tagInput := opslevel.TagInput{
				Key:   key,
				Value: value,
			}

			if opslevel.IsID(resource) {
				input = opslevel.TagAssignInput{
					Id:   opslevel.ID(resource),
					Tags: []opslevel.TagInput{tagInput},
				}
			} else {
				input = opslevel.TagAssignInput{
					Alias: resource,
					Type:  opslevel.TaggableResource(resourceType),
					Tags:  []opslevel.TagInput{tagInput},
				}
			}

			result, err := getClientGQL().AssignTag(input)
			cobra.CheckErr(err)
			fmt.Printf("updated new tag on %s: %s\n", resource, result[0].Id)
		} else {
			var input opslevel.TagCreateInput
			if opslevel.IsID(resource) {
				input = opslevel.TagCreateInput{
					Id:    opslevel.ID(resource),
					Key:   key,
					Value: value,
				}
			} else {
				input = opslevel.TagCreateInput{
					Alias: resource,
					Type:  opslevel.TaggableResource(resourceType),
					Key:   key,
					Value: value,
				}
			}

			result, err := getClientGQL().CreateTag(input)
			cobra.CheckErr(err)
			fmt.Printf("added new tag on %s: %s\n", resource, result.Id)
		}
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

	deleteCmd.AddCommand(deleteTagCmd)
}
