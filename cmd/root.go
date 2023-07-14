/*
Copyright © 2022 Hai.Tran (github.com/epiHATR)
*/
package cmd

import (
	"cloudsync/src/const/text"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-delve/delve/pkg/config"
	"github.com/spf13/cobra"
)

var isDebug bool = false

var (
	version     string = "v0.0"
	build       string = "#"
	commit      string = "#"
	releaseDate string = "#"
)

var rootCmd = &cobra.Command{
	Use:   "cloudsync",
	Short: "CloudSync CLI version",
	Long:  "CloudSync CLI ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(text.CloudSync)
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableSuggestions = false
	rootCmd.PersistentFlags().SortFlags = false

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&isDebug, "debug", "", false, "show debugging information in output windows")
	rootCmd.PersistentFlags().BoolP("help", "", false, "show command help for instructions and examples")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if isDebug {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(ioutil.Discard)
	}
	config.LoadConfig()
}
