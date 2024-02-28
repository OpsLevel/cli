package cmd

import (
	"fmt"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var importFilepath string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to OpsLevel.",
	Long:  "Import data to OpsLevel.",
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.PersistentFlags().StringVarP(&importFilepath, "filepath", "f", "-", "File to read data from. Defaults to reading from stdin.")
}

func readImportFilepathAsCSV() (*common.CSVReader, error) {
	if importFilepath == "" {
		return nil, fmt.Errorf("empty filepath specified")
	}
	if importFilepath == "-" {
		return common.ReadCSVFile("/dev/stdin")
	}
	return common.ReadCSVFile(importFilepath)
}
