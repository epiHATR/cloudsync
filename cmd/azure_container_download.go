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

var saveTo string = "/tmp/cloudsync/containers"
var accountName string = ""
var containerName string = ""
var key string = ""
var connectionString string = ""

var dldActiveFS []string = []string{}
var dldBaseFS []string = []string{"account-name", "container", "key"}
var dldConnStringFS []string = []string{"container", "connection-string"}

// azDownloadCmd represents the download command
var azDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific container from Azure Storage Account.",
	Long:  "Download a specific container(with child blobs) from Azure Storage Account.",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintOut("LOGS", fmt.Sprintf("working on flagset %s", strings.Join(dldActiveFS, ", ")))

		saveTo, err := input.GetInputValue("save-to", saveTo, cmd)
		errorHelper.Handle(err)
		saveTo = saveTo + "/" + containerName

		if common.IsSameArray(dldBaseFS, dldActiveFS) {
			azure.DownloadContainerWithKey(accountName, containerName, key, saveTo)
		} else if common.IsSameArray(dldConnStringFS, dldActiveFS) {
			azure.DownloadContainerWithConnectionString(connectionString, containerName, saveTo)
		} else {
			// we will define another download method here
			output.PrintFormat(text.Azure_Container_Download_HelpText)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		flags, err := input.GetActiveFlagSet(cmd, "", dldBaseFS, dldConnStringFS)
		errorHelper.Handle(err)
		dldActiveFS = flags
	},
}

func init() {
	azureContainerCmd.AddCommand(azDownloadCmd)
	azDownloadCmd.Flags().StringVarP(&accountName, "account-name", "a", "", "Name of storage account where you want to get its container downloaded.")
	azDownloadCmd.Flags().StringVarP(&containerName, "container", "c", "", "Name of container you want to download.")
	azDownloadCmd.Flags().StringVarP(&key, "key", "k", "", "Storage Account key to access Azure storage account")
	azDownloadCmd.Flags().StringVarP(&connectionString, "connection-string", "", "", "Storage account connection string")
	azDownloadCmd.Flags().StringVarP(&saveTo, "save-to", "", saveTo, "Location where container and its blobs will be saved.")
	azDownloadCmd.Flags().SortFlags = false
}
