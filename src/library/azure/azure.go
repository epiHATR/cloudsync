package azurelib

import (
	"bytes"
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/file"
	"cloudsync/src/helpers/output"
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func VerifyStorageAccountWithKey(accountName, key string) (*azblob.Client, error) {
	serviceUrl := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, key)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceUrl, credential, nil)
	errorHelper.Handle(err)
	return client, nil
}

func VerifyStorageAccountWithConnectionString(connectionString string) (*azblob.Client, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	errorHelper.Handle(err)
	return client, nil
}

func GetBlobsInContainer(client azblob.Client, containerName string) ([]string, error) {
	blobs := []string{}
	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		errorHelper.Handle(err)

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
		errorHelper.Handle(err)

		for _, container := range resp.ContainerItems {
			if *container.Name == containerName {
				isContainerExist = true
			}
		}
	}
	if isContainerExist == false {
		_, err := client.CreateContainer(ctx, containerName, nil)
		errorHelper.Handle(err)
	}
	return nil
}

func DownloadBlobs(client *azblob.Client, containerName, path string) {
	ctx := context.Background()

	blobs, err := GetBlobsInContainer(*client, containerName)
	errorHelper.Handle(err)

	var wg sync.WaitGroup
	for _, blob := range blobs {
		wg.Add(1)
		go func(blobName string) {
			defer wg.Done()
			DownloadBlob(ctx, client, containerName, blobName, path)
		}(blob)
	}
	wg.Wait()
}

func DownloadBlob(ctx context.Context, client *azblob.Client, containerName, blobName string, path string) {
	get, err := client.DownloadStream(ctx, containerName, blobName, nil)
	errorHelper.Handle(err)

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	errorHelper.Handle(err)

	err = retryReader.Close()
	errorHelper.Handle(err)

	err = file.SaveToLocalFile(downloadedData.String(), fmt.Sprintf("%s/%s", path, blobName))
	errorHelper.Handle(err)
	output.PrintOut("LOGS", fmt.Sprintf("downloaded blob %s", blobName))
}

func CopyBlobs(ctx context.Context, sourceClient *azblob.Client, destClient *azblob.Client, srcContainer string, destContainer string, sourceBlobs []string) int {
	totalFile := 0

	var wg sync.WaitGroup
	for _, blob := range sourceBlobs {
		wg.Add(1)
		go func(blob string) {
			defer wg.Done()

			DownloadBlob(ctx, sourceClient, srcContainer, blob, "/tmp")
			filePath := "/tmp/" + blob
			file, _ := os.OpenFile(filePath, os.O_RDONLY, 0)
			defer file.Close()
			_, err := destClient.UploadFile(context.TODO(), destContainer, blob, file, nil)
			errorHelper.Handle(err)
			_ = os.Remove(filePath)
			output.PrintOut("INFO", "copied blob "+blob)
			totalFile = totalFile + 1
		}(blob)
	}
	wg.Wait()
	return totalFile
}
