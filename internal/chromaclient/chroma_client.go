package chromaclient

import (
	"chroma-db/internal/constants"
	"chroma-db/pkg/logger"
	"context"

	chromago "github.com/amikos-tech/chroma-go"
	chromaAPI "github.com/amikos-tech/chroma-go/swagger"
)

type Metadata map[string]interface{}

var log = logger.Log

func GetChromaClient(ctx context.Context, url string) (*chromago.Client, error) {
	// Create a new client with url
	chromaClient, err := chromago.NewClient(
		url)
	if err != nil {
		return nil, err
	}

	// Confirm that the client can access the server
	if _, errHb := chromaClient.Heartbeat(ctx); errHb != nil {
		return nil, errHb
	}

	return chromaClient, err
}

// GetChromaClient creates a new **chromago.Client** with tenant and database
func GetChromaClientWithOptions(ctx context.Context, url string) (*chromago.Client, error) {
	// Create a new client with options
	chromaClient, err := chromago.NewClient(
		url,
		chromago.WithTenant(constants.TenantName),
		chromago.WithDatabase(constants.Database),
		chromago.WithDebug(true),
		// chromago.WithDefaultHeaders(map[string]string{"Authorization": "Bearer my token"}),
		// chromago.WithSSLCert("path/to/cert.pem"),
	)
	if err != nil {
		return nil, err
	}

	if _, errHb := chromaClient.Heartbeat(ctx); errHb != nil {
		return nil, errHb
	}

	return chromaClient, err
}

// GetOrCreateTenant creates a new **openapi.Tenant** if it does not exist
func GetOrCreateTenant(ctx context.Context, client *chromago.Client, tenantName string) (*chromaAPI.Tenant, error) {

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
	tenantName *string) (*chromaAPI.Database, error) {

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
