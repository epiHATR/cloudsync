/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudsync/src/const/text"
	"cloudsync/src/helpers/common"
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/input"
	"cloudsync/src/helpers/output"
	"cloudsync/src/helpers/provider/azure"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var azUplAccount string = ""
var azUplContainer string = ""
var azUplKey string = ""
var azUplConn string = ""
var azUplPath string = ""

var uldActiveFS []string = []string{}
var azUplKeyFS []string = []string{"account", "container", "key", "path"}
var azUplConnFS []string = []string{"connection-string", "container", "path"}

// uploadCmd represents the upload command
var azUplCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a folder/file to Azure storage account container. ",
	Long:  "Upload a folder/file to Azure storage account container. ",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintOut("LOGS", fmt.Sprintf("working on flagset %s", strings.Join(uldActiveFS, ", ")))

		if common.IsSameArray(azUplKeyFS, uldActiveFS) {
			azure.UploadToContainerWithKey(azUplAccount, azUplContainer, azUplKey, azUplPath)
		} else {
			azure.UploadToContainerWithConnectionString(azUplContainer, azUplConn, azUplPath)
		}

	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// bind viper command flags
		viper.BindPFlag("account", cmd.Flags().Lookup("account"))
		viper.BindPFlag("container", cmd.Flags().Lookup("container"))
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("connection-string", cmd.Flags().Lookup("connection-string"))
		viper.BindPFlag("path", cmd.Flags().Lookup("path"))

		azUplAccount = viper.GetString("account")
		azUplContainer = viper.GetString("container")
		azUplKey = viper.GetString("key")
		azUplConn = viper.GetString("connection-string")
		azUplPath = viper.GetString("path")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		activeFlagSet, err := input.GetActiveFlagSet(cmd, azUplKeyFS, azUplConnFS)
		if IsShowingExample {
			output.PrintFormat(text.Azure_Container_Upload_HelpText)
			os.Exit(0)
		} else {
			errorHelper.Handle(err, true)
		}
		uldActiveFS = activeFlagSet
	},
}

func init() {
	// describe supported flag for command
	azUplCmd.Flags().StringVarP(&azUplAccount, "account", "a", "", "Name of storage account contains the uploading files/folder.")
	azUplCmd.Flags().StringVarP(&azUplContainer, "container", "c", "", "Name of container.")
	azUplCmd.Flags().StringVarP(&azUplKey, "key", "k", "", "Storage Account key to access Azure storage account.")
	azUplCmd.Flags().StringVarP(&azUplConn, "connection-string", "", "", "Storage account connection string.")
	azUplCmd.Flags().StringVarP(&azUplPath, "path", "", "", "Folder/file path need to be uploaded.")

	//disable flag sorting
	azUplCmd.Flags().SortFlags = false
	azureContainerCmd.AddCommand(azUplCmd)
}
