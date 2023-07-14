package input

import (
	"fmt"
	"os"
	"strings"
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
			return flagEnvValue, nil
		}
	}
}
