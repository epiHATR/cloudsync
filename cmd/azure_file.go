/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var azureFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Working with Azure File shares.",
	Long:  "Working with Azure File shares.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
}

func init() {
	azureCmd.AddCommand(azureFileCmd)
}
