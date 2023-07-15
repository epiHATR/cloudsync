package helpers

import (
	"cloudsync/src/helpers/output"
)

func HandleError(err error) {
	if err != nil {
		output.PrintError(err.Error())
	}
}
