package secatest

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/eu-sovereign-cloud/go-sdk/secapi"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestRegionalSuites(t *testing.T) {
	ctx := context.Background()

	// Initialize global client
	globalClient, err := secapi.NewGlobalClient(&secapi.GlobalConfig{
		AuthToken: config.clientAuthToken,
		Endpoints: secapi.GlobalEndpoints{
			RegionV1:        config.providerRegionV1,
			AuthorizationV1: config.providerAuthorizationV1,
		},
	})
	if err != nil {
		slog.Error("Failed to create global client", "error", err)
		os.Exit(1)
	}

	// Initialize regional client
	regionalClient, err := globalClient.NewRegionalClient(ctx, config.clientRegion)
	if err != nil {
		slog.Error("Failed to create regional client", "error", err)
		os.Exit(1)
	}

	// Load region available zones
	regionResp, err := globalClient.RegionV1.GetRegion(ctx, config.clientRegion)
	if err != nil {
		slog.Error("Failed to get region", "error", err)
		os.Exit(1)
	}
	regionZones := regionResp.Spec.AvailableZones

	// Load available instance skus
	instanceSkus, err := loadInstanceSkus(ctx, regionalClient)
	if err != nil {
		slog.Error("Failed to list instance skus", "error", err)
		os.Exit(1)
	}

	// Load available storage skus
	storageSkus, err := loadStorageSkus(ctx, regionalClient)
	if err != nil {
		slog.Error("Failed to list storage skus", "error", err)
		os.Exit(1)
	}

	// Load available network skus
	networkSkus, err := loadNetworkSkus(ctx, regionalClient)
	if err != nil {
		slog.Error("Failed to list network skus", "error", err)
		os.Exit(1)
	}

	// Run test suites

	suite.RunNamedSuite(t, "Workspace V1", &WorkspaceV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			region: config.clientRegion,
			client: regionalClient,
		},
	})

	suite.RunNamedSuite(t, "Storage V1", &StorageV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			region: config.clientRegion,
			client: regionalClient,
		},
		storageSkus: storageSkus,
	})

	suite.RunNamedSuite(t, "Compute V1", &ComputeV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			region: config.clientRegion,
			client: regionalClient,
		},
		availableZones: regionZones,
		instanceSkus:   instanceSkus,
		storageSkus:    storageSkus,
	})

	suite.RunNamedSuite(t, "Network V1", &NetworkV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			region: config.clientRegion,
			client: regionalClient,
		},
		networkCidr:    config.scenarioCidr,
		publicIpsRange: config.scenarioPublicIps,
		regionZones:    regionZones,
		instanceSkus:   instanceSkus,
		storageSkus:    storageSkus,
		networkSkus:    networkSkus,
	})

	suite.RunNamedSuite(t, "Usages V1", &UsagesV1TestSuite{
		mixedTestSuite: mixedTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: &config.mockServerURL,
			},
			globalClient:   globalClient,
			regionalClient: regionalClient,
			region:         config.clientRegion,
		},
		users:          config.scenarioUsers,
		networkCidr:    config.scenarioCidr,
		publicIpsRange: config.scenarioPublicIps,
		regionZones:    regionZones,
		instanceSkus:   instanceSkus,
		storageSkus:    storageSkus,
		networkSkus:    networkSkus,
	})
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
