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
	"strings"

	"github.com/spf13/cobra"
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

		// handle save-to flag
		saveTo, err := input.GetInputValue("save-to", azDlSaveTo, cmd)
		errorHelper.Handle(err, false)
		if len(azDlBlobPath) <= 0 {
			saveTo = saveTo + "/" + azDlContainer
		}

		if common.IsSameArray(dldKeyFS, dldActiveFS) {
			azure.DownloadContainerWithKey(azDlAccountName, azDlContainer, azDlKey, azDlBlobPath, saveTo)
		} else if common.IsSameArray(dldConnFS, dldActiveFS) {
			azure.DownloadContainerWithConnectionString(azDlConn, azDlContainer, azDlBlobPath, saveTo)
		} else {
			// we will define another download method here
			output.PrintFormat(text.Azure_Container_Download_HelpText)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		flags, err := input.GetActiveFlagSet(cmd, dldKeyFS, dldConnFS)
		errorHelper.Handle(err, true, text.Azure_Container_Download_HelpText)
		dldActiveFS = flags
	},
}

func init() {
	azureContainerCmd.AddCommand(azDlCmd)
	azDlCmd.Flags().StringVarP(&azDlAccountName, "account", "a", "", "Name of storage account where you want to get its container downloaded.")
	azDlCmd.Flags().StringVarP(&azDlContainer, "container", "c", "", "Name of container you want to download.")
	azDlCmd.Flags().StringVarP(&azDlBlobPath, "blob", "b", azDlBlobPath, "Path to the blob you want to download.")
	azDlCmd.Flags().StringVarP(&azDlKey, "key", "k", "", "Storage Account key to access Azure storage account.")
	azDlCmd.Flags().StringVarP(&azDlConn, "connection-string", "", "", "Storage account connection string.")
	azDlCmd.Flags().StringVarP(&azDlSaveTo, "save-to", "", azDlSaveTo, "Location where container and its blobs will be saved.")
	azDlCmd.Flags().SortFlags = false
}
