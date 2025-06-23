package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2025"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var exampleInfraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Example infrastructure resource",
	Long:  `Example infrastructure resource`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample(opslevel.InfrastructureResourceInput{
			Schema: &opslevel.InfrastructureResourceSchemaInput{
				Type: "example_schema",
			},
			ProviderData: &opslevel.InfrastructureResourceProviderDataInput{
				AccountName:  "example_account",
				ExternalUrl:  opslevel.RefOf("example_external_url"),
				ProviderName: opslevel.RefOf("example_provider"),
			},
			ProviderResourceType: opslevel.RefOf("example_provider_resource_type"),
			OwnerId:              opslevel.RefOf(opslevel.ID("Z2lkOi8vc2VydmljZS8xMjM0NTY3ODk")),
		}))
	},
}

var createInfraCmd = &cobra.Command{
	Use:   "infra",
	Short: "Create an infrastructure resource",
	Long:  `Create an infrastructure resource`,
	Example: `
cat << EOF | opslevel create infra -f -
schema: "Database"
provider:
    account: "Dev - 123456789"
    name: "GCP"
    type: "BigQuery"
    url: "https://google.com"
data: 
    name: "my-big-query"
    endpoint: "https://google.com"
    engine: "BigQuery"
    replica: false
    storage_size: |-
        value: 1024,
        unit: "GB"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readInfraInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateInfrastructure(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

var getInfraSchemaCmd = &cobra.Command{
	Use:        "infra-schema ALIAS",
	Aliases:    []string{"infraschema"},
	Short:      "Get infrastructure schemas",
	Long:       `Get infrastructure schemas`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ALIAS"},
	Example: `
opslevel --log-level=ERROR get infra-schema Database > ~/.opslevel/schemas/database.json
opslevel --log-level=ERROR get infra-schema Network > ~/.opslevel/schemas/network.json
opslevel --log-level=ERROR get infra-schema Compute > ~/.opslevel/schemas/compute.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		client := getClientGQL()
		opslevel.Cache.CacheInfraSchemas(client)
		schema, found := opslevel.Cache.TryGetInfrastructureSchema(key)
		if !found {
			log.Error().Msgf("unable to find infrastructure schema '%s'", key)
			return
		}
		// TODO: this sucks we need to make the JSON type in opslevel-go more robust
		dto := map[string]any{}
		for k, v := range schema.Schema {
			dto[k] = v
		}
		if isYamlOutput() {
			common.YamlPrint(dto)
		} else {
			common.PrettyPrint(dto)
		}
	},
}

var getInfraCmd = &cobra.Command{
	Use:        "infra ID|ALIAS",
	Short:      "Get details about an infrastructure resource",
	Long:       `Get details about an infrastructure resource`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Example:    `opslevel get infra my-infra-alias-or-id`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		result, err := getClientGQL().GetInfrastructure(key)
		cobra.CheckErr(err)
		common.WasFound(result.Id == "", key)
		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
	},
}

var listInfraSchemasCmd = &cobra.Command{
	Use:     "infra-schema",
	Aliases: []string{"infra-schemas", "infraschema", "infraschemas"},
	Short:   "List infrastructure schemas",
	Long:    `List infrastructure schemas`,
	Example: `
mkdir -p ~/.opslevel/schemas/
for DATA in $(opslevel list infra-schemas -o json | jq -r '.[] | @base64');
do
	SCHEMA="$(echo ${DATA} | base64 --decode | jq -r '.schema')"
	TYPE="$(echo ${DATA} | base64 --decode | jq -r '.type')"
	FILENAME="$(echo ${TYPE} | awk '{print tolower($0)}' | sed 's/ /_/g')"
	echo "${SCHEMA}" > ~/.opslevel/schemas/${FILENAME}.json
	echo "[opslevel] wrote infra-schema '${FILENAME}'"
done`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListInfrastructureSchemas(nil)
		list := resp.Nodes
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"TYPE"})
			for _, item := range list {
				w.Write([]string{item.Type})
			}
			w.Flush()
		} else {
			w := common.NewTabWriter("TYPE")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t\n", item.Type)
			}
			w.Flush()
		}
	},
}

var listInfraCmd = &cobra.Command{
	Use:     "infra",
	Aliases: []string{"infras"},
	Short:   "List infrastructure resources",
	Long:    `List infrastructure resources`,
	Example: `
# list all my unique network CIDRs
opslevel list infra -o json | jq 'map(select(.type == "Network") | .data | fromjson | .ipv4_cidr) | unique'
# list all my networks with only the information i care about
opslevel list infra -o json | jq 'map(select(.type == "Network") | .data | fromjson | {name: .name, cidr: .ipv4_cidr, region: .zone})'
# list all my unique compute image ids
opslevel list infra -o json | jq 'map(select(.type == "Compute") | .data | fromjson | .image_id) | unique'
# list all my database to see if they are public
opslevel list infra -o json | jq 'map(select(.type == "Database") | .data | fromjson | {name: .name, public: .publicly_accessible})'
# list all my databases to see their storage size
opslevel list infra -o json | jq 'map(select(.type == "Database") | .data | fromjson | {name: .name, size: "\(.storage_size.value) \(.storage_size.unit)"})'  
`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListInfrastructure(nil)
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

var updateInfraCmd = &cobra.Command{
	Use:        "infra ID|ALIAS",
	Short:      "Update an infrastructure resource",
	Long:       `Update an infrastructure resource`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Example: `
cat << EOF | opslevel update infra my-database -f -
owner: "devs"
provider:
    account: "dev"
    url: "https://google.com"
data:
    name: "my-big-query"
    endpoint: "https://google.com"
    engine: "BigQuery"
    replica: false
    storage_size:
		unit: "GB"
		value: 100
	storage_iops:
		unit: "IOPS"
		value: 12000
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readInfraInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().UpdateInfrastructure(key, *input)
		cobra.CheckErr(err)
		fmt.Println(string(result.Id))
	},
}

var deleteInfraCmd = &cobra.Command{
	Use:        "infra ID|ALIAS",
	Short:      "Delete an infrastructure resource",
	Long:       `Delete an infrastructure resource`,
	Example:    `opslevel delete system my-system-alias-or-id`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteInfrastructure(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' infrastructure resource\n", key)
	},
}

func init() {
	exampleCmd.AddCommand(exampleInfraCmd)
	createCmd.AddCommand(createInfraCmd)
	getCmd.AddCommand(getInfraSchemaCmd)
	getCmd.AddCommand(getInfraCmd)
	listCmd.AddCommand(listInfraSchemasCmd)
	listCmd.AddCommand(listInfraCmd)
	updateCmd.AddCommand(updateInfraCmd)
	deleteCmd.AddCommand(deleteInfraCmd)
}

func readInfraInput() (*opslevel.InfraInput, error) {
	file, err := io.ReadAll(os.Stdin)
	cobra.CheckErr(err)
	evt := &opslevel.InfraInput{}
	cobra.CheckErr(yaml.Unmarshal(file, &evt))
	return evt, nil
}
