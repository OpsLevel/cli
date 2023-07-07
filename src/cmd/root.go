package cmd

import (
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var (
	_clientRest *resty.Client
	_clientGQL  *opslevel.Client
)

var rootCmd = &cobra.Command{
	Use:   "opslevel",
	Short: "Opslevel Commandline Tool",
	Long:  `Opslevel Commandline Tool`,
}

func Execute(v string, currentCommit string) {
	version = v
	commit = currentCommit
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().String("log-format", "TEXT", "overrides environment variable 'OPSLEVEL_LOG_FORMAT' (options [\"JSON\", \"TEXT\"])")
	rootCmd.PersistentFlags().String("log-level", "INFO", "overrides environment variable 'OPSLEVEL_LOG_LEVEL' (options [\"ERROR\", \"WARN\", \"INFO\", \"DEBUG\"])")
	rootCmd.PersistentFlags().String("api-url", "https://app.opslevel.com", "The OpsLevel API Url. Overrides environment variable 'OPSLEVEL_API_URL'")
	rootCmd.PersistentFlags().String("api-token", "", "The OpsLevel API Token. Overrides environment variable 'OPSLEVEL_API_TOKEN'")
	rootCmd.PersistentFlags().Bool("no-headers", false, "If --output=text and this flag is set the headers will be skip from being output")
	rootCmd.PersistentFlags().Lookup("no-headers").NoOptDefVal = "true"
	rootCmd.PersistentFlags().Int("api-timeout", 10, "The number of seconds to timeout of the request. Overrides environment variable 'OPSLEVEL_API_TIMEOUT'")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindEnv("log-format", "OPSLEVEL_LOG_FORMAT", "OL_LOG_FORMAT", "OL_LOGFORMAT")
	viper.BindEnv("log-level", "OPSLEVEL_LOG_LEVEL", "OL_LOG_LEVEL", "OL_LOGLEVEL")
	viper.BindEnv("api-url", "OPSLEVEL_API_URL", "OL_API_URL", "OPSLEVEL_APP_URL", "OL_APP_URL")
	viper.BindEnv("api-token", "OPSLEVEL_API_TOKEN", "OL_API_TOKEN", "OL_APITOKEN")
	viper.BindEnv("api-timeout", "OPSLEVEL_API_TIMEOUT")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("OPSLEVEL")
	viper.AutomaticEnv()
	setupLogging()
}

func setupLogging() {
	logFormat := strings.ToLower(viper.GetString("log-format"))
	logLevel := strings.ToLower(viper.GetString("log-level"))

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if logFormat == "text" {
		output := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(output)
	}

	switch {
	case logLevel == "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case logLevel == "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case logLevel == "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

func getClientRest() *resty.Client {
	if _clientRest == nil {
		_clientRest = opslevel.NewRestClient(opslevel.SetURL(viper.GetString("api-url")))
	}
	return _clientRest
}

func getClientGQL(options ...opslevel.Option) *opslevel.Client {
	if _clientGQL == nil {
		_clientGQL = common.NewGraphClient(version, options...)
	}
	return _clientGQL
}
