package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

const ImageSourceTypeAzureSAS ImageSourceType = "azure_sas"

func init() {
	RegisterSource(ImageSourceTypeAzureSAS, NewAzureSASImageSource)
}

type AzureSASImageSource struct {
	Config *SourceConfig
}

func NewAzureSASImageSource(config *SourceConfig) ImageSource {
	return &AzureSASImageSource{Config: config}
}

func (s *AzureSASImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodGet && parseAzureSASBlobURL(r) != ""
}

func (s *AzureSASImageSource) GetImage(r *http.Request) ([]byte, error) {
	sasURL := parseAzureSASBlobURL(r)
	sasURL, err := url.QueryUnescape(sasURL)
	if err != nil {
		return nil, fmt.Errorf("azure_sas: error reverting query: %w", err)
	}
	fmt.Printf("\n\nsas url: %s\n\n\n", sasURL)

	u, err := url.Parse(sasURL)
	if err != nil {
		return nil, fmt.Errorf("azure_sas: error parsing url: %w", err)
	}

	blobURL := azblob.NewBlobURL(
		*u,
		azblob.NewPipeline(
			azblob.NewAnonymousCredential(),
			azblob.PipelineOptions{},
		),
	)

	dlResp, err := blobURL.Download(r.Context(), 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, fmt.Errorf("azure_sas: error downloading blob: %w", err)
	}

	data := &bytes.Buffer{}
	bodyData := dlResp.Body(azblob.RetryReaderOptions{})
	defer bodyData.Close()

	if _, err := data.ReadFrom(bodyData); err != nil {
		return nil, fmt.Errorf("azure_sas: error reading data: %w", err)
	}

	return data.Bytes(), nil
}

func uploadBufferToAzureSAS(data []byte, sasURL string) error {
	sasURL, err := url.QueryUnescape(sasURL)
	if err != nil {
		return fmt.Errorf("azure_sas: error reverting query: %w", err)
	}

	u, err := url.Parse(sasURL)
	if err != nil {
		return fmt.Errorf("azure_sas: error parsing url: %w", err)
	}

	blobURL := azblob.NewBlobURL(
		*u,
		azblob.NewPipeline(
			azblob.NewAnonymousCredential(),
			azblob.PipelineOptions{},
		),
	).ToBlockBlobURL()

	if _, err := blobURL.Upload(
		context.Background(),
		bytes.NewReader(data),
		azblob.BlobHTTPHeaders{},
		azblob.Metadata{},
		azblob.BlobAccessConditions{},
	); err != nil {
		return fmt.Errorf("azure_sas: uploading image failed: %w", err)
	}

	return nil
}

func parseAzureSASBlobURL(request *http.Request) string {
	return request.URL.Query().Get("azureSASBlobURL")
}
