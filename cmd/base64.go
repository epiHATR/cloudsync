/*
Copyright © 2023 Hai Tran <hidetran@gmail.com>
*/
package cmd

import (
	"cloudsync/src/helpers/input"
	"cloudsync/src/helpers/output"
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
			if err != nil {
				panic(err.Error())
			} else {
				rawDecodedText, err := base64.StdEncoding.DecodeString(inputString)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(rawDecodedText))
			}

		} else {
			inputString, err := input.GetInputValue("input", cmdInputString)
			if err != nil {
				panic(err.Error())
			} else {
				rawEncodedText := base64.StdEncoding.EncodeToString([]byte(inputString))
				if err != nil {
					panic(err)
				}
				fmt.Println(string(rawEncodedText))
			}
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		requiredFlags := []string{"input"}
		output.PrintRequiredFlags(requiredFlags, cmd)
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)
	base64Cmd.Flags().BoolVarP(&base64Decoded, "decode", "d", false, "Decode base64 input string to plain data.")
	base64Cmd.Flags().StringVarP(&cmdInputString, "input", "i", "", "String need to be encoded/decoded.")
}
