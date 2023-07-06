package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Create a domain",
	Long:  `Create a domain`,
	Example: `

cat << EOF | opslevel create domain -f -
name: "My Domain"
description: "Hello World Domain"
owner: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83NjY"
note: "Additional details"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readDomainInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateDomain(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

var deleteDomainCmd = &cobra.Command{
	Use:   "domain ID|ALIAS",
	Short: "Delete a domain",
	Long:  `Delete a domain`,
	Example: `
opslevel delete domain my_domain
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteDomain(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' domain\n", key)
	},
}

var getDomainCmd = &cobra.Command{
	Use:   "domain ID|ALIAS",
	Short: "Get details about a domain",
	Long:  `Get details about a domain`,
	Example: `
opslevel get domain my_domain
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		result, err := getClientGQL().GetDomain(key)
		cobra.CheckErr(err)
		common.WasFound(result.Id == "", key)
		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
	},
}

var listDomainCmd = &cobra.Command{
	Use:     "domain",
	Aliases: []string{"domains"},
	Short:   "Lists the domains",
	Long:    `Lists the domains`,
	Example: `
opslevel list domain
`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListDomains(nil)
		list := resp.Nodes
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			if err := w.Write([]string{"NAME", "ID", "ALIASES"}); err != nil {
				panic(err)
			}
			for _, item := range list {
				err := w.Write([]string{item.Name, string(item.Id), strings.Join(item.Aliases, "/")})
				if err != nil {
					panic(err)
				}
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

var updateDomainCmd = &cobra.Command{
	Use:   "domain ID|ALIAS",
	Short: "Update a domain",
	Long: `Update a domain

cat << EOF | opslevel update domain my_domain -f -
name: "My New Domain"
description: "Hello World New Domain"
owner: "Z2lkOi8vb3BzbGV2ZWwvVGVhbS83ODk"
note: "Additional details for my new domain"
EOF
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readDomainInput()
		cobra.CheckErr(err)
		domain, err := getClientGQL().UpdateDomain(key, *input)
		cobra.CheckErr(err)
		fmt.Println(domain.Id)
	},
}

func init() {
	createCmd.AddCommand(createDomainCmd)
	deleteCmd.AddCommand(deleteDomainCmd)
	getCmd.AddCommand(getDomainCmd)
	listCmd.AddCommand(listDomainCmd)
	updateCmd.AddCommand(updateDomainCmd)
}

func readDomainInput() (*opslevel.DomainInput, error) {
	readCreateConfigFile()
	evt := &opslevel.DomainInput{}
	if err := viper.Unmarshal(&evt); err != nil {
		return nil, err
	}
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
