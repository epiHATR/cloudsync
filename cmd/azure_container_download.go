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

var azDlAccountName string = ""
var azDlContainer string = ""
var azDlKey string = ""
var azDlConn string = ""
var azDlBlobPath string = ""
var azDlSaveTo string = "/tmp/cloudsync/containers"

var dldActiveFS []string = []string{}
var dldKeyFS []string = []string{"account", "container", "key"}
var dldConnFS []string = []string{"container", "connection-string"}

// azDlCmd represents the download command
var azDlCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific container/blob from Azure Storage Account.",
	Long:  "Download a specific container/blob from Azure Storage Account to local.",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintOut("LOGS", fmt.Sprintf("working on flagset %s", strings.Join(dldActiveFS, ", ")))
		if common.IsSameArray(dldKeyFS, dldActiveFS) {
			azure.DownloadContainerWithKey(azDlAccountName, azDlContainer, azDlKey, azDlBlobPath, azDlSaveTo)
		} else {
			azure.DownloadContainerWithConnectionString(azDlConn, azDlContainer, azDlBlobPath, azDlSaveTo)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("account", cmd.Flags().Lookup("account"))
		viper.BindPFlag("container", cmd.Flags().Lookup("container"))
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("blob", cmd.Flags().Lookup("blob"))
		viper.BindPFlag("connection-string", cmd.Flags().Lookup("connection-string"))
		viper.BindPFlag("save-to", cmd.Flags().Lookup("save-to"))

		azDlAccountName = viper.GetString("account")
		azDlContainer = viper.GetString("container")
		azDlKey = viper.GetString("key")
		azDlConn = viper.GetString("connection-string")
		azDlBlobPath = viper.GetString("blob")
		azDlSaveTo = viper.GetString("save-to")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command

		activeFlagSet, err := input.GetActiveFlagSet(cmd, dldKeyFS, dldConnFS)
		if IsShowingExample {
			output.PrintFormat(text.Azure_Container_Download_HelpText)
			os.Exit(0)
		} else {
			errorHelper.Handle(err, true)
		}
		dldActiveFS = activeFlagSet
		if len(azDlBlobPath) <= 0 {
			azDlSaveTo = azDlSaveTo + "/" + azDlContainer
		}
	},
}

func init() {
	azDlCmd.Flags().StringVarP(&azDlAccountName, "account", "a", "", "Name of storage account where you want to get its container downloaded.")
	azDlCmd.Flags().StringVarP(&azDlContainer, "container", "c", "", "Name of container you want to download.")
	azDlCmd.Flags().StringVarP(&azDlKey, "key", "k", "", "Storage Account key to access Azure storage account.")
	azDlCmd.Flags().StringVarP(&azDlConn, "connection-string", "", "", "Storage account connection string.")
	azDlCmd.Flags().StringVarP(&azDlBlobPath, "blob", "b", "", "Path to the blob you want to download.")
	azDlCmd.Flags().StringVarP(&azDlSaveTo, "save-to", "", azDlSaveTo, "Location where container and its blobs will be saved.")

	azureContainerCmd.AddCommand(azDlCmd)
	azDlCmd.Flags().SortFlags = false
}
