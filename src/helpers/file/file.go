package file

import (
	"cloudsync/src/helpers/errorHelper"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func GetFileNameFromPath(path string) (string, error) {
	_, fileName := filepath.Split(path)
	return fileName, nil
}

func SaveStringToFile(input, filePath string) error {
	// Check if the file already exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		dirPath := filepath.Dir(filePath)
		err := os.MkdirAll(dirPath, 0755)
		errorHelper.Handle(err)
	}

	// Write the input string to the file
	err := ioutil.WriteFile(filePath, []byte(input), 0644)
	errorHelper.Handle(err)
	return nil
}

func SaveToLocalFile(content, filePath string) error {
	return SaveStringToFile(content, filePath)
}

func GetCurrentUserHomePath() (string, error) {
	currentUser, err := user.Current()
	errorHelper.Handle(err)
	return currentUser.HomeDir, nil
}
