package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "opslevel",
	Short: "Opslevel Commandline Tools",
	Long:  `Opslevel Commandline Tools`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().String("logFormat", "TEXT", "overrides environment variable 'OPSLEVEL_LOGFORMAT' (options [\"JSON\", \"TEXT\"])")
	rootCmd.PersistentFlags().String("logLevel", "INFO", "overrides environment variable 'OPSLEVEL_LOGLEVEL' (options [\"ERROR\", \"WARN\", \"INFO\", \"DEBUG\"])")
	rootCmd.PersistentFlags().String("apiurl", "https://api.opslevel.com/graphql", "The OpsLevel API Url. Overrides environment variable 'OPSLEVEL_APITOKEN'")
	rootCmd.PersistentFlags().String("apitoken", "", "The OpsLevel API Token. Overrides environment variable 'OPSLEVEL_APITOKEN'")

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.BindEnv("apiurl", "OPSLEVEL_API_URL", "OL_API_URL")
	viper.BindEnv("apitoken", "OPSLEVEL_API_TOKEN", "OL_API_TOKEN", "OL_APITOKEN")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("OPSLEVEL")
	viper.AutomaticEnv()
	setupLogging()
}

func setupLogging() {
	logFormat := strings.ToLower(viper.GetString("logFormat"))
	logLevel := strings.ToLower(viper.GetString("logLevel"))

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
