package secatest

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/eu-sovereign-cloud/conformance/internal/mock"
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

	// Create Scenario
	err := mock.CreateWorkspaceScenario(mock.Workspace{
		WireMockURL:   config.MockServerURL,
		TenantName:    tenant1Name,
		WorkspaceName: workspace1Name,
		Region:        config.ClientRegion,
		Token:         config.ClientAuthToken,
	})
	if err != nil {
		slog.Error("Failed to create workspace scenario", "error", err)
		os.Exit(1)
	}

	// Initialize global client
	globalClient, err := secapi.NewGlobalClient(&secapi.GlobalConfig{
		AuthToken: config.ClientAuthToken,
		Endpoints: secapi.GlobalEndpoints{
			RegionV1:        config.ProviderRegionV1,
			AuthorizationV1: config.ProviderAuthorizationV1,
		},
	})
	if err != nil {
		slog.Error("Failed to create global client", "error", err)
		os.Exit(1)
	}

	// Initialize regional client
	regionalClient, err := globalClient.NewRegionalClient(ctx, config.ClientRegion)
	if err != nil {
		slog.Error("Failed to create regional client", "error", err)
		os.Exit(1)
	}

	// Run test suites
	suite.RunNamedSuite(t, "Workspace V1", &WorkspaceV1TestSuite{
		client: regionalClient,
	})
}

func configureReports(config *Config) {
	resultsPath := config.ReportResultsPath

	outputPath := filepath.Dir(resultsPath)
	os.Setenv("ALLURE_OUTPUT_PATH", outputPath)

	outputFolder := filepath.Base(resultsPath)
	os.Setenv("ALLURE_OUTPUT_FOLDER", outputFolder)
}
