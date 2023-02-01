/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

Use '-q' to specify the graphql request body.
Pass "-" to read from standard input.
If the value starts with "@" it is interpreted as a filename to read from.
`,
	Example: `opslevel graphql --paginate -q='
    query($endCursor: String) {
      account {
        services(first: 100, after: $endCursor) {
          nodes {
            name
            aliases
            owner { name }
          }
          pageInfo {
            hasNextPage
            endCursor
          }
        }
      }
    }
  '

opslevel graphql -f owner='my-team' -f tier="tier_1" -q='
    query($endCursor: String, $owner: String!, $tier: String!) {
      account {
        services(first: 100, after: $endCursor, ownerAlias: $owner, tierAlias: $tier) {
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
    }
  '
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("graphql called")
	},
}

func init() {
	rootCmd.AddCommand(graphqlCmd)

	graphqlCmd.Flags().StringArrayP("header", "H", nil, "Add a HTTP request header in `key:value` format")
	graphqlCmd.Flags().BoolP("paginate", "p", false, "Automatically make additional requests to fetch all pages of results")
	graphqlCmd.Flags().StringP("query", "q", "", "The query or mutation body to use")
	graphqlCmd.Flags().StringP("operationName", "o", "", "The query or mutation 'operation name' to use")
	graphqlCmd.Flags().StringArrayP("field", "f", nil, "Add a variable in `key=value` format")

}
