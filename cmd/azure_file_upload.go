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

var azFileUlAccount = ""
var azFileUlKey = ""
var azFileUlShareName = ""
var azFileUlUploadPath = ""

var azFileUlActiveFS []string = []string{}
var azFileUlKeyFS []string = []string{"account", "key", "share-name", "path"}

// azFileUlCmd represents the upload command
var azFileUlCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a file to Azure file shares.",
	Long:  "Upload a file to Azure file shares.",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintOut("LOGS", fmt.Sprintf("working on flagset %s", strings.Join(azFileUlActiveFS, ", ")))
		if common.IsSameArray(azFileUlKeyFS, azFileUlActiveFS) {
			azure.UploadToAzureFile(azFileUlAccount, azFileUlKey, azFileUlShareName, azFileUlUploadPath)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("account", cmd.Flags().Lookup("account"))
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("share-name", cmd.Flags().Lookup("share-name"))
		viper.BindPFlag("path", cmd.Flags().Lookup("path"))

		azFileUlAccount = viper.GetString("account")
		azFileUlKey = viper.GetString("key")
		azFileUlShareName = viper.GetString("share-name")
		azFileUlUploadPath = viper.GetString("path")
	},
	PreRun: func(cmd *cobra.Command, args []string) {

		activeFlagSet, err := input.GetActiveFlagSet(cmd, azFileUlKeyFS)
		if IsShowingExample {
			output.PrintFormat(text.Azure_FileShare_Upload_HelpText)
			os.Exit(0)
		} else {
			errorHelper.Handle(err, true)
		}
		azFileUlActiveFS = activeFlagSet
	},
}

func init() {
	azureFileCmd.AddCommand(azFileUlCmd)

	azFileUlCmd.Flags().StringVarP(&azFileUlAccount, "account", "a", "", "Name of storage account where file share located.")
	azFileUlCmd.Flags().StringVarP(&azFileUlKey, "key", "k", "", "Storage account authentiation key")
	azFileUlCmd.Flags().StringVarP(&azFileUlShareName, "share-name", "n", "", "Name of share file where file will be uploaded to.")
	azFileUlCmd.Flags().StringVarP(&azFileUlUploadPath, "path", "", "", "Folder/file path need to be uploaded.")
}
