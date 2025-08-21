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

	// Run test suites

	/*suite.RunNamedSuite(t, "Workspace V1", &WorkspaceV1TestSuite{
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
	})*/

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
	})
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
