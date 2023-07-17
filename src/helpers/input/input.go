package input

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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

func IsFlagSetHasData(flagSet []string, cmd *cobra.Command) ([]string, bool) {
	flag := []string{}
	for _, flagName := range flagSet {
		_, err := GetInputValue(flagName, "", cmd)
		if err == nil {
			flag = append(flag, flagName)
			return flag, true
		}
	}
	return flag, false
}

func GetActiveFlagSet(cmd *cobra.Command, cmdHelpText string, flagSets ...[]string) ([]string, error) {
	errorText := []string{}
	selectedFS := []string{}

	if len(flagSets) <= 0 {
		return []string{}, fmt.Errorf("no flag set was provided")
	}

	noFlagProvided := true

	for _, flagSet := range flagSets {
		err := ValidateRequireFlags(flagSet, "", cmd)
		if err == nil {
			return flagSet, nil
		} else {
			flagHasData, isData := IsFlagSetHasData(flagSet, cmd)
			if isData && len(errorText) <= 0 {
				errorText = []string{err.Error()}
				noFlagProvided = false
				selectedFS = flagHasData
			}
		}
	}

	if noFlagProvided {
		errorText = []string{PrintOutFlagAndShortFlag(flagSets[0], *cmd)}
	} else {
		errorText = []string{PrintOutFlagAndShortFlag(GetShortestArray(selectedFS, true, flagSets...), *cmd)}
	}

	// return if error
	if len(cmdHelpText) > 0 {
		return []string{}, fmt.Errorf("the following arguments are required: %s\n%s", strings.Join(errorText, ", "), cmdHelpText)
	} else {
		return []string{}, fmt.Errorf("the following arguments are required: %s", strings.Join(errorText, ", "))
	}
}

func SortStrings(arr []string) {
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

	SortStrings(sortedArr1)
	SortStrings(sortedArr2)

	// Compare each element of the sorted arrays
	for i := range sortedArr1 {
		if sortedArr1[i] != sortedArr2[i] {
			return false
		}
	}
	return true
}

func PrintOutFlagAndShortFlag(flagSet []string, cmd cobra.Command) string {
	output := []string{}
	for _, flagName := range flagSet {
		flag := cmd.Flag(flagName)
		if len(flag.Shorthand) > 0 {
			output = append(output, fmt.Sprintf("--%s/%s", flag.Name, flag.Shorthand))
		} else {
			output = append(output, fmt.Sprintf("--%s", flag.Name))
		}
	}

	return strings.Join(output, ", ")
}

func GetShortestArray(input []string, excludeInput bool, arrays ...[]string) []string {
	if len(input) == 0 {
		return nil
	}

	var shortestArray []string
	shortestLength := -1

	for _, arr := range arrays {
		for _, item := range input {
			if contains(arr, item) && (shortestLength == -1 || len(arr) < shortestLength) {
				shortestArray = arr
				shortestLength = len(arr)
			}
		}
	}

	if excludeInput {
		// Filter out input elements from the shortestArray
		filteredArray := make([]string, 0, len(shortestArray))
		for _, item := range shortestArray {
			if !contains(input, item) {
				filteredArray = append(filteredArray, item)
			}
		}
		return filteredArray
	}

	return shortestArray
}

func contains(arr []string, item string) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}
