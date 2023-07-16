package output

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const colorError = "\033[0;31m"
const colorNone = "\033[0m"

func printLine(line, trimmedLine, color string) {
	fmt.Fprintf(os.Stderr, "%s%s%s\n", color, trimmedLine, colorNone)
}

func PrintFormat(input string) {
	const colorLink = "\033[36m"
	const colorGray = "\033[0;90m"
	const colorCommand = "\033[0;34m"

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		trimmedLine := strings.TrimLeft(line, "#>|~")
		if trimmedLine == "" {
			fmt.Println(line)
		} else {
			switch {
			case strings.HasPrefix(line, "~"):
				printLine(line, trimmedLine, colorLink)
			case strings.HasPrefix(line, "|"):
				printLine(line, trimmedLine, colorNone)
			case strings.HasPrefix(line, "#"):
				printLine(line, trimmedLine, colorGray)
			case strings.HasPrefix(line, ">"):
				printLine(line, trimmedLine, colorCommand)
			default:
				printLine(line, trimmedLine, colorError)
			}
		}
	}
}

func PrintError(input string) {
	PrintFormat(input)
	os.Exit(1)
}

func PrintLog(input string) {
	dt := time.Now()
	fmt.Println(fmt.Sprintf("%s %s", dt.Format("2006-01-02 15:04:05"), input))
}
