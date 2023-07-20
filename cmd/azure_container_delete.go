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

var azDelAccountName = ""
var azDelContainer = ""
var azDelKey = ""
var azDelBlob = ""
var azDelConn = ""
var azDelForce = false
var azDelConfirmText string = "no"

var delActiveFS []string = []string{}
var delKeyFS []string = []string{"account", "container", "key"}
var delConnFS []string = []string{"container", "connection-string"}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete container/blob in an Azure storage account.",
	Long:  `Delete container/blob in an Azure storage account.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !azDelForce {
			for {
				fmt.Print("Deleting container blobs, please confirm your Delete action (yes/no): ")
				fmt.Scan(&azDelConfirmText)
				if strings.ToLower(azDelConfirmText) == "yes" || strings.ToLower(azDelConfirmText) == "no" {
					break
				}
			}
		} else {
			azDelConfirmText = "yes"
		}

		if azDelConfirmText == "yes" {
			if common.IsSameArray(delActiveFS, delKeyFS) {
				azure.DeleteBlobWithKey(azDelAccountName, azDelContainer, azDelKey, azDelBlob)
			} else {
				azure.DeleteBlobWithConnectionString(azDelConn, azDelContainer, azDelBlob)
			}
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		viper.BindPFlag("account", cmd.Flags().Lookup("account"))
		viper.BindPFlag("container", cmd.Flags().Lookup("container"))
		viper.BindPFlag("key", cmd.Flags().Lookup("key"))
		viper.BindPFlag("blob", cmd.Flags().Lookup("blob"))
		viper.BindPFlag("connection-string", cmd.Flags().Lookup("connection-string"))

		azDelAccountName = viper.GetString("account")
		azDelContainer = viper.GetString("container")
		azDelKey = viper.GetString("key")
		azDelBlob = viper.GetString("blob")
		azDelConn = viper.GetString("connection-string")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		activeFlagSet, err := input.GetActiveFlagSet(cmd, delKeyFS, delConnFS)
		if IsShowingExample {
			output.PrintFormat(text.Azure_Container_Delete_HelpText)
			os.Exit(0)
		} else {
			errorHelper.Handle(err, true)
		}
		delActiveFS = activeFlagSet
	},
}

func init() {
	azureContainerCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&azDelAccountName, "account", "a", "", "Name of storage account contain deleting objects.")
	deleteCmd.Flags().StringVarP(&azDelContainer, "container", "c", "", "Name of storage container.")
	deleteCmd.Flags().StringVarP(&azDelKey, "key", "k", "", "Storage account access key.")
	deleteCmd.Flags().StringVarP(&azDelBlob, "blob", "b", "", "Blob name to be deleted.")
	deleteCmd.Flags().StringVarP(&azDelConn, "connection-string", "", "", "Storage account connection string.")
	deleteCmd.Flags().BoolVarP(&azDelForce, "force", "", false, "force delete container/blob without confirmation.")
}
