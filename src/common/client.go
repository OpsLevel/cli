package common

import (
	"fmt"
	"github.com/opslevel/opslevel-go/v2022"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func NewGraphClient(version string) *opslevel.Client {
	client := opslevel.NewGQLClient(
		opslevel.SetAPIToken(viper.GetString("api-token")),
		opslevel.SetURL(viper.GetString("api-url")),
		opslevel.SetUserAgentExtra(fmt.Sprintf("cli-%s", version)))
	cobra.CheckErr(client.Validate())
	timeout := time.Second * time.Duration(viper.GetInt("api-timeout"))
	client := opslevel.NewGQLClient(opslevel.SetAPIToken(viper.GetString("api-token")),
		opslevel.SetURL(viper.GetString("api-url")),
		opslevel.SetTimeout(timeout),
		opslevel.SetUserAgentExtra(fmt.Sprintf("cli-%s", version)))

	clientErr := client.Validate()
	if clientErr != nil {
		if strings.Contains(clientErr.Error(), "Please provide a valid OpsLevel API token") {
			cobra.CheckErr(fmt.Errorf("%s via 'export OPSLEVEL_API_TOKEN=XXX' or '--api-token=XXX'", clientErr.Error()))
		} else {
			cobra.CheckErr(clientErr)
		}
	}
	cobra.CheckErr(clientErr)

	return client
}
