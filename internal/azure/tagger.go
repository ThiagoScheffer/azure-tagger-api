package azure

import (
	"context"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

var ErrMissingSubscription = errors.New("AZURE_SUBSCRIPTION_ID is missing")

type Tagger struct {
	subscriptionID string
}

func NewTagger() (*Tagger, error) {
	sub := os.Getenv("AZURE_SUBSCRIPTION_ID")
	if sub == "" {
		return nil, ErrMissingSubscription
	}
	return &Tagger{subscriptionID: sub}, nil
}

// ApplyTags applies tags to a resourceID (full Azure resource ID).
func (t *Tagger) ApplyTags(ctx context.Context, resourceID string, apiVersion string, tags map[string]string) error {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	client, err := armresources.NewClient(t.subscriptionID, cred, nil)
	if err != nil {
		return err
	}

	azureTags := make(map[string]*string, len(tags))
	for k, v := range tags {
		azureTags[k] = to.Ptr(v)
	}

	poller, err := client.BeginUpdateByID(ctx, resourceID, apiVersion, armresources.GenericResource{
		Tags: azureTags,
	}, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, nil)
	return err
}

//obs new sdk version has different method signature for ApplyTags
