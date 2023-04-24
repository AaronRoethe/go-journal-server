package storage

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func SaveMessageToBlob(message []byte) error {

	accountName := "devstoreaccount1"
	accountKey := "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="
	containerName := "journal-test"
	blobName := "test.txt"

	// Create a credential object using account name and account key
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Failed to create a credential: %v\n", err)
		return err
	}

	// Create a pipeline object using the credential
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a container URL object
	urlStr := fmt.Sprintf("http://127.0.0.1:10000/devstoreaccount1/%s", containerName)
	url, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Failed to parse container URL: %v\n", err)
		return err
	}
	containerURL := azblob.NewContainerURL(*url, p)

	// Create a blob URL object
	blobURL := containerURL.NewBlockBlobURL(blobName)

	// Create the container (if it doesn't already exist)
	_, err = containerURL.Create(context.Background(), azblob.Metadata{}, azblob.PublicAccessNone)
	if err != nil {
		if stgErr, ok := err.(azblob.StorageError); ok && stgErr.ServiceCode() == azblob.ServiceCodeContainerAlreadyExists {
			log.Printf("Container already exists: %s\n", containerName)
		} else {
			log.Printf("Failed to create container: %v\n", err)
			return err
		}
	}

	// Upload the message to the blob
	_, err = azblob.UploadBufferToBlockBlob(context.Background(), message, blobURL, azblob.UploadToBlockBlobOptions{})
	if err != nil {
		log.Printf("Failed to upload message to blob: %v\n", err)
		return err
	}

	log.Printf("Successfully saved message to blob: %s\n", message)
	return nil
}
