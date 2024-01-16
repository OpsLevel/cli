package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	keyValueExp    = regexp.MustCompile(`([\w-]+)=(.*)`)
	hasNextPageExp = regexp.MustCompile(`"hasNextPage":([\w]+)`)
	endCursorExp   = regexp.MustCompile(`"endCursor":\"([\w]+)\"`)
)

// graphqlCmd represents the graphql command
var graphqlCmd = &cobra.Command{
	Use:   "graphql",
	Short: "Make authenticated raw GraphQL requests",
	Long: `Make authenticated raw GraphQL requests.

Pass one or more '-f/--field' values in "key=value" format to add graphql variables
to the request payload.

In '--paginate' mode, all pages of results will sequentially be requested until
there are no more pages of results. This requires that the
original query accepts an '$endCursor: String' variable and that it fetches the
'pageInfo{ hasNextPage, endCursor }' set of fields from a collection.

Note that only the first 'endCursor' value found in the response body will be used
so ensure you are only paginating on 1 resource.  Nested resources pagination will
not work and will cause odd results or errors.

Use '--aggregate'' to specify a JQ expression to use as an function
for the results.  In '--paginate' mode it will use this expression to aggregate
the results into a JSON list.

Use '-q' to specify the graphql request body.
Pass "-" to read from standard input.
If the value starts with "@" it is interpreted as a filename to read from.
`,
	Example: `opslevel graphql --paginate -a=".account.services.nodes[]" -q='
query ($endCursor: String) {
  account {
    services(first: 5, after: $endCursor) {
      nodes {
        name
        aliases
        owner {
          name
        }
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
}'

opslevel graphql -f "owner=platform" -f "tier=tier_1" --paginate -a=".account.services.nodes[]" -q='
query ($endCursor: String, $owner: String!, $tier: String!) {
  account {
    services(first: 1, after: $endCursor, ownerAlias: $owner, tierAlias: $tier) {
      nodes {
        name
        aliases
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
}'

opslevel graphql -f "id=XXXXXX" -H "GraphQL-Visibility=internal" -a=".account.configFile.yaml" -q='
query ($id: ID!){
  account {
    configFile(id: $id) {
      yaml
    }
  }
}' | jq -r '.[0]' > opslevel.yml
`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		headersValue, err := flags.GetStringArray("header")
		headers := map[string]string{}
		for _, value := range headersValue {
			matches := keyValueExp.FindStringSubmatch(value)
			headers[matches[1]] = matches[2]
		}

		cobra.CheckErr(err)
		paginate, err := flags.GetBool("paginate")
		handleErr("error getting paginate flag", err)

		aggregate, err := flags.GetString("aggregate")
		jq, err := gojq.Parse(aggregate)
		handleErr("error parsing pagination flag value", err)
		aggregation, err := gojq.Compile(jq)

		handleErr("error compiling pagination flag value", err)
		queryValue, err := flags.GetString("query")
		handleErr("error getting query flag value", err)
		queryParsed, err := convert(queryValue)
		query, ok := queryParsed.(string)
		if !ok {
			handleErr("error parsing query flag value", fmt.Errorf("'%#v' is not a string", queryParsed))
		}
		operationName, err := flags.GetString("operationName")
		cobra.CheckErr(err)
		fields, err := flags.GetStringArray("field")
		handleErr("error getting field flag value", err)

		variables := map[string]interface{}{}
		for _, field := range fields {
			matches := keyValueExp.FindStringSubmatch(field)
			value, err := convert(matches[2])
			handleErr(fmt.Sprintf("error parsing variable '%s'", field), err)
			variables[matches[1]] = value
		}

		client := getClientGQL(opslevel.SetHeaders(headers))
		var output []interface{}

		hasNextPage := true
		for hasNextPage {
			data, err := client.ExecRaw(query, variables, opslevel.WithName(operationName))
			handleErr("error making graphql api call", err)
			output = append(output, handleAggregate(data, aggregation)...)

			if paginate {
				hasNextPage, err = strconv.ParseBool(string(hasNextPageExp.FindSubmatch(data)[1]))
				handleErr("error parsing bool for has next page", err)
				// don't try to parse endCursor unless we know there's another page
				if hasNextPage {
					variables["endCursor"] = string(endCursorExp.FindSubmatch(data)[1])
				}
			} else {
				hasNextPage = false
			}
		}

		json, err := json.Marshal(output)
		handleErr("error marshaling output to json", err)

		fmt.Println(string(json))
	},
}

func init() {
	rootCmd.AddCommand(graphqlCmd)

	graphqlCmd.Flags().StringArrayP("header", "H", nil, "Add a HTTP request header in `key=value` format")
	graphqlCmd.Flags().BoolP("paginate", "p", false, "Automatically make additional requests to fetch all pages of results")
	graphqlCmd.Flags().StringP("aggregate", "a", ".", "JQ expression to use to aggregate results")
	graphqlCmd.Flags().StringP("query", "q", "", "The query or mutation body to use")
	graphqlCmd.Flags().StringP("operationName", "o", "", "The query or mutation 'operation name' to use")
	graphqlCmd.Flags().StringArrayP("field", "f", nil, "Add a variable in `key=value` format")
}

func handleErr(msg string, err error) {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		os.Exit(1)
	}
}

func convert(v string) (interface{}, error) {
	if v == "-" {
		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		return convert(data)
	}

	if strings.HasPrefix(v, "@") {
		b, err := os.ReadFile(v[1:])
		if err != nil {
			return "", err
		}
		return convert(string(b))
	}

	if n, err := strconv.Atoi(v); err == nil {
		return n, nil
	}

	switch v {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	}
	return v, nil
}

func handleAggregate(data []byte, aggregation *gojq.Code) []interface{} {
	var parsed map[string]interface{}
	err := json.Unmarshal(data, &parsed)
	handleErr("error parsing graphql response to json", err)
	iter := aggregation.Run(parsed)
	var output []interface{}
	for {
		value, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := value.(error); ok {
			handleErr("error running aggregation function", err)
		}
		output = append(output, value)
	}
	return output
}
