package common

import (
	"fmt"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func NewGraphClient(version string, options ...opslevel.Option) *opslevel.Client {
	timeout := time.Second * time.Duration(viper.GetInt("api-timeout"))
	options = append(options, opslevel.SetAPIToken(viper.GetString("api-token")))
	options = append(options, opslevel.SetURL(viper.GetString("api-url")))
	options = append(options, opslevel.SetTimeout(timeout))
	options = append(options, opslevel.SetUserAgentExtra(fmt.Sprintf("cli-%s", version)))
	options = append(options, opslevel.SetAPIVisibility("internal"))
	client := opslevel.NewGQLClient(options...)

	clientErr := client.Validate()
	cobra.CheckErr(clientErr)

	return client
}
