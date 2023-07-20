/*
Copyright © 2022 Hai.Tran (github.com/epiHATR)
*/
package cmd

import (
	"cloudsync/src/const/text"
	"cloudsync/src/helpers/output"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var IsDebug bool = false

var (
	version     string = "v0.0.8"
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

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var IsShowingExample bool = false
var replacer = strings.NewReplacer("-", "_")

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&IsDebug, "debug", "", false, "show debugging information in output windows")
	rootCmd.PersistentFlags().BoolP("help", "", false, "show command help for instructions and examples")
	rootCmd.PersistentFlags().BoolVarP(&IsShowingExample, "example", "", IsShowingExample, "show command implementations & examples.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	configFileName := ".cloudsync"
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)
	viper.SetConfigName(configFileName)

	path := filepath.Join(home, configFileName)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	viper.SetEnvPrefix("CLOUDSYNC_ENV")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.WriteConfig()
	if IsDebug {
		output.IsDebug = true
	} else {
		output.IsDebug = false
	}
}
