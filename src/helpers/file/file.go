package file

import (
	"cloudsync/src/helpers/errorHelper"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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

func IsFilePath(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to identify path, %s", err)
	}

	if fileInfo.IsDir() {
		return false, nil
	}

	return true, nil
}

func GetFiles(path string) ([]string, error) {
	var fileNames []string

	// Check if the path exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("Path does not exist: %s", path)
	}

	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileNames = append(fileNames, filePath)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to read files: %s", err)
	}

	return fileNames, nil
}

func GetFilePathFromFolder(folderPath, filePath string) string {
	relPath, err := filepath.Rel(folderPath, filePath)
	if err != nil || strings.HasPrefix(relPath, "..") {
		return filePath
	}

	return relPath
}
