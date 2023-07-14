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

func DownloadContainerToLocal(accountName, containerName, sasKey, path string) {
	output.PrintLog(fmt.Sprintf("Start downloading blobs in %s(account: %s) to %s", containerName, accountName, path))
	ctx := context.Background()
	client, err := azurelib.VerifySourceAccount(accountName, sasKey)
	helpers.HandleError(err)

	blobs, err := azurelib.GetBlobsInContainer(*client, containerName)
	helpers.HandleError(err)

	var wg sync.WaitGroup
	for _, blob := range blobs {
		wg.Add(1)
		go func(blobName string) {
			defer wg.Done()
			output.PrintLog(fmt.Sprintf("transfering blob %s", blobName))

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
	output.PrintLog(fmt.Sprintf("All blobs in %s has been transfered to %s", containerName, path))
}
