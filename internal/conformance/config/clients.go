package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/eu-sovereign-cloud/conformance/internal/conformance/params"
	"github.com/eu-sovereign-cloud/conformance/internal/constants"
	"github.com/eu-sovereign-cloud/conformance/internal/mock"
	mockscenarios "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios"
	mockclients "github.com/eu-sovereign-cloud/conformance/internal/mock/scenarios/clients"
	"github.com/eu-sovereign-cloud/go-sdk/secapi"
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
	clientsLock.Lock()
	defer clientsLock.Unlock()

	if Clients != nil {
		return nil
	}

	// Setup mock, if configured to use
	var mockScenario *mockscenarios.Scenario
	if Parameters.MockEnabled {
		mockScenario = mockscenarios.NewScenario(constants.ClientsInitScenarioName,
			&mock.MockParams{
				ServerURL: Parameters.MockServerURL,
				AuthToken: Parameters.ClientAuthToken,
			},
		)
		params := &params.ClientsInitParams{
			Region: Parameters.ClientRegion,
		}
		err := mockclients.ConfigureInitScenarioV1(mockScenario, params)
		if err != nil {
			return fmt.Errorf("failed to configure mock scenario: %w", err)
		}
	}

	var err error
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
	_, isComputeV1Unavailable := Clients.RegionalClient.ComputeV1.(*secapi.ComputeV1Unavailable)
	if !isComputeV1Unavailable {
		Clients.InstanceSkus, err = loadInstanceSkus(ctx, Clients.RegionalClient)
		if err != nil {
			return fmt.Errorf("failed to list instance skus: %w", err)
		}
	}

	// Load available storage skus, if storage provider is available
	_, isStorageV1Unavailable := Clients.RegionalClient.StorageV1.(*secapi.StorageV1Unavailable)
	if !isStorageV1Unavailable {
		Clients.StorageSkus, err = loadStorageSkus(ctx, Clients.RegionalClient)
		if err != nil {
			return fmt.Errorf("failed to list storage skus: %w", err)
		}
	}

	// Load available network skus, if storage network is available
	_, isNetworkV1Unavailable := Clients.RegionalClient.NetworkV1.(*secapi.NetworkV1Unavailable)
	if !isNetworkV1Unavailable {
		Clients.NetworkSkus, err = loadNetworkSkus(ctx, Clients.RegionalClient)
		if err != nil {
			return fmt.Errorf("failed to list network skus: %w", err)
		}
	}

	// Cleanup configured mock scenarios
	if Parameters.MockEnabled {
		if err := mockScenario.ResetScenario(); err != nil {
			return err
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
