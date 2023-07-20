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

var azCpSrcAccount string = ""
var azCpSrcContainer string = ""
var azCpSrcKey string = ""
var azCpSrcBlob string = ""
var azCpSrcConn string = ""

var azCpDestAccount string = ""
var azCpDestContainer string = ""
var azCpDestKey string = ""
var azCpDestConn string = ""

var cpActiveFS []string = []string{}
var azCpKeyFS []string = []string{"source-account", "source-container", "source-key", "destination-container"}
var azCpConnFS []string = []string{"source-connection-string", "source-container", "destination-connection-string", "destination-container"}

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
			azure.CopyContainerWithConnectionString(azCpSrcConn, azCpSrcContainer, azCpSrcBlob, azCpDestConn, azCpDestContainer)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("source-account", cmd.Flags().Lookup("source-account"))
		viper.BindPFlag("source-container", cmd.Flags().Lookup("source-container"))
		viper.BindPFlag("source-key", cmd.Flags().Lookup("source-key"))
		viper.BindPFlag("source-blob", cmd.Flags().Lookup("source-blob"))
		viper.BindPFlag("source-connection-string", cmd.Flags().Lookup("source-connection-string"))

		viper.BindPFlag("destination-account", cmd.Flags().Lookup("destination-account"))
		viper.BindPFlag("destination-container", cmd.Flags().Lookup("destination-container"))
		viper.BindPFlag("destination-key", cmd.Flags().Lookup("destination-key"))
		viper.BindPFlag("destination-blob", cmd.Flags().Lookup("destination-blob"))
		viper.BindPFlag("destination-connection-string", cmd.Flags().Lookup("destination-connection-string"))

		azCpSrcAccount = viper.GetString("source-account")
		azCpSrcContainer = viper.GetString("source-container")
		azCpSrcKey = viper.GetString("source-key")
		azCpSrcConn = viper.GetString("source-connection-string")
		azCpSrcBlob = viper.GetString("source-blob")

		azCpDestAccount = viper.GetString("destination-account")
		azCpDestContainer = viper.GetString("destination-container")
		azCpDestKey = viper.GetString("destination-key")
		azCpDestConn = viper.GetString("destination-connection-string")

		if len(azCpDestAccount) <= 0 {
			azCpDestAccount = azCpSrcAccount
		}

		if len(azCpDestContainer) <= 0 {
			azCpDestContainer = azCpSrcContainer
		}

		if len(azCpDestKey) <= 0 {
			azCpDestKey = azCpSrcKey
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		activeFlagSet, err := input.GetActiveFlagSet(cmd, azCpKeyFS, azCpConnFS)
		if IsShowingExample {
			output.PrintFormat(text.Azure_Container_Copy_HelpText)
			os.Exit(0)
		} else {
			errorHelper.Handle(err, true)
		}
		cpActiveFS = activeFlagSet
	},
}

func init() {

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
	azureContainerCmd.AddCommand(azCpCmd)
}
