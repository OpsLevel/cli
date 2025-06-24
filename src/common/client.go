package common

import (
	"fmt"
	"time"

	"github.com/opslevel/opslevel-go/v2025"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGraphClient(version string, options ...opslevel.Option) *opslevel.Client {
	timeout := time.Second * time.Duration(viper.GetInt("api-timeout"))
	options = append(
		options,
		opslevel.SetAPIToken(viper.GetString("api-token")),
		opslevel.SetURL(viper.GetString("api-url")),
		opslevel.SetTimeout(timeout),
		opslevel.SetUserAgentExtra(fmt.Sprintf("cli-%s", version)),
	)
	client := opslevel.NewGQLClient(options...)

	clientErr := client.Validate()
	cobra.CheckErr(clientErr)

	return client
}
