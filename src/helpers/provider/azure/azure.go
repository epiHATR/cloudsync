package azure

import (
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/file"
	"cloudsync/src/helpers/output"
	azurelib "cloudsync/src/library/azure"
	"context"
	"fmt"
	"os"
	"strconv"
)

// Download container in a specific Storage account to a specific path, given access by storage account connection string
func DownloadContainerWithConnectionString(connectionString, containerName, path string) {
	output.PrintOut("INFO", fmt.Sprintf("Start downloading blobs in container %s to %s", containerName, path))
	client, err := azurelib.VerifyStorageAccountWithConnectionString(connectionString)
	errorHelper.Handle(err)
	azurelib.DownloadBlobs(client, containerName, path)
	output.PrintOut("INFO", fmt.Sprintf("All blobs in %s has been transferred to %s", containerName, path))
}

// Download container in a specific Storage account to a specific path, given access by storage account key
func DownloadContainerWithKey(accountName, containerName, key, path string) {
	output.PrintOut("INFO", fmt.Sprintf("Start downloading blobs in %s/%s to %s", containerName, accountName, path))
	client, err := azurelib.VerifyStorageAccountWithKey(accountName, key)
	errorHelper.Handle(err)
	azurelib.DownloadBlobs(client, containerName, path)
	output.PrintOut("INFO", fmt.Sprintf("All blobs in %s has been transferred to %s", containerName, path))
}

// Copy a container from a specific Storage account to the other Storage Account, givent access by keys
func CopyContainerWithKey(srcAccount, srcContainer, srcKey, destAccount, destContainer, destKey string) {
	ctx := context.Background()
	if len(destContainer) <= 0 {
		destContainer = srcContainer
	}

	// Verify source storage account
	sourceClient, err := azurelib.VerifyStorageAccountWithKey(srcAccount, srcKey)
	errorHelper.Handle(err)

	// Verify if destination container exists
	destClient, err := azurelib.VerifyStorageAccountWithKey(destAccount, destKey)
	errorHelper.Handle(err)

	// Create destination container
	azurelib.CreateContainer(ctx, *destClient, destContainer)

	// Start copying blobs
	output.PrintOut("INFO", fmt.Sprintf("Copying all blobs from storage account %s (container %s) to storage account %s (container %s)", srcAccount, srcContainer, destAccount, destContainer))
	sourceBlobs, err := azurelib.GetBlobsInContainer(*sourceClient, srcContainer)
	errorHelper.Handle(err)
	totalFile := azurelib.CopyBlobs(ctx, sourceClient, destClient, srcContainer, destContainer, sourceBlobs)
	output.PrintOut("INFO", "total", strconv.Itoa(totalFile), "blobs was copied to", destContainer)
}

func UploadToContainerWithKey(accountName, containerName, key, pathToUpload string) {
	isFileUpload, err := file.IsFilePath(pathToUpload)
	destClient, err := azurelib.VerifyStorageAccountWithKey(accountName, key)
	errorHelper.Handle(err)

	if isFileUpload {
		output.PrintOut("INFO", fmt.Sprintf("Start uploading file to the container %s", containerName))
		blobName, err := file.GetFileNameFromPath(pathToUpload)
		errorHelper.Handle(err)
		file, _ := os.OpenFile(pathToUpload, os.O_RDONLY, 0)
		defer file.Close()
		_, err = destClient.UploadFile(context.TODO(), containerName, blobName, file, nil)
		errorHelper.Handle(err)
		output.PrintOut("INFO", fmt.Sprintf("uploaded file %s to blob %s", pathToUpload, blobName))

	} else {
		output.PrintOut("INFO", fmt.Sprintf("Start uploading folder to the container %s", containerName))
		fileList, err := file.GetFiles(pathToUpload)
		errorHelper.Handle(err)
		azurelib.UploadBlobs(fileList, pathToUpload, containerName, destClient)
	}
}
