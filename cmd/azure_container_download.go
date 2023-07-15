/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudsync/src/const/text"
	helpers "cloudsync/src/helpers/error"
	"cloudsync/src/helpers/file"
	"cloudsync/src/helpers/input"
	"cloudsync/src/helpers/output"
	"cloudsync/src/helpers/provider/azure"

	"github.com/spf13/cobra"
)

var saveTo string = ""
var accountName string = ""
var containerName string = ""
var sasKey string = ""

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific container from Azure Storage Account.",
	Long:  "Download a specific container from Azure Storage Account.",
	Run: func(cmd *cobra.Command, args []string) {
		accountName, err := input.GetInputValue("account-name", accountName)
		helpers.HandleError(err)

		sasKey, err := input.GetInputValue("saskey", sasKey)
		helpers.HandleError(err)

		containerName, err := input.GetInputValue("container", containerName)
		helpers.HandleError(err)

		if len(accountName) > 0 && len(sasKey) > 0 && len(containerName) > 0 {
			if len(saveTo) <= 0 {
				homeDir, err := file.GetCurrentUserHomePath()
				helpers.HandleError(err)
				saveTo = homeDir + "/downloads/" + accountName + "/" + containerName
			}
			azure.DownloadContainerToLocal(accountName, containerName, sasKey, saveTo)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		requiredFlags := []string{"account-name", "container", "saskey"}
		output.PrintRequiredFlags(requiredFlags, text.Azure_Container_Download, cmd)
	},
}

func init() {
	containerCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&accountName, "account-name", "a", "", "Name of storage account where you want to get its container downloaded.")
	downloadCmd.Flags().StringVarP(&containerName, "container", "c", "", "Name of container you want to download.")
	downloadCmd.Flags().StringVarP(&sasKey, "saskey", "", "", "SAS key to access Azure storage account")
	downloadCmd.Flags().StringVarP(&saveTo, "save-to", "", "", "Location where contains and its blobs will be saved.")
}
