package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/opslevel/opslevel-go/v2025"

	"github.com/rs/zerolog/log"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var exampleServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Example service",
	Long:    `Example service`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample(opslevel.ServiceCreateInput{
			Name:           "example_name",
			Description:    opslevel.RefOf("example_description"),
			Framework:      opslevel.RefOf("example_framework"),
			Language:       opslevel.RefOf("example_language"),
			LifecycleAlias: opslevel.RefOf("example_lifecycle"),
			OwnerAlias:     opslevel.RefOf("example_owner"),
			Parent:         opslevel.NewIdentifier("example_parent"),
			Product:        opslevel.RefOf("example_product"),
			TierAlias:      opslevel.RefOf("example_tier"),
			Type:           opslevel.NewIdentifier("example_type"),
		}))
	},
}

var createServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Create a service",
	Long: `Create a service

cat << EOF | opslevel create service -f -
name: "hello world"
description: "Hello World Service"
framework: "fasthttp"
language: "Go"
lifecycle: beta
owner:
  alias: "platform"
parent:
  alias: "my_system"
product: "OSS"
tier: "tier_4"
type:
  alias: "mobile-app"
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
	Aliases:    []string{"svc"},
	Short:      "Get details about a service",
	Long:       `Get details about a service`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		key := args[0]
		client := getClientGQL()
		service, err := getService(key)
		cobra.CheckErr(err)

		// Extra fields only displayed in JSON format
		if isJsonOutput() {
			_, err = service.GetDependents(client, nil)
			cobra.CheckErr(err)
			_, err = service.GetDependencies(client, nil)
			cobra.CheckErr(err)
			_, err = service.GetProperties(client, nil)
			cobra.CheckErr(err)
		}

		common.WasFound(service.Id == "", key)
		common.PrettyPrint(service)
	},
}

var listServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"services", "svc", "svcs"},
	Short:   "Lists services",
	Long:    `Lists services`,
	Run: func(cmd *cobra.Command, args []string) {
		list := []opslevel.Service{}
		client := getClientGQL()
		resp, err := client.ListServices(nil)
		cobra.CheckErr(err)
		for _, service := range resp.Nodes {
			if !isJsonOutput() {
				list = append(list, service)
				continue
			}

			// Extra fields only displayed in JSON format
			if ok, _ := cmd.Flags().GetBool("dependencies"); ok {
				_, err = service.GetDependencies(client, nil)
				cobra.CheckErr(err)
			}
			if ok, _ := cmd.Flags().GetBool("dependents"); ok {
				_, err = service.GetDependents(client, nil)
				cobra.CheckErr(err)
			}
			if ok, _ := cmd.Flags().GetBool("properties"); ok {
				_, err = service.GetProperties(client, nil)
				cobra.CheckErr(err)
			}
			list = append(list, service)
		}
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"NAME", "ID", "TYPE", "ALIASES"})
			for _, item := range list {
				w.Write([]string{item.Name, string(item.Id), item.Type.Aliases[0], strings.Join(item.Aliases, "/")})
			}
			w.Flush()
		} else {
			w := common.NewTabWriter("NAME", "ID", "TYPE", "ALIASES")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", item.Name, item.Id, item.Type.Aliases[0], strings.Join(item.Aliases, ","))
			}
			w.Flush()
		}
	},
}

var updateServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc"},
	Short:   "Update a service",
	Long: `Update a service

cat << EOF | opslevel update service -f -
name: "hello world"
alias: "hello_world"
description: "Hello World Service Updated"
framework: "fasthttp"
language: "Go"
lifecycle: beta
owner:
  alias: "platform"
parent:
  alias: "my_system"
product: "OSS"
tier: "tier_3"
type:
  alias: "mobile-app"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readResourceInput[opslevel.ServiceUpdateInput]()
		if len(args) == 1 {
			key := args[0]
			if opslevel.IsID(key) {
				input.Id = opslevel.RefOf(opslevel.ID(key))
			} else {
				input.Alias = opslevel.RefOf(key)
			}
		}
		cobra.CheckErr(err)
		service, err := getClientGQL().UpdateService(*input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(service, "", "    "))
	},
}

var deleteServiceCmd = &cobra.Command{
	Use:        "service ID|ALIAS",
	Aliases:    []string{"svc"},
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
	Aliases: []string{"services", "svc", "svcs"},
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
					input.TierAlias = opslevel.RefOf(item.Alias)
				}
			}
			lifecycle := reader.Text("Lifecycle")
			if lifecycle != "" {
				if item, ok := opslevel.Cache.Lifecycles[lifecycle]; ok {
					input.LifecycleAlias = opslevel.RefOf(item.Alias)
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

	listServiceCmd.PersistentFlags().Bool("dependencies", false, "Include dependencies of each service")
	listServiceCmd.PersistentFlags().Bool("dependents", false, "Include dependents of each service")
	listServiceCmd.PersistentFlags().Bool("properties", false, "Include properties of each service")

	importCmd.AddCommand(importServicesCmd)
}

func NullableString(value *string) *opslevel.Nullable[string] {
	if value == nil {
		return nil
	}
	if *value == "" {
		return opslevel.NewNull()
	}
	return opslevel.NewNullableFrom(*value)
}

func getService(identifier string) (*opslevel.Service, error) {
	service, err := getClientGQL().GetService(identifier)
	if service == nil || !opslevel.IsID(string(service.Id)) {
		err = fmt.Errorf("service with identifier '%s' not found", identifier)
	}

	return service, err
}
