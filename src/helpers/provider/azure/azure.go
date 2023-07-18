package azure

import (
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/output"
	azurelib "cloudsync/src/library/azure"
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
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
	totalFile := 0
	sourceBlobs, err := azurelib.GetBlobsInContainer(*sourceClient, srcContainer)
	errorHelper.Handle(err)

	var wg sync.WaitGroup
	for _, blob := range sourceBlobs {
		wg.Add(1)
		go func(blob string) {
			defer wg.Done()

			azurelib.DownloadBlob(ctx, sourceClient, srcContainer, blob, "/tmp")
			filePath := "/tmp/" + blob
			file, _ := os.OpenFile(filePath, os.O_RDONLY, 0)
			defer file.Close()
			_, err = destClient.UploadFile(context.TODO(), destContainer, blob, file, nil)
			_ = os.Remove(filePath)
			output.PrintOut("INFO", "copied blob "+blob)
			totalFile = totalFile + 1
		}(blob)
	}

	wg.Wait()
	output.PrintOut("INFO", "total", strconv.Itoa(totalFile), "blobs was copied to", destContainer)
}
