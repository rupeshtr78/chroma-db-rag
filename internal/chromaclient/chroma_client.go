package chromaclient

import (
	"chroma-db/internal/constants"
	"chroma-db/pkg/logger"
	"context"
	"sync"

	chromago "github.com/amikos-tech/chroma-go"
	chromaAPI "github.com/amikos-tech/chroma-go/swagger"
)

type Metadata map[string]interface{}

type ChromaClient struct {
	Client *chromago.Client
}

var chromaClient *ChromaClient
var once sync.Once
var log = logger.Log

// GetChromaClientInstance returns a singleton instance of the ChromaClient
func GetChromaClientInstance(ctx context.Context, url string, tenantName string, databaseName string) (*ChromaClient, error) {
	var err error
	once.Do(func() {
		chromaClient = &ChromaClient{}
		chromaClient.Client, err = InitializeChroma(ctx, url, tenantName, databaseName)
	})
	return chromaClient, err
}

// InitializeClient initializes the Chroma client and sets the tenant and database.
func InitializeChroma(ctx context.Context, chromaUrl string, tenantName string, databaseName string) (*chromago.Client, error) {
	// Initialize the chroma client
	client, err := GetChromaClient(ctx, constants.ChromaUrl)
	if err != nil {
		log.Debug().Msgf("Error getting chroma client: %v\n", err)
		return nil, err
	}

	// // Get or create the tenant
	_, err = GetOrCreateTenant(ctx, client, constants.TenantName)
	if err != nil {
		log.Debug().Msgf("Error getting or creating tenant: %v\n", err)
		return nil, err
	}

	// Set the tenant for the client
	client.SetTenant(constants.TenantName)

	// // Get or create the database
	_, err = GetOrCreateDatabase(ctx, client, constants.Database, &constants.TenantName)
	if err != nil {
		log.Debug().Msgf("Error getting or creating database: %v\n", err)
		return nil, err
	}

	// Set the database for the client
	client.SetDatabase(constants.Database)

	// client.SetDatabase(constants.Database)
	log.Debug().Msgf("Client Tenant: %v\n", client.Tenant)
	log.Debug().Msgf("Client Database: %v\n", client.Database)

	return client, nil

}

func GetChromaClient(ctx context.Context, url string) (*chromago.Client, error) {
	// Create a new client with url
	chromaClient, err := chromago.NewClient(
		url,
		chromago.WithDebug(false),
		chromago.WithTenant(constants.TenantName),
		chromago.WithDatabase(constants.Database),
	)
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
func GetChromaClientWithOptions(ctx context.Context, url string, tenant string, database string) (*chromago.Client, error) {
	// Create a new client with options
	chromaClient, err := chromago.NewClient(
		url,
		chromago.WithTenant(tenant),
		chromago.WithDatabase(database),
		chromago.WithDebug(false),
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
	// Get the tenant
	tenant, res, err := client.ApiClient.DefaultApi.GetTenant(ctx, tenantName).Execute()
	if err != nil || res.StatusCode != 200 {
		log.Debug().Msgf("Failed to get tenant %v\n", tenantName)
	}
	if tenant != nil && res.StatusCode == 200 {
		log.Debug().Msgf("Tenant %v already exists\n", tenantName)
		return tenant, nil
	}

	t, err := client.CreateTenant(ctx, tenantName)
	if err != nil || t == nil {
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
