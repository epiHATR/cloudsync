/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	helpers "cloudsync/src/helpers/error"
	"cloudsync/src/helpers/input"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

// base64Cmd represents the base64 command

var base64Decoded bool = false
var cmdInputString string = ""

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Base64 encoding & decoding",
	Long:  "Base64 encoding & decoding",
	Run: func(cmd *cobra.Command, args []string) {
		if base64Decoded {
			inputString, err := input.GetInputValue("input", cmdInputString)
			helpers.HandleError(err)

			rawDecodedText, err := base64.StdEncoding.DecodeString(inputString)
			helpers.HandleError(err)
			fmt.Println(string(rawDecodedText))

		} else {
			inputString, err := input.GetInputValue("input", cmdInputString)
			helpers.HandleError(err)
			rawEncodedText := base64.StdEncoding.EncodeToString([]byte(inputString))
			helpers.HandleError(err)
			fmt.Println(string(rawEncodedText))
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		requiredFlags := []string{"input"}
		err := input.ValidateRequireFlags(requiredFlags, "", cmd)
		helpers.HandleError(err)
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)
	base64Cmd.Flags().BoolVarP(&base64Decoded, "decode", "d", false, "Decode base64 input string to plain data.")
	base64Cmd.Flags().StringVarP(&cmdInputString, "input", "i", "", "String need to be encoded/decoded.")
}
