/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudsync/src/helpers/common"
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/input"
	"cloudsync/src/helpers/output"
	"cloudsync/src/helpers/provider/azure"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var azUploadCmdAccount string = ""
var azUploadCmdContainer string = ""
var azUploadCmdKey string = ""
var azUploadCmdConn string = ""
var azUploadCmdPath string = ""

var uldActiveFS []string = []string{}
var uldBaseFS []string = []string{"account-name", "container", "key", "path"}
var uldConnStringFS []string = []string{"container", "connection-string", "path"}

// uploadCmd represents the upload command
var azUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a folder/file to Azure storage account container. ",
	Long:  "Upload a folder/file to Azure storage account container. ",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintOut("LOGS", fmt.Sprintf("working on flagset %s", strings.Join(uldBaseFS, ", ")))

		if common.IsSameArray(uldBaseFS, uldActiveFS) {
			azure.UploadToContainerWithKey(azUploadCmdAccount, azUploadCmdContainer, azUploadCmdKey, azUploadCmdPath)
		} else {
		}

	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		flags, err := input.GetActiveFlagSet(cmd, "", uldBaseFS, uldConnStringFS)
		errorHelper.Handle(err)
		uldActiveFS = flags
	},
}

func init() {
	azureContainerCmd.AddCommand(azUploadCmd)
	azUploadCmd.Flags().StringVarP(&azUploadCmdAccount, "account-name", "a", "", "Name of storage account contains the uploading files/folder.")
	azUploadCmd.Flags().StringVarP(&azUploadCmdContainer, "container", "c", "", "Name of container.")
	azUploadCmd.Flags().StringVarP(&azUploadCmdKey, "key", "k", "", "Storage Account key to access Azure storage account.")
	azUploadCmd.Flags().StringVarP(&azUploadCmdConn, "connection-string", "", "", "Storage account connection string")
	azUploadCmd.Flags().StringVarP(&azUploadCmdPath, "path", "", saveTo, "Folder/file path need to be uploaded.")
	azUploadCmd.Flags().SortFlags = false
}
