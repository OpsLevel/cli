package cmd

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var createDocumentCmd = &cobra.Command{
	Use:   "document",
	Short: "Upload Swagger API documents via a file",
	Long: `Upload Swagger API documents via a file:

opslevel create document my-service -i xxxxx -f swagger.json
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serviceAlias := args[0]
		integrationID, err := cmd.Flags().GetString("integration-id")
		cobra.CheckErr(err)
		integrationURL := fmt.Sprintf("https://app.opslevel.com/integrations/api_docs/%s/%s", integrationID, serviceAlias)
		fileContents, err := ioutil.ReadFile(createDataFile)
		cobra.CheckErr(err)
		var resp struct {
			Result string `json:"result"`
		}
		err = getClientRest().Do("POST", "application/octet-stream", integrationURL, fileContents, &resp)
		cobra.CheckErr(err)
		log.Info().Msgf("Successfully registered api-doc for '%s'", serviceAlias)
		log.Info().Msgf("%v", resp)

	},
}

func init() {
	createCmd.AddCommand(createDocumentCmd)

	createDocumentCmd.Flags().StringP("integration-id", "i", "", "OpsLevel integration ID")

	//	createCmd.PersistentFlags().StringVarP(&createDataFile, "file", "f", "-", "File to read data from. If '.' then reads from './data.yaml'. Defaults to reading from stdin.")
	//	viper.BindPFlags(createCmd.Flags())
}
