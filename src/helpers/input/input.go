package input

import (
	"cloudsync/src/helpers/common"
	"cloudsync/src/helpers/output"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Check if provided Flagset has value (both of input value and environment variables).
// Return error if one of flag in that Flagset has no value provided
func ValidateRequireFlags(flagSet []string, cmd *cobra.Command) error {
	errText := []string{}
	for _, flagName := range flagSet {
		flagValue := viper.GetString(flagName)
		output.PrintOut("LOGS", (fmt.Sprintf("Flag %s -> %s", flagName, flagValue)))
		if len(flagValue) <= 0 {
			flag := cmd.Flag(flagName)
			if len(flag.Shorthand) <= 0 {
				errText = append(errText, fmt.Sprintf("--%s", flag.Name))
			} else {
				errText = append(errText, fmt.Sprintf("--%s/-%s", flag.Name, flag.Shorthand))
			}
		}
	}
	// only print errText if it has any value
	if len(errText) > 0 {
		message := fmt.Sprintf("the following arguments are required: %s", strings.Join(errText, ", "))
		// otherwise, return error with errText only
		return fmt.Errorf(message)
	} else {
		return nil
	}
}

// Check if a Flagset has one flag contains value. Return list of flags in specific flagset has data
func GetFlagsHasValueInFlagSet(flagSet []string, cmd *cobra.Command) []string {
	flags := []string{}
	for _, flagName := range flagSet {
		flagValue := viper.GetString(flagName)
		if len(flagValue) > 0 {
			flags = append(flags, flagName)
		}
	}
	if len(flags) > 0 {
		return flags
	}
	return []string{}
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

// Return the flag set and error(if failed) in the list of input flag set if it has value (in both of command flag value and environment value).
func GetActiveFlagSet(cmd *cobra.Command, flagSets ...[]string) ([]string, error) {
	errorText := []string{}
	flagsHasValue := []string{}
	noFlagProvided := true
	showFlagShorthand := true
	excludedInput := true
	maxRatio := float64(0)

	output.PrintOut("LOGS", "getting active Flagset for command", cmd.Use)
	for _, flagSet := range flagSets {
		err := ValidateRequireFlags(flagSet, cmd)
		if err == nil {
			return flagSet, nil
		} else {
			valuedFlags := GetFlagsHasValueInFlagSet(flagSet, cmd)
			currentRatio := float64(len(valuedFlags)) / float64(len(flagSet))
			if len(valuedFlags) > 0 {
				errorText = []string{err.Error()}
				noFlagProvided = false
				if currentRatio > maxRatio {
					maxRatio = currentRatio
					flagsHasValue = valuedFlags
				}
			}
		}
	}

	if noFlagProvided {
		errorText = []string{GetFlagsetString(common.GetShortestArray([]string{}, excludedInput, flagSets...), showFlagShorthand, *cmd)}
	} else {
		errorText = []string{GetFlagsetString(common.GetShortestArray(flagsHasValue, excludedInput, flagSets...), showFlagShorthand, *cmd)}
	}

	// return if error
	return []string{}, fmt.Errorf(strings.Join(errorText, ", "))
}
