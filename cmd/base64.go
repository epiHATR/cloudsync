/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudsync/src/const/text"
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/input"
	"cloudsync/src/helpers/output"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// base64Cmd represents the base64 command

var base64Decoded bool = false
var cmdInputString string = ""

var requiredFlags = []string{"input"}

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Base64 encoding & decoding",
	Long:  "Base64 encoding & decoding",
	Run: func(cmd *cobra.Command, args []string) {
		if base64Decoded {
			rawDecodedText, err := base64.StdEncoding.DecodeString(cmdInputString)
			errorHelper.Handle(err, false)
			fmt.Println(string(rawDecodedText))

		} else {
			rawEncodedText := base64.StdEncoding.EncodeToString([]byte(cmdInputString))
			fmt.Println(string(rawEncodedText))
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		cmdInputString = viper.GetString("input")
		err := input.ValidateRequireFlags(requiredFlags, cmd)
		if IsShowingExample {
			output.PrintFormat(text.Base64_HelpText)
			os.Exit(0)
		} else {
			errorHelper.Handle(err, true)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)
	base64Cmd.Flags().BoolVarP(&base64Decoded, "decode", "d", false, "Decode base64 input string to plain data.")
	base64Cmd.Flags().StringVarP(&cmdInputString, "input", "i", "", "String need to be encoded/decoded.")

	viper.BindPFlag("input", base64Cmd.Flags().Lookup("input"))
}
