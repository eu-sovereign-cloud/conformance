package secatest

import (
	"context"
	"fmt"
	"sync"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/wiremock/go-wiremock"
)

type ClientsHolder struct {
	globalClient   *secapi.GlobalClient
	regionalClient *secapi.RegionalClient

	regionZones  []string
	instanceSkus []string
	storageSkus  []string
	networkSkus  []string
}

var (
	clients     *ClientsHolder
	clientsLock sync.Mutex
)

func initClients(ctx context.Context) error {
	var err error

	clientsLock.Lock()
	defer clientsLock.Unlock()

	if clients != nil {
		return nil
	}

	// Setup mock, if configured to use
	var wm *wiremock.Client
	if config.mockEnabled {
		params := mock.ClientsInitParams{
			Params: &mock.Params{
				MockURL:   config.mockServerURL,
				AuthToken: config.clientAuthToken,
				Region:    config.clientRegion,
			},
		}
		wm, err = mock.ConfigClientsInitScenario(&params)
		if err != nil {
			return fmt.Errorf("failed to configure mock scenario: %w", err)
		}
	}

	clients = &ClientsHolder{}

	// Initialize global client
	clients.globalClient, err = secapi.NewGlobalClient(&secapi.GlobalConfig{
		AuthToken: config.clientAuthToken,
		Endpoints: secapi.GlobalEndpoints{
			RegionV1:        config.providerRegionV1,
			AuthorizationV1: config.providerAuthorizationV1,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create global client: %w", err)
	}

	// Initialize regional client
	clients.regionalClient, err = clients.globalClient.NewRegionalClient(ctx, config.clientRegion)
	if err != nil {
		return fmt.Errorf("failed to create regional client: %w", err)
	}

	// Load region available zones
	regionResp, err := clients.globalClient.RegionV1.GetRegion(ctx, config.clientRegion)
	if err != nil {
		return fmt.Errorf("failed to get region: %w", err)
	}
	clients.regionZones = regionResp.Spec.AvailableZones

	// Load available instance skus
	clients.instanceSkus, err = loadInstanceSkus(ctx, clients.regionalClient)
	if err != nil {
		return fmt.Errorf("failed to list instance skus: %w", err)
	}

	// Load available storage skus
	clients.storageSkus, err = loadStorageSkus(ctx, clients.regionalClient)
	if err != nil {
		return fmt.Errorf("failed to list storage skus: %w", err)
	}

	// Load available network skus
	clients.networkSkus, err = loadNetworkSkus(ctx, clients.regionalClient)
	if err != nil {
		return fmt.Errorf("failed to list network skus: %w", err)
	}

	// Cleanup configured mock scenarios
	if config.mockEnabled {
		if err := wm.ResetAllScenarios(); err != nil {
			fmt.Errorf("Failed to reset scenarios: %w", err)
		}
	}

	return nil
}

// TODO Convert these load skus functions to a generic one
func loadInstanceSkus(ctx context.Context, regionalClient *secapi.RegionalClient) ([]string, error) {
	resp, err := regionalClient.ComputeV1.ListSkus(ctx, secapi.TenantID(config.clientTenant))
	if err != nil {
		return nil, err
	}

	skus, err := resp.All(ctx)
	if err != nil {
		return nil, err
	}

	var available []string
	for _, sku := range skus {
		available = append(available, sku.Metadata.Name)
	}

	return available, nil
}

func loadStorageSkus(ctx context.Context, regionalClient *secapi.RegionalClient) ([]string, error) {
	resp, err := regionalClient.StorageV1.ListSkus(ctx, secapi.TenantID(config.clientTenant))
	if err != nil {
		return nil, err
	}

	skus, err := resp.All(ctx)
	if err != nil {
		return nil, err
	}

	var available []string
	for _, sku := range skus {
		available = append(available, sku.Metadata.Name)
	}

	return available, nil
}

func loadNetworkSkus(ctx context.Context, regionalClient *secapi.RegionalClient) ([]string, error) {
	resp, err := regionalClient.NetworkV1.ListSkus(ctx, secapi.TenantID(config.clientTenant))
	if err != nil {
		return nil, err
	}

	skus, err := resp.All(ctx)
	if err != nil {
		return nil, err
	}

	var available []string
	for _, sku := range skus {
		available = append(available, sku.Metadata.Name)
	}

	return available, nil
}
