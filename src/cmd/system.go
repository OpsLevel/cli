package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var createSystemCmd = &cobra.Command{
	Use:   "system",
	Short: "Create a system",
	Long:  `Create a system`,
	Example: `
		cat << EOF | opslevel create system -f -
		name: "My System"
		description: "Hello World System"
		owner: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NjY"
		parent:
			alias: "alias of domain"
		note: "Additional system details"
		EOF
		`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readSystemInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateSystem(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

var getSystemCmd = &cobra.Command{
	Use:        "system ID|ALIAS",
	Short:      "Get details about a system",
	Long:       `Get details about a system`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Example: `
		opslevel get system my-system-alias-or-id
		`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		result, err := getClientGQL().GetSystem(key)
		cobra.CheckErr(err)
		common.WasFound(result.Id == "", key)
		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
	},
}

var listSystemCmd = &cobra.Command{
	Use:     "system",
	Aliases: []string{"systems"},
	Short:   "Lists the systems",
	Long:    `Lists the systems`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListSystems(nil)
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
			w := common.NewTabWriter("NAME", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Id, strings.Join(item.Aliases, ","))
			}
			w.Flush()
		}
	},
}

// TODO: bug in API prevents use of alias in this function.  Adding full functionality for now.
var updateSystemCmd = &cobra.Command{
	Use:   "system ID|ALIAS",
	Short: "Update an OpsLevel system",
	Long:  "Update an OpsLevel system",
	Example: `
		cat << EOF | opslevel update system my-system-alias-or-id -f -
		name: "My Updated System"
		description: "Hello Updated System"
		owner: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NjY"
		parent:
			alias: "my_domain"
		note: "Additional system details for my updated system"
		EOF
		`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readSystemInput()
		cobra.CheckErr(err)
		system, err := getClientGQL().UpdateSystem(key, *input)
		cobra.CheckErr(err)
		fmt.Println(system.Id)
	},
}

var deleteSystemCmd = &cobra.Command{
	Use:   "system ID|ALIAS",
	Short: "Delete a system",
	Long:  "Delete a system from OpsLevel",
	Example: `
		opslevel delete system my-system-alias-or-id
		`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteSystem(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' system\n", key)
	},
}

func init() {
	createCmd.AddCommand(createSystemCmd)
	getCmd.AddCommand(getSystemCmd)
	listCmd.AddCommand(listSystemCmd)
	updateCmd.AddCommand(updateSystemCmd)
	deleteCmd.AddCommand(deleteSystemCmd)
}

func readSystemInput() (*opslevel.SystemInput, error) {
	readCreateConfigFile()
	evt := &opslevel.SystemInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
