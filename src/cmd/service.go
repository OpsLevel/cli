package cmd

import (
	"encoding/json"
	"fmt"
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
		service, err := getClientGQL().CreateService(*input)
		cobra.CheckErr(err)
		fmt.Println(service.Id)
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
		var service *opslevel.Service
		var err error
		if common.IsID(key) {
			service, err = getClientGQL().GetService(key)
			cobra.CheckErr(err)
		} else {
			service, err = getClientGQL().GetServiceWithAlias(key)
			cobra.CheckErr(err)
		}
		cobra.CheckErr(err)
		common.PrettyPrint(service)
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
