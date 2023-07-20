package azurelib

import (
	"bytes"
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/file"
	"cloudsync/src/helpers/output"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Get a Storage Account client with key
func VerifyStorageAccountWithKey(accountName, key string) (*azblob.Client, error) {
	serviceUrl := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)
	credential, err := azblob.NewSharedKeyCredential(accountName, key)
	errorHelper.Handle(err, false)
	client, err := azblob.NewClientWithSharedKeyCredential(serviceUrl, credential, nil)
	errorHelper.Handle(err, false)
	return client, nil
}

// Get a Storage Account client with connection string
func VerifyStorageAccountWithConnectionString(connectionString string) (*azblob.Client, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	errorHelper.Handle(err, false)
	return client, nil
}

// Get blobs of specific container in a specific Storage Account
func GetBlobsInContainer(client azblob.Client, containerName string) ([]string, error) {
	blobs := []string{}
	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		errorHelper.Handle(err, false)

		for _, blob := range resp.Segment.BlobItems {
			blobs = append(blobs, *blob.Name)
		}
	}
	return blobs, nil
}

// Create a container in a specific Storage Account
func CreateContainer(ctx context.Context, client azblob.Client, containerName string) error {
	isContainerExist := false
	pager := client.NewListContainersPager(&azblob.ListContainersOptions{
		Include: azblob.ListContainersInclude{},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		errorHelper.Handle(err, false)

		for _, container := range resp.ContainerItems {
			if *container.Name == containerName {
				isContainerExist = true
			}
		}
	}
	if isContainerExist == false {
		_, err := client.CreateContainer(ctx, containerName, nil)
		errorHelper.Handle(err, false)
	}
	return nil
}

// Download all blobs in a specific Storage Account
func DownloadBlobs(client *azblob.Client, containerName, path string) {
	ctx := context.Background()

	blobs, err := GetBlobsInContainer(*client, containerName)
	errorHelper.Handle(err, false)

	var wg sync.WaitGroup
	for _, blob := range blobs {
		wg.Add(1)
		go func(blobName string) {
			defer wg.Done()
			DownloadBlob(ctx, client, containerName, blobName, path, true)
		}(blob)
	}
	wg.Wait()
}

// Download a specific blob in a Storage Account
func DownloadBlob(ctx context.Context, client *azblob.Client, containerName, blobName string, path string, showOutput bool) {
	get, err := client.DownloadStream(ctx, containerName, blobName, nil)
	errorHelper.Handle(err, false)

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	errorHelper.Handle(err, false)

	err = retryReader.Close()
	errorHelper.Handle(err, false)

	err = file.SaveToLocalFile(downloadedData.String(), fmt.Sprintf("%s/%s", path, blobName))
	errorHelper.Handle(err, false)
	if showOutput == true {
		output.PrintOut("INFO", fmt.Sprintf("downloaded blob %s", blobName))
	} else {
		output.PrintOut("LOGS", fmt.Sprintf("downloaded blob %s", blobName))
	}
}

// Copy blobs between Storage accounts
func CopyBlobs(ctx context.Context, sourceClient *azblob.Client, destClient *azblob.Client, srcContainer string, destContainer string, sourceBlobs []string) int {
	totalFile := 0

	var wg sync.WaitGroup
	for _, blob := range sourceBlobs {
		wg.Add(1)
		go func(blob string) {
			defer wg.Done()

			DownloadBlob(ctx, sourceClient, srcContainer, blob, "/tmp", false)
			filePath := "/tmp/" + blob
			file, _ := os.OpenFile(filePath, os.O_RDONLY, 0)
			defer file.Close()
			_, err := destClient.UploadFile(context.TODO(), destContainer, blob, file, nil)
			errorHelper.Handle(err, false)
			_ = os.Remove(filePath)
			output.PrintOut("INFO", "copied blob "+blob)
			totalFile = totalFile + 1
		}(blob)
	}
	wg.Wait()
	return totalFile
}

// Upload files to a specific container in a Storage Account
func UploadBlobs(fileList []string, fromPath string, toContainer string, client *azblob.Client) {
	output.PrintOut("LOGS", fmt.Sprintf("found %d files in the path %s", len(fileList), fromPath))

	pathInfo, err := os.Stat(fromPath)
	errorHelper.Handle(err, false)

	if pathInfo.IsDir() {
		selectedFolder := filepath.Base(fromPath)
		var wg sync.WaitGroup
		for _, filePath := range fileList {
			wg.Add(1)

			go func(filePath string) {
				defer wg.Done()
				blobName := file.GetFilePathFromFolder(fromPath, filePath)
				fileName, err := file.GetFileNameFromPath(filePath)
				errorHelper.Handle(err, false)

				file, _ := os.OpenFile(filePath, os.O_RDONLY, 0)
				defer file.Close()
				_, err = client.UploadFile(context.TODO(), toContainer, fmt.Sprintf("%s/%s", selectedFolder, blobName), file, nil)
				errorHelper.Handle(err, false)

				output.PrintOut("INFO", fmt.Sprintf("uploaded file %s to blob %s/%s", fileName, selectedFolder, blobName))
			}(filePath)
		}
		wg.Wait()
		output.PrintOut("INFO", fmt.Sprintf("Folder %s has been upload to %s", selectedFolder, toContainer))
	} else {
		blobName, err := file.GetFileNameFromPath(fromPath)
		file, _ := os.OpenFile(fromPath, os.O_RDONLY, 0)
		defer file.Close()
		_, err = client.UploadFile(context.TODO(), toContainer, blobName, file, nil)
		errorHelper.Handle(err, false)
		output.PrintOut("INFO", fmt.Sprintf("uploaded file %s to blob at %s/%s", fromPath, toContainer, blobName))
	}
}

// return Storage Account from SAS Url
func GetStorageAccountNameFromSasURL(sasURL string) (string, error) {
	parsedURL, err := url.Parse(sasURL)
	errorHelper.Handle(err, false)

	// The SAS token is the query part of the URL.
	queryParams, err := url.ParseQuery(parsedURL.RawQuery)
	errorHelper.Handle(err, false)

	// The "sv" query parameter contains the storage account name.
	storageAccountName := queryParams.Get("sv")
	if storageAccountName == "" {
		return "", fmt.Errorf("storage account name not found in the SAS URL")
	}

	// If the storage account name contains the path, remove it.
	storageAccountName = strings.Split(storageAccountName, ".")[0]

	return storageAccountName, nil
}

// delete blobs in a specific storage container
func DeleteContainerBlobs(ctx context.Context, client *azblob.Client, containerName string, blobs []string) error {
	for _, blob := range blobs {
		_, err := client.DeleteBlob(context.TODO(), containerName, blob, nil)
		if err != nil {
			return err
		} else {
			output.PrintOut("INFO", "deleted blob", blob)
		}
	}
	return nil
}
