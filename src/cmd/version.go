package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Build struct {
	Version         string          `json:"version,omitempty"`
	Commit          string          `json:"git,omitempty"`
	GoInfo          GoInfo          `json:"go,omitempty"`
	OpslevelVersion OpslevelVersion `json:"opslevel,omitempty"`
}

type OpslevelVersion struct {
	Commit  string `json:"app_commit,omitempty"`
	Version string `json:"app_version,omitempty"`
}

type GoInfo struct {
	Version  string `json:"version,omitempty"`
	Compiler string `json:"compiler,omitempty"`
	OS       string `json:"os,omitempty"`
	Arch     string `json:"arch,omitempty"`
}

var (
	version = "development"
	commit  = "none"
	build   Build
)

func initBuild() {
	build.Version = version
	if len(commit) >= 12 {
		build.Commit = commit[:12]
	} else {
		build.Commit = commit
	}

	build.GoInfo = getGoInfo()
	build.OpslevelVersion = getOpslevelVersion()
}

func getGoInfo() GoInfo {
	return GoInfo{
		Version:  runtime.Version(),
		Compiler: runtime.Compiler,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}
}

func getOpslevelVersion() OpslevelVersion {
	// Need to update all of this when we switch over to resty client
	url, err := url.Parse(viper.GetString("app-url"))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	url.Path = "/api/ping"
	response, err := http.Get(url.String())
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var opslevelVersion OpslevelVersion
	json.Unmarshal(responseData, &opslevelVersion)
	if len(opslevelVersion.Commit) >= 12 {
		opslevelVersion.Commit = opslevelVersion.Commit[:12]
	}

	return opslevelVersion
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version information`,
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	initBuild()
	versionInfo, err := json.MarshalIndent(build, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(versionInfo))
	return nil
}
