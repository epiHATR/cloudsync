package azure

import (
	helpers "cloudsync/src/helpers/error"
	"cloudsync/src/helpers/output"
	azurelib "cloudsync/src/library/azure"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
)

func DownloadContainerWithConnectionString(connectionString, containerName, path string) {
	output.PrintLog(fmt.Sprintf("Start downloading blobs in container %s to %s", containerName, path))
	client, err := azurelib.VerifyStorageAccountWithConnectionString(connectionString)
	helpers.HandleError(err)

	azurelib.DownloadBlobs(client, containerName, path)

	output.PrintLog(fmt.Sprintf("All blobs in %s has been transferred to %s", containerName, path))
}

func DownloadContainerWithKey(accountName, containerName, key, path string) {
	output.PrintLog(fmt.Sprintf("Start downloading blobs in %s/%s to %s", containerName, accountName, path))
	client, err := azurelib.VerifyStorageAccountWithKey(accountName, key)
	helpers.HandleError(err)

	azurelib.DownloadBlobs(client, containerName, path)

	output.PrintLog(fmt.Sprintf("All blobs in %s has been transferred to %s", containerName, path))
}

func CopyContainerWithKey(srcAccount, srcContainer, srcKey, destAccount, destContainer, destKey string) {
	ctx := context.Background()
	if len(destContainer) <= 0 {
		destContainer = srcContainer
	}

	// Verify source storage account
	sourceClient, err := azurelib.VerifyStorageAccountWithKey(srcAccount, srcKey)
	helpers.HandleError(err)

	// Verify if destination container exists
	destClient, err := azurelib.VerifyStorageAccountWithKey(destAccount, destKey)
	helpers.HandleError(err)

	// Create destination container
	azurelib.CreateContainer(ctx, *destClient, destContainer)

	// Start copying blobs
	log.Println(fmt.Sprintf("Copying all blobs from storage account %s (container %s) to storage account %s (container %s)", srcAccount, srcContainer, destAccount, destContainer))
	totalFile := 0
	sourceBlobs, err := azurelib.GetBlobsInContainer(*sourceClient, srcContainer)
	helpers.HandleError(err)

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
			output.PrintLog("copied blob " + blob)
			totalFile = totalFile + 1
		}(blob)
	}

	wg.Wait()
	log.Println("total", totalFile, "blobs was copied to", destContainer)
}
