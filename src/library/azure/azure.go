package azurelib

import (
	"bytes"
	helpers "cloudsync/src/helpers/error"
	"cloudsync/src/helpers/file"
	"cloudsync/src/helpers/output"
	"context"
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func VerifyStorageAccountWithKey(accountName, key string) (*azblob.Client, error) {
	serviceUrl := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, key)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceUrl, credential, nil)
	helpers.HandleError(err)
	return client, nil
}

func VerifyStorageAccountWithConnectionString(connectionString string) (*azblob.Client, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	helpers.HandleError(err)
	return client, nil
}

func GetBlobsInContainer(client azblob.Client, containerName string) ([]string, error) {
	blobs := []string{}
	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		helpers.HandleError(err)

		for _, blob := range resp.Segment.BlobItems {
			blobs = append(blobs, *blob.Name)
		}
	}
	return blobs, nil
}

func CreateContainer(ctx context.Context, client azblob.Client, containerName string) error {
	isContainerExist := false
	pager := client.NewListContainersPager(&azblob.ListContainersOptions{
		Include: azblob.ListContainersInclude{},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		helpers.HandleError(err)

		for _, container := range resp.ContainerItems {
			if *container.Name == containerName {
				isContainerExist = true
			}
		}
	}
	if isContainerExist == false {
		_, err := client.CreateContainer(ctx, containerName, nil)
		helpers.HandleError(err)
	}
	return nil
}

func DownloadBlobs(client *azblob.Client, containerName, path string) {
	ctx := context.Background()

	blobs, err := GetBlobsInContainer(*client, containerName)
	helpers.HandleError(err)

	var wg sync.WaitGroup
	for _, blob := range blobs {
		wg.Add(1)
		go func(blobName string) {
			defer wg.Done()
			DownloadBlob(ctx, client, containerName, blobName, path)
			output.PrintLog("downloaded blob " + blobName)
		}(blob)
	}

	wg.Wait()
}

func DownloadBlob(ctx context.Context, client *azblob.Client, containerName, blobName string, path string) {
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
}
