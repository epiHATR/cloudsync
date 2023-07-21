package fileHelper

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
		errorHelper.Handle(err, false)
	}

	// Write the input string to the file
	err := ioutil.WriteFile(filePath, []byte(input), 0644)
	errorHelper.Handle(err, false)
	return nil
}

func SaveToLocalFile(content, filePath string) error {
	return SaveStringToFile(content, filePath)
}

func GetCurrentUserHomePath() (string, error) {
	currentUser, err := user.Current()
	errorHelper.Handle(err, false)
	return currentUser.HomeDir, nil
}

func GetPathType(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to identify path, %s", err)
	}

	if fileInfo.IsDir() {
		return "DIR", nil
	}

	return "FILE", nil
}

func GetFiles(path string) ([]string, error) {
	var fileNames []string

	// Check if the path exists
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read path information: %s", err)
	}

	// If the path is a file, return a single-element array with the file path
	if !fileInfo.IsDir() {
		return []string{path}, nil
	}

	// If the path is a directory, use filepath.Walk to get all files in the directory
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

func GetParentFolder(filePath string) (string, error) {
	// Getting the absolute path of the file to handle relative paths
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	// Getting the directory part of the absolute path
	parentDir := filepath.Dir(absPath)

	// Getting the last part of the parent directory, which is the folder name
	parentFolder := filepath.Base(parentDir)

	return parentFolder, nil
}

func GetRelativeDir(parentPath, fullPath string) (string, error) {
	relativePath := strings.TrimPrefix(fullPath, parentPath)
	// Trim any leading path separator (e.g., '/' or '\')
	relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator))
	relativeDir := filepath.Dir(relativePath)
	return relativeDir, nil
}
