package secatest

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/eu-sovereign-cloud/go-sdk/secapi"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

var config *Config

func TestMain(m *testing.M) {
	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	var err error
	config, err = loadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Configure allure
	configureReports(config)

	// Run tests
	os.Exit(m.Run())
}

func TestSuites(t *testing.T) {
	ctx := context.Background()

	setupLogger()

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

	suite.RunNamedSuite(t, "Authorization V1", &AuthorizationV1TestSuite{
		globalTestSuite: globalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: config.mockServerURL,
			},
			client: globalClient,
		},
		users: config.scenarioUsers,
	})

	suite.RunNamedSuite(t, "Workspace V1", &WorkspaceV1TestSuite{
		regionalTestSuite: regionalTestSuite{
			testSuite: testSuite{
				tenant:        config.clientTenant,
				authToken:     config.clientAuthToken,
				mockEnabled:   config.mockEnabled,
				mockServerURL: config.mockServerURL,
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
				mockServerURL: config.mockServerURL,
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
				mockServerURL: config.mockServerURL,
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
				mockServerURL: config.mockServerURL,
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
}

func setupLogger() {
	// TODO Configure handler type and log level via env variables
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)
}

func configureReports(config *Config) {
	resultsPath := config.reportResultsPath

	outputPath := filepath.Dir(resultsPath)
	if err := os.Setenv("ALLURE_OUTPUT_PATH", outputPath); err != nil {
		slog.Error("Failed to configure reports", "error", err)
	}

	outputFolder := filepath.Base(resultsPath)
	if err := os.Setenv("ALLURE_OUTPUT_FOLDER", outputFolder); err != nil {
		slog.Error("Failed to configure reports", "error", err)
	}
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
