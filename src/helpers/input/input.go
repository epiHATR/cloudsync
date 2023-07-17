package input

import (
	"cloudsync/src/helpers/common"
	"cloudsync/src/helpers/output"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Get a Flag value in both of command flag and environment variable
func GetInputValue(flagName string, value string, cmd *cobra.Command) (string, error) {
	errText := ""
	if len(value) > 0 {
		return value, nil
	} else {
		flag := cmd.Flag(flagName)
		// if flag value was not provided or just an empty string
		if flag == nil || flag.Changed == false || len(flag.Value.String()) <= 0 {
			// then we'll check environment variables starts with CLOUDSYNC_ENV_
			flagEnvName := "CLOUDSYNC_ENV_" + strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_"))
			flagEnvValue := os.Getenv(flagEnvName)
			if len(flagEnvValue) <= 0 {
				// add to errText that no value found for that flag
				if len(flag.Shorthand) <= 0 {
					errText = fmt.Sprintf("--%s", flag.Name)
				} else {
					errText = fmt.Sprintf("--%s/-%s", flag.Name, flag.Shorthand)
				}
				return "", fmt.Errorf(errText)
			} else {
				return flagEnvName, nil
			}
		} else if flag.Changed {
			return flag.Value.String(), nil
		}
		return "", nil
	}
}

// Check if provided Flagset has value (both of input value and environment variables).
// Return error if one of flag in that Flagset has no value provided
func ValidateRequireFlags(flagSet []string, commandHelpText string, cmd *cobra.Command) error {
	errText := []string{}
	for _, flagName := range flagSet {
		_, err := GetInputValue(flagName, "", cmd)
		if err != nil {
			errText = append(errText, err.Error())
		}
	}

	// only print errText if it has any value
	if len(errText) > 0 {
		message := fmt.Sprintf("the following arguments are required: %s", strings.Join(errText, ", "))
		if len(commandHelpText) > 0 {
			// return error with command help text if provided
			return fmt.Errorf(fmt.Sprintf("%s\n\n%s", message, commandHelpText))
		} else {
			// otherwise, return error with errText only
			return fmt.Errorf(message)
		}
	} else {
		return nil
	}
}

// Check if a Flagset has one flag contains value. Return list of flags in specific flagset has data
func IsOneFlagValueProvided(flagSet []string, cmd *cobra.Command) ([]string, bool) {
	flag := []string{}
	for _, flagName := range flagSet {
		_, err := GetInputValue(flagName, "", cmd)
		if err == nil {
			flag = append(flag, flagName)
		}
	}
	if len(flag) > 0 {
		return flag, true
	}
	return []string{}, false
}

// Return the flag set and error(if failed) in the list of input flag set if it has value (in both of command flag value and environment value).
func GetActiveFlagSet(cmd *cobra.Command, cmdHelpText string, allFS ...[]string) ([]string, error) {
	output.PrintLog("getting active Flagset for command", cmd.Use)

	errorText := []string{}
	flagsHasValue := []string{}

	if len(allFS) <= 0 {
		return []string{}, fmt.Errorf("please provide at least one flag set.")
	}

	noFlagProvided := true

	for _, flagSet := range allFS {
		err := ValidateRequireFlags(flagSet, "", cmd)
		if err == nil {
			return flagSet, nil
		} else {
			flagName, isProvided := IsOneFlagValueProvided(flagSet, cmd)
			if isProvided && len(errorText) <= 0 {
				errorText = []string{err.Error()}
				noFlagProvided = false
				flagsHasValue = flagName
			}
		}
	}

	if noFlagProvided {
		output.PrintLog("no Flag provided in command", cmd.Use)
		errorText = []string{GetFlagsetString(allFS[0], true, *cmd)}
	} else {
		output.PrintLog("some Flags has value but other mandatory flags value was not provided.")
		errorText = []string{GetFlagsetString(common.GetShortestArray(flagsHasValue, true, allFS...), true, *cmd)}
	}

	// return if error
	if len(cmdHelpText) > 0 {
		return []string{}, fmt.Errorf("the following arguments are required: %s\n%s", strings.Join(errorText, ", "), cmdHelpText)
	} else {
		return []string{}, fmt.Errorf("the following arguments are required: %s", strings.Join(errorText, ", "))
	}
}

// Get all flag in a flagset printed as a string include both of flag and its short version if exists
func GetFlagsetString(flagSet []string, showShort bool, cmd cobra.Command) string {
	var output []string

	for _, flagName := range flagSet {
		flag := cmd.Flag(flagName)

		flagStr := fmt.Sprintf("--%s", flagName)
		if showShort && len(flag.Shorthand) > 0 {
			flagStr = fmt.Sprintf("--%s/-%s", flagName, flag.Shorthand)
		}

		output = append(output, flagStr)
	}

	return strings.Join(output, ", ")
}
