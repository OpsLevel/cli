package cmd

import (
	"encoding/json"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/spf13/cobra"
)

type NullArguments struct{}

type LightweightComponent struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
	URL   string `json:"url"`
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP Server",
	Long:  "MCP Server",

	RunE: func(cmd *cobra.Command, args []string) error {
		done := make(chan struct{})

		// transport := http.NewHTTPTransport("/mcp")
		// transport.WithAddr(":8080")
		// server := mcp_golang.NewServer(transport)
		server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

		// Register Teams
		if err := server.RegisterTool("teams", "Get all the team names, identifiers and metadata for the opslevel account.  Teams are owners of other objects in opslevel. Only use this if you need to search all teams.", func(args NullArguments) (*mcp_golang.ToolResponse, error) {
			client := getClientGQL()
			resp, err := client.ListTeams(nil)
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(resp.Nodes)
			if err != nil {
				return nil, err
			}
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(data))), nil
		}); err != nil {
			panic(err)
		}

		// Register Users
		if err := server.RegisterTool("users", "Get all the user names, e-mail addresses and metadata for the opslevel account.  Users are the people in opslevel. Only use this if you need to search all users.", func(args NullArguments) (*mcp_golang.ToolResponse, error) {
			client := getClientGQL()
			resp, err := client.ListUsers(nil)
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(resp.Nodes)
			if err != nil {
				return nil, err
			}
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(data))), nil
		}); err != nil {
			panic(err)
		}

		// Register Actions
		if err := server.RegisterTool("actions", "Get all the information about actions the user can run in the opslevel account", func(args NullArguments) (*mcp_golang.ToolResponse, error) {
			client := getClientGQL()
			resp, err := client.ListTriggerDefinitions(nil)
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(resp.Nodes)
			if err != nil {
				return nil, err
			}
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(data))), nil
		}); err != nil {
			panic(err)
		}

		// Register Filters
		if err := server.RegisterTool("filters", "Get all the rubric filter names and which predicates they have for the opslevel account", func(args NullArguments) (*mcp_golang.ToolResponse, error) {
			client := getClientGQL()
			resp, err := client.ListFilters(nil)
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(resp.Nodes)
			if err != nil {
				return nil, err
			}
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(data))), nil
		}); err != nil {
			panic(err)
		}

		// Register Components
		if err := server.RegisterTool("components", "Get all the components in the opslevel account.  Components are objects in opslevel that represent things like apis, libraries, services, frontends, backends, etc.", func(args NullArguments) (*mcp_golang.ToolResponse, error) {
			client := getClientGQL()
			resp, err := client.ListServices(nil)
			if err != nil {
				return nil, err
			}
			var components []LightweightComponent
			for _, node := range resp.Nodes {
				components = append(components, LightweightComponent{
					Id:    string(node.Id),
					Name:  node.Name,
					Owner: node.Owner.Alias,
					URL:   node.HtmlURL,
				})
			}
			data, err := json.Marshal(components)
			if err != nil {
				return nil, err
			}
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(data))), nil
		}); err != nil {
			panic(err)
		}

		// Register Infra
		if err := server.RegisterTool("infrastructure", "Get all the infrastructure in the opslevel account.  Infrastructure are objects in opslevel that represent cloud provider resources like vpc, databases, caches, networks, vms, etc.", func(args NullArguments) (*mcp_golang.ToolResponse, error) {
			client := getClientGQL()
			resp, err := client.ListInfrastructure(nil)
			if err != nil {
				return nil, err
			}
			data, err := json.Marshal(resp.Nodes)
			if err != nil {
				return nil, err
			}
			return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(data))), nil
		}); err != nil {
			panic(err)
		}

		if err := server.Serve(); err != nil {
			panic(err)
		}
		<-done
		return nil
	},
}

func init() {
	betaCmd.AddCommand(mcpCmd)
}
