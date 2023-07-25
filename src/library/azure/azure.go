package azurelib

import (
	"bytes"
	"cloudsync/src/helpers/errorHelper"
	"cloudsync/src/helpers/fileHelper"
	"cloudsync/src/helpers/output"
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-storage-file-go/azfile"
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

	err = fileHelper.SaveToLocalFile(downloadedData.String(), fmt.Sprintf("%s/%s", path, blobName))
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
				blobName := fileHelper.GetFilePathFromFolder(fromPath, filePath)
				fileName, err := fileHelper.GetFileNameFromPath(filePath)
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
		blobName, err := fileHelper.GetFileNameFromPath(fromPath)
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

func GetServiceUrl(accountName, accountKey string) azfile.ServiceURL {
	credential, err := azfile.NewSharedKeyCredential(accountName, accountKey)
	errorHelper.Handle(err, false)
	serviceURL, err := url.Parse(fmt.Sprintf("https://%s.file.core.windows.net/", accountName))
	errorHelper.Handle(err, false)
	p := azfile.NewPipeline(credential, azfile.PipelineOptions{})
	serviceUrl := azfile.NewServiceURL(*serviceURL, p)
	return serviceUrl
}

func GetShareUrl(accountName, accountKey, shareName string) azfile.ShareURL {
	serviceUrl := GetServiceUrl(accountName, accountKey)
	shareUrl := serviceUrl.NewShareURL(shareName)
	return shareUrl
}

func GetDirectoryUrl(accountName, accountKey, shareName, directoryName string) azfile.DirectoryURL {
	shareUrl := GetShareUrl(accountName, accountKey, shareName)
	directoryUrl := shareUrl.NewDirectoryURL(directoryName)
	return directoryUrl
}

func GetFileUrl(accountName, accountKey, shareName, directoryName, fileName string) azfile.FileURL {
	directoryUrl := GetDirectoryUrl(accountName, accountKey, shareName, directoryName)
	fileUrl := directoryUrl.NewFileURL(fileName)
	return fileUrl
}

func CreateFileShareDirectory(accountName, accountKey, shareName, shareLocation string) error {
	output.PrintOut("LOGS", fmt.Sprintf("creating folder at %s", shareLocation))

	// Split the shareLocation by '/'
	directories := strings.Split(shareLocation, "/")

	// Removing any empty strings that might result from leading/trailing slashes or double slashes
	var cleanedDirectories []string
	for _, dir := range directories {
		if dir != "" {
			cleanedDirectories = append(cleanedDirectories, dir)
		}
	}

	// Initialize the parent directory URL
	parentDirectoryUrl := GetDirectoryUrl(accountName, accountKey, shareName, "")

	// Check if the parent directory exists
	_, err := parentDirectoryUrl.GetProperties(context.TODO())
	if err != nil {
		// If the parent directory doesn't exist, return an error
		return fmt.Errorf("parent directory does not exist: %s", parentDirectoryUrl.String())
	}

	// Create parent and child directories
	for i, dir := range cleanedDirectories {
		// Append the current directory to the parent URL
		parentDirectoryUrl = parentDirectoryUrl.NewDirectoryURL(dir)

		// Check if the current directory exists
		_, err := parentDirectoryUrl.GetProperties(context.TODO())
		if err == nil {
			// If the current directory already exists, move to the next iteration
			continue
		}

		// Create the directory if it doesn't exist
		_, err = parentDirectoryUrl.Create(context.TODO(), azfile.Metadata{}, azfile.SMBProperties{})
		if err != nil {
			return err
		}

		// If this is the last directory, set the parent URL to the current URL
		// to ensure the next iteration starts from the correct parent.
		if i == len(cleanedDirectories)-1 {
			parentDirectoryUrl = parentDirectoryUrl.NewDirectoryURL("")
		}
	}

	return nil
}

func CheckShareExistance(accountName, accountKey, shareName string) error {
	shareUrl := GetShareUrl(accountName, accountKey, shareName)
	_, err := shareUrl.GetProperties(context.TODO())
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("No share called %s exists in storage account %s", shareName, accountName))
	} else {
		return nil
	}

}

func UploadToAzureFile(accountName, accountKey, shareName, fileToUpload, shareLocation string) {
	file, err := os.Open(fileToUpload)
	defer file.Close()
	_, err = file.Stat()
	errorHelper.Handle(err, false)

	shareFileLocation := ""
	fileName, _ := fileHelper.GetFileNameFromPath(fileToUpload)
	if len(shareLocation) <= 0 {
		shareFileLocation, _ = fileHelper.GetFileNameFromPath(fileToUpload)
	} else {
		shareFileLocation = fmt.Sprintf("%s/%s", shareLocation, fileName)
	}

	output.PrintOut("LOGS", fmt.Sprintf("uploading file %s to file share %s at path %s", fileToUpload, shareName, shareFileLocation))
	if len(shareFileLocation) > 0 {
		err := CreateFileShareDirectory(accountName, accountKey, shareName, shareLocation)
		errorHelper.Handle(err, false)
	}

	fileShareUrl := GetFileUrl(accountName, accountKey, shareName, shareLocation, fileName)
	err = azfile.UploadFileToAzureFile(context.TODO(), file, fileShareUrl,
		azfile.UploadToAzureFileOptions{
			FileHTTPHeaders: azfile.FileHTTPHeaders{
				CacheControl: "no-transform",
			},
		})
	errorHelper.Handle(err, false)
	output.PrintOut("INFO", fmt.Sprintf("uploaded file %s to file share path %s", fileToUpload, shareFileLocation))
}
