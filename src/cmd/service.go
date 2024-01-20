package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/rs/zerolog/log"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var exampleServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: common.GetAliases("Service"),
	Short:   "Example service",
	Long:    `Example service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.ServiceCreateInput]())
	},
}

var createServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: common.GetAliases("Service"),
	Short:   "Create a service",
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
		input, err := readResourceInput[opslevel.ServiceCreateInput]()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateService(*input)
		cobra.CheckErr(err)
		common.PrettyPrint(result.Id)
	},
}

var getServiceCmd = &cobra.Command{
	Use:        "service ID|ALIAS",
	Aliases:    common.GetAliases("Service"),
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

var listServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: common.GetAliases("Service"),
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
	Use:     "service",
	Aliases: common.GetAliases("Service"),
	Short:   "Update a service",
	Long: `Update a service

cat << EOF | opslevel update service -f -
alias: "hello_world"
description: "Hello World Service Updated"
tier: "tier_3"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readResourceInput[opslevel.ServiceUpdateInput]()
		cobra.CheckErr(err)
		service, err := getClientGQL().UpdateService(*input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(service, "", "    "))
	},
}

var deleteServiceCmd = &cobra.Command{
	Use:        "service ID|ALIAS",
	Aliases:    common.GetAliases("Service"),
	Short:      "Delete a service",
	Long:       `Delete a service`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteService(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' service\n", key)
	},
}

var importServicesCmd = &cobra.Command{
	Use:     "service",
	Aliases: common.GetAliases("Service"),
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
				Description: opslevel.RefOf(reader.Text("Description")),
				Product:     opslevel.RefOf(reader.Text("Product")),
				Language:    opslevel.RefOf(reader.Text("Language")),
				Framework:   opslevel.RefOf(reader.Text("Framework")),
			}
			tier := reader.Text("Tier")
			if tier != "" {
				if item, ok := opslevel.Cache.Tiers[tier]; ok {
					input.TierAlias = &item.Alias
				}
			}
			lifecycle := reader.Text("Lifecycle")
			if lifecycle != "" {
				if item, ok := opslevel.Cache.Lifecycles[lifecycle]; ok {
					input.LifecycleAlias = &item.Alias
				}
			}
			owner := reader.Text("Owner")
			if owner != "" {
				if item, ok := opslevel.Cache.Teams[owner]; ok {
					input.OwnerInput = opslevel.NewIdentifier(item.Alias)
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
	exampleCmd.AddCommand(exampleServiceCmd)
	createCmd.AddCommand(createServiceCmd)
	getCmd.AddCommand(getServiceCmd)
	listCmd.AddCommand(listServiceCmd)
	updateCmd.AddCommand(updateServiceCmd)
	deleteCmd.AddCommand(deleteServiceCmd)

	importCmd.AddCommand(importServicesCmd)
}
