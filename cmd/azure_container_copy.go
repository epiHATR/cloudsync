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

var cpActiveFS []string = []string{}

var src_account string = ""
var src_container string = ""
var src_key string = ""

var dest_account string = ""
var dest_container string = ""
var dest_key string = ""

var cpBaseFS []string = []string{"source-account", "source-container", "source-key", "destination-account", "destination-key"}

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy blobs and containers between Azure storage accounts.",
	Long:  "Copy blobs and containers between Azure storage accounts.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(fmt.Sprintf("working on flagset %s", strings.Join(cpActiveFS, ", ")))

		if input.AreFlagSetsEqual(cpActiveFS, cpBaseFS) {
			src_account, err := input.GetInputValue("source-account", src_account)
			src_key, err = input.GetInputValue("source-key", src_key)
			src_container, err = input.GetInputValue("source-container", src_container)
			dest_account, err := input.GetInputValue("source-account", dest_account)
			dest_key, err = input.GetInputValue("destination-key", dest_key)
			dest_container, err = input.GetInputValue("destination-container", dest_container)
			helpers.HandleError(err)

			azure.CopyContainerWithKey(src_account, src_container, src_key, dest_account, dest_container, dest_key)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// we need to verify what is the current cmd flag set user want to provided to the command
		err, flags := input.GetActiveFlagSet(cmd, text.Azure_Container_Copy_HelpText, cpBaseFS)
		//print if there's any error
		helpers.HandleError(err)
		cpActiveFS = flags
	},
}

func init() {
	azureContainerCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVarP(&src_account, "source-account", "", "", "Name of storage account where you want to get its container copied.")
	copyCmd.Flags().StringVarP(&src_container, "source-container", "", "", "Name of container you want to copy.")
	copyCmd.Flags().StringVarP(&src_key, "source-key", "", "", "Source storage account key.")

	copyCmd.Flags().StringVarP(&dest_account, "destination-account", "", "", "Name of storage account where you want to get its container copied.")
	copyCmd.Flags().StringVarP(&dest_container, "destination-container", "", "", "Name of storage account where you want to get its container copied.")
	copyCmd.Flags().StringVarP(&dest_key, "destination-key", "", "", "Destination storage account key.")
}
