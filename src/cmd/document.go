package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var createDocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Upload Swagger API documents via a file",
	Long: `Upload Swagger API documents via a file:

opslevel create document my-service -i xxxxx -f swagger.json

opslevel create document my-service -r services -t openapi -i xxxxx -f swagger.json
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceAlias := args[0]
		flags := cmd.Flags()
		resourceType := getResourceType(flags)
		documentType := getDocumentType(flags)
		integrationID, err := flags.GetString("integration-id")
		cobra.CheckErr(err)
		integrationURL := fmt.Sprintf("integrations/document/%s/%s/%s/%s", integrationID, resourceType, serviceAlias, documentType)
		fileContents, err := os.ReadFile(dataFile)
		cobra.CheckErr(err)
		var result opslevel.RestResponse
		response, err := getClientRest().R().
			SetResult(&result).
			SetBody(fileContents).
			SetHeader("Content-Type", "application/octet-stream").
			Post(integrationURL)
		cobra.CheckErr(err)
		if result.Result == "invalid_format" {
			log.Warn().Msgf("%s", result.Message)
		} else if response.IsSuccess() {
			log.Info().Msgf("Successfully registered api-doc for '%s'", serviceAlias)
		} else {
			log.Error().Msgf("%s", response)
		}
	},
}

func init() {
	createCmd.AddCommand(createDocumentCmd)

	createDocumentCmd.Flags().StringP("integration-id", "i", "", "OpsLevel integration ID")
	createDocumentCmd.Flags().StringP("resource-type", "r", "services", "OpsLevel Resource Type (options [\\\"services\\\"])")
	createDocumentCmd.Flags().StringP("document-type", "t", "openapi", "API Document Type (options [\\\"openapi\\\", \\\"swagger\\\"])")
}

func getResourceType(flags *pflag.FlagSet) string {
	resourceType, err := flags.GetString("resource-type")
	if err != nil {
		return "services"
	}
	switch strings.ToLower(resourceType) {
	case "services":
		return "services"
	case "service":
		return "services"
	default:
		return "services"
	}
}

func getDocumentType(flags *pflag.FlagSet) string {
	documentType, err := flags.GetString("document-type")
	if err != nil {
		return "openapi"
	}
	switch strings.ToLower(documentType) {
	case "openapi":
		return "openapi"
	case "swagger":
		return "openapi"
	default:
		return "openapi"
	}
}
