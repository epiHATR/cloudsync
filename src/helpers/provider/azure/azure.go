package azure

import (
	"cloudcync/src/helpers/output"
	azurelib "cloudcync/src/library/azure"
)

func DownloadContainerToLocal(accountName, containerName, sasKey, path string) error {
	output.PrintLog("Start downloading blobs in container " + containerName + "(account: " + accountName + ") to path: " + path)
	err := azurelib.ListBlobsUsingSAS(accountName, containerName, sasKey)
	if err != nil {
		return err
	} else {
		return nil
	}
}
