package cmd

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/opslevel/opslevel-go/v2023"

	"github.com/rs/zerolog/log"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Create a service",
	Long: `Create a service

cat << EOF | opslevel create service -f -
name: "hello world"
description: "Hello World Service"
product: "OSS"
language: "Go"
tier: "tier_4"
framework: "fasthttp"
owner:
  alias: "Platform"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readServiceCreateInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateService(*input)
		cobra.CheckErr(err)
		common.PrettyPrint(result.Id)
	},
}

var createServiceTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Create a service tag",
	Long:  `Create a service tag`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("This command is deprecated! Please use \nopslevel create tag <args>")
		cobra.CheckErr(err)
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
		client := getClientGQL()
		var service *opslevel.Service
		var err error
		if opslevel.IsID(key) {
			service, err = getClientGQL().GetService(opslevel.ID(key))
			cobra.CheckErr(err)
		} else {
			service, err = getClientGQL().GetServiceWithAlias(key)
			cobra.CheckErr(err)
		}
		_, err = service.GetDependents(client, nil)
		cobra.CheckErr(err)
		_, err = service.GetDependencies(client, nil)
		cobra.CheckErr(err)
		common.WasFound(service.Id == "", key)
		common.PrettyPrint(service)
	},
}

var getServiceTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Get a service's tag",
	Long:  `Get a service's tag`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("This command is deprecated! Please use \nopslevel get tag <args>")
		cobra.CheckErr(err)
	},
}

var listServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Lists services",
	Long:    `Lists services`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListServices(nil)
		list := resp.Nodes
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"NAME", "ID", "ALIASES"})
			for _, item := range list {
				w.Write([]string{item.Name, string(item.Id), strings.Join(item.Aliases, "/")})
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
	Long: `Update a service

cat << EOF | opslevel update service -f -
alias: "hello_world"
description: "Hello World Service Updated"
tier: "tier_3"
EOF`,
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
		if opslevel.IsID(key) {
			err = getClientGQL().DeleteService(opslevel.ServiceDeleteInput{
				Id: opslevel.ID(key),
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
	Use:   "tag",
	Short: "Delete a service's tag",
	Long:  `Delete a service's tag'`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("This command is deprecated! Please use \nopslevel delete tag <args>")
		cobra.CheckErr(err)
	},
}

var importServicesCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services"},
	Short:   "Imports services from a CSV",
	Long: `Imports a list of services from a CSV file with the column headers:
Name,Description,Product,Language,Framework,Tier,Lifecycle,Owner

Example:

cat << EOF | opslevel import services -f -
Name,Description,Product,Language,Framework,Tier,Lifecycle,Owner
Service A,,,Go,Cobra,tier_1,pre_alpha,
Service B,,,Python,Django,tier_3,beta,sales
Service C,Test,Home,,,,,platform
EOF
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
			if lifecycle != "" {
				if item, ok := opslevel.Cache.Lifecycles[lifecycle]; ok {
					input.Lifecycle = item.Alias
				}
			}
			owner := reader.Text("Owner")
			if owner != "" {
				if item, ok := opslevel.Cache.Teams[owner]; ok {
					input.Owner = opslevel.NewIdentifier(item.Alias)
				}
			}
			service, err := getClientGQL().CreateService(input)
			if err != nil {
				log.Error().Err(err).Msgf("error creating service '%s'", name)
				continue
			}
			log.Info().Msgf("created service '%s' with id '%s'", service.Name, service.Id)
		}
		reader.Close()
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
}

func readServiceCreateInput() (*opslevel.ServiceCreateInput, error) {
	readInputConfig()
	evt := &opslevel.ServiceCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}

func readServiceUpdateInput() (*opslevel.ServiceUpdateInput, error) {
	readInputConfig()
	evt := &opslevel.ServiceUpdateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
