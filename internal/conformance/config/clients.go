package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockclients "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/clients"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/wiremock/go-wiremock"
)

type ClientsHolder struct {
	GlobalClient   *secapi.GlobalClient
	RegionalClient *secapi.RegionalClient

	RegionZones  []string
	InstanceSkus []string
	StorageSkus  []string
	NetworkSkus  []string
}

var (
	Clients     *ClientsHolder
	clientsLock sync.Mutex
)

func InitClients(ctx context.Context) error {
	var err error

	clientsLock.Lock()
	defer clientsLock.Unlock()

	if Clients != nil {
		return nil
	}

	// Setup mock, if configured to use
	var wm *wiremock.Client
	if Parameters.MockEnabled {
		params := params.ClientsInitParams{
			MockParams: &mock.MockParams{
				ServerURL: Parameters.MockServerURL,
				AuthToken: Parameters.ClientAuthToken,
			},
			Region: Parameters.ClientRegion,
		}
		wm, err = mockclients.ConfigureInitScenarioV1(&params)
		if err != nil {
			return fmt.Errorf("failed to configure mock scenario: %w", err)
		}
	}

	Clients = &ClientsHolder{}

	// Initialize global client
	Clients.GlobalClient, err = secapi.NewGlobalClient(&secapi.GlobalConfig{
		AuthToken: Parameters.ClientAuthToken,
		Endpoints: secapi.GlobalEndpoints{
			RegionV1:        Parameters.ProviderRegionV1,
			AuthorizationV1: Parameters.ProviderAuthorizationV1,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create global client: %w", err)
	}

	// Initialize regional client
	Clients.RegionalClient, err = Clients.GlobalClient.NewRegionalClient(ctx, Parameters.ClientRegion)
	if err != nil {
		return fmt.Errorf("failed to create regional client: %w", err)
	}

	// Load region available zones
	regionResp, err := Clients.GlobalClient.RegionV1.GetRegion(ctx, Parameters.ClientRegion)
	if err != nil {
		return fmt.Errorf("failed to get region: %w", err)
	}
	Clients.RegionZones = regionResp.Spec.AvailableZones

	// Load available instance skus, if compute provider is available
	if Clients.RegionalClient.ComputeV1 != nil {
		Clients.InstanceSkus, err = loadInstanceSkus(ctx, Clients.RegionalClient)
		if err != nil {
			return fmt.Errorf("failed to list instance skus: %w", err)
		}
	}

	// Load available storage skus, if storage provider is available
	if Clients.RegionalClient.StorageV1 != nil {
		Clients.StorageSkus, err = loadStorageSkus(ctx, Clients.RegionalClient)
		if err != nil {
			return fmt.Errorf("failed to list storage skus: %w", err)
		}
	}

	// Load available network skus, if storage network is available
	if Clients.RegionalClient.NetworkV1 != nil {
		Clients.NetworkSkus, err = loadNetworkSkus(ctx, Clients.RegionalClient)
		if err != nil {
			return fmt.Errorf("failed to list network skus: %w", err)
		}
	}

	// Cleanup configured mock scenarios
	if Parameters.MockEnabled {
		if err := wm.ResetAllScenarios(); err != nil {
			return fmt.Errorf("failed to reset scenarios: %w", err)
		}
	}

	return nil
}

// TODO Convert these load skus functions to a generic one
func loadInstanceSkus(ctx context.Context, regionalClient *secapi.RegionalClient) ([]string, error) {
	resp, err := regionalClient.ComputeV1.ListSkus(ctx, secapi.TenantID(Parameters.ClientTenant))
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
	resp, err := regionalClient.StorageV1.ListSkus(ctx, secapi.TenantID(Parameters.ClientTenant))
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
	resp, err := regionalClient.NetworkV1.ListSkus(ctx, secapi.TenantID(Parameters.ClientTenant))
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
