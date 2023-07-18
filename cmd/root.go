/*
Copyright © 2022 Hai.Tran (github.com/epiHATR)
*/
package cmd

import (
	"cloudsync/src/const/text"
	"cloudsync/src/helpers/output"
	"os"

	"github.com/spf13/cobra"
)

var IsDebug bool = false

var (
	version     string = "v0.0.5"
	build       string = "0"
	commit      string = "0"
	releaseDate string = "0000-00-00 00:00:00"
)

var rootCmd = &cobra.Command{
	Use:   "cloudsync",
	Short: "CloudSync CLI version",
	Long:  "CloudSync CLI ",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintFormat(text.CloudSync)
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().SortFlags = false
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	// Hide debug flag at root
	//rootCmd.PersistentFlags().MarkHidden("debug")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&IsDebug, "debug", "", false, "show debugging information in output windows")
	rootCmd.PersistentFlags().BoolP("help", "", false, "show command help for instructions and examples")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if IsDebug {
		output.IsDebug = true
	} else {
		output.IsDebug = false
	}
}
