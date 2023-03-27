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
	Long: `Create a system
cat << EOF | opslevel create system -f -
name: "My System"
description: "Hello World System"
ownerId: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NjY"
parent:
	alias: "Name of parent domain"
note: "Additional system details"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readSystemCreateInput()
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
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		result, err := getClientGQL().GetSystem(key)
		cobra.CheckErr(err)
		common.WasFound(result.Id == "", key)
		common.PrettyPrint(result)
	},
}

// The story for this seems to be the need to retrieve all the parent/attached domains for a given system
var getSystemDomainCmd = &cobra.Command{
	Use:        "domain ID|ALIAS",
	Aliases:    []string{"systems"},
	Short:      "Get domains attached to a given system",
	Long:       `Get domains attached to a given system.`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		system, err := getClientGQL().GetSystem(key) //it gets wonky here
		cobra.CheckErr(err)
		common.WasFound(domain.Id == "", key)
		resp, err := domain.ChildSystems(getClientGQL(), nil)
		systems := resp.Nodes
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.PrettyPrint(systems)
		} else {
			w := common.NewTabWriter("Name", "ID")
			for _, item := range systems {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Name, item.Id)
			}
			w.Flush()
		}
	},
}

var getSystemTagCmd = &cobra.Command{
	Use:     "tag ID|ALIAS TAG_KEY",
	Aliases: []string{"tags"},
	Short:   "Get a system's tags",
	Long: `Get a system's' tags
opslevel get system tag my_system | jq 'from_entries'
opslevel get system tag my_system my-tag
`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"ID", "ALIAS", "TAG_KEY"},
	Run: func(cmd *cobra.Command, args []string) {
		systemKey := args[0]
		singleTag := len(args) == 2
		var tagKey string
		if singleTag {
			tagKey = args[1]
		}

		system, err := getClientGQL().GetSystem(systemKey)
		cobra.CheckErr(err)
		if system.Id == "" {
			cobra.CheckErr(fmt.Errorf("system '%s' not found", systemKey))
		}
		var output []opslevel.Tag
		tags, err := system.Tags(getClientGQL(), nil)
		cobra.CheckErr(err)
		for _, tag := range tags.Nodes {
			if singleTag == false || tagKey == tag.Key {
				output = append(output, tag)
			}
		}
		if len(output) == 0 {
			cobra.CheckErr(fmt.Errorf("tag with key '%s' not found on system '%s'", tagKey, systemKey))
		}
		common.PrettyPrint(output)
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

func init() {
	createCmd.AddCommand(createSystemCmd)
	getCmd.AddCommand(getSystemCmd)
	getSystemCmd.AddCommand(getSystemDomainCmd)
	getSystemCmd.AddCommand(getSystemTagCmd)
	listCmd.AddCommand(listSystemCmd)
}

func readSystemCreateInput() (*opslevel.SystemCreateInput, error) {
	readCreateConfigFile()
	evt := &opslevel.SystemCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
