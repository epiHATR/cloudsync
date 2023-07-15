package azure

import (
	"bytes"
	helpers "cloudsync/src/helpers/error"
	"cloudsync/src/helpers/file"
	"cloudsync/src/helpers/output"
	azurelib "cloudsync/src/library/azure"
	"context"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func DownloadContainerWithConnectionString(connectionString, containerName, path string) {
	output.PrintLog(fmt.Sprintf("Start downloading blobs in container %s to %s", containerName, path))
	client, err := azurelib.VerifySourceAccountWithConnectionString(connectionString)
	helpers.HandleError(err)

	downloadBlobs(client, containerName, path)

	output.PrintLog(fmt.Sprintf("All blobs in %s has been transferred to %s", containerName, path))
}

func DownloadContainerWithKey(accountName, containerName, key, path string) {
	output.PrintLog(fmt.Sprintf("Start downloading blobs in %s/%s to %s", containerName, accountName, path))
	client, err := azurelib.VerifySourceAccountWithKey(accountName, key)
	helpers.HandleError(err)

	downloadBlobs(client, containerName, path)

	output.PrintLog(fmt.Sprintf("All blobs in %s has been transferred to %s", containerName, path))
}

func downloadBlobs(client *azblob.Client, containerName, path string) {
	ctx := context.Background()

	blobs, err := azurelib.GetBlobsInContainer(*client, containerName)
	helpers.HandleError(err)

	var wg sync.WaitGroup
	for _, blob := range blobs {
		wg.Add(1)
		go func(blobName string) {
			defer wg.Done()
			fileName, _ := file.GetFileNameFromPath(blobName)
			output.PrintLog(fmt.Sprintf("transferring blob %s", fileName))

			get, err := client.DownloadStream(ctx, containerName, blobName, nil)
			helpers.HandleError(err)

			downloadedData := bytes.Buffer{}
			retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
			_, err = downloadedData.ReadFrom(retryReader)
			helpers.HandleError(err)

			err = retryReader.Close()
			helpers.HandleError(err)

			err = file.SaveToLocalFile(downloadedData.String(), fmt.Sprintf("%s/%s", path, blobName))
			helpers.HandleError(err)
		}(blob)
	}

	wg.Wait()
}
