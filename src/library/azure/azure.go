package azurelib

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func ListBlobsUsingSAS(accountName, containerName, sasToken string) error {
	serviceUrl := fmt.Sprintf("https://%s.blob.core.windows.net", accountName)

	// Create a new credential with your SAS token
	credential, err := azblob.NewSharedKeyCredential(accountName, sasToken)
	if err != nil {
		return fmt.Errorf("failed to create shared key credential: %w", err)
	}
	client, err := azblob.NewClientWithSharedKeyCredential(serviceUrl, credential, nil)
	if err != nil {
		return err
	}

	fmt.Println("Listing the blobs in the container:")

	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return err
		}

		for _, blob := range resp.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}

	return nil
}
