/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudcync/src/helpers/input"
	"cloudcync/src/helpers/output"
	"cloudcync/src/helpers/provider/azure"

	"github.com/spf13/cobra"
)

var savePath string = "/tmp/"
var accountName string = ""
var containerName string = ""
var sasKey string = ""

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific container from Azure Storage Account",
	Long:  "Download a specific container from Azure Storage Account",
	Run: func(cmd *cobra.Command, args []string) {
		accountName, err := input.GetInputValue("account-name", accountName)
		if err != nil {
			output.PrintError(err.Error())
		}

		sasKey, err := input.GetInputValue("saskey", sasKey)
		if err != nil {
			output.PrintError(err.Error())
		}

		containerName, err := input.GetInputValue("container", containerName)
		if err != nil {
			output.PrintError(err.Error())
		}

		if len(accountName) > 0 && len(sasKey) > 0 && len(containerName) > 0 {
			err := azure.DownloadContainerToLocal(accountName, containerName, sasKey, savePath)
			if err != nil {
				output.PrintError(err.Error())
			}
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		requiredFlags := []string{"account-name", "container", "saskey"}
		output.PrintRequiredFlags(requiredFlags, cmd)
	},
}

func init() {
	containerCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&accountName, "account-name", "a", "", "Name of storage account where you want to get its container downloaded.")
	downloadCmd.Flags().StringVarP(&containerName, "container", "c", "", "Name of container you want to download.")
	downloadCmd.Flags().StringVarP(&sasKey, "saskey", "", "", "SAS key to access Azure storage account")
	downloadCmd.Flags().StringVarP(&savePath, "path-to-save", "p", "/temp/", "Location where contains and its blobs will be saved.")
}
