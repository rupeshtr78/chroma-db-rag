package chromaclient

import (
	"chroma-db/pkg/logger"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	openapiclient "github.com/amikos-tech/chroma-go/swagger"
)

type Metadata map[string]interface{}

var log = logger.Log

// GetChromaClient creates a new **chromago.Client** and confirms that it can access the server
func GetChromaClient(ctx context.Context, url string) (*chromago.Client, error) {
	// create the client connection and confirm that we can access the server with it
	chromaClient, err := chromago.NewClient(url)
	if err != nil {
		return nil, err
	}

	if _, errHb := chromaClient.Heartbeat(ctx); errHb != nil {
		return nil, errHb
	}

	return chromaClient, err
}

// GetOrCreateTenant creates a new **openapi.Tenant** if it does not exist
func GetOrCreateTenant(ctx context.Context, client *chromago.Client, tenantName string) (*openapiclient.Tenant, error) {

	if t, err := client.GetTenant(ctx, tenantName); err == nil {
		log.Debug().Msgf("Tenant %v already exists\n", tenantName)
		return t, nil
	}

	t, err := client.CreateTenant(ctx, tenantName)
	if err != nil {
		log.Debug().Msgf("Failed to create tenant %v\n", tenantName)
		return nil, err
	}
	return t, nil
}

func GetOrCreateDatabase(ctx context.Context,
	client *chromago.Client,
	dbName string,
	tenantName *string) (*openapiclient.Database, error) {

	if d, err := client.GetDatabase(ctx, dbName, tenantName); err == nil {
		log.Debug().Msgf("Database %v already exists\n", dbName)
		return d, nil
	}

	d, err := client.CreateDatabase(ctx, dbName, tenantName)
	if err != nil {
		log.Debug().Msgf("Failed to create database %v\n", dbName)
		return nil, err
	}
	return d, nil
}
