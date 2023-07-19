package azure

import (
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/file"
	"cloudsync/src/helpers/output"
	azurelib "cloudsync/src/library/azure"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Download container in a specific Storage account to a specific path, given access by storage account connection string
func DownloadContainerWithConnectionString(connectionString, containerName, blobPath, path string) {
	client, err := azurelib.VerifyStorageAccountWithConnectionString(connectionString)
	errorHelper.Handle(err, false)
	if len(blobPath) <= 0 {
		output.PrintOut("INFO", fmt.Sprintf("Start downloading blobs in container %s to %s with connection string", containerName, path))
		azurelib.DownloadBlobs(client, containerName, path)
		output.PrintOut("INFO", fmt.Sprintf("All blobs in %s has been transferred to %s", containerName, path))
	} else {
		output.PrintOut("INFO", fmt.Sprintf("Start downloading blob %s in container %s to %s with connection string", blobPath, containerName, path))
		azurelib.DownloadBlob(context.Background(), client, containerName, blobPath, path, true)
		output.PrintOut("INFO", fmt.Sprintf("Blob %s in container %s has been transferred to %s", blobPath, containerName, path))
	}
}

// Download container in a specific Storage account to a specific path, given access by storage account key
func DownloadContainerWithKey(accountName, containerName, key, blobPath, path string) {
	client, err := azurelib.VerifyStorageAccountWithKey(accountName, key)
	errorHelper.Handle(err, false)
	if len(blobPath) <= 0 {
		output.PrintOut("INFO", fmt.Sprintf("Start downloading blobs in container %s to %s with key", containerName, path))
		azurelib.DownloadBlobs(client, containerName, path)
		output.PrintOut("INFO", fmt.Sprintf("All blobs in container %s has been transferred to %s", containerName, path))
	} else {
		output.PrintOut("INFO", fmt.Sprintf("Start downloading blob %s in container %s to %s with key", blobPath, containerName, path))
		azurelib.DownloadBlob(context.Background(), client, containerName, blobPath, path, true)
		output.PrintOut("INFO", fmt.Sprintf("Blob %s in container %s has been transferred to %s", blobPath, containerName, path))
	}
}

// Copy a container from a specific Storage account to the other Storage Account, givent access by keys
func CopyContainerWithKey(srcAccount, srcContainer, srcKey, sourceBlobs, destAccount, destContainer, destKey string) {
	ctx := context.Background()
	if len(destContainer) <= 0 {
		destContainer = srcContainer
	}

	// Verify source storage account
	sourceClient, err := azurelib.VerifyStorageAccountWithKey(srcAccount, srcKey)
	errorHelper.Handle(err, false)

	// Verify if destination container exists
	destClient, err := azurelib.VerifyStorageAccountWithKey(destAccount, destKey)
	errorHelper.Handle(err, false)

	// Create destination container
	azurelib.CreateContainer(ctx, *destClient, destContainer)

	// Start copying blobs
	blobsToCopy := []string{}
	if len(sourceBlobs) <= 0 {
		blobs, err := azurelib.GetBlobsInContainer(*sourceClient, srcContainer)
		errorHelper.Handle(err, false)
		blobsToCopy = blobs
	} else {
		blobs := strings.Split(sourceBlobs, ",")
		for _, blob := range blobs {
			trimmed := strings.TrimLeft(blob, " ")
			trimmed = strings.TrimLeft(trimmed, " ")
			if len(trimmed) > 0 {
				blobsToCopy = append(blobsToCopy, trimmed)
			}
		}
	}

	output.PrintOut("INFO", fmt.Sprintf("Copying all blobs from storage account %s (container %s) to storage account %s (container %s)", srcAccount, srcContainer, destAccount, destContainer))
	totalFile := azurelib.CopyBlobs(ctx, sourceClient, destClient, srcContainer, destContainer, blobsToCopy)
	output.PrintOut("INFO", "total", strconv.Itoa(totalFile), "blobs was copied to", destContainer)
}

// Upload object to a Azure storage container with connection string authentiation method
func UploadToContainerWithConnectionString(containerName, connectionString, pathToUpload string) {
	uploadType, err := file.GetPathType(pathToUpload)
	errorHelper.Handle(err, false)
	destClient, err := azurelib.VerifyStorageAccountWithConnectionString(connectionString)
	errorHelper.Handle(err, false)
	UploadToContainer(uploadType, destClient, containerName, pathToUpload)
}

// Upload object to a Azure storage container with storage account key authentiation method
func UploadToContainerWithKey(accountName, containerName, key, pathToUpload string) {
	uploadType, err := file.GetPathType(pathToUpload)
	errorHelper.Handle(err, false)
	destClient, err := azurelib.VerifyStorageAccountWithKey(accountName, key)
	errorHelper.Handle(err, false)
	UploadToContainer(uploadType, destClient, containerName, pathToUpload)
}

// Upload an object in pathToUpload to a specific Azure storage container
func UploadToContainer(uploadType string, client *azblob.Client, containerName, pathToUpload string) {
	if uploadType == "FILE" {
		output.PrintOut("INFO", fmt.Sprintf("Start uploading file %s to the container %s", pathToUpload, containerName))
	} else {
		output.PrintOut("INFO", fmt.Sprintf("Start uploading folder to the container %s", containerName))
	}
	fileList, err := file.GetFiles(pathToUpload)
	errorHelper.Handle(err, false)
	azurelib.UploadBlobs(fileList, pathToUpload, containerName, client)
}

// Download container blob with Storage account key
func DownloadBlobWithKey(accountName, containerName, blobName, key, path string) {
	output.PrintOut("INFO", fmt.Sprintf("Start downloading blob %s in container %s to %s with KEY", blobName, containerName, path))
	client, err := azurelib.VerifyStorageAccountWithKey(accountName, key)
	errorHelper.Handle(err, false)
	azurelib.DownloadBlob(context.Background(), client, containerName, blobName, path, true)
	output.PrintOut("INFO", fmt.Sprintf("Blob %s in container %s has been transferred to %s", blobName, containerName, path))
}

// Download container blob with Storage account key
func DownloadBlobWithConnectionString(connectionString, containerName, blobName, key, path string) {
	output.PrintOut("INFO", fmt.Sprintf("Start downloading blob %s in container %s to %s with connection string", blobName, containerName, path))
	client, err := azurelib.VerifyStorageAccountWithConnectionString(connectionString)
	errorHelper.Handle(err, false)
	azurelib.DownloadBlob(context.Background(), client, containerName, blobName, path, true)
	output.PrintOut("INFO", fmt.Sprintf("Blob %s in container %s has been transferred to %s", blobName, containerName, path))
}
