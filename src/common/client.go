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
	timeout := time.Second * time.Duration(viper.GetInt("api-timeout"))
	client := opslevel.NewGQLClient(opslevel.SetAPIToken(viper.GetString("api-token")),
		opslevel.SetURL(viper.GetString("api-url")),
		opslevel.SetTimeout(timeout),
		opslevel.SetUserAgentExtra(fmt.Sprintf("cli-%s", version)))

	clientErr := client.Validate()
	cobra.CheckErr(clientErr)

	return client
}
