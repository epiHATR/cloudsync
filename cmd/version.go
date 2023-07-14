/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudcync/src/const/text"
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of the current CloudSync CLI",
	Long:  "Print version of the current CloudSync CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(text.CliVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
