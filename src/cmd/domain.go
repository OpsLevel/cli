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

var createDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Create a domain",
	Long: `Create a domain

cat << EOF | opslevel create domain -f -
name: "My Domain"
description: "Hello World Domain"
ownerId: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NjY"
note: "Additional details"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readDomainCreateInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateDomain(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

var getDomainCmd = &cobra.Command{
	Use:        "domain ID|ALIAS",
	Short:      "Get details about a domain",
	Long:       `Get details about a domain`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		result, err := getClientGQL().GetDomain(key)
		cobra.CheckErr(err)
		common.WasFound(result.Id == "", key)
		common.PrettyPrint(result)
	},
}

var getDomainSystemCmd = &cobra.Command{
	Use:        "system ID|ALIAS",
	Aliases:    []string{"systems"},
	Short:      "Get systems for a domain",
	Long:       `The systems that belong to a domain.`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		domain, err := getClientGQL().GetDomain(key)
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

var getDomainTagCmd = &cobra.Command{
	Use:     "tag ID|ALIAS TAG_KEY",
	Aliases: []string{"tags"},
	Short:   "Get a domain's tag",
	Long: `Get a domain's' tag

opslevel get domain tag my_domain | jq 'from_entries'
opslevel get domain tag my_domain my-tag
`,
	Args:       cobra.MinimumNArgs(1),
	ArgAliases: []string{"ID", "ALIAS", "TAG_KEY"},
	Run: func(cmd *cobra.Command, args []string) {
		domainKey := args[0]
		singleTag := len(args) == 2
		var tagKey string
		if singleTag {
			tagKey = args[1]
		}

		domain, err := getClientGQL().GetDomain(domainKey)
		cobra.CheckErr(err)
		if domain.Id == "" {
			cobra.CheckErr(fmt.Errorf("domain '%s' not found", domainKey))
		}
		var output []opslevel.Tag
		tags, err := domain.Tags(getClientGQL(), nil)
		cobra.CheckErr(err)
		for _, tag := range tags.Nodes {
			if singleTag == false || tagKey == tag.Key {
				output = append(output, tag)
			}
		}
		if len(output) == 0 {
			cobra.CheckErr(fmt.Errorf("tag with key '%s' not found on domain '%s'", tagKey, domainKey))
		}
		common.PrettyPrint(output)
	},
}

var listDomainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"domains"},
	Short:   "Lists the domains",
	Long:    `Lists the domains`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListDomains(nil)
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
	createCmd.AddCommand(createDomainCmd)
	getCmd.AddCommand(getDomainCmd)
	getDomainCmd.AddCommand(getDomainSystemCmd)
	getDomainCmd.AddCommand(getDomainTagCmd)
	listCmd.AddCommand(listDomainCmd)
}

func readDomainCreateInput() (*opslevel.DomainCreateInput, error) {
	readCreateConfigFile()
	evt := &opslevel.DomainCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
