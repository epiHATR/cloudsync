package output

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var IsDebug = false

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
				printLine(line, fmt.Sprintf("#%s", trimmedLine), colorGray)
			case strings.HasPrefix(line, ">"):
				printLine(line, trimmedLine, colorCommand)
			default:
				printLine(line, trimmedLine, colorError)
			}
		}
	}
}

func PrintOut(logType string, input ...string) {
	switch logType {
	case "INFO":
		{
			if IsDebug {
				logger := log.New(os.Stdout, "INFO ", 3)
				logger.Println(strings.Join(input, " "))
			} else {
				logger := log.New(os.Stdout, "", 0)
				logger.Println(strings.Join(input, " "))
			}
		}
	case "ERROR":
		{
			PrintFormat(fmt.Sprintf("%s", strings.Join(input, " ")))
			os.Exit(1)
		}
	case "LOGS":
		{
			if IsDebug {
				logger := log.New(os.Stdout, "LOGS ", 3)
				logger.Println(fmt.Sprintf("%s", strings.Join(input, " ")))
			}
		}
	default:
		{
			if IsDebug {
				logger := log.New(os.Stdout, "INFO ", 3)
				logger.Println(strings.Join(input, " "))
			}
		}
	}
}
