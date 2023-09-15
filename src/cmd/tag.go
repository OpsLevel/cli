package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2023"
	"golang.org/x/exp/slices"

	"github.com/spf13/cobra"
)

// TODO: combine TaggableResourceFetchFunction and TaggableResourceFetchAliasFunction once the ObjectGet functions have been
// harmonized to support both get by ID and get by Alias
type TaggableResourceFetchFunction func(id opslevel.ID) (any, error)

var TaggableResourceFetchFunctions = map[opslevel.TaggableResource]TaggableResourceFetchFunction{
	opslevel.TaggableResourceService:                func(id opslevel.ID) (any, error) { return getClientGQL().GetService(id) },
	opslevel.TaggableResourceRepository:             func(id opslevel.ID) (any, error) { return getClientGQL().GetRepository(id) },
	opslevel.TaggableResourceTeam:                   func(id opslevel.ID) (any, error) { return getClientGQL().GetTeam(id) },
	opslevel.TaggableResourceUser:                   func(id opslevel.ID) (any, error) { return getClientGQL().GetUser(string(id)) },
	opslevel.TaggableResourceDomain:                 func(id opslevel.ID) (any, error) { return getClientGQL().GetDomain(string(id)) },
	opslevel.TaggableResourceSystem:                 func(id opslevel.ID) (any, error) { return getClientGQL().GetSystem(string(id)) },
	opslevel.TaggableResourceInfrastructureresource: func(id opslevel.ID) (any, error) { return getClientGQL().GetInfrastructure(string(id)) },
}

type TaggableResourceFetchAliasFunction func(alias string) (any, error)

var TaggableResourceFetchAliasFunctions = map[opslevel.TaggableResource]TaggableResourceFetchAliasFunction{
	opslevel.TaggableResourceService:                func(alias string) (any, error) { return getClientGQL().GetServiceWithAlias(alias) },
	opslevel.TaggableResourceRepository:             func(alias string) (any, error) { return getClientGQL().GetRepositoryWithAlias(alias) },
	opslevel.TaggableResourceTeam:                   func(alias string) (any, error) { return getClientGQL().GetTeamWithAlias(alias) },
	opslevel.TaggableResourceUser:                   func(alias string) (any, error) { return getClientGQL().GetUser(alias) },
	opslevel.TaggableResourceDomain:                 func(alias string) (any, error) { return getClientGQL().GetDomain(alias) },
	opslevel.TaggableResourceSystem:                 func(alias string) (any, error) { return getClientGQL().GetSystem(alias) },
	opslevel.TaggableResourceInfrastructureresource: func(alias string) (any, error) { return getClientGQL().GetInfrastructure(alias) },
}

func GetTags(obj interface{}) (*opslevel.TagConnection, error) {
	// call Elem because input is a ptr
	v := reflect.ValueOf(obj).Elem()
	tagsField := v.FieldByName("Tags")

	if tagsField.IsValid() {
		return tagsField.Interface().(*opslevel.TagConnection), nil
	} else {
		return nil, errors.New("reflection error")
	}
}

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
	Use:   "tag --type=RESOURCE_TYPE RESOURCE_ID [KEY]",
	Short: "Get values of tags on objects",
	Long:  "Get values of tags on objects",
	Example: `
opslevel get tag --type=Service ID|ALIAS KEY                # search for values of a specific key
opslevel get tag --type=Team ID|ALIAS | jq 'from_entries'   # values of all keys
`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"RESOURCE_ID", "KEY"},
	Run: func(cmd *cobra.Command, args []string) {
		err := validateResourceTypeArg()
		cobra.CheckErr(err)

		resource := args[0]
		singleTag := len(args) == 2
		var tagKey string
		if singleTag {
			tagKey = args[1]
		}

		var result any
		if opslevel.IsID(resource) {
			id := opslevel.ID(resource)
			result, err = TaggableResourceFetchFunctions[opslevel.TaggableResource(resourceType)](id)
		} else {
			alias := args[0]
			result, err = TaggableResourceFetchAliasFunctions[opslevel.TaggableResource(resourceType)](alias)
		}

		tags, err := GetTags(result)
		cobra.CheckErr(err)

		var output []opslevel.Tag
		for _, tag := range tags.Nodes {
			if !singleTag || tagKey == tag.Key {
				output = append(output, tag)
			}
		}

		if len(output) == 0 {
			err := fmt.Errorf("tag with key '%s' not found on %s '%s'", tagKey, resourceType, resource)
			cobra.CheckErr(err)
		}
		common.PrettyPrint(output)
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
