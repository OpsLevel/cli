package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
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
	Use:   "tag ID|ALIAS",
	Short: "Create a service tag",
	Long: `Create a service tag
	
cat << EOF | opslevel create service tag my-service
key: "foo"
value: "bar"
EOF
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readTagCreateInput()
		if common.IsID(key) {
			input.Id = key
		} else {
			input.Alias = key
		}
		input.Type = opslevel.TaggableResourceService
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateTag(*input)
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
		common.PrettyPrint(result)
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
				w.Write([]string{item.Name, item.Id.(string), strings.Join(item.Aliases, ",")})
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

func init() {
	createCmd.AddCommand(createServiceCmd)
	getCmd.AddCommand(getServiceCmd)
	listCmd.AddCommand(listServiceCmd)
	updateCmd.AddCommand(updateServiceCmd)
	deleteCmd.AddCommand(deleteServiceCmd)

	createServiceCmd.AddCommand(createServiceTagCmd)
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

func readTagCreateInput() (*opslevel.TagCreateInput, error) {
	readCreateConfigFile()
	evt := &opslevel.TagCreateInput{}
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
