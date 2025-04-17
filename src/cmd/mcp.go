package cmd

import (
	"context"
	"encoding/json"
	
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
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

		s := server.NewMCPServer(
			"OpsLevel",
			"1.0.0",
		)

		// Register Teams
		s.AddTool(
			mcp.NewTool("teams",
				mcp.WithDescription("Get all the team names, identifiers and metadata for the opslevel account.  Teams are owners of other objects in opslevel. Only use this if you need to search all teams.")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client := getClientGQL()
				resp, err := client.ListTeams(nil)
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp.Nodes)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(data)), nil
			})

		// Register Users
		s.AddTool(
			mcp.NewTool("users", mcp.WithDescription("Get all the user names, e-mail addresses and metadata for the opslevel account.  Users are the people in opslevel. Only use this if you need to search all users.")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client := getClientGQL()
				resp, err := client.ListUsers(nil)
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp.Nodes)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(data)), nil
			})

		// Register Actions
		s.AddTool(
			mcp.NewTool("actions", mcp.WithDescription("Get all the information about actions the user can run in the opslevel account")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client := getClientGQL()
				resp, err := client.ListTriggerDefinitions(nil)
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp.Nodes)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(data)), nil
			})

		// Register Filters
		s.AddTool(
			mcp.NewTool("filters", mcp.WithDescription("Get all the rubric filter names and which predicates they have for the opslevel account")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client := getClientGQL()
				resp, err := client.ListFilters(nil)
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp.Nodes)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(data)), nil
			})

		// Register Components
		s.AddTool(
			mcp.NewTool("components", mcp.WithDescription("Get all the components in the opslevel account.  Components are objects in opslevel that represent things like apis, libraries, services, frontends, backends, etc.")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
				return mcp.NewToolResultText(string(data)), nil
			})

		// Register Infra
		s.AddTool(
			mcp.NewTool("infrastructure", mcp.WithDescription("Get all the infrastructure in the opslevel account.  Infrastructure are objects in opslevel that represent cloud provider resources like vpc, databases, caches, networks, vms, etc.")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client := getClientGQL()
				resp, err := client.ListInfrastructure(nil)
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp.Nodes)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(data)), nil
			})

		// Register Domains
		s.AddTool(
			mcp.NewTool("domains", mcp.WithDescription("Get all the domains in the opslevel account.  Domains are objects in opslevel that represent a top-level abstraction used to organize and categorize software systems.")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client := getClientGQL()
				resp, err := client.ListDomains(nil)
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp.Nodes)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(data)), nil
			})

		// Register Systems
		s.AddTool(
			mcp.NewTool("systems", mcp.WithDescription("Get all the systems in the opslevel account.  Systems are objects in opslevel that represent a grouping of services or components that act together to serve a business function or process.")),
			func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				client := getClientGQL()
				resp, err := client.ListSystems(nil)
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp.Nodes)
				if err != nil {
					return nil, err
				}
				return mcp.NewToolResultText(string(data)), nil
			})

		log.Info().Msg("Starting MCP server...")
		if err := server.ServeStdio(s); err != nil {
			panic(err)
		}
		<-done

		return nil
	},
}

func init() {
	betaCmd.AddCommand(mcpCmd)
}
