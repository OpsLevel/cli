package cmd

import (
	"encoding/json"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/spf13/cobra"
)

type NullArguments struct{}

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

		listCommands := []struct {
			name        string
			description string
			call       	func() (any, error)
		}{
			{
				"teams",
				"List all teams",
				func() (any, error) { return getClientGQL().ListTeams(nil) },
			},
		}

		for _, cmd := range listCommands {
			if err := server.RegisterTool(cmd.name, cmd.description, func(args NullArguments) (*mcp_golang.ToolResponse, error) {
				resp, err := cmd.call()
				if err != nil {
					return nil, err
				}
				data, err := json.Marshal(resp)
				if err != nil {
					return nil, err
				}
				return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(string(data))), nil
			}); err != nil {
				panic(err)
			}
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
