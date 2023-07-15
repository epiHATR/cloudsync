package input

import (
	"cloudsync/src/helpers/output"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func GetInputValue(flagName string, value string) (string, error) {
	if len(value) > 0 {
		return value, nil
	} else {
		flagEnvName := "CLOUDSYNC_ENV_" + strings.ToUpper(strings.ReplaceAll(flagName, "-", "_"))
		flagEnvValue := os.Getenv(flagEnvName)
		if len(flagEnvValue) <= 0 {
			return "", fmt.Errorf("no value found for flag %s (also environment variable %s)", flagName, flagEnvName)
		} else {
			output.PrintLog(fmt.Sprintf("found environment variables %s with value %s", flagEnvName, flagEnvValue))
			return flagEnvValue, nil
		}
	}
}

func ValidateRequireFlags(flagSet []string, commandHelpText string, cmd *cobra.Command) error {
	errText := []string{}
	for _, flagName := range flagSet {
		flag := cmd.Flag(flagName)
		// if flag value was not provided
		if flag == nil || flag.Changed == false {
			// then we'll check environment variables starts with CLOUDSYNC_ENV_
			flagEnvName := "CLOUDSYNC_ENV_" + strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_"))
			flagEnvValue := os.Getenv(flagEnvName)
			if len(flagEnvValue) <= 0 {
				// add to errText that no value found for that flag
				if len(flag.Shorthand) <= 0 {
					errText = append(errText, fmt.Sprintf("--%s", flag.Name))
				} else {
					errText = append(errText, fmt.Sprintf("--%s/-%s", flag.Name, flag.Shorthand))
				}
			}
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
		// exit if Required Flags values not found
	} else {
		return nil
	}
}

func GetActiveFlagSet(cmd *cobra.Command, cmdHelpText string, flagSets ...[]string) (error, []string) {
	errorText := ""

	if len(flagSets) <= 0 {
		return fmt.Errorf("no flag set was provided"), []string{}
	}

	for _, flagSet := range flagSets {
		err := ValidateRequireFlags(flagSet, "", cmd)
		if err == nil {
			return nil, flagSet
		} else {
			errorText = err.Error()
		}
	}
	// return any error
	if len(cmdHelpText) > 0 {
		return fmt.Errorf("%s\n%s", errorText, cmdHelpText), []string{}
	} else {
		return fmt.Errorf(errorText), []string{}
	}
}

func sortStrings(arr []string) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

func AreFlagSetsEqual(arr1, arr2 []string) bool {
	// Check if the lengths of the arrays are equal
	if len(arr1) != len(arr2) {
		return false
	}

	// Sort both arrays to ensure consistent order
	// before comparing the elements
	sortedArr1 := make([]string, len(arr1))
	copy(sortedArr1, arr1)
	sortedArr2 := make([]string, len(arr2))
	copy(sortedArr2, arr2)
	sortStrings(sortedArr1)
	sortStrings(sortedArr2)

	// Compare each element of the sorted arrays
	for i := range sortedArr1 {
		if sortedArr1[i] != sortedArr2[i] {
			return false
		}
	}
	return true
}
