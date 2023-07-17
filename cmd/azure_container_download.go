/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudsync/src/const/text"
	helpers "cloudsync/src/helpers/error"
	"cloudsync/src/helpers/input"
	"cloudsync/src/helpers/provider/azure"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var dldActiveFS []string = []string{}
var saveTo string = "/tmp/cloudsync/containers"
var accountName string = ""
var containerName string = ""
var key string = ""
var connectionString string = ""

var dldBaseFS []string = []string{"account-name", "container", "key"}
var dldConnStringFS []string = []string{"container", "connection-string"}

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a specific container from Azure Storage Account.",
	Long:  "Download a specific container from Azure Storage Account.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(fmt.Sprintf("working on flagset %s", strings.Join(dldActiveFS, ", ")))

		saveTo, err := input.GetInputValue("save-to", saveTo, cmd)
		helpers.HandleError(err)
		saveTo = saveTo + "/" + containerName

		if input.AreFlagSetsEqual(dldBaseFS, dldActiveFS) {
			// download container with account-name, container, key
			accountName, err := input.GetInputValue("account-name", accountName, cmd)
			containerName, err := input.GetInputValue("container", containerName, cmd)
			key, err := input.GetInputValue("key", key, cmd)
			helpers.HandleError(err)
			if len(accountName) > 0 && len(containerName) > 0 && len(key) > 0 {
				azure.DownloadContainerWithKey(accountName, containerName, key, saveTo)
			}
		} else if input.AreFlagSetsEqual(dldConnStringFS, dldActiveFS) {
			//download container with connection-string, container
			containerName, err := input.GetInputValue("container", containerName, cmd)
			connectionString, err := input.GetInputValue("connection-string", connectionString, cmd)
			helpers.HandleError(err)
			if len(containerName) > 0 && len(connectionString) > 0 {
				azure.DownloadContainerWithConnectionString(connectionString, containerName, saveTo)
			}
		} else {
			// we will define another download method here
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		flags, err := input.GetActiveFlagSet(cmd, text.Azure_Container_Download_HelpText, dldBaseFS, dldConnStringFS)
		//print if there's any error
		helpers.HandleError(err)
		dldActiveFS = flags
	},
}

func init() {
	azureContainerCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&accountName, "account-name", "a", "", "Name of storage account where you want to get its container downloaded.")
	downloadCmd.Flags().StringVarP(&key, "key", "", "", "Storage Account key to access Azure storage account")
	downloadCmd.Flags().StringVarP(&connectionString, "connection-string", "", "", "Storage account connection string")
	downloadCmd.Flags().StringVarP(&containerName, "container", "c", "", "Name of container you want to download.")
	downloadCmd.Flags().StringVarP(&saveTo, "save-to", "", saveTo, "Location where container and its blobs will be saved.")
}
