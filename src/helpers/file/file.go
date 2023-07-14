package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func SaveStringToFile(input, filePath string) error {
	// Check if the file already exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		dirPath := filepath.Dir(filePath)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
	}

	// Write the input string to the file
	err := ioutil.WriteFile(filePath, []byte(input), 0644)
	if err != nil {
		return err
	}
	return nil
}

func SaveToLocalFile(content, filePath string) error {
	return SaveStringToFile(content, filePath)
}
