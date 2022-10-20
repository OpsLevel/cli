package common

import (
	"fmt"
	"github.com/opslevel/opslevel-go/v2022"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGraphClient(version string) *opslevel.Client {
	client := opslevel.NewGQLClient(
		opslevel.SetAPIToken(viper.GetString("api-token")),
		opslevel.SetURL(viper.GetString("api-url")),
		opslevel.SetUserAgentExtra(fmt.Sprintf("cli-%s", version)))
	cobra.CheckErr(client.Validate())
	return client
}
