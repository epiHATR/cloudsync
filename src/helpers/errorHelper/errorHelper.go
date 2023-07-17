package errorHelper

import "cloudsync/src/helpers/output"

func Handle(err error) {
	if err != nil {
		output.PrintError(err.Error())
	}
}
