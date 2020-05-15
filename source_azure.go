package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

const ImageSourceTypeAzure ImageSourceType = "azure"

var (
	credential *azblob.TokenCredential
)

func init() {
	RegisterSource(ImageSourceTypeAzure, NewAzureImageSource)
}

func newAzureSession(container string) (*azblob.ContainerURL, error) {
	if credential == nil {
		if err := initAzure(); err != nil {
			return nil, err
		}
	}

	accountName := os.Getenv("AZURE_ACCOUNT_NAME")

	p := azblob.NewPipeline(*credential, azblob.PipelineOptions{})
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", accountName))
	containerURL := azblob.NewServiceURL(*u, p).NewContainerURL(container)

	return &containerURL, nil
}

type AzureImageSource struct {
	Config *SourceConfig
}

func NewAzureImageSource(config *SourceConfig) ImageSource {
	return &AzureImageSource{Config: config}
}

func (s *AzureImageSource) Matches(r *http.Request) bool {
	return r.Method == http.MethodGet && parseAzureBlobKey(r) != ""
}

func (s *AzureImageSource) GetImage(r *http.Request) ([]byte, error) {
	key, container := parseAzureBlobKey(r), parseAzureContainer(r)

	session, err := newAzureSession(container)
	if err != nil {
		return nil, fmt.Errorf("azure: error getting azure session: %w", err)
	}

	dlResp, err := session.NewBlobURL(key).
		Download(r.Context(), 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return nil, fmt.Errorf("azure: error downloading blob: %w", err)
	}

	data := &bytes.Buffer{}
	bodyData := dlResp.Body(azblob.RetryReaderOptions{})
	defer bodyData.Close()

	if _, err := data.ReadFrom(bodyData); err != nil {
		return nil, fmt.Errorf("azure: error reading data: %w", err)
	}

	return data.Bytes(), nil
}

func uploadBufferToAzure(data []byte, outputBlobKey, container string) error {
	session, err := newAzureSession(container)
	if err != nil {
		return fmt.Errorf("azure: error getting azure session: %w", err)
	}

	if _, err := session.
		NewBlockBlobURL(outputBlobKey).
		Upload(
			context.Background(),
			bytes.NewReader(data),
			azblob.BlobHTTPHeaders{},
			azblob.Metadata{},
			azblob.BlobAccessConditions{},
		); err != nil {
		return fmt.Errorf("azure: uploading image failed: %w", err)
	}

	return nil
}

func parseAzureBlobKey(request *http.Request) string {
	return request.URL.Query().Get("azureBlobKey")
}

func parseAzureBlobOutputKey(request *http.Request) string {
	return request.URL.Query().Get("azureOutputBlobKey")
}

func parseAzureContainer(request *http.Request) string {
	return request.URL.Query().Get("azureContainer")
}

func initAzure() error {
	azureEnv, err := azure.EnvironmentFromName("AZUREPUBLICCLOUD")
	if err != nil {
		return fmt.Errorf("azure/init: error getting environment from name: %s", err)
	}

	azureTenantID := os.Getenv("AZURE_TENANT_ID")
	azureOAuthConfig, err := adal.NewOAuthConfig(azureEnv.ActiveDirectoryEndpoint, azureTenantID)
	if err != nil {
		return fmt.Errorf("azure/init: error in new oauth config: %s", err)
	}

	if azureOAuthConfig.IsZero() {
		return fmt.Errorf("azure/init: error configuring oauth for tenant")
	}

	azureClientID := os.Getenv("AZURE_CLIENT_ID")
	azureClientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	spt, err := adal.NewServicePrincipalToken(
		*azureOAuthConfig,
		azureClientID,
		azureClientSecret,
		azureEnv.ResourceIdentifiers.Storage,
	)
	if err != nil {
		return fmt.Errorf("azure: error getting service principal auth: %w", err)
	}

	tokenRefresher := func(spt *adal.ServicePrincipalToken) func(credential azblob.TokenCredential) time.Duration {
		lock := sync.Mutex{}
		return func(credential azblob.TokenCredential) time.Duration {
			// This is possible data race for token refresh so we lock it down.
			lock.Lock()
			defer lock.Unlock()

			if err := spt.Refresh(); err != nil {
				panic(fmt.Errorf("azure: error refreshing token: %s", err))
			}

			token := spt.Token()
			credential.SetToken(token.AccessToken)

			expiresIn, err := token.ExpiresIn.Int64()
			if err != nil {
				panic(err)
			}

			// We reduce the given time by 2 minutes so instead that token is active 10min we
			// refresh it after 8.
			return time.Duration(expiresIn-2*60) * time.Second
		}
	}(spt)

	c := azblob.NewTokenCredential("", tokenRefresher)
	credential = &c

	return nil
}
