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

var cpActiveFS []string = []string{}

var azCpSrcAccount string = ""
var azCpSrcContainer string = ""
var azCpSrcKey string = ""
var azCpSrcBlob string = ""
var azCpSrcConn string = ""

var azCpDestAccount string = ""
var azCpDestContainer string = ""
var azCpDestKey string = ""
var azCpDestConn string = ""

var azCpKeyFS []string = []string{"source-account", "source-container", "source-key", "destination-account", "destination-key"}
var azCpConnFS []string = []string{"source-container", "source-connection-string", "destination-container", "destination-connection-string"}

// azCpCmd represents the copy command
var azCpCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy container/blob between Azure storage accounts.",
	Long:  "Copy container/blob between Azure storage accounts.",
	Run: func(cmd *cobra.Command, args []string) {
		output.PrintOut("LOGS", fmt.Sprintf("working on flagset %s", strings.Join(cpActiveFS, ", ")))

		if common.IsSameArray(cpActiveFS, azCpKeyFS) {
			//copy from source to destination with storage account key
			azure.CopyContainerWithKey(azCpSrcAccount, azCpSrcContainer, azCpSrcKey, azCpSrcBlob, azCpDestAccount, azCpDestContainer, azCpDestKey)
		} else if common.IsSameArray(cpActiveFS, azCpConnFS) {
			//copy from source to destination with connection string
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		flags, err := input.GetActiveFlagSet(cmd, azCpKeyFS, azCpConnFS)
		errorHelper.Handle(err, true, text.Azure_Container_Copy_HelpText)
		cpActiveFS = flags
	},
}

func init() {
	azureContainerCmd.AddCommand(azCpCmd)

	azCpCmd.Flags().StringVarP(&azCpSrcAccount, "source-account", "", "", "Name of storage account where you want to get its container copied.")
	azCpCmd.Flags().StringVarP(&azCpSrcContainer, "source-container", "", "", "Name of container you want to copy.")
	azCpCmd.Flags().StringVarP(&azCpSrcKey, "source-key", "", "", "Source storage account key.")
	azCpCmd.Flags().StringVarP(&azCpSrcBlob, "source-blob", "", "", "Blobs to copy, separated by commas ','.")
	azCpCmd.Flags().StringVarP(&azCpSrcConn, "source-connection-string", "", "", "Source storage account connection string.")

	azCpCmd.Flags().StringVarP(&azCpDestAccount, "destination-account", "", "", "Name of storage account where you want to get its container copied.")
	azCpCmd.Flags().StringVarP(&azCpDestContainer, "destination-container", "", "", "Name of storage account where you want to get its container copied.")
	azCpCmd.Flags().StringVarP(&azCpDestKey, "destination-key", "", "", "Destination storage account key.")
	azCpCmd.Flags().StringVarP(&azCpDestConn, "destination-connection-string", "", "", "Destination storage account connection string.")

	azCpCmd.Flags().SortFlags = false

	azureContainerCmd.SuggestFor = append(azureContainerCmd.SuggestFor, "copy")
}
