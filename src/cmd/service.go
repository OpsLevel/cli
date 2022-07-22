package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2022"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Create a service",
	Long:  `Create a service`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readServiceCreateInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateService(*input)
		cobra.CheckErr(err)
		common.PrettyPrint(result.Id)
	},
}

var createServiceTagCmd = &cobra.Command{
	Use:   "tag ID|ALIAS TAG_KEY TAG_VALUE",
	Short: "Create a service tag",
	Long: `Create a service tag
	
opslevel create service tag my-service foo bar
opslevel create service tag --assign my-service foo bar
`,
	Args:       cobra.ExactArgs(3),
	ArgAliases: []string{"ID", "ALIAS", "TAG_KEY", "TAG_VALUE"},
	Run: func(cmd *cobra.Command, args []string) {
		var result interface{}
		var err error
		serviceKey := args[0]
		tagKey := args[1]
		tagValue := args[2]
		tagAssign, err := cmd.Flags().GetBool("assign")
		cobra.CheckErr(err)
		if tagAssign {
			input := opslevel.TagAssignInput{
				Tags: []opslevel.TagInput{
					{Key: tagKey, Value: tagValue},
				},
			}
			if common.IsID(serviceKey) {
				input.Id = serviceKey
			} else {
				input.Alias = serviceKey
			}
			input.Type = opslevel.TaggableResourceService
			result, err = getClientGQL().AssignTags(input)
		} else {
			input := opslevel.TagCreateInput{
				Key:   tagKey,
				Value: tagValue,
			}
			if common.IsID(serviceKey) {
				input.Id = serviceKey
			} else {
				input.Alias = serviceKey
			}
			input.Type = opslevel.TaggableResourceService
			result, err = getClientGQL().CreateTag(input)
		}
		cobra.CheckErr(err)
		common.PrettyPrint(result)
	},
}

var getServiceCmd = &cobra.Command{
	Use:        "service ID|ALIAS",
	Short:      "Get details about a service",
	Long:       `Get details about a service`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var result *opslevel.Service
		var err error
		if common.IsID(key) {
			result, err = getClientGQL().GetService(key)
			cobra.CheckErr(err)
		} else {
			result, err = getClientGQL().GetServiceWithAlias(key)
			cobra.CheckErr(err)
		}
		cobra.CheckErr(err)
		common.WasFound(result.Id == nil, key)
		common.PrettyPrint(result)
	},
}

var getServiceTagCmd = &cobra.Command{
	Use:   "tag ID|ALIAS TAG_KEY",
	Short: "Get a service's tag",
	Long: `Get a service's' tag

opslevel get service tag my-service | jq 'from_entries'
opslevel get service tag my-service my-tag
`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"ID", "ALIAS", "TAG_KEY"},
	Run: func(cmd *cobra.Command, args []string) {
		serviceKey := args[0]
		singleTag := len(args) == 2
		var tagKey string
		if singleTag {
			tagKey = args[1]
		}

		var result *opslevel.Service
		var err error
		if common.IsID(serviceKey) {
			result, err = getClientGQL().GetService(serviceKey)
			cobra.CheckErr(err)
		} else {
			result, err = getClientGQL().GetServiceWithAlias(serviceKey)
			cobra.CheckErr(err)
		}
		if result.Id == nil {
			cobra.CheckErr(fmt.Errorf("service '%s' not found", serviceKey))
		}
		output := []opslevel.Tag{}
		for _, tag := range result.Tags.Nodes {
			if singleTag == false || tagKey == tag.Key {
				output = append(output, tag)
			}
		}
		if len(output) == 0 {
			cobra.CheckErr(fmt.Errorf("tag with key '%s' not found on service '%s'", tagKey, serviceKey))
		}
		common.PrettyPrint(output)
	},
}

var listServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Lists services",
	Long:    `Lists services`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := getClientGQL().ListServices()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"NAME", "ID", "ALIASES"})
			for _, item := range list {
				w.Write([]string{item.Name, item.Id.(string), strings.Join(item.Aliases, "/")})
			}
			w.Flush()
		} else {
			w := common.NewTabWriter("NAME", "ID", "ALIASES")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Id, strings.Join(item.Aliases, ","))
			}
			w.Flush()
		}
	},
}

var updateServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Update a service",
	Long:  `Update a service`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readServiceUpdateInput()
		cobra.CheckErr(err)
		service, err := getClientGQL().UpdateService(*input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(service, "", "    "))
	},
}

var deleteServiceCmd = &cobra.Command{
	Use:        "service ID|ALIAS",
	Short:      "Delete a service",
	Long:       `Delete a service`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var err error
		if common.IsID(key) {
			err = getClientGQL().DeleteService(opslevel.ServiceDeleteInput{
				Id: key,
			})
			cobra.CheckErr(err)
		} else {
			err = getClientGQL().DeleteServiceWithAlias(key)
			cobra.CheckErr(err)
		}
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' service\n", key)
	},
}

var deleteServiceTagCmd = &cobra.Command{
	Use:        "tag ID|ALIAS TAG_KEY",
	Short:      "Delete a service's tag",
	Long:       `Delete a service's tag'`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"ID", "ALIAS", "TAG_KEY"},
	Run: func(cmd *cobra.Command, args []string) {
		serviceKey := args[0]
		tagKey := args[1]
		var result *opslevel.Service
		var err error
		if common.IsID(serviceKey) {
			result, err = getClientGQL().GetService(serviceKey)
			cobra.CheckErr(err)
		} else {
			result, err = getClientGQL().GetServiceWithAlias(serviceKey)
			cobra.CheckErr(err)
		}
		if result.Id == nil {
			cobra.CheckErr(fmt.Errorf("service '%s' not found", serviceKey))
		}
		for _, tag := range result.Tags.Nodes {
			if tagKey == tag.Key {
				getClientGQL().DeleteTag(tag.Id)
				fmt.Println("Deleted Tag")
				common.PrettyPrint(tag)
			}
		}
	},
}

var importServicesCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Imports services from a CSV",
	Long: `Imports a list of services from a CSV file with the column headers

    Name,Description,Product,Language,Framework,Tier,Lifecycle,Owner
`,
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := readImportFilepathAsCSV()
		client := getClientGQL()
		opslevel.Cache.CacheLifecycles(client)
		opslevel.Cache.CacheTiers(client)
		opslevel.Cache.CacheTeams(client)
		cobra.CheckErr(err)
		for reader.Rows() {
			name := reader.Text("Name")
			input := opslevel.ServiceCreateInput{
				Name:        name,
				Description: reader.Text("Description"),
				Product:     reader.Text("Product"),
				Language:    reader.Text("Language"),
				Framework:   reader.Text("Framework"),
			}
			tier := reader.Text("Tier")
			if tier != "" {
				if item, ok := opslevel.Cache.Tiers[tier]; ok {
					input.Tier = item.Alias
				}
			}
			lifecycle := reader.Text("Lifecycle")
			if tier != "" {
				if item, ok := opslevel.Cache.Lifecycles[lifecycle]; ok {
					input.Lifecycle = item.Alias
				}
			}
			owner := reader.Text("Owner")
			if tier != "" {
				if item, ok := opslevel.Cache.Teams[owner]; ok {
					input.Owner = item.Alias
				}
			}
			service, err := getClientGQL().CreateService(input)
			if err != nil {
				log.Error().Err(err).Msgf("error creating service '%s'", name)
				continue
			}
			log.Info().Msgf("created service '%s' with id '%s'\n", service.Name, service.Id)
		}
	},
}

func init() {
	createCmd.AddCommand(createServiceCmd)
	getCmd.AddCommand(getServiceCmd)
	listCmd.AddCommand(listServiceCmd)
	updateCmd.AddCommand(updateServiceCmd)
	deleteCmd.AddCommand(deleteServiceCmd)

	createServiceCmd.AddCommand(createServiceTagCmd)
	getServiceCmd.AddCommand(getServiceTagCmd)
	deleteServiceCmd.AddCommand(deleteServiceTagCmd)

	importCmd.AddCommand(importServicesCmd)

	createServiceTagCmd.Flags().Bool("assign", false, "Use the `tagAssign` mutation instead of `tagCreate`")
}

func readServiceCreateInput() (*opslevel.ServiceCreateInput, error) {
	readCreateConfigFile()
	evt := &opslevel.ServiceCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}

func readServiceUpdateInput() (*opslevel.ServiceUpdateInput, error) {
	readUpdateConfigFile()
	evt := &opslevel.ServiceUpdateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
