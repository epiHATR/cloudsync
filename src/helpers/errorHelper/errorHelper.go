package errorHelper

import (
	"cloudsync/src/helpers/output"
	"fmt"
)

func Handle(err error, flagRequired bool, extraContent ...string) {
	if err != nil {
		errorText := err.Error()
		if len(extraContent) > 0 {
			for _, item := range extraContent {
				errorText = fmt.Sprintf("%s\n%s", errorText, item)
			}
		}
		if flagRequired {
			output.PrintOut("ERROR", "the following arguments are required:", errorText)
		} else {
			output.PrintOut("ERROR", errorText)
		}
	}
}
