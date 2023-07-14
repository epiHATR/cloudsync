/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of the current CloudSync CLI",
	Long:  "Print version of the current CloudSync CLI",
	Run: func(cmd *cobra.Command, args []string) {
		shortTag, _ := cmd.Flags().GetBool("short")
		if shortTag {
			fmt.Println(version + "." + build)
		} else {
			fmt.Println("Cloudsync CLI version ", version+"."+build)
			fmt.Println("Build", build)
			fmt.Println("Release Date", releaseDate)
			fmt.Println("Commit", commit)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("short", "s", false, "Display short description for current cloudsync CLI.")
}
